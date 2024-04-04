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
	server    string // server UUID
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
			Vendor: strings.ToLower(flagsDefinedListFirmware.vendor),
			// TODO - if we really want to search using multiple models
			//
			//  fix the the firmware search in fleetdb, its currently useless
			//  because fleetdb queries the data using an 'AND' instead of an 'OR'
			Model:   []string{strings.ToLower(flagsDefinedListFirmware.model)},
			Version: flagsDefinedListFirmware.version,
			Pagination: &fleetdbapi.PaginationParams{
				Limit: flagsDefinedListFirmware.limit,
				Page:  flagsDefinedListFirmware.page,
			},
		}

		firmware, _, err := client.ListServerComponentFirmware(ctx, &filterParams)
		if err != nil {
			log.Fatal("fleetdb API client returned error: ", err)
		}

		if output == mctl.OutputTypeJSON.String() {
			printJSON(firmware)
			os.Exit(0)
		}

		// the built in filter only filters out vendor, model, and version, will have to filter out the other columns manually
		if flagsDefinedListFirmware.server != "" || flagsDefinedListFirmware.component != "" {
			filteredFirmware := make([]fleetdbapi.ComponentFirmwareVersion, 0)
			for _, f := range firmware {
				if (flagsDefinedListFirmware.server == "" || f.UUID.String() == flagsDefinedListFirmware.server) &&
					(flagsDefinedListFirmware.component == "" || f.Component == flagsDefinedListFirmware.component) {
					filteredFirmware = append(filteredFirmware, f)
				}
			}
			firmware = filteredFirmware
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

	mctl.AddServerFlag(listFirmware, &flagsDefinedListFirmware.server)
	mctl.AddVendorFlag(listFirmware, &flagsDefinedListFirmware.vendor)
	mctl.AddModelFlag(listFirmware, &flagsDefinedListFirmware.model)
	mctl.AddComponentTypeFlag(listFirmware, &flagsListComponent.slug)
	mctl.AddFirmwareVersionFlag(listFirmware, &flagsDefinedListFirmware.version)
	mctl.AddPageLimitFlag(listFirmware, &flagsDefinedListFirmware.limit)
	mctl.AddPageFlag(listFirmware, &flagsDefinedListFirmware.page)
}
