package list

import (
	"log"
	"os"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
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

	listComponent.PersistentFlags().BoolVar(&flagsListComponent.records, "records", false, "print record count found with pagination info and return")
	listComponent.PersistentFlags().StringVar(&flagsListComponent.slug, "slug", "", "filter by component slug (nic/drive/bmc/bios...)")
	listComponent.PersistentFlags().StringVar(&flagsListComponent.vendor, "vendor", "", "filter by component vendor")
	listComponent.PersistentFlags().StringVar(&flagsListComponent.model, "model", "", "filter by one or more component models")
	listComponent.PersistentFlags().IntVar(&flagsListComponent.page, "page", 0, "limit results to page (for use with --limit)")
	listComponent.PersistentFlags().IntVar(&flagsListComponent.limit, "limit", 10, "limit results returned") // nolint:gomnd // value is obvious as is

	if err := listComponent.MarkPersistentFlagRequired("slug"); err != nil {
		log.Fatal(err)
	}
}
