package list

import (
	"context"
	"log"
	"os"
	"strings"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
)

type listFirmwareFlags struct {
	vendor    string
	model     string
	component string
	version   string
	limit     int
	page      int
}

var (
	flagsDefinedListFirmware *listFirmwareFlags
)

// List
var listFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "List firmware",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		filterParams := fleetdbapi.ComponentFirmwareVersionListParams{
			Vendor:    strings.ToLower(flagsDefinedListFirmware.vendor),
			Version:   flagsDefinedListFirmware.version,
			Component: flagsDefinedListFirmware.component,
			Pagination: &fleetdbapi.PaginationParams{
				Limit: flagsDefinedListFirmware.limit,
				Page:  flagsDefinedListFirmware.page,
			},
		}

		if flagsDefinedListFirmware.model != "" {
			// TODO - if we really want to search using multiple models
			//
			//  fix the the firmware search in fleetdb, its currently useless
			//  because fleetdb queries the data using an 'AND' instead of an 'OR'
			filterParams.Model = []string{strings.ToLower(flagsDefinedListFirmware.model)}
		}

		if flagsDefinedListFirmware.component != "" {
			filterParams.Component = strings.ToLower(flagsDefinedListFirmware.component)
		}

		firmware, _, err := client.ListServerComponentFirmware(ctx, &filterParams)
		if err != nil {
			log.Fatal("fleetdb API client returned error: ", err)
		}

		if output == mctl.OutputTypeJSON.String() {
			printJSON(firmware)
			os.Exit(0)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"UUID", "Vendor", "Model", "Component", "Version"})
		for _, f := range firmware {
			table.Append([]string{f.UUID.String(), f.Vendor, strings.Join(f.Model, ","), f.Component, f.Version})
		}
		table.Render()
	},
}

func init() {
	flagsDefinedListFirmware = &listFirmwareFlags{limit: 10}

	mctl.AddVendorFlag(listFirmware, &flagsDefinedListFirmware.vendor)
	mctl.AddModelFlag(listFirmware, &flagsDefinedListFirmware.model)
	mctl.AddComponentTypeFlag(listFirmware, &flagsDefinedListFirmware.component)
	mctl.AddFirmwareVersionFlag(listFirmware, &flagsDefinedListFirmware.version)
	mctl.AddPageLimitFlag(listFirmware, &flagsDefinedListFirmware.limit)
	mctl.AddPageFlag(listFirmware, &flagsDefinedListFirmware.page)
}
