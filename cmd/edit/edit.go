package edit

import (
	"github.com/metal-toolbox/mctl/cmd"

	"github.com/spf13/cobra"
)

var edit = &cobra.Command{
	Use:   "edit",
	Short: "Edit resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(edit)
	edit.AddCommand(editFirmwareSet)
}
