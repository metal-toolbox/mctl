package install

import (
	"log"
	"strings"

	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
)

var install = &cobra.Command{
	Use:   "install",
	Short: "Install actions",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware-set"}
		log.Fatal("A valid list command parameter was expected: " + strings.Join(commands, ", "))
	},
}

func init() {
	cmd.RootCmd.AddCommand(install)
}
