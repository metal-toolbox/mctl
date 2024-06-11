package create

import (
	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "Create resources",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(create)
	create.AddCommand(createFirmware)
	create.AddCommand(createFirmwareSet)
	create.AddCommand(uploadBomFile)
	create.AddCommand(serverEnroll)
}
