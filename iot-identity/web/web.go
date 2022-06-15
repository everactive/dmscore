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
	"github.com/everactive/dmscore/iot-identity/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Logger is a logger specific to the web package and can be swapped out
var Logger = log.StandardLogger()

// Web is the interface for the web API
type Web interface {
	Run() error
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
