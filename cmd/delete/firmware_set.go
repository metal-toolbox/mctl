package deleteresource

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
)

var (
	deleteFWSetFlags mctl.FirmwareSetFlags
)

var deleteFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Delete a firmware set",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(deleteFWSetFlags.ID)
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.DeleteServerComponentFirmwareSet(cmd.Context(), id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("firmware set deleted: " + id.String())
	},
}

func init() {
	mctl.AddFirmwareSetFlag(deleteFirmwareSet, &deleteFWSetFlags.ID)
	mctl.RequireFlag(deleteFirmwareSet, mctl.FirmwareSetFlag)
}
