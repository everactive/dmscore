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

	"github.com/everactive/dmscore/iot-devicetwin/config"
	"github.com/everactive/dmscore/iot-devicetwin/config/keys"
	"github.com/spf13/viper"

	"github.com/everactive/dmscore/iot-devicetwin/service/controller"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	"github.com/everactive/dmscore/iot-devicetwin/service/mqtt"
)

func testController() controller.Controller {
	return &controller.Service{
		MQTT:       &mqtt.MockConnect{},
		DeviceTwin: &devicetwin.MockDeviceTwin{},
	}
}

func newService() *Service {
	config.LoadConfig(viper.GetString(keys.ConfigPath))
	return NewService(viper.GetString(keys.ServicePort), testController())
}

func TestService_SnapList(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		code   int
		result string
	}{
		{"valid", "/v1/device/abc/a111/snaps", 200, ""},
		{"invalid", "/v1/device/abc/invalid/snaps", 400, "SnapList"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			viper.Set(keys.ConfigPath, tmpDir)
			wb := newService()

			w := sendRequest("GET", tt.url, nil, wb)
			if w.Code != tt.code {
				t.Errorf("Web.SnapList() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseSnapsResponse(w.Body)
			if err != nil {
				t.Errorf("Web.SnapList() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.SnapList() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_SnapActions(t *testing.T) {
	settings1 := `{"title": "Hello World!}"`
	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid-install", "/v1/device/abc/a111/snaps/helloworld", "POST", nil, 200, ""},
		{"invalid-install", "/v1/device/abc/invalid/snaps/helloworld", "POST", nil, 400, "SnapInstall"},

		{"valid-remove", "/v1/device/abc/a111/snaps/helloworld", "DELETE", nil, 200, ""},
		{"invalid-remove", "/v1/device/abc/invalid/snaps/helloworld", "DELETE", nil, 400, "SnapRemove"},

		{"valid-update-enable", "/v1/device/abc/a111/snaps/helloworld/enable", "PUT", nil, 200, ""},
		{"invalid-update-enable", "/v1/device/abc/invalid/snaps/helloworld/enable", "PUT", nil, 400, "SnapUpdate"},
		{"valid-update-disable", "/v1/device/abc/a111/snaps/helloworld/disable", "PUT", nil, 200, ""},
		{"invalid-update-disable", "/v1/device/abc/invalid/snaps/helloworld/disable", "PUT", nil, 400, "SnapUpdate"},
		{"valid-update-refresh", "/v1/device/abc/a111/snaps/helloworld/refresh", "PUT", nil, 200, ""},
		{"invalid-update-refresh", "/v1/device/abc/invalid/snaps/helloworld/refresh", "PUT", nil, 400, "SnapUpdate"},
		{"invalid-update-invalid", "/v1/device/abc/a111/snaps/helloworld/invalid", "PUT", nil, 400, "SnapUpdate"},
		{"valid-update-settings", "/v1/device/abc/a111/snaps/helloworld/settings", "PUT", strings.NewReader(settings1), 200, ""},
		{"invalid-update-settings", "/v1/device/abc/invalid/snaps/helloworld/settings", "PUT", strings.NewReader(settings1), 400, "SnapSetConf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			viper.Set(keys.ConfigPath, tmpDir)
			wb := newService()
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.SnapActions() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.SnapActions() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.SnapActions() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_SnapServiceAction(t *testing.T) {
	snapServices := `{"services":["hello_world", "goodbye_world"]}`
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
		{"valid-start", "/v1/device/abc/a111/services/helloworld/start", "POST", strings.NewReader(emptyBody), 200, ""},
		{"valid-start-specific-services", "/v1/device/abc/a111/services/helloworld/start", "POST", strings.NewReader(snapServices), 200, ""},
		{"valid-stop", "/v1/device/abc/a111/services/helloworld/stop", "POST", strings.NewReader(emptyBody), 200, ""},
		{"valid-stop-specific-services", "/v1/device/abc/a111/services/helloworld/stop", "POST", strings.NewReader(snapServices), 200, ""},
		{"valid-restart", "/v1/device/abc/a111/services/helloworld/restart", "POST", strings.NewReader(emptyBody), 200, ""},
		{"valid-restart-specific-services", "/v1/device/abc/a111/services/helloworld/restart", "POST", strings.NewReader(snapServices), 200, ""},

		{"empty-start", "/v1/device/abc/a111/services/helloworld/start", "POST", strings.NewReader(noBody), 200, ""},
		{"empty-stop", "/v1/device/abc/a111/services/helloworld/stop", "POST", strings.NewReader(noBody), 200, ""},
		{"empty-restart", "/v1/device/abc/a111/services/helloworld/restart", "POST", strings.NewReader(noBody), 200, ""},

		{"invalid-action", "/v1/device/abc/a111/services/helloworld/invalid", "POST", strings.NewReader(emptyBody), 400, "SnapServiceAction"},
		//  /v1/device/{orgid}/{id}/snaps/{snap}/{action}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			viper.Set(keys.ConfigPath, tmpDir)
			wb := newService()
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.SnapActions() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.SnapActions() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.SnapActions() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_SnapSnapshot(t *testing.T) {
	s3Data := `{"url":"https://somelongurl.from.s3.test.com/bucket/something"}`
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
		{"valid", "/v1/device/abc/a111/snaps/helloworld/snapshot", "POST", strings.NewReader(s3Data), 202, ""},
		{"empty", "/v1/device/abc/a111/snaps/helloworld/snapshot", "POST", strings.NewReader(emptyBody), 400, "SnapSnapshot"},
		{"noBody", "/v1/device/abc/a111/snaps/helloworld/snapshot", "POST", strings.NewReader(noBody), 400, "SnapSnapshot"},
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
