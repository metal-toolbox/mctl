package cmd

import (
	"log"

	"github.com/google/uuid"
	"github.com/metal-toolbox/conditionorc/pkg/api/v1/conditions/types"
	"github.com/spf13/cobra"

	"github.com/metal-toolbox/mctl/internal/app"
)

type validationFlags struct {
	srvIDStr, fwSetIDStr, output string
}

var (
	// incoming command line parameters
	fwvFlags = &validationFlags{}
)

var validateFirmwareCmd = &cobra.Command{
	Use:   "validate-firmware",
	Short: "validate a firmware set",
	Run: func(c *cobra.Command, _ []string) {
		theApp := MustCreateApp(c.Context())

		client, err := app.NewConditionsClient(c.Context(), theApp.Config.Conditions, theApp.Reauth)
		if err != nil {
			log.Fatalf("creating app structure: %s", err.Error())
		}

		srvID, err := uuid.Parse(fwvFlags.srvIDStr)
		if err != nil {
			log.Fatalf("parsing server id: %s", err.Error())
		}

		fwSetID, err := uuid.Parse(fwvFlags.fwSetIDStr)
		if err != nil {
			log.Fatalf("parsing firmware set id: %s", err.Error())
		}

		fvr := &types.FirmwareValidationRequest{
			ServerID:      srvID,
			FirmwareSetID: fwSetID,
		}

		resp, err := client.ValidateFirmwareSet(c.Context(), fvr)
		if err != nil {
			log.Fatalf("making validate firmware call: %s", err.Error())
		}

		PrintResults(fwvFlags.output, resp)
	},
}

func init() {
	RootCmd.AddCommand(validateFirmwareCmd)

	AddOutputFlag(validateFirmwareCmd, &fwvFlags.output)
	AddFirmwareSetFlag(validateFirmwareCmd, &fwvFlags.fwSetIDStr)
	AddServerFlag(validateFirmwareCmd, &fwvFlags.srvIDStr)
	RequireFlag(validateFirmwareCmd, ServerFlag)
	RequireFlag(validateFirmwareCmd, FirmwareSetFlag)
}
