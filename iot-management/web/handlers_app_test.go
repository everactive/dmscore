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

package web

import (
	"net/http"
	"testing"

	"github.com/everactive/dmscore/iot-management/crypt"
	"github.com/everactive/dmscore/iot-management/datastore/memory"
	"github.com/everactive/dmscore/iot-management/service/manage"
)

func TestService_IndexHandler(t *testing.T) {
	tests := []struct {
		name     string
		template string
		want     int
	}{
		{"valid", "/../static/app.html", http.StatusOK},
		{"invalid-template", "/does-not-exist.html", http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexTemplate = tt.template

			secret, err := crypt.CreateSecret(32)
			if err != nil {
				t.Fatalf("Error generating JWT secret: %s", err)
				return
			}

			db := memory.NewStore()
			wb := NewService(manage.NewMockManagement(db))
			w := sendRequest("GET", "/", nil, wb, "jamesj", secret, 100)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}
		})
	}
}

func TestService_VersionTokenHandler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"valid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret, err := crypt.CreateSecret(32)
			if err != nil {
				t.Fatalf("Error generating JWT secret: %s", err)
				return
			}

			db := memory.NewStore()
			wb := NewService(manage.NewMockManagement(db))
			w := sendRequest("GET", "/v1/versions", nil, wb, "jamesj", secret, 100)
			if w.Code != http.StatusOK {
				t.Errorf("Expected HTTP status '%d', got: %v", http.StatusOK, w.Code)
			}

			w = sendRequest("GET", "/v1/token", nil, wb, "jamesj", secret, 100)
			if w.Code != http.StatusOK {
				t.Errorf("Expected HTTP status '%d', got: %v", http.StatusOK, w.Code)
			}
		})
	}
}
