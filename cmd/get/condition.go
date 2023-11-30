package get

import (
	"log"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

type getConditionFlags struct {
	// server UUID
	id string
}

var (
	flagsDefinedGetCondition *getConditionFlags
)

var getCondition = &cobra.Command{
	Use:   "condition",
	Short: "get the last server conditions performed",
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

		response, err := client.ServerConditionStatus(cmd.Context(), id)
		if err != nil {
			log.Fatal(err)
		}

		mctl.PrintResults(output, response)
	},
}

func init() {
	flagsDefinedGetCondition = &getConditionFlags{}

	getCondition.Flags().StringVarP(&flagsDefinedGetCondition.id, "server", "s", "", "server UUID")

	if err := getCondition.MarkFlagRequired("server"); err != nil {
		log.Fatal(err)
	}

}
