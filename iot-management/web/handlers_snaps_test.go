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
	"github.com/everactive/dmscore/iot-devicetwin/web"
	"github.com/everactive/dmscore/iot-management/service/manage/mocks"
	"net/http"
	"path"
	"testing"
)

func TestService_SnapListHandler(t *testing.T) {
	tests := []struct {
		name        string
		orgID       string
		deviceID    string
		username    string
		url         string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "abc", "a111", "everactive", "/v1/device/%s/%s/snaps", 300, http.StatusOK, ""},
		{"invalid-permissions", "abc", "a111", "everactive", "/v1/device/%s/%s/snaps", 0, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			manageMock.On("SnapList", tt.orgID, tt.username, tt.permissions, tt.deviceID).Return(web.SnapsResponse{})

			w := sendRequest("GET", fmt.Sprintf(tt.url, tt.orgID, tt.deviceID), nil, wb, tt.username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.SnapListHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_SnapWorkflow_SnapListOnDevice(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		orgID       string
		deviceID    string
		username    string
		url         string
		body        []byte
		permissions int
		want        int
		wantErr     string
	}{
		{"list-valid", "POST", "abc", "a111", "everactive", "/v1/snaps/%s/%s/list", nil, 300, http.StatusOK, ""},
		{"list-invalid-permissions", "POST", "abc", "a111", "everactive", "/v1/snaps/%s/%s/list", nil, 0, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			manageMock.On("SnapListOnDevice", tt.orgID, tt.username, tt.permissions, tt.deviceID).Return(web.StandardResponse{})

			w := sendRequest(tt.method, fmt.Sprintf(tt.url, tt.orgID, tt.deviceID), bytes.NewReader(tt.body), wb, tt.username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.SnapInstallHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_SnapWorkflow_SnapInstall(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		orgID       string
		deviceID    string
		username    string
		snap        string
		url         string
		body        []byte
		permissions int
		want        int
		wantErr     string
	}{
		{"install-valid", "POST", "abc", "a111", "everactive", "helloworld", "/v1/snaps/%s/%s/%s", nil, 300, http.StatusOK, ""},
		{"install-invalid-permissions", "POST", "abc", "a111", "everactive", "helloworld", "/v1/snaps/%s/%s/%s", nil, 0, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			manageMock.On("SnapInstall", tt.orgID, tt.username, tt.permissions, tt.deviceID, tt.snap).Return(web.StandardResponse{})

			w := sendRequest(tt.method, fmt.Sprintf(tt.url, tt.orgID, tt.deviceID, tt.snap), bytes.NewReader(tt.body), wb, tt.username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.SnapInstallHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_SnapWorkflow_SnapDelete(t *testing.T) {
	username := "everactive"
	method := "DELETE"
	orgID := "abc"
	deviceID := "a111"
	snap := "helloworld"

	tests := []struct {
		name        string
		url         string
		body        []byte
		permissions int
		want        int
		wantErr     string
	}{
		{"delete-valid", "/v1/snaps/%s/%s/%s", nil, 300, http.StatusOK, ""},
		{"delete-invalid-permissions", "/v1/snaps/%s/%s/%s", nil, 0, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			manageMock.On("SnapRemove", orgID, username, tt.permissions, deviceID, snap).Return(web.StandardResponse{})

			w := sendRequest(method, fmt.Sprintf(tt.url, orgID, deviceID, snap), bytes.NewReader(tt.body), wb, username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.SnapInstallHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_SnapWorkflow_SnapUpdate(t *testing.T) {
	username := "everactive"
	method := "PUT"
	orgID := "abc"
	deviceID := "a111"
	snap := "helloworld"

	tests := []struct {
		name        string
		url         string
		body        []byte
		permissions int
		want        int
		wantErr     string
	}{
		{"update-valid-refresh", "/v1/snaps/%s/%s/%s/refresh", nil, 300, http.StatusOK, ""},
		{"update-valid-enable", "/v1/snaps/%s/%s/%s/enable", nil, 300, http.StatusOK, ""},
		{"update-valid-disable", "/v1/snaps/%s/%s/%s/disable", nil, 300, http.StatusOK, ""},
		{"update-valid-switch", "/v1/snaps/%s/%s/%s/switch", []byte("{}"), 300, http.StatusOK, ""},
		{"update-action-invalid", "/v1/snaps/%s/%s/%s/invalid", nil, 300, http.StatusBadRequest, "SnapUpdate"},
		{"update-invalid-permissions", "/v1/snaps/%s/%s/%s/refresh", nil, 0, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			_, action := path.Split(tt.url)
			body := tt.body
			if tt.body == nil {
				body = []byte("{}")
			}
			manageMock.On("SnapUpdate", orgID, username, tt.permissions, deviceID, snap, action, body).Return(web.StandardResponse{Code: tt.wantErr})

			w := sendRequest(method, fmt.Sprintf(tt.url, orgID, deviceID, snap), bytes.NewReader(tt.body), wb, username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.SnapInstallHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_SnapWorkflow_SnapConfigUpdate(t *testing.T) {
	username := "everactive"
	method := "PUT"
	orgID := "abc"
	deviceID := "a111"
	snap := "helloworld"
	url := fmt.Sprintf("/v1/snaps/%s/%s/%s/settings", orgID, deviceID, snap)
	tests := []struct {
		name        string
		body        []byte
		permissions int
		want        int
		wantErr     string
	}{
		{"config-valid", []byte("{}"), 300, http.StatusOK, ""},
		{"config-valid-empty", nil, 300, http.StatusOK, ""},
		{"config-invalid-permissions", []byte("{}"), 0, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			body := tt.body
			if body == nil {
				body = []byte("{}")
			}
			manageMock.On("SnapConfigSet", orgID, username, tt.permissions, deviceID, snap, body).Return(web.StandardResponse{})

			w := sendRequest(method, url, bytes.NewReader(tt.body), wb, username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.SnapInstallHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_SnapWorkflow_SnapSnapshot(t *testing.T) {
	username := "everactive"
	method := "POST"
	orgID := "abc"
	deviceID := "a111"
	snap := "helloworld"
	url := fmt.Sprintf("/v1/snaps/%s/%s/%s/snapshot", orgID, deviceID, snap)

	tests := []struct {
		name        string
		body        []byte
		permissions int
		want        int
		wantErr     string
	}{
		{"send-snapshot-valid", []byte("{}"), 300, http.StatusOK, ""},
		{"send-snapshot-invalid-permissions", []byte("{}"), 0, http.StatusUnauthorized, "UserAuth"},
		{"send-snapshot-valid-empty", nil, 300, http.StatusBadRequest, "SnapSnapshot"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			manageMock.On("SnapSnapshot", orgID, username, tt.permissions, deviceID, snap, tt.body).Return(web.StandardResponse{})

			w := sendRequest(method, url, bytes.NewReader(tt.body), wb, username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.SnapInstallHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}
