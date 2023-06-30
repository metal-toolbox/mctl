package deleteResource

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

var (
	deleteFWSetFlags mctl.FirmwareSetFlags
)

var deleteFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Delete a firmware set",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		c, err := app.NewServerserviceClient(cmd.Context(), theApp)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(deleteFWSetFlags.ID)
		if err != nil {
			log.Fatal(err)
		}

		_, err = c.DeleteServerComponentFirmwareSet(cmd.Context(), id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("firmware set deleted: " + id.String())
	},
}

func init() {
	deleteFirmwareSet.PersistentFlags().StringVar(&deleteFWSetFlags.ID, "uuid", "", "UUID of firmware set to be deleted")

	if err := deleteFirmwareSet.MarkPersistentFlagRequired("uuid"); err != nil {
		log.Fatal(err)
	}
}
