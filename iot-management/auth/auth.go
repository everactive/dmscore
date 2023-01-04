// Package auth provides functionality that allows the service to use and consume tokens from various sources
package auth

import (
	"errors"
	"fmt"
	"github.com/everactive/dmscore/config/keys"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"net/url"

	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/domain"
	"github.com/everactive/dmscore/iot-management/web"
	"github.com/everactive/ginkeycloak"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	tokenGetter *ginkeycloak.TokenGetter
	// Logger is a logger that can be set/specified independently for logging auth related issues
	Logger = log.StandardLogger()
)

// TokenGetter gets a TokenGetter and creates one if it does not already exist
func TokenGetter() *ginkeycloak.TokenGetter {
	if tokenGetter == nil {
		log.Tracef("tokenGetter == nil")
		clientID := viper.GetString(keys.OAuth2ClientID)
		clientSecret := viper.GetString(keys.OAuth2ClientSecret)

		scheme := "https"
		port := "443"
		if viper.GetString(keys.OAuth2HostScheme) != "" {
			scheme = viper.GetString(keys.OAuth2HostScheme)
		}
		if viper.GetString(keys.OAuth2HostPort) != "" {
			port = viper.GetString(keys.OAuth2HostPort)
		}

		host := viper.GetString(keys.OAuth2HostName)
		tokenPath := viper.GetString(keys.OAuth2AccessTokenPath)

		tokenAccessURL := url.URL{
			Scheme: scheme,
			Host:   fmt.Sprintf("%s:%s", host, port),
			Path:   tokenPath,
		}
		log.Tracef("Access Token URL: %s", tokenAccessURL.String())
		tokenGetter = ginkeycloak.NewGetter(clientID, clientSecret, tokenAccessURL.String(), Logger)
	}

	return tokenGetter
}

// CreateServiceClientUser creates a user account for a service-account if it does not exist previously
func CreateServiceClientUser(ms manage.Manage, clientName string) {
	log.Infof("Using an access token, checking to see if %s user exists", clientName)
	user, err := ms.GetUser(clientName)
	if err != nil {
		log.Infof("%s does not exist, creating", clientName)
		err := ms.CreateUser(domain.User{Username: clientName, Role: datastore.Superuser})
		if err != nil {
			panic(err)
		} else {
			log.Infof("user %s created", clientName)
		}
	} else {
		log.Infof("user exists: %+v", user)
	}
}

// VerifyKeycloakTokenWithAuth takes the Authorization string from the header and validates token and user, creating if necessary for service-accounts
func VerifyKeycloakTokenWithAuth(a *ginkeycloak.Auth) func(authorizationToken string, wb web.Service) (datastore.User, error) {
	return func(authorizationToken string, wb web.Service) (datastore.User, error) {
		valid, clientDetails, err := a.VerifyTokenFromHeader(authorizationToken)
		if err != nil {
			log.Error(err)
			return datastore.User{}, err
		}

		if valid && clientDetails != nil {
			_, err := wb.Manage.GetUser(clientDetails.ClientID)

			// User for this client exists, done and return
			if err == nil {
				return datastore.User{
					Username: clientDetails.ClientID,
					Role:     datastore.Superuser,
				}, nil
			}

			// If user doesn't exist, try to create one
			log.Error(err)

			err = wb.Manage.CreateUser(domain.User{
				Username: clientDetails.ClientID,
				Name:     clientDetails.ClientID,
				Role:     datastore.Superuser,
			})

			if err != nil {
				log.Error(err)
				return datastore.User{}, err
			}

			// We were able to create the user and they are valid and authorized
			return datastore.User{
				Username: clientDetails.ClientID,
				Role:     datastore.Superuser,
			}, nil
		}

		// We were able to create the user and they are valid and authorized
		return datastore.User{
			Username: clientDetails.ClientID,
			Role:     datastore.Superuser,
		}, nil
	}
}

var ErrUnauthorizedStaticClient = errors.New("authorization type static client but provided token was invalid, empty or token not set")

// nolint
// Deprecated: VerifyStaticClientToken verifies that a static token provided in the Authorization header if valid
func VerifyStaticClientToken(authorizationToken string, wb web.Service) (datastore.User, error) {
	staticClientToken := viper.GetString(keys.StaticClientToken)
	if authorizationToken == "" || len(staticClientToken) == 0 {
		return datastore.User{}, ErrUnauthorizedStaticClient
	}

	if authorizationToken != staticClientToken {
		return datastore.User{}, ErrUnauthorizedStaticClient
	}

	// we expect the static-client to exist before it is used
	_, err := wb.Manage.GetUser("static-client")
	if err != nil {
		return datastore.User{}, err
	}

	return datastore.User{
		Username: "static-client",
		Role:     datastore.Superuser,
	}, nil
}
