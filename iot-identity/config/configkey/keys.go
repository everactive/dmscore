// Package configkey contains the constant string literal keys for configuration
package configkey

const (
	// AuthProvider is the string name of the provider for the factory to use (keycloak)
	AuthProvider = "service.auth.provider"
	// DatabaseDriver is the driver (database type) to use (Postgres, memory, etc.)
	DatabaseDriver = "database.driver"
	// DatabaseConnectionString is the data source string for the given database
	DatabaseConnectionString = "database.connection.string"
	// OAuth2ClientID is the client id for the client credential grant
	OAuth2ClientID = "service.oauth2.client.id"
	// OAuth2ClientSecret is the client secret for the client credential grant
	OAuth2ClientSecret = "service.oauth2.client.secret" //nolint:gosec
	// OAuth2TokenIntrospectPath is the path part of the URL for getting a token verified/decoded
	OAuth2TokenIntrospectPath = "service.oauth2.token.introspect.path" //nolint:gosec
	// OAuth2ClientRequiredScope is the scope the client is expected to have to use the Identity API
	OAuth2ClientRequiredScope = "service.oauth2.client.required.scope"
	// OAuth2HostName is the host name or IP of the auth provider
	OAuth2HostName = "service.oauth2.host.name"
	// OAuth2HostPort is the port for the auth provider (optional to override default for scheme)
	OAuth2HostPort = "service.oauth2.host.port"
	// OAuth2HostScheme the scheme for the auth provider (only https should be used)
	OAuth2HostScheme = "service.oauth2.host.scheme"
	// ServicePortInternal is the port for the internal/private only part of the API
	ServicePortInternal = "service.port.internal"
	// ServicePortEnroll is the port for the HTTP service that is exposed externally for clients to use for enrolling
	ServicePortEnroll = "service.port.enroll"
	// MQTTHostAddress is the IP or FQDN of the MQTT broker, i.e. mqtt.somedomain.com or 192.168.13.13
	MQTTHostAddress = "mqtt.host.address"
	// MQTTHostPort is the port of the MQTT broker
	MQTTHostPort = "mqtt.host.port"
	// MQTTCertificatePath is the path to the certificates the service should use when connecting to the broker
	MQTTCertificatePath = "mqtt.certificates.path"
)
