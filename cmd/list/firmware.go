package list

import (
	"context"
	"log"
	"os"
	"strings"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

type listFirmwareFlags struct {
	server    string // server UUID
	vendor    string
	models    []string
	component string
	version   string
}

var (
	flagsDefinedListFirmware *listFirmwareFlags
)

// List
var listFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "List firmware",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		lowerCasedModels := func(models []string) []string {
			lowered := []string{}
			for _, m := range models {
				lowered = append(lowered, strings.ToLower(m))
			}

			return lowered
		}

		filterParams := serverservice.ComponentFirmwareVersionListParams{
			Vendor:  strings.ToLower(flagsDefinedListFirmware.vendor),
			Model:   lowerCasedModels(flagsDefinedListFirmware.models),
			Version: flagsDefinedListFirmware.version,
		}

		firmware, _, err := client.ListServerComponentFirmware(ctx, &filterParams)
		if err != nil {
			log.Fatal("serverservice client returned error: ", err)
		}

		if outputJSON {
			printJSON(firmware)
			os.Exit(0)
		}

		// the built in filter only filters out vendor, model, and version, will have to filter out the other columns manually
		if flagsDefinedListFirmware.server != "" || flagsDefinedListFirmware.component != "" {
			filteredFirmware := make([]serverservice.ComponentFirmwareVersion, 0)
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
	flagsDefinedListFirmware = &listFirmwareFlags{}

	listFirmware.PersistentFlags().StringVar(&flagsDefinedListFirmware.server, "server", "", "server UUID")
	listFirmware.PersistentFlags().StringVar(&flagsDefinedListFirmware.vendor, "vendor", "", "vendor name")
	listFirmware.PersistentFlags().StringSliceVar(&flagsDefinedListFirmware.models, "models", nil, "one or more comma separated models numbers")
	listFirmware.PersistentFlags().StringVar(&flagsDefinedListFirmware.component, "component", "", "component type")
	listFirmware.PersistentFlags().StringVar(&flagsDefinedListFirmware.version, "version", "", "version number")
}
