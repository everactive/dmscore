package devicetwin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/everactive/dmscore/config/keys"
	"github.com/everactive/dmscore/iot-devicetwin/config"
	devicetwindatastore "github.com/everactive/dmscore/iot-devicetwin/datastore"
	"github.com/everactive/dmscore/iot-devicetwin/service/controller"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	"github.com/everactive/dmscore/iot-devicetwin/service/factory"
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
	devicetwinweb "github.com/everactive/dmscore/iot-devicetwin/web"
	"github.com/everactive/dmscore/iot-management/datastore"
	mocks "github.com/everactive/dmscore/mocks/external/mqtt"
	"github.com/everactive/dmscore/models"
	"github.com/everactive/dmscore/pkg/datastores"
	"github.com/everactive/dmscore/pkg/messages"
	migrate2 "github.com/everactive/dmscore/pkg/migrate"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thejerf/suture/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		coreDB datastore.DataStore
		db     *gorm.DB
	}
	tests := []struct {
		name           string
		args           args
		wantSupervisor *suture.Supervisor
		wantController controller.Controller
	}{
		{
			name:           "valid",
			args:           args{coreDB: &datastore.MockDataStore{}},
			wantSupervisor: suture.NewSimple("devicetwin"),
			wantController: &controller.MockController{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			factory.CreateDataStore = testDeviceTwinDatastore(&devicetwindatastore.MockDataStore{}, nil)

			fs := afero.NewMemMapFs()
			AFS = &afero.Afero{Fs: fs}

			rootCAFilename := "ca.crt"
			clientCertificateFilename := "client.crt"
			clientKeyFileName := "client.key"
			viper.Set(keys.MQTTRootCAFilename, rootCAFilename)
			viper.Set(keys.MQTTClientCertificateFilename, clientCertificateFilename)
			viper.Set(keys.MQTTClientKeyFilename, clientKeyFileName)

			_, err := fs.Create(rootCAFilename)
			if err != nil {
				t.Error(err)
			}
			_, err = fs.Create(clientCertificateFilename)
			if err != nil {
				t.Error(err)
			}
			_, err = fs.Create(clientKeyFileName)
			if err != nil {
				t.Error(err)
			}

			GetMQTTConnection = func(url, port string, connect *config.MQTTConnect, service *Service) (*mqtt.Connection, error) {
				return nil, nil
			}

			newSuperVisor = func(name string) *suture.Supervisor { return tt.wantSupervisor }

			expectedPort := "1000"
			viper.Set(keys.GetDeviceTwinKey(keys.ServicePort), expectedPort)
			expectedController := tt.wantController
			service := devicetwinweb.NewService(expectedPort, expectedController)
			tt.wantController = expectedController
			CreateDeviceTwinWebService = func(port string, ctrl controller.Controller) *devicetwinweb.Service {
				return service
			}

			got, got1 := New(nil, &datastores.DataStores{
				IdentityStore:   nil,
				DeviceTwinStore: nil,
				ManagementStore: nil,
				DataStore:       nil,
			})
			if !reflect.DeepEqual(got, tt.wantSupervisor) {
				t.Errorf("New() got = %v, want %v", got, tt.wantSupervisor)
			}
			if !reflect.DeepEqual(got1, tt.wantController) {
				t.Errorf("New() got1 = %v, want %v", got1, tt.wantController)
			}
		})
	}
}

