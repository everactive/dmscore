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
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/everactive/dmscore/iot-identity/service"
)

// RegisterOrganization registers a new organization with the identity service
func (i IdentityService) RegisterOrganization(context *gin.Context) {
	// Decode the JSON body
	req, err := decodeOrganizationRequest(context.Request, context.Writer)
	if err != nil {
		return
	}

	id, err := i.Identity.RegisterOrganization(req)
	if err != nil {
		log.Println("Error registering organization:", err)
		formatStandardResponse("RegOrg", err.Error(), context.Writer)
		return
	}
	formatRegisterResponse(id, context.Writer)
}

// OrganizationList fetches organizations
func (i IdentityService) OrganizationList(context *gin.Context) {
	orgs, err := i.Identity.OrganizationList()
	if err != nil {
		log.Println("Error listing organizations:", err)
		formatStandardResponse("OrgList", err.Error(), context.Writer)
		return
	}
	formatOrganizationsResponse(orgs, context.Writer)
}

func decodeOrganizationRequest(r *http.Request, w http.ResponseWriter) (*service.RegisterOrganizationRequest, error) { // Decode the REST request
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Logger.Error(err)
		}
	}(r.Body)

	// Decode the JSON body
	org := service.RegisterOrganizationRequest{}
	err := json.NewDecoder(r.Body).Decode(&org)
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
	return &org, err
}
