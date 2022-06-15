package models

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	OrganizationID string `gorm:"column:code"`
	Name           string
}

func (Organization) TableName() string {
	return "organization"
}