func TestService_Serve(t *testing.T) {
	type fields struct{}
	type args struct {
		ctx                   context.Context
		expectedTopic         string
		expectedPayload       string
		expectedHealthMessage *messages.Health
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "valid done",
			fields:  fields{},
			args:    args{},
			wantErr: false,
		},
		{
			name:   "legacy channel publish, no err",
			fields: fields{},
			args: args{
				expectedTopic:   "device/something/or/another",
				expectedPayload: "this is the payload",
			},
			wantErr: false,
		},
		{
			name:   "health channel receive, no err",
			fields: fields{},
			args: args{
				expectedHealthMessage: &messages.Health{
					DeviceId:           "1",
					InstalledSnapsHash: "",
					OrgId:              "12",
					Refresh:            rtClock.Now().String(),
					SnapListHash:       "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{
				legacyPublishChan: make(chan mqtt.PublishMessage),
				healthChan:        make(chan MQTT.Message),
				legacyHealthChan:  make(chan MQTT.Message),
			}

			ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)

			mockedMQTTConnect := mqtt.MockConnect{}
			mockedMessage := mocks.Message{}
			if tt.args.expectedPayload != "" && tt.args.expectedTopic != "" {
				mockedMQTTConnect.On("Publish", tt.args.expectedTopic, tt.args.expectedPayload).Return(nil)
				ctx, cancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
			}

			srv.MQTT = &mockedMQTTConnect

			if tt.args.expectedHealthMessage != nil {
				bytes, err := json.Marshal(tt.args.expectedHealthMessage)
				assert.Nil(t, err)
				mockedMessage.On("Topic").Return(fmt.Sprintf("some/topic/%s", tt.args.expectedHealthMessage.DeviceId))
				mockedMessage.On("Payload").Return(bytes)
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				wg.Done()
				if err := srv.Serve(ctx); (err != nil) != tt.wantErr {
					t.Errorf("Serve() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			if tt.args.expectedPayload != "" && tt.args.expectedTopic != "" {
				srv.legacyPublishChan <- mqtt.PublishMessage{
					Topic:   tt.args.expectedTopic,
					Payload: tt.args.expectedPayload,
				}
				mockedMQTTConnect.AssertExpectations(t)
				cancelFunc()
			} else if tt.args.expectedHealthMessage != nil {
				srv.healthChan <- &mockedMessage
			}

			wg.Wait()
			cancelFunc()

			assert.Equal(t, context.Canceled, ctx.Err())
		})
	}
}

