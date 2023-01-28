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

package web

import (
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"

	"github.com/everactive/dmscore/iot-devicetwin/domain"
)

// StandardResponse is the JSON response from an API method, indicating success or failure.
type StandardResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SnapsResponse is the JSON response to list snaps
type SnapsResponse struct {
	StandardResponse
	Snaps []messages.DeviceSnap `json:"snaps"`
}

// DeviceResponse is the JSON response to get a device
type DeviceResponse struct {
	StandardResponse
	Device messages.Device `json:"device"`
}

// DevicesResponse is the JSON response to list devices
type DevicesResponse struct {
	StandardResponse
	Devices []messages.Device `json:"devices"`
}

// ActionsResponse is the JSON response to list actions for a device
type ActionsResponse struct {
	StandardResponse
	Actions []domain.Action `json:"actions"`
}

// GroupsResponse is the JSON response to list groups
type GroupsResponse struct {
	StandardResponse
	Groups []domain.Group `json:"groups"`
}

// GroupResponse is the JSON response to list groups
type GroupResponse struct {
	StandardResponse
	Group domain.Group `json:"group"`
}
