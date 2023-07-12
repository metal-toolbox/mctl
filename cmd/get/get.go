package get

import (
	"log"
	"strings"

	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var (
	output string
)

var get = &cobra.Command{
	Use:   "get",
	Short: "Get resource",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware", "firmware-set", "component", "server", "attributes", "versioned-attributes"}
		log.Fatal("A valid get command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	cmd.RootCmd.AddCommand(get)
	get.AddCommand(getComponent)
	get.AddCommand(getCondition)
	get.AddCommand(getServerFirmware)

	get.PersistentFlags().StringVarP(&output, "output", "o", "json", "{json|text}")
}
