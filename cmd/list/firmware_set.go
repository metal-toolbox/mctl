package list

import (
	"log"
	"os"
	"strings"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
	rfleetdb "github.com/metal-toolbox/rivets/fleetdb"
)

type listFirmwareSetFlags struct {
	vendor string
	model  string
	labels map[string]string
}

var (
	flagsDefinedListFwSet *listFirmwareSetFlags
)

func sendListFirmwareRequest(client *fleetdbapi.Client, cmd *cobra.Command) ([]fleetdbapi.ComponentFirmwareSet, error) {
	if flagsDefinedListFwSet.vendor == "" && flagsDefinedListFwSet.model == "" {
		fwSet, _, err := client.ListServerComponentFirmwareSet(cmd.Context(), &fleetdbapi.ComponentFirmwareSetListParams{})
		return fwSet, err
	}

	if len(flagsDefinedListFwSet.labels) != 0 {
		return rfleetdb.FirmwareSetByLabels(cmd.Context(), flagsDefinedListFwSet.vendor, flagsDefinedListFwSet.model, flagsDefinedListFwSet.labels, client)
	}

	return rfleetdb.FirmwareSetByVendorModel(cmd.Context(), flagsDefinedListFwSet.vendor, flagsDefinedListFwSet.model, client)
}

// List
var listFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "List firmware",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		fwSet, err := sendListFirmwareRequest(client, cmd)
		if err != nil {
			log.Fatal(err)
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
	mctl.AddLabelsFlag(listFirmwareSet, &flagsDefinedListFwSet.labels, "Labels to from the firmware set - 'foo=bar,foo2=bar2'")
}
