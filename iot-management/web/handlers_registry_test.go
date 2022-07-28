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
	"fmt"
	"github.com/everactive/dmscore/iot-identity/web"
	"github.com/everactive/dmscore/iot-management/service/manage/mocks"
	"net/http"
	"testing"
)

func TestService_RegDeviceList(t *testing.T) {
	tests := []struct {
		name        string
		orgID       string
		url         string
		username    string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "abc", "/v1/%s/register/devices", "jamesj", 300, http.StatusOK, ""},
		{"invalid-permissions", "abc", "/v1/%s/register/devices", "jamesj", 0, http.StatusBadRequest, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			manageMock.On("RegDeviceList", tt.orgID, tt.username, tt.permissions).Return(web.DevicesResponse{})

			w := sendRequest("GET", fmt.Sprintf(tt.url, tt.orgID), nil, wb, tt.username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.RegDeviceList() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_RegDeviceGet(t *testing.T) {
	tests := []struct {
		name        string
		orgID       string
		deviceID    string
		url         string
		username    string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "abc", "a111", "/v1/%s/register/devices/%s", "jamesj", 300, http.StatusOK, ""},
		{"invalid-org", "abc", "invalid", "/v1/%s/register/devices/%s", "jamesj", 300, http.StatusBadRequest, "RegDeviceAuth"},
		{"invalid-permissions", "abc", "a111", "/v1/%s/register/devices/%s", "jamesj", 0, http.StatusBadRequest, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			if tt.wantErr == "" {
				manageMock.On("RegDeviceGet", tt.orgID, tt.username, tt.permissions, tt.deviceID).Return(web.EnrollResponse{})
			} else {
				manageMock.On("RegDeviceGet", tt.orgID, tt.username, tt.permissions, tt.deviceID).Return(web.EnrollResponse{
					StandardResponse: web.StandardResponse{Code: tt.wantErr},
				})
			}

			w := sendRequest("GET", fmt.Sprintf(tt.url, tt.orgID, tt.deviceID), nil, wb, tt.username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.RegDeviceGet() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_RegisterDevice(t *testing.T) {
	d1 := []byte(`{"orgid":"abc", "brand":"deviceinc", "model":"A1000", "serial":"d1234"}`)
	tests := []struct {
		name        string
		orgID       string
		url         string
		data        []byte
		username    string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "abc", "/v1/%s/register/devices", d1, "jamesj", 300, http.StatusOK, ""},
		{"invalid-org", "bbb", "/v1/%s/register/devices", d1, "jamesj", 300, http.StatusBadRequest, "RegDevice"},
		{"invalid-permissions", "abc", "/v1/%s/register/devices", d1, "jamesj", 0, http.StatusBadRequest, "UserAuth"},
		{"invalid-data", "abc", "/v1/%s/register/devices", []byte(`\u1000`), "jamesj", 300, http.StatusBadRequest, "RegDevice"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			if tt.wantErr == "" {
				manageMock.On("RegisterDevice", tt.orgID, tt.username, tt.permissions, tt.data).Return(web.RegisterResponse{})
			} else {
				manageMock.On("RegisterDevice", tt.orgID, tt.username, tt.permissions, tt.data).Return(web.RegisterResponse{
					StandardResponse: web.StandardResponse{
						Code: tt.wantErr,
					},
				})
			}

			w := sendRequest("POST", fmt.Sprintf(tt.url, tt.orgID), bytes.NewReader(tt.data), wb, tt.username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.RegisterDevice() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_RegDeviceUpdate(t *testing.T) {
	d1 := []byte(`{"status":3}`)
	tests := []struct {
		name        string
		orgID       string
		deviceID    string
		url         string
		data        []byte
		username    string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "abc", "a111", "/v1/%s/register/devices/%s", d1, "jamesj", 300, http.StatusOK, ""},
		{"invalid-device", "abc", "invalid", "/v1/%s/register/devices/%s", d1, "jamesj", 300, http.StatusBadRequest, "RegDeviceUpdate"},
		{"invalid-permissions", "abc", "a111", "/v1/%s/register/devices/%s", d1, "jamesj", 0, http.StatusBadRequest, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			response := web.StandardResponse{Code: tt.wantErr}
			manageMock.On("RegDeviceUpdate", tt.orgID, tt.username, tt.permissions, tt.deviceID, tt.data).Return(response)

			w := sendRequest("PUT", fmt.Sprintf(tt.url, tt.orgID, tt.deviceID), bytes.NewReader(tt.data), wb, tt.username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.RegDeviceUpdate() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_RegDeviceGetDownload(t *testing.T) {
	tests := []struct {
		name        string
		orgID       string
		deviceID    string
		url         string
		username    string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "abc", "a111", "/v1/%s/register/devices/%s/download", "jamesj", 300, http.StatusOK, ""},
		{"invalid-org", "abc", "invalid", "/v1/%s/register/devices/%s/download", "jamesj", 300, http.StatusBadRequest, "RegDeviceAuth"},
		{"invalid-permissions", "abc", "a111", "/v1/%s/register/devices/%s/download", "jamesj", 0, http.StatusBadRequest, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			response := web.EnrollResponse{
				StandardResponse: web.StandardResponse{Code: tt.wantErr},
			}
			manageMock.On("RegDeviceGet", tt.orgID, tt.username, tt.permissions, tt.deviceID).Return(response)

			w := sendRequest("GET", fmt.Sprintf(tt.url, tt.orgID, tt.deviceID), nil, wb, tt.username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}
		})
	}
}
