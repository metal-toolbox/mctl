package cmd

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
	"gopkg.in/yaml.v3"

	coApi "github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
	coTyp "github.com/metal-toolbox/conditionorc/pkg/types"
)

var (
	cmdTimeout = 20 * time.Second
)

// List
var cmdListFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "List firmware",
	Run: func(cmd *cobra.Command, args []string) {
		mctl, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(cmd.Context(), cmdTimeout)
		defer cancel()

		c, err := newServerserviceClient(ctx, mctl)
		if err != nil {
			log.Fatal("error initializing serverservice client: ", err)
		}

		firmware, _, err := c.ListServerComponentFirmware(cmd.Context(), nil)
		if err != nil {
			log.Fatal("serverservice client returned error: ", err)
		}

		if outputJSON {
			printJSON(firmware)
			os.Exit(0)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"UUID", "Vendor", "Model", "Component", "Version"})
		for _, f := range firmware {
			table.Append([]string{f.UUID.String(), f.Vendor, strings.Join(f.Model, ","), f.Component, f.Version})
		}
		table.Render()
	},
}

// Create
type createFirmwareFlags struct {
	// file containing firmware configuration
	firmwareConfigFile string
}

var (
	flagsDefinedCreateFirmware *createFirmwareFlags
)

var cmdCreateFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "Create firmware",
	Run: func(cmd *cobra.Command, args []string) {
		mctl, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			log.Fatal(err)
		}

		client, err := newServerserviceClient(cmd.Context(), mctl)
		if err != nil {
			log.Fatal(err)
		}

		firmwareConfig := &model.FirmwareConfig{}
		fbytes, err := os.ReadFile(flagsDefinedCreateFirmware.firmwareConfigFile)
		if err != nil {
			log.Fatal(err)
		}

		if err = yaml.Unmarshal(fbytes, firmwareConfig); err != nil {
			log.Fatal(err)
		}

		for _, config := range firmwareConfig.Firmwares {
			c := serverservice.ComponentFirmwareVersion{
				Vendor:        config.Vendor,
				RepositoryURL: config.RepositoryURL,
				Model:         config.Model,
				UpstreamURL:   config.UpstreamURL,
				Version:       config.Version,
				Filename:      config.FileName,
				Checksum:      config.Checksum,
				Component:     config.Component,
			}

			id, _, err := client.CreateServerComponentFirmware(cmd.Context(), c)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(id)
		}
	},
}

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

var (
	serverIDStr   string
	firmwareIDStr string
	firmwareSet   bool
)

// install firmware on a server
var installFirmware = &cobra.Command{
	Use:     "install --server-id server-uuid --firmware-id firmware-uuid [--set]",
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

		fmwID, err := uuid.Parse(firmwareIDStr)
		if err != nil {
			log.Fatalf("firmware id invalid: %s", err.Error())
		}

		params := installParams{
			FirmwareID: fmwID,
			IsSet:      firmwareSet,
		}

		create := coApi.ConditionCreate{
			Parameters: params.MustBytes(),
		}

		resp, err := client.ServerConditionCreate(ctx, srvID, coTyp.FirmwareInstall, create)
		if err != nil {
			log.Printf("Error returned from creating the server condition: %s", err.Error())
		}
		if resp != nil {
			log.Printf("Message: %s", resp.Message)
			if resp.Records == nil {
				log.Printf("no condition records returned")
			} else {
				retCondID := "not returned"
				retSrvID := resp.Records.ServerID
				if len(resp.Records.Conditions) > 0 {
					retCondID = resp.Records.Conditions[0].ID.String()
				}
				log.Printf("Server => %s\nCondition =>%s\n", retSrvID, retCondID)
			}
		}
	},
}

func init() {
	flagsDefinedCreateFirmware = &createFirmwareFlags{}

	cmdCreateFirmware.PersistentFlags().StringVar(
		&flagsDefinedCreateFirmware.firmwareConfigFile,
		"from-file", "", "YAML file with firmware configuration data")
	if err := cmdCreateFirmware.MarkPersistentFlagRequired("from-file"); err != nil {
		log.Fatal(err)
	}

	installFirmware.PersistentFlags().StringVar(
		&serverIDStr, "server-id", "", "server uuid string")
	installFirmware.PersistentFlags().StringVar(
		&firmwareIDStr, "firmware-id", "", "firmware uuid string")
	installFirmware.PersistentFlags().BoolVar(
		&firmwareSet, "is-set", false, "designates the firmware-id as a firmware set to be applied as a unit")

	if err := installFirmware.MarkPersistentFlagRequired("server-id"); err != nil {
		log.Fatalf("make server-id required: %s", err.Error())
	}
	if err := installFirmware.MarkPersistentFlagRequired("firmware-id"); err != nil {
		log.Fatalf("make firmware-id required: %s", err.Error())
	}

	rootCmd.AddCommand(installFirmware)
}
