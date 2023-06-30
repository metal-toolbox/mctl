package get

import (
	"log"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

type getComponentFlags struct {
	// server UUID
	id string
}

var (
	flagsDefinedGetComponent *getComponentFlags
)

var getComponent = &cobra.Command{
	Use:   "component",
	Short: "get server components",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		c, err := app.NewServerserviceClient(cmd.Context(), theApp)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(flagsDefinedGetComponent.id)
		if err != nil {
			log.Fatal(err)
		}

		components, _, err := c.GetComponents(cmd.Context(), id, nil)
		if err != nil {
			log.Fatal(err)
		}

		writeResults(components)
	},
}

func init() {
	flagsDefinedGetComponent = &getComponentFlags{}

	getComponent.PersistentFlags().StringVar(&flagsDefinedGetComponent.id, "server-id", "", "server UUID")

	if err := getComponent.MarkPersistentFlagRequired("server-id"); err != nil {
		log.Fatal(err)
	}
}
