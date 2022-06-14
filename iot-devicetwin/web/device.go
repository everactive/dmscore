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
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
)

// DeviceGet is the API call to get a device
func (wb Service) DeviceGet(w http.ResponseWriter, r *http.Request, vars varLookup) {

	device, err := wb.Controller.DeviceGet(vars("orgid"), vars("id"))
	if err != nil {
		log.Printf("Error fetching the device `%s`: %v", vars("id"), err)
		formatStandardResponse("DeviceGet", "Error fetching the device", w)
		return
	}

	formatDeviceResponse(device, w)
}

// DeviceList is the API call to get devices
func (wb Service) DeviceList(w http.ResponseWriter, r *http.Request, vars varLookup) {

	devices, err := wb.Controller.DeviceList(vars("orgid"))
	if err != nil {
		log.Printf("Error fetching the device list for `%s`: %v", vars("orgid"), err)
		formatStandardResponse("DeviceList", "Error fetching devices", w)
		return
	}

	formatDevicesResponse(devices, w)
}

// DeviceDelete is the API call to delete devices
func (wb Service) DeviceDelete(w http.ResponseWriter, r *http.Request, vars varLookup) {

	deviceID := vars("id")
	orgID := vars("orgid")

	_, err := wb.Controller.DeviceGet(orgID, deviceID)
	if errors.As(err, &sql.ErrNoRows) {
		log.Warnf("When deleting device %s, it was not found (never connected).", deviceID)
		formatStandardResponse("DeviceDelete", "No device found", w)
		return
	}

	err = wb.Controller.DeviceDelete(deviceID)
	if err != nil {
		log.Printf("Error deleting the device `%s`: %v", deviceID, err)
		formatStandardResponse("DeviceDelete", "Error deleting device", w)
		return
	}

	w.Header().Set("Content-Type", JSONHeader)
	response := StandardResponse{
		Message: "device deleted",
	}

	// Encode the response as JSON
	encodeResponse(w, response)
}

// DeviceLogs is the API call to upload syslog logs from a device to S3
func (wb Service) DeviceLogs(w http.ResponseWriter, r *http.Request, vars varLookup) {

	if r == nil {
		log.Error("error in json decoding for DeviceLogs: nil request")
		formatStandardResponse("DeviceLogs", "invalid request", w)
		return
	}

	var data *messages.DeviceLogs
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		log.Error("error in json decoding for DeviceLogs: ", err)
		formatStandardResponse("DeviceLogs", "invalid json", w)
		return
	}

	if data.Url == "" {
		log.Error("error in JSON body: missing url field")
		formatStandardResponse("DeviceLogs", "invalid json", w)
		return
	}

	err = wb.Controller.DeviceLogs(vars("orgid"), vars("id"), data)
	if err != nil {
		log.Printf("Error pulling logs from device `%s`: %v", vars("id"), err)
		formatStandardResponse("DeviceLogs", "Error pulling logs from device", w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", JSONHeader)
	response := StandardResponse{
		Message: "device logs request sent",
	}

	// Encode the response as JSON
	encodeResponse(w, response)
}
