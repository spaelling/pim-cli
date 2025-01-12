package cmd

import (
	"fmt"

	"pim-cli/src/pim-cli/util"

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
		user, _ := util.Whoami()
		fmt.Printf("%s\n", (*user).String())
	},
}
