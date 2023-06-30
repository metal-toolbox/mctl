package create

import (
	"log"
	"os"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/spf13/cobra"
	ss "go.hollow.sh/serverservice/pkg/api/v1"
	"gopkg.in/yaml.v2"
)

// Create
type createFirmwareFlags struct {
	// file containing firmware configuration
	firmwareConfigFile string
}

var (
	flagsDefinedCreateFirmware *createFirmwareFlags
)

var createFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "Create firmware",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewServerserviceClient(cmd.Context(), theApp)
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

		for _, config := range firmwareConfig.Firmwares {
			c := ss.ComponentFirmwareVersion{
				Vendor:        config.Vendor,
				RepositoryURL: config.RepositoryURL,
				Model:         config.Model,
				UpstreamURL:   config.UpstreamURL,
				Version:       config.Version,
				Filename:      config.FileName,
				Checksum:      config.Checksum,
				Component:     config.Component,
			}

			id, _, err := client.CreateServerComponentFirmware(cmd.Context(), c)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(id)
		}
	},
}

func init() {
	flagsDefinedCreateFirmware = &createFirmwareFlags{}

	createFirmware.PersistentFlags().StringVar(
		&flagsDefinedCreateFirmware.firmwareConfigFile,
		"from-file", "", "YAML file with firmware configuration data")
	if err := createFirmware.MarkPersistentFlagRequired("from-file"); err != nil {
		log.Fatal(err)
	}
}
