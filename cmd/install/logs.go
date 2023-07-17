package install

import (
	"context"
	"log"

	"github.com/google/uuid"
	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/logindex"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type installLogsFlags struct {
	serverID string
}

var (
	flagsDefinedInstallLogs *installLogsFlags
	errInvalidServerID      = errors.New("invalid server UUID")
)

var cmdInstallLogs = &cobra.Command{
	Use:   "logs --server",
	Short: "check logs for a firmware install",
	Run: func(cmd *cobra.Command, _ []string) {
		queryLogs(cmd.Context())
	},
}

func queryLogs(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	index, err := logindex.NewQueryor(theApp.Config.Splunk)
	if err != nil {
		log.Fatal(err)
	}

	serverUUID, err := uuid.Parse(flagsDefinedInstallLogs.serverID)
	if err != nil {
		log.Fatal(errors.Wrap(errInvalidServerID, err.Error()))
	}

	if err := index.SearchByAssetID(ctx, serverUUID, uuid.Nil); err != nil {
		log.Fatal(err)
	}
}

func init() {
	flagsDefinedInstallLogs = &installLogsFlags{}

	install.AddCommand(cmdInstallLogs)
	cmdInstallLogs.PersistentFlags().StringVar(&flagsDefinedInstallLogs.serverID, "server", "", "server UUID")

	if err := cmdInstallLogs.MarkPersistentFlagRequired("server"); err != nil {
		log.Fatal(err)
	}
}
