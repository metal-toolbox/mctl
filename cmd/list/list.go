package list

import (
	"github.com/spf13/cobra"

	"github.com/metal-toolbox/mctl/cmd"
)

var (
	output string
)

var list = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(list)
	list.AddCommand(listFirmware)
	list.AddCommand(listFirmwareSet)
	list.AddCommand(listComponent)
	list.AddCommand(cmdListServer)

	cmd.AddOutputFlag(list, &output)
}
