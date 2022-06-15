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

// Package management is the top-level and subcommands for the management binary\
package management

import (
	"fmt"
	"github.com/everactive/dmscore/versions"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/everactive/dmscore/iot-management/auth"
	"github.com/everactive/ginkeycloak"
	"github.com/spf13/cobra"

	"github.com/everactive/dmscore/iot-management/crypt"

	"github.com/everactive/dmscore/iot-management/config/configkey"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/spf13/viper"

	"github.com/everactive/dmscore/iot-management/config"
	"github.com/everactive/dmscore/iot-management/identityapi"
	"github.com/everactive/dmscore/iot-management/service/factory"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/everactive/dmscore/iot-management/twinapi"
	"github.com/everactive/dmscore/iot-management/web"
	log "github.com/sirupsen/logrus"
)

const (
	secretLength = 32
)

func init() {
	Command.AddCommand(version)
	Command.AddCommand(run)
}

// Command is the cobra.Command for the management service
var Command = &cobra.Command{
	Use: "management",
}

var version = &cobra.Command{
	Use:   "version",
	Short: "Prints the version to STDIN and exits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versions.Version)
	},
}

var run = &cobra.Command{
	Use:   "run",
	Short: "Runs the management server process",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versions.Version)
		main()
	},
}

func setupLogging() {
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
}

func setupJWTSecret(db datastore.DataStore) {
	jwtSecret := viper.GetString(configkey.JwtSecret)
	if len(jwtSecret) == 0 {
		// Create it the first time, then store it
		secret, errInt := crypt.CreateSecret(secretLength)
		if errInt != nil {
			log.Fatalf("Error generating JWT secret: %s", errInt)
			return
		}

		viper.Set(configkey.JwtSecret, secret)
		err := db.Set(configkey.JwtSecret, secret)
		if err != nil {
			log.Errorf("Unable to set jwt secret in datastore: %s", err)
		}
	} else {
		log.Tracef("JWT Secret already existed, using database stored value: %s", jwtSecret)
	}
}

func setupExternalAPIs(db datastore.DataStore) *manage.Management {
	deviceTwinAPIURL := viper.GetString(configkey.DeviceTwinAPIURL)
	// Initialize the device twin client
	twinAPI, err := twinapi.NewClientAdapter(deviceTwinAPIURL)
	if err != nil {
		log.Fatalf("Error connecting to the device twin service: %v", err)
	}

	identityAPIURL := viper.GetString(configkey.IdentityAPIURL)
	// Initialize the identity client
	idAPI, err := identityapi.NewClientAdapter(identityAPIURL)
	if err != nil {
		log.Fatalf("Error connecting to the identity service: %v", err)
	}

	// Create the main services
	srv := manage.NewManagement(db, twinAPI, idAPI)

	// Figure out what our auth provider is (keycloak or legacy static client token)
	authProvider := strings.ToLower(viper.GetString(configkey.AuthProvider))
	log.Infof("Auth provider: %s", authProvider)
	if authProvider == "static-client" {
		staticClientToken := viper.GetString(configkey.StaticClientToken)
		if staticClientToken != "" {
			auth.CreateServiceClientUser(db, "static-client")
			web.VerifyTokenAndUser = auth.VerifyStaticClientToken //nolint
		}
	} else if authProvider == "keycloak" {
		clientID := viper.GetString(configkey.OAuth2ClientID)
		secret := viper.GetString(configkey.OAuth2ClientSecret)
		host := viper.GetString(configkey.OAuth2HostName)
		port := viper.GetString(configkey.OAuth2HostPort)
		scheme := viper.GetString(configkey.OAuth2HostScheme)
		tokenIntrospectPath := viper.GetString(configkey.OAuth2TokenIntrospectPath)
		requiredScope := viper.GetString(configkey.OAuth2ClientRequiredScope)

		a := ginkeycloak.New(clientID, secret, host, port, scheme, requiredScope, tokenIntrospectPath, log.StandardLogger())
		web.VerifyTokenAndUser = auth.VerifyKeycloakTokenWithAuth(a)
	}

	return srv
}

func main() {
	setupLogging()

	var configFilePath string
	if filePath, ok := os.LookupEnv("CONFIG_FILE_PATH"); ok {
		configFilePath = filePath
	}

	config.LoadManagementConfig(configFilePath)

	// Open the connection to the local database
	databaseDriver := viper.GetString(configkey.DatabaseDriver)
	dataSource := viper.GetString(configkey.DatabaseConnectionString)
	log.Infof("Connecting to %s with connection string %s", databaseDriver, dataSource)
	db, err := factory.CreateDataStore(databaseDriver, dataSource)
	if err != nil || db == nil {
		log.Fatalf("Error accessing data store: %v, with database source %s", databaseDriver, dataSource)
		// This satisfies an IDE's nil logic check
		return
	}

	// Once we have a database connection, we need to load any settings that were persisted
	// in the database and make those available via Viper
	settings, err := db.GetSettings()
	if err != nil {
		log.Fatalf("Cannot continue without loading settings")
	}

	for _, setting := range settings {
		viper.Set(setting.Key, setting.Value)
	}

	setupJWTSecret(db)

	srv := setupExternalAPIs(db)

	tokenProvider := strings.ToLower(viper.GetString(configkey.ClientTokenProvider))
	log.Infof("Token provider is \"%s\"", tokenProvider)
	if tokenProvider == "keycloak" {
		tokenGetter := auth.TokenGetter()

		_, err = tokenGetter.GetToken()
		if err != nil {
			log.Errorf("Could not retrieve a token: %s", err.Error())
		}

		twinapi.AddAuthorization = addKeycloakAuthorizationToken
		identityapi.AddAuthorization = addKeycloakAuthorizationToken
	}

	// Start the web service
	www := web.NewService(srv)
	www.Run()
}

// Temporary function to avoid duplication but also avoid import cycles
func addKeycloakAuthorizationToken(req *resty.Request) error {
	token, err := auth.TokenGetter().GetToken()
	if err != nil {
		log.Error(err)
		return err
	}

	req.SetHeader("Auth-Type", "keycloak")
	req.SetAuthToken(token.KeycloakToken.AccessToken)
	return nil
}
