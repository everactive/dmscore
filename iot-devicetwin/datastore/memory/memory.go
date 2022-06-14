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

// Package memory is the Datastore implementation for in-memory
package memory

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
)

// Store implements an in-memory store for testing
type Store struct {
	Devices        []datastore.Device
	Snaps          []datastore.DeviceSnap
	Actions        []datastore.Action
	DeviceVersions []datastore.DeviceVersion
	Groups         []datastore.Group
	GroupLinks     []datastore.GroupDeviceLink
	lock           sync.RWMutex
}

const (
	testID1           = 1
	testID2           = 2
	testID3           = 3
	testInstalledSize = 2000
	invalidString     = "invalid"
)

// NewStore creates a new memory store
func NewStore() *Store {
	d1 := datastore.Device{Model: gorm.Model{ID: testID1}, OrganisationID: "abc", DeviceID: "a111", Brand: "example", DeviceModel: "drone-1000", SerialNumber: "DR1000A111", DeviceKey: "AAAAAAAAA", StoreID: "example-store", Active: true}
	d2 := datastore.Device{Model: gorm.Model{ID: testID2}, OrganisationID: "abc", DeviceID: "b222", Brand: "example", DeviceModel: "drone-1000", SerialNumber: "DR1000B222", DeviceKey: "BBBBBBBBB", StoreID: "example-store", Active: true}
	d3 := datastore.Device{Model: gorm.Model{ID: testID3}, OrganisationID: "abc", DeviceID: "c333", Brand: "canonical", DeviceModel: "ubuntu-core-18-amd64", SerialNumber: "d75f7300-abbf-4c11-bf0a-8b7103038490", DeviceKey: "CCCCCCCCC", Active: true}

	return &Store{
		Devices: []datastore.Device{d1, d2, d3},
		Snaps: []datastore.DeviceSnap{
			{DeviceID: testID1, Name: "example-snap", InstalledSize: testInstalledSize, Status: "active"},
		},
		Actions: []datastore.Action{
			{ID: testID1, OrganizationID: "abc", DeviceID: "c333", Action: "list", Status: ""},
			{ID: testID2, OrganizationID: "abc", DeviceID: "c333", Action: "list", Status: ""},
		},
		DeviceVersions: []datastore.DeviceVersion{
			{Model: gorm.Model{ID: testID1}, DeviceID: testID3, KernelVersion: "kernel-123", OSVersionID: "core-123", Series: "16"},
		},
		Groups:     []datastore.Group{{ID: testID1, OrganisationID: "abc", Name: "workshop"}},
		GroupLinks: []datastore.GroupDeviceLink{{ID: testID1, OrganisationID: "abc", GroupID: testID1, DeviceID: testID1}},
	}
}

// Unscoped get an unscoped instance
func (mem *Store) Unscoped() datastore.UnscopedDataStore {
	return mem
}

// DeviceDelete is no-op
func (mem *Store) DeviceDelete(_ string) error {
	return nil
}

// DeviceList fetches existing devices
func (mem *Store) DeviceList(orgID string) ([]datastore.Device, error) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	if orgID == invalidString {
		return nil, fmt.Errorf("MOCK list error")
	}

	devices := []datastore.Device{}

	for _, d := range mem.Devices {
		if d.OrganisationID == orgID {
			devices = append(devices, d)
		}
	}
	return devices, nil
}

// DeviceGet fetches an existing device
func (mem *Store) DeviceGet(id string) (datastore.Device, error) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	for _, d := range mem.Devices {
		if d.DeviceID == id {
			return d, nil
		}
	}
	return datastore.Device{}, fmt.Errorf("device with ID `%s` not found", id)
}

// deviceGetByID fetches an existing device
func (mem *Store) deviceGetByID(id int64) (datastore.Device, error) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	for _, d := range mem.Devices {
		if d.ID == uint(id) {
			return d, nil
		}
	}
	return datastore.Device{}, fmt.Errorf("device with ID `%d` not found", id)
}

// DevicePing updates a device to indicate its health
func (mem *Store) DevicePing(id string, refresh time.Time) error {
	device, err := mem.DeviceGet(id)
	if err != nil {
		return err
	}

	mem.lock.Lock()
	defer mem.lock.Unlock()
	device.LastRefresh = refresh

	for i := range mem.Devices {
		if mem.Devices[i].DeviceID == id {
			mem.Devices[i] = device
		}
	}
	return nil
}

