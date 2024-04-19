package bios

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	rctypes "github.com/metal-toolbox/rivets/condition"
	"github.com/spf13/cobra"
)

var (
	biosFlags *biosActionFlags
)

type biosActionFlags struct {
	serverID string
}

func (f *biosActionFlags) ToCondition() (*types.ConditionCreate, error) {
	id, err := f.ParseServerID()
	if err != nil {
		return nil, err
	}

	biosParams := rctypes.NewBiosControlTaskParameters(id, rctypes.ResetSettings)

	params, err := json.Marshal(biosParams)
	if err != nil {
		return nil, err
	}

	return &types.ConditionCreate{Parameters: params}, nil
}

func (f *biosActionFlags) ParseServerID() (uuid.UUID, error) {
	return uuid.Parse(biosFlags.serverID)
}

var biosCmd = &cobra.Command{
	Use:   "bios",
	Short: "Manage BIOS settings",
}

func init() {
	biosFlags = &biosActionFlags{}
	mctl.RootCmd.AddCommand(biosCmd)
}
