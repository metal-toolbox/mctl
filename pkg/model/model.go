package model

const (
	
)

type Config struct {
	// File is configuration file path
	File string

	// Disable Oauth
	DisableOAuth bool `mapstructure:"disable_oauth"`

	// ServerService is the Hollow server inventory store
	// https://github.com/metal-toolbox/hollow-serverservice
	ServerserviceEndpoint string `mapstructure:"serverservice_endpoint"`
	OidcClientID          string `mapstructure:"oidc_client_id"`
	OidcIssuerEndpoint    string `mapstructure:"oidc_issuer_endpoint"`
	OidcAudience          string `mapstructure:"oidc_audience"`
}
