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

type actionFields struct {
	URL string
}
type actionArgs struct {
	orgID    string
	deviceID string
	body     string
}

type actionTest struct {
	name      string
	fields    actionFields
	args      actionArgs
	want      int
	wantErr   string
	responder ActionResponderFunc
}

var actionBody1 = `{"actions": [{"deviceId":"a111"}]}`

type ActionResponderFunc func(t *actionTest) httpmock.Responder

var actionValidResponder = ActionResponderFunc(func(t *actionTest) httpmock.Responder {
	return httpmock.NewStringResponder(200, t.args.body)
})

var actionFailedResponder = ActionResponderFunc(func(t *actionTest) httpmock.Responder {
	return httpmock.NewErrorResponder(errors.New(t.wantErr))
})

var tests = []actionTest{
	{"valid", actionFields{""}, actionArgs{"abc", "a111", actionBody1}, 1, "", actionValidResponder},
	{"invalid-org", actionFields{""}, actionArgs{"invalid", "a111", actionBody1}, 0, "MOCK error get", actionFailedResponder},
	{"invalid-body", actionFields{""}, actionArgs{"abc", "a111", ""}, 0, "EOF", actionFailedResponder},
}

func TestClientAdapter_ActionList(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(keys.ClientTokenProvider, "disabled")
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())

			url := path.Join("/", "device", tt.args.orgID, tt.args.deviceID, "actions")

			httpmock.RegisterResponder("GET", url, tt.responder(&tt))

			a := &ClientAdapter{
				URL:    tt.fields.URL,
				client: client,
			}

			got := a.ActionList(tt.args.orgID, tt.args.deviceID)
			wantErrActual := fmt.Sprintf("%s \"%s\": %s", "Get", url, tt.wantErr)
			if got.Message != tt.wantErr && got.Message != wantErrActual {
				t.Errorf("ClientAdapter.ActionList() = %v, want %v", got.Message, wantErrActual)
			}
			if len(got.Actions) != tt.want {
				t.Errorf("ClientAdapter.ActionList() = %v, want %v", len(got.Actions), tt.want)
			}
		})
	}
}
