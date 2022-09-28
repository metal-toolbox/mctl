package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

var cmdListFirmware = &cobra.Command{
	Use:   fmt.Sprintf("firmware"),
	Short: "List firmware",
	Run: func(cmd *cobra.Command, args []string) {
		mctl, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		c, err := newServerserviceClient(cmd.Context(), mctl)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}

		firmware, _, err := c.ListServerComponentFirmware(cmd.Context(), nil)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}

		spew.Dump(firmware)
	},
}

func init() {
	cmdList.AddCommand(cmdListFirmware)
}
