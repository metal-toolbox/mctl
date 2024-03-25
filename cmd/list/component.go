package list

import (
	"log"
	"os"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/spf13/cobra"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
)

type listComponentFlags struct {
	slug    string
	vendor  string
	serial  string
	model   string
	records bool
	limit   int
	page    int
}

var (
	flagsListComponent *listComponentFlags
)

// List
var listComponent = &cobra.Command{
	Use:   "component",
	Short: "List Components",
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()
		theApp := mctl.MustCreateApp(ctx)

		if flagsListComponent.limit > fleetdbapi.MaxPaginationSize {
			log.Printf("Notice: Limit was set above max, setting limit to %d. If you want to list more than %d components, please use '--page` to index individual pages", fleetdbapi.MaxPaginationSize, fleetdbapi.MaxPaginationSize)
			flagsListComponent.limit = fleetdbapi.MaxPaginationSize
		}

		client, err := app.NewFleetDBAPIClient(ctx, theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		lp := &fleetdbapi.ServerComponentListParams{
			ServerComponentType: flagsListComponent.slug,
			Vendor:              flagsListComponent.vendor,
			Serial:              flagsListComponent.serial,
			Model:               flagsListComponent.model,
			Pagination: &fleetdbapi.PaginationParams{
				Limit:   flagsListComponent.limit,
				Page:    flagsListComponent.page,
				Preload: false,
			},
		}

		components, res, err := client.ListComponents(ctx, lp)
		if err != nil {
			log.Fatal("fleetdb API query returned error: " + err.Error())
		}

		if flagsListComponent.records {
			d := struct {
				CurrentPage      int
				Limit            int
				TotalPages       int
				TotalRecordCount int64
				Link             string
			}{
				res.Page,
				res.PageCount,
				res.TotalPages,
				res.TotalRecordCount,
				res.Links.Self.Href,
			}

			printJSON(d)

			os.Exit(0)
		}

		printJSON(components)
	},
}

func init() {
	flagsListComponent = &listComponentFlags{}

	mctl.AddWithRecordsFlag(listComponent, &flagsListComponent.records)
	mctl.AddSlugFlag(listComponent, &flagsListComponent.slug, "filter by component slug (nic/drive/bmc/bios...)")
	mctl.AddVendorFlag(listComponent, &flagsListComponent.vendor)
	mctl.AddModelFlag(listComponent, &flagsListComponent.model)
	mctl.AddPageFlag(listComponent, &flagsListComponent.page)
	mctl.AddPageLimitFlag(listComponent, &flagsListComponent.limit)

	mctl.RequireFlag(listComponent, mctl.SlugFlag)
}
