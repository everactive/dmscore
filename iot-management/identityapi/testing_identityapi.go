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

package identityapi

import (
	"github.com/everactive/dmscore/iot-identity/domain"
	"github.com/everactive/dmscore/iot-identity/web"
)

//nolint
type MockIdentity struct{}

//nolint
func (m *MockIdentity) RegDeviceList(orgID string) web.DevicesResponse {
	if orgID == "invalid" {
		return web.DevicesResponse{
			StandardResponse: web.StandardResponse{Code: "RegDeviceAuth", Message: "MOCK error devices"},
			Devices:          nil,
		}
	}
	return web.DevicesResponse{
		StandardResponse: web.StandardResponse{},
		Devices:          []domain.Enrollment{},
	}
}

//nolint
func (m *MockIdentity) RegisterDevice(body []byte) web.RegisterResponse {
	return web.RegisterResponse{
		StandardResponse: web.StandardResponse{},
		ID:               "d444",
	}
}

//nolint
func (m *MockIdentity) RegDeviceGet(orgID, deviceID string) web.EnrollResponse {
	if deviceID == "invalid" {
		return web.EnrollResponse{
			StandardResponse: web.StandardResponse{Code: "RegDeviceAuth", Message: "MOCK error get"},
			Enrollment:       domain.Enrollment{},
		}
	}
	return web.EnrollResponse{
		StandardResponse: web.StandardResponse{},
		Enrollment:       domain.Enrollment{},
	}
}

//nolint
func (m *MockIdentity) RegDeviceUpdate(orgID, deviceID string, body []byte) web.StandardResponse {
	if deviceID == "invalid" {
		return web.StandardResponse{Code: "RegDeviceUpdate", Message: "MOCK error update"}
	}
	return web.StandardResponse{}
}

//nolint
func (m *MockIdentity) RegOrganizationList() web.OrganizationsResponse {
	return web.OrganizationsResponse{
		StandardResponse: web.StandardResponse{},
		Organizations:    []domain.Organization{},
	}
}

//nolint
func (m *MockIdentity) RegisterOrganization(body []byte) web.RegisterResponse {
	return web.RegisterResponse{
		StandardResponse: web.StandardResponse{},
		ID:               "def",
	}
}

//nolint
func (m *MockIdentity) DeviceDelete(orgID string) web.StandardResponse {
	return web.StandardResponse{}
}
