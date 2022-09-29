package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

// List
var cmdListFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "List firmware",
	Run: func(cmd *cobra.Command, args []string) {
		mctl, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			log.Fatal(err)
		}

		c, err := newServerserviceClient(cmd.Context(), mctl)
		if err != nil {
			log.Fatal(err)
		}

		set, _, err := c.ListServerComponentFirmwareSet(cmd.Context(), nil)
		if err != nil {
			log.Fatal(err)
		}

		if outputJSON {
			printJSON(set)
			os.Exit(0)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"UUID", "Name", "Metadata", "firmware UUID", "Vendor", "Model", "Component", "Version"})
		for _, s := range set {
			table.Append([]string{s.UUID.String(), s.Name, string(s.Metadata), "-", "-", "-", "-", "-"})
			for _, f := range s.ComponentFirmware {
				table.Append([]string{s.UUID.String(), "", "", f.UUID.String(), f.Vendor, f.Model, f.Component, f.Version})
			}
		}

		table.SetAutoMergeCells(true)
		table.Render()
	},
}

// Create
type createFirmwareSetFlags struct {
	// comma separated list of firmware UUIDs
	firmwareUUIDs string
	// name for the firmware set to be created
	firmwareSetName string
}

var (
	flagsDefinedCreateFirmwareSet *createFirmwareSetFlags
)

var cmdCreateFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Create a firmware set",
	Run: func(cmd *cobra.Command, args []string) {
		mctl, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			log.Fatal(err)
		}

		c, err := newServerserviceClient(cmd.Context(), mctl)
		if err != nil {
			log.Fatal(err)
		}

		payload := serverservice.ComponentFirmwareSetRequest{
			Name:                   flagsDefinedCreateFirmwareSet.firmwareSetName,
			ComponentFirmwareUUIDs: []string{},
		}

		for _, id := range strings.Split(flagsDefinedCreateFirmwareSet.firmwareUUIDs, ",") {
			_, err = uuid.Parse(id)
			if err != nil {
				log.Println(err.Error())
				os.Exit(1)
			}

			payload.ComponentFirmwareUUIDs = append(payload.ComponentFirmwareUUIDs, id)
		}

		if len(payload.ComponentFirmwareUUIDs) == 0 {
			log.Fatal("one or more firmware UUIDs required to create set")
		}

		id, _, err := c.CreateServerComponentFirmwareSet(cmd.Context(), payload)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id)
	},
}

// Delete
type deleteFirmwareSetFlags struct {
	id string
}

var (
	flagsDefinedDeleteFirmwareSet *deleteFirmwareSetFlags
)

var cmdDeleteFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Delete a firmware set",
	Run: func(cmd *cobra.Command, args []string) {
		mctl, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			log.Fatal(err)
		}

		c, err := newServerserviceClient(cmd.Context(), mctl)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(flagsDefinedDeleteFirmwareSet.id)
		if err != nil {
			log.Fatal(err)
		}

		_, err = c.DeleteServerComponentFirmwareSet(cmd.Context(), id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("firmware set deleted: " + id.String())
	},
}

// Edit

type editFirmwareSetFlags struct {
	// firmware set UUID
	id string
	// comma separated list of firmware UUIDs to add to firmware set
	// addFirmwareUUIDs string
	// comma separated list of firmware UUIDs to remove from a firmware set
	removeFirmwareUUIDs string
	// set a new name for the firmware set
	// newName string
}

var (
	flagsDefinedEditFirmwareSet *editFirmwareSetFlags
)

var cmdEditFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Edit a firmware set",
	Run: func(cmd *cobra.Command, args []string) {
		mctl, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			log.Fatal(err)
		}

		c, err := newServerserviceClient(cmd.Context(), mctl)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(flagsDefinedEditFirmwareSet.id)
		if err != nil {
			log.Fatal(err)
		}

		payload := serverservice.ComponentFirmwareSetRequest{
			ID:                     id,
			ComponentFirmwareUUIDs: []string{},
		}

		for _, id := range strings.Split(flagsDefinedEditFirmwareSet.removeFirmwareUUIDs, ",") {
			_, err = uuid.Parse(id)
			if err != nil {
				log.Println(err.Error())
				os.Exit(1)
			}

			payload.ComponentFirmwareUUIDs = append(payload.ComponentFirmwareUUIDs, id)
		}

		if len(payload.ComponentFirmwareUUIDs) == 0 {
			log.Fatal("one or more firmware UUIDs required to create set")
		}

		_, err = c.RemoveServerComponentFirmwareSetFirmware(cmd.Context(), id, payload)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("firmware set updated: " + id.String())
	},
}

func init() {
	// create
	flagsDefinedCreateFirmwareSet = &createFirmwareSetFlags{}
	cmdCreateFirmwareSet.PersistentFlags().StringVar(&flagsDefinedCreateFirmwareSet.firmwareUUIDs, "firmware-uuids", "", "comma separated list of UUIDs of firmware to be included in the set to be created")
	cmdCreateFirmwareSet.PersistentFlags().StringVar(&flagsDefinedCreateFirmwareSet.firmwareSetName, "name", "", "A name for the firmware set")

	// mark flags as required
	if err := cmdCreateFirmwareSet.MarkPersistentFlagRequired("firmware-uuids"); err != nil {
		log.Fatal(err)
	}

	if err := cmdCreateFirmwareSet.MarkPersistentFlagRequired("name"); err != nil {
		log.Fatal(err)
	}

	// delete
	flagsDefinedDeleteFirmwareSet = &deleteFirmwareSetFlags{}

	cmdDeleteFirmwareSet.PersistentFlags().StringVar(&flagsDefinedDeleteFirmwareSet.id, "uuid", "", "UUID of firmware set to be deleted")

	if err := cmdDeleteFirmwareSet.MarkPersistentFlagRequired("uuid"); err != nil {
		log.Fatal(err)
	}

	// edit
	flagsDefinedEditFirmwareSet = &editFirmwareSetFlags{}

	cmdEditFirmwareSet.PersistentFlags().StringVar(&flagsDefinedEditFirmwareSet.id, "uuid", "", "UUID of firmware set to be deleted")

	if err := cmdEditFirmwareSet.MarkPersistentFlagRequired("uuid"); err != nil {
		log.Fatal(err)
	}

	cmdEditFirmwareSet.PersistentFlags().StringVar(&flagsDefinedEditFirmwareSet.removeFirmwareUUIDs, "remove-firmware-uuids", "", "UUIDs of firmware to be removed from the set")

	if err := cmdEditFirmwareSet.MarkPersistentFlagRequired("remove-firmware-uuids"); err != nil {
		log.Fatal(err)
	}
}
