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
	"errors"
	"fmt"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
)

// DeviceGet fetches a device details from the database cache
func (srv *Service) DeviceGet(orgID, deviceID string) (messages.Device, error) {
	d, isDeleted, err := srv.DeviceGetByID(deviceID)
	if err != nil {
		return messages.Device{}, err
	}

	if isDeleted {
		return messages.Device{}, errors.New("device found but is deleted")
	}

	// Validate the supplied orgid
	if d.OrgId != orgID {
		return messages.Device{}, fmt.Errorf("the organization ID does not match the device")
	}

	return *d, nil
}

// DeviceGetByID gets a device by its id
func (srv *Service) DeviceGetByID(deviceID string) (*messages.Device, bool, error) {
	// Get the device
	d, err := srv.DB.Unscoped().DeviceGet(deviceID)
	if err != nil {
		log.Error(err)
		return nil, false, err
	}

	device := dataToDomainDevice(d)

	device.Version = &messages.DeviceVersion{
		DeviceId:      d.DeviceID,
		Version:       d.DeviceVersion.Version,
		Series:        d.DeviceVersion.Series,
		OsId:          d.DeviceVersion.OSID,
		OsVersionId:   d.DeviceVersion.OSVersionID,
		OnClassic:     d.DeviceVersion.OnClassic,
		KernelVersion: d.DeviceVersion.KernelVersion,
	}

	return &device, d.IsDeleted(), nil
}

// DeviceDelete deletes the device from the database
func (srv *Service) DeviceDelete(deviceID string) (string, error) {
	err := srv.DB.DeviceDelete(deviceID)
	if err != nil {
		return "failed to delete device", err
	}

	return deviceID, nil
}

// DeviceList fetches devices from the database cache
func (srv *Service) DeviceList(orgID string) ([]messages.Device, error) {
	dd, err := srv.DB.DeviceList(orgID)
	if err != nil {
		return nil, err
	}

	devices := []messages.Device{}
	for _, d := range dd {
		devices = append(devices, dataToDomainDevice(d))
	}
	return devices, nil
}

func dataToDomainDevice(d datastore.Device) messages.Device {
	return messages.Device{
		OrgId:       d.OrganisationID,
		DeviceId:    d.DeviceID,
		Brand:       d.Brand,
		Model:       d.DeviceModel,
		Serial:      d.SerialNumber,
		Store:       d.StoreID,
		DeviceKey:   d.DeviceKey,
		Version:     &messages.DeviceVersion{},
		Created:     d.CreatedAt,
		LastRefresh: d.LastRefresh,
	}
}
