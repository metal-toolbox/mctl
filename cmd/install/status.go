package install

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	rctypes "github.com/metal-toolbox/rivets/v2/condition"
	"github.com/spf13/cobra"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
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

	resp, err := client.ServerConditionStatus(ctx, serverID)
	if err != nil {
		log.Fatalf("querying server conditions: %s", err.Error())
	}

	s, err := mctl.FormatConditionResponse(resp, rctypes.FirmwareInstall)
	if err != nil {
		log.Fatalf("condition response error: %s", err.Error())
	}

	fmt.Println(s)
}

func init() {
	mctl.AddServerFlag(installStatus, &serverIDStr)
	mctl.RequireFlag(installStatus, mctl.ServerFlag)
}
