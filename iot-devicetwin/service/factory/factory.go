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

// Package factory creates the Datastore as configured
package factory

import (
	"fmt"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"
	"github.com/everactive/dmscore/iot-devicetwin/datastore/memory"
	"github.com/everactive/dmscore/iot-devicetwin/datastore/postgres"
)

// CreateDataStore is the factory method to create a data store
var CreateDataStore = createDataStore

func createDataStore(databaseDriver string, dataStoreSource string) (datastore.DataStore, error) {
	var db datastore.DataStore
	switch databaseDriver {
	case "memory":
		db = memory.NewStore()
	case "postgres":
		db = postgres.OpenDataStore(databaseDriver, dataStoreSource)
	default:
		return nil, fmt.Errorf("unknown data store driver: %v", databaseDriver)
	}

	return db, nil
}
