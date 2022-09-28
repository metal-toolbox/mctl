package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var cmdList = &cobra.Command{
	Use:   fmt.Sprintf("list"),
	Short: "List resources",
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{"firmware", "firmware-set", "component", "server", "attributes", "versioned-attributes"}

		log.Println("A valid list command parameter was expected")
		log.Println("supported command parameters: " + strings.Join(commands, ", "))
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(cmdList)
}
