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
	"path"
	"testing"

	"github.com/everactive/dmscore/iot-management/config/configkey"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
)

const (
	emptyResponseBody = `{"code":"", "message":""}`
)

func TestClientAdapter_SnapList(t *testing.T) {
	b1 := `{"snaps": [{"deviceId":"a111", "name":"helloworld"}]}`
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
		want      int
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
		{"valid", fields{""}, args{"abc", "a111", b1}, 1, "", validResponder},
		{"invalid-org", fields{""}, args{"invalid", "a111", b1}, 0, "MOCK error get", failedResponder},
		{"invalid-body", fields{""}, args{"abc", "a111", ""}, 0, "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(configkey.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()

			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID, "snaps")

			httpmock.RegisterResponder("GET", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}

			got := a.SnapList(tt.args.orgID, tt.args.deviceID)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Get", url, tt.wantErr)
			if got.Message != wantErrActual && got.Message != tt.wantErr {
				t.Errorf("ClientAdapter.SnapList() = %v, want %v", got.Message, tt.wantErr)
			}
			if len(got.Snaps) != tt.want {
				t.Errorf("ClientAdapter.SnapList() = %v, want %v", len(got.Snaps), tt.want)
			}
		})
	}
}

func TestClientAdapter_SnapInstall(t *testing.T) {
	type fields struct {
		URL string
	}
	type args struct {
		orgID    string
		deviceID string
		snap     string
		body     string
	}
	type test struct {
		name      string
		fields    fields
		args      args
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
		{"valid", fields{""}, args{"abc", "a111", "helloworld", emptyResponseBody}, "", validResponder},
		{"invalid-org", fields{""}, args{"invalid", "a111", "helloworld", emptyResponseBody}, "MOCK error post", failedResponder},
		{"invalid-body", fields{""}, args{"abc", "a111", "helloworld", ""}, "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(configkey.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()

			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID, "snaps", tt.args.snap)

			httpmock.RegisterResponder("POST", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}
			got := a.SnapInstall(tt.args.orgID, tt.args.deviceID, tt.args.snap)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Post", url, tt.wantErr)
			if got.Message != wantErrActual && got.Message != tt.wantErr {
				t.Errorf("ClientAdapter.SnapInstall() = %v, want %v", got.Message, tt.wantErr)
			}
		})
	}
}

func TestClientAdapter_SnapRemove(t *testing.T) {
	type fields struct {
		URL string
	}
	type args struct {
		orgID    string
		deviceID string
		snap     string
		body     string
	}
	type test struct {
		name      string
		fields    fields
		args      args
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
		{"valid", fields{""}, args{"abc", "a111", "helloworld", emptyResponseBody}, "", validResponder},
		{"invalid-org", fields{""}, args{"invalid", "a111", "helloworld", emptyResponseBody}, "MOCK error delete", failedResponder},
		{"invalid-body", fields{""}, args{"abc", "a111", "helloworld", ""}, "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(configkey.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()

			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID, "snaps", tt.args.snap)
			httpmock.RegisterResponder("DELETE", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}
			got := a.SnapRemove(tt.args.orgID, tt.args.deviceID, tt.args.snap)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Delete", url, tt.wantErr)
			if got.Message != wantErrActual && got.Message != tt.wantErr {
				t.Errorf("ClientAdapter.SnapRemove() = %v, want %v", got.Message, tt.wantErr)
			}
		})
	}
}

func TestClientAdapter_SnapUpdate(t *testing.T) {
	type fields struct {
		URL string
	}
	type args struct {
		orgID    string
		deviceID string
		snap     string
		action   string
		body     string
		data     []byte
	}
	type test struct {
		name      string
		fields    fields
		args      args
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
		{"valid", fields{""}, args{"abc", "a111", "helloworld", "refresh", emptyResponseBody, []byte("{}")}, "", validResponder},
		{"invalid-org", fields{""}, args{"invalid", "a111", "helloworld", "refresh,", emptyResponseBody, []byte("{}")}, "MOCK error put", failedResponder},
		{"invalid-body", fields{""}, args{"abc", "a111", "helloworld", "refresh", "", []byte("{}")}, "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(configkey.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()

			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID, "snaps", tt.args.snap, tt.args.action)
			httpmock.RegisterResponder("PUT", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}
			got := a.SnapUpdate(tt.args.orgID, tt.args.deviceID, tt.args.snap, tt.args.action, tt.args.data)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Put", url, tt.wantErr)
			if got.Message != wantErrActual && got.Message != tt.wantErr {
				t.Errorf("ClientAdapter.SnapUpdate() = %v, want %v", got.Message, tt.wantErr)
			}
		})
	}
}

