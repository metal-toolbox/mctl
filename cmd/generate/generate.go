package generate

import (
	"log"
	"os"

	"github.com/metal-toolbox/mctl/cmd"
	"github.com/spf13/cobra"
	cobradoc "github.com/spf13/cobra/doc"
)

var (
	output string
)

var cmdGenerate = &cobra.Command{
	Use:   "generate",
	Short: "Generate CLI docs",
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
	cmd.RootCmd.AddCommand(cmdGenerate)
}
