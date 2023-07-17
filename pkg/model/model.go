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
	Serverservice *ConfigOIDC     `mapstructure:"serverservice_api"`
	Conditions    *ConfigOIDC     `mapstructure:"conditions_api"`
	Splunk        *ConfigLogIndex `mapstructure:"splunk"`
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

type ConfigLogIndex struct {
	Endpoint string `mapstructure:"endpoint"`
	Token    string `mapstructure:"token"`
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
