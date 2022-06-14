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

// Package datastore is the types and interface for the device datastore
package datastore

import (
	"time"

	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// Action is the log of an action request
type Action struct {
	ID             int64
	Created        time.Time
	Modified       time.Time
	OrganizationID string
	DeviceID       string
	ActionID       string
	Action         string
	Status         string
	Message        string
}

// Device the repository definition of a device
type Device struct {
	gorm.Model
	LastRefresh    time.Time     `gorm:"column:lastrefresh"`
	OrganisationID string        `gorm:"column:org_id"`
	DeviceID       string        `gorm:"column:device_id"`
	Brand          string        `gorm:"column:brand"`
	DeviceModel    string        `gorm:"column:model"`
	SerialNumber   string        `gorm:"column:serial"`
	DeviceKey      string        `gorm:"column:device_key"`
	StoreID        string        `gorm:"column:store_id"`
	Active         bool          `gorm:"column:active"`
	DeviceVersion  DeviceVersion `gorm:"constraint:OnDelete:CASCADE"`
	DeviceSnaps    []*DeviceSnap `gorm:"constraint:OnDelete:CASCADE"`
}

// TableName is the Postgres table name to use
func (Device) TableName() string {
	return "device"
}

// AfterDelete is used to take additional database actions after a device is deleted
func (d *Device) AfterDelete(tx *gorm.DB) (err error) {
	logrus.Tracef("Trying to delete device_version where device is = %+v\n", d)
	deviceVersion := DeviceVersion{}

	tx.Where(&DeviceVersion{DeviceID: int64(d.ID)}).Find(&deviceVersion)
	tx.Delete(&DeviceVersion{}, deviceVersion.ID)

	deviceSnaps := []DeviceSnap{}
	tx.Where(&DeviceSnap{DeviceID: int64(d.ID)}).Find(&deviceSnaps)
	for _, ds := range deviceSnaps {
		logrus.Tracef("Deleting snap %+v for device", ds)
		tx.Delete(&DeviceSnap{}, ds.ID)
	}

	return
}

// DeviceSnap holds the details of snap on a device
type DeviceSnap struct {
	gorm.Model
	DeviceID        int64            `gorm:"column:device_id"`
	Name            string           `gorm:"column:name"`
	InstalledSize   int64            `gorm:"column:installed_size"`
	InstalledDate   time.Time        `gorm:"column:installed_date"`
	Status          string           `gorm:"column:status"`
	Channel         string           `gorm:"column:channel"`
	Confinement     string           `gorm:"column:confinement"`
	Version         string           `gorm:"column:version"`
	Revision        int              `gorm:"column:revision"`
	Devmode         bool             `gorm:"column:devmode"`
	Config          string           `gorm:"column:config"`
	ServiceStatuses []*ServiceStatus `gorm:"constraint:OnDelete:CASCADE"`
}

// AfterDelete is used to take additional database actions after a device snap is deleted
func (d *DeviceSnap) AfterDelete(tx *gorm.DB) (err error) {
	serviceStatuses := []ServiceStatus{}
	tx.Where(&ServiceStatus{DeviceSnapID: int64(d.ID)}).Find(&serviceStatuses)
	for _, ss := range serviceStatuses {
		logrus.Tracef("Deleting service status %+v for device", ss)
		tx.Delete(&ServiceStatus{}, ss.ID)
	}

	return
}

// TableName specifies a custom table name for legacy tables to use with gorm
func (DeviceSnap) TableName() string {
	return "device_snap"
}

// DeviceVersion holds the details of the OS details on the device
type DeviceVersion struct {
	gorm.Model
	DeviceID      int64
	Version       string
	Series        string
	OSID          string
	OSVersionID   string
	OnClassic     bool
	KernelVersion string
}

// TableName specifies a custom table name for legacy tables to use with gorm
func (DeviceVersion) TableName() string {
	return "device_version"
}

// Group is the record for grouping devices
type Group struct {
	ID             int64
	Created        time.Time
	Modified       time.Time
	OrganisationID string
	Name           string
}

// GroupDeviceLink is the record for linking devices to groups
type GroupDeviceLink struct {
	ID             int64
	Created        time.Time
	OrganisationID string
	GroupID        int64
	DeviceID       int64
}

// ServiceStatus is the status of a service for a snap when a list of snaps or info for a single snap is retrieved
type ServiceStatus struct {
	gorm.Model
	DeviceSnapID int64
	Name         string
	Daemon       string
	Enabled      bool
	Active       bool
}
