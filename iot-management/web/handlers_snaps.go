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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/everactive/dmscore/iot-devicetwin/web"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	snapRefreshURI = "refresh"
	snapEnableURI  = "enable"
	snapDisableURI = "disable"
	snapSwitchURI  = "switch"
)

// SnapListHandler fetches the list of installed snaps from the device
func (wb Service) SnapListHandler(c *gin.Context) {
	// Get the device snaps
	orgID := c.Param("orgid")
	deviceID := c.Param("deviceid")

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	response := wb.Manage.SnapList(orgID, user.Username, user.Role, deviceID)
	c.JSON(http.StatusOK, response)
}

// SnapListOnDevice lists snaps on the device
func (wb Service) SnapListOnDevice(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", JSONHeader)

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	response := wb.Manage.SnapListOnDevice(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"))
	err = encodeResponse(response, w)
	if err != nil {
		log.Error(err)
	}
}

// SnapInstallHandler installs a snap on the device
func (wb Service) SnapInstallHandler(c *gin.Context) { //nolint
	w := c.Writer
	w.Header().Set("Content-Type", JSONHeader)
	// Install a snap on a device
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	response := wb.Manage.SnapInstall(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"), c.Param("snap"))
	err = encodeResponse(response, w)
	if err != nil {
		log.Error(err)
	}
}

//nolint
// SnapDeleteHandler uninstalls a snap from the device
func (wb Service) SnapDeleteHandler(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", JSONHeader)

	// Uninstall a snap on a device
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	response := wb.Manage.SnapRemove(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"), c.Param("snap"))
	err = encodeResponse(response, w)
	if err != nil {
		log.Error(err)
	}
}

// SnapUpdateHandler updates a snap on the device
// Permitted actions are: enable, disable, refresh, or switch
func (wb Service) SnapUpdateHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	w.Header().Set("Content-Type", JSONHeader)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		formatStandardResponse("SnapUpdate", err.Error(), c)
		return
	}

	if len(body) == 0 {
		body = []byte("{}")
	}

	defer func(Body io.ReadCloser) {
		errInt := Body.Close()
		if errInt != nil {
			log.Error(errInt)
		}
	}(r.Body)

	var response web.StandardResponse

	switch c.Param("action") {
	case snapEnableURI, snapDisableURI, snapRefreshURI, snapSwitchURI:
		response = wb.Manage.SnapUpdate(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"), c.Param("snap"), c.Param("action"), body)
	default:
		w.WriteHeader(http.StatusBadRequest)
		response = web.StandardResponse{Code: "SnapUpdate", Message: fmt.Sprintf("Invalid action provided: %s", c.Param("action"))}
	}
	err = encodeResponse(response, w)
	if err != nil {
		log.Error(err)
	}
}

func (wb Service) snapSnapshotHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	w.Header().Set("Content-Type", JSONHeader)

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		formatStandardResponse("SnapPost", err.Error(), c)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response := web.StandardResponse{Code: "SnapSnapshot", Message: "Body is empty"}
		err2 := encodeResponse(response, w)
		if err2 != nil {
			log.Error(err2)
		}
		return
	}

	defer func(Body io.ReadCloser) {
		errInt := Body.Close()
		if errInt != nil {
			log.Error(errInt)
		}
	}(r.Body)

	response := wb.Manage.SnapSnapshot(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"), c.Param("snap"), body)

	err = encodeResponse(response, w)
	if err != nil {
		log.Error(err)
	}
}

// SnapConfigUpdateHandler gets the config for a snap on a device
func (wb Service) SnapConfigUpdateHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	w.Header().Set("Content-Type", JSONHeader)

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		formatStandardResponse("SnapUpdate", err.Error(), c)
		return
	}

	if len(body) == 0 {
		body = []byte("{}")
	}

	defer func(Body io.ReadCloser) {
		errInt := Body.Close()
		if errInt != nil {
			log.Error(errInt)
		}
	}(r.Body)

	// Update a snap's config on a device
	response := wb.Manage.SnapConfigSet(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"), c.Param("snap"), body)
	err = encodeResponse(response, w)
	if err != nil {
		log.Error(err)
	}
}

// SnapServiceAction start/stop/restart a snap on device
func (wb Service) SnapServiceAction(c *gin.Context) {
	w := c.Writer
	r := c.Request
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	w.Header().Set("Content-Type", JSONHeader)
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		formatStandardResponse("SnapUpdate", err.Error(), c)
		return
	}

	if len(body) == 0 {
		body = []byte("{}")
	}

	defer func(Body io.ReadCloser) {
		errInt := Body.Close()
		if errInt != nil {
			log.Error(errInt)
		}
	}(r.Body)

	// Update a snap's config on a device
	response := wb.Manage.SnapServiceAction(c.Param("orgid"), user.Username, user.Role, c.Param("deviceid"), c.Param("snap"), c.Param("action"), body)
	err = encodeResponse(response, w)
	if err != nil {
		log.Error(err)
	}
}
