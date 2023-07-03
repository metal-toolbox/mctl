package get

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	attr "github.com/metal-toolbox/mctl/pkg/attributes"
	"github.com/spf13/cobra"
	ss "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	cmdTimeout  = 2 * time.Minute
	serverIDStr string
	serverID    uuid.UUID
)

var getServerFirmware = &cobra.Command{
	Use:   "firmware {-s | --server-id} <server uuid>",
	Short: "Get all firmware components on a server",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), cmdTimeout)
		defer cancel()

		c, err := app.NewServerserviceClient(cmd.Context(), theApp)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(serverIDStr)
		if err != nil {
			log.Fatal(err)
		}

		params := &ss.PaginationParams{}
		cmps, resp, err := c.GetComponents(ctx, id, params)
		if err != nil {
			log.Fatalf("error on initial component query: %s", err.Error())
		}
		currentPage := resp.Page
		stopAt := resp.PageCount + 1
		for currentPage < stopAt {
			params.Page = (currentPage + 1)
			next, resp, err := c.GetComponents(ctx, id, params)
			if err != nil {
				log.Printf("component iteration interrupted by an error: %s", err)
				break
			}
			cmps = append(cmps, next...)
			currentPage = resp.Page
			log.Printf("retrieved page: %d", currentPage)
		}
		log.Printf("retrieved %d components", len(cmps))
		// select only those components that have a firmware attribute
		fwset := attr.FirmwareFromComponents(cmps)
		writeResults(fwset)
	},
}

func init() {
	flags := getServerFirmware.PersistentFlags()

	flags.StringVarP(&serverIDStr, "server-id", "s", "", "the server id to look up")

	if err := getServerFirmware.MarkFlagRequired("server-id"); err != nil {
		log.Fatalf("set server-id required: %w", err)
	}
}
