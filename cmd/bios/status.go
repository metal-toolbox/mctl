package bios

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	rctypes "github.com/metal-toolbox/rivets/condition"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get condition status of server",
	Run: func(cmd *cobra.Command, _ []string) {
		getConditionStatus(cmd.Context())
	},
}

func getConditionStatus(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	id, err := uuid.Parse(biosFlags.serverID)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.ServerConditionStatus(ctx, id)
	if err != nil {
		log.Fatalf("querying server conditions: %s", err.Error())
	}

	formattedResponse, err := mctl.FormatConditionResponse(resp, rctypes.BiosControl)
	if err != nil {
		log.Fatalf("condition response error: %s", err.Error())
	}

	fmt.Println(formattedResponse)
}

func init() {
	mctl.AddServerFlag(statusCmd, &biosFlags.serverID)

	mctl.RequireFlag(statusCmd, mctl.ServerFlag)

	biosCmd.AddCommand(statusCmd)
}
