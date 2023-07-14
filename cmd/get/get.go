package get

import (
	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var (
	output string
)

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get resource",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:errcheck // returns nil
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(cmdGet)
	cmdGet.AddCommand(getComponent)
	cmdGet.AddCommand(getCondition)
	cmdGet.AddCommand(getFirmwareAvailable)
	cmdGet.AddCommand(getServerFirmware)
	cmdGet.AddCommand(getFirmwareSet)
	cmdGet.AddCommand(getBiosConfig)

	cmdGet.PersistentFlags().StringVarP(&output, "output", "o", "json", "{json|text}")
}
