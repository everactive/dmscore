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
	"encoding/json"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/actions"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
)

// DeviceGetByID gets the device from the database cache
func (srv *Service) DeviceGetByID(clientID string) (*messages.Device, bool, error) {
	return srv.DeviceTwin.Unscoped().DeviceGetByID(clientID)
}

// DeviceGet gets the device from the database cache
func (srv *Service) DeviceGet(orgID, clientID string) (messages.Device, error) {
	return srv.DeviceTwin.DeviceGet(orgID, clientID)
}

// DeviceList gets the devices from the database cache
func (srv *Service) DeviceList(orgID string) ([]messages.Device, error) {
	return srv.DeviceTwin.DeviceList(orgID)
}

// DeviceDelete deletes the device from the database cache
func (srv *Service) DeviceDelete(deviceID string) error {
	_, err := srv.DeviceTwin.DeviceDelete(deviceID)
	return err
}

// deviceSnapAction triggers a device action on a device
func (srv *Service) deviceAction(orgID, clientID string, action messages.SubscribeAction) error {
	// Validate the org and device ID
	device, err := srv.DeviceTwin.DeviceGet(orgID, clientID)
	if err != nil {
		return err
	}

	// Trigger the action on the device
	err = srv.triggerActionOnDevice(device.OrgId, device.DeviceId, action)
	if err != nil {
		return err
	}
	return err
}

// DeviceLogs triggers an upload of snapd logs to S3
func (srv *Service) DeviceLogs(orgID, clientID string, logData *messages.DeviceLogs) error {
	jsonBytes, err := json.Marshal(logData)
	if err != nil {
		return err
	}

	act := messages.SubscribeAction{
		Action: actions.Logs,
		Data:   string(jsonBytes),
	}

	return srv.deviceSnapAction(orgID, clientID, act)
}
