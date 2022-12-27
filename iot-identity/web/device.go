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
	"fmt"
	"io"
	"net/http"

	"github.com/everactive/dmscore/iot-identity/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/snapcore/snapd/asserts"
)

// EnrollDevice connects an IoT device with the identity service
func (i IdentityService) EnrollDevice(c *gin.Context) {
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

	log.Tracef("Model and serial asertions decoded")

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

	log.Tracef("Attempting to enroll device")

	en, err := i.Identity.EnrollDevice(&req)
	if err != nil {
		log.Error("enrolling device: ", err)
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
