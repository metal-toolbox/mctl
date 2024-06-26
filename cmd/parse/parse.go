package parse

import (
	"github.com/spf13/cobra"

	"github.com/metal-toolbox/mctl/cmd"
)

type parseActionFlags struct {
	JsonFileToParse string
	OutputCSVFile string
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse Stuff",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

var flagsDefinedParseAction *parseActionFlags

func init() {
	cmd.RootCmd.AddCommand(parseCmd)
	parseCmd.AddCommand(flasherCmd)

	flagsDefinedParseAction = &parseActionFlags{}
	parseCmd.PersistentFlags().StringVarP(&flagsDefinedParseAction.JsonFileToParse, "json-in", "i", "", "input json")
	parseCmd.PersistentFlags().StringVarP(&flagsDefinedParseAction.OutputCSVFile, "csv-out", "o", "", "output csv")
}