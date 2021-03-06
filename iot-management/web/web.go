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
	"github.com/everactive/dmscore/iot-management/config/configkey"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/juju/usso"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

var ssoServer UbuntuSingleSignOnServer

// JSONHeader is the content-type header for JSON responses
const JSONHeader = "application/json; charset=UTF-8"

// Service is the implementation of the web API
type Service struct {
	Manage manage.Manage
}

// NewService returns a new web controller
func NewService(srv manage.Manage) *Service {
	s := &Service{
		Manage: srv,
	}

	AuthMiddleWare = s.authMiddleWare()

	return s
}

// Run starts the web service
func (wb Service) Run() {
	servicePort := viper.GetString(configkey.ServicePort)
	log.Info("Starting service on port : ", servicePort)

	ssoServer = &usso.ProductionUbuntuSSOServer

	r := gin.Default()

	r.Use(gin.Logger())

	wb.router(r)

	err := r.Run(":" + servicePort)
	if err != nil {
		panic(err)
	}
}
