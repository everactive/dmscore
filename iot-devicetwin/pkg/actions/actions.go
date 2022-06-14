// Package actions contains the valid actions
package actions

const (
	// Ack is the action to add an assertion to the device
	Ack = "ack"
	// Conf is the action to get a snap's configuration
	Conf = "conf"
	// Device is the action to get a device's info
	Device = "device"
	// Disable is the action for disabling a snap
	Disable = "disable"
	// Enable is the action for enabling a snap
	Enable = "enable"
	// Info is the action to get info about a snap
	Info = "info"
	// Install is the action for installing a snap
	Install = "install"
	// List is the action for getting a list of snaps
	List = "list"
	// Logs is the action for getting logs from a device from its snapd endpoint
	Logs = "logs"
	// Refresh is the action for refreshing a snap
	Refresh = "refresh"
	// Remove is the action for removing a snap
	Remove = "remove"
	// Revert is the action for reverting a snap
	Revert = "revert"
	// Restart is the action for restarting a snap or snap service
	Restart = "restart"
	// Server is the action to get details of the device version
	Server = "server"
	// SetConf is the action for setting snap configuration
	SetConf = "setconf"
	// Snapshot is the action for creating and sending a snap snapshot to an S3 storage service
	Snapshot = "snapshot"
	// Start is the action for starting a snap or snap service
	Start = "start"
	// Stop is the action for stopping a snap or snap service
	Stop = "stop"
	// Switch is the action to switch a snap to another channel
	Switch = "switch"
	// Unregister is the action for unregistering a device from the service
	Unregister = "unregister"
	// User is the action for adding/removing a user from the device
	User = "user"
)
