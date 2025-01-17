package cmd

import (
	"fmt"

	"github.com/spaelling/pim-cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listTenantsCmd)
	listCmd.AddCommand(listAzureCmd)
	listCmd.AddCommand(listEntraCmd)

	listAzureCmd.AddCommand(listAzureEligibleCmd)
	listEntraCmd.AddCommand(listEntraEligibleCmd)
	listEntraCmd.AddCommand(listEntraRolesCmd)
}

// MARK: list
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

// MARK: list tenants
var listTenantsCmd = &cobra.Command{
	Use:   "tenants",
	Short: "List tenants that the currently signed in users is a member of",
	Long: `
	
Usage example:
pim list tenants

Name                                          ID                                           
Acme Corporation                              123e4567-e89b-12d3-a456-426614174000         
Globex Corporation                            234e5678-f90c-23d4-b567-526715274001         
Soylent Corp                                  345f6789-0a1b-34e5-c678-627816374002         
Initech                                       456g7890-1b2c-45f6-d789-728917474003         
Umbrella Corporation                          567h8901-2c3d-56g7-e890-829018574004         
Hooli                                         678i9012-3d4e-67h8-f901-930119674005
`,
	Run: func(cmd *cobra.Command, args []string) {
		tenants, err := util.ListTenants()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		fmt.Printf("%-45s %-45s\n", "Name", "ID")
		for _, tenant := range tenants {
			fmt.Printf("%-45s %-45s\n", *tenant.DisplayName, *tenant.TenantID)
			// TODO: mark the currently active tenant
		}
	},
}

// MARK: list az
var listAzureCmd = &cobra.Command{
	Use:   "az",
	Short: "",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

// MARK: list entra
var listEntraCmd = &cobra.Command{
	Use:   "entra",
	Short: "",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

// MARK: list az eligible
var listAzureEligibleCmd = &cobra.Command{
	Use:   "eligible",
	Short: "List eligibility for azure resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Not implemented yet\n")
	},
}

// MARK: list entra roles
var listEntraRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List Entra ID roles",
	Long: `
List all active Entra ID roles

Usage example:
pim list entra roles

Name                                          ID                                           
Directory Readers                             88d8e3e3-8f55-4a1e-953a-9b9898b8876b         
Azure AD Joined Device Local Administrator    9f06204d-73c1-4d4c-880a-6edb90606fd8         
Application Developer                         cf1c38e5-3621-4004-a7cb-879624dced7c         
Attribute Assignment Administrator            58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d         
Teams Administrator                           69091246-20e8-4a56-aa4d-066075b2a7a8         
Attribute Assignment Reader                   ffd52fa5-98dc-465c-991d-fc073eb59f8f         
Exchange Administrator                        29232cdf-9323-42fd-ade2-1d097af3e4de         
Application Administrator                     9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3         
Knowledge Administrator                       b5a8dcf3-09d5-43a9-a639-8e29ef291470         
Global Administrator                          62e90394-69f5-4237-9190-012177145e10  
`,
	Run: func(cmd *cobra.Command, args []string) {
		roleDefinitions, err := util.ListEntraIdRoleDefinitions()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		// fmt.Printf("Found %d roles\n", len(*roleDefinitions))

		fmt.Printf("%-45s %-45s\n", "Name", "ID")
		for _, role := range *roleDefinitions {
			// only care about the DisplayName and RoleTemplateID
			fmt.Printf("%-45s %-45s\n", role.DisplayName, role.RoleTemplateID)
		}
	},
}

// MARK: list entra eligible
var listEntraEligibleCmd = &cobra.Command{
	Use:   "eligible",
	Short: "List eligible Entra ID roles",
	Long: `
List eligible Entra ID roles

Usage example:
pim list entra eligible

Principal ID: 65cfef21-f882-40a8-acc4-e00eeb156088

Role:   Application Developer
Id:     cf1c38e5-3621-4004-a7cb-879624dced7c

Role:   Attribute Assignment Administrator
Id:     58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d
`,
	Run: func(cmd *cobra.Command, args []string) {
		eligibleRoles, err := util.ListEntraIdEligibleRoles()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		if len(*eligibleRoles) == 0 {
			fmt.Println("No eligible roles found")
			return
		}
		fmt.Printf("Principal ID: %s\n\n", (*eligibleRoles)[0].PrincipalID)

		for _, role := range *eligibleRoles {
			roleDefinition, _ := util.GetRoleDefinitionByID(role.RoleDefinitionID)
			fmt.Printf("Role:\t%s\n", roleDefinition.DisplayName)
			fmt.Printf("Id:\t%s\n", roleDefinition.RoleTemplateID)
			fmt.Print("\n")
		}
	},
}
