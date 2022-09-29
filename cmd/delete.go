package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware", "firmware-set"}
		log.Fatal("A valid delete command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	rootCmd.AddCommand(cmdDelete)
	cmdDelete.AddCommand(cmdDeleteFirmwareSet)
}
