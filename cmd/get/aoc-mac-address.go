package get

import (
	"log"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

type getBomInfoByAocMacAddressFlags struct {
	macAddr string
}

var (
	flagsGetBomByAocMacAddress *getBomInfoByAocMacAddressFlags
)

var getBomInfoByAocMacAddress = &cobra.Command{
	Use:   "aocmacaddress",
	Short: "Get bom object by aocMacAddr",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		bomInfo, _, err := client.GetBomInfoByAOCMacAddr(cmd.Context(), flagsGetBomByAocMacAddress.macAddr)
		if err != nil {
			log.Fatal(err)
		}

		writeResults(bomInfo)
	},
}

func init() {
	flagsGetBomByAocMacAddress = &getBomInfoByAocMacAddressFlags{}

	getBomInfoByAocMacAddress.PersistentFlags().StringVar(&flagsGetBomByAocMacAddress.macAddr, "get-bom-by-aoc-mac-address", "", "get bom info by aocMacAddr")

	if err := getBomInfoByAocMacAddress.MarkPersistentFlagRequired("get-bom-by-aoc-mac-address"); err != nil {
		log.Fatal(err)
	}
}
