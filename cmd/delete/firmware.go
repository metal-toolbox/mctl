package deleteresource

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"

	ss "go.hollow.sh/serverservice/pkg/api/v1"
)

var deleteFW string

var deleteFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "Delete a firmware",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(deleteFW)
		if err != nil {
			log.Fatal(err)
		}

		cfv := ss.ComponentFirmwareVersion{
			UUID: id,
		}

		_, err = client.DeleteServerComponentFirmware(cmd.Context(), cfv)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("firmware deleted: " + id.String())
	},
}

func init() {
	deleteFirmware.PersistentFlags().StringVar(&deleteFW, "uuid", "", "UUID of firmware to be deleted")

	if err := deleteFirmware.MarkPersistentFlagRequired("uuid"); err != nil {
		log.Fatal(err)
	}
}
