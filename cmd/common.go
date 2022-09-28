package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/metal-toolbox/mctl/internal/app"
	"golang.org/x/net/context"

	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

func newServerserviceClient(ctx context.Context, mctl *app.App) (*serverservice.Client, error) {
	token, err := mctl.RefreshToken(
		ctx,
		mctl.Config.OidcClientID,
		mctl.Config.OidcIssuerEndpoint,
	)
	if err != nil {
		if strings.Contains(err.Error(), "secret not found in keyring") {
			log.Println("please run `mctl auth` and try your command again: " + err.Error())
			os.Exit(1)
		}

		log.Println("authentication error: " + err.Error())
		os.Exit(1)
	}

	return serverservice.NewClientWithToken(token.AccessToken, mctl.Config.ServerserviceEndpoint, nil)
}
