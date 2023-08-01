package model

const (
	AttributeNSFirmwareSetLabels = "sh.hollow.firmware_set.labels"
)

type (
	APIKind string
)

const (
	ServerserviceAPI APIKind = "serverservice"
	ConditionsAPI    APIKind = "conditions"
)

// Config struct holds the mctl configuration parameters
type Config struct {

	// File is configuration file path
	File          string
	Serverservice *ConfigOIDC `mapstructure:"serverservice_api"`
	Conditions    *ConfigOIDC `mapstructure:"conditions_api"`
}

type ConfigOIDC struct {
	// ServerService is the Hollow server inventory store,
	// https://github.com/metal-toolbox/hollow-serverservice
	Endpoint string `mapstructure:"endpoint"`

	// Disable skips OAuth setup
	Disable bool `mapstructure:"disable"`

	// ServerService OAuth2 parameters
	ClientID         string   `mapstructure:"oidc_client_id"`
	IssuerEndpoint   string   `mapstructure:"oidc_issuer_endpoint"`
	AudienceEndpoint string   `mapstructure:"oidc_audience_endpoint"`
	Scopes           []string `mapstructure:"oidc_scopes"`
	PkceCallbackURL  string   `mapstructure:"oidc_pkce_callback_url"`
}
