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
	"errors"
	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
)

func getDeviceIDIfSerial(db *DataStore, deviceID string) (string, error) {
	device := datastore.Device{}
	res := db.gormDB.Where("serial = ?", deviceID).Find(&device)
	if res.RowsAffected == 1 {
		return device.DeviceID, nil
	}
	return "", errors.New("device not found")
}

// ActionCreate log an new action
func (db *DataStore) ActionCreate(act datastore.Action) (int64, error) {
	res := db.gormDB.Create(&act)
	if res.Error != nil {
		log.Error(res.Error)
		return 0, res.Error
	}

	return int64(act.ID), nil
}

// ActionUpdate updates an action record
func (db *DataStore) ActionUpdate(actionID, status, message string) error {
	res := db.gormDB.Model(&datastore.Action{}).
		Where("action_id = ?", actionID).
		Updates(&datastore.Action{Status: status, Message: message})

	if res.Error != nil {
		log.Error(res.Error)
		return res.Error
	}

	return nil
}

// ActionListForDevice lists the actions for a device
func (db *DataStore) ActionListForDevice(orgID, deviceID string) ([]datastore.Action, error) {
	newDeviceID, err := getDeviceIDIfSerial(db, deviceID)
	if err == nil {
		deviceID = newDeviceID
	}

	actions := []datastore.Action{}
	res := db.gormDB.Where("org_id = ? AND device_id = ?", orgID, deviceID).Order("updated_at").Find(&actions)
	if res.Error != nil {
		log.Error(res.Error)
		return actions, res.Error
	}

	return actions, nil
}
