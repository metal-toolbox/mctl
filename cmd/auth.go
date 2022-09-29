package cmd

import (
	"fmt"
	"os"

	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Handles getting auth tokens to talk to services",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := app.New(cmd.Context(), cfgFile)
		if err != nil {
			return err
		}

		_, err = m.GetOAuth2Token(
			cmd.Context(),
			m.Config.OidcClientID,
			m.Config.OidcIssuerEndpoint,
			m.Config.OidcAudience,
		)

		if err != nil {
			fmt.Println("authentication error: " + err.Error())
			os.Exit(1)
		}

		fmt.Println("authentication successful, auth token stored in keyring.")
		return nil

	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
