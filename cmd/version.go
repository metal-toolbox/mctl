package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/metal-toolbox/mctl/internal/version"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print mctl version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mctl -- brought to you by Fleet Services")
		fmt.Printf(
			"version: %s\ncommit: %s\nbranch: %s\ngo version: %s\nbuilt-on: %s",
			version.AppVersion, version.GitCommit, version.GitBranch,
			version.GoVersion, version.BuildDate,
		)
	},
}

func init() {
	RootCmd.AddCommand(cmdVersion)
}
