package get

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

type getFirmwareFlags struct {
	id string
}

var (
	flagsDefinedGetFirmware *getFirmwareFlags
)

// Get firmware info
var getFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "Get information for given firmware identifier",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		fwID, err := uuid.Parse(flagsDefinedGetFirmware.id)
		if err != nil {
			log.Fatal(err)
		}

		firmware, _, err := client.GetServerComponentFirmware(ctx, fwID)
		if err != nil {
			log.Fatal("serverservice client returned error: ", err)
		}

		writeResults(firmware)
		os.Exit(0)
	},
}

func init() {
	flagsDefinedGetFirmware = &getFirmwareFlags{}

	getFirmware.PersistentFlags().StringVarP(&flagsDefinedGetFirmware.id, "firmware-id", "f", "", "firmware UUID")

	if err := getFirmware.MarkPersistentFlagRequired("firmware-id"); err != nil {
		log.Fatal(err)
	}
}
