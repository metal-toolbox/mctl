package app

import (
	"context"
	"fmt"
	"strings"

	co "github.com/metal-toolbox/conditionorc/pkg/api/v1/client"
	"github.com/pkg/errors"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	ErrNoTokenInRing = errors.New("secret not found in keyring")
	ErrAuth          = errors.New("authentication error")
)

func getAuthTokenFromKeyring(ctx context.Context, mctl *App) (string, error) {
	accessToken := "fake"

	if !mctl.Config.DisableOAuth {
		token, err := mctl.RefreshToken(
			ctx,
			mctl.Config.OidcClientID,
			mctl.Config.OidcIssuerEndpoint,
		)
		if err != nil {
			if strings.Contains(err.Error(), "secret not found in keyring") {
				return "", ErrNoTokenInRing
			}

			return "", fmt.Errorf("%w: %s", ErrAuth, err.Error())
		}
		accessToken = token.AccessToken
	}
	return accessToken, nil
}

func NewServerserviceClient(ctx context.Context, mctl *App) (*serverservice.Client, error) {
	accessToken, err := getAuthTokenFromKeyring(ctx, mctl)
	if err != nil {
		return nil, err
	}

	return serverservice.NewClientWithToken(accessToken, mctl.Config.ServerserviceEndpoint, nil)
}

func NewConditionsClient(ctx context.Context, mctl *App) (*co.Client, error) {
	accessToken, err := getAuthTokenFromKeyring(ctx, mctl)
	if err != nil {
		return nil, err
	}

	return co.NewClient(mctl.Config.ConditionsEndpoint,
		co.WithAuthToken(accessToken),
	)
}
