// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Identity Service
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

package postgres

import (
	"fmt"
	"github.com/everactive/dmscore/iot-identity/datastore"
	"github.com/everactive/dmscore/iot-identity/domain"
	"github.com/everactive/dmscore/iot-identity/models"
)

// DeviceNew creates a new device registration
func (s *Store) DeviceNew(d datastore.DeviceNewRequest) (string, error) {
	var id int64
	var deviceID = datastore.GenerateID()

	err := s.QueryRow(createDeviceSQL, deviceID, d.OrganizationID, d.Brand, d.Model, d.SerialNumber, d.Credentials.PrivateKey, d.Credentials.Certificate, d.Credentials.MQTTURL, d.Credentials.MQTTPort, d.DeviceData).Scan(&id)
	if err != nil {
		datastore.Logger.Error("Error creating device: ", err)
	}

	return deviceID, err
}

// DeviceDelete deletes a device from the database
func (s *Store) DeviceDelete(deviceID string) (string, error) {
	datastore.Logger.Tracef("Deleting device: %s", deviceID)
	_, err := s.Exec(deleteDeviceByID, deviceID)
	if err != nil {
		datastore.Logger.Error("Error deleting device: ", err)
		return deviceID, fmt.Errorf("error deleting device: %w", err)
	}

	return deviceID, nil
}

// DeviceGet fetches a device registration
func (s *Store) DeviceGet(brand, model, serial string) (*domain.Enrollment, error) {
	d := domain.Enrollment{
		Device:       domain.Device{},
		Organization: domain.Organization{},
		Credentials:  domain.Credentials{},
	}

	err := s.QueryRow(getDeviceSQL, brand, model, serial).Scan(
		&d.ID, &d.Organization.ID, &d.Device.Brand, &d.Device.Model, &d.Device.SerialNumber,
		&d.Credentials.PrivateKey, &d.Credentials.Certificate, &d.Credentials.MQTTURL, &d.Credentials.MQTTPort,
		&d.Device.StoreID, &d.Device.DeviceKey, &d.Status, &d.DeviceData)
	if err != nil {
		datastore.Logger.Error("Error retrieving device:", err)
		return &d, fmt.Errorf("error retrieving device: %w", err)
	}

	// Get the organization details for the device
	org, err := s.OrganizationGet(d.Organization.ID)
	if err != nil {
		datastore.Logger.Error("Error retrieving device organization:", err)
		return &d, fmt.Errorf("error retrieving device organization: %w", err)
	}
	d.Organization = *org

	return &d, err
}

// DeviceGetEnrollmentByID fetches a device registration
func (s *Store) DeviceGetEnrollmentByID(deviceID string) (*domain.Enrollment, error) {
	d := domain.Enrollment{
		Device:       domain.Device{},
		Organization: domain.Organization{},
		Credentials:  domain.Credentials{},
	}

	registeredDevice := models.RegisteredDevice{}
	res := s.gormDB.Where("device_id = ?", deviceID).Or("serial_number = ?", deviceID).Find(&registeredDevice)
	if res.Error != nil {
		panic(res.Error)
	}

	domain.Enrollment{}.FromRegisteredDeviceModel(&registeredDevice, &d)

	// Get the organization details for the device
	org, err := s.OrganizationGet(d.Organization.ID)
	if err != nil {
		datastore.Logger.Error("Error retrieving device organization:", err)
		return &d, fmt.Errorf("error retrieving device organization: %w", err)
	}
	d.Organization = *org

	return &d, err
}

// DeviceEnroll enrolls a device with the IoT service
func (s *Store) DeviceEnroll(d datastore.DeviceEnrollRequest) (*domain.Enrollment, error) {
	_, err := s.Exec(enrollDeviceSQL, d.Brand, d.Model, d.SerialNumber, d.StoreID, d.DeviceKey, models.StatusEnrolled)
	if err != nil {
		datastore.Logger.Error("Error updating the device:", err)
	}

	return s.DeviceGet(d.Brand, d.Model, d.SerialNumber)
}

// DeviceList fetches the device registrations for an organization
func (s *Store) DeviceList(orgID string) ([]domain.Enrollment, error) {
	devices := []domain.Enrollment{}

	dbDevices := []models.RegisteredDevice{}
	res := s.gormDB.Model(&models.RegisteredDevice{}).Where("org_id = ?", orgID).Find(&dbDevices)
	if res.Error != nil {
		return devices, res.Error
	}

	if res.RowsAffected == 0 {
		// May be using org name instead of id
		org := models.Organization{}
		res = s.gormDB.Model(&models.Organization{}).Where("name = ?", orgID).Find(&org)
		if res.RowsAffected == 0 {
			// If we didn't find the devices by org id and the org doesn't exist by name, then it's an error
			return devices, res.Error
		}

		res := s.gormDB.Model(&models.RegisteredDevice{}).Where("org_id = ?", org.OrgId).Find(&dbDevices)
		if res.Error != nil {
			return devices, res.Error
		}
	}

	for _, d := range dbDevices {
		device := domain.Enrollment{}
		domain.Enrollment{}.FromRegisteredDeviceModel(&d, &device)
		devices = append(devices, device)
	}

	return devices, nil
}

// DeviceUpdate updates a device registration
func (s *Store) DeviceUpdate(deviceID string, status models.Status, deviceData string) error {
	res := s.gormDB.Model(&models.RegisteredDevice{}).
		Where("device_id = ?", deviceID).
		Or("serial_number = ?", deviceID).
		Updates(&models.RegisteredDevice{
			Status:     status,
			DeviceData: deviceData,
		})

	return res.Error
}
