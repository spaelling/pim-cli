package cmd

import (
	"fmt"

	"github.com/spaelling/pim-cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	// rootCmd.AddCommand(authCmd)
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "get the currently signed in user's profile",
	Long: `Uses Microsoft Graph API to get the currently signed in user's profile.
Usage: pim whoami`,
	Run: func(cmd *cobra.Command, args []string) {
		// this is the user logged into az cli
		user, error := util.AzWhoami()
		if error != nil {
			fmt.Printf("failed to get user: %v\n", error)
			return
		}
		fmt.Printf("az cli context: %s (%s)\n", (*user).ID, (*user).DisplayName)
		// this is the user logged into mgc cli
		user, _ = util.MgcWhoami()
		fmt.Printf("mgc cli context: %s (%s)\n", (*user).ID, (*user).DisplayName)
	},
}
