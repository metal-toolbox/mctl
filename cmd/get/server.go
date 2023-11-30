package get

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	rts "github.com/metal-toolbox/rivets/serverservice"
	rt "github.com/metal-toolbox/rivets/types"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	ss "go.hollow.sh/serverservice/pkg/api/v1"
)

type getServerFlags struct {
	// server UUID
	id             string
	component      string
	listComponents bool
	biosconfig     bool
	table          bool
}

var (
	cmdArgs    *getServerFlags
	cmdTimeout = 2 * time.Minute
)

var getServer = &cobra.Command{
	Use:   "server",
	Short: "Get server information",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(cmd.Context(), cmdTimeout)
		defer cancel()

		theApp := mctl.MustCreateApp(ctx)

		client, err := app.NewServerserviceClient(ctx, theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(cmdArgs.id)
		if err != nil {
			log.Fatal(err)
		}

		withComponents := cmdArgs.listComponents || cmdArgs.component != ""
		server, err := server(ctx, client, id, withComponents)
		if err != nil {
			log.Fatal(err)
		}

		switch {
		case cmdArgs.listComponents:
			renderComponentListTable(server.Components)
			os.Exit(0)
		case cmdArgs.component != "":
			printComponent(server.Components, cmdArgs.component)
			os.Exit(0)
		case cmdArgs.biosconfig:
			mctl.PrintResults(output, server.BIOSCfg)
			os.Exit(0)
		case cmdArgs.table:
			renderServerTable(server)
			os.Exit(0)
		}

		fmt.Println("here")
		mctl.PrintResults(output, server)
	},
}

func printComponent(components []*rt.Component, slug string) {
	got := []*rt.Component{}

	for _, c := range components {
		c := c
		if strings.EqualFold(slug, c.Name) {
			got = append(got, c)
		}
	}

	mctl.PrintResults(output, got)
}

func renderServerTable(server *rt.Server) {
	tableServer := tablewriter.NewWriter(os.Stdout)
	tableServer.Append([]string{"ID", server.ID})
	tableServer.Append([]string{"Name", server.Name})
	tableServer.Append([]string{"Model", server.Model})
	tableServer.Append([]string{"Vendor", server.Vendor})
	tableServer.Append([]string{"Serial", server.Serial})
	tableServer.Append([]string{"BMC", server.BMCAddress})
	tableServer.Append([]string{"Facility", server.Facility})
	tableServer.Append([]string{"Reported", humanize.Time(server.UpdatedAt)})

	tableServer.Render()
}

func renderComponentListTable(components []*rt.Component) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Component", "Vendor", "Model", "Serial", "FW", "Status", "Reported"})
	for _, c := range components {
		vendor := "-"
		model := "-"
		serial := "-"
		installed := "-"
		status := "-"

		if c.Firmware != nil && c.Firmware.Installed != "" {
			installed = c.Firmware.Installed
		}

		if c.Status != nil {
			if c.Status.Health != "" {
				status = c.Status.Health
			} else if c.Status.State != "" {
				status = c.Status.State
			}
		}

		if c.Vendor != "" {
			vendor = c.Vendor
		}

		if c.Model != "" {
			model = c.Model
		}

		if c.Serial != "" {
			serial = c.Serial
		}

		table.Append([]string{c.Name, vendor, model, serial, installed, status, humanize.Time(c.UpdatedAt)})
	}

	table.Render()
}

func server(ctx context.Context, client *ss.Client, id uuid.UUID, withComponents bool) (*rt.Server, error) {
	server, _, err := client.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	cserver := rts.ConvertServer(server)
	if withComponents {
		var err error
		cserver.Components, err = components(ctx, client, id)
		if err != nil {
			return nil, err
		}
	}

	return cserver, nil
}

func components(ctx context.Context, c *ss.Client, id uuid.UUID) ([]*rt.Component, error) {
	params := &ss.PaginationParams{}

	components, resp, err := c.GetComponents(ctx, id, params)
	if err != nil {
		return nil, err
	}

	// XXX: the default result-set size is 100, so more than 100 components will trip
	// the following error.
	if resp.TotalPages > 0 {
		return nil, errors.New("too many components -- add pagination")
	}

	return rts.ConvertComponents(components), nil
}

func init() {
	cmdArgs = &getServerFlags{}

	cmdPFlags := getServer.PersistentFlags()

	cmdPFlags.StringVar(&cmdArgs.id, "id", "", "server UUID")
	cmdPFlags.StringVar(&cmdArgs.component, "component", "c", "component slug")
	cmdPFlags.BoolVarP(&cmdArgs.listComponents, "list-components", "l", false, "include component data")
	cmdPFlags.BoolVarP(&cmdArgs.biosconfig, "bioscfg", "b", false, "print bios configuration")
	cmdPFlags.BoolVarP(&cmdArgs.table, "table", "t", false, "format output in a table")

	if err := getServer.MarkPersistentFlagRequired("id"); err != nil {
		log.Fatal(err)
	}
}
