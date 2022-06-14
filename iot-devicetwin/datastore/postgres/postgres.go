// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Device Twin Service
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

// Package postgres is the Datastore implementation for Postgres
package postgres

import (
	"context"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"os"
	"time"

	"github.com/everactive/dmscore/iot-devicetwin/datastore"

	"gorm.io/driver/postgres"

	_ "github.com/lib/pq" // postgresql driver

	"gorm.io/gorm"
)

// DataStore is the postgreSQL implementation of a data store
type DataStore struct {
	driver string
	*sql.DB
	gormDB   *gorm.DB
	unscoped bool
}

var pgStore *DataStore

// OpenDataStore returns an open database connection
func OpenDataStore(driver, dataSource string) *DataStore {
	if pgStore != nil {
		return pgStore
	}

	// Open the database
	pgStore = openDatabase(driver, dataSource)

	return pgStore
}

type GormLogrusAdapter struct {
	logger *log.Logger
}

func (g GormLogrusAdapter) New(logger *log.Logger) *GormLogrusAdapter {
	return &GormLogrusAdapter{logger: logger}
}

func (g *GormLogrusAdapter) LogMode(level logger.LogLevel) logger.Interface {
	switch level {
	case logger.Silent:
		fallthrough
	case logger.Error:
		g.logger.SetLevel(log.ErrorLevel)
		break
	case logger.Info:
		g.logger.SetLevel(log.InfoLevel)
	case logger.Warn:
		g.logger.SetLevel(log.WarnLevel)
	}

	return g
}

func (g *GormLogrusAdapter) Info(ctx context.Context, s string, i ...interface{}) {
	g.logger.Info(s, i)
}

func (g *GormLogrusAdapter) Warn(ctx context.Context, s string, i ...interface{}) {
	g.logger.Warn(s, i)
}

func (g *GormLogrusAdapter) Error(ctx context.Context, s string, i ...interface{}) {
	g.logger.Error(s, i)
}

func (g *GormLogrusAdapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sqlString, rowsAffected := fc()
	g.logger.Tracef("SQL: %s, rows affected: %d", sqlString, rowsAffected)
}

// openDatabase return an open database connection for a postgreSQL database
func openDatabase(driver, dataSource string) *DataStore {
	// Open the database connection
	newLogger := log.New()
	newLogger.SetOutput(os.Stdout)
	newLogger.SetLevel(log.TraceLevel)
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{
		Logger: GormLogrusAdapter{}.New(newLogger),
	})

	if err != nil {
		log.Fatalf("Error opening the database: %v\n", err)
	}

	// db.Logger.LogMode(logger.Info)

	// Check that we have a valid database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error accessing the database: %v\n", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Error accessing the database: %v\n", err)
	}

	return &DataStore{
		driver: driver,
		DB:     sqlDB,
		gormDB: db,
	}
}

// Unscoped gets an unscoped instance of the datastore
func (db *DataStore) Unscoped() datastore.UnscopedDataStore {
	unscopedDS := *db

	unscopedDS.unscoped = true
	return &unscopedDS
}
