package model

const (
	AttributeNSFirmwareSetLabels = "sh.hollow.firmware_set.labels"
)

type (
	APIKind string
)

const (
	FleetDBAPI     APIKind = "fleetdbapi"
	ConditionsAPI  APIKind = "conditions"
	BomsServiceAPI APIKind = "bomservice"
)

// Config struct holds the mctl configuration parameters
type Config struct {

	// File is configuration file path
	File       string
	FleetDBAPI *ConfigOIDC `mapstructure:"serverservice_api"` // TODO: implement backwards compatibility and rename.
	Conditions *ConfigOIDC `mapstructure:"conditions_api"`
	BomService *ConfigOIDC `mapstructure:"bomservice_api"`
}

type ConfigOIDC struct {
	// FleetDBAPI is the Hollow server inventory store,
	// https://github.com/metal-toolbox/fleetdb
	Endpoint string `mapstructure:"endpoint"`

	// Disable skips OAuth setup
	Disable bool `mapstructure:"disable"`

	// FleetDBAPI OAuth2 parameters
	ClientID         string   `mapstructure:"oidc_client_id"`
	IssuerEndpoint   string   `mapstructure:"oidc_issuer_endpoint"`
	AudienceEndpoint string   `mapstructure:"oidc_audience_endpoint"`
	Scopes           []string `mapstructure:"oidc_scopes"`
	PkceCallbackURL  string   `mapstructure:"oidc_pkce_callback_url"`
}
