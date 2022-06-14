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
	"gorm.io/gorm"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
)

// DeviceGet fetches a device details from the database cache
func (srv *Service) DeviceGet(orgID, deviceID string) (messages.Device, error) {
	d, deleted, err := srv.deviceGetByID(deviceID, srv.unscoped)
	if err != nil {
		return messages.Device{}, err
	}

	if deleted && !srv.unscoped {
		return messages.Device{}, errors.New("device not found")
	}

	// If the OrgIds don't match it's possible we are using the org name, so check that
	if d.OrgId != orgID {
		// Only do this lookup if we need to
		organization, err := srv.CoreDB.OrganizationGet(orgID)
		if err != nil {
			return messages.Device{}, err
		}
		if d.OrgId != organization.OrganizationID {
			return messages.Device{}, fmt.Errorf("the organization ID does not match the device")
		}
	}

	return *d, nil
}

func (srv *Service) deviceGetByID(deviceID string, unscoped bool) (*messages.Device, bool, error) {
	// Get the device
	var d datastore.Device
	var err error
	if unscoped {
		d, err = srv.DB.Unscoped().DeviceGet(deviceID)
		if err != nil {
			return nil, false, err
		}
	} else {
		d, err = srv.DB.DeviceGet(deviceID)
		if err != nil {
			return nil, false, err
		}
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

	return &device, d.DeletedAt != gorm.DeletedAt{}, nil
}

// DeviceGetByID gets a device by its id
func (srv *Service) DeviceGetByID(deviceID string) (*messages.Device, bool, error) {
	return srv.deviceGetByID(deviceID, srv.unscoped)
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
