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

package web

import (
	"github.com/everactive/dmscore/api"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var AuthMiddleWare gin.HandlerFunc

func (wb Service) authMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := wb.checkIsStandardAndGetUserFromJWT(ctx)
		if err != nil {
			log.Error(err)
			response := api.StandardResponse{Code: "UserAuth", Message: err.Error()}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &response)
			return
		}

		ctx.Set("USER", &user)
	}
}

func (wb Service) addWebAppGroup(webAppGroup *gin.RouterGroup) {
	// OpenID routes: using Ubuntu SSO
	webAppGroup.GET("/login", wb.LoginHandler)
	webAppGroup.GET("/logout", wb.LogoutHandler)
}

func (wb Service) addAPI(apiRouter *gin.RouterGroup) {
	apiRouter.Use(AuthMiddleWare)

	// API routes: accounts
	apiRouter.GET("/organizations", wb.OrganizationListHandler)
	apiRouter.GET("/organizations/:id", wb.OrganizationGetHandler)
	apiRouter.PUT("/organizations/:id", wb.OrganizationUpdateHandler)
	apiRouter.POST("/organizations", wb.OrganizationCreateHandler)

	// API routes: users
	apiRouter.GET("/users", wb.UserListHandler)
	apiRouter.POST("/users", wb.UserCreateHandler)
	apiRouter.GET("/users/:username", wb.UserGetHandler)
	apiRouter.PUT("/users/:username", wb.UserUpdateHandler)
	apiRouter.DELETE("/users/:username", wb.UserDeleteHandler)
	apiRouter.GET("/users/:username/organizations", wb.OrganizationsForUserHandler)
	apiRouter.POST("/users/:username/organizations/:orgid", wb.OrganizationUpdateForUserHandler)

	//// API routes: registered devices
	apiRouter.GET("/:orgid/register/devices", wb.RegDeviceList)
	apiRouter.POST("/:orgid/register/devices", wb.RegisterDevice)
	apiRouter.GET("/:orgid/register/devices/:device", wb.RegDeviceGet)
	apiRouter.PUT("/:orgid/register/devices/:device", wb.RegDeviceUpdate)
	apiRouter.GET("/:orgid/register/devices/:device/download", wb.RegDeviceGetDownload)

	//// API routes: devices
	apiRouter.GET("/:orgid/devices", wb.DevicesListHandler)
	apiRouter.GET("/:orgid/devices/:deviceid", wb.DeviceGetHandler)
	apiRouter.GET("/:orgid/devices/:deviceid/actions", wb.ActionListHandler)
	apiRouter.DELETE("/:orgid/devices/:deviceid", wb.DeviceDeleteHandler)
	apiRouter.POST("/:orgid/devices/:deviceid/logs", wb.DeviceLogsHandler)
	apiRouter.POST("/:orgid/devices/:deviceid/users", wb.DeviceUsersActionHandler)

	//// API routes: snap functionality
	apiRouter.GET("/device/:orgid/:deviceid/snaps", wb.SnapListHandler)

	apiRouter.POST("/snaps/:orgid/:deviceid/list", wb.SnapListOnDevice)
	apiRouter.POST("/snaps/:orgid/:deviceid/:snap", wb.SnapInstallHandler)
	apiRouter.DELETE("/snaps/:orgid/:deviceid/:snap", wb.SnapDeleteHandler)
	apiRouter.PUT("/snaps/:orgid/:deviceid/:snap/settings", wb.SnapConfigUpdateHandler)
	apiRouter.POST("/snaps/:orgid/:deviceid/services/:snap/:action", wb.SnapServiceAction)
	apiRouter.PUT("/snaps/:orgid/:deviceid/:snap/:action", wb.SnapUpdateHandler)
	apiRouter.POST("/snaps/:orgid/:deviceid/:snap/snapshot", wb.snapSnapshotHandler)

	//// API routes: store functionality
	apiRouter.GET("/store/snaps/:model/:snapName", wb.StoreSearchHandler)
}

func (wb Service) router(e *gin.Engine) {
	wb.addWebAppGroup(e.Group("/"))

	nonAuthAPIGroup := e.Group("/v1")
	{
		// API routes: CSRF token and auth token
		nonAuthAPIGroup.GET("/token", wb.TokenHandler)
		nonAuthAPIGroup.GET("/authtoken", wb.TokenHandler)
		nonAuthAPIGroup.GET("/versions", wb.VersionHandler)

		// API routes: login
		nonAuthAPIGroup.Any("/login", wb.LoginHandlerAPIClient)
	}

	wb.addAPI(e.Group("/v1"))
}
