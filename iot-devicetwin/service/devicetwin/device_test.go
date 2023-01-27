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

package devicetwin

import (
	"github.com/everactive/dmscore/iot-management/datastore"
	"testing"

	"github.com/everactive/dmscore/iot-devicetwin/datastore/memory"
)

func TestService_DeviceGet(t *testing.T) {
	type args struct {
		orgID    string
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "a111"}, false},
		{"valid-with-server", args{"abc", "c333"}, false},
		{"invalid", args{"abc", "invalid"}, true},
		{"invalid-orgid", args{"invalid", "a111"}, true},
	}
	for _, tt := range tests {
		localtt := tt
		t.Run(localtt.name, func(t *testing.T) {
			srv := NewService(memory.NewStore(), &datastore.MockDataStore{})
			got, err := srv.DeviceGet(localtt.args.orgID, localtt.args.clientID)
			if (err != nil) != localtt.wantErr {
				t.Errorf("Service.DeviceGet() error = %v, wantErr %v", err, localtt.wantErr)
				return
			}
			if localtt.wantErr {
				return
			}
			if got.DeviceId != localtt.args.clientID {
				t.Errorf("Service.DeviceGet() = %v, want %v", got.DeviceId, localtt.args.clientID)
			}
		})
	}
}

func TestService_DeviceList(t *testing.T) {
	type args struct {
		orgID string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"valid", args{"abc"}, 3, false},
		{"valid-no-devices", args{"none"}, 0, false},
		{"invalid", args{"invalid"}, 0, true},
	}
	for _, tt := range tests {
		localtt := tt
		t.Run(localtt.name, func(t *testing.T) {
			srv := NewService(memory.NewStore(), &datastore.MockDataStore{})
			got, err := srv.DeviceList(localtt.args.orgID)
			if (err != nil) != localtt.wantErr {
				t.Errorf("Service.DeviceList() error = %v, wantErr %v", err, localtt.wantErr)
				return
			}
			if len(got) != localtt.want {
				t.Errorf("Service.DeviceList() = %v, want %v", len(got), localtt.want)
			}
		})
	}
}
