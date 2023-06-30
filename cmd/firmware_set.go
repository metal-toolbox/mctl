package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// List
var cmdListFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "List firmware",
	Run: func(cmd *cobra.Command, args []string) {
		mctl, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			log.Fatal(err)
		}

		c, err := newServerserviceClient(cmd.Context(), mctl)
		if err != nil {
			log.Fatal(err)
		}

		set, _, err := c.ListServerComponentFirmwareSet(cmd.Context(), nil)
		if err != nil {
			log.Fatal(err)
		}

		if outputJSON {
			printJSON(set)
			os.Exit(0)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"UUID", "Name", "Labels", "firmware UUID", "Vendor", "Model", "Component", "Version"})
		for _, s := range set {
			var labels string
			if len(s.Attributes) > 0 {
				attr := findAttribute(model.AttributeNSFirmwareSetLabels, s.Attributes)
				if attr != nil {
					labels = string(attr.Data)
				}
			}
			table.Append([]string{s.UUID.String(), s.Name, labels, "-", "-", "-", "-", "-"})
			for _, f := range s.ComponentFirmware {
				table.Append([]string{s.UUID.String(), "", "", f.UUID.String(), f.Vendor, strings.Join(f.Model, ","), f.Component, f.Version})
			}
		}

		table.SetAutoMergeCells(true)
		table.Render()
	},
}
