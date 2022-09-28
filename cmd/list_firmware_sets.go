package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdListFirmwareSets = &cobra.Command{
	Use:   fmt.Sprintf("firmware"),
	Short: "List firmware Sets",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cmdList.AddCommand(cmdListFirmwareSets)
}
