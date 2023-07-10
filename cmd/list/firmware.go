package list

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	cmdTimeout = 20 * time.Second
)

// List
var listFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "List firmware",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), cmdTimeout)
		defer cancel()

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		firmware, _, err := client.ListServerComponentFirmware(ctx, nil)
		if err != nil {
			log.Fatal("serverservice client returned error: ", err)
		}

		if outputJSON {
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
