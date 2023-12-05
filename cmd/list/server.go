package list

import (
	"fmt"
	"log"
	"os"
	"strings"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	rts "github.com/metal-toolbox/rivets/serverservice"
	rt "github.com/metal-toolbox/rivets/types"
	ss "go.hollow.sh/serverservice/pkg/api/v1"
)

type listServerFlags struct {
	records   bool
	bmcerrors bool
	creds     bool
	table     bool
	vendor    string
	model     string
	serial    string
	facility  string
	limit     int
	page      int
}

var (
	fdlServer *listServerFlags
)

// List
var cmdListServer = &cobra.Command{
	Use:   "server",
	Short: "List servers",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		theApp := mctl.MustCreateApp(ctx)

		client, err := app.NewServerserviceClient(ctx, theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		lsp := &ss.ServerListParams{
			FacilityCode:        fdlServer.facility,
			AttributeListParams: attributeParamsFromFlags(fdlServer),
			PaginationParams: &ss.PaginationParams{
				Limit:   fdlServer.limit,
				Page:    fdlServer.page,
				Preload: false,
			},
		}

		servers, res, err := client.List(ctx, lsp)
		if err != nil {
			log.Fatal(err)
		}

		if fdlServer.records {
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

		if len(servers) == 0 {
			fmt.Println("no servers matched filters")
			os.Exit(0)
		}

		rtServers := make([]*rt.Server, 0, len(servers))
		for _, s := range servers {
			s := s
			rtServers = append(rtServers, rts.ConvertServer(&s))
		}

		if fdlServer.creds {
			for idx := range rtServers {
				if err := mctl.ServerBMCCredentials(ctx, client, rtServers[idx]); err != nil {
					log.Fatal(err)
				}
			}
		}

		if fdlServer.table {
			serversTable(rtServers, fdlServer)
			os.Exit(0)
		}

		printJSON(rtServers)
	},
}

func serversTable(servers []*rt.Server, fl *listServerFlags) {
	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"UUID", "Name", "Vendor", "Model", "Serial", "BMCAddr"}

	if fl.creds {
		headers = append(headers, []string{"BMCUser", "BMCPass"}...)
	}

	table.SetHeader(headers)
	for _, server := range servers {
		row := []string{
			server.ID,
			server.Name,
			server.Vendor,
			server.Model,
			server.Serial,
			server.BMCAddress,
		}

		if fl.creds {
			row = append(row, []string{server.BMCUser, server.BMCPassword}...)
		}

		table.Append(row)
	}

	table.Render()
}

func attributeParamsFromFlags(fl *listServerFlags) []ss.AttributeListParams {
	alp := []ss.AttributeListParams{}

	// match by vendor, model attributes
	if fl.vendor != "" {
		alp = append(
			alp,
			ss.AttributeListParams{
				Namespace: rts.ServerAttributeNSVendor,
				Keys:      []string{"vendor"},
				Operator:  "eq",
				Value:     strings.ToLower(fdlServer.vendor),
			},
		)
	}

	if fl.model != "" {
		alp = append(
			alp,
			ss.AttributeListParams{
				Namespace: rts.ServerAttributeNSVendor,
				Keys:      []string{"model"},
				Operator:  "like",
				Value:     strings.ToLower(fdlServer.model),
			},
		)
	}

	if fl.serial != "" {
		alp = append(
			alp,
			ss.AttributeListParams{
				Namespace: rts.ServerAttributeNSVendor,
				Keys:      []string{"serial"},
				Operator:  "eq",
				Value:     strings.ToLower(fdlServer.serial),
			},
		)
	}

	if fl.bmcerrors {
		alp = append(
			alp,
			ss.AttributeListParams{
				Namespace: rts.ServerNSBMCErrorsAttribute,
			},
		)
	}

	return alp
}

func init() {
	fdlServer = &listServerFlags{}

	cmdListServer.PersistentFlags().BoolVar(&fdlServer.records, "records", false, "only print record count matching filters")
	cmdListServer.PersistentFlags().BoolVar(&fdlServer.table, "table", false, "print records in a table format")
	cmdListServer.PersistentFlags().BoolVar(&fdlServer.bmcerrors, "bmcerrors", false, "list servers with BMC errors")
	cmdListServer.PersistentFlags().BoolVar(&fdlServer.creds, "creds", false, "list BMC credentials in ")
	cmdListServer.PersistentFlags().StringVar(&fdlServer.vendor, "vendor", "", "filter by server vendor")
	cmdListServer.PersistentFlags().StringVar(&fdlServer.model, "model", "", "filter by server model")
	cmdListServer.PersistentFlags().StringVar(&fdlServer.facility, "facility", "", "filter by facility code")
	cmdListServer.PersistentFlags().StringVar(&fdlServer.serial, "serial", "", "filter by server serial")
	cmdListServer.PersistentFlags().IntVar(&fdlServer.page, "page", 0, "limit results to page (for use with --limit)")
	cmdListServer.PersistentFlags().IntVar(&fdlServer.limit, "limit", 10, "limit results returned") // nolint:gomnd // value is obvious as is
}
