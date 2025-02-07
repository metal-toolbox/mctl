package create

import (
	"context"
	"encoding/json"
	"log"
	"os"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

var fromFile string

var createServerBiosConfigSet = &cobra.Command{
	Use:   "bios-config-set",
	Short: "Create a bios config set",
	PreRun: func(cmd *cobra.Command, _ []string) {
		fromFile, err := cmd.Flags().GetString(mctl.FromFileFlag.Name())
		if err != nil {
			log.Fatal(err)
		}

		if fromFile == "" {
			mctl.RequireFlag(cmd, mctl.FromFileFlag)
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		err = createServerBiosConfigSetFromFile(cmd.Context(), client)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func createServerBiosConfigSetFromFile(ctx context.Context, client *fleetdbapi.Client) (err error) {
	biosconfigset := &fleetdbapi.BiosConfigSet{}

	fbytes, err := os.ReadFile(fromFile)
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(fbytes, &biosconfigset); err != nil {
		log.Fatal(err)
	}

	_, err = client.CreateServerBiosConfigSet(ctx, *biosconfigset)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func init() {
	mctl.AddFromFileFlag(createServerBiosConfigSet, &fromFile, "path to JSON file containing bios config set")
}
