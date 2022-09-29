package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var cmdEdit = &cobra.Command{
	Use:   "edit",
	Short: "Edit resources",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware-set"}
		log.Fatal("A valid edit command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	rootCmd.AddCommand(cmdEdit)
	cmdEdit.AddCommand(cmdEditFirmwareSet)
}
