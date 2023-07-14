package install

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	coapi "github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
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

	fmt.Println(formatCondition(resp))
}

type conditionDisplay struct {
	ID         uuid.UUID              `json:"id"`
	Kind       cotypes.ConditionKind  `json:"kind"`
	State      cotypes.ConditionState `json:"state"`
	Parameters json.RawMessage        `json:"parameters"`
	Status     json.RawMessage        `json:"status"`
}

// XXX: this logs Fatal on errors, and I don't love it but the choice is boilerplate error
// definitions and handling or do something expediant in a one-file command function.
func formatCondition(resp *coapi.ServerResponse) string {
	if resp.Records == nil {
		log.Fatal("no records returned")
	}

	if len(resp.Records.Conditions) == 0 {
		log.Fatal("no install condition found")
	}

	inc := resp.Records.Conditions[0]

	display := &conditionDisplay{
		ID:         inc.ID,
		Kind:       inc.Kind,
		Parameters: inc.Parameters,
		State:      inc.State,
		Status:     inc.Status,
	}

	// XXX: seems highly unlikely that we get a response that deserializes cleanly and doesn't
	// re-serialize.
	b, err := json.MarshalIndent(display, "", "  ")
	if err != nil {
		log.Fatalf("bad json in response: %s", err.Error())
	}

	return string(b)
}

func init() {
	flags := installStatus.Flags()
	flags.StringVarP(&serverIDStr, "server", "s", "", "server id (typically a UUID)")

	if err := installStatus.MarkFlagRequired("server"); err != nil {
		log.Fatalf("marking server flag as required: %s", err.Error())
	}

	install.AddCommand(installStatus)
}
