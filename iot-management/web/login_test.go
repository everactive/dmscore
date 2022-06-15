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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/everactive/dmscore/iot-management/config"
	"github.com/everactive/dmscore/iot-management/config/configkey"
	"github.com/everactive/dmscore/iot-management/datastore/memory"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/juju/usso"
	"github.com/spf13/viper"
)

type TestSSOServer struct {
	TokenIsValid bool
	Accounts     string
	ReturnError  error
}

func (t *TestSSOServer) IsTokenValid(_ *usso.SSOData) (bool, error) {
	return t.TokenIsValid, t.ReturnError
}

func (t *TestSSOServer) GetAccounts(_ *usso.SSOData) (string, error) {
	return t.Accounts, t.ReturnError
}

func TestLoginHandlerAPIClient(t *testing.T) {
	// Mock the services
	config.LoadConfig("../testing/memory.yaml")
	db := memory.NewStore()
	m := manage.NewMockManagement(db)
	wb := NewService(m)

	ssodata := usso.SSOData{
		ConsumerKey:    "consumer-key",
		ConsumerSecret: "consumer-secret",
		Realm:          "real",
		TokenKey:       "token-key",
		TokenName:      "token-name",
		TokenSecret:    "token-secret",
	}

	ssodatabytes, _ := json.Marshal(&ssodata)
	bodyReader := bytes.NewReader(ssodatabytes)

	ts := TestSSOServer{
		TokenIsValid: true,
		Accounts:     `{ "username": "jamesj", "email": "jamesj@example.com", "displayname": "James J" }`,
	}

	ssoServer = &ts
	w := sendRequestWithoutAuth("GET", "/v1/login", bodyReader, wb)

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status '200', got: %v", w.Code)
	}
}

func TestLoginHandlerUserNotFound(t *testing.T) {
	// Mock the services
	config.LoadConfig("../testing/memory.yaml")
	db := memory.NewStore()
	m := manage.NewMockManagement(db)
	wb := NewService(m)

	ssodata := usso.SSOData{
		ConsumerKey:    "consumer-key",
		ConsumerSecret: "consumer-secret",
		Realm:          "real",
		TokenKey:       "token-key",
		TokenName:      "token-name",
		TokenSecret:    "token-secret",
	}

	ssodatabytes, _ := json.Marshal(&ssodata)
	bodyReader := bytes.NewReader(ssodatabytes)

	ts := TestSSOServer{
		TokenIsValid: true,
		Accounts:     `{ "username": "franktester", "email": "frank@example.com", "displayname": "Frank Tester" }`,
	}

	ssoServer = &ts
	w := sendRequest("GET", "/v1/login", bodyReader, wb, "jamesj", viper.GetString(configkey.JwtSecret), 100)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected HTTP status '200', got: %v", w.Code)
	}

	if w.Header().Get("Location") != "/notfound" {
		t.Errorf("Expected /notfound for Location, got: %v", w.Header().Get("Location"))
	}
}

func TestLoginHandlerAccountsUserDNE(t *testing.T) {
	// Mock the services
	config.LoadConfig("../testing/memory.yaml")
	db := memory.NewStore()
	m := manage.NewMockManagement(db)
	wb := NewService(m)

	ssodata := usso.SSOData{
		ConsumerKey:    "consumer-key",
		ConsumerSecret: "consumer-secret",
		Realm:          "real",
		TokenKey:       "token-key",
		TokenName:      "token-name",
		TokenSecret:    "token-secret",
	}

	ssodatabytes, _ := json.Marshal(&ssodata)
	bodyReader := bytes.NewReader(ssodatabytes)

	ts := TestSSOServer{
		TokenIsValid: true,
		Accounts:     `{ "username": "frank", "email": "frank@example.com" }`,
	}

	ssoServer = &ts
	w := sendRequest("GET", "/v1/login", bodyReader, wb, "frank", viper.GetString(configkey.JwtSecret), 100)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected HTTP status '307', got: %v", w.Code)
	}
}

func TestLoginHandlerAccountsError(t *testing.T) {
	// Mock the services
	config.LoadConfig("../testing/memory.yaml")
	db := memory.NewStore()
	m := manage.NewMockManagement(db)
	wb := NewService(m)

	ssodata := usso.SSOData{
		ConsumerKey:    "consumer-key",
		ConsumerSecret: "consumer-secret",
		Realm:          "real",
		TokenKey:       "token-key",
		TokenName:      "token-name",
		TokenSecret:    "token-secret",
	}

	ssodatabytes, _ := json.Marshal(&ssodata)
	bodyReader := bytes.NewReader(ssodatabytes)

	ts := TestSSOServer{
		TokenIsValid: true,
	}

	ssoServer = &ts
	ts.ReturnError = errors.New("there was an error getting accounts")
	w := sendRequest("GET", "/v1/login", bodyReader, wb, "jamesj", viper.GetString(configkey.JwtSecret), 100)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected HTTP status '400', got: %v", w.Code)
	}

	bodyReader = bytes.NewReader(ssodatabytes)
	ts.ReturnError = nil
	ts.Accounts = "......"
	w = sendRequest("GET", "/v1/login", bodyReader, wb, "jamesj", viper.GetString(configkey.JwtSecret), 100)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected HTTP status '401', got: %v", w.Code)
	}
}

