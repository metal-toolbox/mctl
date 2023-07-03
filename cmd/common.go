package cmd

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	co "github.com/metal-toolbox/conditionorc/pkg/api/v1/client"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	ErrAttributeFromLabel = errors.New("error creating Attribute from Label")
	ErrLabelFromAttribute = errors.New("error creating Label from Attribute")
)

func MustCreateApp(ctx context.Context) *app.App {
	mctl, err := app.New(ctx, cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	return mctl
}

func getAuthToken(ctx context.Context, mctl *app.App) string {
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
	return accessToken
}

func newServerserviceClient(ctx context.Context, mctl *app.App) (*serverservice.Client, error) {
	accessToken := getAuthToken(ctx, mctl)
	return serverservice.NewClientWithToken(accessToken, mctl.Config.ServerserviceEndpoint, nil)
}

func newConditionsClient(ctx context.Context, mctl *app.App) (*co.Client, error) {
	accessToken := getAuthToken(ctx, mctl)
	return co.NewClient(mctl.Config.ConditionsEndpoint,
		co.WithAuthToken(accessToken),
	)
}

func AttributeFromLabels(ns string, labels map[string]string) (*serverservice.Attributes, error) {
	data, err := json.Marshal(labels)
	if err != nil {
		return nil, errors.Wrap(ErrAttributeFromLabel, err.Error())
	}

	return &serverservice.Attributes{
		Namespace: ns,
		Data:      data,
	}, nil
}
