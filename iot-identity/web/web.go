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

// Package web implements the REST API handling and routing
package web

import (
	"net/http"
	"os"
	"strings"

	"github.com/everactive/dmscore/iot-identity/auth"

	"github.com/everactive/dmscore/iot-identity/config/configkey"
	"github.com/everactive/dmscore/iot-identity/middleware/logger"
	"github.com/everactive/dmscore/iot-identity/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

// Logger is a logger specific to the web package and can be swapped out
var Logger = log.StandardLogger()

// Web is the interface for the web API
type Web interface {
	Run() error
	InternalRouter(engine *gin.Engine)
	EnrollRouter(engine *gin.Engine)
	RegisterOrganization(w http.ResponseWriter, r *http.Request)
	RegisterDevice(w http.ResponseWriter, r *http.Request)
	OrganizationList(w http.ResponseWriter, r *http.Request)
	DeviceList(w http.ResponseWriter, r *http.Request)

	EnrollDevice(w http.ResponseWriter, r *http.Request)
}

// IdentityService is the implementation of the web API
type IdentityService struct {
	Identity service.Identity
	logger   *log.Logger
}

// NewIdentityService returns a new web controller
func NewIdentityService(id service.Identity, l *log.Logger) *IdentityService {
	return &IdentityService{
		Identity: id,
		logger:   l,
	}
}

// Run starts the web service
func (wb IdentityService) Run() error {
	internalPort := viper.GetString(configkey.ServicePortInternal)
	enrollPort := viper.GetString(configkey.ServicePortEnroll)

	log.Info("Starting service (internal) on port : ", internalPort)
	log.Info("Starting service (enroll) on port : ", enrollPort)

	internalRouter := gin.New()
	enrollRouter := gin.New()

	logFormat := os.Getenv("LOG_FORMAT")
	if strings.ToUpper(logFormat) == "JSON" {
		log.Infof("Setting up JSON log format for logger middleware")

		middlewareLogger := logger.New(log.StandardLogger(), logger.LogOptions{EnableStarting: true})

		internalRouter.Use(middlewareLogger.HandleFunc)
		enrollRouter.Use(middlewareLogger.HandleFunc)

	} else {
		internalRouter.Use(gin.Logger())
		enrollRouter.Use(gin.Logger())
	}

	internalRouter.Use(auth.Factory(viper.GetString(configkey.AuthProvider)))

	wb.internalRouter(internalRouter)
	wb.enrollRouter(enrollRouter)

	// Use a goroutine for the internal serve, we'll block with the enroll serve
	go func() {
		log.Info("Listening and serving internal on :" + internalPort)

		err := internalRouter.Run(":" + internalPort)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Info("Listening and serving enroll on :" + enrollPort)

	err := enrollRouter.Run(":" + enrollPort)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