func TestClientAdapter_SnapConfigSet(t *testing.T) {
	config := []byte(`{"title":"Hello World!"}`)
	b1 := `{"code":"", "message":""}`
	type fields struct {
		URL string
	}
	type args struct {
		orgID    string
		deviceID string
		snap     string
		config   []byte
		body     string
	}
	type test struct {
		name      string
		fields    fields
		args      args
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
		{"valid", fields{""}, args{"abc", "a111", "helloworld", config, b1}, "", validResponder},
		{"invalid-org", fields{""}, args{"invalid", "a111", "helloworld", config, b1}, "MOCK error put", failedResponder},
		{"invalid-body", fields{""}, args{"abc", "a111", "helloworld", config, ""}, "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(configkey.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()

			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID, "snaps", tt.args.snap, "settings")
			httpmock.RegisterResponder("PUT", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}
			got := a.SnapConfigSet(tt.args.orgID, tt.args.deviceID, tt.args.snap, tt.args.config)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Put", url, tt.wantErr)
			if got.Message != wantErrActual && got.Message != tt.wantErr {
				t.Errorf("ClientAdapter.SnapConfigSet() = %v, want %v", got.Message, tt.wantErr)
			}
		})
	}
}

func TestClientAdapter_SnapListOnDevice(t *testing.T) {
	type fields struct {
		URL string
	}
	type args struct {
		orgID    string
		deviceID string
		snap     string
		action   string
		body     string
	}
	type test struct {
		name      string
		fields    fields
		args      args
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
		{"valid", fields{""}, args{"abc", "a111", "helloworld", "refresh", emptyResponseBody}, "", validResponder},
		{"invalid-org", fields{""}, args{"invalid", "a111", "helloworld", "refresh", emptyResponseBody}, "MOCK error post", failedResponder},
		{"invalid-body", fields{""}, args{"abc", "a111", "helloworld", "refresh", ""}, "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(configkey.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()

			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID, "snaps", "list")
			httpmock.RegisterResponder("POST", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}

			got := a.SnapListOnDevice(tt.args.orgID, tt.args.deviceID)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Post", url, tt.wantErr)
			if got.Message != wantErrActual && got.Message != tt.wantErr {
				t.Errorf("ClientAdapter.SnapListOnDevice() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestClientAdapter_SnapSnapshot(t *testing.T) {
	b1 := `{"url":"", "https://upload.com/upload":""}`
	type fields struct {
		URL string
	}
	type args struct {
		orgID    string
		deviceID string
		snap     string
		body     string
		data     []byte
	}
	type test struct {
		name      string
		fields    fields
		args      args
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
		{"valid", fields{""}, args{"abc", "a111", "helloworld", b1, []byte("{}")}, "", validResponder},
		{"invalid-org", fields{""}, args{"invalid", "a111", "helloworld", b1, []byte("{}")}, "MOCK error post", failedResponder},
		{"invalid-body", fields{""}, args{"abc", "a111", "helloworld", "", []byte("{}")}, "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(configkey.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()

			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID, "snaps", tt.args.snap, "snapshot")

			httpmock.RegisterResponder("POST", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}
			got := a.SnapSnapshot(tt.args.orgID, tt.args.deviceID, tt.args.snap, tt.args.data)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Post", url, tt.wantErr)
			if got.Message != wantErrActual && got.Message != tt.wantErr {
				t.Errorf("ClientAdapter.SnapSnapshot() = %v, want %v", got.Message, tt.wantErr)
			}
		})
	}
}
