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
	"encoding/json"
	"path"

	"github.com/everactive/dmscore/iot-devicetwin/web"
)

const (
	snapUpdate        = "SnapUpdate"
	deviceUserLiteral = "DeviceUser"
)

// SnapList lists the snaps for a device
func (a *ClientAdapter) SnapList(orgID, deviceID string) web.SnapsResponse {
	r := web.SnapsResponse{}
	p := path.Join("device", orgID, deviceID, "snaps")

	resp, err := a.get(a.urlPath(p))
	if err != nil {
		r.StandardResponse.Message = err.Error()
		return r
	}

	// Parse the response
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		r.StandardResponse.Message = err.Error()
	}
	return r
}

// SnapListOnDevice triggers snap list on a device
func (a *ClientAdapter) SnapListOnDevice(orgID, deviceID string) web.StandardResponse {
	r := web.StandardResponse{}
	p := path.Join("device", orgID, deviceID, "snaps", "list")

	resp, err := a.post(a.urlPath(p), nil)
	if err != nil {
		r.Code = "SnapList"
		r.Message = err.Error()
		return r
	}

	// Parse the response
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		r.Code = "SnapList"
		r.Message = err.Error()
	}
	return r
}

// SnapInstall installs a snap on a device
func (a *ClientAdapter) SnapInstall(orgID, deviceID, snap string) web.StandardResponse {
	r := web.StandardResponse{}
	p := path.Join("device", orgID, deviceID, "snaps", snap)

	resp, err := a.post(a.urlPath(p), nil)
	if err != nil {
		r.Code = "SnapInstall"
		r.Message = err.Error()
		return r
	}

	// Parse the response
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		r.Code = "SnapInstall"
		r.Message = err.Error()
	}
	return r
}

// SnapRemove uninstalls a snap on a device
func (a *ClientAdapter) SnapRemove(orgID, deviceID, snap string) web.StandardResponse {
	r := web.StandardResponse{}
	p := path.Join("device", orgID, deviceID, "snaps", snap)

	resp, err := a.delete(a.urlPath(p))
	if err != nil {
		r.Code = "SnapRemove"
		r.Message = err.Error()
		return r
	}

	// Parse the response
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		r.Code = "SnapRemove"
		r.Message = err.Error()
	}
	return r
}

// SnapUpdate enables/disables/refreshes/switch a snap on a device
func (a *ClientAdapter) SnapUpdate(orgID, deviceID, snap, action string, body []byte) web.StandardResponse { //nolint
	r := web.StandardResponse{}
	p := path.Join("device", orgID, deviceID, "snaps", snap, action)

	resp, err := a.put(a.urlPath(p), body)
	if err != nil {
		r.Code = snapUpdate
		r.Message = err.Error()
		return r
	}

	// Parse the response
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		r.Code = snapUpdate
		r.Message = err.Error()
	}
	return r
}

// SnapConfigSet sets a snap config on a device
func (a *ClientAdapter) SnapConfigSet(orgID, deviceID, snap string, config []byte) web.StandardResponse {
	r := web.StandardResponse{}
	p := path.Join("device", orgID, deviceID, "snaps", snap, "settings")

	resp, err := a.put(a.urlPath(p), config)
	if err != nil {
		r.Code = snapUpdate
		r.Message = err.Error()
		return r
	}

	// Parse the response
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		r.Code = snapUpdate
		r.Message = err.Error()
	}
	return r
}

func (a *ClientAdapter) SnapServiceAction(orgID, deviceID, snap, action string, body []byte) web.StandardResponse { //nolint
	r := web.StandardResponse{}
	p := path.Join("device", orgID, deviceID, "services", snap, action)

	resp, err := a.post(a.urlPath(p), body)
	if err != nil {
		r.Code = "SnapServiceAction"
		r.Message = err.Error()
		return r
	}

	// Parse the response
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		r.Code = "SnapServiceAction"
		r.Message = err.Error()
	}
	return r
}

// SnapSnapshot sends a snap snapshot request to a device and decodes it into a StandardResponse
func (a *ClientAdapter) SnapSnapshot(orgID, deviceID, snap string, body []byte) web.StandardResponse {
	r := web.StandardResponse{}
	p := path.Join("device", orgID, deviceID, "snaps", snap, "snapshot")

	resp, err := a.post(a.urlPath(p), body)
	if err != nil {
		r.Code = "SnapPost"
		r.Message = err.Error()
		return r
	}

	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		r.Code = "SnapPost"
		r.Message = err.Error()
		return r
	}

	return r

}
