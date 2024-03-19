package app

import (
	"context"
	"net/url"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	ErrConfig = errors.New("configuration error")
)

// Config holds configuration data when running mctl
// App holds attributes for the mtl application
type App struct {
	Config *model.Config
	// Force client to re-authenticate with Oauth services.
	Reauth bool
}

func New(_ context.Context, cfgFile string, reauth bool) (app *App, err error) {
	cfg, err := loadConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	err = validateClientParams(cfg)
	if err != nil {
		return nil, err
	}

	return &App{Config: cfg, Reauth: reauth}, nil
}

func openConfig(path string) (*os.File, error) {
	if path != "" {
		return os.Open(path)
	}
	path = viper.GetString("mctlconfig")
	if path != "" {
		return os.Open(path)
	}

	path = filepath.Join(xdg.Home, ".mctl.yml")
	f, err := os.Open(path)
	if err == nil {
		return f, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	path, err = xdg.ConfigFile("mctl/config.yaml")
	if err != nil {
		return nil, err
	}

	return os.Open(path)
}

func loadConfig(cfgFile string) (*model.Config, error) {
	cfg := &model.Config{}
	viper.AutomaticEnv()
	h, err := openConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	cfg.File = h.Name()
	viper.SetConfigFile(cfg.File)

	err = viper.ReadConfig(h)
	if err != nil {
		return nil, errors.Wrap(err, cfg.File)
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// validateClientParams checks required downstream service configuration parameters are present
func validateClientParams(cfg *model.Config) error {
	if cfg.FleetDBAPI != nil {
		if err := validateConfigOIDC(cfg.FleetDBAPI); err != nil {
			return errors.Wrap(err, "fleetdb API API config")
		}
	}

	if cfg.Conditions != nil {
		err := validateConfigOIDC(cfg.Conditions)
		if err != nil {
			return errors.Wrap(err, "conditions API config")
		}
	}

	if cfg.BomService != nil {
		err := validateConfigOIDC(cfg.BomService)
		if err != nil {
			return errors.Wrap(err, "bomservice API config")
		}
	}

	return nil
}

func validateConfigOIDC(cfg *model.ConfigOIDC) error {
	errConfigOIDC := errors.New("OIDC configuration error")

	if cfg == nil {
		return errors.Wrap(errConfigOIDC, "configuration not defined")
	}

	if cfg.Endpoint == "" {
		return errors.Wrap(errConfigOIDC, "endpoint not defined")
	}

	_, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return errors.Wrap(errConfigOIDC, "endpoint URL error: "+err.Error())
	}

	if cfg.Disable {
		return nil
	}

	if cfg.IssuerEndpoint == "" {
		return errors.Wrap(errConfigOIDC, "Issuer endpoint not defined")
	}

	if cfg.AudienceEndpoint == "" {
		return errors.Wrap(errConfigOIDC, "Audience endpoint not defined")
	}

	return nil
}
