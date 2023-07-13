package deleteresource

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/metal-toolbox/conditionorc/pkg/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

type deleteConditionFlags struct {
	// server UUID
	id   string
	kind string
}

var (
	flagsDefinedGetCondition *deleteConditionFlags
)

var deleteCondition = &cobra.Command{
	Use:   "condition",
	Short: "delete server condition",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewConditionsClient(cmd.Context(), theApp.Config.Conditions, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(flagsDefinedGetCondition.id)
		if err != nil {
			log.Fatal(err)
		}

		response, err := client.ServerConditionDelete(cmd.Context(), id, types.ConditionKind(flagsDefinedGetCondition.kind))
		if err != nil {
			log.Fatal(err)
		}

		spew.Dump(response)
	},
}

func init() {
	flagsDefinedGetCondition = &deleteConditionFlags{}

	deleteCondition.PersistentFlags().StringVar(&flagsDefinedGetCondition.id, "server-id", "", "server UUID")
	deleteCondition.PersistentFlags().StringVar(&flagsDefinedGetCondition.kind, "kind", "", "condition kind")

	if err := deleteCondition.MarkPersistentFlagRequired("server-id"); err != nil {
		log.Fatal(err)
	}

	if err := deleteCondition.MarkPersistentFlagRequired("kind"); err != nil {
		log.Fatal(err)
	}

}
