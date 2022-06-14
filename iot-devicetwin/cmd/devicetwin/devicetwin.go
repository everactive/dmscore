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

// Package devicetwin provides the root command and subcommands for the devicetwin binary
package devicetwin

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/everactive/dmscore/iot-devicetwin/version"

	"github.com/everactive/dmscore/iot-devicetwin/config/keys"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-devicetwin/config"
	"github.com/everactive/dmscore/iot-devicetwin/service/controller"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	"github.com/everactive/dmscore/iot-devicetwin/service/factory"
	"github.com/everactive/dmscore/iot-devicetwin/web"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/spf13/cobra"

	// this is needed for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// this is needed for migrate
	_ "github.com/golang-migrate/migrate/v4/source/github"
	// this is needed for migrate
	_ "github.com/lib/pq"
)

func init() {
	Root.AddCommand(run)
	Root.AddCommand(versionCmd)
}

// Root is the top-level cobra.Command for the identity binary
var Root = &cobra.Command{
	Use:   "devicetwin",
	Short: "The devicetwin server",
}

var run = &cobra.Command{
	Use:   "run",
	Short: "Runs the main devicetwin server process",
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
	log.SetLevel(log.TraceLevel)
	log.Println("Server starting")

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
	if logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	var configFilePath string
	if filePath, ok := os.LookupEnv("CONFIG_FILE_PATH"); ok {
		configFilePath = filePath
	}

	connect := config.LoadDeviceTwinConfig(configFilePath)

	databaseDriver := viper.GetString(keys.DatabaseDriver)
	dataStoreSource := viper.GetString(keys.DatastoreSource)

	// Run migrations before doing anything else
	err := runMigrations(dataStoreSource)

	// Cannot continue if migrations fail for some reason
	if err != nil {
		log.Error(err)
		return
	}

	db, err := factory.CreateDataStore(databaseDriver, dataStoreSource)
	if err != nil {
		log.Fatalf("Error connecting to data store: %v", err)
	}

	URL := viper.GetString(keys.MQTTURL)
	port := viper.GetString(keys.MQTTPort)

	twin := devicetwin.NewService(db)
	ctrl := controller.NewService(URL, port, connect, twin)

	servicePort := viper.GetString(keys.ServicePort)

	// Start the web API service
	w := web.NewService(servicePort, ctrl)
	log.Fatal(w.Run())
}

func runMigrations(datasource string) error {
	db, err := sql.Open("postgres", datasource)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres",
		driver)

	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
		return err
	}

	return nil
}
