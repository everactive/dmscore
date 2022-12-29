package devicetwin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everactive/dmscore/config/keys"
	"github.com/everactive/dmscore/iot-devicetwin/datastore"
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
	"github.com/everactive/dmscore/models"
	"github.com/everactive/dmscore/pkg/datastores"
	"github.com/everactive/dmscore/pkg/messages"
	ksuid2 "github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thejerf/suture/v4"
	"regexp"
	"sync"
	"time"
)

type InstallService struct {
	heartbeatInterval   time.Duration
	interval            time.Duration
	supervisor          *suture.Supervisor
	checkerServiceToken *suture.ServiceToken
	checkerService      *Checker
	checkerMutex        sync.Mutex
	stores              *datastores.DataStores
	legacyPublishChan   chan mqtt.PublishMessage
}

func NewInstallService(stores *datastores.DataStores, legacyPublishChan chan mqtt.PublishMessage) *suture.Supervisor {
	sup := suture.NewSimple("install")
	hearbeatInterval := viper.GetDuration(keys.DefaultServiceHeartbeat)
	interval := viper.GetDuration(keys.RequiredSnapsInstallServiceCheckInterval)
	i := &InstallService{
		heartbeatInterval: hearbeatInterval,
		interval:          interval,
		supervisor:        sup,
		stores:            stores,
		legacyPublishChan: legacyPublishChan,
	}

	sup.Add(i)

	return sup
}

func (i *InstallService) String() string {
	return "InstallService"
}

func (i *InstallService) Serve(ctx context.Context) error {
	intervalTicker := time.NewTicker(i.heartbeatInterval)
	installTicker := time.NewTicker(i.interval)

	for {
		select {
		case <-ctx.Done():
			logger.Errorf("We're done: %s", ctx.Err())
			return nil
		case <-intervalTicker.C:
			logger.Infof("%s still ticking", i.String())
		case <-installTicker.C:
			logger.Info("Starting check to see if any devices need required snaps")
			err := i.createOrRefreshChecker(i.legacyPublishChan)
			if err != nil {
				return err
			}

		}
	}
}

func (i *InstallService) createOrRefreshChecker(legacyPublishChan chan mqtt.PublishMessage) error {
	i.checkerMutex.Lock()
	defer i.checkerMutex.Unlock()
	if i.checkerServiceToken != nil {
		// a checker service is already running, if it is just refresh
		log.Trace("Checker service already exists, just refresh")
		err := i.checkerService.RefreshDevices(i.stores)
		if err != nil {
			return err
		}
		return nil
	}

	checkIntervalDuration := viper.GetDuration(keys.RequiredSnapsCheckInterval)

	i.checkerService = NewChecker(legacyPublishChan, checkIntervalDuration)

	serviceToken := i.supervisor.Add(i.checkerService)
	i.checkerServiceToken = &serviceToken

	err := i.checkerService.RefreshDevices(i.stores)
	if err != nil {
		i.checkerService = nil
		i.checkerServiceToken = nil
		return err
	}
	return nil
}

func NewChecker(legacyPublishChan chan mqtt.PublishMessage, checkIntervalDuration time.Duration) *Checker {
	return &Checker{
		legacyPublishChan: legacyPublishChan,
		deviceList:        map[string]*datastore.Device{},
		currentModels:     map[string][]models.DeviceModelRequiredSnap{},
		deviceCheckTicker: time.NewTicker(checkIntervalDuration),
	}
}

type Checker struct {
	devicesIDs        []string
	deviceList        map[string]*datastore.Device
	currentModels     map[string][]models.DeviceModelRequiredSnap
	deviceMutex       sync.Mutex
	legacyPublishChan chan mqtt.PublishMessage
	deviceCheckTicker *time.Ticker
}

