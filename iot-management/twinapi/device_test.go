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

package twinapi

import (
	"errors"
	"fmt"
	"github.com/everactive/dmscore/config/keys"
	"path"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
)

func TestClientAdapter_DeviceList(t *testing.T) {
	b1 := `{"devices": [{"deviceId":"a111"}]}`
	type fields struct {
		URL string
	}
	type args struct {
		orgID string
		body  string
	}
	type test struct {
		name          string
		fields        fields
		args          args
		want          int
		wantErr       string
		responderFunc func(t *test) httpmock.Responder
	}
	validResponder := func(t *test) httpmock.Responder {
		return httpmock.NewStringResponder(200, t.args.body)
	}
	failedResponder := func(t *test) httpmock.Responder {
		return httpmock.NewErrorResponder(errors.New(t.wantErr))
	}
	tests := []test{
		{"valid", fields{""}, args{"abc", b1}, 1, "", validResponder},
		{"invalid-org", fields{""}, args{"invalid", b1}, 0, "MOCK error get", failedResponder},
		{"invalid-body", fields{""}, args{"abc", ""}, 0, "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(keys.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())

			url := path.Join("/", "device", tt.args.orgID)

			httpmock.RegisterResponder("GET", url, tt.responderFunc(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}
			got := a.DeviceList(tt.args.orgID)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Get", url, tt.wantErr)
			if got.Message != tt.wantErr && got.Message != wantErrActual {
				t.Errorf("ClientAdapter.DeviceList() = %v, want %v", got.Message, tt.wantErr)
			}
			if len(got.Devices) != tt.want {
				t.Errorf("ClientAdapter.DeviceList() = %v, want %v", len(got.Devices), tt.want)
			}
		})
	}
}

func TestClientAdapter_DeviceGet(t *testing.T) {
	b1 := `{"device": {"deviceId":"a111"}}`
	type fields struct {
		URL string
	}
	type args struct {
		orgID    string
		deviceID string
		body     string
	}
	type test struct {
		name      string
		fields    fields
		args      args
		want      string
		wantErr   string
		responder func(t *test) httpmock.Responder
	}
	validResponder := func(t *test) httpmock.Responder {
		return httpmock.NewStringResponder(200, t.args.body)
	}
	failedResponder := func(t *test) httpmock.Responder {
		return httpmock.NewErrorResponder(errors.New(t.wantErr))
	}
	tests := []test{
		{"valid", fields{""}, args{"abc", "a111", b1}, "a111", "", validResponder},
		{"invalid-org", fields{""}, args{"invalid", "a111", b1}, "", "MOCK error get", failedResponder},
		{"invalid-body", fields{""}, args{"abc", "a111", ""}, "", "unexpected end of JSON input", validResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(keys.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())

			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID)

			httpmock.RegisterResponder("GET", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}

			got := a.DeviceGet(tt.args.orgID, tt.args.deviceID)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Get", url, tt.wantErr)
			if got.Message != tt.wantErr && got.Message != wantErrActual {
				t.Errorf("ClientAdapter.DeviceGet() = %v, want %v", got.Message, wantErrActual)
			}
			if got.Device.DeviceId != tt.want {
				t.Errorf("ClientAdapter.DeviceGet() = %v, want %v", got.Device.DeviceId, tt.want)
			}
		})
	}
}

func TestClientAdapter_DeviceLogs(t *testing.T) {
	b1 := `{"url": "https://upload.com/upload", "limit": 10}`
	type fields struct {
		URL string
	}
	type args struct {
		orgID    string
		deviceID string
		body     string
		data     []byte
	}
	type test struct {
		name          string
		fields        fields
		args          args
		wantErr       string
		responderFunc func(t *test) httpmock.Responder
	}
	validResponder := func(t *test) httpmock.Responder {
		return httpmock.NewStringResponder(200, t.args.body)
	}
	failedResponder := func(t *test) httpmock.Responder {
		return httpmock.NewErrorResponder(errors.New(t.wantErr))
	}
	tests := []test{
		{"valid", fields{""}, args{"abc", "a111", b1, []byte("{}")}, "", validResponder},
		{"invalid-org", fields{""}, args{"invalid", "a111", b1, []byte("{}")}, "MOCK error post", failedResponder},
		{"invalid-body", fields{""}, args{"abc", "a111", "", []byte("{}")}, "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(keys.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()

			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID, "logs")
			httpmock.RegisterResponder("POST", url, tt.responderFunc(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}

			got := a.DeviceLogs(tt.args.orgID, tt.args.deviceID, []byte(tt.args.body))
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Post", url, tt.wantErr)
			if got.Message != tt.wantErr && got.Message != wantErrActual {
				t.Errorf("ClientAdapter.DeviceLogs() = %v, want %v", got.Message, tt.wantErr)
			}
		})
	}
}

//func TestClientAdapter_DeviceUsersAction(t *testing.T) {
//	type fields struct {
//		URL string
//	}
//	type args struct {
//		orgID    string
//		deviceID string
//		user     DeviceUser
//	}
//	type test struct {
//		name      string
//		fields    fields
//		args      args
//		wantErr   string
//		responder func(t *test) httpmock.Responder
//	}
//	validResponder := func(t *test) httpmock.Responder {
//		return httpmock.NewStringResponder(200, "{}")
//	}
//	failedResponder := func(t *test) httpmock.Responder {
//		return httpmock.NewErrorResponder(errors.New(t.wantErr))
//	}
//	tests := []test{
//		{"valid", fields{""}, args{"abc", "a111", DeviceUser{}}, "", validResponder},
//		{"invalid-org", fields{""}, args{"invalid", "a111", DeviceUser{}}, "MOCK error post", failedResponder},
//		{"invalid-body", fields{""}, args{"abc", "a111", DeviceUser{}}, "unexpected end of JSON input", validResponder},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			viper.Set(keys.ClientTokenProvider, "disabled")
//			client := resty.New()
//			httpmock.ActivateNonDefault(client.GetClient())
//
//			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID, "users")
//			httpmock.RegisterResponder("POST", url, tt.responder(&tt))
//
//			a := &ClientAdapter{
//				URL:    tt.fields.URL,
//				client: client,
//			}
//
//			got := a.DeviceUsersAction(tt.args.orgID, tt.args.deviceID, tt.args.user)
//			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Post", url, tt.wantErr)
//			if got.Message != tt.wantErr && got.Message != wantErrActual {
//				t.Errorf("ClientAdapter.DeviceUsersAction() = %v, want %v", got.Message, wantErrActual)
//			}
//		})
//	}
//}
