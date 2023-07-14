package get

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
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
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		if flagsDefinedGetFirmwareSet.id == "" && flagsDefinedGetFirmwareSet.serverID == "" {
			//nolint:errcheck // returns nil
			cmd.Help()
			os.Exit(1)
		}

		var firmwareSet *serverservice.ComponentFirmwareSet

		if flagsDefinedGetFirmwareSet.serverID != "" {
			firmwareSet, err = firmwareSetForServer(ctx, client, flagsDefinedGetFirmwareSet.serverID)
			if err != nil {
				log.Fatal("serverservice client returned error: ", err)
			}
		} else {
			fwsID, err := uuid.Parse(flagsDefinedGetFirmwareSet.id)
			if err != nil {
				log.Fatal(err)
			}

			firmwareSet, _, err = client.GetServerComponentFirmwareSet(ctx, fwsID)
			if err != nil {
				log.Fatal("serverservice client returned error: ", err)
			}
		}

		writeResults(firmwareSet)
		os.Exit(0)
	},
}

func firmwareSetForServer(ctx context.Context, client *serverservice.Client, serverID string) (*serverservice.ComponentFirmwareSet, error) {
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

	// identify firmware set by vendor, model attributes
	fwSet, err := mctl.FirmwareSetByVendorModel(ctx, vendor, model, client)
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

	getFirmwareSet.PersistentFlags().StringVar(&flagsDefinedGetFirmwareSet.id, "id", "", "firmware set UUID")
	getFirmwareSet.PersistentFlags().StringVar(&flagsDefinedGetFirmwareSet.serverID, "server", "", "server UUID")

	getFirmwareSet.MarkFlagsMutuallyExclusive("id", "server")

}
