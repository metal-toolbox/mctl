package deleteResource

import (
	"log"
	"strings"

	"github.com/metal-toolbox/mctl/cmd"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware", "firmware-set"}
		log.Fatal("A valid delete command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	cmd.RootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteFirmwareSet)
}
