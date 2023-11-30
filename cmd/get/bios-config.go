package get

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
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
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
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
func biosConfigFromNamespaces(ctx context.Context, serverID uuid.UUID, client *serverservice.Client) ([]serverservice.VersionedAttributes, error) {
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

	getBiosConfig.PersistentFlags().StringVar(&flagsDefinedGetBiosConfig.serverID, "server", "", "server UUID")

	if err := getBiosConfig.MarkPersistentFlagRequired("server"); err != nil {
		log.Fatal(err)
	}
}
