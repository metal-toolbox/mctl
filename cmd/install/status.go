package install

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	cotypes "github.com/metal-toolbox/conditionorc/pkg/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

var serverIDStr string

var installStatus = &cobra.Command{
	Use:   "status --server | -s <server uuid>",
	Short: "check the progress of a firmware install on a server",
	Run: func(cmd *cobra.Command, _ []string) {
		statusCheck(cmd.Context())
	},
}

func statusCheck(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	serverID, err := uuid.Parse(serverIDStr)
	if err != nil {
		log.Fatalf("parsing server id: %s", err.Error())
	}

	resp, err := client.ServerConditionGet(ctx, serverID, cotypes.FirmwareInstall)
	if err != nil {
		log.Fatalf("querying server conditions: %s", err.Error())
	}

	fmt.Println(mctl.FormatConditionResponse(resp))
}

func init() {
	flags := installStatus.Flags()
	flags.StringVarP(&serverIDStr, "server", "s", "", "server id (typically a UUID)")

	if err := installStatus.MarkFlagRequired("server"); err != nil {
		log.Fatalf("marking server flag as required: %s", err.Error())
	}

	install.AddCommand(installStatus)
}
