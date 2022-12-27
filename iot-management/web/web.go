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

// Package web provides types and functionality for the REST API and web app static files
package web

import (
	"context"
	"github.com/everactive/dmscore/config/keys"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/juju/usso"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"

	"github.com/gin-gonic/gin"
)

var ssoServer UbuntuSingleSignOnServer

// JSONHeader is the content-type header for JSON responses
const JSONHeader = "application/json; charset=UTF-8"

// Service is the implementation of the web API
type Service struct {
	Manage manage.Manage
	engine *gin.Engine
	runErr error
}

// NewService returns a new web controller
func NewService(srv manage.Manage) *Service {
	engine := gin.Default()

	s := &Service{
		Manage: srv,
		engine: engine,
	}

	AuthMiddleWare = s.authMiddleWare()

	s.router(engine)

	servicePort := viper.GetString(keys.ServicePort)
	log.Info("Starting service on port : ", servicePort)

	ssoServer = &usso.ProductionUbuntuSSOServer

	go func() {
		s.runErr = s.engine.Run(":" + servicePort)
	}()

	return s
}

// Run starts the web service
func (wb Service) Serve(ctx context.Context) error {
	intervalTicker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Infof("We're done: %s", ctx.Err())
			return wb.runErr
		case <-intervalTicker.C:
			log.Infof("%s still ticking", "DeviceTwinService")
		}
	}
}
