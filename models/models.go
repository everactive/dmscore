package models

import "gorm.io/gorm"

type DeviceModel struct {
	gorm.Model
	Name                     string
	DeviceModelRequiredSnaps []DeviceModelRequiredSnap
}

type DeviceModelRequiredSnap struct {
	gorm.Model
	DeviceModelID uint
	Name          string `gorm:"uniqueIndex"`
}
