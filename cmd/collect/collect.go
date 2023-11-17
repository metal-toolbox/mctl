package collect

import (
	"log"

	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var serverIDStr string

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

	pflags := collect.PersistentFlags()
	pflags.StringVarP(&serverIDStr, "server", "s", "", "server id (typically a UUID)")

	if err := collect.MarkPersistentFlagRequired("server"); err != nil {
		log.Fatalf("marking server flag as required: %s", err.Error())
	}
}
