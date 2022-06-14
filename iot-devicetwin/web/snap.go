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
	"encoding/json"
	"io/ioutil"

	"net/http"

	"github.com/everactive/dmscore/iot-devicetwin/pkg/messages"
	log "github.com/sirupsen/logrus"
)

// SnapList is the API call to list snaps for a device
func (wb Service) SnapList(w http.ResponseWriter, r *http.Request, vars varLookup) {
	orgID := vars("orgid")
	deviceID := vars("id")
	snap := vars("snap")
	action := vars("action")

	log.Tracef("SnapList: orgid=%s device id=%s, snap=%s, action=%s",
		orgID, deviceID, snap, action)

	installed, err := wb.Controller.DeviceSnaps(orgID, deviceID)
	if err != nil {
		log.Println("Error fetching snaps for a device:", err)
		formatStandardResponse("SnapList", "Error fetching snaps for the device", w)
		return
	}

	formatSnapsResponse(installed, w)
}

// SnapListPublish is the API call to trigger a snap list on a device
func (wb Service) SnapListPublish(w http.ResponseWriter, r *http.Request, vars varLookup) {

	if err := wb.Controller.DeviceSnapList(vars("orgid"), vars("id")); err != nil {
		log.Println("Error requesting snap list for the device:", err)
		formatStandardResponse("SnapList", "Error requesting snap list for the device", w)
		return
	}

	formatStandardResponse("", "", w)
}

// SnapInstall is the API call to install a snap for a device
func (wb Service) SnapInstall(w http.ResponseWriter, r *http.Request, vars varLookup) {
	if err := wb.Controller.DeviceSnapInstall(vars("orgid"), vars("id"), vars("snap")); err != nil {
		log.Println("Error requesting snap install for the device:", err)
		formatStandardResponse("SnapInstall", "Error requesting snap install for the device", w)
		return
	}

	formatStandardResponse("", "", w)
}

// SnapServiceAction is the API call to start,stop, or restart a snap for a device
func (wb Service) SnapServiceAction(w http.ResponseWriter, r *http.Request, vars varLookup) {
	if r == nil {
		log.Error("error in json decoding for SnapServiceAction: nil request")
		formatStandardResponse("SnapServiceAction", "invalid request", w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		body = []byte("{}")
	}

	var services *messages.SnapService
	err = json.Unmarshal(body, &services)

	if err != nil {
		log.Error("error in json decoding for SnapServiceAction: ", err)
		formatStandardResponse("SnapServiceAction", "invalid json", w)
		return
	}

	if err := wb.Controller.DeviceSnapServiceAction(vars("orgid"), vars("id"), vars("snap"), vars("action"), services); err != nil {
		log.Println("Error requesting snap start for the device:", err)
		formatStandardResponse("SnapServiceAction", "Error requesting snap start for the device", w)
		return
	}

	formatStandardResponse("", "", w)
}

// SnapRemove is the API call to uninstall a snap for a device
func (wb Service) SnapRemove(w http.ResponseWriter, r *http.Request, vars varLookup) {

	if err := wb.Controller.DeviceSnapRemove(vars("orgid"), vars("id"), vars("snap")); err != nil {
		log.Println("Error requesting snap remove for the device:", err)
		formatStandardResponse("SnapRemove", "Error requesting snap remove for the device", w)
		return
	}

	formatStandardResponse("", "", w)
}

// SnapUpdateAction is the API call to update a snap for a device (enable, disable, refresh)
func (wb Service) SnapUpdateAction(w http.ResponseWriter, r *http.Request, vars varLookup) {

	orgID := vars("orgid")
	deviceID := vars("id")
	snap := vars("snap")
	action := vars("action")

	var snapUpdate *messages.SnapUpdate
	var bodyString string
	if r.ContentLength > 0 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading snap update action body:", err)
			formatStandardResponse("SnapUpdate", "Error reading body", w)
			return
		}
		defer func() {
			err := r.Body.Close()
			if err != nil {
				log.Printf("Error trying r.Body.Close(): %+v", err)
			}
		}()

		if len(body) > 0 {
			bodyString = string(body)
			var snapUpdateLoc messages.SnapUpdate
			log.Tracef("Body: %s", string(body))
			err := json.Unmarshal(body, &snapUpdateLoc)
			if err != nil {
				log.Error(err)
				formatStandardResponse("SnapUpdate", "error trying to unmarshal string: "+string(body), w)
				return
			}

			snapUpdate = &snapUpdateLoc
		}
	}

	log.Tracef("SnapUpdateAction: orgid=%s device id=%s, snap=%s, action=%s, data=%s",
		orgID, deviceID, snap, action, bodyString)

	if err := wb.Controller.DeviceSnapUpdate(orgID, deviceID, snap, action, snapUpdate); err != nil {
		log.Println("Error requesting snap update for the device:", err)
		formatStandardResponse("SnapUpdate", "Error requesting snap update for the device", w)
		return
	}

	formatStandardResponse("", "", w)
}

// SnapUpdateConf is the API call to update a snap for a device (settings)
func (wb Service) SnapUpdateConf(w http.ResponseWriter, r *http.Request, vars varLookup) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading snap config body:", err)
		formatStandardResponse("SnapSetConf", "Error requesting snap settings update for the device", w)
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("Error trying r.Body.Close(): %+v", err)
		}
	}()

	if err := wb.Controller.DeviceSnapConf(vars("orgid"), vars("id"), vars("snap"), string(body)); err != nil {
		log.Println("Error requesting snap settings update for the device:", err)
		formatStandardResponse("SnapSetConf", "Error requesting snap settings update for the device", w)
		return
	}

	formatStandardResponse("", "", w)
}

// SnapSnapshot is the API call to upload a snapshot of a snap to S3 storage
func (wb Service) SnapSnapshot(w http.ResponseWriter, r *http.Request, vars varLookup) {

	if r == nil {
		log.Error("error in json decoding for SnapSnapshot: nil request")
		formatStandardResponse("SnapSnapshot", "invalid request", w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		log.Error("missing expected body to request")
		formatStandardResponse("SnapSnapshot", "empty body", w)
		return
	}

	var data *messages.SnapSnapshot
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Error("error in json decoding for SnapSnapshot: ", err)
		formatStandardResponse("SnapSnapshot", "invalid json", w)
		return
	}

	if data.Url == "" {
		log.Error("error in JSON body: missing url field")
		formatStandardResponse("SnapSnapshot", "invalid json", w)
		return
	}

	if err := wb.Controller.DeviceSnapSnapshot(vars("orgid"), vars("id"), vars("snap"), data); err != nil {
		log.Println("Error requesting snaphot for the device:", err)
		formatStandardResponse("SnapSnapshot", "Error requesting snapshot for the device", w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", JSONHeader)
	response := StandardResponse{
		Message: "Snap snapshot request sent",
	}

	// Encode the response as JSON
	encodeResponse(w, response)
}
