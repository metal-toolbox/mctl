package get

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	attr "github.com/metal-toolbox/mctl/pkg/attributes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	ss "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	errInitialCall = errors.New("error reaching out to server-service")
	errIteration   = errors.New("component results interrupted by error")
	cmdTimeout     = 2 * time.Minute
	serverIDStr    string
	onePage        bool
	page           int
)

var getServerFirmware = &cobra.Command{
	Use:   "firmware {-s | --server-id} <server uuid>",
	Short: "Get all firmware components on a server",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), cmdTimeout)
		defer cancel()

		c, err := app.NewServerserviceClient(ctx, theApp)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(serverIDStr)
		if err != nil {
			log.Fatal(err)
		}

		var cmps []ss.ServerComponent
		if onePage {
			cmps, err = getSingleComponentsPage(ctx, c, id)
		} else {
			cmps, err = getAllComponents(ctx, c, id)
		}
		if err != nil {
			log.Fatalf("error getting firmware: %s", err.Error())
		}
		log.Printf("retrieved %d components", len(cmps))
		// select only those components that have a firmware attribute
		fwset := attr.FirmwareFromComponents(cmps)
		writeResults(fwset)
	},
}

func getSingleComponentsPage(ctx context.Context, c *ss.Client, id uuid.UUID) ([]ss.ServerComponent, error) {
	params := &ss.PaginationParams{
		Page: page,
	}

	cmps, _, err := c.GetComponents(ctx, id, params)
	if err != nil {
		return nil, err
	}

	return cmps, nil
}

func getAllComponents(ctx context.Context, c *ss.Client, id uuid.UUID) ([]ss.ServerComponent, error) {
	params := &ss.PaginationParams{}

	cmps, resp, err := c.GetComponents(ctx, id, params)
	if err != nil {
		return nil, errors.Wrap(errInitialCall, err.Error())
	}

	currentPage := resp.Page
	stopAt := resp.TotalPages + 1
	for currentPage < stopAt {
		params.Page = (currentPage + 1)
		next, resp, err := c.GetComponents(ctx, id, params)
		if err != nil {
			return nil, errors.Wrap(errIteration, err.Error())
		}
		cmps = append(cmps, next...)
		currentPage = resp.Page
		log.Printf("Debug -- retrieved page: %d", currentPage)
	}

	return cmps, nil
}

func init() {
	flags := getServerFirmware.Flags()

	flags.StringVarP(&serverIDStr, "server-id", "s", "", "the server id to look up")

	if err := getServerFirmware.MarkFlagRequired("server-id"); err != nil {
		log.Fatalf("getServerFirmware -- set server-id required: %s", err.Error())
	}

	flags.BoolVarP(&onePage, "limit-one", "1", false, "return only a single page of results")
	flags.IntVarP(&page, "page-number", "n", 1, "the results page to retrieve (only valid with --limit-one")

	getServerFirmware.MarkFlagsRequiredTogether("limit-one", "page-number")
}
