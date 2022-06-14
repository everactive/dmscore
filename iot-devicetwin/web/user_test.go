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

func TestService_UserAdd(t *testing.T) {
	validAddUser := `{"email":"testemailstring", "force-managed":true, "action":"create"}`

	// missing email
	invalidAddUser1 := `{"force-managed":true, "action":"create"}`
	// missing action=create
	invalidAddUser2 := `{"email":"testemailstring", "force-managed":true}`

	emptyBody := "{}"
	noBody := ""

	endpoint := "/v1/device/abc/a111/users"
	responseSubject := "UserAdd"

	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid", endpoint, "POST", strings.NewReader(validAddUser), 200, ""},
		{"invalid_missing_email", endpoint, "POST", strings.NewReader(invalidAddUser1), 400, responseSubject},
		{"invalid_missing_action", endpoint, "POST", strings.NewReader(invalidAddUser2), 400, responseSubject},
		{"empty", endpoint, "POST", strings.NewReader(emptyBody), 400, responseSubject},
		{"noBody", endpoint, "POST", strings.NewReader(noBody), 400, responseSubject},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			viper.Set(keys.ConfigPath, tmpDir)
			wb := newService()
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.UserAdd() code got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.UserAdd() got response error = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.UserAdd() got parsed code = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_UserRemove(t *testing.T) {
	validRemoveUser := `{"username":"testusername", "action":"remove"}`

	// missing username
	invalidRemoveUser1 := `{"action":"remove"}`
	// missing action=remove
	invalidRemoveUser2 := `{"username":"testusername"}`

	emptyBody := "{}"
	noBody := ""

	endpoint := "/v1/device/abc/a111/users"
	method := "DELETE"
	responseSubject := "UserRemove"

	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid", endpoint, method, strings.NewReader(validRemoveUser), 200, ""},
		{"invalid_missing_email", endpoint, method, strings.NewReader(invalidRemoveUser1), 400, responseSubject},
		{"invalid_missing_action", endpoint, method, strings.NewReader(invalidRemoveUser2), 400, responseSubject},
		{"empty", endpoint, method, strings.NewReader(emptyBody), 400, responseSubject},
		{"noBody", endpoint, method, strings.NewReader(noBody), 400, responseSubject},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			viper.Set(keys.ConfigPath, tmpDir)
			wb := newService()
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.UserRemove() code got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.UserRemove() got response error = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.UserRemove() got parsed code = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}
