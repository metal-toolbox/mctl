package model

// Config struct holds the mctl configuration parameters
type Config struct {
	// File is configuration file path
	File string
	// ServerService is the Hollow server inventory store,
	// https://github.com/metal-toolbox/hollow-serverservice
	ServerserviceEndpoint string `mapstructure:"serverservice_endpoint"`

	// ServerService OAuth2 parameters
	OidcClientID       string `mapstructure:"oidc_client_id"`
	OidcIssuerEndpoint string `mapstructure:"oidc_issuer_endpoint"`
	OidcAudience       string `mapstructure:"oidc_audience"`

	// Disable Oauth
	DisableOAuth bool `mapstructure:"disable_oauth"`
}

// Firmware includes a firmware version attributes and is part of FirmwareConfig
type Firmware struct {
	Version       string `yaml:"version"`
	UpstreamURL   string `yaml:"upstreamURL"`
	FileName      string `yaml:"filename"`
	Utility       string `yaml:"utility"`
	Model         string `yaml:"model"`
	ComponentSlug string `yaml:"componentslug"`
	Checksum      string `yaml:"checksum"`
}

// FirmwareProviders struct holds firmwares for firmware providers is part of FirmwareConfig
type FirmwareProviders struct {
	Vendor           string      `yaml:"vendor"`
	RepositoryURL    string      `yaml:"repositoryURL"`
	RepositoryRegion string      `yaml:"repositoryRegion"`
	Firmwares        []*Firmware `yaml:"firmwares"`
}

// FirmwareConfig struct holds firmware configuration data
type FirmwareConfig struct {
	Providers []*FirmwareProviders `yaml:"providers"`
}
