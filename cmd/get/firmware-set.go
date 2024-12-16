package get

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
)

type getFirmwareSetFlags struct {
	id       string
	serverID string
}

var (
	flagsDefinedGetFirmwareSet *getFirmwareSetFlags
)

// Get firmware set
var getFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Get information for given firmware set identifier",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		if flagsDefinedGetFirmwareSet.id == "" && flagsDefinedGetFirmwareSet.serverID == "" {
			_ = cmd.Help()
			os.Exit(1)
		}

		var firmwareSet *fleetdbapi.ComponentFirmwareSet

		if flagsDefinedGetFirmwareSet.serverID != "" {
			firmwareSet, err = firmwareSetForServer(ctx, client, flagsDefinedGetFirmwareSet.serverID)
			if err != nil {
				log.Fatal("fleetdb API client returned error: ", err)
			}
		} else {
			fwsID, err := uuid.Parse(flagsDefinedGetFirmwareSet.id)
			if err != nil {
				log.Fatal(err)
			}

			firmwareSet, _, err = client.GetServerComponentFirmwareSet(ctx, fwsID)
			if err != nil {
				log.Fatal("fleetdb API client returned error: ", err)
			}
		}

		mctl.PrintResults(output, firmwareSet)
		os.Exit(0)
	},
}

func firmwareSetForServer(ctx context.Context, client *fleetdbapi.Client, serverID string) (*fleetdbapi.ComponentFirmwareSet, error) {
	errNoVendorAttrs := errors.New("unable to determine server vendor, model attributes")
	errNotFound := errors.New("no firmware sets identified for server")

	serverUUID, err := uuid.Parse(serverID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid server ID")
	}

	// validate server exists
	server, _, err := client.Get(ctx, serverUUID)
	if err != nil {
		if strings.Contains(err.Error(), "resource not found") {
			return nil, errors.Wrap(err, "invalid server ID")
		}

		return nil, errors.Wrap(err, "failed to retrieve server object")
	}

	// identify vendor, model attributes
	vendor, model := mctl.VendorModelFromAttrs(server.Attributes)
	if vendor == "" || model == "" {
		return nil, errNoVendorAttrs
	}

	params := &fleetdbapi.ComponentFirmwareSetListParams{
		Vendor: strings.TrimSpace(vendor),
		Model:  strings.TrimSpace(model),
		Labels: "default=true,latest=true",
	}

	// identify firmware set by vendor, model attributes
	fwSet, _, err := client.ListServerComponentFirmwareSet(ctx, params)
	if err != nil {
		return nil, err
	}

	if fwSet == nil {
		return nil, errors.Wrap(
			errNotFound,
			fmt.Sprintf("vendor: %s, model: %s", vendor, model),
		)
	}

	return &fwSet[0], nil
}

func init() {
	flagsDefinedGetFirmwareSet = &getFirmwareSetFlags{}

	mctl.AddServerFlag(getFirmwareSet, &flagsDefinedGetFirmwareSet.serverID)
	mctl.AddFirmwareSetFlag(getFirmwareSet, &flagsDefinedGetFirmwareSet.id)

	mctl.MutuallyExclusiveFlags(getFirmwareSet, mctl.ServerFlag, mctl.FirmwareSetFlag)
}
