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

package postgres

import (
	"time"

	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
)

// DeviceCreate adds a new record to device database table, returning the record ID
func (db *DataStore) DeviceCreate(device datastore.Device) (int64, error) {
	log.Tracef("Attempting to create device: %+v", device)
	res := db.gormDB.Create(&device)
	if res.Error != nil {
		log.Error(res.Error)
		return 0, res.Error
	}

	return int64(device.ID), nil
}

// DeviceGet fetches a device from the database
func (db *DataStore) DeviceGet(deviceID string) (datastore.Device, error) {
	item := datastore.Device{}

	tx := db.gormDB
	if db.unscoped {
		tx = tx.Unscoped()
	}

	res := tx.Preload("DeviceVersion").Where("device_id = ?", deviceID).Or("serial = ?", deviceID).First(&item)
	if res.RowsAffected == 1 {
		return item, res.Error
	}

	return item, res.Error
}

// DevicePing updates the last ping time from a device
func (db *DataStore) DevicePing(deviceID string, refresh time.Time) error {
	res := db.gormDB.Model(&datastore.Device{}).Where("device_id = ?", deviceID).Update("lastrefresh", refresh)
	if res.Error != nil {
		log.Error(res.Error)
		return res.Error
	}
	return nil
}

// DeviceDelete deletes the device
func (db *DataStore) DeviceDelete(deviceID string) error {

	device := datastore.Device{}
	db.gormDB.Where(&datastore.Device{DeviceID: deviceID}).Find(&device)
	db.gormDB.Delete(&datastore.Device{Model: gorm.Model{ID: device.ID}})

	return nil
}

// DeviceList fetches the devices for an organization from the database
func (db *DataStore) DeviceList(orgID string) ([]datastore.Device, error) {
	devices := []datastore.Device{}
	db.gormDB.Where(&datastore.Device{OrganisationID: orgID}).Find(&devices)

	log.Tracef("Devices: %+v", devices)

	return devices, nil
}
