package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everactive/dmscore/config"
	"github.com/everactive/dmscore/config/keys"
	identitydatastore "github.com/everactive/dmscore/iot-identity/datastore"
	"github.com/everactive/dmscore/iot-identity/service"
	identityfactory "github.com/everactive/dmscore/iot-identity/service/factory"
	identityweb "github.com/everactive/dmscore/iot-identity/web"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/identityapi"
	"github.com/everactive/dmscore/iot-management/service/factory"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/everactive/dmscore/iot-management/twinapi"
	migrate2 "github.com/everactive/dmscore/pkg/migrate"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"

	// this is needed for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// this is needed for migrate
	_ "github.com/lib/pq"

	web2 "github.com/everactive/dmscore/pkg/web"
	"github.com/thejerf/suture/v4"

	devicetwin2 "github.com/everactive/dmscore/pkg/devicetwin"
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

		db, ds, err := createManagementDatastore()

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

		dts, controller := devicetwin2.New(ds)

		// Create the main services
		srv := manage.NewManagement(db, ds, twinAPI, identityAPI, controller, ids.Identity)

		supervisorSpec := suture.Spec{}
		sup := suture.New("dmscore", supervisorSpec)

		sup.Add(ids)
		sup.Add(dts)
		sup.Add(web2.New(srv, db))

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

func createManagementDatastore() (*gorm.DB, datastore.DataStore, error) {
	// Open the connection to the local database
	databaseDriver := viper.GetString(keys.DatabaseDriver)
	dataSource := viper.GetString(keys.DatabaseConnectionString)
	log.Infof("Connecting to %s with connection string %s", databaseDriver, dataSource)

	// Create gorm database to use with the datastore
	// Open the database connection
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening the database: %v\n", err)
	}

	// Check that we have a valid database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error accessing the database: %v\n", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Error accessing the database: %v\n", err)
	}

	ds, err := factory.CreateDataStoreWithDB(db, databaseDriver)
	if err != nil || db == nil {
		log.Fatalf("Error accessing data store: %v, with database source %s", databaseDriver, dataSource)
		return nil, nil, err
	}

	return db, ds, err
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

	return wb
}
