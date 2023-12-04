package create

import (
	"context"
	"encoding/json"
	"log"

	"github.com/spf13/cobra"

	"github.com/metal-toolbox/mctl/internal/app"

	coapiv1 "github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
	mctl "github.com/metal-toolbox/mctl/cmd"
)

type serverEnrollParams struct {
	serverID string
	facility string
	ip       string
	username string
	password string
}

var (
	serverEnrollFlags *serverEnrollParams
)

var serverEnroll = &cobra.Command{
	Use:   "server",
	Short: "Enroll server and publish conditions",
	Run: func(cmd *cobra.Command, args []string) {
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
		Exclusive:  false,
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

	serverEnroll.PersistentFlags().StringVar(&serverEnrollFlags.serverID, "server-id", "", "server id to be created. New id will be created if null.")
	serverEnroll.PersistentFlags().StringVar(&serverEnrollFlags.facility, "facility", "", "facility of the server")
	serverEnroll.PersistentFlags().StringVar(&serverEnrollFlags.ip, "ip", "", "ip of the server")
	serverEnroll.PersistentFlags().StringVar(&serverEnrollFlags.username, "user", "", "username of the server")
	serverEnroll.PersistentFlags().StringVar(&serverEnrollFlags.password, "pwd", "", "password of the server")

	if err := serverEnroll.MarkPersistentFlagRequired("facility"); err != nil {
		log.Fatal(err)
	}

	if err := serverEnroll.MarkPersistentFlagRequired("ip"); err != nil {
		log.Fatal(err)
	}

	if err := serverEnroll.MarkPersistentFlagRequired("user"); err != nil {
		log.Fatal(err)
	}

	if err := serverEnroll.MarkPersistentFlagRequired("pwd"); err != nil {
		log.Fatal(err)
	}
}
