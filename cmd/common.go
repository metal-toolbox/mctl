package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/metal-toolbox/mctl/internal/app"
	"golang.org/x/net/context"

	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

func newServerserviceClient(ctx context.Context, mctl *app.App) (*serverservice.Client, error) {
	accessToken := "fake"

	if !mctl.Config.DisableOAuth {
		token, err := mctl.RefreshToken(
			ctx,
			mctl.Config.OidcClientID,
			mctl.Config.OidcIssuerEndpoint,
		)
		if err != nil {
			if strings.Contains(err.Error(), "secret not found in keyring") {
				log.Println("please run `mctl auth` and retry your command: " + err.Error())
				os.Exit(1)
			}

			log.Println("authentication error: " + err.Error())
			os.Exit(1)
		}

		accessToken = token.AccessToken
	}

	return serverservice.NewClientWithToken(accessToken, mctl.Config.ServerserviceEndpoint, nil)
}

func printJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
