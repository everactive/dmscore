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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/everactive/dmscore/iot-identity/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/snapcore/snapd/asserts"
)

type devicesField struct {
	ID     string
	Serial string
}

// DeviceList fetches device registrations
func (wb IdentityService) DeviceList(context *gin.Context) {
	devices, err := wb.Identity.DeviceList(context.Param("orgid"))
	if err != nil {
		log.Println("Error fetching devices:", err)
		formatStandardResponse("DeviceList", err.Error(), context.Writer)
		return
	}

	devicesFields := []devicesField{}
	for _, d := range devices {
		devicesFields = append(devicesFields, devicesField{
			ID:     d.ID,
			Serial: d.Device.SerialNumber,
		})
	}
	wb.logger.WithField("devices", devicesFields)

	formatDevicesResponse(devices, context.Writer)
}

// DeviceGet fetches a device registration
func (wb IdentityService) DeviceGet(context *gin.Context) {
	en, err := wb.Identity.DeviceGet(context.Param("orgid"), context.Param("device"))
	if err != nil {
		log.Printf("Error fetching device `%s`: %v\n", context.Param("device"), err)
		formatStandardResponse("DeviceGet", err.Error(), context.Writer)
		return
	}
	formatEnrollResponse(*en, context.Writer)
}

// DeviceUpdate updates a device registration
func (wb IdentityService) DeviceUpdate(context *gin.Context) {
	orgid := context.Param("orgid")
	device := context.Param("device")
	request := &service.DeviceUpdateRequest{}
	err := decodeRequest(context.Writer, context.Request, request)
	if err != nil {
		return
	}

	err = wb.Identity.DeviceUpdate(orgid, device, request)
	if err != nil {
		log.Printf("Error updating device `%s`: %v\n", device, err)
		formatStandardResponse("DeviceUpdate", err.Error(), context.Writer)
		return
	}
	formatStandardResponse("", "", context.Writer)
}

// DeleteDevice unregisters a new device with the identity service
func (wb IdentityService) DeleteDevice(context *gin.Context) {
	id, err := wb.Identity.DeleteDevice(context.Param("deviceid"))
	if err != nil {
		log.Println("Error deleting device:", err)
		formatStandardResponse("DeleteDevice", err.Error(), context.Writer)
		return
	}
	formatRegisterResponse(id, context.Writer)
}

// RegisterDevice registers a new device with the identity service
func (wb IdentityService) RegisterDevice(context *gin.Context) {
	// Decode the JSON body
	request := &service.RegisterDeviceRequest{}
	err := decodeRequest(context.Writer, context.Request, request)
	if err != nil {
		return
	}

	id, err := wb.Identity.RegisterDevice(request)
	if err != nil {
		log.Println("Error registering device:", err)
		formatStandardResponse("RegDevice", err.Error(), context.Writer)
		return
	}
	formatRegisterResponse(id, context.Writer)
}

// EnrollDevice connects an IoT device with the identity service
func (wb IdentityService) EnrollDevice(c *gin.Context) {
	// Decode the assertions from the request
	assertion1, assertion2, err := decodeEnrollRequest(c.Request)
	if err != nil {
		formatStandardResponse("EnrollDevice", err.Error(), c.Writer)
		return
	}
	if assertion1 == nil || assertion2 == nil {
		formatStandardResponse("EnrollDevice", "A model and serial assertion is required", c.Writer)
		return
	}

	req := service.EnrollDeviceRequest{}

	if assertion1.Type().Name == asserts.ModelType.Name && assertion2.Type().Name == asserts.SerialType.Name {
		req.Model = assertion1
		req.Serial = assertion2
	} else if assertion1.Type().Name == asserts.SerialType.Name && assertion2.Type().Name == asserts.ModelType.Name {
		req.Model = assertion2
		req.Serial = assertion1
	}
	if req.Model == nil || req.Serial == nil {
		log.Println("A model and serial assertion must be provided")
	}

	en, err := wb.Identity.EnrollDevice(&req)
	if err != nil {
		log.Println("Error enrolling device:", err)
		formatStandardResponse("EnrollDevice", err.Error(), c.Writer)
		return
	}

	formatEnrollResponse(*en, c.Writer)
}

func decodeEnrollRequest(r *http.Request) (asserts.Assertion, asserts.Assertion, error) {
	// Use snapd assertion module to decode the assertions in the request stream
	dec := asserts.NewDecoder(r.Body)
	assertion1, err := dec.Decode()
	if err == io.EOF {
		return nil, nil, fmt.Errorf("no data supplied")
	}
	if err != nil {
		return nil, nil, err
	}

	// Decode the second assertion
	assertion2, err := dec.Decode()
	if err != nil && err != io.EOF {
		return nil, nil, err
	}

	// Stream must be ended now
	_, err = dec.Decode()
	if err != io.EOF {
		if err == nil {
			return nil, nil, fmt.Errorf("unexpected assertion in the request stream")
		}
		return nil, nil, err
	}

	return assertion1, assertion2, nil
}

func decodeRequest(w http.ResponseWriter, r *http.Request, i interface{}) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Logger.Error(err)
		}
	}(r.Body)

	// Decode the JSON body
	err := json.NewDecoder(r.Body).Decode(i)
	switch {
	// Check we have some data
	case err == io.EOF:
		formatStandardResponse("NoData", "No data supplied.", w)
		log.Println("No data supplied.")
		// Check for parsing errors
	case err != nil:
		formatStandardResponse("BadData", err.Error(), w)
		log.Println(err)
	}
	return err
}
