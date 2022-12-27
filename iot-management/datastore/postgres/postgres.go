// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Management Service
 * Copyright 2019 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU Affero General Public License version 3, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY,
 * SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.
 * See the GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package postgres

import (
	"database/sql"
	"fmt"
	"github.com/everactive/dmscore/config/keys"
	"github.com/everactive/dmscore/pkg/migrate"
	"github.com/spf13/viper"
	"log"

	_ "github.com/lib/pq" // postgresql driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Store implements an in-memory store for testing
type Store struct {
	*sql.DB
	gormDB *gorm.DB
}

var pgStore *Store

// OpenStore returns an open database connection
func OpenStore(driver, dataSource string) *Store {
	if pgStore != nil {
		return pgStore
	}

	// Open the database
	pgStore = openDatabase(driver, dataSource)

	return pgStore
}

func OpenStoreWithDB(db *gorm.DB) *Store {
	// Check that we have a valid database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error accessing the database: %v\n", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Error accessing the database: %v\n", err)
	}

	databaseName := viper.GetString(keys.DatabaseName)
	migrationsPath := viper.GetString(keys.MigrationsSourceURL)
	err = migrate.RunWithDB(sqlDB, fmt.Sprintf("file://%s", migrationsPath), databaseName)
	if err != nil {
		log.Fatalf("Error during migrations, need to manually intervene: %s", err.Error())
	}

	return &Store{
		DB:     sqlDB,
		gormDB: db,
	}
}

// openDatabase return an open database connection for a postgreSQL database
func openDatabase(driver, dataSource string) *Store {
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

	databaseName := viper.GetString(keys.DatabaseName)
	migrationsPath := viper.GetString(keys.MigrationsSourceURL)
	err = migrate.Run(dataSource, driver, fmt.Sprintf("file://%s", migrationsPath), databaseName)
	if err != nil {
		log.Fatalf("Error during migrations, need to manually intervene: %s", err.Error())
	}

	return &Store{
		DB:     sqlDB,
		gormDB: db,
	}
}
