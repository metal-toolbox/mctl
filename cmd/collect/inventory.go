package collect

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	coapiv1 "github.com/metal-toolbox/conditionorc/pkg/api/v1/conditions/types"
	rctypes "github.com/metal-toolbox/rivets/v2/condition"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
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
	Run: func(cmd *cobra.Command, _ []string) {
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
		Parameters: params,
	}

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.ServerConditionCreate(ctx, serverID, rctypes.Inventory, conditionCreate)
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

	mctl.AddServerFlag(collectInventoryCmd, &flagsDefinedCollectInventory.serverID)
	mctl.AddSkipFWStatusFlag(collectInventoryCmd, &flagsDefinedCollectInventory.skipFirmwareStatusCollect)
	mctl.AddSkipBiosConfigFlag(collectInventoryCmd, &flagsDefinedCollectInventory.skipBiosConfigCollect)

	mctl.RequireFlag(collectInventoryCmd, mctl.ServerFlag)
}
