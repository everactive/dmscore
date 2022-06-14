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
	"encoding/json"
	"github.com/everactive/dmscore/iot-identity/service"
	"github.com/everactive/dmscore/iot-identity/web"
)

// RegDeviceList gets the registered devices a user can access for an organization
func (srv *Management) RegDeviceList(orgID, username string, role int) web.DevicesResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.DevicesResponse{
			StandardResponse: web.StandardResponse{
				Code:    "RegDevicesAuth",
				Message: "the user does not have permissions for the organization",
			},
		}
	}

	return srv.IdentityAPI.RegDeviceList(orgID)
}

// RegisterDevice registers a new device
func (srv *Management) RegisterDevice(orgID, username string, role int, body []byte) web.RegisterResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.RegisterResponse{
			StandardResponse: web.StandardResponse{
				Code:    "RegDeviceAuth",
				Message: "the user does not have permissions for the organization",
			},
		}
	}

	request := service.RegisterDeviceRequest{}
	err := json.Unmarshal(body, &request)
	if err != nil {
		return web.RegisterResponse{
			StandardResponse: web.StandardResponse{
				Code:    "RegDevice",
				Message: err.Error(),
			},
		}
	}

	// return srv.IdentityAPI.RegisterDevice(body)
	device, err := srv.IdentityService.Identity.RegisterDevice(&request)
	if err != nil {
		return web.RegisterResponse{
			StandardResponse: web.StandardResponse{
				Code:    "RegDevice",
				Message: err.Error(),
			},
		}
	}

	return web.RegisterResponse{
		ID: device,
	}
}

// RegDeviceGet fetches a device registration
func (srv *Management) RegDeviceGet(orgID, username string, role int, deviceID string) web.EnrollResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.EnrollResponse{
			StandardResponse: web.StandardResponse{
				Code:    "RegDeviceAuth",
				Message: "the user does not have permissions for the organization",
			},
		}
	}

	// return srv.IdentityAPI.RegDeviceGet(orgID, deviceID)
	enrollment, err := srv.IdentityService.Identity.DeviceGet(orgID, deviceID)
	if err != nil {
		return web.EnrollResponse{
			StandardResponse: web.StandardResponse{
				Code:    "RegDeviceGet",
				Message: err.Error(),
			},
		}
	}

	return web.EnrollResponse{
		Enrollment: *enrollment,
	}
}

// RegDeviceUpdate updates a device registration
func (srv *Management) RegDeviceUpdate(orgID, username string, role int, deviceID string, body []byte) web.StandardResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "RegDeviceAuth",
			Message: "the user does not have permissions for the organization",
		}
	}
	// return srv.IdentityAPI.RegDeviceUpdate(orgID, deviceID, body)
	request := service.DeviceUpdateRequest{}
	err := json.Unmarshal(body, &request)
	if err != nil {
		return web.StandardResponse{
			Code:    "RegDeviceUpdate",
			Message: err.Error(),
		}
	}

	err = srv.IdentityService.Identity.DeviceUpdate(orgID, deviceID, &request)
	if err != nil {
		return web.StandardResponse{
			Code:    "RegDeviceUpdate",
			Message: err.Error(),
		}
	}

	return web.StandardResponse{}
}
