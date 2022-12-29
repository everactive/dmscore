package devicetwin

import (
	"context"
	"errors"
	"fmt"
	"github.com/everactive/dmscore/config/keys"
	devicetwindatastore "github.com/everactive/dmscore/iot-devicetwin/datastore"
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
	identitydatastore "github.com/everactive/dmscore/iot-identity/datastore"
	"github.com/everactive/dmscore/iot-identity/domain"
	"github.com/everactive/dmscore/models"
	"github.com/everactive/dmscore/pkg/datastores"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/thejerf/suture/v4"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
	"sync"
	"testing"
	"time"
)

func TestChecker_RefreshDevices(t *testing.T) {
	type fields struct {
		devicesIDs        []string
		deviceList        map[string]*devicetwindatastore.Device
		currentModels     map[string][]models.DeviceModelRequiredSnap
		deviceMutex       *sync.Mutex
		legacyPublishChan chan mqtt.PublishMessage
		deviceCheckTicker *time.Ticker
	}
	type args struct {
		dss                 datastores.DataStores
		organizationList    []domain.Organization
		organizationListErr error
		deviceList          []devicetwindatastore.Device
		deviceListErr       error
		deviceModel         *models.DeviceModel
		deviceSnapList      []devicetwindatastore.DeviceSnap
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error on organization list",
			fields: fields{
				deviceMutex: &sync.Mutex{},
			},
			args: args{
				organizationListErr: errors.New("this is an error"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				// return assert.Errorf(t, err, "this is an error")
				return assert.Equal(t, errors.New("this is an error"), err)
			},
		},
		{
			name: "error on device list",
			fields: fields{
				deviceMutex: &sync.Mutex{},
			},
			args: args{
				organizationList: []domain.Organization{{ID: "1", Name: "something", RootCert: []byte{}, RootKey: []byte{}}},
				deviceListErr:    errors.New("this is an error"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equalf(t, errors.New("this is an error"), err, "Errors did not match")
			},
		},
		{
			name: "valid device",
			fields: fields{
				deviceMutex:       &sync.Mutex{},
				deviceList:        map[string]*devicetwindatastore.Device{},
				deviceCheckTicker: time.NewTicker(1 * time.Second),
			},
			args: args{
				organizationList: []domain.Organization{{ID: "1", Name: "something", RootCert: []byte{}, RootKey: []byte{}}},
				deviceModel: &models.DeviceModel{
					Name:                     "test-model-2000",
					DeviceModelRequiredSnaps: []models.DeviceModelRequiredSnap{},
				},
				deviceList: []devicetwindatastore.Device{
					{
						Model:       gorm.Model{ID: 1},
						DeviceModel: "test-model-2000",
					},
				},
				deviceSnapList: []devicetwindatastore.DeviceSnap{
					{},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		viper.Set(keys.RequiredSnapsCheckInterval, 1*time.Second)

		mockedIdentityStore := identitydatastore.MockDataStore{}
		mockedIdentityStore.On("OrganizationList").Return(tt.args.organizationList, tt.args.organizationListErr)

		mockedDevicetwinStore := devicetwindatastore.MockDataStore{}
		mockedDataStore := datastores.MockDataStore{}
		if tt.args.deviceModel != nil {
			mockedDataStore.On("GetModelRequiredSnaps", tt.args.deviceModel.Name).Return(tt.args.deviceModel, nil)
		}
		if len(tt.args.organizationList) > 0 {
			mockedDevicetwinStore.On("DeviceList", tt.args.organizationList[0].ID).Return(tt.args.deviceList, tt.args.deviceListErr)
		}
		if len(tt.args.deviceSnapList) > 0 {
			mockedDevicetwinStore.On("DeviceSnapList", int64(tt.args.deviceList[0].ID)).Return(tt.args.deviceSnapList, nil)
		}

		tt.args.dss.IdentityStore = &mockedIdentityStore
		tt.args.dss.DeviceTwinStore = &mockedDevicetwinStore
		tt.args.dss.DataStore = &mockedDataStore

		t.Run(tt.name, func(t *testing.T) {
			c := &Checker{
				devicesIDs:        tt.fields.devicesIDs,
				deviceList:        tt.fields.deviceList,
				currentModels:     tt.fields.currentModels,
				deviceMutex:       *tt.fields.deviceMutex,
				legacyPublishChan: tt.fields.legacyPublishChan,
				deviceCheckTicker: tt.fields.deviceCheckTicker,
			}
			tt.wantErr(t, c.RefreshDevices(&tt.args.dss), fmt.Sprintf("RefreshDevices(%v)", tt.args.dss))
			assert.Equal(t, len(tt.args.deviceList), len(c.currentModels))
			assert.Equal(t, len(tt.args.deviceList), len(maps.Keys(c.deviceList)))
		})
	}
}

func TestChecker_Serve(t *testing.T) {
	type fields struct {
		devicesIDs                []string
		deviceList                map[string]*devicetwindatastore.Device
		currentModels             map[string][]models.DeviceModelRequiredSnap
		deviceMutex               sync.Mutex
		legacyPublishChan         chan mqtt.PublishMessage
		deviceCheckTickerDuration time.Duration
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "done",
			fields: fields{
				deviceCheckTickerDuration: 5 * time.Second,
			},
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Checker{
				devicesIDs:        tt.fields.devicesIDs,
				deviceList:        tt.fields.deviceList,
				currentModels:     tt.fields.currentModels,
				deviceMutex:       tt.fields.deviceMutex,
				legacyPublishChan: tt.fields.legacyPublishChan,
				deviceCheckTicker: time.NewTicker(tt.fields.deviceCheckTickerDuration),
			}

			viper.Set(keys.DefaultServiceHeartbeat, 5*time.Second)

			ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				wg.Done()
				if err := c.Serve(ctx); (err != nil) != tt.wantErr {
					t.Errorf("Serve() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			wg.Wait()
			cancelFunc()

			assert.Equal(t, context.Canceled, ctx.Err())
		})
	}
}

func TestChecker_String(t *testing.T) {
	type fields struct {
		devicesIDs        []string
		deviceList        map[string]*devicetwindatastore.Device
		currentModels     map[string][]models.DeviceModelRequiredSnap
		deviceMutex       sync.Mutex
		legacyPublishChan chan mqtt.PublishMessage
		deviceCheckTicker *time.Ticker
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Checker{
				devicesIDs:        tt.fields.devicesIDs,
				deviceList:        tt.fields.deviceList,
				currentModels:     tt.fields.currentModels,
				deviceMutex:       tt.fields.deviceMutex,
				legacyPublishChan: tt.fields.legacyPublishChan,
				deviceCheckTicker: tt.fields.deviceCheckTicker,
			}
			assert.Equalf(t, tt.want, c.String(), "String()")
		})
	}
}

func TestChecker_checkNextDevice(t *testing.T) {
	type fields struct {
		devicesIDs        []string
		deviceList        map[string]*devicetwindatastore.Device
		currentModels     map[string][]models.DeviceModelRequiredSnap
		deviceMutex       sync.Mutex
		legacyPublishChan chan mqtt.PublishMessage
		deviceCheckTicker *time.Ticker
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Checker{
				devicesIDs:        tt.fields.devicesIDs,
				deviceList:        tt.fields.deviceList,
				currentModels:     tt.fields.currentModels,
				deviceMutex:       tt.fields.deviceMutex,
				legacyPublishChan: tt.fields.legacyPublishChan,
				deviceCheckTicker: tt.fields.deviceCheckTicker,
			}
			tt.wantErr(t, c.checkNextDevice(), fmt.Sprintf("checkNextDevice()"))
		})
	}
}

