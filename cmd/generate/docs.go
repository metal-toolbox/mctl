package generate

import "github.com/spf13/cobra"

// List
var listFirmware = &cobra.Command{
	Use:   "docs",
	Short: "Generate markdown docs",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
