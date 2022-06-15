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
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
)

// DeviceSnaps fetches the snaps for a device
func (srv *Service) DeviceSnaps(orgID, clientID string) ([]messages.DeviceSnap, error) {
	device, err := srv.DB.DeviceGet(clientID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	snaps, err := srv.DB.DeviceSnapList(int64(device.ID))
	if err != nil {
		return nil, err
	}

	installed := []messages.DeviceSnap{}
	for _, s := range snaps {
		snap := messages.DeviceSnap{
			DeviceId:      device.DeviceID,
			Name:          s.Name,
			InstalledSize: s.InstalledSize,
			InstalledDate: s.InstalledDate,
			Status:        s.Status,
			Channel:       s.Channel,
			Confinement:   s.Confinement,
			Version:       s.Version,
			Revision:      s.Revision,
			Devmode:       s.Devmode,
			Config:        s.Config,
		}

		services := []*messages.ServiceStatus{}
		bytes, _ := json.Marshal(&s.ServiceStatuses)
		fmt.Printf("marshaled service statuses: %s", string(bytes))
		for _, service := range s.ServiceStatuses {
			logrus.Tracef("Adding service %s to snap %s. [%+v]", service.Name, snap.Name, service)
			services = append(services, &messages.ServiceStatus{
				Name:    service.Name,
				Active:  service.Active,
				Enabled: service.Enabled,
				Daemon:  service.Daemon,
			})
		}

		snap.Services = services

		installed = append(installed, snap)
	}
	return installed, nil
}
