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
	"github.com/everactive/dmscore/config/keys"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/web/usso"
	"github.com/gin-gonic/gin"
)

func (wb Service) checkIsStandardAndGetUserFromJWT(c *gin.Context) (datastore.User, error) {
	return wb.checkPermissionsAndGetUserFromJWT(c, datastore.Standard)
}

// VerifyTokenAndUser is a variable function to verify the token and extract the user based on the current provider
var VerifyTokenAndUser = func(authorizationToken string, wb Service) (datastore.User, error) {
	return datastore.User{}, errors.New("service account authorization not configured")
}

func (wb Service) checkPermissionsAndGetUserFromJWT(c *gin.Context, minimumAuthorizedRole int) (datastore.User, error) {
	authType := c.GetHeader("Auth-Type")
	authProvider := strings.ToLower(viper.GetString(keys.AuthProvider))
	log.Tracef("Auth provider: %s, Auth-Type: %s", authProvider, authType)

	if (authProvider == "static-client" || authProvider == "keycloak") && (authType == "static-client" || authType == "keycloak") {
		token := c.GetHeader("Authorization")
		log.Tracef("Authorization token: %s", token)
		return VerifyTokenAndUser(token, wb)
	}

	user, err := wb.getUserFromJWT(c)
	if err != nil {
		return user, err
	}
	err = checkUserPermissions(user, minimumAuthorizedRole)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (wb Service) getUserFromJWT(c *gin.Context) (datastore.User, error) {
	token, err := wb.JWTCheck(c)
	if err != nil {
		return datastore.User{}, err
	}

	// Null token is invalid
	if token == nil {
		return datastore.User{}, errors.New("no JWT provided")
	}

	claims := token.Claims.(jwt.MapClaims)
	username := claims[usso.ClaimsUsername].(string)
	role := int(claims[usso.ClaimsRole].(float64))

	return datastore.User{
		Username: username,
		Role:     role,
	}, nil
}

func checkUserPermissions(user datastore.User, minimumAuthorizedRole int) error {
	if user.Role < minimumAuthorizedRole {
		return errors.New("the user is not authorized")
	}
	return nil
}

func GetUserFromContextAndCheckPermissions(c *gin.Context, minimumAuthorizableRole int) (*datastore.User, error) {
	return getUserFromContextAndCheckPermissions(c, minimumAuthorizableRole)
}

func getUserFromContextAndCheckPermissions(c *gin.Context, minimumAuthorizableRole int) (*datastore.User, error) {
	userInterface, exists := c.Get("USER")
	if exists {
		if user, ok := userInterface.(*datastore.User); ok {
			err := checkUserPermissions(*user, minimumAuthorizableRole)
			return user, err
		}
	}

	return nil, errors.New("user context not found")
}
