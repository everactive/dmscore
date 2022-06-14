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

package web

import (
	"github.com/gin-gonic/gin"
)

func (wb IdentityService) internalRouter(engine *gin.Engine) {
	engine.POST("/v1/organization", wb.RegisterOrganization)
	engine.GET("/v1/organizations", wb.OrganizationList)
	engine.POST("/v1/device", wb.RegisterDevice)
	engine.DELETE("/v1/device/:deviceid", wb.DeleteDevice)
	engine.GET("/v1/devices/:orgid", wb.DeviceList)
	engine.GET("/v1/devices/:orgid/:device", wb.DeviceGet)
	engine.PUT("/v1/devices/:orgid/:device", wb.DeviceUpdate)
}

func (wb IdentityService) enrollRouter(engine *gin.Engine) {
	engine.POST("/v1/device/enroll", wb.EnrollDevice)
}

func (wb IdentityService) SetRouters(internal, enroll *gin.Engine) {
	wb.internalRouter(internal)
	wb.enrollRouter(enroll)
}
