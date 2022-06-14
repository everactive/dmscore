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

package controller

import (
	"encoding/json"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/actions"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
)

// User sends a user action to the device, which will either add or remove a user.
func (srv *Service) User(orgID, clientID string, user messages.DeviceUser) error {
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	act := messages.SubscribeAction{
		Action: actions.User,
		Data:   string(jsonBytes),
	}
	return srv.deviceSnapAction(orgID, clientID, act)
}
