package app

import (
	"context"
	"net/url"
	"os"

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

func loadConfig(cfgFile string) (*model.Config, error) {
	cfg := &model.Config{}

	if cfgFile != "" {
		cfg.File = cfgFile
	} else {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		cfg.File = homedir + "/" + ".mctl.yml"
	}

	viper.AutomaticEnv()
	if viper.GetString("mctlconfig") != "" {
		cfg.File = viper.GetString("mctlconfig")
	}

	h, err := os.Open(cfg.File)
	if err != nil {
		return nil, err
	}

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
	if cfg.Serverservice != nil {
		if err := validateConfigOIDC(cfg.Serverservice); err != nil {
			return errors.Wrap(err, "serverservice API config")
		}
	}

	if cfg.Conditions != nil {
		err := validateConfigOIDC(cfg.Conditions)
		if err != nil {
			return errors.Wrap(err, "conditions API config")
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
