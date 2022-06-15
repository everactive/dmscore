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
	"fmt"
	"log"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
)

// GroupCreate creates a new group for an organization
func (db *DataStore) GroupCreate(orgID, name string) (int64, error) {
	var id int64
	err := db.QueryRow(createOrgGroupSQL, orgID, name).Scan(&id)
	if err != nil {
		log.Printf("Error creating group %s/%s: %v\n", orgID, name, err)
	}

	return id, err
}

// GroupList lists the groups for an organization
func (db *DataStore) GroupList(orgID string) ([]datastore.Group, error) {
	rows, err := db.Query(listOrgGroupSQL, orgID)
	if err != nil {
		log.Printf("Error retrieving groups: %v\n", err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Error attempting to close rows: %+v", err)
		}
	}()

	groups := []datastore.Group{}
	for rows.Next() {
		item := datastore.Group{}
		err := rows.Scan(&item.ID, &item.Created, &item.Modified, &item.OrganisationID, &item.Name)
		if err != nil {
			return nil, err
		}
		groups = append(groups, item)
	}

	return groups, nil
}

// GroupGet fetches a group
func (db *DataStore) GroupGet(orgID, name string) (datastore.Group, error) {
	item := datastore.Group{}
	row := db.QueryRow(getOrgGroupSQL, orgID, name)
	err := row.Scan(&item.ID, &item.Created, &item.Modified, &item.OrganisationID, &item.Name)
	if err != nil {
		log.Printf("Error retrieving group `%s`: %v\n", name, err)
	}
	return item, err
}

// GroupLinkDevice links a device to a group
func (db *DataStore) GroupLinkDevice(orgID, name, deviceID string) error {
	// Get the group record
	grp, err := db.GroupGet(orgID, name)
	if err != nil {
		return fmt.Errorf("error finding group: %v", err)
	}

	// Get the device
	device, err := db.DeviceGet(deviceID)
	if err != nil {
		return fmt.Errorf("error finding device: %v", err)
	}

	// Create the group link record
	_, err = db.Exec(createGroupDeviceLinkSQL, orgID, grp.ID, device.ID)
	return err
}

// GroupUnlinkDevice unlinks a device from a group
func (db *DataStore) GroupUnlinkDevice(orgID, name, deviceID string) error {
	// Get the group record
	grp, err := db.GroupGet(orgID, name)
	if err != nil {
		return fmt.Errorf("error finding group: %v", err)
	}

	// Get the device
	device, err := db.DeviceGet(deviceID)
	if err != nil {
		return fmt.Errorf("error finding device: %v", err)
	}

	// Delete the group link record
	_, err = db.Exec(deleteGroupDeviceLinkSQL, grp.ID, device.ID)
	return err
}

func (db *DataStore) getGroupDevices(sqlString, orgID, name string) ([]datastore.Device, error) {
	// Get the group record
	grp, err := db.GroupGet(orgID, name)
	if err != nil {
		return nil, fmt.Errorf("error finding group: %v", err)
	}

	// Get the devices for the group
	rows, err := db.Query(sqlString, orgID, grp.ID)
	if err != nil {
		log.Printf("Error retrieving devices for group: %v\n", err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Error attempting to close rows: %+v", err)
		}
	}()

	devices := []datastore.Device{}
	for rows.Next() {
		item := datastore.Device{}
		err := rows.Scan(&item.ID, &item.CreatedAt, &item.LastRefresh, &item.OrganisationID, &item.DeviceID, &item.Brand, &item.Model, &item.SerialNumber, &item.StoreID, &item.DeviceKey, &item.Active)
		if err != nil {
			return nil, err
		}
		devices = append(devices, item)
	}

	return devices, nil
}

// GroupGetDevices retrieves the devices for a group
func (db *DataStore) GroupGetDevices(orgID, name string) ([]datastore.Device, error) {
	return db.getGroupDevices(listGroupDeviceLinkSQL, orgID, name)
}

// GroupGetExcludedDevices retrieves the devices not in a group
func (db *DataStore) GroupGetExcludedDevices(orgID, name string) ([]datastore.Device, error) {
	return db.getGroupDevices(listGroupDeviceExcludedLinkSQL, orgID, name)
}
