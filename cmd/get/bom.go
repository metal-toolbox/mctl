package get

import (
	"log"

	"github.com/spf13/cobra"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
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

	mctl.AddMacAOCFlag(getBomInfoByMacAddress, &flagsGetBomByMacAddress.aocMacAddr)
	mctl.AddMacBMCFlag(getBomInfoByMacAddress, &flagsGetBomByMacAddress.bmcMacAddr)

	mctl.MutuallyExclusiveFlags(getBomInfoByMacAddress, mctl.MacAOCFlag, mctl.MacBMCFlag)
}
