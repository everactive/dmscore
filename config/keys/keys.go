package keys

import "fmt"

const (
	// ModelKeyTemplate is the template for getting a model key based on the model string
	ModelKeyTemplate = "store.model.%s"
)

const (
	DeviceTwinKeyPrefix = "devicetwin"
	IdentityKeyPrefix   = "identity"

	// AuthProvider is the string name of the provider for the factory to use (keycloak)
	AuthProvider = "service.auth.provider"
	// CertificatesPath is the path to the Identity component certificates
	CertificatesPath = "certificates.path"
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
	// MQTTCertificatesPath is the file path to the certificates used for MQTT communication
	MQTTCertificatesPath = "mqtt.certificates.path"
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
	// MigrationsSourceURL is the source URL of the migrations for each service sub-component
	MigrationsSourceURL = "service.migrations.source"
	// DatabaseName the name of the database for each service sub-component
	DatabaseName = "database.name"
	// ServiceScheme is the service scheme to use
	ServiceScheme = "service.scheme"
	// ServiceHost is the service host name
	ServiceHost = "service.host"
	// AutoRegistrationEnabled determines whether devices will be registered when they try to enroll
	AutoRegistrationEnabled = "identity.auto.registration.enabled"
	// DefaultOrganization is default organization used by auto registration
	DefaultOrganization = "identity.default.organization"
	// ValidSHA384Keys is and array of the SHA384 of public keys that are acceptable to have signed model and serial
	// assertions for auto-registration during enrollment
	ValidSHA384Keys = "identity.assertions.valid.key.signatures"
	// StoreURL is the URL to Canonicals Snap Store API
	StoreURL = "store.url"
	// StoreIDs is a JSON serialized array of modelStoreIds
	StoreIDs = "store.ids"
	// ClientTokenProvider is the provider for access tokens for external APIs
	ClientTokenProvider = "service.client.token.provider"
	// FrontEndScheme is the protocol scheme for the frontend
	FrontEndScheme = "frontend.scheme"
	// FrontEndHost is the host for the frontend (ex. localhost:80, localhost:8080, etc.)
	FrontEndHost = "frontend.host"
	// ComponentVersionsCacheDuration is a time.Duration of the time between when component version information is read from disk
	ComponentVersionsCacheDuration = "service.component.versions.cache.duration"
	// JwtSecret is the JWT secret value
	JwtSecret = "service.jwtsecret" //nolint
	// OAuth2ClientID is the client id for the client credential grant
	OAuth2ClientID = "service.oauth2.client.id"
	// OAuth2ClientSecret is the client secret for the client credential grant
	OAuth2ClientSecret = "service.oauth2.client.secret" //nolint:gosec
	// OAuth2AccessTokenPath is the path part of the URL for getting a token verified/decoded
	OAuth2AccessTokenPath = "service.oauth2.token.access.path" //nolint:gosec
	// OAuth2HostName is the host name or IP of the auth provider
	OAuth2HostName = "service.oauth2.host.name"
	// OAuth2HostPort is the port for the auth provider (optional to override default for scheme)
	OAuth2HostPort = "service.oauth2.host.port"
	// OAuth2HostScheme the scheme for the auth provider (only https should be used)
	OAuth2HostScheme = "service.oauth2.host.scheme"
	// DefaultServiceHeartbeat is the default duration between service heartbeats in the logs (info)
	DefaultServiceHeartbeat = "service.default.heartbeat"
	// RequiredSnapsInstallServiceCheckInterval is the interval in which the install service will start or
	// refresh the checker service which in turn iterates over all the devices to see if they are missing snaps
	RequiredSnapsInstallServiceCheckInterval = "service.install.interval"
	// RefreshSnapListOnAnyChange controls whether a snap list is requested if any changes in snaps are detected
	RefreshSnapListOnAnyChange = "service.refresh-snaps-any-change"
	// RequiredSnapsCheckInterval is the time, if the required snaps checker is running, between each device that it
	// checks to see if it needs snaps that are required (it doesn't have them installed),
	// 10 devices checked = 1s if this value is 100ms
	RequiredSnapsCheckInterval = "service.install-required-snaps-check.interval"
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
