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

// Package devicetwin is the ManualMockDeviceTwin and associated pieces for testing
package devicetwin

import (
	"fmt"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"

	"github.com/everactive/dmscore/iot-devicetwin/domain"
)

const (
	mockInstalledSize     = 2000
	invalidDeviceIDString = "invalid"
)

// ManualMockDeviceTwin mocks a device twin service
type ManualMockDeviceTwin struct {
	Actions                 []string
	ReturnSoftDeletedDevice bool
}

// Unscoped gets an instance of the DeviceTwin that is able to access (soft) deleted data
func (twin *ManualMockDeviceTwin) Unscoped() UnscopedDeviceTwin {
	return twin
}

// DeviceGetByID gets a device just by its id without org
func (twin *ManualMockDeviceTwin) DeviceGetByID(clientID string) (*messages.Device, bool, error) {
	if twin.ReturnSoftDeletedDevice {
		d, err := twin.DeviceGet("", clientID)
		return &d, true, err
	}

	return nil, false, nil
}

// HealthHandler mocks the health handler
func (twin *ManualMockDeviceTwin) HealthHandler(payload messages.Health) error {
	if payload.DeviceId == invalidDeviceIDString || payload.DeviceId == "new-device" {
		return fmt.Errorf("MOCK error in health handler")
	}
	return nil
}

// ActionResponse mocks the action handler
func (twin *ManualMockDeviceTwin) ActionResponse(clientID, actionID, action string, payload []byte) error {
	if action == invalidDeviceIDString {
		return fmt.Errorf("MOCK error in action")
	}
	return nil
}

// DeviceSnaps mocks the snap list
func (twin *ManualMockDeviceTwin) DeviceSnaps(orgID, clientID string) ([]messages.DeviceSnap, error) {
	if clientID == invalidDeviceIDString {
		return nil, fmt.Errorf("MOCK snaps list")
	}
	return []messages.DeviceSnap{
		{Name: "example-snap", InstalledSize: mockInstalledSize, Status: "active"},
	}, nil
}

// ActionCreate mocks the action log creation
func (twin *ManualMockDeviceTwin) ActionCreate(orgID, deviceID string, act messages.SubscribeAction) error {
	if deviceID == invalidDeviceIDString {
		return fmt.Errorf("MOCK action log create")
	}
	if twin.Actions == nil {
		twin.Actions = []string{}
	}
	twin.Actions = append(twin.Actions, act.Id)
	return nil
}

// ActionUpdate mocks the action log update
func (twin *ManualMockDeviceTwin) ActionUpdate(actionID, status, message string) error {
	return nil
}

// ActionList mocks the action log list
func (twin *ManualMockDeviceTwin) ActionList(orgID, clientID string) ([]domain.Action, error) {
	if clientID == invalidDeviceIDString {
		return nil, fmt.Errorf("MOCK error action list")
	}
	return []domain.Action{}, nil
}

// DeviceGet mocks fetching a device
func (twin *ManualMockDeviceTwin) DeviceGet(orgID, clientID string) (messages.Device, error) {
	if clientID == invalidDeviceIDString {
		return messages.Device{}, fmt.Errorf("MOCK error device get")
	}
	return messages.Device{
		OrgId:     "abc",
		DeviceId:  "c333",
		Brand:     "canonical",
		Model:     "ubuntu-core-18-amd64",
		Serial:    "d75f7300-abbf-4c11-bf0a-8b7103038490",
		DeviceKey: "CCCCCCCCC",
	}, nil
}

// DeviceList mocks fetching devices for an organization
func (twin *ManualMockDeviceTwin) DeviceList(orgID string) ([]messages.Device, error) {
	if orgID == invalidDeviceIDString {
		return nil, fmt.Errorf("MOCK error device list")
	}

	return []messages.Device{
		{OrgId: "abc",
			DeviceId:  "c333",
			Brand:     "canonical",
			Model:     "ubuntu-core-18-amd64",
			Serial:    "d75f7300-abbf-4c11-bf0a-8b7103038490",
			DeviceKey: "CCCCCCCCC",
		},
	}, nil
}

// DeviceDelete mocks deleting a device
func (twin *ManualMockDeviceTwin) DeviceDelete(deviceID string) (string, error) {
	return "c333", nil
}

// GroupCreate mocks creating a group
func (twin *ManualMockDeviceTwin) GroupCreate(orgID, name string) error {
	if orgID == invalidDeviceIDString {
		return fmt.Errorf("MOCK error group create")
	}
	return nil
}

// GroupList mocks listing groups
func (twin *ManualMockDeviceTwin) GroupList(orgID string) ([]domain.Group, error) {
	if orgID == invalidDeviceIDString {
		return nil, fmt.Errorf("MOCK error group list")
	}
	return []domain.Group{
		{OrganizationID: "abc", Name: "workshop"},
	}, nil
}

// GroupGet mocks fetching a group
func (twin *ManualMockDeviceTwin) GroupGet(orgID, name string) (domain.Group, error) {
	if orgID == invalidDeviceIDString || name == invalidDeviceIDString {
		return domain.Group{}, fmt.Errorf("MOCK error group device unlink")
	}
	return domain.Group{
		OrganizationID: "abc", Name: "workshop",
	}, nil
}

// GroupLinkDevice mocks linking a device to a group
func (twin *ManualMockDeviceTwin) GroupLinkDevice(orgID, name, clientID string) error {
	if orgID == invalidDeviceIDString || name == invalidDeviceIDString || clientID == invalidDeviceIDString {
		return fmt.Errorf("MOCK error group device link")
	}
	return nil
}

// GroupUnlinkDevice mocks unlinking a device from a group
func (twin *ManualMockDeviceTwin) GroupUnlinkDevice(orgID, name, clientID string) error {
	if orgID == invalidDeviceIDString || name == invalidDeviceIDString || clientID == invalidDeviceIDString {
		return fmt.Errorf("MOCK error group device unlink")
	}
	return nil
}

// GroupGetDevices mocks retrieving the devices for a group
func (twin *ManualMockDeviceTwin) GroupGetDevices(orgID, name string) ([]messages.Device, error) {
	if orgID == invalidDeviceIDString || name == invalidDeviceIDString {
		return nil, fmt.Errorf("MOCK error group devices")
	}
	return []messages.Device{
		{OrgId: "abc",
			DeviceId:  "c333",
			Brand:     "canonical",
			Model:     "ubuntu-core-18-amd64",
			Serial:    "d75f7300-abbf-4c11-bf0a-8b7103038490",
			DeviceKey: "CCCCCCCCC",
		},
	}, nil
}

// GroupGetExcludedDevices mocks retrieving the devices not in a group
func (twin *ManualMockDeviceTwin) GroupGetExcludedDevices(orgID, name string) ([]messages.Device, error) {
	if orgID == invalidDeviceIDString || name == invalidDeviceIDString {
		return nil, fmt.Errorf("MOCK error group excluded devices")
	}
	return []messages.Device{
		{OrgId: "abc",
			DeviceId:  "b222",
			Brand:     "example",
			Model:     "drone-1000",
			Serial:    "DR1000B222",
			DeviceKey: "BBBBBBBBB",
		},
		{OrgId: "abc",
			DeviceId:  "c333",
			Brand:     "canonical",
			Model:     "ubuntu-core-18-amd64",
			Serial:    "d75f7300-abbf-4c11-bf0a-8b7103038490",
			DeviceKey: "CCCCCCCCC",
		},
	}, nil
}
