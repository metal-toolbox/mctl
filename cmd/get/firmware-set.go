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

type getFirmwareSetFlags struct {
	id string
}

var (
	flagsDefinedGetFirmwareSet *getFirmwareSetFlags
)

// Get firmware set
var getFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Get information for given firmware set identifier",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		fwsID, err := uuid.Parse(flagsDefinedGetFirmwareSet.id)
		if err != nil {
			log.Fatal(err)
		}

		firmwareSet, _, err := client.GetServerComponentFirmwareSet(ctx, fwsID)
		if err != nil {
			log.Fatal("serverservice client returned error: ", err)
		}

		writeResults(firmwareSet)
		os.Exit(0)
	},
}

func init() {
	flagsDefinedGetFirmwareSet = &getFirmwareSetFlags{}

	getFirmwareSet.PersistentFlags().StringVar(&flagsDefinedGetFirmwareSet.id, "id", "", "firmware set UUID")

	if err := getFirmwareSet.MarkPersistentFlagRequired("id"); err != nil {
		log.Fatal(err)
	}
}
