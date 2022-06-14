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

package identityapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"testing"

	"github.com/everactive/dmscore/iot-identity/service"
	"github.com/everactive/dmscore/iot-management/config/configkey"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
)

func TestClientAdapter_RegisterOrganization(t *testing.T) {
	b1 := `{"id": "def", "message": ""}`
	type fields struct {
		URL string
	}
	type args struct {
		name    string
		country string
		body    string
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
		{"valid", fields{""}, args{"Test Inc", "GB", b1}, "def", "", validResponder},
		{"invalid-org", fields{"invalid"}, args{"invalid", "GB", b1}, "", "no responder found", failedResponder},
		{"invalid-body", fields{""}, args{"abc", "GB", ""}, "", "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := json.Marshal(service.RegisterOrganizationRequest{Name: tt.args.name, CountryName: tt.args.country})

			viper.Set(configkey.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()

			var url string
			url = path.Join("/", "organization")
			if tt.fields.URL != "" {
				url = path.Join("/", tt.fields.URL, "organization")
			}

			httpmock.RegisterResponder("POST", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}
			got := a.RegisterOrganization(b)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Post", url, tt.wantErr)
			if got.Message != wantErrActual && got.Message != tt.wantErr {
				t.Errorf("ClientAdapter.RegisterOrganization() = %v, want %v", got.Message, tt.wantErr)
			}
			if got.ID != tt.want {
				t.Errorf("ClientAdapter.RegisterOrganization() = %v, want ID %v", got.Code, tt.want)
			}
		})
	}
}

func TestClientAdapter_RegOrganizationList(t *testing.T) {
	b1 := `{"organizations": [{"id":"abc", "name":"Test Org Ltd"}]}`
	type fields struct {
		URL string
	}
	type test struct {
		name      string
		body      string
		fields    fields
		want      int
		wantErr   string
		responder func(t *test) httpmock.Responder
	}
	validResponder := func(t *test) httpmock.Responder {
		return httpmock.NewStringResponder(200, t.body)
	}
	failedResponder := func(t *test) httpmock.Responder {
		return httpmock.NewErrorResponder(errors.New(t.wantErr))
	}
	tests := []test{
		{"valid", b1, fields{""}, 1, "", validResponder},
		{"invalid-org", "", fields{"invalid"}, 0, "MOCK error get", failedResponder},
		{"invalid-body", "", fields{""}, 0, "EOF", failedResponder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(configkey.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()

			var url string
			url = path.Join("/", "organizations")
			if tt.fields.URL != "" {
				url = path.Join("/", tt.fields.URL, "organizations")
			}

			httpmock.RegisterResponder("GET", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}

			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Get", url, tt.wantErr)

			got := a.RegOrganizationList()
			if got.Message != wantErrActual && got.Message != tt.wantErr {
				t.Errorf("ClientAdapter.RegisterOrganization() = %v, want %v", got.Message, tt.wantErr)
			}
			if len(got.Organizations) != tt.want {
				t.Errorf("ClientAdapter.RegisterOrganization() = %v, want ID %v", len(got.Organizations), tt.want)
			}
		})
	}
}
