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
	"github.com/everactive/dmscore/config/keys"
	"github.com/everactive/dmscore/iot-devicetwin/web"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"

	"github.com/everactive/dmscore/iot-management/crypt"
	"github.com/spf13/viper"
)

func TestService_DeviceHandlers(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "/v1/abc/devices", 300, http.StatusOK, ""},
		{"invalid-permissions", "/v1/abc/devices", 0, http.StatusUnauthorized, "UserAuth"},

		{"valid", "/v1/abc/devices/a111", 300, http.StatusOK, ""},
		{"invalid-permissions", "/v1/abc/devices/a111", 0, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret, err := crypt.CreateSecret(32)
			if err != nil {
				t.Fatalf("Error generating JWT secret: %s", err)
				return
			}
			viper.Set(keys.JwtSecret, secret)

			manageMock := &manage.MockManage{}
			manageMock.On("DeviceList", mock.Anything, mock.Anything, mock.Anything).Return(web.DevicesResponse{})
			manageMock.On("DeviceGet", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(web.DeviceResponse{})

			wb := NewService(manageMock, gin.Default())
			w := sendRequest("GET", tt.url, nil, wb, "jamesj", secret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.DeviceHandlers() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_ActionListHandler(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "/v1/abc/devices/a111/actions", 300, http.StatusOK, ""},
		{"invalid-permissions", "/v1/abc/devices/a111/actions", 0, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret, err := crypt.CreateSecret(32)
			if err != nil {
				t.Fatalf("Error generating JWT secret: %s", err)
				return
			}
			viper.Set(keys.JwtSecret, secret)

			manageMock := &manage.MockManage{}
			manageMock.On("ActionList", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(web.ActionsResponse{})

			wb := NewService(manageMock, gin.Default())
			w := sendRequest("GET", tt.url, nil, wb, "jamesj", secret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.ActionListHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}
