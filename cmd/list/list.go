package list

import (
	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var (
	outputJSON bool
)

var list = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:errcheck // returns nil
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(list)
	list.AddCommand(listFirmware)
	list.AddCommand(listFirmwareSet)
	list.AddCommand(listComponent)

	list.PersistentFlags().BoolVarP(&outputJSON, "output-json", "j", false, "Output listing as JSON")
}
