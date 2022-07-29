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
	"github.com/everactive/dmscore/iot-management/datastore/mocks"
	"testing"

	"github.com/everactive/dmscore/iot-devicetwin/datastore/memory"
)

func TestService_GroupGetExcludedDevices(t *testing.T) {
	type args struct {
		orgID string
		name  string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"valid", args{"abc", "workshop"}, 2, false},
		{"invalid", args{"abc", "does-not-exist"}, 0, true},
	}
	for _, tt := range tests {
		localtt := tt
		t.Run(localtt.name, func(t *testing.T) {
			srv := NewService(memory.NewStore(), &mocks.DataStore{})
			got, err := srv.GroupGetExcludedDevices(localtt.args.orgID, localtt.args.name)
			if (err != nil) != localtt.wantErr {
				t.Errorf("Service.GroupGetExcludedDevices() error = %v, wantErr %v", err, localtt.wantErr)
				return
			}
			if len(got) != localtt.want {
				t.Errorf("Service.GroupGetExcludedDevices() = %v, want %v", len(got), localtt.want)
			}
		})
	}
}
