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

// Package datastore provides the DataStore interface and entity types
package datastore

import (
	"github.com/everactive/dmscore/iot-management/datastore/models"
	"github.com/juju/usso/openid"
)

// DataStore is the interfaces for the data repository
type DataStore interface {
	OpenIDNonceStore() openid.NonceStore
	CreateUser(user User) (int64, error)
	GetUser(username string) (User, error)
	UserList() ([]User, error)
	UserUpdate(user User) error
	UserDelete(username string) error

	OrgUserAccess(orgID, username string, role int) bool
	OrganizationsForUser(username string) ([]Organization, error)
	OrganizationForUserToggle(orgID, username string) error
	OrganizationGet(orgIDOrName string) (Organization, error)
	OrganizationCreate(org Organization) error
	OrganizationUpdate(org Organization) error

	GetSettings() ([]models.Setting, error)
	Set(key string, value string) error
}
