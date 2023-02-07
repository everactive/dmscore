package devicetwin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/benbjohnson/clock"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/everactive/dmscore/config/keys"
	devicetwinconfig "github.com/everactive/dmscore/iot-devicetwin/config"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/actions"
	"github.com/everactive/dmscore/iot-devicetwin/service/controller"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
	"github.com/everactive/dmscore/iot-devicetwin/web"
	devicetwinweb "github.com/everactive/dmscore/iot-devicetwin/web"
	"github.com/everactive/dmscore/iot-identity/service/cert"
	identityweb "github.com/everactive/dmscore/iot-identity/web"
	"github.com/everactive/dmscore/models"
	datastore2 "github.com/everactive/dmscore/pkg/datastores"
	"github.com/everactive/dmscore/pkg/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/thejerf/suture/v4"
	"gorm.io/gorm"
	"path"
	"strings"
	"time"
)

const (
	clientIDMQTTTopicPartsCount = 4
	defaultChannelBufferSize    = 10
	serviceName                 = "DeviceTwin"
)

var (
	FS                         = afero.NewOsFs()
	AFS                        = &afero.Afero{Fs: FS}
	serviceLogger              = log.StandardLogger()
	logger                     = serviceLogger.WithFields(map[string]interface{}{"Service": serviceName})
	CreateDeviceTwinWebService = createDeviceTwinWebService
	rtClock                    = clock.New()
)

func createDeviceTwinWebService(port string, ctrl controller.Controller) *devicetwinweb.Service {
	return devicetwinweb.NewService(port, ctrl)
}

type Service struct {
	deviceTwinWebService *web.Service
	MQTT                 mqtt.Connect
	legacyHealthChan     chan MQTT.Message
	legacyActionChan     chan MQTT.Message
	legacyPublishChan    chan mqtt.PublishMessage
	healthChan           chan MQTT.Message
	actionChan           chan MQTT.Message
	publishChan          chan mqtt.PublishMessage
	db                   *gorm.DB
	controller           controller.Controller
	twin                 devicetwin.DeviceTwin
	identity             *identityweb.IdentityService
	datastore            datastore2.ConsolidatedDataStore
}

var GetMQTTConnection = getMQTTConnection

func getMQTTConnection(url, port string, connect *devicetwinconfig.MQTTConnect, service *Service) (*mqtt.Connection, error) {
	return mqtt.GetConnection(url, port, connect, func(c MQTT.Client) {
		logger.Println("Connecting to MQTT, subscribing to actions")
		if err := service.SubscribeToActions(); err != nil {
			logger.Fatalf("error establishing MQTT subscriptions: %s", err)
		}
	})
}

var newSuperVisor = func(name string) *suture.Supervisor { return suture.NewSimple(name) }

