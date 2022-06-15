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

// Package web contains the code for handling the REST API
package web

import (
	"github.com/everactive/dmscore/iot-devicetwin/service/controller"
)

// Service is the implementation of the web API
type Service struct {
	port       string
	Controller controller.Controller
}

// NewService returns a new web controller
func NewService(port string, ctrl controller.Controller) *Service {
	return &Service{
		port:       port,
		Controller: ctrl,
	}
}
