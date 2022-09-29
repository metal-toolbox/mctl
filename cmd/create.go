package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create resources",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware", "firmware-set"}
		log.Fatal("A valid create command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	rootCmd.AddCommand(cmdCreate)
	cmdCreate.AddCommand(cmdCreateFirmware)
	cmdCreate.AddCommand(cmdCreateFirmwareSet)
}
