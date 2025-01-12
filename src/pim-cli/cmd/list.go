package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"pim-cli/src/pim-cli/util"
)

func init() {
	listCmd.AddCommand(listTenantsCmd)
	listCmd.AddCommand(listAzureCmd)
	listCmd.AddCommand(listEntraCmd)

	listAzureCmd.AddCommand(listAzureEligibleCmd)
	listEntraCmd.AddCommand(listEntraEligibleCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

var listTenantsCmd = &cobra.Command{
	Use:   "tenants",
	Short: "List tenants that the currently signed in users is a member of",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: handle error
		tenants, _ := util.ListTenants()
		for _, tenant := range tenants {
			fmt.Printf("%s (%s)\n", *tenant.TenantID, *tenant.DisplayName)
			// TODO: mark the currently active tenant
		}
	},
}

var listAzureCmd = &cobra.Command{
	Use:   "az",
	Short: "List azure resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: print help
		fmt.Printf("Usage: pim-cli list az [eligible]\n")
	},
}

var listEntraCmd = &cobra.Command{
	Use:   "entra",
	Short: "List Entra ID roles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: print help
		fmt.Printf("Usage: pim-cli list entra [eligible]\n")
	},
}

var listAzureEligibleCmd = &cobra.Command{
	Use:   "eligible",
	Short: "List eligibility for azure resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tenantId := util.GetTenantId()
		accessToken := util.GetAccessToken(tenantId, util.AZ_MGMT_SCOPE)
		fmt.Printf("token: %s\n", accessToken)

	},
}

var listEntraEligibleCmd = &cobra.Command{
	Use:   "eligible",
	Short: "List eligibility for Entra ID roles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		eligibleRoles, err := util.ListEntraIdEligibleRoles()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		// fmt.Printf("Found %d eligible roles\n", len(*eligibleRoles))

		for _, role := range *eligibleRoles {
			roleDefinition, _ := util.GetRoleDefinitionByID(role.RoleDefinitionID)
			// fmt.Printf("%s\n%s", roleDefinition.DisplayName, roleDefinition.Description)
			fmt.Printf("%s\nPrincipal: %s\nRoleId: %s", roleDefinition.DisplayName, role.PrincipalID, role.RoleDefinitionID, roleDefinition.ID)
		}
	},
}