func New(identity *identityweb.IdentityService, dss *datastore2.DataStores) (*suture.Supervisor, controller.Controller) {
	sup := newSuperVisor("devicetwin") //suture.NewSimple("devicetwin")

	URL := viper.GetString(keys.MQTTURL)
	port := viper.GetString(keys.MQTTPort)

	certsDir := viper.GetString(keys.MQTTCertificatesPath)
	logger.Tracef("MQTT Certs dir: %s", certsDir)

	rootCAFilename := viper.GetString(keys.MQTTRootCAFilename)
	clientCertFilename := viper.GetString(keys.MQTTClientCertificateFilename)
	clientKeyFilename := viper.GetString(keys.MQTTClientKeyFilename)

	c := devicetwinconfig.MQTTConnect{}

	// nolint: gosec
	rootCA, err := AFS.ReadFile(path.Join(certsDir, rootCAFilename))
	if err != nil {
		panic(err)
	}

	// nolint: gosec
	certFile, err := AFS.ReadFile(path.Join(certsDir, clientCertFilename))
	if err != nil {
		panic(err)
	}

	// nolint: gosec
	key, err := AFS.ReadFile(path.Join(certsDir, clientKeyFilename))

	c.RootCA = rootCA
	c.ClientKey = key
	c.ClientCert = certFile

	prefix := viper.GetString(keys.MQTTClientIDPrefix)

	// Generate a random string
	s, err := cert.CreateSecret(6)
	if err != nil {
		logger.Printf("Error creating client ID: %v", err)
	}

	c.ClientID = fmt.Sprintf("%s-%s", prefix, s)

	legacyHealthChan := make(chan MQTT.Message, defaultChannelBufferSize)
	legacyActionChan := make(chan MQTT.Message, defaultChannelBufferSize)
	legacyPublishChan := make(chan mqtt.PublishMessage, defaultChannelBufferSize)

	twin := devicetwin.NewService(dss.DeviceTwinStore, dss.ManagementStore)
	ctrl := controller.NewService(legacyHealthChan, legacyActionChan, legacyPublishChan, twin)

	servicePort := viper.GetString(keys.GetDeviceTwinKey(keys.ServicePort))

	w := CreateDeviceTwinWebService(servicePort, ctrl)

	service := &Service{
		deviceTwinWebService: w,
		healthChan:           make(chan MQTT.Message, defaultChannelBufferSize),
		publishChan:          make(chan mqtt.PublishMessage, defaultChannelBufferSize),
		actionChan:           make(chan MQTT.Message, defaultChannelBufferSize),
		db:                   dss.GetDatabase(),
		twin:                 twin,
		controller:           ctrl,
		identity:             identity,
	}

	// Set up the MQTT client and handle pub/sub from here... as the MQTT and DeviceTwin services are mutually dependent
	// The onconnect handler will subscribe on both new connection and reconnection
	m, err := GetMQTTConnection(URL, port, &c, service)

	if err != nil {
		logger.Fatalf("Error connecting to MQTT broker: %v", err)
	}

	service.MQTT = m

	service.legacyHealthChan = legacyHealthChan
	service.legacyActionChan = legacyActionChan
	service.legacyPublishChan = legacyPublishChan

	installService := NewInstallService(dss, legacyPublishChan)

	sup.Add(service)
	sup.Add(ctrl)
	sup.Add(installService)

	return sup, w.Controller
}

func (srv *Service) String() string {
	return serviceName
}

