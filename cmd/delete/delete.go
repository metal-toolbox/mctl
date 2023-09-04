package deleteresource

import (
	"github.com/metal-toolbox/mctl/cmd"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:errcheck // returns nil
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteFirmwareSet)
	deleteCmd.AddCommand(deleteCondition)
	deleteCmd.AddCommand(deleteFirmware)
}