func TestInstallService_Serve(t *testing.T) {
	type fields struct {
		heartbeatInterval   time.Duration
		interval            time.Duration
		supervisor          *suture.Supervisor
		checkerServiceToken *suture.ServiceToken
		checkerService      *Checker
		checkerMutex        sync.Mutex
		stores              *datastores.DataStores
		legacyPublishChan   chan mqtt.PublishMessage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InstallService{
				heartbeatInterval:   tt.fields.heartbeatInterval,
				interval:            tt.fields.interval,
				supervisor:          tt.fields.supervisor,
				checkerServiceToken: tt.fields.checkerServiceToken,
				checkerService:      tt.fields.checkerService,
				checkerMutex:        tt.fields.checkerMutex,
				stores:              tt.fields.stores,
				legacyPublishChan:   tt.fields.legacyPublishChan,
			}
			tt.wantErr(t, i.Serve(tt.args.ctx), fmt.Sprintf("Serve(%v)", tt.args.ctx))
		})
	}
}

func TestInstallService_String(t *testing.T) {
	type fields struct {
		heartbeatInterval   time.Duration
		interval            time.Duration
		supervisor          *suture.Supervisor
		checkerServiceToken *suture.ServiceToken
		checkerService      *Checker
		checkerMutex        sync.Mutex
		stores              *datastores.DataStores
		legacyPublishChan   chan mqtt.PublishMessage
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InstallService{
				heartbeatInterval:   tt.fields.heartbeatInterval,
				interval:            tt.fields.interval,
				supervisor:          tt.fields.supervisor,
				checkerServiceToken: tt.fields.checkerServiceToken,
				checkerService:      tt.fields.checkerService,
				checkerMutex:        tt.fields.checkerMutex,
				stores:              tt.fields.stores,
				legacyPublishChan:   tt.fields.legacyPublishChan,
			}
			assert.Equalf(t, tt.want, i.String(), "String()")
		})
	}
}

func TestInstallService_createOrRefreshChecker(t *testing.T) {
	type fields struct {
		heartbeatInterval   time.Duration
		interval            time.Duration
		supervisor          *suture.Supervisor
		checkerServiceToken *suture.ServiceToken
		checkerService      *Checker
		checkerMutex        sync.Mutex
		stores              *datastores.DataStores
		legacyPublishChan   chan mqtt.PublishMessage
	}
	type args struct {
		legacyPublishChan chan mqtt.PublishMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InstallService{
				heartbeatInterval:   tt.fields.heartbeatInterval,
				interval:            tt.fields.interval,
				supervisor:          tt.fields.supervisor,
				checkerServiceToken: tt.fields.checkerServiceToken,
				checkerService:      tt.fields.checkerService,
				checkerMutex:        tt.fields.checkerMutex,
				stores:              tt.fields.stores,
				legacyPublishChan:   tt.fields.legacyPublishChan,
			}
			tt.wantErr(t, i.createOrRefreshChecker(tt.args.legacyPublishChan), fmt.Sprintf("createOrRefreshChecker(%v)", tt.args.legacyPublishChan))
		})
	}
}

func TestNewInstallService(t *testing.T) {
	type args struct {
		stores            *datastores.DataStores
		legacyPublishChan chan mqtt.PublishMessage
	}
	tests := []struct {
		name string
		args args
		want *suture.Supervisor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewInstallService(tt.args.stores, tt.args.legacyPublishChan), "NewInstallService(%v, %v)", tt.args.stores, tt.args.legacyPublishChan)
		})
	}
}

func Test_sanitizeSerial(t *testing.T) {
	type args struct {
		serial string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, sanitizeSerial(tt.args.serial), "sanitizeSerial(%v)", tt.args.serial)
		})
	}
}
