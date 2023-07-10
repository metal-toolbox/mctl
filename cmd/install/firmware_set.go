package install

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	cotypesv1 "github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
	cotypes "github.com/metal-toolbox/conditionorc/pkg/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

type installFirmwareSetFlags struct {
	firmwareSetID string
	serverID      string
	forceInstall  bool
	skipBMCReset  bool
}

var (
	flagsDefined *installFirmwareSetFlags
)

// List
var installFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Install firmware set",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		fwSetID, err := uuid.Parse(flagsDefined.firmwareSetID)
		if err != nil {
			log.Fatal(err)
		}

		ssc, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		_, _, err = ssc.GetServerComponentFirmwareSet(cmd.Context(), fwSetID)
		if err != nil {
			log.Fatal(err)
		}

		client, err := app.NewConditionsClient(cmd.Context(), theApp.Config.Conditions, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		serverID, err := uuid.Parse(flagsDefined.serverID)
		if err != nil {
			log.Fatal(err)
		}

		b, _ := json.Marshal(parameters{
			AssetID:               serverID,
			FirmwareSetID:         fwSetID,
			ResetBMCBeforeInstall: !flagsDefined.skipBMCReset,
			ForceInstall:          flagsDefined.forceInstall,
		})

		co := cotypesv1.ConditionCreate{
			Exclusive:  true,
			Parameters: json.RawMessage(b),
		}

		response, err := client.ServerConditionCreate(cmd.Context(), serverID, cotypes.FirmwareInstall, co)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(response)
	},
}

func init() {
	flagsDefined = &installFirmwareSetFlags{}

	install.AddCommand(installFirmwareSet)
	installFirmwareSet.PersistentFlags().StringVar(&flagsDefined.serverID, "server", "", "server UUID")
	installFirmwareSet.PersistentFlags().StringVar(&flagsDefined.firmwareSetID, "id", "", "firmware set UUID")
	installFirmwareSet.PersistentFlags().BoolVar(&flagsDefined.forceInstall, "force", false, "force install (skips firmware version check)")
	installFirmwareSet.PersistentFlags().BoolVar(&flagsDefined.skipBMCReset, "skip-bmc-reset", false, "skip BMC reset before firmware install")

	if err := installFirmwareSet.MarkPersistentFlagRequired("server"); err != nil {
		log.Fatal(err)
	}

	if err := installFirmwareSet.MarkPersistentFlagRequired("id"); err != nil {
		log.Fatal(err)
	}
}
