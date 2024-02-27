package create

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
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

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		var firmwares []*serverservice.ComponentFirmwareVersion
		fbytes, err := os.ReadFile(flagsDefinedCreateFirmware.firmwareConfigFile)
		if err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(fbytes, &firmwares); err != nil {
			log.Fatal(err)
		}

		for _, fw := range firmwares {
			id, _, err := client.CreateServerComponentFirmware(cmd.Context(), *fw)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(id)
		}
	},
}

func init() {
	flagsDefinedCreateFirmware = &createFirmwareFlags{}
	usage := "JSON file with firmware configuration data"

	mctl.AddFromFileFlag(createFirmware, &flagsDefinedCreateFirmware.firmwareConfigFile, usage)
	mctl.RequireFlag(createFirmware, mctl.FromFileFlag)
}