func TestLoginHandlerTokenInvalidOrError(t *testing.T) {
	// Mock the services
	config.LoadConfig("../testing/memory.yaml")
	db := memory.NewStore()
	m := manage.NewMockManagement(db)
	wb := NewService(m)

	ssodata := usso.SSOData{
		ConsumerKey:    "consumer-key",
		ConsumerSecret: "consumer-secret",
		Realm:          "real",
		TokenKey:       "token-key",
		TokenName:      "token-name",
		TokenSecret:    "token-secret",
	}

	ssodatabytes, _ := json.Marshal(&ssodata)
	bodyReader := bytes.NewReader(ssodatabytes)

	ts := TestSSOServer{
		TokenIsValid: false,
		Accounts:     "",
	}

	ssoServer = &ts
	w := sendRequest("GET", "/v1/login", bodyReader, wb, "jamesj", viper.GetString(configkey.JwtSecret), 100)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected HTTP status '400', got: %v", w.Code)
	}

	bodyReader = bytes.NewReader(ssodatabytes)
	ts.ReturnError = errors.New("there was an error checking if token is invalid")
	w = sendRequest("GET", "/v1/login", bodyReader, wb, "jamesj", viper.GetString(configkey.JwtSecret), 100)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected HTTP status '401', got: %v", w.Code)
	}
}

func TestLoginHandlerAPIClientNoBodyOrMalformed(t *testing.T) {
	// Mock the services
	config.LoadConfig("../testing/memory.yaml")
	db := memory.NewStore()
	m := manage.NewMockManagement(db)
	wb := NewService(m)

	ts := TestSSOServer{
		TokenIsValid: false,
		Accounts:     "",
	}

	ssoServer = &ts

	w := sendRequest("GET", "/v1/login", nil, wb, "jamesj", viper.GetString(configkey.JwtSecret), 100)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP status '400', got: %v", w.Code)
	}

	body := `................`
	bodyReader := bytes.NewReader([]byte(body))

	w = sendRequest("GET", "/v1/login", bodyReader, wb, "jamesj", viper.GetString(configkey.JwtSecret), 100)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP status '400', got: %v", w.Code)
	}

	body = `{ "malformed": "body-doesnt-make-sense" }`
	bodyReader = bytes.NewReader([]byte(body))

	w = sendRequest("GET", "/v1/login", bodyReader, wb, "jamesj", viper.GetString(configkey.JwtSecret), 100)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP status '400', got: %v", w.Code)
	}

	ts.Accounts = `{ "username": "jamesj", "email": "jamesj@example.com" }`
	ts.TokenIsValid = true

	ssodata := usso.SSOData{
		ConsumerKey:    "consumer-key",
		ConsumerSecret: "consumer-secret",
		Realm:          "real",
		TokenKey:       "token-key",
		TokenName:      "token-name",
		TokenSecret:    "token-secret",
	}
	ssodatabytes, _ := json.Marshal(&ssodata)
	bodyReader = bytes.NewReader(ssodatabytes)
	w = sendRequest("GET", "/v1/login", bodyReader, wb, "jamesj", viper.GetString(configkey.JwtSecret), 100)
	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP status '200', got: %v", w.Code)
	}
}

func TestLoginHandlerUSSORedirect(t *testing.T) {
	// Mock the services
	config.LoadConfig("../testing/memory.yaml")
	db := memory.NewStore()
	m := manage.NewMockManagement(db)
	wb := NewService(m)

	w := sendRequest("GET", "/login", nil, wb, "jamesj", viper.GetString(configkey.JwtSecret), 100)

	if w.Code != http.StatusFound {
		t.Errorf("Expected HTTP status '302', got: %v", w.Code)
	}

	u, err := url.Parse(w.Header().Get("Location"))
	if err != nil {
		t.Errorf("Error Parsing the redirect URL: %v", u)
		return
	}

	// Check that the redirect is to the Ubuntu SSO service
	ssoURL := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	if ssoURL != usso.ProductionUbuntuSSOServer.LoginURL() {
		t.Errorf("Unexpected redirect URL: %v", ssoURL)
	}
}
