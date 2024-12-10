package list

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
)

type listFirmwareSetFlags struct {
	vendor string
	model  string
	labels map[string]string
}

var (
	flags *listFirmwareSetFlags
)

func sendListFirmwareRequest(client *fleetdbapi.Client, cmd *cobra.Command) ([]fleetdbapi.ComponentFirmwareSet, error) {
	params := &fleetdbapi.ComponentFirmwareSetListParams{
		Vendor: strings.TrimSpace(flags.vendor),
		Model:  strings.TrimSpace(flags.model),
	}

	labelParts := make([]string, 0)
	for k, v := range flags.labels {
		labelParts = append(labelParts, fmt.Sprintf("%s=%s", k, v))
	}
	params.Labels = strings.Join(labelParts, ",")

	fwSet, _, err := client.ListServerComponentFirmwareSet(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("retrieving firmware sets: %w", err)
	}

	if len(fwSet) == 0 {
		return nil, errors.New("no fw sets identified")
	}

	return fwSet, nil
}

// List
var listFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "List firmware",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		fwSet, err := sendListFirmwareRequest(client, cmd)
		if err != nil {
			log.Fatal(err)
		}

		if output == mctl.OutputTypeJSON.String() {
			printJSON(fwSet)
			os.Exit(0)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"UUID", "Name", "Labels", "firmware UUID", "Vendor", "Model", "Component", "Version"})
		for _, s := range fwSet {
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

func init() {
	flags = &listFirmwareSetFlags{}

	mctl.AddModelFlag(listFirmwareSet, &flags.model)
	mctl.AddVendorFlag(listFirmwareSet, &flags.vendor)
	mctl.AddLabelsFlag(listFirmwareSet, &flags.labels,
		"Labels to identify the firmware set - e.g. 'key=value,default=true,latest=true'")
	mctl.RequireFlag(listFirmwareSet, mctl.VendorFlag)
	mctl.RequireFlag(listFirmwareSet, mctl.ModelFlag)
}
