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

package config

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/everactive/dmscore/config/keys"
	legacykeys "github.com/everactive/dmscore/iot-devicetwin/config/keys"

	"github.com/everactive/dmscore/iot-identity/service/cert"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var drivers = []string{"memory", "postgres"}

// MQTTConnect holds the credentials for MQTT connection
type MQTTConnect struct {
	ClientID   string
	RootCA     []byte
	ClientCert []byte
	ClientKey  []byte
}

var defaultValues = map[string]interface{}{
	legacykeys.ConfigPath:                    "/srv/config",
	legacykeys.DatabaseDriver:                "postgres",
	legacykeys.DatastoreSource:               "dbname=management host=localhost user=manager password=abc1234 sslmode=disable",
	legacykeys.MQTTClientCertificateFilename: "server.crt",
	legacykeys.MQTTClientKeyFilename:         "server.key",
	legacykeys.MQTTRootCAFilename:            "ca.crt",
	legacykeys.MQTTClientIDPrefix:            "devicetwin",
	legacykeys.MQTTHealthTopic:               "devices/health/+",
	legacykeys.MQTTPort:                      "8883",
	legacykeys.MQTTPubTopic:                  "devices/pub/+",
	legacykeys.MQTTURL:                       "localhost",
	legacykeys.ServicePort:                   "8040",
}

const (
	envPrefix = "IOTDEVICETWIN"
)

// LoadDeviceTwinConfig loads all configuration for the service, including base Viper LoadConfig
func LoadDeviceTwinConfig(configFilePath string) *MQTTConnect {
	LoadConfig(configFilePath)

	databaseDriver := viper.GetString(legacykeys.DatabaseDriver)
	found := false
	for i := range drivers {
		if drivers[i] == databaseDriver {
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("The database driver must be one of: %s", strings.Join(drivers, ", "))
	}

	certsDir := viper.GetString(keys.MQTTCertificatesPath)

	// Get the certificates for the MQTT broker
	m, err := readCerts(certsDir)
	if err != nil {
		log.Fatalf("Error reading certificates: %v", err)
	}

	return &m
}

// LoadConfig handles loading configuration for Viper using configuration file, environment variables and default values
func LoadConfig(configFilePath string) {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if len(configFilePath) > 0 {
		viper.SetConfigFile(configFilePath)
	} else {
		viper.SetConfigFile("config.yaml") // name of config file (without extension)
		viper.AddConfigPath(".")           // path to look for the config file in
	}

	// set defaults first
	for key, val := range defaultValues {
		viper.SetDefault(key, val)
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Warn("Config file not found, using defaults")
	}

	for _, key := range viper.AllKeys() {
		log.Tracef("%s = %+v", key, viper.Get(key))
	}
}

// readCerts reads the certificates from the file system
func readCerts(certsDir string) (MQTTConnect, error) {
	rootCAFilename := viper.GetString(legacykeys.MQTTRootCAFilename)
	clientCertFilename := viper.GetString(legacykeys.MQTTClientCertificateFilename)
	clientKeyFilename := viper.GetString(legacykeys.MQTTClientKeyFilename)

	c := MQTTConnect{}

	// nolint: gosec
	rootCA, err := ioutil.ReadFile(path.Join(certsDir, rootCAFilename))
	if err != nil {
		return c, err
	}

	// nolint: gosec
	certFile, err := ioutil.ReadFile(path.Join(certsDir, clientCertFilename))
	if err != nil {
		return c, err
	}

	// nolint: gosec
	key, err := ioutil.ReadFile(path.Join(certsDir, clientKeyFilename))

	c.RootCA = rootCA
	c.ClientKey = key
	c.ClientCert = certFile
	c.ClientID = generateClientID()

	return c, err
}

func generateClientID() string {
	prefix := viper.GetString(legacykeys.MQTTClientIDPrefix)

	// Generate a random string
	s, err := cert.CreateSecret(6)
	if err != nil {
		log.Printf("Error creating client ID: %v", err)
	}

	return fmt.Sprintf("%s-%s", prefix, s)
}
