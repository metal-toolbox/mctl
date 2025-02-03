package create

import (
	"log"
	"os"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

// Create Bom informations.
type uploadBomFileFlags struct {
	// Xlsx file containing one or multiple boms information.
	bomXlsxFile string
}

var flagsUploadBomFileFlags *uploadBomFileFlags

var uploadBomFile = &cobra.Command{
	Use:   "bom",
	Short: "Upload Bom File",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewBomServiceClient(cmd.Context(), theApp.Config.BomService, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		fBytes, err := os.ReadFile(flagsUploadBomFileFlags.bomXlsxFile)
		if err != nil {
			log.Fatal(err)
		}

		serverResp, err := client.XlsxFileUpload(cmd.Context(), fBytes)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(serverResp)
	},
}

func init() {
	flagsUploadBomFileFlags = &uploadBomFileFlags{}
	usage := "xlsx file with BOM information"

	mctl.AddFromFileFlag(uploadBomFile, &flagsUploadBomFileFlags.bomXlsxFile, usage)
	mctl.RequireFlag(uploadBomFile, mctl.FromFileFlag)
}
