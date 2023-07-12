package get

import (
	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var (
	output string
)

var get = &cobra.Command{
	Use:   "get",
	Short: "Get resource",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(cmdGet)
	cmdGet.AddCommand(getComponent)
	cmdGet.AddCommand(getCondition)
	cmdGet.AddCommand(getFirmware)
	cmdGet.AddCommand(getFirmwareSet)

	get.PersistentFlags().StringVarP(&output, "output", "o", "json", "{json|text}")
}
