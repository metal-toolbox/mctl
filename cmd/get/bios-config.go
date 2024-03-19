package get

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/spf13/cobra"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
)

type getBiosConfigFlags struct {
	serverID string
}

var (
	flagsDefinedGetBiosConfig *getBiosConfigFlags
)

// Get BIOS configuration
var getBiosConfig = &cobra.Command{
	Use:   "bios-config",
	Short: "Get bios configuration information for a server",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		serverID, err := uuid.Parse(flagsDefinedGetBiosConfig.serverID)
		if err != nil {
			log.Fatal(err)
		}

		biosCfg, err := biosConfigFromNamespaces(ctx, serverID, client)
		if err != nil {
			log.Fatal(err)
		}

		if biosCfg == nil {
			log.Println("no bios configuration data found")
			os.Exit(0)
		}

		mctl.PrintResults(output, biosCfg[0])
	},
}

// returns bios configuration data
func biosConfigFromNamespaces(ctx context.Context, serverID uuid.UUID, client *fleetdbapi.Client) ([]fleetdbapi.VersionedAttributes, error) {
	namespaces := []string{
		"sh.hollow.alloy.inband.bios_configuration",
		"sh.hollow.alloy.outofband.bios_configuration",
	}

	for _, ns := range namespaces {
		biosCfg, _, err := client.GetVersionedAttributes(ctx, serverID, ns)
		if err != nil {
			if strings.Contains(err.Error(), "resource not found") {
				continue
			}

			return nil, err
		}

		return biosCfg, nil
	}

	return nil, nil
}

func init() {
	flagsDefinedGetBiosConfig = &getBiosConfigFlags{}

	mctl.AddServerFlag(getBiosConfig, &flagsDefinedGetBiosConfig.serverID)
	mctl.RequireFlag(getBiosConfig, mctl.ServerFlag)
}
