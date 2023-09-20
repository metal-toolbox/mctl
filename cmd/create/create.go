package create

import (
	"log"
	"strings"

	"github.com/metal-toolbox/mctl/cmd"

	"github.com/spf13/cobra"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "Create resources",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware", "firmware-set", "uploadbom"}
		log.Fatal("A valid create command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	cmd.RootCmd.AddCommand(create)
	create.AddCommand(createFirmware)
	create.AddCommand(createFirmwareSet)
	create.AddCommand(uploadBomFile)
}
