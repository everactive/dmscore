package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/everactive/dmscore/config"
	"github.com/everactive/dmscore/config/keys"
	devicetwinconfig "github.com/everactive/dmscore/iot-devicetwin/config"
	"github.com/everactive/dmscore/iot-devicetwin/service/controller"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	devicetwinfactory "github.com/everactive/dmscore/iot-devicetwin/service/factory"
	devicetwinweb "github.com/everactive/dmscore/iot-devicetwin/web"
	identityauth "github.com/everactive/dmscore/iot-identity/auth"
	"github.com/everactive/dmscore/iot-identity/config/configkey"
	identitydatastore "github.com/everactive/dmscore/iot-identity/datastore"
	"github.com/everactive/dmscore/iot-identity/middleware/logger"
	"github.com/everactive/dmscore/iot-identity/service"
	"github.com/everactive/dmscore/iot-identity/service/cert"
	identityfactory "github.com/everactive/dmscore/iot-identity/service/factory"
	identityweb "github.com/everactive/dmscore/iot-identity/web"
	"github.com/everactive/dmscore/iot-management/auth"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/identityapi"
	"github.com/everactive/dmscore/iot-management/service/factory"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/everactive/dmscore/iot-management/twinapi"
	"github.com/everactive/dmscore/iot-management/web"
	"github.com/everactive/ginkeycloak"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func init() {
	Root.AddCommand(&runCommand)
}

var Root = cobra.Command{
	Use:   "dmscore",
	Short: "dmscore",
	Long:  "dmscore",
}

var runCommand = cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  "run",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.TraceLevel)

		var configFilePath string
		if filePath, ok := os.LookupEnv("CONFIG_FILE_PATH"); ok {
			configFilePath = filePath
		}

		config.LoadConfig(configFilePath)

		// Open the connection to the local database
		databaseDriver := viper.GetString(keys.DatabaseDriver)
		dataSource := viper.GetString(keys.DatabaseConnectionString)
		log.Infof("Connecting to %s with connection string %s", databaseDriver, dataSource)
		db, err := factory.CreateDataStore(databaseDriver, dataSource)
		if err != nil || db == nil {
			log.Fatalf("Error accessing data store: %v, with database source %s", databaseDriver, dataSource)
			return
		}

		deviceTwinAPIURL := viper.GetString(keys.DeviceTwinAPIURL)
		// Initialize the device twin client
		twinAPI, err := twinapi.NewClientAdapter(deviceTwinAPIURL)
		if err != nil {
			log.Fatalf("Error connecting to the device twin service: %v", err)
		}

		identityAPIURL := viper.GetString(keys.IdentityAPIURL)
		// Initialize the identity client
		identityAPI, err := identityapi.NewClientAdapter(identityAPIURL)
		if err != nil {
			log.Fatalf("Error connecting to the identity service: %v", err)
		}

		ids := createIdentityService()

		dts := createDeviceTwinService(db)

		// Create the main services
		srv := manage.NewManagement(db, twinAPI, identityAPI, dts, ids)

		// Figure out what our auth provider is (keycloak or legacy static client token)
		authProvider := strings.ToLower(viper.GetString(keys.AuthProvider))
		authDisabled := viper.GetBool(keys.DisableAuth)
		if authProvider == "disabled" && authDisabled {
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
					auth.CreateServiceClientUser(db, "static-client")
					web.VerifyTokenAndUser = auth.VerifyStaticClientToken //nolint
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

		// Start the web service
		www := web.NewService(srv)

		web.AuthMiddleWare = func(c *gin.Context) {
			user := datastore.User{
				Username: "static-client",
				Role:     datastore.Superuser,
			}

			c.Set("USER", &user)
		}

		www.Run()
	},
}

