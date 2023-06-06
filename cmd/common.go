package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	ErrAttributeFromLabel = errors.New("error creating Attribute from Label")
	ErrLabelFromAttribute = errors.New("error creating Label from Attribute")
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

func findAttribute(ns string, attributes []serverservice.Attributes) *serverservice.Attributes {
	for _, attribute := range attributes {
		if attribute.Namespace == ns {
			return &attribute
		}
	}

	return nil
}

func attributeFromLabels(ns string, labels map[string]string) (*serverservice.Attributes, error) {
	data, err := json.Marshal(labels)
	if err != nil {
		return nil, errors.Wrap(ErrAttributeFromLabel, err.Error())
	}

	return &serverservice.Attributes{
		Namespace: ns,
		Data:      data,
	}, nil
}

func printJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
