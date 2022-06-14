package models

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	OrgId       string `gorm:"unique"`
	Name        string
	CountryName string
	RootCert    string
	RootKey     string
}

func (Organization) TableName() string {
	return "organization"
}
