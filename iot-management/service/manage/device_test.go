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

package manage

import (
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	mocks2 "github.com/everactive/dmscore/iot-devicetwin/service/controller/mocks"
	"github.com/everactive/dmscore/iot-management/datastore"
	mocks3 "github.com/everactive/dmscore/iot-management/datastore/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestManagement_DeviceList(t *testing.T) {
	type args struct {
		orgID    string
		username string
		role     int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr string
	}{
		{"valid", args{"abc", "jamesj", 300}, 3, ""},
		{"invalid-user", args{"abc", "invalid", 200}, 0, "DevicesAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manageDataStore := &mocks3.DataStore{}
			manageDataStore.On("OrganizationsForUser", mock.Anything).Return([]datastore.Organization{}, nil)

			if tt.wantErr == "DevicesAuth" {
				manageDataStore.On("OrgUserAccess", mock.Anything, mock.Anything, mock.Anything).Return(false)
			} else {
				manageDataStore.On("OrgUserAccess", mock.Anything, mock.Anything, mock.Anything).Return(true)
			}

			deviceTwinController := &mocks2.Controller{}

			devices := []messages.Device{}
			if tt.want > 0 {
				for v := 0; v < tt.want; v++ {
					devices = append(devices, messages.Device{})
				}
			}

			deviceTwinController.On("DeviceList", mock.Anything).Return(devices, nil)

			srv := Management{
				DS:                   manageDataStore,
				DB:                   nil,
				TwinAPI:              nil,
				IdentityAPI:          nil,
				DeviceTwinController: deviceTwinController,
				Identity:             nil,
			}

			got := srv.DeviceList(tt.args.orgID, tt.args.username, tt.args.role)
			if got.Code != tt.wantErr {
				t.Errorf("Management.DeviceList() = %v, want %v", got.Code, tt.wantErr)
			}
			if len(got.Devices) != tt.want {
				t.Errorf("Management.DeviceList() = %v, want %v", len(got.Devices), tt.want)
			}
		})
	}
}

func TestManagement_DeviceGet(t *testing.T) {
	type args struct {
		orgID    string
		username string
		role     int
		deviceID string
	}
	tests := []struct {
		name       string
		args       args
		wantSerial string
		wantErr    string
	}{
		{"valid", args{"abc", "jamesj", 200, "b222"}, "DR1000B222", ""},
		{"invalid-user", args{"abc", "invalid", 200, "b222"}, "", "DeviceAuth"},
	}
	for _, tt := range tests {
		manageDataStoreMock := &mocks3.DataStore{}
		deviceTwinControllerMock := &mocks2.Controller{}
		manageDataStoreMock.On("OrganizationsForUser", mock.Anything).Return([]datastore.Organization{}, nil)

		hasAccess := false
		if tt.wantErr == "" {
			hasAccess = true
		}
		manageDataStoreMock.On("OrgUserAccess", mock.Anything, mock.Anything, mock.Anything).Return(hasAccess)

		deviceTwinControllerMock.On("DeviceGet", mock.Anything, mock.Anything).Return(messages.Device{Serial: tt.wantSerial}, nil)

		t.Run(tt.name, func(t *testing.T) {
			srv := Management{
				DS:                   manageDataStoreMock,
				DB:                   nil,
				TwinAPI:              nil,
				IdentityAPI:          nil,
				DeviceTwinController: deviceTwinControllerMock,
				Identity:             nil,
			}

			got := srv.DeviceGet(tt.args.orgID, tt.args.username, tt.args.role, tt.args.deviceID)
			if got.Code != tt.wantErr {
				t.Errorf("Management.DeviceGet() = %v, want %v", got.Code, tt.wantErr)
			}
			if got.Device.Serial != tt.wantSerial {
				t.Errorf("Management.DeviceGet() = %v, want %v", got.Device.Serial, tt.wantSerial)
			}
		})
	}
}

func TestManagement_DeviceLogs(t *testing.T) {
	type args struct {
		orgID    string
		username string
		role     int
		deviceID string
		body     []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{"valid-enable", args{"abc", "jamesj", 300, "a111", []byte("{}")}, ""},
		{"invalid-user", args{"abc", "invalid", 200, "a111", []byte("{}")}, "DeviceAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manageDataStoreMock := &mocks3.DataStore{}
			deviceTwinControllerMock := &mocks2.Controller{}
			manageDataStoreMock.On("OrganizationsForUser", mock.Anything).Return([]datastore.Organization{}, nil)

			hasAccess := false
			if tt.wantErr == "" {
				hasAccess = true
			}
			manageDataStoreMock.On("OrgUserAccess", mock.Anything, mock.Anything, mock.Anything).Return(hasAccess)

			deviceTwinControllerMock.On("DeviceLogs", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			srv := Management{
				DS:                   manageDataStoreMock,
				DB:                   nil,
				TwinAPI:              nil,
				IdentityAPI:          nil,
				DeviceTwinController: deviceTwinControllerMock,
				Identity:             nil,
			}
			got := srv.DeviceLogs(tt.args.orgID, tt.args.username, tt.args.role, tt.args.deviceID, &messages.DeviceLogs{})
			if got.Code != tt.wantErr {
				t.Errorf("Management.DeviceLogs() = %v, want %v", got.Code, tt.wantErr)
			}
		})
	}
}

func TestManagement_DeviceUsersAction(t *testing.T) {
	type args struct {
		orgID    string
		username string
		role     int
		deviceID string
		body     []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{"valid-enable", args{"abc", "jamesj", 300, "a111", []byte("{}")}, ""},
		{"invalid-user", args{"abc", "invalid", 200, "a111", []byte("{}")}, "DeviceAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manageDataStoreMock := &mocks3.DataStore{}
			deviceTwinControllerMock := &mocks2.Controller{}
			manageDataStoreMock.On("OrganizationsForUser", mock.Anything).Return([]datastore.Organization{}, nil)

			hasAccess := false
			if tt.wantErr == "" {
				hasAccess = true
			}
			manageDataStoreMock.On("OrgUserAccess", mock.Anything, mock.Anything, mock.Anything).Return(hasAccess)

			deviceTwinControllerMock.On("User", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			srv := Management{
				DS:                   manageDataStoreMock,
				DB:                   nil,
				TwinAPI:              nil,
				IdentityAPI:          nil,
				DeviceTwinController: deviceTwinControllerMock,
				Identity:             nil,
			}

			got := srv.DeviceUsersAction(tt.args.orgID, tt.args.username, tt.args.role, tt.args.deviceID, messages.DeviceUser{})
			if got.Code != tt.wantErr {
				t.Errorf("Management.DeviceUsersAction() = %v, want %v", got.Code, tt.wantErr)
			}
		})
	}
}
