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
	"log"

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
	var id int64
	err := db.QueryRow(createActionSQL, act.OrganizationID, act.DeviceID, act.ActionID, act.Action, act.Status, act.Message).Scan(&id)
	if err != nil {
		log.Printf("Error creating action %s/%s: %v\n", act.DeviceID, act.ActionID, err)
	}

	return id, err
}

// ActionUpdate updates an action record
func (db *DataStore) ActionUpdate(actionID, status, message string) error {
	_, err := db.Exec(updateActionSQL, actionID, status, message)
	if err != nil {
		log.Printf("Error updating the device: %v\n", err)
	}

	return err
}

// ActionListForDevice lists the actions for a device
func (db *DataStore) ActionListForDevice(orgID, deviceID string) ([]datastore.Action, error) {
	newDeviceID, err := getDeviceIDIfSerial(db, deviceID)
	if err == nil {
		deviceID = newDeviceID
	}

	rows, err := db.Query(listActionSQL, orgID, deviceID)
	if err != nil {
		log.Printf("Error retrieving actions: %v\n", err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Error attempting to close rows: %+v", err)
		}
	}()

	actions := []datastore.Action{}
	for rows.Next() {
		item := datastore.Action{}
		err := rows.Scan(&item.ID, &item.Created, &item.Modified, &item.OrganizationID, &item.DeviceID, &item.ActionID, &item.Action, &item.Status, &item.Message)
		if err != nil {
			return nil, err
		}
		actions = append(actions, item)
	}

	return actions, nil
}
