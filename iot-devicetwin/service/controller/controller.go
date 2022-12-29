// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Device Twin Service
 * Copyright 2019 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU Affero General Public License version 3, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY,
 * SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.
 * See the GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/actions"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/everactive/dmscore/iot-devicetwin/domain"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
	"github.com/segmentio/ksuid"
)

// UnscopedController is an instance of the controller for accessing (soft) deleted data
type UnscopedController interface {
	DeviceGetByID(clientID string) (*messages.Device, bool, error)
}

// Controller interface for the service
type Controller interface {
	// MQTT handlers
	HealthHandler(msg MQTT.Message)
	ActionHandler(msg MQTT.Message)

	// Passthrough to the device twin service
	DeviceSnaps(orgID, clientID string) ([]messages.DeviceSnap, error)
	DeviceList(orgID string) ([]messages.Device, error)
	DeviceDelete(deviceID string) error
	DeviceGet(orgID, clientID string) (messages.Device, error)
	DeviceLogs(orgID, clientID string, logData *messages.DeviceLogs) error
	GroupCreate(orgID, name string) error
	GroupList(orgID string) ([]domain.Group, error)
	GroupGet(orgID, name string) (domain.Group, error)
	GroupLinkDevice(orgID, name, clientID string) error
	GroupUnlinkDevice(orgID, name, clientID string) error
	GroupGetDevices(orgID, name string) ([]messages.Device, error)
	GroupGetExcludedDevices(orgID, name string) ([]messages.Device, error)

	// Actions on a device
	DeviceSnapList(orgID, clientID string) error
	DeviceSnapInstall(orgID, clientID, snap string) error
	DeviceSnapServiceAction(orgID, clientID, snap, action string, services *messages.SnapService) error
	DeviceSnapRemove(orgID, clientID, snap string) error
	DeviceSnapUpdate(orgID, clientID, snap, action string, snapUpdate *messages.SnapUpdate) error
	DeviceSnapConf(orgID, clientID, snap, settings string) error
	DeviceSnapSnapshot(orgID, clientID, snap string, s3data *messages.SnapSnapshot) error
	ActionList(orgID, clientID string) ([]domain.Action, error)
	User(orgID, clientID string, user messages.DeviceUser) error
}

const (
	clientIDMQTTTopicPartsCount = 4
)

// Unscoped gets an Unscoped instance of the service for accessing (soft) deleted data
func (srv *Service) Unscoped() UnscopedController {
	return &Service{DeviceTwin: srv.DeviceTwin, unscoped: true}
}

// Service implementation of the devicetwin service use cases
type Service struct {
	DeviceTwin  devicetwin.DeviceTwin
	unscoped    bool
	healthChan  <-chan MQTT.Message
	actionChan  <-chan MQTT.Message
	publishChan chan<- mqtt.PublishMessage
}

// NewService creates an implementation of the devicetwin use cases
func NewService(healthChan chan MQTT.Message, actionChan <-chan MQTT.Message, publishChan chan<- mqtt.PublishMessage, twin devicetwin.DeviceTwin) *Service {
	srv := &Service{
		DeviceTwin:  twin,
		healthChan:  healthChan,
		actionChan:  actionChan,
		publishChan: publishChan,
	}

	return srv
}

func (srv *Service) Serve(ctx context.Context) error {
	intervalTicker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Errorf("We're done: %s", ctx.Err())
			return nil
		case <-intervalTicker.C:
			log.Infof("%s still ticking", "DeviceTwinControllerService")
		case a := <-srv.actionChan:
			log.Infof("Processing action channel message: %s", string(a.Payload()))
			srv.ActionHandler(a)
		case h := <-srv.healthChan:
			log.Infof("Processing health channel message: %s", string(h.Payload()))
			srv.HealthHandler(h)
		}
	}

	return nil
}

// ActionHandler is the handler for the main subscription topic
func (srv *Service) ActionHandler(msg MQTT.Message) {
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

	// Handle the action
	if err := srv.DeviceTwin.ActionResponse(clientID, a.Id, a.Action, msg.Payload()); err != nil {
		log.Printf("Error with action `%s`: %v", a.Action, err)
	}
}

// HealthHandler is the handler for the devices health messages
func (srv *Service) HealthHandler(msg MQTT.Message) {
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

	// Check to make sure it's not a device we've deleted (soft)
	// If it isn't then we can handle it even if we don't know about the device yet (expected)
	device, isDeleted, err := srv.Unscoped().DeviceGetByID(clientID)

	// If there's an error AND it's something other than record not found, it's a real error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("Error trying to get device by id = %s: %v. Will try to handle it anyway.", clientID, err)
	}

	if isDeleted {
		// We've previously soft-deleted this device, which means we need to ignore it now
		log.Tracef("Dropping messages for %s (%s), previously deleted", device.Serial, device.DeviceId)
		return
	}

	// Update the device record
	if err2 := srv.DeviceTwin.HealthHandler(h); err2 == nil {
		// Exit if successful
		return
	} else {
		log.Error(err2)
	}

	// We don't have the device details, so request them from the device
	act := messages.SubscribeAction{
		Action: actions.Device,
	}
	if err := srv.triggerActionOnDevice(h.OrgId, h.DeviceId, act); err != nil {
		log.Printf("Triggering action: %v", err)
	}

	// Get the snaps from the device
	act = messages.SubscribeAction{
		Action: actions.List,
	}
	if err := srv.triggerActionOnDevice(h.OrgId, h.DeviceId, act); err != nil {
		log.Printf("Triggering action: %v", err)
	}
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

var generateKSUID = ksuid.New

// triggerActionOnDevice triggers an action on the device via MQTT
func (srv *Service) triggerActionOnDevice(orgID, deviceID string, act messages.SubscribeAction) error {
	// Generate a request ID
	id := generateKSUID()
	act.Id = id.String()

	// Serialize the action
	data, err := serializePayload(act)
	if err != nil {
		log.Printf("Error in action serialization: %v", err)
		return err
	}

	// Publish the request
	t := fmt.Sprintf("devices/sub/%s", deviceID)
	pubMessage := mqtt.PublishMessage{Topic: t, Payload: string(data)}
	log.Infof("Triggering action on device, sending pubMessage: %s", string(data))
	srv.publishChan <- pubMessage
	if err != nil {
		log.Printf("Error in publish: %v", err)
		return fmt.Errorf("error in publish: %v", err)
	}

	// Log the request
	return srv.DeviceTwin.ActionCreate(orgID, deviceID, act)
}

func serializePayload(act messages.SubscribeAction) ([]byte, error) {
	return json.Marshal(act)
}
