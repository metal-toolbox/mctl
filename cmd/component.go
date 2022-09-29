package cmd

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

// Get

type getComponentFlags struct {
	// server UUID
	id string
}

var (
	flagsDefinedGetComponent *getComponentFlags
)

var cmdGetComponent = &cobra.Command{
	Use:   "component",
	Short: "Get Component",
	Run: func(cmd *cobra.Command, args []string) {
		mctl, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			log.Fatal(err)
		}

		c, err := newServerserviceClient(cmd.Context(), mctl)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(flagsDefinedGetComponent.id)
		if err != nil {
			log.Fatal(err)
		}

		components, _, err := c.GetComponents(cmd.Context(), id, nil)
		if err != nil {
			log.Fatal(err)
		}

		if outputJSON {
			printJSON(components)
			os.Exit(0)
		}

		spew.Dump(components)

		//	table := tablewriter.NewWriter(os.Stdout)
		//	table.SetHeader([]string{"UUID", "Vendor", "Model", "Component", "Version"})
		//	for _, c := range components {
		//		c.
		//			table.Append([]string{f.UUID.String(), f.Vendor, f.Model, f.Component, f.Version})
		//	}
		//	table.Render()
	},
}

func init() {
	flagsDefinedGetComponent = &getComponentFlags{}

	cmdGetComponent.PersistentFlags().StringVar(&flagsDefinedGetComponent.id, "server-uuid", "", "server UUID")

	if err := cmdGetComponent.MarkPersistentFlagRequired("server-uuid"); err != nil {
		log.Fatal(err)
	}
}
