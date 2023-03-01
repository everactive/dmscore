package web

import (
	"context"
	"github.com/everactive/dmscore/config/keys"
	"github.com/everactive/dmscore/iot-identity/config/configkey"
	"github.com/everactive/dmscore/iot-management/auth"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/everactive/dmscore/iot-management/web"
	"github.com/everactive/ginkeycloak"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thejerf/suture/v4"
	"gorm.io/gorm"
	"strings"
	"time"
)

type HandlerService struct {
	manage    manage.Manage
	legacyWeb *web.Service
	engine    *gin.Engine
	db        *gorm.DB
}

func New(srv manage.Manage, db *gorm.DB) *suture.Supervisor {
	sup := suture.NewSimple("webhandlers")

	// Figure out what our auth provider is (keycloak or legacy static client token)
	authProvider := strings.ToLower(viper.GetString(keys.AuthProvider))
	authDisabled := viper.GetBool(keys.DisableAuth)
	if authProvider == "disabled" && authDisabled {
		log.Infof("Auth is disabled and auth provider is set to disabled, using static-client for requests with no auth checking")
		web.VerifyTokenAndUser = func(authorizationToken string, wb web.Service) (datastore.User, error) {
			return datastore.User{
				Username: "static-client",
				Role:     datastore.Superuser,
			}, nil
		}
	} else {
		log.Infof("Auth provider: %s", authProvider)
		if authProvider == "static-client" {
			staticClientToken := viper.GetString(keys.StaticClientToken)
			if staticClientToken != "" {
				auth.CreateServiceClientUser(srv, "static-client")
				web.VerifyTokenAndUser = auth.VerifyStaticClientToken //nolint
			} else {
				log.Error("Static client token is empty, not properly configured for using static client")
			}
		} else if authProvider == "keycloak" {
			clientID := viper.GetString(configkey.OAuth2ClientID)
			secret := viper.GetString(configkey.OAuth2ClientSecret)
			host := viper.GetString(configkey.OAuth2HostName)
			port := viper.GetString(configkey.OAuth2HostPort)
			scheme := viper.GetString(configkey.OAuth2HostScheme)
			tokenIntrospectPath := viper.GetString(configkey.OAuth2TokenIntrospectPath)
			requiredScope := viper.GetString(configkey.OAuth2ClientRequiredScope)

			a := ginkeycloak.New(clientID, secret, host, port, scheme, requiredScope, tokenIntrospectPath, log.StandardLogger())
			web.VerifyTokenAndUser = auth.VerifyKeycloakTokenWithAuth(a)
		}
	}

	engine := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	legacyWeb := web.NewService(srv, engine)
	hs := &HandlerService{srv, legacyWeb, engine, db}

	sup.Add(hs)
	sup.Add(legacyWeb)

	group := engine.Group("/v1")
	group.Use(web.AuthMiddleWare)
	group.POST("/:orgid/models/:model/required", hs.AddRequiredModelSnap)
	group.DELETE("/:orgid/models/:model/required", hs.DeleteRequiredModelSnap)
	group.GET("/:orgid/models/:model/required", hs.RequiredModelSnaps)

	return sup
}

func (h *HandlerService) Serve(ctx context.Context) error {
	intervalTicker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Errorf("We're done: %s", ctx.Err())
			return nil
		case <-intervalTicker.C:
			log.Infof("%s still ticking", "ManagementHandlerService")
		}
	}

	return nil
}
