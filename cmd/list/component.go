package list

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"

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
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		theApp := mctl.MustCreateApp(ctx)

		client, err := app.NewServerserviceClient(ctx, theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		lp := &serverservice.ServerComponentListParams{
			ServerComponentType: flagsListComponent.slug,
			Vendor:              flagsListComponent.vendor,
			Serial:              flagsListComponent.serial,
			Model:               flagsListComponent.model,
			Pagination: &serverservice.PaginationParams{
				Limit:   flagsListComponent.limit,
				Page:    flagsListComponent.page,
				Preload: false,
			},
		}

		components, res, err := client.ListComponents(ctx, lp)
		if err != nil {
			log.Fatal("serverservice query returned error: " + err.Error())
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
