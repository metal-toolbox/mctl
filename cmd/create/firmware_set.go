package create

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
)

var (
	definedfirmwareSetFlags *mctl.FirmwareSetFlags
	errFwSetUUIDs           = errors.New("one or more firmware UUIDs required to create set")
)

var createFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Create a firmware set",
	PreRun: func(cmd *cobra.Command, _ []string) {
		// set required flags if from-file flag is not passed in
		fromFile, err := cmd.Flags().GetString(mctl.FromFileFlag.Name())
		if err != nil {
			log.Fatal(err)
		}

		if fromFile == "" {
			mctl.RequireFlag(cmd, mctl.FirmwareIDsFlag)
			mctl.RequireFlag(cmd, mctl.NameFlag)
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		if definedfirmwareSetFlags.CreateFromFile != "" {
			err = createFWSetsFromFile(cmd.Context(), client, definedfirmwareSetFlags)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = createFWSetFromCLI(cmd.Context(), client, definedfirmwareSetFlags)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func createFWSetsFromFile(ctx context.Context, client *fleetdbapi.Client, flgs *mctl.FirmwareSetFlags) (err error) {
	var fwsets []*fleetdbapi.ComponentFirmwareSet

	fbytes, err := os.ReadFile(flgs.CreateFromFile)
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(fbytes, &fwsets); err != nil {
		log.Fatal(err)
	}

	firmwareAdded := map[string]bool{}
	for _, set := range fwsets {
		if len(set.ComponentFirmware) == 0 {
			continue
		}

		// create firmware
		setFwUUIDs := []string{}

		// nolint:gocritic // this is fine
		for _, fw := range set.ComponentFirmware {
			setFwUUIDs = append(setFwUUIDs, fw.UUID.String())

			_, exists := firmwareAdded[fw.UUID.String()]
			if exists {
				continue
			}

			log.Printf("Adding firmware object: " + fw.UUID.String())
			_, _, err = client.CreateServerComponentFirmware(ctx, fw)
			if err != nil {
				log.Fatal("error adding firmware object: ", err)
			}

			firmwareAdded[fw.UUID.String()] = true
		}

		log.Printf("Adding firmware-set object: " + set.UUID.String())
		req := fleetdbapi.ComponentFirmwareSetRequest{
			ID:                     set.UUID,
			Attributes:             set.Attributes,
			Name:                   set.Name,
			ComponentFirmwareUUIDs: setFwUUIDs,
		}

		_, _, err = client.CreateServerComponentFirmwareSet(ctx, req)
		if err != nil {
			log.Fatal("error adding firmware-set object: ", err)
		}
	}

	return nil
}

func createFWSetFromCLI(ctx context.Context, client *fleetdbapi.Client, flgs *mctl.FirmwareSetFlags) (err error) {
	payload := fleetdbapi.ComponentFirmwareSetRequest{
		Name:                   flgs.FirmwareSetName,
		ComponentFirmwareUUIDs: []string{},
	}

	if len(definedfirmwareSetFlags.Labels) > 0 {
		var attrs *fleetdbapi.Attributes
		attrs, err = mctl.AttributeFromLabels(model.AttributeNSFirmwareSetLabels, flgs.Labels)
		if err != nil {
			return err
		}

		payload.Attributes = []fleetdbapi.Attributes{*attrs}
	}

	for _, id := range flgs.AddFirmwareUUIDs {
		_, err = uuid.Parse(id)
		if err != nil {
			return err
		}

		payload.ComponentFirmwareUUIDs = append(payload.ComponentFirmwareUUIDs, id)
	}

	if len(payload.ComponentFirmwareUUIDs) == 0 {
		return errFwSetUUIDs
	}

	id, _, err := client.CreateServerComponentFirmwareSet(ctx, payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(id)

	return nil
}

func init() {
	definedfirmwareSetFlags = &mctl.FirmwareSetFlags{}

	mctl.AddFirmwareIDsFlag(createFirmwareSet, &definedfirmwareSetFlags.AddFirmwareUUIDs)
	mctl.AddNameFlag(createFirmwareSet, &definedfirmwareSetFlags.FirmwareSetName, "A name for the firmware set")
	mctl.AddLabelsFlag(createFirmwareSet, &definedfirmwareSetFlags.Labels,
		"Labels to assign to the firmware set - 'vendor=foo,model=bar'")

	mctl.AddFromFileFlag(createFirmwareSet, &definedfirmwareSetFlags.CreateFromFile, "JSON file with firmware configuration data")
}
