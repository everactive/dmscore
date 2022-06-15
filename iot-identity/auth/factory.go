// Package auth provides a factory to create a gin.HandlerFunc based on specified auth provider
package auth

import (
	"net/http"
	"strings"

	"github.com/everactive/dmscore/iot-identity/config/configkey"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	keycloakauth "github.com/everactive/ginkeycloak"
	"github.com/gin-gonic/gin"
)

// Factory returns a gin middleware HandlerFunc from the provider specified by the provider string
func Factory(provider string) gin.HandlerFunc {
	if strings.ToLower(provider) == "keycloak" {
		clientID := viper.GetString(configkey.OAuth2ClientID)
		secret := viper.GetString(configkey.OAuth2ClientSecret)
		host := viper.GetString(configkey.OAuth2HostName)
		port := viper.GetString(configkey.OAuth2HostPort)
		scheme := viper.GetString(configkey.OAuth2HostScheme)
		tokenIntrospectPath := viper.GetString(configkey.OAuth2TokenIntrospectPath)
		requiredScope := viper.GetString(configkey.OAuth2ClientRequiredScope)

		a := keycloakauth.New(clientID, secret, host, port, scheme, requiredScope, tokenIntrospectPath, log.StandardLogger())
		return a.HandleFunc
	} else if strings.ToLower(provider) == "disabled" {
		log.Errorf("Running with API authentication disabled!")
		return func(context *gin.Context) {
			// disabled does nothing except call next
			context.Next()
		}
	}

	return failedAuth
}

func failedAuth(context *gin.Context) {
	context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed Authorization, No Auth Provider"})
}
