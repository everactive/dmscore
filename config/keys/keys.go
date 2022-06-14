package keys

import "fmt"

const (
	DeviceTwinKeyPrefix = "devicetwin"
	IdentityKeyPrefix   = "identity"

	// AuthProvider is the string name of the provider for the factory to use (keycloak)
	AuthProvider = "service.auth.provider"
	// DisableAuth is a toggle for checking for auth, if not it will use static-client for every request (auth provider must also be "disabled")
	DisableAuth = "service.auth.disabled"
	// DatabaseDriver is the driver (database type) to use (Postgres, memory, etc.)
	DatabaseDriver = "database.driver"
	// DatabaseConnectionString is the data source string for the given database
	DatabaseConnectionString = "database.connection.string"
	// DeviceTwinAPIURL is the URL to the DeviceTwin service REST API
	DeviceTwinAPIURL = "devicetwin.api.url"
	// IdentityAPIURL is the URL to the Identity service REST API
	IdentityAPIURL = "identity.api.url"
	// MQTTURL is the URL for the MQTT broker without scheme or port (i.e. broker.example.com)
	MQTTURL = "mqtt.url"
	// MQTTPort is the port for the MQTT broker
	MQTTPort = "mqtt.port"
	// CertificatesPath is the file path to the certificates used for MQTT communication
	CertificatesPath = "mqtt.certificates.path"
	// MQTTClientCertificateFilename is the filename of the client certificate to use for MQTT
	MQTTClientCertificateFilename = "mqtt.client.certificate.filename"
	// MQTTClientKeyFilename is the filename of the client key to use for MQTT
	MQTTClientKeyFilename = "mqtt.client.key.filename"
	// MQTTRootCAFilename is the filename of the root certificate authority certificate to use
	MQTTRootCAFilename = "mqtt.root.ca.filename"
	// MQTTClientIDPrefix is the prefix to use for the MQTT client
	MQTTClientIDPrefix = "mqtt.client.id.prefix"
	// ServicePort is port to run the service (REST API) on
	ServicePort = "service.port"
	// MQTTHealthTopic is the health topic to use to get health messages from devices
	MQTTHealthTopic = "mqtt.topic.health"
	// MQTTPubTopic is publish topic to use for sending devices actions
	MQTTPubTopic = "mqtt.topic.pub"
	// ServicePortInternal is the port for the internal/private only part of the API
	ServicePortInternal = "service.port.internal"
	// ServicePortEnroll is the port for the HTTP service that is exposed externally for clients to use for enrolling
	ServicePortEnroll = "service.port.enroll"
	// StaticClientToken is the static client token used to authenticate with the other services
	StaticClientToken = "static.client.token"
)

func GetIdentityKey(key string) string {
	return getServiceKey(IdentityKeyPrefix, key)
}

func GetDeviceTwinKey(key string) string {
	return getServiceKey(DeviceTwinKeyPrefix, key)
}

func getServiceKey(prefix, key string) string {
	return fmt.Sprintf("%s.%s", prefix, key)
}
