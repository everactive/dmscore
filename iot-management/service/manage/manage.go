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

// Package manage provides the interface and implementation of the Manage service
package manage

import (
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	"github.com/everactive/dmscore/iot-devicetwin/service/controller"
	"github.com/everactive/dmscore/iot-devicetwin/web"
	"github.com/everactive/dmscore/iot-identity/service"
	idweb "github.com/everactive/dmscore/iot-identity/web"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/domain"
	"github.com/everactive/dmscore/iot-management/identityapi"
	"github.com/everactive/dmscore/iot-management/twinapi"
	"github.com/everactive/dmscore/models"
	"github.com/juju/usso/openid"
	"gorm.io/gorm"
)

// Manage interface for the service
type Manage interface {
	OpenIDNonceStore() openid.NonceStore
	CreateUser(user domain.User) error
	GetUser(username string) (domain.User, error)
	UserList() ([]domain.User, error)
	UserUpdate(user domain.User) error
	UserDelete(username string) error

	RegDeviceList(orgID, username string, role int) idweb.DevicesResponse
	RegisterDevice(orgID, username string, role int, body []byte) idweb.RegisterResponse
	RegDeviceGet(orgID, username string, role int, deviceID string) idweb.EnrollResponse
	RegDeviceUpdate(orgID, username string, role int, deviceID string, body []byte) idweb.StandardResponse

	DeviceList(orgID, username string, role int) web.DevicesResponse
	DeviceGet(orgID, username string, role int, deviceID string) web.DeviceResponse
	DeviceDelete(orgID, username string, role int, deviceID string) web.StandardResponse
	DeviceLogs(orgID, username string, role int, deviceID string, logs *messages.DeviceLogs) web.StandardResponse
	DeviceUsersAction(orgID, username string, role int, deviceID string, deviceUser messages.DeviceUser) web.StandardResponse
	ActionList(orgID, username string, role int, deviceID string) web.ActionsResponse

	SnapSnapshot(orgID, username string, role int, deviceID, snap string, body []byte) web.StandardResponse
	SnapList(orgID, username string, role int, deviceID string) web.SnapsResponse
	SnapListOnDevice(orgID, username string, role int, deviceID string) web.StandardResponse
	SnapInstall(orgID, username string, role int, deviceID, snap string) web.StandardResponse
	SnapRemove(orgID, username string, role int, deviceID, snap string) web.StandardResponse
	SnapUpdate(orgID, username string, role int, deviceID, snap, action string, body []byte) web.StandardResponse
	SnapConfigSet(orgID, username string, role int, deviceID, snap string, config []byte) web.StandardResponse
	SnapServiceAction(orgID, username string, role int, deviceID, snap, action string, body []byte) web.StandardResponse

	GroupList(orgID, username string, role int) web.GroupsResponse
	GroupCreate(orgID, username string, role int, body []byte) web.StandardResponse
	GroupDevices(orgID, username string, role int, name string) web.DevicesResponse
	GroupExcludedDevices(orgID, username string, role int, name string) web.DevicesResponse
	GroupDeviceLink(orgID, username string, role int, name, deviceID string) web.StandardResponse
	GroupDeviceUnlink(orgID, username string, role int, name, deviceID string) web.StandardResponse

	OrganizationsForUser(username string) ([]domain.Organization, error)
	OrganizationForUserToggle(orgID, username string) error
	OrganizationGet(orgID string) (domain.Organization, error)
	OrganizationCreate(org domain.OrganizationCreate) error
	OrganizationUpdate(org domain.Organization) error

	AddModelRequiredSnap(orgID, username, modelName, snapName string, role int) (*models.DeviceModelRequiredSnap, error)
	GetModelRequiredSnaps(orgID, username, modelName string, role int) (*models.DeviceModel, error)
	DeleteModelRequiredSnap(orgID, username, modelName, snapName string, role int) error
}

// Management implementation of the management service use cases
type Management struct {
	DS                   datastore.DataStore
	DB                   *gorm.DB
	TwinAPI              twinapi.Client
	IdentityAPI          identityapi.Client
	DeviceTwinController controller.Controller
	Identity             service.Identity
}

// NewManagement creates an implementation of the management use cases
func NewManagement(db *gorm.DB, ds datastore.DataStore, api twinapi.Client, id identityapi.Client, dtc controller.Controller, ids service.Identity) *Management {
	return &Management{
		DS:                   ds,
		TwinAPI:              api,
		IdentityAPI:          id,
		DeviceTwinController: dtc,
		Identity:             ids,
		DB:                   db,
	}
}