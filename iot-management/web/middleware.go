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
	"errors"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/everactive/dmscore/iot-management/config/configkey"
	"github.com/everactive/dmscore/iot-management/web/usso"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// JWTCheck extracts the JWT from the request, validates it and returns the token
func (wb Service) JWTCheck(c *gin.Context) (*jwt.Token, error) {
	// Get the JWT from the header or cookie
	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		var err error
		jwtToken, err = c.Cookie(usso.JWTCookie)
		if err != nil {
			log.Println("Error in JWT extraction:", err.Error())
			return nil, errors.New("error in retrieving the authentication token")
		}
	}

	trimmedToken := strings.TrimPrefix(jwtToken, "Bearer ")
	jwtSecret := viper.GetString(configkey.JwtSecret)

	// Verify the JWT string
	token, err := usso.VerifyJWT(jwtSecret, trimmedToken)
	if err != nil {
		log.Printf("JWT fails verification: %v", err.Error())
		return nil, errors.New("the authentication token is invalid")
	}

	if !token.Valid {
		log.Println("Invalid JWT")
		return nil, errors.New("the authentication token is invalid")
	}

	// Set up the bearer token in the header
	c.Header("Authorization", "Bearer "+jwtToken)

	return token, nil
}
