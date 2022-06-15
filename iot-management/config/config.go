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

// Package config provides the functionality and types to load configuration from file or environment with the latter override the former
package config

import (
	"encoding/json"
	"fmt"
	"github.com/everactive/dmscore/versions"
	"os"
	"strings"

	"github.com/everactive/dmscore/iot-management/config/configkey"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var defaultValues = map[string]interface{}{
	configkey.DatabaseConnectionString:       "dbname=management host=localhost user=manager password=abc1234 sslmode=disable",
	configkey.ServiceScheme:                  "http",
	configkey.ServicePort:                    "8010",
	configkey.ServiceVersion:                 "0.0",
	configkey.ComponentVersionsCacheDuration: "60s",
}

type modelStoreID struct {
	Model   string `json:"model"`
	StoreID string `json:"storeid"`
}

const (
	envPrefix = "IOTMGMT"
	// ModelKeyTemplate is the template for getting a model key based on the model string
	ModelKeyTemplate = "store.model.%s"
)

var keys []string

func makeEnvironmentVariableKey(key string) string {
	k := strings.ToUpper(key)
	k = strings.ReplaceAll(k, ".", "_")
	k = envPrefix + "_" + k
	return k
}

// LoadManagementConfig loads all of the generic file and environment settings as well as service specific items
func LoadManagementConfig(configFilePath string) {
	keys = []string{
		configkey.DeviceTwinAPIURL,
		configkey.IdentityAPIURL,
		configkey.StoreURL,
		configkey.StoreIDs,

		configkey.DatabaseDriver,
		configkey.DatabaseConnectionString,
		configkey.StaticClientToken,

		configkey.ServicePort,
		configkey.ServiceScheme,
		configkey.ServiceHost,

		configkey.ClientTokenProvider,
		configkey.OAuth2ClientRequiredScope,
		configkey.OAuth2ClientID,
		configkey.OAuth2ClientSecret,
		configkey.OAuth2HostScheme,
		configkey.OAuth2HostName,
		configkey.OAuth2HostPort,
		configkey.OAuth2TokenIntrospectPath,
	}

	for _, k := range keys {
		envKey := makeEnvironmentVariableKey(k)
		log.Tracef("%s = %s", envKey, os.Getenv(envKey))
	}

	LoadConfig(configFilePath)

	modelStoreIds := []modelStoreID{}
	storeIDsString := viper.GetString(configkey.StoreIDs)
	log.Tracef("Store IDs string: %s", storeIDsString)
	if storeIDsString != "" {
		err := json.Unmarshal([]byte(storeIDsString), &modelStoreIds)
		if err != nil {
			log.Error(err)
		}
		for _, item := range modelStoreIds {
			log.Tracef("For model %s, loading store id %s", item.Model, item.StoreID)
			viper.Set(fmt.Sprintf(ModelKeyTemplate, item.Model), item.StoreID)
		}
	} else {
		log.Errorf("%s was empty or not set, could not load store ids", configkey.StoreIDs)
	}

	var componentVersions string
	versions := versions.GetComponentVersions()
	log.Tracef("%+v", versions)
	for name, version := range versions {
		componentVersions += fmt.Sprintf("%s : %s ", name, version)
	}

	viper.Set(configkey.ServiceVersion, componentVersions)
}

// LoadConfig loads all of the generic file and environment settings
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
