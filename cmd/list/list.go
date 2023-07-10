package list

import (
	"log"
	"strings"

	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var (
	outputJSON bool
)

var list = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware", "firmware-set", "component", "server", "attributes", "versioned-attributes"}
		log.Fatal("A valid list command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	cmd.RootCmd.AddCommand(list)
	list.AddCommand(listFirmware)
	list.AddCommand(listFirmwareSet)
	list.AddCommand(listCondition)

	list.PersistentFlags().BoolVar(&outputJSON, "output-json", false, "Output listing as JSON")
}
