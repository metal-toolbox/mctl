package model

const (
	AttributeNSFirmwareSetLabels = "sh.hollow.firmware_set.labels"
)

// Config struct holds the mctl configuration parameters
type Config struct {
	// File is configuration file path
	File string
	// ServerService is the Hollow server inventory store,
	// https://github.com/metal-toolbox/hollow-serverservice
	ServerserviceEndpoint string `mapstructure:"serverservice_endpoint"`

	// ConditionsEndpoint is the URL for the Condition Orchestrator API
	ConditionsEndpoint string `mapstructure:"conditions_endpoint"`

	// ServerService OAuth2 parameters
	OidcClientID       string `mapstructure:"oidc_client_id"`
	OidcIssuerEndpoint string `mapstructure:"oidc_issuer_endpoint"`
	OidcAudience       string `mapstructure:"oidc_audience"`

	// Disable Oauth
	DisableOAuth bool `mapstructure:"disable_oauth"`
}

// Firmware includes a firmware version attributes and is part of FirmwareConfig
type Firmware struct {
	Vendor        string   `yaml:"vendor"`
	Version       string   `yaml:"version"`
	UpstreamURL   string   `yaml:"upstreamURL"`
	RepositoryURL string   `yaml:"repositoryURL"`
	FileName      string   `yaml:"filename"`
	Utility       string   `yaml:"utility"`
	Component     string   `yaml:"component"`
	Checksum      string   `yaml:"checksum"`
	Model         []string `yaml:"model"`
}

// FirmwareConfig struct holds firmware configuration data
type FirmwareConfig struct {
	Firmwares []*Firmware `yaml:"firmwares"`
}
