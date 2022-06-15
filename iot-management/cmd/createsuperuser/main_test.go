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

package main

import (
	"os"
	"testing"

	"github.com/everactive/dmscore/iot-management/datastore/memory"
)

func Test_run_success(t *testing.T) {
	tests := []struct {
		name    string
		user    string
		wantErr bool
	}{
		{"valid", "john", false},
		{"valid-no-user", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := memory.NewStore()
			err := run(db, tt.user, "User", "j@example.com")
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSuperUser.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_main(t *testing.T) {
	type args struct {
		user string
	}
	tests := []struct {
		name string
		args args
	}{
		{"valid", args{"john"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// config.LoadConfig("../testing/memory.yaml")
			os.Setenv("CONFIG_FILE_PATH", "../../testing/memory.yaml")
			parseFlags = func() {}
			username = tt.args.user
			main()
		})
	}
}
