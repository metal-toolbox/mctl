package get

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	rts "github.com/metal-toolbox/rivets/serverservice"
	rt "github.com/metal-toolbox/rivets/types"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	ss "go.hollow.sh/serverservice/pkg/api/v1"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
)

type getServerFlags struct {
	// server UUID
	id             string
	component      string
	listComponents bool
	biosconfig     bool
	table          bool
	creds          bool
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
		server, err := server(ctx, client, id, withComponents, cmdArgs.creds)
		if err != nil {
			log.Fatal(err)
		}

		if cmdArgs.table {
			switch {
			case cmdArgs.listComponents:
				renderComponentListTable(server.Components)
			default:
				renderServerTable(server, cmdArgs.creds)
			}

			os.Exit(0)
		}

		if cmdArgs.component != "" {
			printComponent(server.Components, cmdArgs.component)
			os.Exit(0)
		}

		if cmdArgs.biosconfig {
			mctl.PrintResults(output, server.BIOSCfg)
			os.Exit(0)
		}

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

func renderServerTable(server *rt.Server, withCreds bool) {
	tableServer := tablewriter.NewWriter(os.Stdout)
	tableServer.Append([]string{"ID", server.ID})
	tableServer.Append([]string{"Name", server.Name})
	tableServer.Append([]string{"Model", server.Model})
	tableServer.Append([]string{"Vendor", server.Vendor})
	tableServer.Append([]string{"Serial", server.Serial})
	tableServer.Append([]string{"BMCAddr", server.BMCAddress})
	if withCreds {
		tableServer.Append([]string{"BMCUser", server.BMCUser})
		tableServer.Append([]string{"BMCPassword", server.BMCPassword})
	}
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

func server(ctx context.Context, client *ss.Client, id uuid.UUID, withComponents, withCreds bool) (*rt.Server, error) {
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

	if withCreds {
		if err := mctl.ServerBMCCredentials(ctx, client, cserver); err != nil {
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

	mctl.AddServerFlag(getServer, &cmdArgs.id)
	mctl.AddSlugFlag(getServer, &cmdArgs.component, "list component on server by slug (drive/nic/cpu..)")
	mctl.AddWithCredsFlag(getServer, &cmdArgs.creds)
	mctl.AddPrintTableFlag(getServer, &cmdArgs.table)
	mctl.AddBIOSConfigFlag(getServer, &cmdArgs.biosconfig)
	mctl.AddListComponentsFlag(getServer, &cmdArgs.listComponents)

	mctl.RequireFlag(getServer, mctl.ServerFlag)
}
