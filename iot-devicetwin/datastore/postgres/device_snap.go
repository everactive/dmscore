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
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
	"gorm.io/gorm/clause"
)

// DeviceSnapUpsert creates or updates a device snap record
func (db *DataStore) DeviceSnapUpsert(ds datastore.DeviceSnap) error {
	tx := db.gormDB.Begin()
	tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "device_id"}},
		UpdateAll: true,
	}).Omit("ServiceStatuses").Create(&ds)

	if len(ds.ServiceStatuses) > 0 {
		for _, ss := range ds.ServiceStatuses {
			ss.DeviceSnapID = int64(ds.ID)
		}

		tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "device_snap_id"}},
			UpdateAll: true,
		}).Create(ds.ServiceStatuses)
	} else {
		log.Tracef("Snap %s has no services, nothing to update.", ds.Name)
	}

	tx.Commit()
	return nil
}

// DeviceSnapList lists the snaps for a device
func (db *DataStore) DeviceSnapList(deviceID int64) ([]datastore.DeviceSnap, error) {
	snaps := []datastore.DeviceSnap{}

	res := db.gormDB.Where(&datastore.DeviceSnap{DeviceID: deviceID}).Preload(clause.Associations).Find(&snaps)

	if log.IsLevelEnabled(log.TraceLevel) {
		bytes, err := json.Marshal(&snaps)
		if err != nil {
			log.Errorf("While trying to marshal snaps: %s", err.Error())
		} else {
			log.Tracef("%s", string(bytes))
		}
	}

	if res.Error != nil {
		return nil, res.Error
	}

	return snaps, nil
}

// DeviceSnapDelete removes a snap for a device
func (db *DataStore) DeviceSnapDelete(id int64) error {
	tx := db.gormDB.Begin()
	tx.Where(&datastore.DeviceSnap{DeviceID: id}).Delete(&datastore.DeviceSnap{})
	tx.Commit()

	return tx.Error
}
