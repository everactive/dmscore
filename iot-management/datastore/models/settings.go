// Package models provides the models for database types with GORM
package models

import "gorm.io/gorm"

// Setting is a model of a setting used by GORM to store them in a database
type Setting struct {
	gorm.Model
	Key   string
	Value string
}
