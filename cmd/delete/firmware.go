package deleteresource

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
)

type deleteFirmwareFlags struct {
	// firmware UUID
	id string
}

var (
	flagsDefinedDeleteFirmware *deleteFirmwareFlags
)
var deleteFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "Delete a firmware object",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(flagsDefinedDeleteFirmware.id)
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.DeleteServerComponentFirmware(cmd.Context(), serverservice.ComponentFirmwareVersion{UUID: id})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("firmware deleted: " + id.String())
	},
}

func init() {
	flagsDefinedDeleteFirmware = &deleteFirmwareFlags{}

	mctl.AddFirmwareFlag(deleteFirmware, &flagsDefinedDeleteFirmware.id)
	mctl.RequireFlag(deleteFirmware, mctl.FirmwareFlag)
}