func (c *Checker) RefreshDevices(dss *datastores.DataStores) error {
	log.Tracef("Refreshing devices")

	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()
	orgs, err := dss.IdentityStore.OrganizationList()
	if err != nil {
		log.Errorf("error getting organization list: %v", err)
		return err
	}

	c.currentModels = make(map[string][]models.DeviceModelRequiredSnap)

	log.Tracef("Starting deviceList size %d", len(c.deviceList))
	log.Tracef("Starting deviceIDs size %d", len(c.devicesIDs))

	// get all the devices, for all orgs
	for _, org := range orgs {
		log.Tracef("Getting device list for org %s", org.Name)
		list, err2 := dss.DeviceTwinStore.DeviceList(org.ID)
		if err2 != nil {
			return err2
		}

		deviceIDsToBeAdded := []string{}
		for _, device := range list {
			// If we've already added this model to the current models, then we don't
			// need to do it again
			if _, ok := c.currentModels[device.DeviceModel]; !ok {
				deviceModel, err3 := dss.DataStore.GetModelRequiredSnaps(device.DeviceModel)
				if err3 != nil {
					log.Errorf("error trying to get required snaps for model %s", device.DeviceModel)
					break
				}

				log.Infof("Required snaps for model %s are: %+v", deviceModel.Name, deviceModel.DeviceModelRequiredSnaps)
				c.currentModels[device.DeviceModel] = deviceModel.DeviceModelRequiredSnaps
			} else {
				log.Tracef("Model %s already exists in the current models", device.DeviceModel)
			}

			deviceToAdd := device
			if d, ok := c.deviceList[device.DeviceID]; !ok {
				// Get snaps for device first
				snapList, err3 := dss.DeviceTwinStore.DeviceSnapList(int64(device.ID))
				if err3 != nil {
					log.Errorf("error trying to get snaps for device %s: %s", d.SerialNumber, err3)
					break
				}

				for _, s := range snapList {
					snapToAdd := s
					log.Tracef("Adding snap %s to device %s", snapToAdd.Name, deviceToAdd.SerialNumber)
					deviceToAdd.DeviceSnaps = append(deviceToAdd.DeviceSnaps, &snapToAdd)
				}

				log.Tracef("Adding device to device list: %s, %+v", device.DeviceID, deviceToAdd)
				c.deviceList[device.DeviceID] = &deviceToAdd
				deviceIDsToBeAdded = append(deviceIDsToBeAdded, device.DeviceID)
			} else {
				log.Tracef("device %s already exists in device list", d.DeviceID)
			}
		}

		log.Tracef("Devices to be added for this org:")
		log.Tracef("%+v", deviceIDsToBeAdded)
		c.devicesIDs = append(c.devicesIDs, deviceIDsToBeAdded...)
	}

	log.Tracef("Device ids:")
	log.Tracef("%+v", c.devicesIDs)
	log.Tracef("Device list:")
	log.Tracef("%+v", c.deviceList)

	// Make sure our checker ticker is running on the interval we want
	// (it could have been stopped if we ran out of devices to check last time)
	checkIntervalDuration := viper.GetDuration(keys.RequiredSnapsCheckInterval)
	c.deviceCheckTicker = time.NewTicker(checkIntervalDuration)

	return nil
}

func (c *Checker) String() string {
	return "Required Snaps Checker"
}

func (c *Checker) Serve(ctx context.Context) error {
	heartbeatInterval := viper.GetDuration(keys.DefaultServiceHeartbeat)
	intervalTicker := time.NewTicker(heartbeatInterval)
	for {
		select {
		case <-ctx.Done():
			logger.Errorf("We're done: %s", ctx.Err())
			return nil
		case <-intervalTicker.C:
			logger.Infof("%s still ticking", c.String())
		case <-c.deviceCheckTicker.C:
			logger.Trace("Checking next device")
			err := c.checkNextDevice()
			if err != nil && err != ErrNoMoreDevices {
				return err
			}
			if err == ErrNoMoreDevices {
				// If we ran out of devices to check, stop checking and wait for the next time
				// the install service runs the checkers
				log.Infof("No more devices to check, stopping timer")
				c.deviceCheckTicker.Stop()
			}
		}
	}
}

var ErrNoMoreDevices = errors.New("no more devices to check")

func sanitizeSerial(serial string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9_-]*")
	serialProcessed := re.ReplaceAllString(serial, "")
	return serialProcessed
}

func (c *Checker) checkNextDevice() error {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()
	if len(c.devicesIDs) == 0 {
		log.Infof("No more devives to check, returning ErrNoMoreDevices")
		// Basically don't tick any more and it will get reset when the new check run is done
		return ErrNoMoreDevices
	}

	nextDevice := c.deviceList[c.devicesIDs[0]]

	log.Tracef("Checking device id=%s, serial=%s", nextDevice.DeviceID, nextDevice.SerialNumber)

	requiredSnaps := c.currentModels[nextDevice.DeviceModel]
	requiredForThisDevice := []*messages.SnapsItems{}
	for _, requiredSnap := range requiredSnaps {
		found := false
		for _, snap := range nextDevice.DeviceSnaps {
			if snap.Name == requiredSnap.Name {
				found = true
				break
			}
		}

		if !found {
			requiredForThisDevice = append(requiredForThisDevice, &messages.SnapsItems{
				Channel: "latest",
				Name:    requiredSnap.Name,
				Track:   "stable",
			})
		}
	}

	if len(requiredForThisDevice) == 0 {
		log.Tracef("No required snaps for device id=%s, serial=%s", nextDevice.DeviceID, nextDevice.SerialNumber)
		c.devicesIDs = c.devicesIDs[1:]
		return nil
	}

	t := fmt.Sprintf("devices/actions/%s/required-install", sanitizeSerial(nextDevice.SerialNumber))
	m := messages.RequiredInstall{
		Id:    ksuid2.New().String(),
		Snaps: requiredForThisDevice,
	}

	log.Infof("Publishing required-install action for %s: %+v", nextDevice.DeviceID, m)

	bytes, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	pubMessage := mqtt.PublishMessage{Topic: t, Payload: string(bytes)}
	c.legacyPublishChan <- pubMessage

	delete(c.deviceList, c.devicesIDs[0])
	c.devicesIDs = c.devicesIDs[1:]

	return nil
}
