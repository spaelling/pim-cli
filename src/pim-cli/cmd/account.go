package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"pim-cli/src/pim-cli/util"
)

func init() {

	accountCmd.AddCommand(accountShowCmd)
	accountCmd.AddCommand(accountListCmd)
	accountListCmd.AddCommand(accountListUsersCmd)
	accountListCmd.AddCommand(accountListTenantsCmd)
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

var accountShowCmd = &cobra.Command{
	Use:   "show",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		accounts, _ := util.AzAccount("show")
		account := (*accounts)[0]
		fmt.Print(account.String())
	},
}

var accountListCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		accounts, _ := util.AzAccount("list")
		account := (*accounts)[0]
		fmt.Print(account.String())
	},
}

var accountListUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		accounts, _ := util.AzAccount("list")
		for _, user := range *accounts.GetUsers() {
			fmt.Printf("%s\n", user)
		}
	},
}

var accountListTenantsCmd = &cobra.Command{
	Use:   "tenants",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		accounts, _ := util.AzAccount("list")
		for _, tenant := range *accounts.GetTenants() {
			fmt.Printf("%s\n", tenant)
		}
	},
}
