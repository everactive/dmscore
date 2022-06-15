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

package manage

import (
	"testing"

	"github.com/everactive/dmscore/iot-management/identityapi"

	"github.com/everactive/dmscore/iot-management/datastore/memory"
	"github.com/everactive/dmscore/iot-management/domain"
	"github.com/everactive/dmscore/iot-management/twinapi"
)

func TestManagement_OrganizationsForUser(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"valid", args{"jamesj"}, 1, false},
		{"valid-not-found", args{"unknown"}, 0, false},
		{"invalid", args{"invalid"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewManagement(memory.NewStore(), &twinapi.MockClient{}, &identityapi.MockIdentity{})
			got, err := srv.OrganizationsForUser(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("Management.OrganizationsForUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Management.OrganizationsForUser() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestManagement_OrganizationGet(t *testing.T) {
	type args struct {
		orgID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc"}, false},
		{"invalid", args{"invalid"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewManagement(memory.NewStore(), &twinapi.MockClient{}, &identityapi.MockIdentity{})
			got, err := srv.OrganizationGet(tt.args.orgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Management.OrganizationGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.OrganizationID != tt.args.orgID {
				t.Errorf("Management.OrganizationGet() = %v, want %v", got.OrganizationID, tt.args.orgID)
			}
		})
	}
}

func TestManagement_OrganizationCreate(t *testing.T) {
	type args struct {
		org domain.OrganizationCreate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{domain.OrganizationCreate{Name: "Test Inc", Country: "GB"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewManagement(memory.NewStore(), &twinapi.MockClient{}, &identityapi.MockIdentity{})
			if err := srv.OrganizationCreate(tt.args.org); (err != nil) != tt.wantErr {
				t.Errorf("Management.OrganizationCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManagement_OrganizationUpdate(t *testing.T) {
	type args struct {
		org domain.Organization
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{domain.Organization{OrganizationID: "abc", Name: "Test Inc"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewManagement(memory.NewStore(), &twinapi.MockClient{}, &identityapi.MockIdentity{})
			if err := srv.OrganizationUpdate(tt.args.org); (err != nil) != tt.wantErr {
				t.Errorf("Management.OrganizationUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManagement_OrganizationForUserToggle(t *testing.T) {
	type args struct {
		orgID    string
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "jamesj"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewManagement(memory.NewStore(), &twinapi.MockClient{}, &identityapi.MockIdentity{})
			if err := srv.OrganizationForUserToggle(tt.args.orgID, tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("Management.OrganizationForUserToggle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
