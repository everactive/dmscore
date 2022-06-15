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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/everactive/dmscore/iot-management/config/configkey"
	"github.com/spf13/viper"

	"github.com/everactive/dmscore/iot-management/config"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/service/factory"
)

var username, name, email string

func main() {
	// Get the command line parameters
	parseFlags()

	var configFilePath string
	if filePath, ok := os.LookupEnv("CONFIG_FILE_PATH"); ok {
		configFilePath = filePath
	}

	config.LoadConfig(configFilePath)

	// Open the connection to the local database
	databaseDriver := viper.GetString(configkey.DatabaseDriver)
	dataSource := viper.GetString(configkey.DatabaseConnectionString)
	db, err := factory.CreateDataStore(databaseDriver, dataSource)
	if err != nil || db == nil {
		log.Fatalf("Error accessing data store: %v", databaseDriver)
		return
	}

	// Once we have a database connection, we need to load any settings that were persisted
	// in the database and make those available via Viper
	settings, err := db.GetSettings()
	if err != nil {
		log.Fatalf("Cannot continue without loading settings")
		return
	}

	for _, setting := range settings {
		viper.Set(setting.Key, setting.Value)
	}

	// Create the user
	err = run(db, username, name, email)
	if err != nil {
		fmt.Println("Error creating user:", err.Error())
		os.Exit(1)
	}
}

func run(db datastore.DataStore, username, name, email string) error {
	if len(username) == 0 {
		return fmt.Errorf("the username must be supplied")
	}

	// Create the user
	user := datastore.User{
		Username: username,
		Name:     name,
		Email:    email,
		Role:     datastore.Superuser,
	}
	_, err := db.CreateUser(user)
	return err
}

var parseFlags = func() {
	flag.StringVar(&username, "username", "", "Ubuntu SSO username of the user (https://login.ubuntu.com/)")
	flag.StringVar(&name, "name", "Super User", "Full name of the user")
	flag.StringVar(&email, "email", "user@example.com", "Email address of the user")
	flag.Parse()
}
