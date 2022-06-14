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
	"encoding/json"
	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	"net/http"

	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/gin-gonic/gin"

	dtwin "github.com/everactive/dmscore/iot-devicetwin/web"
	log "github.com/sirupsen/logrus"
)

func formatStandardResponse(errorCode, message string, c *gin.Context) {
	response := dtwin.StandardResponse{Code: errorCode, Message: message}
	if len(errorCode) > 0 {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DevicesListHandler is the API method to list the registered devices
func (wb Service) DevicesListHandler(c *gin.Context) {
	w := c.Writer

	w.Header().Set("Content-Type", JSONHeader)
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	// Get the devices
	response := wb.Manage.DeviceList(c.Param("orgid"), user.Username, user.Role)
	log.Tracef("Sending response back: %+v", response)
	_ = encodeResponse(response, w)

}

// DeviceDeleteHandler is the API method to delete a registered device
func (wb Service) DeviceDeleteHandler(c *gin.Context) {
	w := c.Writer

	w.Header().Set("Content-Type", JSONHeader)
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Admin)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	// Delete the device
	response := wb.Manage.DeviceDelete(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"))
	_ = encodeResponse(response, w)
}

// DeviceGetHandler is the API method to get a registered device
func (wb Service) DeviceGetHandler(c *gin.Context) {
	w := c.Writer

	w.Header().Set("Content-Type", JSONHeader)
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Admin)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	// Get the device
	response := wb.Manage.DeviceGet(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"))
	_ = encodeResponse(response, w)
}

// ActionListHandler is the API method to get actions for a device
func (wb Service) ActionListHandler(c *gin.Context) {
	w := c.Writer

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Admin)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	// Get the device
	response := wb.Manage.ActionList(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"))
	_ = encodeResponse(response, w)
}

//nolint
// DeviceLogsHandler is the API method to get logs for a device
func (wb Service) DeviceLogsHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	w.Header().Set("Content-Type", JSONHeader)
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Admin)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	deviceLogs := messages.DeviceLogs{}
	err = json.NewDecoder(r.Body).Decode(&deviceLogs)
	if err != nil {
		formatStandardResponse("DeviceLogs", err.Error(), c)
		return
	}

	// Get the device
	response := wb.Manage.DeviceLogs(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"), &deviceLogs)
	_ = encodeResponse(response, w)
}

//nolint
// DeviceUsersActionHandler is the API method to create a user for a device
func (wb Service) DeviceUsersActionHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	w.Header().Set("Content-Type", JSONHeader)
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Admin)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	deviceUser := messages.DeviceUser{}
	err = json.NewDecoder(r.Body).Decode(&deviceUser)
	if err != nil {
		formatStandardResponse("DeviceUsersAction", err.Error(), c)
		return
	}

	// Get the device
	response := wb.Manage.DeviceUsersAction(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"), deviceUser)
	_ = encodeResponse(response, w)
}
