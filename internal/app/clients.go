package app

import (
	"context"

	co "github.com/metal-toolbox/conditionorc/pkg/api/v1/client"
	"github.com/metal-toolbox/mctl/internal/auth"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/pkg/errors"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	ErrNoTokenInRing = errors.New("secret not found in keyring")
	ErrAuth          = errors.New("authentication error")
	ErrNilConfig     = errors.New("no configuration defined")
)

func NewServerserviceClient(ctx context.Context, cfg *model.ConfigOIDC, reauth bool) (*serverservice.Client, error) {
	accessToken := "fake"

	if cfg == nil {
		return nil, errors.Wrap(ErrNilConfig, "missing serverservice API client configuration")
	}

	if cfg.Disable {
		return serverservice.NewClientWithToken(
			accessToken,
			cfg.Endpoint,
			nil,
		)
	}

	token, err := auth.AccessToken(ctx, model.ServerserviceAPI, cfg, reauth)
	if err != nil {
		return nil, errors.Wrap(ErrAuth, string(model.ServerserviceAPI)+err.Error())
	}

	return serverservice.NewClientWithToken(
		token,
		cfg.Endpoint,
		nil,
	)
}

func NewConditionsClient(ctx context.Context, cfg *model.ConfigOIDC, reauth bool) (*co.Client, error) {
	if cfg == nil {
		return nil, errors.Wrap(ErrNilConfig, "missing conditions API client configuration")
	}

	if cfg.Disable {
		return co.NewClient(
			cfg.Endpoint,
		)
	}

	token, err := auth.AccessToken(ctx, model.ConditionsAPI, cfg, reauth)
	if err != nil {
		return nil, errors.Wrap(ErrAuth, string(model.ConditionsAPI)+err.Error())
	}

	return co.NewClient(
		cfg.Endpoint,
		co.WithAuthToken(token),
	)
}
