package list

import (
	"log"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
)

type listFirmwareSetFlags struct {
	vendor string
	model  string
}

var (
	flagsDefinedListFwSet *listFirmwareSetFlags
)

// List
var listFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "List firmware",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		var fwSet []serverservice.ComponentFirmwareSet

		if flagsDefinedListFwSet.vendor != "" || flagsDefinedListFwSet.model != "" {
			fwSet, err = mctl.FirmwareSetByVendorModel(cmd.Context(), flagsDefinedListFwSet.vendor, flagsDefinedListFwSet.model, client)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fwSet, _, err = client.ListServerComponentFirmwareSet(cmd.Context(), &serverservice.ComponentFirmwareSetListParams{})
			if err != nil {
				log.Fatal(err)
			}
		}

		if output == mctl.OutputTypeJSON.String() {
			printJSON(fwSet)
			os.Exit(0)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"UUID", "Name", "Labels", "firmware UUID", "Vendor", "Model", "Component", "Version"})
		for _, s := range fwSet {
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

func init() {
	flagsDefinedListFwSet = &listFirmwareSetFlags{}

	mctl.AddModelFlag(listFirmwareSet, &flagsDefinedListFwSet.model)
	mctl.AddVendorFlag(listFirmwareSet, &flagsDefinedListFwSet.vendor)
}
