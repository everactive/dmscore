package migrate

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	log "github.com/sirupsen/logrus"
)

func RunWithDB(db *sql.DB, sourceURL, databaseName string) error {
	driverInst, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		databaseName,
		driverInst)

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

func Run(datasource, driver, sourceURL, databaseName string) error {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		log.Fatal(err)
	}

	return RunWithDB(db, sourceURL, databaseName)
}
