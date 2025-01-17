package cmd

import (
	"fmt"

	"github.com/spaelling/pim-cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	// rootCmd.AddCommand(authCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to Entra ID",
	Long:  `Authenticate to Entra ID`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("Getting credential")
		tenantId := args[0]
		util.AzLogin(tenantId)
		accounts, _ := util.AzAccount("show")
		account := (*accounts)[0]
		fmt.Print(account.String())
	},
}