func createDeviceTwinService(coreDB datastore.DataStore) *devicetwinweb.Service {
	databaseDriver := viper.GetString(keys.GetDeviceTwinKey(keys.DatabaseDriver))
	dataStoreSource := viper.GetString(keys.GetDeviceTwinKey(keys.DatabaseConnectionString))

	db, err := devicetwinfactory.CreateDataStore(databaseDriver, dataStoreSource)
	if err != nil {
		log.Fatalf("Error connecting to data store: %v", err)
	}

	URL := viper.GetString(keys.MQTTURL)
	port := viper.GetString(keys.MQTTPort)

	certsDir := viper.GetString(keys.CertificatesPath)

	rootCAFilename := viper.GetString(keys.MQTTRootCAFilename)
	clientCertFilename := viper.GetString(keys.MQTTClientCertificateFilename)
	clientKeyFilename := viper.GetString(keys.MQTTClientKeyFilename)

	c := devicetwinconfig.MQTTConnect{}

	// nolint: gosec
	rootCA, err := ioutil.ReadFile(path.Join(certsDir, rootCAFilename))
	if err != nil {
		panic(err)
	}

	// nolint: gosec
	certFile, err := ioutil.ReadFile(path.Join(certsDir, clientCertFilename))
	if err != nil {
		panic(err)
	}

	// nolint: gosec
	key, err := ioutil.ReadFile(path.Join(certsDir, clientKeyFilename))

	c.RootCA = rootCA
	c.ClientKey = key
	c.ClientCert = certFile

	prefix := viper.GetString(keys.MQTTClientIDPrefix)

	// Generate a random string
	s, err := cert.CreateSecret(6)
	if err != nil {
		log.Printf("Error creating client ID: %v", err)
	}

	c.ClientID = fmt.Sprintf("%s-%s", prefix, s)

	twin := devicetwin.NewService(db, coreDB)
	ctrl := controller.NewService(URL, port, &c, twin)

	servicePort := viper.GetString(keys.GetDeviceTwinKey(keys.ServicePort))

	// Start the web API service
	w := devicetwinweb.NewService(servicePort, ctrl)
	// log.Fatal(w.Run())
	return w
}

func runMigrations(datasource, sourceURL string) error {
	db, err := sql.Open("postgres", datasource)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver)

	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
		return err
	}

	return nil
}

func CreateIdentityDataStore() (identitydatastore.DataStore, error) {
	// Open the connection to the local database
	databaseDriver := viper.GetString(keys.GetIdentityKey(keys.DatabaseDriver))
	dataStoreSource := viper.GetString(keys.GetIdentityKey(keys.DatabaseConnectionString))

	log.Infof("Connecting to %s with connection string %s", databaseDriver, dataStoreSource)
	identitydatastore.Logger = log.StandardLogger()
	db, err := identityfactory.CreateDataStore(databaseDriver, dataStoreSource)
	if err != nil || db == nil {
		log.Fatalf("Error accessing data store: %v, with database source %s", databaseDriver, dataStoreSource)
		// This satisfies an IDE's nil logic check
		return nil, err
	}

	return db, nil
}
func createIdentityService() *identityweb.IdentityService {
	db, err := CreateIdentityDataStore()
	if err != nil || db == nil {
		log.Fatalf("Error accessing data store: %v", err)
		// This satisfies an IDE's nil logic check
		return nil
	}

	service.Logger = log.StandardLogger()
	srv := service.NewIdentityService(db)

	// Start the web service
	identityweb.Logger = log.StandardLogger()
	wb := identityweb.NewIdentityService(srv, log.StandardLogger())

	// log.Fatal(w.Run())
	internalPort := viper.GetString(keys.GetIdentityKey(keys.ServicePortInternal))
	enrollPort := viper.GetString(keys.GetIdentityKey(keys.ServicePortEnroll))

	log.Info("Starting service (internal) on port : ", internalPort)
	log.Info("Starting service (enroll) on port : ", enrollPort)

	internalRouter := gin.New()
	enrollRouter := gin.New()

	logFormat := os.Getenv("LOG_FORMAT")
	if strings.ToUpper(logFormat) == "JSON" {
		log.Infof("Setting up JSON log format for logger middleware")

		middlewareLogger := logger.New(log.StandardLogger(), logger.LogOptions{EnableStarting: true})

		internalRouter.Use(middlewareLogger.HandleFunc)
		enrollRouter.Use(middlewareLogger.HandleFunc)

	} else {
		internalRouter.Use(gin.Logger())
		enrollRouter.Use(gin.Logger())
	}

	internalRouter.Use(identityauth.Factory(viper.GetString(configkey.AuthProvider)))

	wb.SetRouters(internalRouter, enrollRouter)

	// Use a goroutine for the internal serve, we'll block with the enroll serve
	//go func() {
	//	log.Info("Listening and serving internal on :" + internalPort)
	//
	//	err := internalRouter.Run(":" + internalPort)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	go func() {
		log.Info("Listening and serving enroll on :" + enrollPort)

		err = enrollRouter.Run(":" + enrollPort)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return wb
}
