package get

import (
	"log"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

type getBomInfoByBmcMacAddressFlags struct {
	macAddr string
}

var (
	flagsGetBomByBmcMacAddress *getBomInfoByBmcMacAddressFlags
)

var getBomInfoByBmcMacAddress = &cobra.Command{
	Use:   "bmcmacaddress",
	Short: "Get bom object by bmcMacAddr",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewServerserviceClient(cmd.Context(), theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		bomInfo, _, err := client.GetBomInfoByBMCMacAddr(cmd.Context(), flagsGetBomByBmcMacAddress.macAddr)
		if err != nil {
			log.Fatal(err)
		}

		writeResults(bomInfo)
	},
}

func init() {
	flagsGetBomByBmcMacAddress = &getBomInfoByBmcMacAddressFlags{}

	getBomInfoByBmcMacAddress.PersistentFlags().StringVar(&flagsGetBomByBmcMacAddress.macAddr, "get-bom-by-bmc-mac-address", "", "get bom info by bmcMacAddr")

	if err := getBomInfoByBmcMacAddress.MarkPersistentFlagRequired("get-bom-by-bmc-mac-address"); err != nil {
		log.Fatal(err)
	}
}
