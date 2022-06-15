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

// Package identityapi provides an interface and types for interacting with the Identity REST API
package identityapi

import (
	"bytes"
	"github.com/everactive/dmscore/iot-identity/web"
	"net/url"
	"path"

	"github.com/everactive/dmscore/iot-management/requests"

	"github.com/go-resty/resty/v2"

	log "github.com/sirupsen/logrus"
)

// Client is a client for the identity API
type Client interface {
	RegDeviceList(orgID string) web.DevicesResponse
	RegisterDevice(body []byte) web.RegisterResponse
	RegDeviceGet(orgID, deviceID string) web.EnrollResponse
	RegDeviceUpdate(orgID, deviceID string, body []byte) web.StandardResponse
	//DeviceDelete(orgID string) web.StandardResponse

	RegisterOrganization(body []byte) web.RegisterResponse
	RegOrganizationList() web.OrganizationsResponse
}

// ClientAdapter adapts our expectations to device twin API
type ClientAdapter struct {
	URL    string
	client *resty.Client
}

// NewClientAdapter creates an adapter to access the device twin service
func NewClientAdapter(u string) (*ClientAdapter, error) {
	adapter := &ClientAdapter{
		URL:    u,
		client: resty.New(),
	}

	return adapter, nil
}

func (a *ClientAdapter) urlPath(p string) string {
	u, _ := url.Parse(a.URL)
	u.Path = path.Join(u.Path, p)
	return u.String()
}

// AddAuthorization is a function to add an Authorization header to requests before making them for the IdentityAPI
var AddAuthorization = requests.DefaultAddAuthorization

func addAuthorization(req *resty.Request) {
	err := AddAuthorization(req)
	if err != nil {
		log.Error(err)
	}
}

func (a *ClientAdapter) get(p string) (*resty.Response, error) {
	req := a.client.R()
	addAuthorization(req)
	resp, err := req.Get(p)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resp, nil
}

func (a *ClientAdapter) delete(p string) (*resty.Response, error) {
	req := a.client.R()
	addAuthorization(req)
	resp, err := req.Delete(p)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resp, nil
}

func (a *ClientAdapter) post(p string, data []byte) (*resty.Response, error) {
	req := a.client.R()
	addAuthorization(req)
	req.SetHeader("Content-Type", "application/json")
	req.SetBody(bytes.NewReader(data))
	resp, err := req.Post(p)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resp, nil
}

func (a *ClientAdapter) put(p string, data []byte) (*resty.Response, error) {
	req := a.client.R()
	addAuthorization(req)
	req.SetBody(bytes.NewReader(data))
	resp, err := req.Put(p)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resp, nil
}
