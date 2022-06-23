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

package devicetwin

import (
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
	"github.com/everactive/dmscore/iot-devicetwin/domain"
)

// ActionCreate logs an action
func (srv *Service) ActionCreate(orgID, deviceID string, action messages.SubscribeAction) error {
	act := datastore.Action{
		OrganizationID: orgID,
		DeviceID:       deviceID,
		ActionID:       action.Id,
		Action:         action.Action,
		Status:         "requested",
	}
	_, err := srv.DB.ActionCreate(act)
	return err
}

// ActionUpdate updates action
func (srv *Service) ActionUpdate(actionID, status, message string) error {
	return srv.DB.ActionUpdate(actionID, status, message)
}

// ActionList lists actions for a device
func (srv *Service) ActionList(orgID, deviceID string) ([]domain.Action, error) {
	list := []domain.Action{}
	actions, err := srv.DB.ActionListForDevice(orgID, deviceID)
	if err != nil {
		return list, err
	}

	// Map the database item to the domain item
	for _, act := range actions {
		list = append(list, domain.Action{
			Created:        act.CreatedAt,
			Modified:       act.UpdatedAt,
			OrganizationID: act.OrganizationID,
			DeviceID:       act.DeviceID,
			ActionID:       act.ActionID,
			Action:         act.Action,
			Status:         act.Status,
			Message:        act.Message,
		})
	}

	return list, nil
}
