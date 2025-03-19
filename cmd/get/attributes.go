package get

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
)

var attributesServerID string

// Get firmware info
var getAttributes = &cobra.Command{
	Use:   "attributes",
	Short: "Get attributes",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		fleetClient, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		serverID, err := uuid.Parse(attributesServerID)
		if err != nil {
			log.Fatal(err)
		}

		params := &fleetdbapi.PaginationParams{
			Limit:   1000,
			Page:    1,
			Preload: false,
		}
		attributes, _, err := fleetClient.ListAttributes(ctx, serverID, params)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Attributes:\n")

		for i := range attributes {
			log.Printf("%24s %s\n", attributes[i].Namespace, string(attributes[i].Data))
		}

		log.Printf("Versioned Attributes\n")

		versionedAttributes, _, err := fleetClient.ListVersionedAttributes(ctx, serverID)
		if err != nil {
			log.Fatal(err)
		}

		for i := range versionedAttributes {
			log.Printf("%d: %24s %s\n", versionedAttributes[i].Tally, versionedAttributes[i].Namespace, string(versionedAttributes[i].Data))
		}

		os.Exit(0)
	},
}

func init() {
	cmdGet.AddCommand(getAttributes)

	mctl.AddServerFlag(getAttributes, &attributesServerID)
	mctl.RequireFlag(getAttributes, mctl.ServerFlag)
}