package cmd

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"

	coApi "github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
	coTyp "github.com/metal-toolbox/conditionorc/pkg/types"
)

type installParams struct {
	FirmwareID uuid.UUID `json:"firmwareId"`
	IsSet      bool      `json:"isSet"`
}

func (i *installParams) MustBytes() json.RawMessage {
	byt, err := json.Marshal(i)
	if err != nil {
		log.Fatalf("marshaling install parameters: %s", err.Error())
	}
	return byt
}

// fail fatally if we fail any argument sanity checks
func composeInstallParams(uuidStr string, isSet bool) []byte {
	fmwID, err := uuid.Parse(uuidStr)
	if err != nil {
		log.Fatalf("firmware id is not a uuid: %s", err.Error())
	}
	params := installParams{
		FirmwareID: fmwID,
		IsSet:      isSet,
	}
	return params.MustBytes()
}

var (
	fwFlagName       = "firmware-id"
	setFlagName      = "fwset-id"
	serverIDStr      string
	firmwareIDStr    string
	firmwareSetIDStr string
)

// install firmware on a server
var installFirmware = &cobra.Command{
	Use:     "install --server-id server-uuid { --firmware-id firmware-uuid | --fwset-id set-uuid }",
	Aliases: []string{"flash"},
	Short:   "install firmware or a firmware set on a server",
	Args:    cobra.ExactArgs(0),
	Run: func(c *cobra.Command, args []string) {
		ctx := c.Context()
		mctl, err := app.New(ctx, cfgFile)
		if err != nil {
			log.Fatalf("creating app: %s", err.Error())
		}

		client, err := newConditionsClient(ctx, mctl)
		if err != nil {
			log.Fatalf("creating condition client: %s", err.Error())
		}

		srvID, err := uuid.Parse(serverIDStr)
		if err != nil {
			log.Fatalf("server id invalid: %s", err.Error())
		}

		var payload json.RawMessage
		switch {
		case firmwareIDStr != "":
			payload = composeInstallParams(firmwareIDStr, false)
		case firmwareSetIDStr != "":
			payload = composeInstallParams(firmwareSetIDStr, true)
		default:
			log.Fatalf("either %s or %s must be used", fwFlagName, setFlagName)
		}

		create := coApi.ConditionCreate{
			Parameters: payload,
		}

		resp, err := client.ServerConditionCreate(ctx, srvID, coTyp.FirmwareInstall, create)
		if err != nil {
			log.Fatalf("Error returned from creating the server condition: %s", err.Error())
		}

		if resp == nil {
			log.Fatal("nil response from server")
		}

		if resp.Message != "" {
			log.Printf("Message: %s", resp.Message)
		}

		if resp.Records == nil {
			log.Fatal("no condition records returned")
		}

		retCondID := "not returned"
		retSrvID := resp.Records.ServerID
		if len(resp.Records.Conditions) > 0 {
			retCondID = resp.Records.Conditions[0].ID.String()
		}
		log.Printf("Server => %s\nCondition =>%s\n", retSrvID, retCondID)
	},
}

func init() {
	installFirmware.PersistentFlags().StringVar(
		&serverIDStr, "server-id", "", "server uuid string")
	installFirmware.PersistentFlags().StringVarP(
		&firmwareIDStr, fwFlagName, "f", "", "firmware uuid string")
	installFirmware.PersistentFlags().StringVarP(
		&firmwareSetIDStr, setFlagName, "s", "", "uuid of a set of firmware to be applied as a unit")

	if err := installFirmware.MarkPersistentFlagRequired("server-id"); err != nil {
		log.Fatalf("make server-id required: %s", err.Error())
	}

	installFirmware.MarkFlagsMutuallyExclusive(fwFlagName, setFlagName)

	RootCmd.AddCommand(installFirmware)
}
