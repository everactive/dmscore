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

// Package twinapi provides the interface and implementation of Client to interact with the DeviceTwin REST API
package twinapi

import (
	"bytes"
	"github.com/everactive/dmscore/iot-devicetwin/web"
	"net/url"
	"path"

	"github.com/everactive/dmscore/iot-management/requests"
	"github.com/go-resty/resty/v2"

	log "github.com/sirupsen/logrus"
)

// Client is a client for the device twin API
type Client interface {
	//DeviceList(orgID string) web.DevicesResponse
	//DeviceGet(orgID, deviceID string) web.DeviceResponse
	//DeviceDelete(orgID, deviceID string) web.StandardResponse
	//DeviceLogs(orgID, deviceID string, body []byte) web.StandardResponse
	//DeviceUsersAction(orgID, deviceID string, body []byte) web.StandardResponse

	ActionList(orgID, deviceID string) web.ActionsResponse
	SnapList(orgID, deviceID string) web.SnapsResponse

	SnapSnapshot(orgID, deviceID, snap string, body []byte) web.StandardResponse
	SnapListOnDevice(orgID, deviceID string) web.StandardResponse
	SnapInstall(orgID, deviceID, snap string) web.StandardResponse
	SnapRemove(orgID, deviceID, snap string) web.StandardResponse
	SnapUpdate(orgID, deviceID, snap, action string, body []byte) web.StandardResponse
	SnapConfigSet(orgID, deviceID, snap string, config []byte) web.StandardResponse
	SnapServiceAction(orgID, deviceID, snap, action string, body []byte) web.StandardResponse

	GroupList(orgID string) web.GroupsResponse
	GroupCreate(orgID string, body []byte) web.StandardResponse
	GroupDevices(orgID, name string) web.DevicesResponse
	GroupExcludedDevices(orgID, name string) web.DevicesResponse
	GroupDeviceLink(orgID, name, deviceID string) web.StandardResponse
	GroupDeviceUnlink(orgID, name, deviceID string) web.StandardResponse
}

// ClientAdapter adapts our expectations to device twin API
type ClientAdapter struct {
	URL    string
	client *resty.Client
}

// NewClientAdapter creates an adapter to access the device twin service
func NewClientAdapter(u string) (*ClientAdapter, error) {
	return &ClientAdapter{
		URL:    u,
		client: resty.New(),
	}, nil
}

func (a *ClientAdapter) urlPath(p string) string {
	u, _ := url.Parse(a.URL)
	u.Path = path.Join(u.Path, p)
	return u.String()
}

// AddAuthorization is a variable way to add authorization to a request with a
// default that checks for an explicit disabled if no other function is provided and otherwise panics
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

func (a *ClientAdapter) deleteWithBody(p string, data []byte) (*resty.Response, error) {
	req := a.client.R()
	addAuthorization(req)
	req.SetBody(bytes.NewReader(data))
	resp, err := req.Delete(p)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resp, nil
}
