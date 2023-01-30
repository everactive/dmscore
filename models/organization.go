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

type OrganizationUser struct {
	ID       uint `gorm:"primarykey"`
	OrgID    string
	UserName string `gorm:"column:username"`
}

func (OrganizationUser) TableName() string {
	return "organization_user"
}
