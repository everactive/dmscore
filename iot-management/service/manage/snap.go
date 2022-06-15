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
	"fmt"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	"github.com/everactive/dmscore/iot-devicetwin/web"
)

// SnapList lists the snaps for a device
func (srv *Management) SnapList(orgID, username string, role int, deviceID string) web.SnapsResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.SnapsResponse{
			StandardResponse: web.StandardResponse{
				Code:    "SnapsAuth",
				Message: "the user does not have permissions for the organization",
			},
		}
	}

	err, response, match := srv.verifyOrgMatches(orgID, deviceID)
	if !match {
		return web.SnapsResponse{StandardResponse: response}
	}

	snaps, err := srv.DeviceTwinService.Controller.DeviceSnaps(orgID, deviceID)
	if err != nil {
		return web.SnapsResponse{StandardResponse: web.StandardResponse{
			Code:    "SnapList",
			Message: err.Error(),
		}}
	}

	return web.SnapsResponse{Snaps: snaps}
}

func (srv *Management) verifyOrgMatches(orgID string, deviceID string) (error, web.StandardResponse, bool) {
	// Sanity check, verify orgID matches device
	enrollment, err := srv.IdentityService.Identity.DeviceGet(orgID, deviceID)
	if err != nil {
		return nil, web.StandardResponse{
			Code:    "OrgIdOrName",
			Message: err.Error(),
		}, false
	}

	if enrollment.Organization.ID != orgID && enrollment.Organization.Name != orgID {
		return nil, web.StandardResponse{
			Code:    "OrgIdOrName",
			Message: fmt.Sprintf("orgID=%s did not match expect device organization id or name", orgID),
		}, false
	}

	return err, web.StandardResponse{}, true
}

// SnapListOnDevice lists snaps on a device
func (srv *Management) SnapListOnDevice(orgID, username string, role int, deviceID string) web.StandardResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "SnapAuth",
			Message: "the user does not have permissions for the organization",
		}
	}

	err, response, match := srv.verifyOrgMatches(orgID, deviceID)
	if !match {
		return response
	}

	err = srv.DeviceTwinService.Controller.DeviceSnapList(orgID, deviceID)
	if err != nil {
		return web.StandardResponse{
			Code:    "SnapListOnDevice",
			Message: err.Error(),
		}
	}

	return web.StandardResponse{}
}

// SnapInstall installs a snap on a device
func (srv *Management) SnapInstall(orgID, username string, role int, deviceID, snap string) web.StandardResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "SnapAuth",
			Message: "the user does not have permissions for the organization",
		}
	}

	err := srv.DeviceTwinService.Controller.DeviceSnapInstall(orgID, deviceID, snap)
	if err != nil {
		return web.StandardResponse{
			Code:    "SnapInstall",
			Message: err.Error(),
		}
	}

	return web.StandardResponse{}
}

// SnapRemove uninstalls a snap on a device
func (srv *Management) SnapRemove(orgID, username string, role int, deviceID, snap string) web.StandardResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "SnapAuth",
			Message: "the user does not have permissions for the organization",
		}
	}

	err := srv.DeviceTwinService.Controller.DeviceSnapRemove(orgID, deviceID, snap)
	if err != nil {
		return web.StandardResponse{
			Code:    "SnapRemove",
			Message: err.Error(),
		}
	}

	return web.StandardResponse{}
}

// SnapUpdate enables/disables/refreshes/swtich a snap on a device
func (srv *Management) SnapUpdate(orgID, username string, role int, deviceID, snap, action string, body []byte) web.StandardResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "SnapAuth",
			Message: "the user does not have permissions for the organization",
		}
	}

	snapUpdate := messages.SnapUpdate{}
	err := json.Unmarshal(body, &snapUpdate)

	err = srv.DeviceTwinService.Controller.DeviceSnapUpdate(orgID, deviceID, snap, action, &snapUpdate)
	if err != nil {
		return web.StandardResponse{
			Code:    "SnapUpdate",
			Message: err.Error(),
		}
	}

	return web.StandardResponse{}
}

// SnapConfigSet updates a snap config on a device
func (srv *Management) SnapConfigSet(orgID, username string, role int, deviceID, snap string, config []byte) web.StandardResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "SnapAuth",
			Message: "the user does not have permissions for the organization",
		}
	}

	err := srv.DeviceTwinService.Controller.DeviceSnapConf(orgID, deviceID, snap, string(config))
	if err != nil {
		return web.StandardResponse{
			Code:    "SnapConfigSet",
			Message: err.Error(),
		}
	}

	return web.StandardResponse{}
}

// SnapServiceAction requests from the DeviceTwin API that an action be performed on a snap service
func (srv *Management) SnapServiceAction(orgID, username string, role int, deviceID, snap, action string, body []byte) web.StandardResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "SnapAuth",
			Message: "the user does not have permissions for the organization",
		}
	}

	services := messages.SnapService{}
	err := json.Unmarshal(body, &services)
	if err != nil {
		return web.StandardResponse{
			Code:    "SnapServiceAction",
			Message: err.Error(),
		}
	}

	err = srv.DeviceTwinService.Controller.DeviceSnapServiceAction(orgID, deviceID, snap, action, &services)
	if err != nil {
		return web.StandardResponse{
			Code:    "SnapServiceAction",
			Message: err.Error(),
		}
	}

	return web.StandardResponse{}
}

// SnapSnapshot requests from the DeviceTwin API that a snapshot be made of a given snap
func (srv *Management) SnapSnapshot(orgID, username string, role int, deviceID, snap string, body []byte) web.StandardResponse {
	hasAccess := srv.DB.OrgUserAccess(orgID, username, role)
	if !hasAccess {
		return web.StandardResponse{
			Code:    "SnapAuth",
			Message: "the user does not have permissions for the organization",
		}
	}

	snapshot := messages.SnapSnapshot{}
	err := json.Unmarshal(body, &snapshot)
	if err != nil {
		return web.StandardResponse{
			Code:    "SnapSnapshot",
			Message: err.Error(),
		}
	}

	err = srv.DeviceTwinService.Controller.DeviceSnapSnapshot(orgID, deviceID, snap, &snapshot)
	if err != nil {
		return web.StandardResponse{
			Code:    "SnapSnapshot",
			Message: err.Error(),
		}
	}

	return web.StandardResponse{}
}
