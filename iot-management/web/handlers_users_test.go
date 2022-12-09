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
	"github.com/everactive/dmscore/iot-management/domain"
	"github.com/everactive/dmscore/iot-management/service/manage/mocks"
	"github.com/everactive/dmscore/iot-management/web/usso"
	"github.com/juju/usso/openid"
	"github.com/stretchr/testify/mock"
	"net/http"
	"path"
	"testing"
)

var (
	username = "everactive"
)

func TestService_UserListHandler(t *testing.T) {
	username := "everactive"

	tests := []struct {
		name        string
		url         string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "/v1/users", 300, http.StatusOK, ""},
		{"invalid-permissions", "/v1/users", 200, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			manageMock.On("UserList").Return([]domain.User{}, nil)
			w := sendRequest("GET", tt.url, nil, wb, username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.UserListHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_UserCreateHandler(t *testing.T) {
	u1 := []byte(`{"username":"jane", "name":"Jane D", "email":"jd@example.com", "role":200}`)
	u2 := []byte(`{"username":"invalid", "name":"Invalid", "email":"jd@example.com", "role":200}`)
	u3 := []byte(``)
	u4 := []byte(`\u1000`)
	tests := []struct {
		name        string
		url         string
		permissions int
		data        []byte
		want        int
		wantErr     string
	}{
		{"valid", "/v1/users", 300, u1, http.StatusOK, ""},
		{"invalid-user", "/v1/users", 300, u2, http.StatusBadRequest, "UserAuth"},
		{"invalid-permissions", "/v1/users", 200, u1, http.StatusUnauthorized, "UserAuth"},
		{"invalid-empty", "/v1/users", 300, u3, http.StatusBadRequest, "UserAuth"},
		{"invalid-data", "/v1/users", 300, u4, http.StatusBadRequest, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)
			if tt.wantErr == "" {
				manageMock.On("CreateUser", mock.Anything).Return(nil)
			} else {
				manageMock.On("CreateUser", mock.Anything).Return(errors.New("some error text"))
			}

			w := sendRequest("POST", tt.url, bytes.NewReader(tt.data), wb, username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.UserListHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_UserGetHandler(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "/v1/users/jamesj", 300, http.StatusOK, ""},
		{"invalid-permissions", "/v1/users/jamesj", 200, http.StatusUnauthorized, "UserAuth"},
		{"invalid-user", "/v1/users/invalid", 300, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			_, username := path.Split(tt.url)
			if tt.wantErr == "" {
				manageMock.On("GetUser", username).Return(domain.User{}, nil)
			} else {
				manageMock.On("GetUser", username).Return(domain.User{}, errors.New("some error text"))
			}
			w := sendRequest("GET", tt.url, nil, wb, username, jwtSecret, tt.permissions)
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.UserListHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_UserUpdateHandler(t *testing.T) {
	u1 := []byte(`{"username":"jamesj", "name":"James Jone", "email":"jj@example.com", "role":200}`)
	u2 := []byte(`{"username":"invalid", "name":"Invalid", "email":"jd@example.com", "role":200}`)
	u3 := []byte(``)
	u4 := []byte(`\u1000`)
	tests := []struct {
		name          string
		url           string
		permissions   int
		data          []byte
		unmarshalUser bool
		want          int
		wantErr       string
	}{
		{"valid", "/v1/users/jamesj", 300, u1, true, http.StatusOK, ""},
		{"invalid-user", "/v1/users/invalid", 300, u2, true, http.StatusBadRequest, "UserUpdate"},
		{"invalid-permissions", "/v1/users/jamesj", 200, u1, true, http.StatusUnauthorized, "UserAuth"},
		{"invalid-empty", "/v1/users/jamesj", 300, u3, false, http.StatusBadRequest, "UserAuth"},
		{"invalid-data", "/v1/users/jamesj", 300, u4, false, http.StatusBadRequest, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			var user domain.User
			if tt.unmarshalUser {
				err := json.Unmarshal(tt.data, &user)
				if err != nil {
					t.Error(err)
				}
			}

			_, username := path.Split(tt.url)
			if tt.wantErr == "" {
				manageMock.On("UserUpdate", user).Return(nil)
			} else {
				manageMock.On("UserUpdate", user).Return(errors.New("some error text"))
			}

			w := sendRequestWithBeforeServeHook("PUT", tt.url, bytes.NewReader(tt.data), wb, func(request *http.Request) error {
				sreg := map[string]string{"nickname": username, "fullname": user.Name, "email": user.Email}
				resp := openid.Response{ID: "identity", Teams: []string{}, SReg: sreg}
				jwtToken, err := usso.NewJWTToken(jwtSecret, &resp, tt.permissions)
				if err != nil {
					return fmt.Errorf("error creating a JWT: %v", err)
				}
				request.Header.Set("Authorization", "Bearer "+jwtToken)
				return nil
			})

			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.UserUpdateHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func TestService_UserDeleteHandler(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		permissions int
		want        int
		wantErr     string
	}{
		{"valid", "/v1/users/jamesj", 300, http.StatusOK, ""},
		{"invalid-user", "/v1/users/invalid", 300, http.StatusBadRequest, "UserDelete"},
		{"invalid-permissions", "/v1/users/jamesj", 200, http.StatusUnauthorized, "UserAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jwtSecret := createAndSetJWTSecret(t)

			manageMock := &mocks.Manage{}
			wb := NewService(manageMock)

			_, username := path.Split(tt.url)
			user := domain.User{
				Username: username,
				Role:     tt.permissions,
			}

			var returnErr error
			if tt.wantErr != "" {
				returnErr = errors.New("some error text")
			}
			manageMock.On("UserDelete", username).Return(returnErr)

			w := sendRequestWithBeforeServeHook("DELETE", tt.url, nil, wb, setupServeWithUser(user, jwtSecret))
			if w.Code != tt.want {
				t.Errorf("Expected HTTP status '%d', got: %v", tt.want, w.Code)
			}

			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Error parsing response: %v", err)
			}
			if resp.Code != tt.wantErr {
				t.Errorf("Web.UserUpdateHandler() got = %v, want %v", resp.Code, tt.wantErr)
			}
		})
	}
}

func setupServeWithUser(user domain.User, jwtSecret string) func(request *http.Request) error {
	return func(request *http.Request) error {
		sreg := map[string]string{"nickname": username, "fullname": user.Name, "email": user.Email}
		resp := openid.Response{ID: "identity", Teams: []string{}, SReg: sreg}
		jwtToken, err := usso.NewJWTToken(jwtSecret, &resp, user.Role)
		if err != nil {
			return fmt.Errorf("error creating a JWT: %v", err)
		}
		request.Header.Set("Authorization", "Bearer "+jwtToken)
		return nil
	}
}
