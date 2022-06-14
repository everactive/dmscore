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

package devicetwin

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/actions"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
	"github.com/everactive/dmscore/iot-devicetwin/domain"

	managementdatastore "github.com/everactive/dmscore/iot-management/datastore"
)

// UnscopedDeviceTwin gets data and includes (soft) deleted data as well
type UnscopedDeviceTwin interface {
	DeviceGetByID(clientID string) (*messages.Device, bool, error)
}

// DeviceTwin interface for the service
type DeviceTwin interface {
	HealthHandler(payload messages.Health) error
	ActionResponse(clientID, actionID, action string, payload []byte) error // process a response from a device

	ActionCreate(orgID, deviceID string, act messages.SubscribeAction) error
	ActionUpdate(actionID, status, message string) error
	ActionList(orgID, deviceID string) ([]domain.Action, error)

	DeviceSnaps(orgID, clientID string) ([]messages.DeviceSnap, error)

	DeviceList(orgID string) ([]messages.Device, error)
	DeviceGet(orgID, clientID string) (messages.Device, error)
	DeviceDelete(deviceID string) (string, error)

	GroupCreate(orgID, name string) error
	GroupList(orgID string) ([]domain.Group, error)
	GroupGet(orgID, name string) (domain.Group, error)
	GroupLinkDevice(orgID, name, clientID string) error
	GroupUnlinkDevice(orgID, name, clientID string) error
	GroupGetDevices(orgID, name string) ([]messages.Device, error)
	GroupGetExcludedDevices(orgID, name string) ([]messages.Device, error)

	Unscoped() UnscopedDeviceTwin
}

// Service implementation of the identity use cases
type Service struct {
	DB       datastore.DataStore
	unscoped bool
	CoreDB   managementdatastore.DataStore
}

// NewService creates an implementation of the device twin use cases
func NewService(db datastore.DataStore, coreDB managementdatastore.DataStore) *Service {
	return &Service{
		DB:     db,
		CoreDB: coreDB,
	}
}

// HealthHandler handles a health update from a device
func (srv *Service) HealthHandler(payload messages.Health) error {
	// Check that we have the device
	_, err := srv.DB.DeviceGet(payload.DeviceId)
	if err != nil {
		// Request the device details to be published as we don't have it
		log.Tracef("We don't know about this device yet: DeviceId=%s, OrgId=%s", payload.DeviceId, payload.OrgId)
		return err
	}

	// Update the last refresh on the device
	return srv.DB.DevicePing(payload.DeviceId, payload.Refresh)
}

// ActionResponse handles action response from a device
func (srv *Service) ActionResponse(clientID, actionID, action string, payload []byte) error {
	var (
		err     error
		status  = "complete"
		message = ""
	)

	log.Printf("Action: %s", action)

	// Act based on the message action
	switch action {
	case actions.Device:
		err = srv.actionDevice(payload)
	case actions.List:
		err = srv.actionList(clientID, payload)
	case actions.Install, actions.Remove, actions.Refresh, actions.Revert,
		actions.Enable, actions.Disable, actions.SetConf, actions.Start,
		actions.Stop, actions.Switch, actions.Restart:
		message, err = srv.actionForSnap(clientID, action, payload)
	case actions.Conf, actions.Info:
		err = srv.actionConf(clientID, payload)
	// case "ack":
	case actions.Server:
		err = srv.actionServer(clientID, payload)
	case actions.Unregister:
		err = srv.actionUnregister(clientID, payload)
	default:
		return fmt.Errorf("error unhandled action `%s`", action)
	}

	// Update the action status
	if err != nil {
		status = "error"
		message = err.Error()
	}
	e := srv.ActionUpdate(actionID, status, message)
	if e != nil {
		log.Printf("Error updating action `%s`: %v", actionID, e)
	}
	return err // return the response from the original action
}

// Unscoped gets an instance of the service that is able to access (soft) deleted data
func (srv *Service) Unscoped() UnscopedDeviceTwin {
	unscopedService := *srv
	unscopedService.unscoped = true
	return &unscopedService
}
