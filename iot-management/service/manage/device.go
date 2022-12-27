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
	"fmt"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	"github.com/everactive/dmscore/iot-devicetwin/web"
)

func getUserOrgIDIfOrgName(srv *Management, username, orgID string) (string, error) {
	userOrganizations, err := srv.DS.OrganizationsForUser(username)
	if err != nil {
		return "", err
	}

	for _, o := range userOrganizations {
		if o.Name == orgID {
			// If the organization ID supplied matches a name, then use the found ID
			return o.OrganizationID, nil
		}
	}

	return orgID, nil
}

// DeviceList gets the devices a user can access for an organization
func (srv *Management) DeviceList(orgID, username string, role int) web.DevicesResponse {
	newOrgID, err := getUserOrgIDIfOrgName(srv, username, orgID)
	if err != nil {
		return web.DevicesResponse{
			StandardResponse: web.StandardResponse{
				Code:    "Error",
				Message: err.Error(),
			},
		}
	}

	orgID = newOrgID

	hasAccess := srv.DS.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.DevicesResponse{
			StandardResponse: web.StandardResponse{
				Code:    "DevicesAuth",
				Message: "the user does not have permissions for the organization",
			},
		}
	}

	deviceList, err := srv.DeviceTwinController.DeviceList(orgID)
	if err != nil {
		return web.DevicesResponse{
			StandardResponse: web.StandardResponse{
				Code:    "Error",
				Message: err.Error(),
			},
		}
	}

	return web.DevicesResponse{
		StandardResponse: web.StandardResponse{},
		Devices:          deviceList,
	}
}

// DeviceGet gets the device for an organization
func (srv *Management) DeviceGet(orgID, username string, role int, deviceID string) web.DeviceResponse {
	newOrgID, err := getUserOrgIDIfOrgName(srv, username, orgID)
	if err != nil {
		return web.DeviceResponse{
			StandardResponse: web.StandardResponse{
				Code:    "Error",
				Message: err.Error(),
			},
		}
	}

	orgID = newOrgID

	hasAccess := srv.DS.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.DeviceResponse{
			StandardResponse: web.StandardResponse{
				Code:    "DeviceAuth",
				Message: "the user does not have permissions for the organization",
			},
		}
	}

	device, err := srv.DeviceTwinController.DeviceGet(orgID, deviceID)
	if err != nil {
		return web.DeviceResponse{
			StandardResponse: web.StandardResponse{
				Code:    "Error",
				Message: err.Error(),
			},
		}
	}

	return web.DeviceResponse{
		StandardResponse: web.StandardResponse{},
		Device:           device,
	}
}

// DeviceDelete deletes the device from an organization
func (srv *Management) DeviceDelete(orgID, username string, role int, deviceID string) web.StandardResponse {
	newOrgID, err := getUserOrgIDIfOrgName(srv, username, orgID)
	if err != nil {
		return web.StandardResponse{
			Code:    "Error",
			Message: err.Error(),
		}
	}

	orgID = newOrgID

	hasAccess := srv.DS.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "DeviceAuth",
			Message: "the user does not have permissions for the organization",
		}
	}

	twinAPIMessage := "device deleted"
	identityAPIMessage := "device deleted"

	err = srv.DeviceTwinController.DeviceDelete(deviceID)
	if err != nil {
		twinAPIMessage = err.Error()
	}
	_, err = srv.Identity.DeleteDevice(deviceID)
	if err != nil {
		identityAPIMessage = err.Error()
	}

	message := fmt.Sprintf("twinapi: %s, identity: %s", twinAPIMessage, identityAPIMessage)
	return web.StandardResponse{Message: message}
}

// DeviceLogs requests from the DeviceTwin API that logs for a device be sent
func (srv *Management) DeviceLogs(orgID, username string, role int, deviceID string, logs *messages.DeviceLogs) web.StandardResponse {
	newOrgID, err := getUserOrgIDIfOrgName(srv, username, orgID)
	if err != nil {
		return web.StandardResponse{
			Code:    "Error",
			Message: err.Error(),
		}
	}

	orgID = newOrgID

	hasAccess := srv.DS.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "DeviceAuth",
			Message: "the user does not have permissions for the organization",
		}
	}

	err = srv.DeviceTwinController.DeviceLogs(orgID, deviceID, logs)
	if err != nil {
		return web.StandardResponse{
			Code:    "Error",
			Message: err.Error(),
		}
	}

	return web.StandardResponse{}
}

// DeviceUsersAction requests from the DeviceTwin API that a user action be performed on the device
func (srv *Management) DeviceUsersAction(orgID, username string, role int, deviceID string, deviceUser messages.DeviceUser) web.StandardResponse {
	newOrgID, err := getUserOrgIDIfOrgName(srv, username, orgID)
	if err != nil {
		return web.StandardResponse{
			Code:    "Error",
			Message: err.Error(),
		}
	}

	orgID = newOrgID

	hasAccess := srv.DS.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "DeviceAuth",
			Message: "the user does not have permissions for the organization",
		}
	}

	err = srv.DeviceTwinController.User(orgID, deviceID, deviceUser)
	if err != nil {
		return web.StandardResponse{
			Code:    "Error",
			Message: err.Error(),
		}
	}

	return web.StandardResponse{}
}
