package app

import (
	"context"

	bomclient "github.com/metal-toolbox/bomservice/pkg/api/v1/client"
	co "github.com/metal-toolbox/conditionorc/pkg/api/v1/conditions/client"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/metal-toolbox/mctl/internal/auth"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/pkg/errors"
)

var (
	ErrNoTokenInRing = errors.New("secret not found in keyring")
	ErrAuth          = errors.New("authentication error")
	ErrNilConfig     = errors.New("no configuration defined")
)

func NewFleetDBAPIClient(ctx context.Context, cfg *model.ConfigOIDC, reauth bool) (*fleetdbapi.Client, error) {
	accessToken := "fake"

	if cfg == nil {
		return nil, errors.Wrap(ErrNilConfig, "missing fleetdb API API client configuration")
	}

	if cfg.Disable {
		return fleetdbapi.NewClientWithToken(
			accessToken,
			cfg.Endpoint,
			nil,
		)
	}

	token, err := auth.AccessToken(ctx, model.FleetDBAPI, cfg, reauth)
	if err != nil {
		return nil, errors.Wrap(ErrAuth, string(model.FleetDBAPI)+err.Error())
	}

	return fleetdbapi.NewClientWithToken(
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

func NewBomServiceClient(ctx context.Context, cfg *model.ConfigOIDC, reauth bool) (*bomclient.Client, error) {
	if cfg == nil {
		return nil, errors.Wrap(ErrNilConfig, "missing bom service API client configuration")
	}

	if cfg.Disable {
		return bomclient.NewClient(
			cfg.Endpoint,
		)
	}

	token, err := auth.AccessToken(ctx, model.BomsServiceAPI, cfg, reauth)
	if err != nil {
		return nil, errors.Wrap(ErrAuth, string(model.BomsServiceAPI)+err.Error())
	}

	return bomclient.NewClient(
		cfg.Endpoint,
		bomclient.WithAuthToken(token),
	)
}
