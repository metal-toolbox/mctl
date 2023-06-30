package edit

import (
	"log"
	"strings"

	"github.com/metal-toolbox/mctl/cmd"

	"github.com/spf13/cobra"
)

var edit = &cobra.Command{
	Use:   "edit",
	Short: "Edit resources",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware-set"}
		log.Fatal("A valid edit command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	cmd.RootCmd.AddCommand(edit)
	edit.AddCommand(editFirmwareSet)
}
