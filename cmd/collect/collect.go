package collect

import (
	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var collect = &cobra.Command{
	Use:   "collect",
	Short: "Collect current server firmware status and bios configuration",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:errcheck // returns nil
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(collect)
}
