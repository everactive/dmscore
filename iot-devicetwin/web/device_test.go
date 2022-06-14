// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Device Twin Service
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
	"io"
	"strings"
	"testing"

	"github.com/everactive/dmscore/iot-devicetwin/config/keys"
	"github.com/spf13/viper"
)

func TestService_DeviceGet(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		code   int
		result string
	}{
		{"valid", "/v1/device/abc/a111", 200, ""},
		{"invalid", "/v1/device/abc/invalid", 400, "DeviceGet"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			viper.Set(keys.ConfigPath, tmpDir)
			wb := newService()

			w := sendRequest("GET", tt.url, nil, wb)
			if w.Code != tt.code {
				t.Errorf("Web.DeviceGet() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.DeviceGet() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.DeviceGet() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_DeviceList(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		code   int
		result string
	}{
		{"valid", "/v1/device/abc", 200, ""},
		{"invalid", "/v1/device/invalid", 400, "DeviceList"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			viper.Set(keys.ConfigPath, tmpDir)
			wb := newService()

			w := sendRequest("GET", tt.url, nil, wb)
			if w.Code != tt.code {
				t.Errorf("Web.DeviceList() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseDevicesResponse(w.Body)
			if err != nil {
				t.Errorf("Web.DeviceList() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.DeviceList() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_DeviceUnregister(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		method string
		code   int
		result string
	}{
		{"valid", "/v1/device/abc/a111", "DELETE", 200, ""},
		{"invalid", "/v1/device/abc/invalid", "DELETE", 400, "DeviceDelete"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			viper.Set(keys.ConfigPath, tmpDir)
			wb := newService()

			w := sendRequest(tt.method, tt.url, nil, wb)
			if w.Code != tt.code {
				t.Errorf("Web.DeviceDelete() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseDevicesResponse(w.Body)
			if err != nil {
				t.Errorf("Web.DeviceDelete() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.DeviceDelete() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_DeviceLogs(t *testing.T) {
	logData := `{"url":"https://somelongurl.from.s3.test.com/bucket/something", "limit":200}`
	emptyBody := "{}"
	noBody := ""
	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid", "/v1/device/abc/a111/logs", "POST", strings.NewReader(logData), 202, ""},
		{"empty", "/v1/device/abc/a111/logs", "POST", strings.NewReader(emptyBody), 400, "DeviceLogs"},
		{"noBody", "/v1/device/abc/a111/logs", "POST", strings.NewReader(noBody), 400, "DeviceLogs"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			viper.Set(keys.ConfigPath, tmpDir)
			wb := newService()
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.SnapSnapshot() code got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.SnapSnapshot() got response error = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.SnapSnapshot() got parsed code = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}
