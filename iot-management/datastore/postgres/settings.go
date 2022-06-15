package postgres

import (
	"github.com/everactive/dmscore/iot-management/datastore/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// GetSettings gets an array of models.Setting or an error
func (s Store) GetSettings() ([]models.Setting, error) {
	settings := []models.Setting{}
	db := s.gormDB.Find(&settings)
	if db.Error != nil {
		return []models.Setting{}, db.Error
	}

	return settings, nil
}

// Set sets the value of a given key in settings
func (s Store) Set(key string, value string) error {
	setting := models.Setting{
		Key:   key,
		Value: value,
	}
	log.Tracef("Setting %s to %s", key, value)
	tx := s.gormDB.Begin()
	tx.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&setting)

	if tx.Error != nil {
		log.Println(tx.Error)
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}
