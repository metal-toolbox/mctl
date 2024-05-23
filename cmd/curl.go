package cmd

// Adopted from an internal tool, credits to those unnamed authors.
import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"syscall"

	"github.com/metal-toolbox/mctl/internal/auth"
	"github.com/metal-toolbox/mctl/pkg/model"
	"github.com/spf13/cobra"
)

// curlCmd represents the curl command
var curlCmd = &cobra.Command{
	Use:   "curl fleetdbapi -- args",
	Short: "Make a curl request with your auth token",
	Args:  cobra.MinimumNArgs(2), // nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		apiKind := args[0]
		if !slices.Contains(
			[]string{string(model.FleetDBAPI)},
			apiKind,
		) {
			log.Fatal("invalid service parameter: " + apiKind)
		}

		curlArgs := args[1:]
		doCurl(cmd.Context(), model.APIKind(apiKind), curlArgs)
	},
}

func init() {
	RootCmd.AddCommand(curlCmd)
}

func doCurl(ctx context.Context, _ model.APIKind, args []string) {
	mctl := MustCreateApp(ctx)
	ctx, cancel := context.WithTimeout(ctx, CmdTimeout)
	defer cancel()

	token, err := auth.AccessToken(ctx, model.FleetDBAPI, mctl.Config.FleetDBAPI, false)
	if err != nil {
		// nolint:gocritic // its fine if the ctx is not cleaned up we're exiting the app.
		log.Fatal("auth token error: " + err.Error())
	}

	binary, lookErr := exec.LookPath("curl")
	if lookErr != nil {
		panic(lookErr)
	}

	cmd := []string{"curl", "-H", fmt.Sprintf("Authorization: Bearer %s", token)}
	cmd = append(cmd, args...)

	// syscall.Exec changes the current running process to curl.
	// this allows curl to take over the same pid and allows us to exit
	if err := syscall.Exec(binary, cmd, os.Environ()); err != nil {
		log.Fatal(err)
	}
}