func TestService_String(t *testing.T) {
	type fields struct{}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "valid",
			fields: fields{},
			want:   serviceName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{}
			if got := srv.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SubscribeToActions(t *testing.T) {
	type fields struct {
		expectedHealthTopic  string
		expectedPublishTopic string
	}
	type args struct {
		expectedHealthTopicReturn  error
		expectedPublishTopicReturn error
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				expectedHealthTopic:  "health-topic",
				expectedPublishTopic: "publish-topic",
			},
		},
		{
			name: "error subscribing to publish topic",
			fields: fields{
				expectedHealthTopic:  "health-topic",
				expectedPublishTopic: "publish-topic",
			},
			args: args{
				expectedPublishTopicReturn: errors.New("this is an error"),
			},
			wantErr: true,
		},
		{
			name: "error subscribing to health topic",
			fields: fields{
				expectedHealthTopic:  "health-topic",
				expectedPublishTopic: "publish-topic",
			},
			args: args{
				expectedHealthTopicReturn: errors.New("this is an error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			viper.Set(keys.MQTTHealthTopic, tt.fields.expectedHealthTopic)
			viper.Set(keys.MQTTPubTopic, tt.fields.expectedPublishTopic)

			srv := &Service{}

			mqttMock := &mqtt.MockConnect{}
			mqttMock.On("Subscribe", tt.fields.expectedHealthTopic, mock.AnythingOfType("mqtt.MessageHandler")).Return(tt.args.expectedHealthTopicReturn)
			mqttMock.On("Subscribe", tt.fields.expectedPublishTopic, mock.AnythingOfType("mqtt.MessageHandler")).Return(tt.args.expectedPublishTopicReturn)

			srv.MQTT = mqttMock

			if err := srv.SubscribeToActions(); (err != nil) != tt.wantErr {
				t.Errorf("SubscribeToActions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_actionChannelForwarder(t *testing.T) {
	type fields struct {
		actionChan chan MQTT.Message
	}
	type args struct {
		mqttClient      MQTT.Client
		expectedMessage string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "valid",
			fields: fields{},
			args:   args{mqttClient: &mocks.Client{}, expectedMessage: "this is the expected message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{
				actionChan: make(chan MQTT.Message),
			}
			mockMessage := mocks.Message{}
			mockMessage.On("Payload").Return([]byte(tt.args.expectedMessage))

			go func() {
				srv.actionChannelForwarder(tt.args.mqttClient, &mockMessage)
			}()

			msg := <-srv.actionChan

			assert.Equal(t, tt.args.expectedMessage, string(msg.Payload()))
		})
	}
}

func TestService_healthChannelForwarder(t *testing.T) {
	type fields struct {
		healthChan chan MQTT.Message
	}
	type args struct {
		mqttClient      MQTT.Client
		expectedMessage string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "valid",
			fields: fields{},
			args:   args{mqttClient: &mocks.Client{}, expectedMessage: "this is the expected message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{
				healthChan: make(chan MQTT.Message),
			}
			mockMessage := mocks.Message{}
			mockMessage.On("Payload").Return([]byte(tt.args.expectedMessage))

			go func() {
				srv.healthChannelForwarder(tt.args.mqttClient, &mockMessage)
			}()

			msg := <-srv.healthChan

			assert.Equal(t, tt.args.expectedMessage, string(msg.Payload()))
		})
	}
}

func TestService_actionMessageHandler(t *testing.T) {
	type fields struct {
	}
	type args struct {
		expectedTopic               string
		expectedPublishSnapsMessage messages.PublishSnapsV2
		expectedDeviceID            string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid-versioned",
			args: args{
				expectedTopic: "device/health/some-device-id", expectedPublishSnapsMessage: messages.PublishSnapsV2{Version: "2", Action: "list", Success: true, Id: "1029384756"}, expectedDeviceID: "some-device-id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{}

			mockMessage := &mocks.Message{}
			mockMessage.On("Topic").Return(tt.args.expectedTopic)
			bytes, err := json.Marshal(&tt.args.expectedPublishSnapsMessage)
			if err != nil {
				assert.Error(t, err)
			}
			mockMessage.On("Payload").Return(bytes)

			embeddedPostgres, db, err := createDatabase()
			if err != nil {
				t.Errorf("Error creating database: %v", err)
				t.FailNow()
			}

			defer func() {
				err2 := embeddedPostgres.Stop()
				if err2 != nil {
					t.Errorf("error stopping embedded database: %v", err2)
				}
			}()

			dt := &devicetwin.MockDeviceTwin{}

			publishSnaps := messages.PublishSnaps{
				Action:  tt.args.expectedPublishSnapsMessage.Action,
				Id:      tt.args.expectedPublishSnapsMessage.Id,
				Message: tt.args.expectedPublishSnapsMessage.Message,
				Success: tt.args.expectedPublishSnapsMessage.Success,
			}

			if tt.args.expectedPublishSnapsMessage.Result != nil {
				publishSnaps.Result = tt.args.expectedPublishSnapsMessage.Result.Snaps
			}

			payload, err := json.Marshal(&publishSnaps)
			if err != nil {
				t.Errorf("error marshling message: %v", err)
				t.FailNow()
			}
			dt.On("ActionResponse", tt.args.expectedDeviceID, tt.args.expectedPublishSnapsMessage.Id, tt.args.expectedPublishSnapsMessage.Action, payload).Return(nil)

			srv.db = db
			srv.twin = dt

			if err3 := srv.actionMessageHandler(mockMessage); (err3 != nil) != tt.wantErr {
				t.Errorf("actionMessageHandler() error = %v, wantErr %v", err3, tt.wantErr)
			}
		})
	}
}

