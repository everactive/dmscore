package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everactive/dmscore/config"
	"github.com/everactive/dmscore/config/keys"
	devicetwinconfig "github.com/everactive/dmscore/iot-devicetwin/config"
	"github.com/everactive/dmscore/iot-devicetwin/service/controller"
	"github.com/everactive/dmscore/iot-devicetwin/service/devicetwin"
	devicetwinfactory "github.com/everactive/dmscore/iot-devicetwin/service/factory"
	devicetwinweb "github.com/everactive/dmscore/iot-devicetwin/web"
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
	migrate2 "github.com/everactive/dmscore/pkg/migrate"
	web2 "github.com/everactive/dmscore/pkg/web"
	"github.com/everactive/ginkeycloak"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	// this is needed for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// this is needed for migrate
	_ "github.com/lib/pq"

	"github.com/thejerf/suture/v4"
)

func init() {
	Root.AddCommand(&runCommand)
	Root.AddCommand(&version)

	Root.AddCommand(&createSuperuser)
	createSuperuser.Flags().String("username", "", "The username of the user to create (must match Ubuntu SSO)")
	createSuperuser.Flags().String("name", "", "The name of the user to create (must match Ubuntu SSO)")
	createSuperuser.Flags().String("email", "", "The email address of the user to create (must match Ubuntu SSO)")
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
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetLevel(log.TraceLevel)

		var configFilePath string
		if filePath, ok := os.LookupEnv("CONFIG_FILE_PATH"); ok {
			configFilePath = filePath
		}

		config.LoadConfig(configFilePath)

		loadStoreIDs()

		db, err := createManagementDatastore()

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
		srv := manage.NewManagement(db, twinAPI, identityAPI, dts.Controller, ids.Identity)

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

		supervisorSpec := suture.Spec{}
		sup := suture.New("dmscore", supervisorSpec)

		sup.Add(web2.New(srv))

		ctx := context.Background()
		ctx, cancelCtx := context.WithCancel(ctx)

		err = <-sup.ServeBackground(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Errorf("Well, this isn't good, we existed with an error and we were not canceled: %s", err)
		} else {
			log.Infof("Context was canceled, exiting")
		}

		quitSignals := make(chan os.Signal, 1)
		signal.Notify(quitSignals, syscall.SIGINT, syscall.SIGTERM)

		interruptSignal := <-quitSignals

		cancelCtx()

		// SIGINT is expected from systemd and should not result in an error exit
		if interruptSignal == syscall.SIGINT || interruptSignal == syscall.SIGTERM {
			// We expect this signal, it is not an error
			log.Infof("Caught SIGINT or SIGTERM, exiting...")
			return nil
		}

		return fmt.Errorf("caught signal: %v, %d", interruptSignal, interruptSignal)
	},
}

func loadStoreIDs() {
	var modelStoreIds []struct {
		Model   string `json:"model"`
		StoreID string `json:"storeid"`
	}
	storeIDsString := viper.GetString(keys.StoreIDs)
	log.Tracef("Store IDs string: %s", storeIDsString)
	if storeIDsString != "" {
		err := json.Unmarshal([]byte(storeIDsString), &modelStoreIds)
		if err != nil {
			log.Error(err)
		}
		for _, item := range modelStoreIds {
			log.Tracef("For model %s, loading store id %s", item.Model, item.StoreID)
			viper.Set(fmt.Sprintf(keys.ModelKeyTemplate, item.Model), item.StoreID)
		}
	} else {
		log.Errorf("%s was empty or not set, could not load store ids", keys.StoreIDs)
	}
}

func createManagementDatastore() (datastore.DataStore, error) {
	// Open the connection to the local database
	databaseDriver := viper.GetString(keys.DatabaseDriver)
	dataSource := viper.GetString(keys.DatabaseConnectionString)
	log.Infof("Connecting to %s with connection string %s", databaseDriver, dataSource)
	db, err := factory.CreateDataStore(databaseDriver, dataSource)
	if err != nil || db == nil {
		log.Fatalf("Error accessing data store: %v, with database source %s", databaseDriver, dataSource)
		return nil, err
	}

	return db, err
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

	certsDir := viper.GetString(keys.MQTTCertificatesPath)
	log.Tracef("MQTT Certs dir: %s", certsDir)

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
	return w
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

	sourceURL := viper.GetString(keys.GetIdentityKey(keys.MigrationsSourceURL))
	databaseName := viper.GetString(keys.GetIdentityKey(keys.DatabaseName))
	err = migrate2.Run(dataStoreSource, databaseDriver, fmt.Sprintf("file://%s", sourceURL), databaseName)
	if err != nil {
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

	enrollPort := viper.GetString(keys.GetIdentityKey(keys.ServicePortEnroll))

	log.Info("Starting service (enroll) on port : ", enrollPort)

	// internalRouter := gin.New()
	enrollRouter := gin.New()

	logFormat := os.Getenv("LOG_FORMAT")
	if strings.ToUpper(logFormat) == "JSON" {
		log.Infof("Setting up JSON log format for logger middleware")

		middlewareLogger := logger.New(log.StandardLogger(), logger.LogOptions{EnableStarting: true})

		enrollRouter.Use(middlewareLogger.HandleFunc)

	} else {
		enrollRouter.Use(gin.Logger())
	}

	wb.SetRouter(enrollRouter)

	go func() {
		log.Info("Listening and serving enroll on :" + enrollPort)

		err = enrollRouter.Run(":" + enrollPort)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return wb
}
