// Package keys contains the constant values for the configuration keys
package keys

const (
	// AuthProvider is the string name of the provider for the factory to use (keycloak)
	AuthProvider = "service.auth.provider"

	// ConfigPath is the path used to load (if provided) and store configuration for the service
	ConfigPath = "service.config.path"
	// DatabaseDriver is the well known name for a given database type (memory or postgres)
	DatabaseDriver = "database.driver"
	// DatastoreSource is the connection string required for the database driver chosen
	DatastoreSource = "database.source"
	// MQTTClientCertificateFilename is the filename of the client certificate to use for MQTT
	MQTTClientCertificateFilename = "mqtt.client.certificate.filename"
	// MQTTClientKeyFilename is the filename of the client key to use for MQTT
	MQTTClientKeyFilename = "mqtt.client.key.filename"
	// MQTTRootCAFilename is the filename of the root certificate authority certificate to use
	MQTTRootCAFilename = "mqtt.root.ca.filename"
	// MQTTClientIDPrefix is the prefix to use for the MQTT client
	MQTTClientIDPrefix = "mqtt.client.id.prefix"
	// MQTTHealthTopic is the health topic to use to get health messages from devices
	MQTTHealthTopic = "mqtt.topic.health"
	// MQTTPubTopic is publish topic to use for sending devices actions
	MQTTPubTopic = "mqtt.topic.pub"
	// MQTTURL is the URL for the MQTT broker without scheme or port (i.e. broker.example.com)
	MQTTURL = "mqtt.url"
	// MQTTPort is the port for the MQTT broker
	MQTTPort = "mqtt.port"
	// ServicePort is port to run the service (REST API) on
	ServicePort = "service.port"
)
