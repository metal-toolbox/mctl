package generate

import (
	"log"
	"os"

	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
	cobradoc "github.com/spf13/cobra/doc"
)

var cmdGenerateDocs = &cobra.Command{
	Use:   "gendocs",
	Short: "Generate markdown docs for mctl CLI",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.MkdirAll("docs", os.ModeSticky|os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		if err := cobradoc.GenMarkdownTree(cmd.Root(), "./docs"); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(cmdGenerateDocs)
}
