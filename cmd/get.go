package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get resource",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware", "firmware-set", "component", "server", "attributes", "versioned-attributes"}
		log.Fatal("A valid get command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	rootCmd.AddCommand(cmdGet)
	cmdGet.AddCommand(cmdGetComponent)

	cmdGet.PersistentFlags().BoolVar(&outputJSON, "output-json", false, "Output listing as JSON")
}
