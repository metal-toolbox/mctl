package app

import (
	"context"
	"net/url"
	"os"

	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

const (
	EnvDisableAuth = "DISABLE_AUTH"
)

var (
	ErrConfig = errors.New("configuration error")
)

// Config holds configuration data when running mctl
// App holds attributes for the mtl application
type App struct {
	Config *model.Config
	Client *serverservice.Client
}

func New(ctx context.Context, cfgFile string) (app *App, err error) {
	cfg, err := loadConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	_, err = validateServerServiceParams(cfg)
	if err != nil {
		return nil, err
	}

	return &App{Config: cfg}, nil
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

// validateServerServiceParams checks required serverservice configuration parameters are present
// and returns the serverservice URL endpoint
func validateServerServiceParams(cfg *model.Config) (*url.URL, error) {
	if cfg.ServerserviceEndpoint == "" {
		return nil, errors.Wrap(ErrConfig, "Serverservice endpoint not defined")
	}

	endpoint, err := url.Parse(cfg.ServerserviceEndpoint)
	if err != nil {
		return nil, errors.Wrap(ErrConfig, "Serverservice endpoint URL error: "+err.Error())
	}

	if cfg.DisableOAuth {
		return endpoint, nil
	}

	if cfg.OidcIssuerEndpoint == "" {
		return nil, errors.Wrap(ErrConfig, "OIDC issuer endpoint not defined")
	}

	if cfg.OidcAudience == "" {
		return nil, errors.Wrap(ErrConfig, "OIDC Audience not defined")
	}

	return endpoint, nil
}
