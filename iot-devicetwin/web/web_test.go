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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func sendRequest(method, url string, data io.Reader, srv *Service) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, url, data)

	engine := gin.Default()
	srv.AddRoutes(engine)

	engine.ServeHTTP(w, r)

	return w
}

func parseSnapsResponse(r io.Reader) (SnapsResponse, error) {
	// Parse the response
	result := SnapsResponse{}
	err := json.NewDecoder(r).Decode(&result)
	return result, err
}

func parseDevicesResponse(r io.Reader) (DevicesResponse, error) {
	// Parse the response
	result := DevicesResponse{}
	err := json.NewDecoder(r).Decode(&result)
	return result, err
}

func parseStandardResponse(r io.Reader) (StandardResponse, error) {
	// Parse the response
	result := StandardResponse{}
	err := json.NewDecoder(r).Decode(&result)
	return result, err
}
