package create

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/spf13/cobra"
	ss "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	definedfirmwareSetFlags *mctl.FirmwareSetFlags
)

var createFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Create a firmware set",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		payload := ss.ComponentFirmwareSetRequest{
			Name:                   definedfirmwareSetFlags.FirmwareSetName,
			ComponentFirmwareUUIDs: []string{},
		}

		var attrs *ss.Attributes
		if len(definedfirmwareSetFlags.Labels) > 0 {
			attrs, err = mctl.AttributeFromLabels(model.AttributeNSFirmwareSetLabels, definedfirmwareSetFlags.Labels)
			if err != nil {
				log.Fatal(err)
			}

			payload.Attributes = []ss.Attributes{*attrs}
		}

		for _, id := range definedfirmwareSetFlags.AddFirmwareUUIDs {
			_, err = uuid.Parse(id)
			if err != nil {
				log.Fatal(err)
			}

			payload.ComponentFirmwareUUIDs = append(payload.ComponentFirmwareUUIDs, id)
		}

		if len(payload.ComponentFirmwareUUIDs) == 0 {
			log.Fatal("one or more firmware UUIDs required to create set")
		}

		id, _, err := client.CreateServerComponentFirmwareSet(cmd.Context(), payload)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id)
	},
}

func init() {
	definedfirmwareSetFlags = &mctl.FirmwareSetFlags{}

	createFirmwareSet.PersistentFlags().StringSliceVar(&definedfirmwareSetFlags.AddFirmwareUUIDs, "firmware-uuids", []string{}, "comma separated list of UUIDs of firmware to be included in the set to be created")
	createFirmwareSet.PersistentFlags().StringVar(&definedfirmwareSetFlags.FirmwareSetName, "name", "", "A name for the firmware set")
	createFirmwareSet.PersistentFlags().StringToStringVar(&definedfirmwareSetFlags.Labels, "labels", nil, "Labels to assign to the firmware set - 'vendor=foo,model=bar'")

	// mark flags as required
	if err := createFirmwareSet.MarkPersistentFlagRequired("firmware-uuids"); err != nil {
		log.Fatal(err)
	}

	if err := createFirmwareSet.MarkPersistentFlagRequired("name"); err != nil {
		log.Fatal(err)
	}

}
