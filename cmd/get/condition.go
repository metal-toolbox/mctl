package get

import (
	"log"

	"github.com/google/uuid"
	"github.com/metal-toolbox/conditionorc/pkg/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

type getConditionFlags struct {
	// server UUID
	id   string
	kind string
}

var (
	flagsDefinedGetCondition *getConditionFlags
)

var getCondition = &cobra.Command{
	Use:   "condition",
	Short: "get server condition",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewConditionsClient(cmd.Context(), theApp.Config.Conditions)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(flagsDefinedGetCondition.id)
		if err != nil {
			log.Fatal(err)
		}

		response, err := client.ServerConditionGet(cmd.Context(), id, types.ConditionKind(flagsDefinedGetCondition.kind))
		if err != nil {
			log.Fatal(err)
		}

		writeResults(response)
	},
}

func init() {
	flagsDefinedGetCondition = &getConditionFlags{}

	getCondition.PersistentFlags().StringVar(&flagsDefinedGetCondition.id, "server-id", "", "server UUID")
	getCondition.PersistentFlags().StringVar(&flagsDefinedGetCondition.kind, "kind", "", "condition kind")

	if err := getCondition.MarkPersistentFlagRequired("server-id"); err != nil {
		log.Fatal(err)
	}

	if err := getCondition.MarkPersistentFlagRequired("kind"); err != nil {
		log.Fatal(err)
	}

}
