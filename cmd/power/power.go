package power

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	mctl "github.com/metal-toolbox/mctl/cmd"

	"github.com/metal-toolbox/conditionorc/pkg/api/v1/client"
	coapiv1 "github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
	rctypes "github.com/metal-toolbox/rivets/condition"

	"github.com/metal-toolbox/mctl/internal/app"
)

var powerCmd = &cobra.Command{
	Use:   "power",
	Short: "Execute server/bmc power and next boot commands: " + strings.Join(serverPowerActions, ","),
	Run: func(cmd *cobra.Command, args []string) {
		powerAction(cmd.Context())
	},
}

func init() {
	mctl.RootCmd.AddCommand(powerCmd)
}

var (
	flagsDefinedPowerAction *powerActionFlags
	queryActionStatus       bool
	serverPowerActions      = []string{
		"on",
		"off",
		"cycle",
		"reset",
		"soft",
		"status",
		"bmc-reset",
	}
)

type powerActionFlags struct {
	serverID  string
	parameter string
}

func powerAction(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	c, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	serverID, err := uuid.Parse(flagsDefinedPowerAction.serverID)
	if err != nil {
		log.Fatal(err)
	}

	if queryActionStatus {
		actionStatus(ctx, serverID, c)
		return
	}

	controlParams, err := paramsFromFlags(flagsDefinedPowerAction)
	if err != nil {
		log.Fatal(err)
	}

	params, err := json.Marshal(controlParams)
	if err != nil {
		log.Fatal(err)
	}

	conditionCreate := coapiv1.ConditionCreate{
		Parameters: params,
	}

	response, err := c.ServerConditionCreate(ctx, serverID, rctypes.ServerControl, conditionCreate)
	if err != nil {
		log.Fatal(err)
	}

	condition, err := mctl.ConditionFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%d msg=%s conditionID=%s", response.StatusCode, response.Message, condition.ID)
}

func actionStatus(ctx context.Context, serverID uuid.UUID, c *client.Client) {
	resp, err := c.ServerConditionStatus(ctx, serverID)
	if err != nil {
		log.Fatalf("querying server conditions: %s", err.Error())
	}

	s, err := mctl.FormatConditionResponse(resp, rctypes.ServerControl)
	if err != nil {
		log.Fatalf("condition response error: %s", err.Error())
	}

	fmt.Println(s)
}

func paramsFromFlags(f *powerActionFlags) (*rctypes.ServerControlTaskParameters, error) {
	actionParam := strings.ToLower(f.parameter)
	if !slices.Contains(serverPowerActions, actionParam) {
		return nil, errors.New("invalid power action requested: " + actionParam)
	}

	var action rctypes.ServerControlAction

	switch actionParam {
	case "on", "off", "cycle", "reset", "soft":
		action = rctypes.SetPowerState
	case "bmc-reset":
		action = rctypes.PowerCycleBMC
	case "status":
		action = rctypes.GetPowerState
	}

	return rctypes.NewServerControlTaskParameters(
		uuid.MustParse(f.serverID),
		action,
		actionParam,
		false,
		false,
	), nil
}

func init() {
	flagsDefinedPowerAction = &powerActionFlags{}

	mctl.AddServerFlag(powerCmd, &flagsDefinedPowerAction.serverID)
	mctl.AddServerPowerActionFlag(powerCmd, &flagsDefinedPowerAction.parameter, serverPowerActions)
	mctl.AddServerPowerActionStatusFlag(powerCmd, &queryActionStatus)
	mctl.MutuallyExclusiveFlags(powerCmd, mctl.ServerActionPowerActionFlag, mctl.ServerActionPowerActionStatusFlag)
	mctl.RequireOneFlag(powerCmd, mctl.ServerActionPowerActionFlag, mctl.ServerActionPowerActionStatusFlag)
	mctl.RequireFlag(powerCmd, mctl.ServerFlag)
}
