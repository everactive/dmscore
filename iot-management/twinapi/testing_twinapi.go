// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Management Service
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

package twinapi

import (
	"github.com/everactive/dmscore/iot-devicetwin/domain"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	"github.com/everactive/dmscore/iot-devicetwin/web"
)

const (
	testInstalledSize = 3000
)

//nolint
type MockClient struct{}

//nolint
func (m *MockClient) DeviceList(orgID string) web.DevicesResponse {
	return web.DevicesResponse{
		StandardResponse: web.StandardResponse{},
		Devices: []messages.Device{
			{OrgId: "abc", DeviceId: "a111", Brand: "example", Model: "drone-1000", Serial: "DR1000A111", DeviceKey: "AAAAAAAAA", Store: "example-store"},
			{OrgId: "abc", DeviceId: "b222", Brand: "example", Model: "drone-1000", Serial: "DR1000B222", DeviceKey: "BBBBBBBBB", Store: "example-store"},
			{OrgId: "abc", DeviceId: "c333", Brand: "canonical", Model: "ubuntu-core-18-amd64", Serial: "d75f7300-abbf-4c11-bf0a-8b7103038490", DeviceKey: "CCCCCCCCC"},
		},
	}
}

//nolint
func (m *MockClient) DeviceGet(orgID, deviceID string) web.DeviceResponse {
	return web.DeviceResponse{
		StandardResponse: web.StandardResponse{},
		Device:           messages.Device{OrgId: "abc", DeviceId: "b222", Brand: "example", Model: "drone-1000", Serial: "DR1000B222", DeviceKey: "BBBBBBBBB", Store: "example-store"},
	}
}

//nolint
func (m *MockClient) ActionList(orgID, deviceID string) web.ActionsResponse {
	return web.ActionsResponse{
		StandardResponse: web.StandardResponse{},
		Actions: []domain.Action{
			{OrganizationID: "abc", DeviceID: "b222", Action: "list", Status: "complete"},
		},
	}
}

//nolint
func (m *MockClient) SnapList(orgID, deviceID string) web.SnapsResponse {
	return web.SnapsResponse{
		StandardResponse: web.StandardResponse{},
		Snaps:            []messages.DeviceSnap{{Name: "example-snap", InstalledSize: testInstalledSize, Status: "active"}},
	}
}

//nolint
func (m *MockClient) SnapListOnDevice(orgID, deviceID string) web.StandardResponse {
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) SnapInstall(orgID, deviceID, snap string) web.StandardResponse {
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) SnapRemove(orgID, deviceID, snap string) web.StandardResponse {
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) SnapUpdate(orgID, deviceID, snap, action string, body []byte) web.StandardResponse {
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) SnapConfigSet(orgID, deviceID, snap string, config []byte) web.StandardResponse {
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) GroupList(orgID string) web.GroupsResponse {
	return web.GroupsResponse{
		StandardResponse: web.StandardResponse{},
		Groups:           []domain.Group{{OrganizationID: "abc", Name: "workshop"}},
	}
}

//nolint
func (m *MockClient) GroupDevices(orgID, name string) web.DevicesResponse {
	return web.DevicesResponse{
		StandardResponse: web.StandardResponse{},
		Devices: []messages.Device{
			{OrgId: "abc", DeviceId: "a111", Brand: "example", Model: "drone-1000", Serial: "DR1000A111", DeviceKey: "AAAAAAAAA", Store: "example-store"},
		},
	}
}

//nolint
func (m *MockClient) GroupExcludedDevices(orgID, name string) web.DevicesResponse {
	return web.DevicesResponse{
		StandardResponse: web.StandardResponse{},
		Devices: []messages.Device{
			{OrgId: "abc", DeviceId: "b222", Brand: "example", Model: "drone-1000", Serial: "DR1000B222", DeviceKey: "BBBBBBBBB", Store: "example-store"},
			{OrgId: "abc", DeviceId: "c333", Brand: "canonical", Model: "ubuntu-core-18-amd64", Serial: "d75f7300-abbf-4c11-bf0a-8b7103038490", DeviceKey: "CCCCCCCCC"},
		},
	}
}

//nolint
func (m *MockClient) GroupCreate(orgID string, body []byte) web.StandardResponse {
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) GroupDeviceLink(orgID, name, deviceID string) web.StandardResponse {
	if orgID == "invalid" || deviceID == "invalid" {
		return web.StandardResponse{Code: "GroupDevice", Message: "MOCK error link"}
	}
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) GroupDeviceUnlink(orgID, name, deviceID string) web.StandardResponse {
	if orgID == "invalid" || deviceID == "invalid" {
		return web.StandardResponse{Code: "GroupDevice", Message: "MOCK error unlink"}
	}
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) SnapSnapshot(orgID, deviceID, snap string, body []byte) web.StandardResponse {
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) DeviceLogs(orgID, deviceID string, body []byte) web.StandardResponse {
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) DeviceUsersAction(orgID, deviceID string, body []byte) web.StandardResponse {
	return web.StandardResponse{}
}

//nolint
func (m *MockClient) DeviceDelete(orgID, deviceID string) web.StandardResponse {
	panic("implement me")
}

//nolint
func (m *MockClient) SnapServiceAction(orgID, deviceID, snap, action string, body []byte) web.StandardResponse {
	panic("implement me")
}
