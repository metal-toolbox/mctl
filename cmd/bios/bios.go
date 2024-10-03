package bios

import (
	"context"
	"log"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	rctypes "github.com/metal-toolbox/rivets/condition"
	"github.com/spf13/cobra"
)

var (
	biosFlags *biosActionFlags
)

type biosActionFlags struct {
	serverID string
}

func CreateBiosControlCondition(ctx context.Context, action rctypes.BiosControlAction) error {
	theApp := mctl.MustCreateApp(ctx)

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		return err
	}

	serverID, err := biosFlags.ParseServerID()
	if err != nil {
		return err
	}

	params := rctypes.NewBiosControlTaskParameters(serverID, action)

	response, err := client.ServerBiosControl(ctx, params)
	if err != nil {
		return err
	}

	conditionResp, err := mctl.ConditionFromResponse(response)
	if err != nil {
		return err
	}

	log.Printf("status=%d msg=%s conditionID=%s", response.StatusCode, response.Message, conditionResp.ID)

	return err
}

func (f *biosActionFlags) ParseServerID() (uuid.UUID, error) {
	return uuid.Parse(f.serverID)
}

var biosCmd = &cobra.Command{
	Use:   "bios",
	Short: "Manage BIOS settings",
}

func init() {
	biosFlags = &biosActionFlags{}
	mctl.RootCmd.AddCommand(biosCmd)
}