func TestService_healthMessageHandler(t *testing.T) {
	type fields struct {
	}
	type args struct {
		expectedTopic         string
		expectedHealthMessage messages.Health
		previousHealthMessage *messages.Health
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantErr              bool
		needsDatabase        bool
		addHealthHashesFirst bool
	}{
		{
			name:    "valid-no-hashes",
			fields:  fields{},
			args:    args{},
			wantErr: false,
		},
		{
			name:   "valid hashes not changed",
			fields: fields{},
			args: args{
				expectedTopic: "devices/health/1",
				expectedHealthMessage: messages.Health{
					DeviceId:           "1",
					InstalledSnapsHash: "ABC123",
					OrgId:              "RealOrgDotCom",
					Refresh:            time.Now().String(),
					SnapListHash:       "456DEF",
				},
			},
			wantErr:              false,
			needsDatabase:        true,
			addHealthHashesFirst: true,
		},
		{
			name:   "valid hashes did not exist",
			fields: fields{},
			args: args{
				expectedTopic: "devices/health/1",
				expectedHealthMessage: messages.Health{
					DeviceId:           "1",
					InstalledSnapsHash: "ABC123",
					OrgId:              "RealOrgDotCom",
					Refresh:            time.Now().String(),
					SnapListHash:       "456DEF",
				},
			},
			wantErr:       false,
			needsDatabase: true,
		},
		{
			name:   "valid hashes exist but changed",
			fields: fields{},
			args: args{
				expectedTopic: "devices/health/1",
				expectedHealthMessage: messages.Health{
					DeviceId:           "1",
					InstalledSnapsHash: "ABC123",
					OrgId:              "RealOrgDotCom",
					Refresh:            time.Now().String(),
					SnapListHash:       "456DEF",
				},
				previousHealthMessage: &messages.Health{
					DeviceId:           "1",
					InstalledSnapsHash: "XYZ789",
					OrgId:              "RealOrgDotCom",
					Refresh:            time.Now().String(),
					SnapListHash:       "101112TUV",
				},
			},
			wantErr:              false,
			needsDatabase:        true,
			addHealthHashesFirst: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedMessage := &mocks.Message{}
			mockedMessage.On("Topic").Return(tt.args.expectedTopic)
			bytes, err := json.Marshal(&tt.args.expectedHealthMessage)
			if err != nil {
				assert.Error(t, err)
			}
			mockedMessage.On("Payload").Return(bytes)

			healthChan := make(chan MQTT.Message)

			srv := &Service{
				legacyHealthChan: healthChan,
			}

			if tt.needsDatabase {
				embeddedPostgres, db, err := createDatabase()
				if err != nil {
					t.Errorf("Error creating database: %v", err)
					t.FailNow()
				}

				defer func() {
					err2 := embeddedPostgres.Stop()
					if err2 != nil {
						t.Errorf("error stopping embedded database: %v", err2)
					}
				}()

				srv.db = db
			}

			if tt.addHealthHashesFirst {
				var healthHash models.HealthHash
				if tt.args.previousHealthMessage != nil {
					healthHash.OrgID = tt.args.previousHealthMessage.OrgId
					healthHash.DeviceID = tt.args.previousHealthMessage.DeviceId
					healthHash.SnapListHash = tt.args.previousHealthMessage.SnapListHash
					healthHash.InstalledSnapsHash = tt.args.previousHealthMessage.InstalledSnapsHash
				} else {
					healthHash.OrgID = tt.args.expectedHealthMessage.OrgId
					healthHash.DeviceID = tt.args.expectedHealthMessage.DeviceId
					healthHash.SnapListHash = tt.args.expectedHealthMessage.SnapListHash
					healthHash.InstalledSnapsHash = tt.args.expectedHealthMessage.InstalledSnapsHash
				}
			}

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				if err = srv.healthMessageHandler(mockedMessage); (err != nil) != tt.wantErr {
					t.Errorf("healthMessageHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
				wg.Done()
			}()

			msg := <-healthChan
			assert.Equal(t, mockedMessage, msg)
			wg.Wait()

			if tt.needsDatabase {
				var hh models.HealthHash
				srv.db.Find(&hh, &models.HealthHash{DeviceID: tt.args.expectedHealthMessage.DeviceId})
				assert.Equal(t, tt.args.expectedHealthMessage.SnapListHash, hh.SnapListHash)
				assert.Equal(t, tt.args.expectedHealthMessage.InstalledSnapsHash, hh.InstalledSnapsHash)
			}
		})
	}
}

