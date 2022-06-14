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

package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// wrapF is a helper function for wrapping VarHandlerFunc and returns a Gin middleware.
func wrapF(f varHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c.Writer, c.Request, func(k string) string {
			return c.Param(k)
		})
	}
}

type varLookup func(string) string

type varHandlerFunc func(http.ResponseWriter, *http.Request, varLookup)

// AddRoutes adds the routes to the gin Engine
func (wb Service) AddRoutes(engine *gin.Engine) {

	// Actions on a device twin
	engine.GET("/v1/device/:orgid/:id/snaps", wrapF(wb.SnapList))
	engine.GET("/v1/device/:orgid", wrapF(wb.DeviceList))
	engine.GET("/v1/device/:orgid/:id", wrapF(wb.DeviceGet))
	engine.GET("/v1/device/:orgid/:id/actions", wrapF(wb.ActionList))

	// Actions on a device
	engine.DELETE("/v1/device/:orgid/:id", wrapF(wb.DeviceDelete))
	engine.POST("/v1/device/:orgid/:id/snaps/list", wrapF(wb.SnapListPublish))
	engine.POST("/v1/device/:orgid/:id/snaps/:snap", wrapF(wb.SnapInstall))
	engine.DELETE("/v1/device/:orgid/:id/snaps/:snap", wrapF(wb.SnapRemove))
	engine.PUT("/v1/device/:orgid/:id/snaps/:snap/settings", wrapF(wb.SnapUpdateConf))
	engine.PUT("/v1/device/:orgid/:id/snaps/:snap/:action", wrapF(wb.SnapUpdateAction))
	engine.POST("/v1/device/:orgid/:id/snaps/:snap/snapshot", wrapF(wb.SnapSnapshot))
	engine.POST("/v1/device/:orgid/:id/services/:snap/:action", wrapF(wb.SnapServiceAction))
	engine.POST("/v1/device/:orgid/:id/logs", wrapF(wb.DeviceLogs))

	engine.POST("/v1/device/:orgid/:id/users", wrapF(wb.UserAdd))
	engine.DELETE("/v1/device/:orgid/:id/users", wrapF(wb.UserRemove))

	// Actions on a group
	engine.POST("/v1/group/:orgid", wrapF(wb.GroupCreate))
	engine.GET("/v1/group/:orgid", wrapF(wb.GroupList))
	engine.GET("/v1/group/:orgid/:name", wrapF(wb.GroupGet))
	engine.POST("/v1/group/:orgid/:name/:id", wrapF(wb.GroupLinkDevice))
	engine.DELETE("/v1/group/:orgid/:name/:id", wrapF(wb.GroupUnlinkDevice))
	engine.GET("/v1/group/:orgid/:name/devices", wrapF(wb.GroupGetDevices))
	engine.GET("/v1/group/:orgid/:name/devices/excluded", wrapF(wb.GroupGetExcludedDevices))
}
