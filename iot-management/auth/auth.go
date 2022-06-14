// Package auth provides functionality that allows the service to use and consume tokens from various sources
package auth

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/everactive/dmscore/iot-management/config/configkey"
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
		clientID := viper.GetString(configkey.OAuth2ClientID)
		clientSecret := viper.GetString(configkey.OAuth2ClientSecret)

		scheme := "https"
		port := "443"
		if viper.GetString(configkey.OAuth2HostScheme) != "" {
			scheme = viper.GetString(configkey.OAuth2HostScheme)
		}
		if viper.GetString(configkey.OAuth2HostPort) != "" {
			port = viper.GetString(configkey.OAuth2HostPort)
		}

		host := viper.GetString(configkey.OAuth2HostName)
		tokenPath := viper.GetString(configkey.OAuth2AccessTokenPath)

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
func CreateServiceClientUser(ds datastore.DataStore, clientName string) {
	log.Infof("Using an access token, checking to see if %s user exists", clientName)
	user, err := ds.GetUser(clientName)
	if err != nil {
		log.Infof("%s does not exist, creating", clientName)
		createdUser, err := ds.CreateUser(datastore.User{
			Username: clientName,
			Role:     datastore.Superuser,
		})
		if err != nil {
			panic(err)
		} else {
			log.Infof("user created: %+v", createdUser)
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

//nolint
// Deprecated: VerifyStaticClientToken verifies that a static token provided in the Authorization header if valid
func VerifyStaticClientToken(authorizationToken string, wb web.Service) (datastore.User, error) {
	if authorizationToken != "" {
		staticClientToken := viper.GetString(configkey.StaticClientToken)
		if len(staticClientToken) > 0 {
			if authorizationToken == staticClientToken {
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
		}
	}

	return datastore.User{}, errors.New("authorization type static client but token was invalid")
}
