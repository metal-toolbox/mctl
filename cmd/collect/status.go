package collect

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"

	cotypes "github.com/metal-toolbox/conditionorc/pkg/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	rctypes "github.com/metal-toolbox/rivets/condition"
)

var inventoryStatus = &cobra.Command{
	Use:   "status --server | -s <server uuid>",
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

	serverID, err := uuid.Parse(serverIDStr)
	if err != nil {
		log.Fatalf("parsing server id: %s", err.Error())
	}

	// cotypes to be moved into rivets library
	resp, err := client.ServerConditionGet(ctx, serverID, cotypes.ConditionKind(rctypes.Inventory))
	if err != nil {
		log.Fatalf("querying server conditions: %s", err.Error())
	}

	s, err := mctl.FormatConditionResponse(resp)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(s)
}

func init() {
	collect.AddCommand(inventoryStatus)
}