func TestService_sendLegacyActionHandlerUnversionedPayload(t *testing.T) {
	type fields struct {
	}
	type args struct {
		expectedTopic               string
		expectedPublishSnapsMessage messages.PublishSnapsV2
		expectedDeviceID            string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "valid",
			fields: fields{},
			args: args{
				expectedTopic: "device/health/some-device-id", expectedPublishSnapsMessage: messages.PublishSnapsV2{Version: "2", Action: "list", Success: true, Id: "1029384756"}, expectedDeviceID: "some-device-id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{}

			dt := &devicetwin.MockDeviceTwin{}

			publishSnaps := messages.PublishSnaps{
				Action:  tt.args.expectedPublishSnapsMessage.Action,
				Id:      tt.args.expectedPublishSnapsMessage.Id,
				Message: tt.args.expectedPublishSnapsMessage.Message,
				Success: tt.args.expectedPublishSnapsMessage.Success,
			}

			if tt.args.expectedPublishSnapsMessage.Result != nil {
				publishSnaps.Result = tt.args.expectedPublishSnapsMessage.Result.Snaps
			}

			payload, err := json.Marshal(&publishSnaps)
			if err != nil {
				t.Errorf("error marshling message: %v", err)
				t.FailNow()
			}
			dt.On("ActionResponse", tt.args.expectedDeviceID, tt.args.expectedPublishSnapsMessage.Id, tt.args.expectedPublishSnapsMessage.Action, payload).Return(nil)

			srv.twin = dt

			mockedMessage := &mocks.Message{}
			payload2, err2 := json.Marshal(&tt.args.expectedPublishSnapsMessage)
			if err2 != nil {
				t.Errorf("error marshling message: %v", err2)
				t.FailNow()
			}
			mockedMessage.On("Payload").Return(payload2)

			versionedMessage := messages.VersionedMessage{
				Action:  tt.args.expectedPublishSnapsMessage.Action,
				Id:      "1029384756",
				Success: true,
			}
			if err = srv.sendLegacyActionHandlerUnversionedPayload(tt.args.expectedDeviceID, mockedMessage, versionedMessage); (err != nil) != tt.wantErr {
				t.Errorf("sendLegacyActionHandlerUnversionedPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getClientID(t *testing.T) {
	type args struct {
		expectedClientID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{expectedClientID: "1029384756"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedMessage := &mocks.Message{}
			mockedMessage.On("Topic").Return(fmt.Sprintf("device/health/%s", tt.args.expectedClientID))
			if got := getClientID(mockedMessage); got != tt.args.expectedClientID {
				t.Errorf("getClientID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createDatabase() (*embeddedpostgres.EmbeddedPostgres, *gorm.DB, error) {
	embeddedPostgres := embeddedpostgres.NewDatabase()
	err := embeddedPostgres.Start()
	if err != nil {
		return nil, nil, err
	}

	dataSource := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	err = migrate2.Run(dataSource, "postgres", fmt.Sprintf("file://%s", "./../../db/migrations"), "postgres")
	if err != nil {
		return nil, nil, err
	}

	return embeddedPostgres, db, nil
}

func testDeviceTwinDatastore(ds devicetwindatastore.DataStore, err error) func(databaseDriver string, dataStoreSource string) (devicetwindatastore.DataStore, error) {
	return func(databaseDriver string, dataStoreSource string) (devicetwindatastore.DataStore, error) {
		return ds, err
	}
}
