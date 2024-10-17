//nolint:gocritic,err113 // the commented code and dynamic error are intentional
package bios

import (
	"context"
	"log"
	"net/url"

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
	serverID      string
	biosConfigURL string
}

func CreateBiosControlCondition(ctx context.Context, action rctypes.BiosControlAction) error {
	theApp := mctl.MustCreateApp(ctx)

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		return err
	}

	serverID, err := uuid.Parse(biosFlags.serverID)
	if err != nil {
		return err
	}

	var biosURL *url.URL
	if action == rctypes.SetConfig {
		biosURL, err = url.Parse(biosFlags.biosConfigURL)
		if err != nil {
			return err
		}
	}

	params := rctypes.NewBiosControlTaskParameters(serverID, action, biosURL)

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

var biosCmd = &cobra.Command{
	Use:   "bios",
	Short: "Manage BIOS settings",
}

func init() {
	biosFlags = &biosActionFlags{}
	mctl.RootCmd.AddCommand(biosCmd)
}
