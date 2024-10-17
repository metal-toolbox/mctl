package bios

import (
	"errors"

	"github.com/google/uuid"
	"github.com/metal-toolbox/conditionorc/pkg/api/v1/conditions/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var (
	biosFlags *biosActionFlags
)

type biosActionFlags struct {
	serverID string
}

func (f *biosActionFlags) ToCondition() (*types.ConditionCreate, error) {
	/*id, err := f.ParseServerID()
	if err != nil {
		return nil, err
	}*/

	return nil, errors.New("unimplemented")
	/* XXX: updating conditionorc hauled in a change to rivets that makes this no longer compile.
	biosParams := rctypes.NewBiosControlTaskParameters(id, rctypes.ResetSettings)

	params, err := json.Marshal(biosParams)
	if err != nil {
		return nil, err
	}

	return &types.ConditionCreate{Parameters: params}, nil*/
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