// DeviceCreate creates a new device
func (mem *Store) DeviceCreate(device datastore.Device) (int64, error) {
	// Check the device does not exist
	if _, err := mem.DeviceGet(device.DeviceID); err == nil {
		return 0, fmt.Errorf("device with ID `%s` already exists", device.DeviceID)
	}

	mem.lock.Lock()
	defer mem.lock.Unlock()

	device.LastRefresh = time.Now()

	device.ID = uint(int64(len(mem.Devices) + 1))
	mem.Devices = append(mem.Devices, device)
	return int64(device.ID), nil
}

// DeviceSnapUpsert creates or updates a snap for a device
func (mem *Store) DeviceSnapUpsert(ds datastore.DeviceSnap) error {
	mem.lock.Lock()
	defer mem.lock.Unlock()
	// Find the snap
	found := -1
	for i, s := range mem.Snaps {
		if s.Name == ds.Name {
			found = i
		}
	}

	if found < 0 {
		// Not found, so create it
		ds.CreatedAt = time.Now()
		ds.UpdatedAt = time.Now()
		mem.Snaps = append(mem.Snaps, ds)
		return nil
	}

	// Update the existing record
	ds.UpdatedAt = time.Now()
	mem.Snaps[found] = ds
	return nil
}

// DeviceSnapList lists the snaps for a device
func (mem *Store) DeviceSnapList(id int64) ([]datastore.DeviceSnap, error) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()
	snaps := []datastore.DeviceSnap{}

	for _, s := range mem.Snaps {
		if s.DeviceID == id {
			snaps = append(snaps, s)
		}
	}
	return snaps, nil
}

// DeviceSnapDelete deletes the snap records for a device
func (mem *Store) DeviceSnapDelete(id int64) error {
	mem.lock.Lock()
	defer mem.lock.Unlock()
	snaps := []datastore.DeviceSnap{}

	for _, s := range mem.Snaps {
		if s.DeviceID != id {
			snaps = append(snaps, s)
		}
	}
	mem.Snaps = snaps

	return nil
}

// ActionCreate creates an action log
func (mem *Store) ActionCreate(act datastore.Action) (int64, error) {
	mem.lock.Lock()
	defer mem.lock.Unlock()

	act.ID = int64(len(mem.Actions) + 1)
	mem.Actions = append(mem.Actions, act)
	return act.ID, nil
}

// ActionUpdate updates an action log
func (mem *Store) ActionUpdate(actionID, status, message string) error {
	mem.lock.Lock()
	defer mem.lock.Unlock()

	actions := []datastore.Action{}
	for _, a := range mem.Actions {
		if a.ActionID == actionID {
			a.Status = status
			a.Message = message
		}
		a.Modified = time.Now()
		actions = append(actions, a)
	}
	mem.Actions = actions
	return nil
}

// ActionListForDevice fetches the actions for a device
func (mem *Store) ActionListForDevice(orgID, clientID string) ([]datastore.Action, error) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	actions := []datastore.Action{}
	for _, a := range mem.Actions {
		if a.DeviceID == clientID {
			actions = append(actions, a)
		}
	}

	return actions, nil
}

// DeviceVersionGet gets the OS details for a device
func (mem *Store) DeviceVersionGet(deviceID int64) (datastore.DeviceVersion, error) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	for _, d := range mem.DeviceVersions {
		if d.DeviceID == deviceID {
			return d, nil
		}
	}
	return datastore.DeviceVersion{}, fmt.Errorf("device version with device ID `%d` not found", deviceID)
}

// DeviceVersionUpsert creates or updates the device OS details
func (mem *Store) DeviceVersionUpsert(dv datastore.DeviceVersion) error {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	// Find the record
	found := -1
	for i, v := range mem.DeviceVersions {
		if v.DeviceID == dv.DeviceID {
			found = i
		}
	}

	if found < 0 {
		// Not found, so create it
		dv.ID = uint(int64(len(mem.DeviceVersions) + 1))
		mem.DeviceVersions = append(mem.DeviceVersions, dv)
		return nil
	}

	// Update the existing record
	mem.DeviceVersions[found] = dv
	return nil
}

// DeviceVersionDelete removes a OS record
func (mem *Store) DeviceVersionDelete(id int64) error {
	mem.lock.Lock()
	defer mem.lock.Unlock()
	versions := []datastore.DeviceVersion{}

	found := false
	for _, s := range mem.DeviceVersions {
		if s.ID != uint(id) {
			versions = append(versions, s)
		} else {
			found = true
		}
	}
	mem.DeviceVersions = versions

	if !found {
		return fmt.Errorf("cannot find record with ID %d", id)
	}
	return nil
}

