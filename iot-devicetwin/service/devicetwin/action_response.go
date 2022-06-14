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
	"log"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
)

// actionDevice process the device info received from a device
func (srv *Service) actionDevice(payload []byte) error {
	// Parse the payload
	d := messages.PublishDevice{}
	if err := json.Unmarshal(payload, &d); err != nil {
		log.Printf("Error in device action message: %v", err)
		return fmt.Errorf("error in device action message: %v", err)
	}

	// Get the device details and create/update the device
	_, err := srv.DB.DeviceGet(d.Result.DeviceId)
	if err == nil {
		return fmt.Errorf("error in device action: device already exists")
	}

	// Device does not exit, so create
	device := datastore.Device{
		OrganisationID: d.Result.OrgId,
		DeviceID:       d.Result.DeviceId,
		Brand:          d.Result.Brand,
		DeviceModel:    d.Result.Model,
		SerialNumber:   d.Result.Serial,
		DeviceKey:      d.Result.DeviceKey,
		StoreID:        d.Result.Store,
	}
	deviceID, err := srv.DB.DeviceCreate(device)
	if err != nil {
		return fmt.Errorf("error in device action: %v", err)
	}

	if d.Result.Version == nil || d.Result.Version.DeviceId == "" {
		// No device version information
		return nil
	}

	// Create or update the device version record
	version := datastore.DeviceVersion{
		DeviceID:      deviceID,
		Version:       d.Result.Version.Version,
		Series:        d.Result.Version.Series,
		OSID:          d.Result.Version.OsId,
		OSVersionID:   d.Result.Version.OsVersionId,
		OnClassic:     d.Result.Version.OnClassic,
		KernelVersion: d.Result.Version.KernelVersion,
	}
	err = srv.DB.DeviceVersionUpsert(version)
	return err
}

// actionList process the list of snaps received from a device
func (srv *Service) actionList(clientID string, payload []byte) error {
	// Parse the payload
	p := messages.PublishSnaps{}
	if err := json.Unmarshal(payload, &p); err != nil {
		log.Printf("Error in list action message: %v", err)
		return fmt.Errorf("error in list action message: %v", err)
	}

	// Get the device details
	device, err := srv.DB.DeviceGet(clientID)
	if err != nil {
		return fmt.Errorf("cannot find device with ID `%s`", clientID)
	}

	// Add the installed snaps
	for _, s := range p.Result {
		snap := datastore.DeviceSnap{
			DeviceID:      int64(device.ID),
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

		if len(s.Services) > 0 {
			serviceStatuses := []*datastore.ServiceStatus{}
			for _, service := range s.Services {
				serviceStatuses = append(serviceStatuses, &datastore.ServiceStatus{
					Name:    service.Name,
					Daemon:  service.Daemon,
					Enabled: service.Enabled,
					Active:  service.Active,
				})
			}

			snap.ServiceStatuses = serviceStatuses
		}

		if err := srv.DB.DeviceSnapUpsert(snap); err != nil {
			return err
		}

	}

	return nil
}

// actionForSnap process the snap response from an action (install, remove, refresh...)
func (srv *Service) actionForSnap(_, action string, payload []byte) (string, error) {
	// Parse the payload
	p := messages.PublishSnapTask{}
	if err := json.Unmarshal(payload, &p); err != nil {
		log.Printf("Error in %s action message: %v", action, err)
		return "", fmt.Errorf("error in %s action message: %v", action, err)
	}

	// The payload is the task ID of the action, log it
	return p.Result, nil
}

// actionConf process the snap response from a conf action
func (srv *Service) actionConf(clientID string, payload []byte) error {
	// Parse the payload
	p := messages.PublishSnap{}
	if err := json.Unmarshal(payload, &p); err != nil {
		log.Printf("Error in conf action message: %v", err)
		return fmt.Errorf("error in conf action message: %v", err)
	}

	// Get the device details
	device, err := srv.DB.DeviceGet(clientID)
	if err != nil {
		return fmt.Errorf("cannot find device with ID `%s`", clientID)
	}

	// Create/update the installed snap details with the current config
	snap := datastore.DeviceSnap{
		// Created       time.Time
		// Modified      time.Time
		DeviceID:      int64(device.ID),
		Name:          p.Result.Name,
		InstalledSize: p.Result.InstalledSize,
		InstalledDate: p.Result.InstalledDate,
		Status:        p.Result.Status,
		Channel:       p.Result.Channel,
		Confinement:   p.Result.Confinement,
		Version:       p.Result.Version,
		Revision:      p.Result.Revision,
		Devmode:       p.Result.Devmode,
		Config:        p.Result.Config,
	}

	return srv.DB.DeviceSnapUpsert(snap)
}

// actionServer process the response from a server action
func (srv *Service) actionServer(clientID string, payload []byte) error {
	// Parse the payload
	// p := domain.PublishDeviceVersion{}
	p := messages.PublishDeviceVersion{}
	if err := json.Unmarshal(payload, &p); err != nil {
		log.Printf("Error in server action message: %v", err)
		return fmt.Errorf("error in server action message: %v", err)
	}

	// Get the device details
	device, err := srv.DB.DeviceGet(clientID)
	if err != nil {
		return fmt.Errorf("cannot find device with ID `%s`", clientID)
	}

	// Create/update the device OS details
	dv := datastore.DeviceVersion{
		DeviceID:      int64(device.ID),
		Version:       p.Result.Version,
		Series:        p.Result.Series,
		OSID:          p.Result.OsId,
		OSVersionID:   p.Result.OsVersionId,
		OnClassic:     p.Result.OnClassic,
		KernelVersion: p.Result.KernelVersion,
	}

	return srv.DB.DeviceVersionUpsert(dv)
}

// actionUnregister process the response from an unregister action
func (srv *Service) actionUnregister(clientID string, payload []byte) error {
	// Parse the payload
	p := messages.PublishDevice{}
	if err := json.Unmarshal(payload, &p); err != nil {
		log.Printf("Error in unregister action message: %v", err)
		return fmt.Errorf("error in unregister action message: %v", err)
	}

	// Get the device details
	device, err := srv.DB.DeviceGet(clientID)
	if err != nil {
		return fmt.Errorf("cannot find device with ID `%s`", clientID)
	}

	// Delete the device twin from the database
	_, err = srv.DeviceDelete(device.DeviceID)
	if err != nil {
		return fmt.Errorf("error in unregister action when deleting device from database: %v", err)
	}

	return nil
}
