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
	"io"
	"net/http"

	"github.com/everactive/dmscore/iot-devicetwin/web"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/domain"
	"github.com/gin-gonic/gin"
)

// OrganizationsResponse defines the response to list users
type OrganizationsResponse struct {
	web.StandardResponse
	Organizations []domain.Organization `json:"organizations"`
}

// OrganizationResponse defines the response to list users
type OrganizationResponse struct {
	web.StandardResponse
	Organization domain.Organization `json:"organization"`
}

// UserOrganization defines an organization and whether it is selected for a user
type UserOrganization struct {
	domain.Organization
	Selected bool `json:"selected"`
}

// UserOrganizationsResponse defines the response to list users
type UserOrganizationsResponse struct {
	web.StandardResponse
	Organizations []UserOrganization `json:"organizations"`
}

func formatOrganizationsResponse(orgs []domain.Organization, w http.ResponseWriter) {
	response := OrganizationsResponse{Organizations: orgs}
	_ = encodeResponse(response, w)
}

func formatOrganizationResponse(org domain.Organization, w http.ResponseWriter) {
	response := OrganizationResponse{Organization: org}
	_ = encodeResponse(response, w)
}

// OrganizationListHandler returns the list of accounts
func (wb Service) OrganizationListHandler(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", JSONHeader)

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	orgs, err := wb.Manage.OrganizationsForUser(user.Username)
	if err != nil {
		formatStandardResponse("OrgList", err.Error(), c)
		return
	}
	formatOrganizationsResponse(orgs, w)
}

// OrganizationGetHandler fetches an organization
func (wb Service) OrganizationGetHandler(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", JSONHeader)
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Superuser)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	org, err := wb.Manage.OrganizationGet(c.Param("id"))
	if err != nil {
		formatStandardResponse("OrgGet", err.Error(), c)
		return
	}

	formatOrganizationResponse(org, w)
}

// OrganizationCreateHandler creates a new organization
func (wb Service) OrganizationCreateHandler(c *gin.Context) {
	r := c.Request
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Superuser)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	org := domain.OrganizationCreate{}
	err = json.NewDecoder(r.Body).Decode(&org)
	switch {
	// Check we have some data
	case err == io.EOF:
		formatStandardResponse("OrgCreate", "No organization data supplied", c)
		return
		// Check for parsing errors
	case err != nil:
		formatStandardResponse("OrgCreate", err.Error(), c)
		return
	}

	if err = wb.Manage.OrganizationCreate(org); err != nil {
		formatStandardResponse("OrgCreate", err.Error(), c)
		return
	}

	orgLookup, err := wb.Manage.OrganizationGet(org.Name)
	if err != nil {
		formatStandardResponse("OrgUpdate", err.Error(), c)
		return
	}
	formatOrganizationResponse(orgLookup, c.Writer)
}

// OrganizationUpdateHandler updates an organization
func (wb Service) OrganizationUpdateHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	w.Header().Set("Content-Type", JSONHeader)
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Superuser)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	org := domain.Organization{}
	err = json.NewDecoder(r.Body).Decode(&org)
	switch {
	// Check we have some data
	case err == io.EOF:
		formatStandardResponse("OrgUpdate", "No organization data supplied", c)
		return
		// Check for parsing errors
	case err != nil:
		formatStandardResponse("OrgUpdate", err.Error(), c)
		return
	}

	if err = wb.Manage.OrganizationUpdate(org); err != nil {
		formatStandardResponse("OrgUpdate", err.Error(), c)
		return
	}
	formatStandardResponse("", "", c)
}

// OrganizationsForUserHandler fetches the organizations for a user
func (wb Service) OrganizationsForUserHandler(c *gin.Context) {
	w := c.Writer
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Superuser)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}

	// Get the organization a user can access
	userOrgs, err := wb.Manage.OrganizationsForUser(c.Param("username"))
	if err != nil {
		formatStandardResponse("OrgList", err.Error(), c)
		return
	}

	// Get all the organizations
	allOrgs, err := wb.Manage.OrganizationsForUser(user.Username)
	if err != nil {
		formatStandardResponse("OrgList", err.Error(), c)
		return
	}

	oo := []UserOrganization{}
	for _, o := range allOrgs {
		found := false
		for _, u := range userOrgs {
			if o.OrganizationID == u.OrganizationID {
				found = true
				break
			}
		}
		oo = append(oo, UserOrganization{o, found})
	}

	_ = encodeResponse(UserOrganizationsResponse{Organizations: oo}, w)
}

// OrganizationUpdateForUserHandler fetches the organizations for a user
func (wb Service) OrganizationUpdateForUserHandler(c *gin.Context) {
	user, err := getUserFromContextAndCheckPermissions(c, datastore.Superuser)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	if err := wb.Manage.OrganizationForUserToggle(c.Param("orgid"), c.Param("username")); err != nil {
		formatStandardResponse("UserOrg", "", c)
		return
	}
	formatStandardResponse("", "", c)
}
