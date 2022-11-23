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
	"errors"
	"fmt"
	"github.com/everactive/dmscore/config/keys"
	"html/template"
	"net/http"

	"github.com/everactive/dmscore/iot-management/datastore"
	webusso "github.com/everactive/dmscore/iot-management/web/usso"
	"github.com/gin-gonic/gin"
	"github.com/juju/usso"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Account is the account structure for a USSO user based on the token
type Account struct {
	Email       string
	Username    string
	DisplayName string `json:"displayname"`
}

// UbuntuSingleSignOnServer is the interface for minimal USSO functionality
type UbuntuSingleSignOnServer interface {
	IsTokenValid(ssodata *usso.SSOData) (bool, error)
	GetAccounts(ssodata *usso.SSOData) (string, error)
}

func (wb Service) getSSOData(context *gin.Context) (*usso.SSOData, int, error) {
	// need to get body with OAuth token information
	var ssodata usso.SSOData
	if context.Request.Body == nil {
		return nil, http.StatusBadRequest, errors.New("body of request is nil")
	}
	err := json.NewDecoder(context.Request.Body).Decode(&ssodata)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// If any of these values are empty we can't continue.
	if ssodata.TokenName == "" || ssodata.TokenKey == "" || ssodata.ConsumerKey == "" || ssodata.ConsumerSecret == "" || ssodata.Realm == "" {
		return nil, http.StatusBadRequest, errors.New("all values for SSOData must have values")
	}

	// need to validate token
	valid, err2 := ssoServer.IsTokenValid(&ssodata)
	if err2 != nil {
		return nil, http.StatusUnauthorized, err2
	}

	if !valid {
		return nil, http.StatusUnauthorized, errors.New("token is invalid")
	}

	return &ssodata, http.StatusOK, nil
}

// LoginHandlerAPIClient processes the login for API client
func (wb Service) LoginHandlerAPIClient(c *gin.Context) {
	ssodata, status, err := wb.getSSOData(c)
	if err != nil {
		replyHTTPError(c.Writer, status, err)
		return
	}

	accountsString, err3 := ssoServer.GetAccounts(ssodata)
	if err3 != nil {
		replyHTTPError(c.Writer, http.StatusUnauthorized, errors.New("unable to get accounts for token"))
		return
	}

	var accts Account
	err4 := json.Unmarshal([]byte(accountsString), &accts)
	if err4 != nil {
		replyHTTPError(c.Writer, http.StatusUnauthorized, errors.New("unable to get accounts for token"))
		return
	}

	// Check that the user is registered
	user, err := wb.Manage.GetUser(accts.Username)
	if err != nil {
		// Cannot find the user, so redirect to the login page
		log.Printf("Error retrieving user from datastore: %v\n", err)
		http.Redirect(c.Writer, c.Request, "/notfound", http.StatusTemporaryRedirect)
		return
	}

	// Verify role value is valid
	if user.Role != datastore.Standard && user.Role != datastore.Admin && user.Role != datastore.Superuser {
		log.Printf("Role obtained from database for user %v has not a valid value: %v\n", accts.Username, user.Role)
		http.Redirect(c.Writer, c.Request, "/notfound", http.StatusTemporaryRedirect)
		return
	}

	jwtSecret := viper.GetString(keys.JwtSecret)

	// Build the JWT
	jwtToken, err := webusso.NewJWTTokenForClient(jwtSecret, accts.Username, accts.DisplayName, accts.Email, ssodata.TokenKey, user.Role)
	if err != nil {
		// Unexpected that this should occur, so leave the detailed response
		log.Printf("Error creating the JWT: %v", err)
		replyHTTPError(c.Writer, http.StatusBadRequest, err)
		return
	}
	c.Writer.Header().Set(webusso.JWTCookie, jwtToken)
}

// LoginHandler processes the login for Ubuntu SSO
func (wb Service) LoginHandler(c *gin.Context) {
	// Get the openid nonce store
	nonce := wb.Manage.OpenIDNonceStore()

	w := c.Writer
	r := c.Request

	// Call the openid login handler
	resp, req, username, err := webusso.LoginHandler(nonce, w, r)
	if err != nil {
		log.Printf("Error verifying the OpenID response: %v\n", err)
		replyHTTPError(w, http.StatusBadRequest, err)
		return
	}
	if req != nil {
		// Redirect is handled by the SSO handler
		return
	}

	// Check that the user is registered
	user, err := wb.Manage.GetUser(username)
	if err != nil {
		// Cannot find the user, so redirect to the login page
		log.Printf("Error retrieving user from datastore: %v\n", err)
		http.Redirect(w, r, "/notfound", http.StatusTemporaryRedirect)
		return
	}

	// Verify role value is valid
	if user.Role != datastore.Standard && user.Role != datastore.Admin && user.Role != datastore.Superuser {
		log.Printf("Role obtained from database for user %v has not a valid value: %v\n", username, user.Role)
		http.Redirect(w, r, "/notfound", http.StatusTemporaryRedirect)
		return
	}

	jwtSecret := viper.GetString(keys.JwtSecret)

	// Build the JWT
	jwtToken, err := webusso.NewJWTToken(jwtSecret, resp, user.Role)
	if err != nil {
		// Unexpected that this should occur, so leave the detailed response
		log.Printf("Error creating the JWT: %v", err)
		replyHTTPError(w, http.StatusBadRequest, err)
		return
	}

	// Set a cookie with the JWT
	webusso.AddJWTCookie(jwtToken, w)

	frontendHost := viper.GetString(keys.FrontEndHost)
	frontendScheme := viper.GetString(keys.FrontEndScheme)

	// Both values have to exist and have values for the alternate redirect
	if frontendScheme != "" && frontendHost != "" {
		// Redirect to the homepage with the JWT
		http.Redirect(w, r, fmt.Sprintf("%s://%s/", frontendScheme, frontendHost), http.StatusTemporaryRedirect)
	} else {
		// Redirect to the homepage with the JWT
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

// LogoutHandler logs the user out by removing the cookie and the JWT authorization header
func (wb Service) LogoutHandler(c *gin.Context) {
	webusso.LogoutHandler(c.Writer, c.Request)
}

func replyHTTPError(w http.ResponseWriter, returnCode int, err error) {
	w.Header().Set("ContentType", "text/html")
	w.WriteHeader(returnCode)
	err = errorTemplate.Execute(w, err)
	if err != nil {
		log.Error(err)
	}
}

var errorTemplate = template.Must(template.New("failure").Parse(`<html>
 <head><title>Login Error</title></head>
 <body>{{.}}</body>
 </html>
 `))
