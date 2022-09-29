package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var (
	outputJSON bool
)

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware", "firmware-set", "component", "server", "attributes", "versioned-attributes"}
		log.Fatal("A valid list command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	rootCmd.AddCommand(cmdList)
	cmdList.AddCommand(cmdListFirmware)
	cmdList.AddCommand(cmdListFirmwareSet)

	cmdList.PersistentFlags().BoolVar(&outputJSON, "output-json", false, "Output listing as JSON")
}
