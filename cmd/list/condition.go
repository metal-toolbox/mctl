package list

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/metal-toolbox/conditionorc/pkg/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

type listConditionFlags struct {
	// server UUID
	id    string
	state string
}

var (
	flagsDefinedListCondition *listConditionFlags
)

var listCondition = &cobra.Command{
	Use:   "condition",
	Short: "list server conditions by state",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewConditionsClient(cmd.Context(), theApp.Config.Conditions, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(flagsDefinedListCondition.id)
		if err != nil {
			log.Fatal(err)
		}

		response, err := client.ServerConditionList(cmd.Context(), id, types.ConditionState(flagsDefinedListCondition.state))
		if err != nil {
			log.Fatal(err)
		}

		printJSON(response)
		os.Exit(0)
	},
}

func init() {
	flagsDefinedListCondition = &listConditionFlags{}

	listCondition.PersistentFlags().StringVar(&flagsDefinedListCondition.id, "server-id", "", "server UUID")
	listCondition.PersistentFlags().StringVar(&flagsDefinedListCondition.state, "state", "", "condition state")

	if err := listCondition.MarkPersistentFlagRequired("server-id"); err != nil {
		log.Fatal(err)
	}

	if err := listCondition.MarkPersistentFlagRequired("state"); err != nil {
		log.Fatal(err)
	}

}
