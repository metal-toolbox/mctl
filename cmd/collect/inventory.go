package collect

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/metal-toolbox/mctl/internal/app"

	coapiv1 "github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
	"github.com/metal-toolbox/conditionorc/pkg/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	rctypes "github.com/metal-toolbox/rivets/condition"
)

type collectInventoryFlags struct {
	serverID                  string
	skipFirmwareStatusCollect bool
	skipBiosConfigCollect     bool
}

var (
	flagsDefinedCollectInventory *collectInventoryFlags
)

var collectInventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Collect current server firmware status and bios configuration",
	Run: func(cmd *cobra.Command, args []string) {
		collectInventory(cmd.Context())

	},
}

func collectInventory(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	serverID, err := uuid.Parse(flagsDefinedCollectInventory.serverID)
	if err != nil {
		log.Fatal(err)
	}

	params, err := json.Marshal(rctypes.NewInventoryTaskParameters(
		serverID,
		rctypes.OutofbandInventory,
		!flagsDefinedCollectInventory.skipFirmwareStatusCollect,
		!flagsDefinedCollectInventory.skipBiosConfigCollect,
	))
	if err != nil {
		log.Fatal(err)
	}

	conditionCreate := coapiv1.ConditionCreate{
		Exclusive:  false,
		Parameters: params,
	}

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.ServerConditionCreate(ctx, serverID, types.ConditionKind(rctypes.Inventory), conditionCreate)
	if err != nil {
		log.Fatal(err)
	}

	condition, err := mctl.ConditionFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%d msg=%s conditionID=%s", response.StatusCode, response.Message, condition.ID)

}

func init() {
	flagsDefinedCollectInventory = &collectInventoryFlags{}

	collect.AddCommand(collectInventoryCmd)
	collectInventoryCmd.PersistentFlags().StringVar(&flagsDefinedCollectInventory.serverID, "server", "", "server UUID")
	collectInventoryCmd.PersistentFlags().BoolVar(&flagsDefinedCollectInventory.skipFirmwareStatusCollect, "skip-fw-status", false, "Skip firmware status data collection")
	collectInventoryCmd.PersistentFlags().BoolVar(&flagsDefinedCollectInventory.skipBiosConfigCollect, "skip-bios-config", false, "Skip BIOS configuration data collection")

	if err := collectInventoryCmd.MarkPersistentFlagRequired("server"); err != nil {
		log.Fatal(err)
	}
}
