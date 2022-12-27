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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/everactive/dmscore/iot-devicetwin/web"
	"github.com/everactive/dmscore/iot-management/web/usso"
	"github.com/gin-gonic/gin"
	"github.com/juju/usso/openid"
)

func sendRequestWithBeforeServeHook(method, url string, data io.Reader, srv *Service, beforeServeHook func(*http.Request) error) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, data)

	r := gin.Default()

	srv.router(r)

	if err := beforeServeHook(req); err != nil {
		panic(err)
	}

	r.ServeHTTP(w, req)

	return w
}

func sendRequest(method, url string, data io.Reader, srv *Service, username, jwtSecret string, permissions int) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, data)

	if err := createJWTWithRole(username, jwtSecret, req, permissions); err != nil {
		log.Fatalf("Error creating JWT: %v", err)
	}

	//r := gin.Default()
	//
	//srv.router(r)
	//
	//r.ServeHTTP(w, req)
	srv.engine.ServeHTTP(w, req)

	return w
}

func sendRequestWithoutAuth(method, url string, data io.Reader, srv *Service) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, data)

	r := gin.Default()

	srv.router(r)

	r.ServeHTTP(w, req)

	return w
}

func createJWTWithRole(username, jwtSecret string, r *http.Request, role int) error {
	sreg := map[string]string{"nickname": username, "fullname": "JJ", "email": "jj@example.com"}
	resp := openid.Response{ID: "identity", Teams: []string{}, SReg: sreg}
	jwtToken, err := usso.NewJWTToken(jwtSecret, &resp, role)
	if err != nil {
		return fmt.Errorf("error creating a JWT: %v", err)
	}
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	return nil
}

func parseStandardResponse(r io.Reader) (web.StandardResponse, error) {
	// Parse the response
	result := web.StandardResponse{}
	err := json.NewDecoder(r).Decode(&result)
	return result, err
}
