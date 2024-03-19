package list

import (
	"fmt"
	"log"
	"os"
	"strings"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	rfleetdb "github.com/metal-toolbox/rivets/fleetdb"
	rt "github.com/metal-toolbox/rivets/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
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
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()
		theApp := mctl.MustCreateApp(ctx)

		client, err := app.NewFleetDBAPIClient(ctx, theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		lsp := &fleetdbapi.
			ServerListParams{
			FacilityCode:        fdlServer.facility,
			AttributeListParams: attributeParamsFromFlags(fdlServer),
			PaginationParams: &fleetdbapi.
				PaginationParams{
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
			rtServers = append(rtServers, rfleetdb.ConvertServer(&s))
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

func attributeParamsFromFlags(fl *listServerFlags) []fleetdbapi.
	AttributeListParams {
	alp := []fleetdbapi.
		AttributeListParams{}

	// match by vendor, model attributes
	if fl.vendor != "" {
		alp = append(
			alp,
			fleetdbapi.
				AttributeListParams{
				Namespace: rfleetdb.ServerAttributeNSVendor,
				Keys:      []string{"vendor"},
				Operator:  "eq",
				Value:     strings.ToLower(fdlServer.vendor),
			},
		)
	}

	if fl.model != "" {
		alp = append(
			alp,
			fleetdbapi.
				AttributeListParams{
				Namespace: rfleetdb.ServerAttributeNSVendor,
				Keys:      []string{"model"},
				Operator:  "like",
				Value:     strings.ToLower(fdlServer.model),
			},
		)
	}

	if fl.serial != "" {
		alp = append(
			alp,
			fleetdbapi.
				AttributeListParams{
				Namespace: rfleetdb.ServerAttributeNSVendor,
				Keys:      []string{"serial"},
				Operator:  "eq",
				Value:     strings.ToLower(fdlServer.serial),
			},
		)
	}

	if fl.bmcerrors {
		alp = append(
			alp,
			fleetdbapi.
				AttributeListParams{
				Namespace: rfleetdb.ServerNSBMCErrorsAttribute,
			},
		)
	}

	return alp
}

func init() {
	fdlServer = &listServerFlags{}

	mctl.AddWithRecordsFlag(cmdListServer, &fdlServer.records)
	mctl.AddVendorFlag(cmdListServer, &fdlServer.vendor)
	mctl.AddModelFlag(cmdListServer, &fdlServer.model)
	mctl.AddFacilityFlag(cmdListServer, &fdlServer.facility)
	mctl.AddPageFlag(cmdListServer, &fdlServer.page)
	mctl.AddPageLimitFlag(cmdListServer, &fdlServer.limit)
	mctl.AddWithBMCErrorsFlag(cmdListServer, &fdlServer.bmcerrors)
	mctl.AddWithCredsFlag(cmdListServer, &fdlServer.creds)
	mctl.AddPrintTableFlag(cmdListServer, &fdlServer.table)
	mctl.AddServerSerialFlag(cmdListServer, &fdlServer.serial)
}
