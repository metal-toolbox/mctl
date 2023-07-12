package install

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/google/uuid"
	cotypesv1 "github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
	cotypes "github.com/metal-toolbox/conditionorc/pkg/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

type installFirmwareSetFlags struct {
	firmwareSetID string
	serverID      string
	forceInstall  bool
	skipBMCReset  bool
}

var (
	flagsDefinedInstallFwSet *installFirmwareSetFlags
)

// List
var installFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Install firmware set",
	Run: func(cmd *cobra.Command, args []string) {
		installFwSet(cmd.Context())

	},
}

func installFwSet(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	serverID, err := uuid.Parse(flagsDefinedInstallFwSet.serverID)
	if err != nil {
		log.Fatal(err)
	}

	ssclient, err := app.NewServerserviceClient(ctx, theApp.Config.Serverservice, theApp.Reauth)
	if err != nil {
		log.Fatal(errors.Wrap(err, "serverservice client init error"))
	}

	fwSetID, err := firmwareSetForInstall(ctx, ssclient, serverID)
	if err != nil {
		log.Fatal(err)
	}

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(parameters{
		AssetID:               serverID,
		FirmwareSetID:         fwSetID,
		ResetBMCBeforeInstall: !flagsDefinedInstallFwSet.skipBMCReset,
		ForceInstall:          flagsDefinedInstallFwSet.forceInstall,
	})

	co := cotypesv1.ConditionCreate{
		Exclusive:  true,
		Parameters: json.RawMessage(b),
	}

	response, err := client.ServerConditionCreate(ctx, serverID, cotypes.FirmwareInstall, co)
	if err != nil {
		log.Fatal(err)
	}

	condition, err := conditionResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%d msg=%s conditionID=%s", response.StatusCode, response.Message, condition.ID)
}

func firmwareSetForInstall(ctx context.Context, client *serverservice.Client, serverID uuid.UUID) (fwSetID uuid.UUID, err error) {
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

func conditionResponse(response *cotypesv1.ServerResponse) (cotypes.Condition, error) {
	errUnexpectedResponse := errors.New("unexpected response")

	if response.Records == nil || len(response.Records.Conditions) == 0 {
		return cotypes.Condition{}, errors.Wrap(errUnexpectedResponse, "empty records")
	}

	return *response.Records.Conditions[0], nil
}

func init() {
	flagsDefinedInstallFwSet = &installFirmwareSetFlags{}

	install.AddCommand(installFirmwareSet)
	installFirmwareSet.PersistentFlags().StringVar(&flagsDefinedInstallFwSet.serverID, "server", "", "server UUID")
	installFirmwareSet.PersistentFlags().StringVar(&flagsDefinedInstallFwSet.firmwareSetID, "id", "", "firmware set UUID")
	installFirmwareSet.PersistentFlags().BoolVar(&flagsDefinedInstallFwSet.forceInstall, "force", false, "force install (skips firmware version check)")
	installFirmwareSet.PersistentFlags().BoolVar(&flagsDefinedInstallFwSet.skipBMCReset, "skip-bmc-reset", false, "skip BMC reset before firmware install")

	if err := installFirmwareSet.MarkPersistentFlagRequired("server"); err != nil {
		log.Fatal(err)
	}
}
