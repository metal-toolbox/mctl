package serverservice

//import (
//	"context"
//	"net/http"
//	"os"
//
//	"github.com/metal-toolbox/mctl/pkg/model"
//	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
//)
//
//const (
//	// EnvServerserviceSkipOAuth when set to true will skip server service OAuth.
//	EnvVarServerserviceSkipOAuth = "SERVERSERVICE_SKIP_OAUTH"
//)
//
//// NewServerServiceClient instantiates and returns a serverService client
//func NewServerServiceClient(ctx context.Context, cfg *model.Config) (*serverservice.Client, error) {
//	// load configuration parameters from env variables
//	loadServerServiceEnvVars(cfg)
//
//	// validate parameters
//	endpointURL, err := validateServerServiceParams(cfg)
//	if err != nil {
//		return nil, err
//	}
//
//	if os.Getenv(EnvVarServerserviceSkipOAuth) == "true" {
//		return serverservice.NewClientWithToken(
//			"faketoken",
//			endpointURL.String(),
//			http.DefaultClient,
//		)
//	}
//
//	client, err := oauthClient(ctx, cfg)
//	if err != nil {
//		return nil, err
//	}
//
//	return serverservice.NewClientWithToken(
//		cfg.ServerService.ClientSecret,
//		endpointURL.String(),
//		client,
//	)
//}
