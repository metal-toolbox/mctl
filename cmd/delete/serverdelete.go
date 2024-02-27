package deleteresource

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/metal-toolbox/mctl/internal/app"

	mctl "github.com/metal-toolbox/mctl/cmd"
)

type serverDeleteParams struct {
	serverID string
}

var (
	serverDeleteFlags *serverDeleteParams
)

var serverDelete = &cobra.Command{
	Use:   "server",
	Short: "Delete server from fleetDB",
	Run: func(cmd *cobra.Command, args []string) {
		deleteServer(cmd.Context())
	},
}

func deleteServer(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.ServerDelete(ctx, serverDeleteFlags.serverID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("status=%d\nmsg=%s\nserverID=%v", response.StatusCode, response.Message, response.Records.ServerID)

}

func init() {
	serverDeleteFlags = &serverDeleteParams{}

	mctl.AddServerFlag(serverDelete, &serverDeleteFlags.serverID)
	mctl.RequireFlag(serverDelete, mctl.ServerFlag)
}