func (srv *Service) Serve(ctx context.Context) error {
	intervalTicker := rtClock.Ticker(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			logger.Errorf("We're done: %s", ctx.Err())
			return nil
		case <-intervalTicker.C:
			logger.Infof("%s still ticking", "DeviceTwinService")
		case m := <-srv.legacyPublishChan:
			err := srv.MQTT.Publish(m.Topic, m.Payload)
			if err != nil {
				logger.Error(err)
			}
		case m := <-srv.healthChan:
			logger.Infof("Received health message: %s", string(m.Payload()))
			err := srv.healthMessageHandler(m)
			if err != nil {
				logger.Error(err)
			}
		case m := <-srv.actionChan:
			logger.Infof("Received action message: %s", string(m.Payload()))
			err := srv.actionMessageHandler(m)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func (srv *Service) sendLegacyActionHandlerUnversionedPayload(clientID string, msg MQTT.Message, versionedMessage messages.VersionedMessage) error {
	var payload []byte
	switch versionedMessage.Action {
	case actions.List:
		var versionedPublishSnaps messages.PublishSnapsV2
		err := json.Unmarshal(msg.Payload(), &versionedPublishSnaps)
		if err != nil {
			return fmt.Errorf("error trying to unmarshal PublishSnapsV2 message: %w", err)
		}
		// If there was no error then construct the unversioned representative message
		publishSnaps := messages.PublishSnaps{
			Action:  versionedPublishSnaps.Action,
			Id:      versionedPublishSnaps.Id,
			Message: versionedPublishSnaps.Message,
			Success: versionedMessage.Success,
		}

		if versionedPublishSnaps.Result != nil {
			publishSnaps.Result = versionedPublishSnaps.Result.Snaps
		}

		payload, err = json.Marshal(&publishSnaps)
		if err != nil {
			return fmt.Errorf("error trying to unmarshal PublishSnaps generated message: %w", err)
		}
	}

	// If we got here, double-check the payload size and if it's non-zero, send it
	if len(payload) > 0 {
		err := srv.twin.ActionResponse(clientID, versionedMessage.Id, versionedMessage.Action, payload)
		if err != nil {
			return fmt.Errorf("error in ActionResponse: %w", err)
		}
	}

	return nil
}

func (srv *Service) actionMessageHandler(msg MQTT.Message) error {
	clientID := getClientID(msg)
	logger.Printf("Action response from %s", clientID)

	// is this a versioned message?
	var versionedMessage messages.VersionedMessage
	if err := json.Unmarshal(msg.Payload(), &versionedMessage); err != nil {
		return fmt.Errorf("error in action message: %w", err)
	}

	// Check if there is an error and handle it
	if !versionedMessage.Success {
		return fmt.Errorf("error in action `%s`: (%s) %s", versionedMessage.Action, versionedMessage.Id, versionedMessage.Message)
	}

	logger.Infof("Received action message type %s, sending to iot-devicetwin", versionedMessage.Action)

	// If it's unversioned it will be the string zero value (empty), if it's a version we know, handle it, otherwise error
	switch versionedMessage.Version {
	case "":
		err := srv.twin.ActionResponse(clientID, versionedMessage.Id, versionedMessage.Action, msg.Payload())
		if err != nil {
			return fmt.Errorf("error in ActionResponse: %w", err)
		}
	case "2":
		// we should only be expecting the list action
		switch versionedMessage.Action {
		case actions.List:
			// handle it here
			var versionedPublishSnaps messages.PublishSnapsV2
			err := json.Unmarshal(msg.Payload(), &versionedPublishSnaps)
			if err != nil {
				return fmt.Errorf("error trying to unmarshal PublishSnapsV2 message: %w", err)
			}
			// update hashes
			var healthHashes models.HealthHash

			tx := srv.db.Find(&healthHashes, &models.HealthHash{DeviceID: clientID})
			if tx.Error != nil {
				return fmt.Errorf("error trying to find health hash for %s: %w", versionedPublishSnaps.Id, tx.Error)
			}

			if tx.RowsAffected == 0 {
				// If we received a list of snaps for device and we don't have a health has
				// entry for it yet, ignore it for now.
				logger.Infof("Received snap list for clientID=%s, actionID=%s but do not have a health hash entry for it yet; will still process it", clientID, versionedPublishSnaps.Id)
				// build a message payload for ActionResponse
				err = srv.sendLegacyActionHandlerUnversionedPayload(clientID, msg, versionedMessage)
				if err != nil {
					return fmt.Errorf("error trying to send unversioned payload to legacy action handler: %w", err)
				}
				return nil
			}

			// OK- we have our health hash entry, update it
			tx = srv.db.Model(&healthHashes).Updates(&models.HealthHash{
				DeviceID:           healthHashes.DeviceID,
				SnapListHash:       versionedPublishSnaps.Result.SnapListHash,
				InstalledSnapsHash: versionedPublishSnaps.Result.InstalledSnapsHash,
				LastRefresh:        time.Now(),
			})

			if tx.Error != nil {
				return fmt.Errorf("error trying to update health hashes for %s: %w", healthHashes.DeviceID, tx.Error)
			}

			// build a message payload for ActionResponse
			err = srv.sendLegacyActionHandlerUnversionedPayload(clientID, msg, versionedMessage)
			if err != nil {
				return fmt.Errorf("error trying to send unversioned payload to legacy action handler: %w", err)
			}
			return nil
		default:
			logger.Errorf("Received a versioned action message we cannot handled: %s", versionedMessage.Action)
			return fmt.Errorf("error trying to handle versioned action message type %s", versionedMessage.Action)
		}
	default:
		err := fmt.Errorf("received a versioned action message with a version we don't expect: %s", versionedMessage.Version)
		logger.Error(err)
		return err

	}

	return nil
}

func (srv *Service) healthMessageHandler(msg MQTT.Message) error {
	// After sending the message to the legacy handler, do our own processing
	clientID := getClientID(msg)
	logger.Printf("Health update from %s", clientID)

	// Parse the body
	h := messages.Health{}
	if err := json.Unmarshal(msg.Payload(), &h); err != nil {
		return fmt.Errorf("error trying to unmarshal health message payload: %w", err)
	}

	// Check that the client ID matches
	if clientID != h.DeviceId {
		return fmt.Errorf("client/device ID mismatch: %s from the topic and %s from the message; discarding", clientID, h.DeviceId)
	}

	logger.Infof("Received health message %+v, sending to iot-devicetwin", msg)
	srv.legacyHealthChan <- msg

	var healthMessage messages.Health
	err := json.Unmarshal(msg.Payload(), &healthMessage)
	if err != nil {
		return fmt.Errorf("failed unmarshaling a health message: %w", err)
	}

	if healthMessage.InstalledSnapsHash == "" && healthMessage.SnapListHash == "" {
		// No error, could be just a device that does not support this yet
		logger.Infof("Received health message from %s but did not have hashes", healthMessage.DeviceId)
		return nil
	}

	var healthHashes models.HealthHash
	tx := srv.db.Find(&healthHashes, &models.HealthHash{DeviceID: healthMessage.DeviceId})
	if tx.Error != nil {
		return fmt.Errorf("error trying to find health hash for %s: %w", healthMessage.DeviceId, tx.Error)
	}

	if tx.RowsAffected == 0 {
		srv.db.Create(&models.HealthHash{
			LastRefresh:        time.Now(),
			OrgID:              healthMessage.OrgId,
			DeviceID:           healthMessage.DeviceId,
			SnapListHash:       healthMessage.SnapListHash,
			InstalledSnapsHash: healthMessage.InstalledSnapsHash,
		})

		return nil
	}

	refreshOnAnyChanges := viper.GetBool(keys.RefreshSnapListOnAnyChange)

	if healthHashes.SnapListHash == healthMessage.SnapListHash &&
		(healthHashes.InstalledSnapsHash == healthMessage.InstalledSnapsHash || (healthHashes.InstalledSnapsHash != healthMessage.InstalledSnapsHash && !refreshOnAnyChanges)) {
		// Just update and return
		srv.db.Model(&healthHashes).Updates(&models.HealthHash{LastRefresh: time.Now()})
		log.Infof("Update health hash last refresh time for %s to now. SnapListHash=%s and InstallSnapsHash=%s",
			healthHashes.DeviceID, healthHashes.SnapListHash, healthHashes.InstalledSnapsHash)
		return nil
	}

	// if they don't match, then we need to request the updated list
	srv.controller.DeviceSnapList(healthMessage.OrgId, healthMessage.DeviceId)

	return nil
}

// SubscribeToActions subscribes to the published topics from the devices
func (srv *Service) SubscribeToActions() error {
	healthTopic := viper.GetString(keys.MQTTHealthTopic)
	pubTopic := viper.GetString(keys.MQTTPubTopic)

	// Subscribe to the device health messages
	if err := srv.MQTT.Subscribe(healthTopic, srv.healthChannelForwarder); err != nil {
		logger.Printf("Error subscribing to topic `%s`: %v", healthTopic, err)
		return err
	}

	// Subscribe to the device action responses
	if err := srv.MQTT.Subscribe(pubTopic, srv.actionChannelForwarder); err != nil {
		logger.Printf("Error subscribing to topic `%s`: %v", pubTopic, err)
		return err
	}

	return nil
}

// ActionHandler is the handler for the main subscription topic
func (srv *Service) actionChannelForwarder(_ MQTT.Client, msg MQTT.Message) {
	srv.actionChan <- msg
}

func (srv *Service) healthChannelForwarder(_ MQTT.Client, msg MQTT.Message) {
	srv.healthChan <- msg
}

// getClientID sets the client ID from the topic
func getClientID(msg MQTT.Message) string {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) != clientIDMQTTTopicPartsCount-1 {
		logger.Printf("Error in health message: expected %d parts, got %d", clientIDMQTTTopicPartsCount-1, len(parts))
		return ""
	}
	return parts[2]
}
