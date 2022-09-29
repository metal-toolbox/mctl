package cmd

import (
	"log"
	"os"

	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
	"gopkg.in/yaml.v3"
)

// List
var cmdListFirmware = &cobra.Command{
	Use:   "firmware",
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

		firmware, _, err := c.ListServerComponentFirmware(cmd.Context(), nil)
		if err != nil {
			log.Fatal(err)
		}

		if outputJSON {
			printJSON(firmware)
			os.Exit(0)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"UUID", "Vendor", "Model", "Component", "Version"})
		for _, f := range firmware {
			table.Append([]string{f.UUID.String(), f.Vendor, f.Model, f.Component, f.Version})
		}
		table.Render()
	},
}

// Create
type createFirmwareFlags struct {
	// file containing firmware configuration
	firmwareConfigFile string
}

var (
	flagsDefinedCreateFirmware *createFirmwareFlags
)

var cmdCreateFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "Create firmware",
	Run: func(cmd *cobra.Command, args []string) {
		mctl, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			log.Fatal(err)
		}

		client, err := newServerserviceClient(cmd.Context(), mctl)
		if err != nil {
			log.Fatal(err)
		}

		firmwareConfig := &model.FirmwareConfig{}
		fbytes, err := os.ReadFile(flagsDefinedCreateFirmware.firmwareConfigFile)
		if err != nil {
			log.Fatal(err)
		}

		if err = yaml.Unmarshal(fbytes, firmwareConfig); err != nil {
			log.Fatal(err)
		}

		for _, provider := range firmwareConfig.Providers {
			for _, config := range provider.Firmwares {
				c := serverservice.ComponentFirmwareVersion{
					Vendor:        provider.Vendor,
					RepositoryURL: provider.RepositoryURL,
					Model:         config.Model,
					UpstreamURL:   config.UpstreamURL,
					Version:       config.Version,
					Filename:      config.FileName,
					Checksum:      config.Checksum,
					Component:     config.ComponentSlug,
				}

				id, _, err := client.CreateServerComponentFirmware(cmd.Context(), c)
				if err != nil {
					log.Fatal(err)
				}

				log.Println(id)
			}
		}
	},
}

func init() {
	flagsDefinedCreateFirmware = &createFirmwareFlags{}

	cmdCreateFirmware.PersistentFlags().StringVar(&flagsDefinedCreateFirmware.firmwareConfigFile, "from-file", "", "YAML file with firmware configuration data")

	if err := cmdCreateFirmware.MarkPersistentFlagRequired("from-file"); err != nil {
		log.Fatal(err)
	}
}
