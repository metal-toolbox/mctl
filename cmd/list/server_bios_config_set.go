package list

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

//nolint:err113 // brevity is best here
func sendListServerBiosConfigSetRequest(client *fleetdbapi.Client, cmd *cobra.Command) (*[]fleetdbapi.BiosConfigSet, error) {
	params := &fleetdbapi.BiosConfigSetListParams{
		Pagination: fleetdbapi.PaginationParams{
			Preload: true,
		},
	}

	biosConfig, err := client.ListServerBiosConfigSet(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("retrieving bios configs: %w", err)
	}

	if biosConfig.TotalRecordCount == 0 {
		return nil, errors.New("no bios configs identified")
	}

	return biosConfig.Records.(*[]fleetdbapi.BiosConfigSet), nil
}

// List
var listServerBiosConfigSet = &cobra.Command{
	Use:   "bios-config-set",
	Short: "List bios config",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		biosConfig, err := sendListServerBiosConfigSetRequest(client, cmd)
		if err != nil {
			log.Fatal(err)
		}

		if output == mctl.OutputTypeJSON.String() {
			printJSON(biosConfig)
			os.Exit(0)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Version"})
		for _, s := range *biosConfig {
			table.Append([]string{s.ID, s.Name, s.Version})
		}

		table.SetAutoMergeCells(true)
		table.Render()
	},
}

func init() {

}
