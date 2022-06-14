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

// Package config contains code and data to load configuration for the service
package config

import (
	"strings"

	"github.com/everactive/dmscore/iot-identity/config/configkey"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var defaultValues = map[string]interface{}{
	configkey.DatabaseDriver:           "postgres",
	configkey.DatabaseConnectionString: "dbname=identity host=localhost user=manager password=abc1234 sslmode=disable",
	configkey.ServicePortInternal:      "8030",
	configkey.ServicePortEnroll:        "8031",
	configkey.MQTTHostAddress:          "localhost",
	configkey.MQTTHostPort:             "8883",
	configkey.MQTTCertificatePath:      "/srv/certs",
}

const (
	envPrefix = "IOTIDENTITY"
)

// LoadIdentityConfig handles loading all configuration: file, environment and service specific
func LoadIdentityConfig(configFilePath string) {
	loadConfig(configFilePath)
}

func loadConfig(configFilePath string) {
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
