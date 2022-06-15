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
	"net/http"
	"testing"

	"github.com/everactive/dmscore/iot-management/config/configkey"
	"github.com/everactive/dmscore/iot-management/crypt"
	"github.com/spf13/viper"

	"github.com/everactive/dmscore/iot-management/datastore/memory"
	"github.com/everactive/dmscore/iot-management/service/manage"
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
		{"invalid-permissions", "/v1/abc/devices", 0, http.StatusBadRequest, "UserAuth"},

		{"valid", "/v1/abc/devices/a111", 300, http.StatusOK, ""},
		{"invalid-permissions", "/v1/abc/devices/a111", 0, http.StatusBadRequest, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret, err := crypt.CreateSecret(32)
			if err != nil {
				t.Fatalf("Error generating JWT secret: %s", err)
				return
			}
			viper.Set(configkey.JwtSecret, secret)
			db := memory.NewStore()
			wb := NewService(manage.NewMockManagement(db))
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
		{"invalid-permissions", "/v1/abc/devices/a111/actions", 0, http.StatusBadRequest, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret, err := crypt.CreateSecret(32)
			if err != nil {
				t.Fatalf("Error generating JWT secret: %s", err)
				return
			}
			viper.Set(configkey.JwtSecret, secret)
			db := memory.NewStore()
			wb := NewService(manage.NewMockManagement(db))
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
