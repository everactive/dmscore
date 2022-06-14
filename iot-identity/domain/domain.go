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

// Package domain provides the identity service specific data structures
package domain

import (
	"github.com/everactive/dmscore/iot-identity/models"
)

// Organization details for an account
type Organization struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RootCert []byte `json:"rootcert"`
	RootKey  []byte `json:"-,omitempty"`
}

// Device details
type Device struct {
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial"`
	StoreID      string `json:"store,omitempty"`
	DeviceKey    string `json:"deviceKey,omitempty"`
}

// Credentials for accessing the MQTT broker
type Credentials struct {
	PrivateKey  []byte `json:"privateKey,omitempty"`
	Certificate []byte `json:"certificate"`
	MQTTURL     string `json:"mqttUrl"`
	MQTTPort    string `json:"mqttPort"`
}

// Enrollment details for a device
type Enrollment struct {
	ID           string        `json:"id"`
	Device       Device        `json:"device"`
	Credentials  Credentials   `json:"credentials,omitempty"`
	Organization Organization  `json:"organization"`
	Status       models.Status `json:"status"`
	DeviceData   string        `json:"deviceData"`
}

func (Device) FromRegisteredDeviceModel(m *models.RegisteredDevice, d *Device) *Device {
	d.DeviceKey = m.DeviceKey
	d.Model = m.DeviceModel
	d.Brand = m.Brand
	d.StoreID = m.StoreID
	d.SerialNumber = m.SerialNumber
	return d
}

func (Credentials) FromRegisteredDeviceModel(m *models.RegisteredDevice, c *Credentials) *Credentials {
	c.MQTTURL = m.MQTTURL
	c.Certificate = m.Certificate
	c.PrivateKey = m.PrivateKey
	c.MQTTPort = m.MQTTPort
	return c
}

func (Enrollment) FromRegisteredDeviceModel(m *models.RegisteredDevice, e *Enrollment) *Enrollment {
	e.Status = m.Status
	e.DeviceData = m.DeviceData

	Device{}.FromRegisteredDeviceModel(m, &e.Device)
	Credentials{}.FromRegisteredDeviceModel(m, &e.Credentials)

	e.Organization.ID = m.OrgID

	return e
}
