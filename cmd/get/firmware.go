package get

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
	attr "github.com/metal-toolbox/mctl/pkg/attributes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	ss "go.hollow.sh/serverservice/pkg/api/v1"
)

type fwSpecList map[string][]uuid.UUID

var (
	errInitialCall = errors.New("error reaching out to server-service")
	errIteration   = errors.New("component results interrupted by error")
	errFetchFW     = errors.New("fetching firmware failed")
	cmdTimeout     = 2 * time.Minute
	serverIDStr    string
)

var getServerFirmware = &cobra.Command{
	Use:   "firmware {-s | --server-id} <server uuid>",
	Short: "Get all firmware components on a server",
	Run: func(cmd *cobra.Command, args []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), cmdTimeout)
		defer cancel()

		client, err := app.NewServerserviceClient(ctx, theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(serverIDStr)
		if err != nil {
			log.Fatal(err)
		}

		cmps, err := getComponents(ctx, client, id)
		if err != nil {
			log.Fatalf("error getting firmware: %s", err.Error())
		}
		log.Printf("retrieved %d components", len(cmps))
		// select only those components that have a firmware attribute
		fwset := attr.FirmwareFromComponents(cmps)
		writeResults(fwset)
		cmpIDs, err := getFirmwareIDs(ctx, client, fwset)
		if err != nil {
			log.Fatalf("error getting firmware ids: %s", err.Error())
		}
		writeResults(cmpIDs)
	},
}

func getComponents(ctx context.Context, c *ss.Client, id uuid.UUID) ([]ss.ServerComponent, error) {
	params := &ss.PaginationParams{}

	cmps, resp, err := c.GetComponents(ctx, id, params)
	if err != nil {
		return nil, err
	}

	// XXX: the default result-set size is 100, so more than 100 components will trip
	// the following error.
	if resp.TotalPages > 0 {
		return nil, errors.New("too many components -- add pagination")
	}

	return cmps, nil
}

// XXX: Notice that the *component* is not part of the search params. If we have a collision
// on components with identical model/vendor/version we're going to get weird results and
// likely won't know until we try to install that firmware.
func getSearchParams(cmp *attr.ComponentWithFirmware) *ss.ComponentFirmwareVersionListParams {
	return &ss.ComponentFirmwareVersionListParams{
		Vendor: cmp.Vendor,
		Model: []string{
			strings.ToLower(cmp.Model),
		},
		Version: cmp.Firmware.Installed,
	}
}

// XXX: Consider getting all firmware in one shot?

// Call server-service and get ids for any firmware that matches the tuple of
// vendor/component/model/version. We are as permissive as we can be here, if
// Alloy didn't log a given datum, just search with what we have.
func getFirmwareIDs(ctx context.Context, client *ss.Client,
	cmps []*attr.ComponentWithFirmware) (fwSpecList, error) {
	fws := make(map[string][]uuid.UUID)
	for _, cmp := range cmps {
		var ids []uuid.UUID
		params := getSearchParams(cmp)
		log.Printf("DEBUG search params for %s: %#v\n", cmp.Name, params)
		// XXX: we'll need to refactor this if we ever have more than a single page (~100 entries) of
		// results for a single component.
		fwRecords, _, err := client.ListServerComponentFirmware(ctx, params)
		if err != nil {
			return nil, errors.Wrap(errFetchFW,
				fmt.Sprintf("%s:%s:%s : %s", cmp.Name, cmp.Vendor, cmp.Model, err.Error()),
			)
		}
		log.Printf("DEBUG %s search returns %d records\n", cmp.Name, len(fwRecords))
		for _, record := range fwRecords {
			ids = append(ids, record.UUID)
		}
		fws[cmp.Name] = ids
	}
	return fws, nil
}

func init() {
	flags := getServerFirmware.Flags()

	flags.StringVarP(&serverIDStr, "server-id", "s", "", "the server id to look up")

	if err := getServerFirmware.MarkFlagRequired("server-id"); err != nil {
		log.Fatalf("getServerFirmware -- set server-id required: %s", err.Error())
	}
}
