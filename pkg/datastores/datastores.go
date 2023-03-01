package datastores

import (
	"errors"
	"fmt"
	"github.com/everactive/dmscore/config/keys"
	devicetwindatastore "github.com/everactive/dmscore/iot-devicetwin/datastore"
	devicetwinfactory "github.com/everactive/dmscore/iot-devicetwin/service/factory"
	identitydatastore "github.com/everactive/dmscore/iot-identity/datastore"
	identityfactory "github.com/everactive/dmscore/iot-identity/service/factory"
	managementdatastore "github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/service/factory"
	"github.com/everactive/dmscore/models"
	migrate2 "github.com/everactive/dmscore/pkg/migrate"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataStores struct {
	IdentityStore   identitydatastore.DataStore
	DeviceTwinStore devicetwindatastore.DataStore
	ManagementStore managementdatastore.DataStore
	DataStore       DataStore
	db              *gorm.DB
}

func (dss *DataStores) GetDatabase() *gorm.DB {
	return dss.db
}

func New() (*DataStores, error) {
	ids, err := createIdentityDataStore()
	if err != nil {
		return nil, err
	}
	mds, err := createManagementDatastore()
	if err != nil {
		return nil, err
	}

	dtds, err := createDeviceTwinDataStore()
	if err != nil {
		return nil, err
	}

	// Open the connection to the local database
	databaseDriver := viper.GetString(keys.DatabaseDriver)
	dataSource := viper.GetString(keys.DatabaseConnectionString)
	log.Infof("Connecting to %s with connection string %s", databaseDriver, dataSource)

	// Create gorm database to use with the datastore
	// Open the database connection
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	cds := &ConsolidatedDataStore{db: db}

	return &DataStores{
		IdentityStore:   ids,
		DeviceTwinStore: dtds,
		ManagementStore: mds,
		DataStore:       cds,
		db:              db,
	}, nil
}

type DataStore interface {
	CheckAccess(orgID, username string, role int) error
	GetModelRequiredSnaps(modelName string) (*models.DeviceModel, error)
}

type ConsolidatedDataStore struct {
	db *gorm.DB
}

var (
	ErrorFindingOrgUser = errors.New("error finding organization user")
)

func (c *ConsolidatedDataStore) CheckAccess(orgID, username string, role int) error {
	// Superusers can access all accounts
	if role == managementdatastore.Superuser {
		return nil
	}

	var orgUser models.OrganizationUser
	tx := c.db.Find(&orgUser, &models.OrganizationUser{UserName: username, OrgID: orgID})
	if tx.Error != nil {
		return tx.Error
	}

	if orgUser.OrgID == orgID && orgUser.UserName == username {
		return nil
	}

	return ErrorFindingOrgUser
}

func (c *ConsolidatedDataStore) GetModelRequiredSnaps(modelName string) (*models.DeviceModel, error) {
	var deviceModel models.DeviceModel
	tx := c.db.Preload("DeviceModelRequiredSnaps").Find(&deviceModel, &models.DeviceModel{Name: modelName})
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &deviceModel, nil
}

func createIdentityDataStore() (identitydatastore.DataStore, error) {
	// Open the connection to the local database
	databaseDriver := viper.GetString(keys.GetIdentityKey(keys.DatabaseDriver))
	dataStoreSource := viper.GetString(keys.GetIdentityKey(keys.DatabaseConnectionString))

	log.Infof("Connecting to %s with connection string %s", databaseDriver, dataStoreSource)
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

func createManagementDatastore() (managementdatastore.DataStore, error) {
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
		return nil, err
	}

	return ds, err
}

func createDeviceTwinDataStore() (devicetwindatastore.DataStore, error) {
	databaseDriver := viper.GetString(keys.GetDeviceTwinKey(keys.DatabaseDriver))
	dataStoreSource := viper.GetString(keys.GetDeviceTwinKey(keys.DatabaseConnectionString))

	dtds, err := devicetwinfactory.CreateDataStore(databaseDriver, dataStoreSource)
	if err != nil {
		return nil, err
	}

	return dtds, nil
}
