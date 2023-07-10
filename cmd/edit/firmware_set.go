package edit

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/spf13/cobra"
	ss "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	editFWSetFlags mctl.FirmwareSetFlags
)

var editFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Edit a firmware set",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(editFWSetFlags.ID)
		if err != nil {
			log.Fatal(err)
		}

		payload := ss.ComponentFirmwareSetRequest{
			ID:                     id,
			ComponentFirmwareUUIDs: []string{},
		}

		var attrs *ss.Attributes
		if len(editFWSetFlags.Labels) > 0 {
			attrs, err = mctl.AttributeFromLabels(model.AttributeNSFirmwareSetLabels, editFWSetFlags.Labels)
			if err != nil {
				log.Fatal(err)
			}

			payload.Attributes = []ss.Attributes{*attrs}

			_, err = client.UpdateComponentFirmwareSetRequest(cmd.Context(), id, payload)
			if err != nil {
				log.Fatal(err)
			}
		}

		if len(payload.ComponentFirmwareUUIDs) > 0 {
			for _, id := range strings.Split(editFWSetFlags.FirmwareUUIDs, ",") {
				_, err = uuid.Parse(id)
				if err != nil {
					log.Println(err.Error())
					os.Exit(1)
				}

				payload.ComponentFirmwareUUIDs = append(payload.ComponentFirmwareUUIDs, id)
			}

			_, err = client.RemoveServerComponentFirmwareSetFirmware(cmd.Context(), id, payload)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("firmware set updated: " + id.String())
	},
}

func init() {
	cmdFlags := editFirmwareSet.PersistentFlags()
	cmdFlags.StringVar(&editFWSetFlags.ID, "uuid", "", "UUID of firmware set to be edited")
	cmdFlags.StringVar(&editFWSetFlags.FirmwareSetName, "name", "", "Update name for the firmware set")
	cmdFlags.StringToStringVar(&editFWSetFlags.Labels, "labels", nil, "Labels to assign to the firmware set - 'vendor=foo,model=bar'")

	if err := editFirmwareSet.MarkPersistentFlagRequired("uuid"); err != nil {
		log.Fatal(err)
	}

	cmdFlags.StringVar(&editFWSetFlags.FirmwareUUIDs, "remove-firmware-uuids", "", "UUIDs of firmware to be removed from the set")

}
