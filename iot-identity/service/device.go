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

package service

import (
	"fmt"
	"github.com/everactive/dmscore/iot-identity/models"

	"github.com/everactive/dmscore/iot-identity/config/configkey"
	"github.com/everactive/dmscore/iot-identity/datastore"
	"github.com/everactive/dmscore/iot-identity/domain"
	"github.com/everactive/dmscore/iot-identity/service/cert"
	"github.com/spf13/viper"
)

// DeviceList fetches the registered devices
func (id IdentityService) DeviceList(orgID string) ([]domain.Enrollment, error) {
	return id.DB.DeviceList(orgID)
}

// DeviceGet fetches a device registration
func (id IdentityService) DeviceGet(orgID, deviceID string) (*domain.Enrollment, error) {
	return id.DB.DeviceGetEnrollmentByID(deviceID)
}

// DeleteDevice deletes a device already in the service
func (id IdentityService) DeleteDevice(deviceID string) (string, error) {
	return id.DB.DeviceDelete(deviceID)
}

// RegisterDevice registers a new device with the service
func (id IdentityService) RegisterDevice(req *RegisterDeviceRequest) (string, error) {
	// Validate fields
	for k, v := range map[string]string{
		"organization ID": req.OrganizationID,
		"brand":           req.Brand,
		"model name":      req.Model,
		"serial number":   req.SerialNumber,
	} {
		if err := validateNotEmpty(k, v); err != nil {
			return "", err
		}
	}

	// Check that the organization exists
	org, err := id.DB.OrganizationGet(req.OrganizationID)
	if err != nil {
		return "", err
	}

	// Check that the device has not been registered
	if _, err = id.DB.DeviceGet(req.Brand, req.Model, req.SerialNumber); err == nil {
		return "", fmt.Errorf("the device `%s/%s/%s` is already registered", req.Brand, req.Model, req.SerialNumber)
	}

	rootCertsDir := viper.GetString(configkey.MQTTCertificatePath)

	// Create a signed certificate
	deviceID := datastore.GenerateID()
	keyPEM, certPEM, err := cert.CreateClientCert(org, rootCertsDir, deviceID)
	if err != nil {
		return "", err
	}

	MQTTURL := viper.GetString(configkey.MQTTHostAddress)
	MQTTPort := viper.GetString(configkey.MQTTHostPort)

	// Create registration
	d := datastore.DeviceNewRequest{
		ID:             deviceID,
		OrganizationID: req.OrganizationID,
		Brand:          req.Brand,
		Model:          req.Model,
		SerialNumber:   req.SerialNumber,
		Credentials: domain.Credentials{
			PrivateKey:  keyPEM,
			Certificate: certPEM,
			MQTTURL:     MQTTURL, // Using a default URL for all devices
			MQTTPort:    MQTTPort,
		},
		DeviceData: req.DeviceData,
	}
	return id.DB.DeviceNew(d)
}

// DeviceUpdate updates an existing device with the service
// Status changes are limited, depending on whether the device has enrolled with the service. If it has, then it
// already has credentials.
// If a device has not enrolled:
// - Waiting => Disabled
// - Disabled => Waiting
// If a device has enrolled:
// - Enrolled => Disabled (TODO: needs to trigger the removal of credentials from MQTT broker or device or both)
// - Enrolled => Waiting
func (id IdentityService) DeviceUpdate(orgID, deviceID string, req *DeviceUpdateRequest) error {
	// Get the device and check the current status
	device, err := id.DB.DeviceGetEnrollmentByID(deviceID)
	if err != nil {
		return err
	}

	// Update the device data, if it has changed
	//if device.DeviceData != req.DeviceData {
	//	if err := id.DB.DeviceUpdate(deviceID, device.Status, req.DeviceData); err != nil {
	//		return err
	//	}
	//}

	if req.Status == int(models.StatusEnrolled) {
		return fmt.Errorf("cannot change a device status to enrolled. The device itself needs to connect for this")
	}

	switch device.Status {
	case models.StatusWaiting:
		if req.Status == int(models.StatusWaiting) {
			// No change required
			return nil
		}
		device.Status = models.StatusDisabled
	case models.StatusDisabled:
		if req.Status == int(models.StatusDisabled) {
			// No change required
			return nil
		}
		device.Status = models.StatusWaiting
	case models.StatusEnrolled:
		if req.Status == int(models.StatusDisabled) {
			//nolint:godox
			// TODO: trigger the removal of credentials from MQTT broker or device or both
			// -- from original code before linting
			device.Status = models.StatusDisabled
		} else {
			device.Status = models.StatusWaiting
		}
	}

	return id.DB.DeviceUpdate(deviceID, device.Status, req.DeviceData)
}
