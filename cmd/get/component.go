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

type getComponentFlags struct {
	// server UUID
	id     string
	fwOnly bool
}

var (
	cmdArgs    *getComponentFlags
	cmdTimeout = 2 * time.Minute
)

var getComponent = &cobra.Command{
	Use:   "component",
	Short: "get server components",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(cmd.Context(), cmdTimeout)
		defer cancel()

		theApp := mctl.MustCreateApp(ctx)

		client, err := app.NewServerserviceClient(ctx, theApp.Config.Serverservice, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(cmdArgs.id)
		if err != nil {
			log.Fatal(err)
		}

		cmps, err := getComponents(ctx, client, id)
		if err != nil {
			log.Fatalf("error getting firmware: %s", err.Error())
		}

		if cmdArgs.fwOnly {
			fwset := attr.FirmwareFromComponents(cmps)
			writeResults(fwset)
			return
		}

		writeResults(cmps)
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

func init() {
	cmdArgs = &getComponentFlags{}

	cmdPFlags := getComponent.PersistentFlags()

	cmdPFlags.StringVar(&cmdArgs.id, "server", "", "server UUID")
	cmdPFlags.BoolVarP(&cmdArgs.fwOnly, "firmware-only", "f",
		false, "only retrieve components with firmware")

	if err := getComponent.MarkPersistentFlagRequired("server"); err != nil {
		log.Fatal(err)
	}
}
