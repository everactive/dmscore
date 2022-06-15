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
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
)

// DeviceVersionGet fetches a device version details from the database
func (db *DataStore) DeviceVersionGet(deviceID int64) (datastore.DeviceVersion, error) {
	item := datastore.DeviceVersion{}
	res := db.gormDB.Where("device_id = ?", deviceID).Find(&item)
	if res.Error != nil {
		log.Errorf("Error retrieving device version: %v", res.Error)
	}

	return item, res.Error
}

// DeviceVersionUpsert creates or updates a device version record
func (db *DataStore) DeviceVersionUpsert(dv datastore.DeviceVersion) error {
	res := db.gormDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}},
		UpdateAll: true,
	}).Create(&dv)

	return res.Error
}

// DeviceVersionDelete removes a device version
func (db *DataStore) DeviceVersionDelete(id int64) error {
	res := db.gormDB.Where("id = ?", id).Delete(&datastore.DeviceVersion{})
	return res.Error
}