// GroupCreate creates a group record
func (mem *Store) GroupCreate(orgID, name string) (int64, error) {
	mem.lock.Lock()
	defer mem.lock.Unlock()

	if orgID == invalidString {
		return 0, fmt.Errorf("error cannot find organization `%s`", orgID)
	}

	for _, g := range mem.Groups {
		if g.OrganisationID == orgID && g.Name == name {
			return 0, fmt.Errorf("group `%s1` already exists for organization `%s`", name, orgID)
		}
	}

	g := datastore.Group{
		ID:             int64(len(mem.Groups) + 1),
		OrganisationID: orgID,
		Name:           name,
		Created:        time.Now(),
		Modified:       time.Now(),
	}
	mem.Groups = append(mem.Groups, g)
	return g.ID, nil
}

// GroupList lists groups for an organization
func (mem *Store) GroupList(orgID string) ([]datastore.Group, error) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	if orgID == invalidString {
		return nil, fmt.Errorf("error cannot find organization `%s`", orgID)
	}

	groups := []datastore.Group{}
	for _, g := range mem.Groups {
		if g.OrganisationID == orgID {
			groups = append(groups, g)
		}
	}

	return groups, nil
}

// GroupGet fetches a group record
func (mem *Store) GroupGet(orgID, name string) (datastore.Group, error) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	for _, g := range mem.Groups {
		if g.OrganisationID == orgID && g.Name == name {
			return g, nil
		}
	}

	return datastore.Group{}, fmt.Errorf("error cannot find group `%s`", name)
}

// GroupLinkDevice links a device with a group
func (mem *Store) GroupLinkDevice(orgID, name, clientID string) error {
	device, err := mem.DeviceGet(clientID)
	if err != nil {
		return err
	}

	group, err := mem.GroupGet(orgID, name)
	if err != nil {
		return err
	}

	mem.lock.Lock()
	defer mem.lock.Unlock()

	for _, l := range mem.GroupLinks {
		if l.GroupID == group.ID && l.ID == int64(device.ID) {
			// Link already exists, so no more work needed
			return nil
		}
	}

	link := datastore.GroupDeviceLink{
		ID:             int64(len(mem.GroupLinks) + 1),
		OrganisationID: orgID,
		GroupID:        group.ID,
		DeviceID:       int64(device.ID),
	}
	mem.GroupLinks = append(mem.GroupLinks, link)
	return nil
}

// GroupUnlinkDevice unlinks a device from a group
func (mem *Store) GroupUnlinkDevice(orgID, name, clientID string) error {
	device, err := mem.DeviceGet(clientID)
	if err != nil {
		return err
	}

	group, err := mem.GroupGet(orgID, name)
	if err != nil {
		return err
	}

	mem.lock.Lock()
	defer mem.lock.Unlock()

	links := []datastore.GroupDeviceLink{}
	for _, l := range mem.GroupLinks {
		if l.GroupID != group.ID || l.DeviceID != int64(device.ID) {
			links = append(links, l)
		}
	}

	mem.GroupLinks = links
	return nil
}

// GroupGetDevices fetches the devices for a group
func (mem *Store) GroupGetDevices(orgID, name string) ([]datastore.Device, error) {
	group, err := mem.GroupGet(orgID, name)
	if err != nil {
		return nil, err
	}

	mem.lock.RLock()
	defer mem.lock.RUnlock()

	devices := []datastore.Device{}

	for _, l := range mem.GroupLinks {
		if l.GroupID != group.ID {
			continue
		}

		d, err := mem.deviceGetByID(l.DeviceID)
		if err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

// GroupGetExcludedDevices fetches the devices not in a group
func (mem *Store) GroupGetExcludedDevices(orgID, name string) ([]datastore.Device, error) {
	devicesInGroup, err := mem.GroupGetDevices(orgID, name)
	if err != nil {
		return nil, err
	}
	devices := []datastore.Device{}

	mem.lock.RLock()
	defer mem.lock.RUnlock()

	for _, l := range mem.Devices {
		found := false
		for _, d := range devicesInGroup {
			if l.ID == d.ID {
				found = true
				break
			}
		}
		if !found {
			devices = append(devices, l)
		}
	}
	return devices, nil
}
