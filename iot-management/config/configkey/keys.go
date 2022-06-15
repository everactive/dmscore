// Package configkey provides the constant keys for getting a configuration value
package configkey

const (
	// AuthProvider is the provider of authentication for external clients (keycloak, static-client)
	AuthProvider = "service.auth.provider"
	// ClientTokenProvider is the provider for access tokens for external APIs
	ClientTokenProvider = "service.client.token.provider"
	// DatabaseDriver is the driver (database type) to use (Postgres, memory, etc.)
	DatabaseDriver = "database.driver"
	// DatabaseConnectionString is the data source string for the given database
	DatabaseConnectionString = "database.connection.string"
	// DeviceTwinAPIURL is the URL to the DeviceTwin service REST API
	DeviceTwinAPIURL = "devicetwin.api.url"
	// IdentityAPIURL is the URL to the Identity service REST API
	IdentityAPIURL = "identity.api.url"
	// JwtSecret is the JWT secret value
	JwtSecret = "service.jwtsecret" //nolint
	// OAuth2ClientID is the client id for the client credential grant
	OAuth2ClientID = "service.oauth2.client.id"
	// OAuth2ClientSecret is the client secret for the client credential grant
	OAuth2ClientSecret = "service.oauth2.client.secret" //nolint:gosec
	// OAuth2AccessTokenPath is the path part of the URL for getting a token verified/decoded
	OAuth2AccessTokenPath = "service.oauth2.token.access.path" //nolint:gosec
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
	// ServiceHost is the service host name
	ServiceHost = "service.host"
	// ServicePort is the service port to use
	ServicePort = "service.port"
	// ServiceScheme is the service scheme to use
	ServiceScheme = "service.scheme"
	// ServiceVersion is the service version
	ServiceVersion = "service.version"
	// StaticClientToken is the static client token used to authenticate with the other services
	StaticClientToken = "static.client.token"
	// StoreURL is the URL to Canonicals Snap Store API
	StoreURL = "store.url"
	// StoreIDs is a JSON serialized array of modelStoreIds
	StoreIDs = "store.ids"
	// ComponentVersionsCacheDuration is a time.Duration of the time between when component version information is read from disk
	ComponentVersionsCacheDuration = "service.component.versions.cache.duration"
	// FrontEndScheme is the protocol scheme for the frontend
	FrontEndScheme = "frontend.scheme"
	// FrontEndHost is the host for the frontend (ex. localhost:80, localhost:8080, etc.)
	FrontEndHost = "frontend.host"
)
