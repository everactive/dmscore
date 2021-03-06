// Code generated by schema-generate. DO NOT EDIT.

// Package messages contains the structs and types as defined by this schema.
package messages

import "time"

const (
	SchemaID2188300499 = "schemas.json"
)

// IoTDeviceTwinToAgentSchemaDefinitions
type IoTDeviceTwinToAgentSchemaDefinitions interface{}

// Device
type Device struct {
	Brand       string         `json:"brand,omitempty"`
	Created     time.Time      `json:"created,omitempty"`
	DeviceId    string         `json:"deviceId,omitempty"`
	DeviceKey   string         `json:"deviceKey,omitempty"`
	LastRefresh time.Time      `json:"lastRefresh,omitempty"`
	Model       string         `json:"model,omitempty"`
	OrgId       string         `json:"orgId,omitempty"`
	Serial      string         `json:"serial,omitempty"`
	Store       string         `json:"store,omitempty"`
	Version     *DeviceVersion `json:"version,omitempty"`
}

// DeviceSnap
type DeviceSnap struct {
	Channel       string           `json:"channel,omitempty"`
	Config        string           `json:"config,omitempty"`
	Confinement   string           `json:"confinement,omitempty"`
	DeviceId      string           `json:"deviceId,omitempty"`
	Devmode       bool             `json:"devmode,omitempty"`
	InstalledDate time.Time        `json:"installedDate,omitempty"`
	InstalledSize int64            `json:"installedSize,omitempty"`
	Name          string           `json:"name,omitempty"`
	Revision      int              `json:"revision,omitempty"`
	Services      []*ServiceStatus `json:"services,omitempty"`
	Status        string           `json:"status,omitempty"`
	Version       string           `json:"version,omitempty"`
}

// DeviceVersion
type DeviceVersion struct {
	DeviceId      string `json:"deviceId,omitempty"`
	KernelVersion string `json:"kernelVersion,omitempty"`
	OnClassic     bool   `json:"onClassic,omitempty"`
	OsId          string `json:"osId,omitempty"`
	OsVersionId   string `json:"osVersionId,omitempty"`
	Series        string `json:"series,omitempty"`
	Version       string `json:"version,omitempty"`
}

// Health
type Health struct {
	DeviceId string    `json:"deviceId,omitempty"`
	OrgId    string    `json:"orgId,omitempty"`
	Refresh  time.Time `json:"refresh,omitempty"`
}

// PublishDevice
type PublishDevice struct {
	Action  string  `json:"action,omitempty"`
	Id      string  `json:"id,omitempty"`
	Message string  `json:"message,omitempty"`
	Result  *Device `json:"result,omitempty"`
	Success bool    `json:"success,omitempty"`
}

// PublishDeviceVersion
type PublishDeviceVersion struct {
	Action  string         `json:"action,omitempty"`
	Id      string         `json:"id,omitempty"`
	Message string         `json:"message,omitempty"`
	Result  *DeviceVersion `json:"result,omitempty"`
	Success bool           `json:"success,omitempty"`
}

// PublishResponse
type PublishResponse struct {
	Action  string `json:"action,omitempty"`
	Id      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success,omitempty"`
}

// PublishSnap
type PublishSnap struct {
	Action  string      `json:"action,omitempty"`
	Id      string      `json:"id,omitempty"`
	Message string      `json:"message,omitempty"`
	Result  *DeviceSnap `json:"result,omitempty"`
	Success bool        `json:"success,omitempty"`
}

// PublishSnapTask
type PublishSnapTask struct {
	Action  string `json:"action,omitempty"`
	Id      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Result  string `json:"result,omitempty"`
	Success bool   `json:"success,omitempty"`
}

// PublishSnaps
type PublishSnaps struct {
	Action  string        `json:"action,omitempty"`
	Id      string        `json:"id,omitempty"`
	Message string        `json:"message,omitempty"`
	Result  []*DeviceSnap `json:"result,omitempty"`
	Success bool          `json:"success,omitempty"`
}

// ServiceStatus
type ServiceStatus struct {
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Enabled bool   `json:"enabled"`
	Daemon  string `json:"daemon"`
}

// SubscribeAction
type SubscribeAction struct {
	Action string `json:"action,omitempty"`
	Data   string `json:"data,omitempty"`
	Id     string `json:"id,omitempty"`
	Snap   string `json:"snap,omitempty"`
}
