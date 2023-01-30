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
	identityweb "github.com/everactive/dmscore/iot-identity/web"
	"github.com/everactive/dmscore/iot-management/identityapi"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/everactive/dmscore/iot-management/twinapi"
	datastores2 "github.com/everactive/dmscore/pkg/datastores"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		dataStores, err := datastores2.New()
		if err != nil {
			// We are in a catastrophic failure state if we cannot create datastores,
			// the only option here is to exit, hope things resolve or require manual intervention
			return err
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

		ids := createIdentityService(dataStores.IdentityStore)

		dts, controller := devicetwin2.New(ids, dataStores)

		// Create the main services
		srv := manage.NewManagement(dataStores, twinAPI, identityAPI, controller, ids.Identity)

		supervisorSpec := suture.Spec{}
		sup := suture.New("dmscore", supervisorSpec)

		sup.Add(ids)
		sup.Add(dts)
		sup.Add(web2.New(srv, dataStores.GetDatabase()))

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

func createIdentityService(ds identitydatastore.DataStore) *identityweb.IdentityService {
	service.Logger = log.StandardLogger()
	srv := service.NewIdentityService(ds)

	// Start the web service
	identityweb.Logger = log.StandardLogger()
	wb := identityweb.NewIdentityService(srv, log.StandardLogger())

	return wb
}
