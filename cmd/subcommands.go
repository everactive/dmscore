package cmd

import (
	"fmt"
	"github.com/everactive/dmscore/config"
	"github.com/everactive/dmscore/config/keys"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/service/factory"
	"github.com/everactive/dmscore/versions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var version = cobra.Command{
	Use:   "version",
	Short: "version",
	Long:  "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versions.Version)
	},
}
var createSuperuser = cobra.Command{
	Use:   "create-superuser",
	Short: "create-superuser",
	Long:  "create-superuser",
	Run: func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Error(err)
			return
		}

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Error(err)
			return
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Error(err)
			return
		}

		if !cmd.Flags().Changed("username") || !cmd.Flags().Changed("email") || !cmd.Flags().Changed("name") {
			log.Errorf("must set username, email and name")
			return
		}

		var configFilePath string
		if filePath, ok := os.LookupEnv("CONFIG_FILE_PATH"); ok {
			configFilePath = filePath
		}

		config.LoadConfig(configFilePath)

		// Open the connection to the local database
		databaseDriver := viper.GetString(keys.DatabaseDriver)
		dataSource := viper.GetString(keys.DatabaseConnectionString)
		db, err := factory.CreateDataStore(databaseDriver, dataSource)
		if err != nil || db == nil {
			log.Fatalf("Error accessing data store: %v", databaseDriver)
			return
		}

		// Once we have a database connection, we need to load any settings that were persisted
		// in the database and make those available via Viper
		settings, err := db.GetSettings()
		if err != nil {
			log.Fatalf("Cannot continue without loading settings")
			return
		}

		for _, setting := range settings {
			viper.Set(setting.Key, setting.Value)
		}

		// Create the user
		if len(username) == 0 {
			log.Errorf("the username must be supplied")
			return
		}

		// Create the user
		user := datastore.User{
			Username: username,
			Name:     name,
			Email:    email,
			Role:     datastore.Superuser,
		}
		_, err = db.CreateUser(user)
		if err != nil {
			log.Error(err)
			return
		}
	},
}
