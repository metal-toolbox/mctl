package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	definedfirmwareSetFlags *FirmwareSetFlags
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
		table.SetHeader([]string{"UUID", "Name", "Labels", "firmware UUID", "Vendor", "Model", "Component", "Version"})
		for _, s := range set {
			var labels string
			if len(s.Attributes) > 0 {
				attr := findAttribute(model.AttributeNSFirmwareSetLabels, s.Attributes)
				if attr != nil {
					labels = string(attr.Data)
				}
			}
			table.Append([]string{s.UUID.String(), s.Name, labels, "-", "-", "-", "-", "-"})
			for _, f := range s.ComponentFirmware {
				table.Append([]string{s.UUID.String(), "", "", f.UUID.String(), f.Vendor, strings.Join(f.Model, ","), f.Component, f.Version})
			}
		}

		table.SetAutoMergeCells(true)
		table.Render()
	},
}

// Delete

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

		id, err := uuid.Parse(definedfirmwareSetFlags.ID)
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

		id, err := uuid.Parse(definedfirmwareSetFlags.ID)
		if err != nil {
			log.Fatal(err)
		}

		payload := serverservice.ComponentFirmwareSetRequest{
			ID:                     id,
			ComponentFirmwareUUIDs: []string{},
		}

		var attrs *serverservice.Attributes
		if len(definedfirmwareSetFlags.Labels) > 0 {
			attrs, err = AttributeFromLabels(model.AttributeNSFirmwareSetLabels, definedfirmwareSetFlags.Labels)
			if err != nil {
				log.Fatal(err)
			}

			payload.Attributes = []serverservice.Attributes{*attrs}

			_, err = c.UpdateComponentFirmwareSetRequest(cmd.Context(), id, payload)
			if err != nil {
				log.Fatal(err)
			}
		}

		if len(payload.ComponentFirmwareUUIDs) > 0 {
			for _, id := range strings.Split(definedfirmwareSetFlags.FirmwareUUIDs, ",") {
				_, err = uuid.Parse(id)
				if err != nil {
					log.Println(err.Error())
					os.Exit(1)
				}

				payload.ComponentFirmwareUUIDs = append(payload.ComponentFirmwareUUIDs, id)
			}

			_, err = c.RemoveServerComponentFirmwareSetFirmware(cmd.Context(), id, payload)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("firmware set updated: " + id.String())
	},
}

func init() {
	definedfirmwareSetFlags = &FirmwareSetFlags{}
	// delete

	cmdDeleteFirmwareSet.PersistentFlags().StringVar(&definedfirmwareSetFlags.ID, "uuid", "", "UUID of firmware set to be deleted")

	if err := cmdDeleteFirmwareSet.MarkPersistentFlagRequired("uuid"); err != nil {
		log.Fatal(err)
	}

	// edit
	cmdEditFirmwareSet.PersistentFlags().StringVar(&definedfirmwareSetFlags.ID, "uuid", "", "UUID of firmware set to be edited")
	cmdEditFirmwareSet.PersistentFlags().StringVar(&definedfirmwareSetFlags.FirmwareSetName, "name", "", "Update name for the firmware set")
	cmdEditFirmwareSet.PersistentFlags().StringToStringVar(&definedfirmwareSetFlags.Labels, "labels", nil, "Labels to assign to the firmware set - 'vendor=foo,model=bar'")

	if err := cmdEditFirmwareSet.MarkPersistentFlagRequired("uuid"); err != nil {
		log.Fatal(err)
	}

	cmdEditFirmwareSet.PersistentFlags().StringVar(&definedfirmwareSetFlags.FirmwareUUIDs, "remove-firmware-uuids", "", "UUIDs of firmware to be removed from the set")
}
