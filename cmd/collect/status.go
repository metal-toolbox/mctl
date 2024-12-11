package collect

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	rctypes "github.com/metal-toolbox/rivets/v2/condition"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
)

type inventoryStatusParams struct {
	serverID string
}

var inventoryStatusFlags *inventoryStatusParams

var inventoryStatus = &cobra.Command{
	Use:   "status",
	Short: "check the progress of a inventory collection for a server",
	Run: func(cmd *cobra.Command, _ []string) {
		statusCheck(cmd.Context())
	},
}

func statusCheck(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	serverID, err := uuid.Parse(inventoryStatusFlags.serverID)
	if err != nil {
		log.Fatalf("parsing server id: %s", err.Error())
	}

	resp, err := client.ServerConditionStatus(ctx, serverID)
	if err != nil {
		log.Fatalf("querying server conditions: %s", err.Error())
	}

	s, err := mctl.FormatConditionResponse(resp, rctypes.Inventory)
	if err != nil {
		log.Fatalf("condition response error: %s", err.Error())
	}

	fmt.Println(s)
}

func init() {
	inventoryStatusFlags = &inventoryStatusParams{}

	mctl.AddServerFlag(inventoryStatus, &inventoryStatusFlags.serverID)
	mctl.RequireFlag(inventoryStatus, mctl.ServerFlag)
}
