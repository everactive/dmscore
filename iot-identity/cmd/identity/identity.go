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

// Package identity contains the root and subcommands for the identity binary
package identity

import (
	"fmt"
	"os"
	"strings"

	"github.com/everactive/dmscore/iot-identity/version"

	"github.com/everactive/dmscore/iot-identity/datastore"

	"github.com/everactive/dmscore/iot-identity/config/configkey"
	"github.com/everactive/dmscore/iot-identity/service/factory"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-identity/config"
	"github.com/everactive/dmscore/iot-identity/service"
	"github.com/everactive/dmscore/iot-identity/web"

	"github.com/spf13/cobra"
)

func init() {
	Root.AddCommand(run)
	Root.AddCommand(versionCmd)
}

// Root is the top-level cobra.Command for the identity binary
var Root = &cobra.Command{
	Use:   "identity",
	Short: "The identity server",
}

var run = &cobra.Command{
	Use:   "run",
	Short: "Runs the main identity server process",
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version and exits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if len(logLevel) > 0 {
		l, err := log.ParseLevel(logLevel)
		if err != nil {
			log.SetLevel(log.TraceLevel)
			log.Tracef("LOG_LEVEL environment variable is set to %s, could not parse to a valid log level. Using trace logging.", logLevel)
		} else {
			log.SetLevel(l)
			log.Infof("Using LOG_LEVEL %s", logLevel)
		}
	}

	logFormat := os.Getenv("LOG_FORMAT")
	if strings.ToUpper(logFormat) == "JSON" {
		log.Infof("Using JSON log format")
		log.SetFormatter(&log.JSONFormatter{})
	}

	var configFilePath string
	if filePath, ok := os.LookupEnv("CONFIG_FILE_PATH"); ok {
		configFilePath = filePath
	}

	config.LoadIdentityConfig(configFilePath)

	// Open the connection to the local database
	databaseDriver := viper.GetString(configkey.DatabaseDriver)
	dataSource := viper.GetString(configkey.DatabaseConnectionString)
	log.Infof("Connecting to %s with connection string %s", databaseDriver, dataSource)
	datastore.Logger = log.StandardLogger()
	db, err := factory.CreateDataStore(databaseDriver, dataSource)
	if err != nil || db == nil {
		log.Fatalf("Error accessing data store: %v, with database source %s", databaseDriver, dataSource)
		// This satisfies an IDE's nil logic check
		return
	}

	service.Logger = log.StandardLogger()
	srv := service.NewIdentityService(db)

	// Start the web service
	web.Logger = log.StandardLogger()
	w := web.NewIdentityService(srv, log.StandardLogger())
	log.Fatal(w.Run())
}
