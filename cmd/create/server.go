package create

import (
	"context"
	"encoding/json"
	"log"

	coapiv1 "github.com/metal-toolbox/conditionorc/pkg/api/v1/conditions/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

type serverEnrollParams struct {
	serverID string
	facility string
	ip       string
	username string
	password string
}

var serverEnrollFlags *serverEnrollParams

var serverEnroll = &cobra.Command{
	Use:   "server",
	Short: "Enroll server and publish conditions",
	Run: func(cmd *cobra.Command, _ []string) {
		enrollServer(cmd.Context())
	},
}

func enrollServer(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	params, err := json.Marshal(coapiv1.AddServerParams{
		Facility: serverEnrollFlags.facility,
		IP:       serverEnrollFlags.ip,
		Username: serverEnrollFlags.username,
		Password: serverEnrollFlags.password,
	})
	if err != nil {
		log.Fatal(err)
	}

	conditionCreate := coapiv1.ConditionCreate{
		Parameters: params,
	}

	client, err := app.NewConditionsClient(ctx, theApp.Config.Conditions, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.ServerEnroll(ctx, serverEnrollFlags.serverID, conditionCreate)
	if err != nil {
		log.Fatal(err)
	}

	condition, err := mctl.ConditionFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%d\nmsg=%s\nconditionID=%s\nserverID=%v", response.StatusCode, response.Message, condition.ID, response.Records.ServerID)
}

func init() {
	serverEnrollFlags = &serverEnrollParams{}

	mctl.AddBMCAddressFlag(serverEnroll, &serverEnrollFlags.ip)
	mctl.AddBMCUsernameFlag(serverEnroll, &serverEnrollFlags.username)
	mctl.AddBMCPasswordFlag(serverEnroll, &serverEnrollFlags.password)
	mctl.AddFacilityFlag(serverEnroll, &serverEnrollFlags.facility)
	mctl.AddServerFlag(serverEnroll, &serverEnrollFlags.serverID)

	mctl.RequireFlag(serverEnroll, mctl.BMCAddressFlag)
	mctl.RequireFlag(serverEnroll, mctl.BMCUsernameFlag)
	mctl.RequireFlag(serverEnroll, mctl.BMCPasswordFlag)
	mctl.RequireFlag(serverEnroll, mctl.FacilityFlag)
}
