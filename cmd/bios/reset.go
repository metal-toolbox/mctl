package bios

import (
	"context"
	"log"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	rctypes "github.com/metal-toolbox/rivets/condition"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset BIOS settings to default values",
	Run: func(cmd *cobra.Command, _ []string) {
		biosAction(cmd.Context(), biosFlags)
	},
}

func biosAction(ctx context.Context, flags *biosActionFlags) {
	theApp := mctl.MustCreateApp(ctx)

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	serverID, err := flags.ParseServerID()
	if err != nil {
		log.Fatal(err)
	}

	condition, err := flags.ToCondition()
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.ServerConditionCreate(ctx, serverID, rctypes.BiosControl, *condition)
	if err != nil {
		log.Fatal(err)
	}

	conditionResp, err := mctl.ConditionFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%d msg=%s conditionID=%s", response.StatusCode, response.Message, conditionResp.ID)
}

func init() {
	mctl.AddServerFlag(resetCmd, &biosFlags.serverID)

	mctl.RequireFlag(resetCmd, mctl.ServerFlag)

	biosCmd.AddCommand(resetCmd)
}
