package models

import (
	"gorm.io/gorm"
	"time"
)

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

type HealthHash struct {
	gorm.Model
	LastRefresh time.Time 
	OrgID string
	DeviceID string
	SnapListHash string
	InstalledSnapsHash string
}
