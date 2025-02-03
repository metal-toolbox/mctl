package install

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	rctypes "github.com/metal-toolbox/rivets/v2/condition"
)

type installFirmwareSetFlags struct {
	firmwareSetID         string
	serverID              string
	forceInstall          bool
	skipBMCReset          bool
	requireHostPoweredOff bool
	dryRun                bool
}

var flagsDefinedInstallFwSet *installFirmwareSetFlags

// List
var installFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Install firmware set",
	Run: func(cmd *cobra.Command, _ []string) {
		installFwSet(cmd.Context())
	},
}

func installFwSet(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	serverID, err := uuid.Parse(flagsDefinedInstallFwSet.serverID)
	if err != nil {
		log.Fatal(err)
	}

	ssclient, err := app.NewFleetDBAPIClient(ctx, theApp.Config.FleetDBAPI, theApp.Reauth)
	if err != nil {
		log.Fatal(errors.Wrap(err, "fleetdb API client init error"))
	}

	fwSetID, err := firmwareSetForInstall(ctx, ssclient, serverID)
	if err != nil {
		log.Fatal(err)
	}

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	params := &rctypes.FirmwareInstallTaskParameters{
		AssetID:               serverID,
		FirmwareSetID:         fwSetID,
		ResetBMCBeforeInstall: !flagsDefinedInstallFwSet.skipBMCReset,
		ForceInstall:          flagsDefinedInstallFwSet.forceInstall,
		DryRun:                flagsDefinedInstallFwSet.dryRun,
		RequireHostPoweredOff: flagsDefinedInstallFwSet.requireHostPoweredOff,
	}

	response, err := client.ServerFirmwareInstall(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	condition, err := mctl.ConditionFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%d msg=%s conditionID=%s", response.StatusCode, response.Message, condition.ID)
}

func firmwareSetForInstall(ctx context.Context, client *fleetdbapi.Client, serverID uuid.UUID) (fwSetID uuid.UUID, err error) {
	errInvalidFwSetID := errors.New("invalid firmware set ID")
	errNoVendorAttrs := errors.New("unable to determine server vendor, model attributes")

	// validate server exists
	server, _, err := client.Get(ctx, serverID)
	if err != nil {
		if strings.Contains(err.Error(), "resource not found") {
			return uuid.Nil, errors.Wrap(err, "invalid server ID")
		}

		return uuid.Nil, errors.Wrap(err, "failed to retrieve server object")
	}

	// if a firmware set identifier was given, validate and return
	if flagsDefinedInstallFwSet.firmwareSetID != "" {
		fwSetID, err = uuid.Parse(flagsDefinedInstallFwSet.firmwareSetID)
		if err != nil {
			return uuid.Nil, errors.Wrap(errInvalidFwSetID, err.Error())
		}

		_, _, err = client.GetServerComponentFirmwareSet(ctx, fwSetID)
		if err != nil {
			return uuid.Nil, errors.Wrap(errInvalidFwSetID, err.Error())
		}

		return fwSetID, nil
	}

	// identify vendor, model attributes
	vendor, model := mctl.VendorModelFromAttrs(server.Attributes)
	if vendor == "" || model == "" {
		return uuid.Nil, errors.Wrap(errNoVendorAttrs, "specify a firmware set ID with --id instead")
	}

	// identify firmware set by vendor, model attributes
	fwSetID, err = mctl.FirmwareSetIDByVendorModel(ctx, vendor, model, client)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "specify a firmware set ID with --id instead")
	}

	return fwSetID, nil
}

func init() {
	flagsDefinedInstallFwSet = &installFirmwareSetFlags{}

	mctl.AddServerFlag(installFirmwareSet, &flagsDefinedInstallFwSet.serverID)
	mctl.AddFirmwareSetFlag(installFirmwareSet, &flagsDefinedInstallFwSet.firmwareSetID)
	mctl.AddForceFlag(installFirmwareSet, &flagsDefinedInstallFwSet.forceInstall,
		"force install (skips firmware version check)")
	mctl.AddDryRunFlag(installFirmwareSet, &flagsDefinedInstallFwSet.dryRun,
		"Run install process in dry-run (skips firmware install)")
	mctl.AddSkipBmcResetFlag(installFirmwareSet, &flagsDefinedInstallFwSet.skipBMCReset)
	mctl.AddPowerOffRequiredFlag(installFirmwareSet, &flagsDefinedInstallFwSet.requireHostPoweredOff,
		"require host to be powered off before proceeding install")

	mctl.RequireFlag(installFirmwareSet, mctl.ServerFlag)
}
