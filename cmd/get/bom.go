package get

import (
	"log"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

type getBomInfoByBmcMacAddressFlags struct {
	aocMacAddr string
	bmcMacAddr string
}

func (gb *getBomInfoByBmcMacAddressFlags) hasAOCMacAddr() bool {
	return gb.aocMacAddr != ""
}

func (gb *getBomInfoByBmcMacAddressFlags) hasBMCMacAddr() bool {
	return gb.bmcMacAddr != ""
}

var (
	flagsGetBomByMacAddress *getBomInfoByBmcMacAddressFlags
)

var getBomInfoByMacAddress = &cobra.Command{
	Use:   "bom",
	Short: "Get bom object by AOC or BMC Addr",
	Run: func(cmd *cobra.Command, args []string) {
		if !flagsGetBomByMacAddress.hasAOCMacAddr() && !flagsGetBomByMacAddress.hasBMCMacAddr() {
			log.Fatalf("--aoc-mac and --bmc-mac not set")
		}

		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		var bomInfo *serverservice.Bom
		if flagsGetBomByMacAddress.hasAOCMacAddr() {
			bomInfo, _, err = client.GetBomInfoByAOCMacAddr(cmd.Context(), flagsGetBomByMacAddress.aocMacAddr)
		} else {
			bomInfo, _, err = client.GetBomInfoByBMCMacAddr(cmd.Context(), flagsGetBomByMacAddress.bmcMacAddr)
		}
		if err != nil {
			log.Fatal(err)
		}

		mctl.PrintResults(output, bomInfo)
	},
}

func init() {
	flagsGetBomByMacAddress = &getBomInfoByBmcMacAddressFlags{}

	getBomInfoByMacAddress.PersistentFlags().StringVar(&flagsGetBomByMacAddress.aocMacAddr, "aoc-mac", "", "get bom info by aoc mac address")
	getBomInfoByMacAddress.PersistentFlags().StringVar(&flagsGetBomByMacAddress.bmcMacAddr, "bmc-mac", "", "get bom info by bmc mac address")

	getBomInfoByMacAddress.MarkFlagsMutuallyExclusive("aoc-mac", "bmc-mac")
}
