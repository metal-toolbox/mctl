package install

import (
	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var install = &cobra.Command{
	Use:   "install",
	Short: "Install actions",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:errcheck // returns nil
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(install)
}
