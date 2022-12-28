package devicetwin

import (
	"context"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/everactive/dmscore/config/keys"
	devicetwinconfig "github.com/everactive/dmscore/iot-devicetwin/config"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	"github.com/everactive/dmscore/iot-devicetwin/service/controller"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	devicetwinfactory "github.com/everactive/dmscore/iot-devicetwin/service/factory"
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
	"github.com/everactive/dmscore/iot-devicetwin/web"
	devicetwinweb "github.com/everactive/dmscore/iot-devicetwin/web"
	"github.com/everactive/dmscore/iot-identity/service/cert"
	"github.com/everactive/dmscore/iot-management/datastore"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thejerf/suture/v4"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

const (
	clientIDMQTTTopicPartsCount = 4
)

type Service struct {
	deviceTwinWebService *web.Service
	MQTT       mqtt.Connect
	legacyHealthChan chan MQTT.Message
	legacyActionChan chan MQTT.Message
	legacyPublishChan chan mqtt.PublishMessage
}

func New(coreDB datastore.DataStore) (*suture.Supervisor, controller.Controller) {
	sup := suture.NewSimple("devicetwin")

	databaseDriver := viper.GetString(keys.GetDeviceTwinKey(keys.DatabaseDriver))
	dataStoreSource := viper.GetString(keys.GetDeviceTwinKey(keys.DatabaseConnectionString))

	db, err := devicetwinfactory.CreateDataStore(databaseDriver, dataStoreSource)
	if err != nil {
		log.Fatalf("Error connecting to data store: %v", err)
	}

	URL := viper.GetString(keys.MQTTURL)
	port := viper.GetString(keys.MQTTPort)

	certsDir := viper.GetString(keys.MQTTCertificatesPath)
	log.Tracef("MQTT Certs dir: %s", certsDir)

	rootCAFilename := viper.GetString(keys.MQTTRootCAFilename)
	clientCertFilename := viper.GetString(keys.MQTTClientCertificateFilename)
	clientKeyFilename := viper.GetString(keys.MQTTClientKeyFilename)

	c := devicetwinconfig.MQTTConnect{}

	// nolint: gosec
	rootCA, err := ioutil.ReadFile(path.Join(certsDir, rootCAFilename))
	if err != nil {
		panic(err)
	}

	// nolint: gosec
	certFile, err := ioutil.ReadFile(path.Join(certsDir, clientCertFilename))
	if err != nil {
		panic(err)
	}

	// nolint: gosec
	key, err := ioutil.ReadFile(path.Join(certsDir, clientKeyFilename))

	c.RootCA = rootCA
	c.ClientKey = key
	c.ClientCert = certFile

	prefix := viper.GetString(keys.MQTTClientIDPrefix)

	// Generate a random string
	s, err := cert.CreateSecret(6)
	if err != nil {
		log.Printf("Error creating client ID: %v", err)
	}

	c.ClientID = fmt.Sprintf("%s-%s", prefix, s)

	legacyHealthChan := make(chan MQTT.Message, 10)
	legacyActionChan := make(chan MQTT.Message, 10)
	legacyPublishChan := make(chan mqtt.PublishMessage, 10)

	twin := devicetwin.NewService(db, coreDB)
	ctrl := controller.NewService(legacyHealthChan, legacyActionChan, legacyPublishChan, twin)

	servicePort := viper.GetString(keys.GetDeviceTwinKey(keys.ServicePort))

	// Start the web API service
	w := devicetwinweb.NewService(servicePort, ctrl)

	service := &Service{deviceTwinWebService: w}

	// Setup the MQTT client and handle pub/sub from here... as the MQTT and DeviceTwin services are mutually dependent
	// The onconnect handler will subscribe on both new connection and reconnection
	m, err := mqtt.GetConnection(URL, port, &c, func(c MQTT.Client) {
		log.Println("Connecting to MQTT, subscribing to actions")
		if err := service.SubscribeToActions(); err != nil {
			log.Fatalf("error establishing MQTT subscriptions: %s", err)
		}
	})

	if err != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", err)
	}

	service.MQTT = m

	service.legacyHealthChan = legacyHealthChan
	service.legacyActionChan = legacyActionChan
	service.legacyPublishChan = legacyPublishChan

	sup.Add(service)
	sup.Add(ctrl)

	return sup, w.Controller
}

func (srv *Service) Serve(ctx context.Context) error {
	intervalTicker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Errorf("We're done: %s", ctx.Err())
			return nil
		case <-intervalTicker.C:
			log.Infof("%s still ticking", "DeviceTwinService")
		case m := <- srv.legacyPublishChan:
			srv.MQTT.Publish(m.Topic, m.Payload)
		}
	}

	return nil
}

// SubscribeToActions subscribes to the published topics from the devices
func (srv *Service) SubscribeToActions() error {
	healthTopic := viper.GetString(keys.MQTTHealthTopic)
	pubTopic := viper.GetString(keys.MQTTPubTopic)

	// Subscribe to the device health messages
	if err := srv.MQTT.Subscribe(healthTopic, srv.HealthHandler); err != nil {
		log.Printf("Error subscribing to topic `%s`: %v", healthTopic, err)
		return err
	}

	// Subscribe to the device action responses
	if err := srv.MQTT.Subscribe(pubTopic, srv.ActionHandler); err != nil {
		log.Printf("Error subscribing to topic `%s`: %v", pubTopic, err)
		return err
	}

	return nil
}

// ActionHandler is the handler for the main subscription topic
func (srv *Service) ActionHandler(_ MQTT.Client, msg MQTT.Message) {
	clientID := getClientID(msg)
	log.Printf("Action response from %s", clientID)

	// Parse the body
	a := messages.PublishResponse{}
	if err := json.Unmarshal(msg.Payload(), &a); err != nil {
		log.Printf("Error in action message: %v", err)
		return
	}

	// Check if there is an error and handle it
	if !a.Success {
		log.Printf("Error in action `%s`: (%s) %s", a.Action, a.Id, a.Message)
		return
	}

	log.Infof("Received action message %+v, sending to iot-devicetwin", msg)
	srv.legacyActionChan <- msg
}

// HealthHandler is the handler for the devices health messages
func (srv *Service) HealthHandler(client MQTT.Client, msg MQTT.Message) {
	clientID := getClientID(msg)
	log.Printf("Health update from %s", clientID)

	// Parse the body
	h := messages.Health{}
	if err := json.Unmarshal(msg.Payload(), &h); err != nil {
		log.Printf("Error in health message: %v", err)
		return
	}

	// Check that the client ID matches
	if clientID != h.DeviceId {
		log.Printf("Client/device ID mismatch: %s and %s", clientID, h.DeviceId)
		return
	}

	log.Infof("Received health message %+v, sending to iot-devicetwin", msg)
	srv.legacyHealthChan <- msg
}

// getClientID sets the client ID from the topic
func getClientID(msg MQTT.Message) string {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) != clientIDMQTTTopicPartsCount-1 {
		log.Printf("Error in health message: expected 4 parts, got %d", len(parts))
		return ""
	}
	return parts[2]
}
