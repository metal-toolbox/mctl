package deleteresource

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"

	mctl "github.com/metal-toolbox/mctl/cmd"
)

var deleteServerBiosConfigSetID string

var deleteServerBiosConfigSet = &cobra.Command{
	Use:   "bios-config-set",
	Short: "Delete a bios config set",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(deleteServerBiosConfigSetID)
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.DeleteServerBiosConfigSet(cmd.Context(), id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("bios config set deleted: " + id.String())
	},
}

func init() {
	mctl.AddBIOSConfigSetIDFlag(deleteServerBiosConfigSet, &deleteServerBiosConfigSetID)
}
