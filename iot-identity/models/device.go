package models

import (
	"gorm.io/gorm"
)

// Status is a top-level enrollment status classification
type Status int

// Enrollment status classifications
const (
	StatusWaiting Status = iota + 1
	StatusEnrolled
	StatusDisabled
)

type RegisteredDevice struct {
	gorm.Model
	OrgID        string
	Brand        string
	DeviceModel  string `gorm:"column:model"`
	SerialNumber string
	PrivateKey   []byte `gorm:"column:cred_key"`
	Certificate  []byte `gorm:"column:cred_cert"`
	MQTTURL      string `gorm:"column:cred_mqtt"`
	MQTTPort     string `gorm:"column:cred_port"`
	StoreID      string
	DeviceKey    string
	Status       Status
	DeviceData   string
	DeviceID     string
}

func (rd RegisteredDevice) TableName() string {
	return "device"
}
