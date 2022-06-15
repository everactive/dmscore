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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/everactive/dmscore/iot-identity/service"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RegDeviceList is the API method to list the registered devices
func (wb Service) RegDeviceList(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", JSONHeader)

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	response := wb.Manage.RegDeviceList(c.Param("orgid"), user.Username, user.Role)
	if len(response.Code) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	_ = encodeResponse(response, w)
}

// RegDeviceGet is the API method to get a registered device
func (wb Service) RegDeviceGet(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", JSONHeader)

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	response := wb.Manage.RegDeviceGet(c.Param("orgid"), user.Username, user.Role, c.Param("device"))
	if len(response.Code) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	_ = encodeResponse(response, w)
}

// RegDeviceGetDownload provides the download of the device data
func (wb Service) RegDeviceGetDownload(c *gin.Context) {
	w := c.Writer
	r := c.Request
	w.Header().Set("Content-Type", JSONHeader)

	// Fetch the device
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	response := wb.Manage.RegDeviceGet(c.Param("orgid"), user.Username, user.Role, c.Param("device"))
	if len(response.Code) > 0 {
		formatStandardResponse(response.Code, response.Message, c)
		return
	}

	// Set the download headers and body for the file
	w.Header().Set("Content-Disposition", "attachment; filename=devicedata-"+response.Enrollment.ID)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	// Decode the base64 file
	data, err := base64.StdEncoding.DecodeString(response.Enrollment.DeviceData)
	if err != nil {
		log.Println("Error decoding the device data:", err)
		formatStandardResponse("DeviceData", "Error decoding the file", c)
		return
	}

	_, err = io.Copy(w, bytes.NewReader(data))
	if err != nil {
		log.Error(err)
	}
}

// RegDeviceUpdate is the API method to update a registered device status
func (wb Service) RegDeviceUpdate(c *gin.Context) {
	w := c.Writer
	r := c.Request
	w.Header().Set("Content-Type", JSONHeader)

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	defer func(Body io.ReadCloser) {
		errInt := Body.Close()
		if errInt != nil {
			log.Error(errInt)
		}
	}(r.Body)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		formatStandardResponse("RegDevice", "error reading the request", c)
		return
	}

	// Get the devices
	response := wb.Manage.RegDeviceUpdate(c.Param("orgid"), user.Username, user.Role, c.Param("device"), b)
	if len(response.Code) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	_ = encodeResponse(response, w)
}

// RegisterDevice registers a new device with the Identity service
func (wb Service) RegisterDevice(c *gin.Context) {
	w := c.Writer
	r := c.Request
	w.Header().Set("Content-Type", JSONHeader)
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Admin)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	// Read the body of the request

	defer func(Body io.ReadCloser) {
		errInt := Body.Close()
		if errInt != nil {
			log.Error(errInt)
		}
	}(r.Body)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		formatStandardResponse("RegDevice", "error reading the request", c)
		return
	}

	// Decode the body and check we have a valid organization ID
	req, err := decodeDeviceRequest(b)
	if err != nil {
		formatStandardResponse("RegDevice", "error decoding the request", c)
		return
	}
	if req.OrganizationID != c.Param("orgid") {
		formatStandardResponse("RegDevice", "the request organization ID is invalid", c)
		return
	}

	// Register the devices
	response := wb.Manage.RegisterDevice(c.Param("orgid"), user.Username, user.Role, b)
	if len(response.Code) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	_ = encodeResponse(response, w)
}

func decodeDeviceRequest(body []byte) (*service.RegisterDeviceRequest, error) {
	// Decode the JSON body
	dev := service.RegisterDeviceRequest{}
	err := json.Unmarshal(body, &dev)
	return &dev, err
}
