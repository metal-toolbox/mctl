package list

import (
	"log"
	"os"
	"strings"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

type listFirmwareSetFlags struct {
	vendor  string
	model   string
	listAll bool
}

var (
	flags *listFirmwareSetFlags
)

// List
var listFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "List firmware",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		var fwSet []serverservice.ComponentFirmwareSet
		if flags.listAll {
			params := &serverservice.ComponentFirmwareSetListParams{}
			if flags.vendor != "" {
				params.AttributeListParams = append(params.AttributeListParams, serverservice.AttributeListParams{
					Namespace: mctl.FirmwareSetAttributeNS,
					Keys:      []string{"vendor"},
					Operator:  "eq",
					Value:     strings.ToLower(flags.vendor),
				})
			}
			if flags.model != "" {
				params.AttributeListParams = append(params.AttributeListParams, serverservice.AttributeListParams{
					Namespace: mctl.FirmwareSetAttributeNS,
					Keys:      []string{"model"},
					Operator:  "eq",
					Value:     strings.ToLower(flags.model),
				})
			}
			fwSet, _, err = client.ListServerComponentFirmwareSet(cmd.Context(), params)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fwSet, err = mctl.FirmwareSetByVendorModel(cmd.Context(), flags.vendor, flags.model, client)
			if err != nil {
				log.Fatal(err)
			}
		}

		if outputJSON {
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

	listFirmwareSet.PersistentFlags().StringVar(&flags.vendor, "vendor", "", "filter by server vendor")
	listFirmwareSet.PersistentFlags().StringVar(&flags.model, "model", "", "filter by server model")
	listFirmwareSet.PersistentFlags().BoolVar(&flags.listAll, "all", false, "show all firmware sets. By default results are filtered on having labels for vendor, model and latest=true")
}
