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
	"context"
	"github.com/everactive/dmscore/config/keys"
	"github.com/everactive/dmscore/iot-identity/middleware/logger"
	"github.com/everactive/dmscore/iot-identity/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// Logger is a logger specific to the web package and can be swapped out
var Logger = log.StandardLogger()

// IdentityService is the implementation of the web API
type IdentityService struct {
	Identity service.Identity
	logger   *log.Logger
	enrollRouter *gin.Engine
}

// NewIdentityService returns a new web controller
func NewIdentityService(id service.Identity, l *log.Logger) *IdentityService {
	enrollRouter := gin.New()

	logFormat := os.Getenv("LOG_FORMAT")
	if strings.ToUpper(logFormat) == "JSON" {
		log.Infof("Setting up JSON log format for logger middleware")

		middlewareLogger := logger.New(log.StandardLogger(), logger.LogOptions{EnableStarting: true})

		enrollRouter.Use(middlewareLogger.HandleFunc)

	} else {
		enrollRouter.Use(gin.Logger())
	}

	i := &IdentityService{
		Identity: id,
		logger:   l,
		enrollRouter: enrollRouter,
	}

	enrollRouter.POST("/v1/device/enroll", i.EnrollDevice)

	return i
}

func (i *IdentityService) Serve(ctx context.Context) error {

	enrollPort := viper.GetString(keys.GetIdentityKey(keys.ServicePortEnroll))
	log.Info("Starting service (enroll) on port : ", enrollPort)

	log.Info("Listening and serving enroll on :" + enrollPort)

	err := i.enrollRouter.Run(":" + enrollPort)

	<- ctx.Done()

	return err
}