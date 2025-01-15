package cmd

import (
	"fmt"
	"pim-cli/src/pim-cli/util"

	"github.com/spf13/cobra"
)

func init() {
	activateCmd.AddCommand(activateAzureCmd)
	activateCmd.AddCommand(activateEntraCmd)
	activateAzureCmd.AddCommand(activateAzureEligibleCmd)
	activateEntraCmd.AddCommand(activateEntraEligibleCmd)

	activateEntraEligibleCmd.Flags().StringP("role", "r", "", "role (id) to activate")
	activateEntraEligibleCmd.MarkFlagRequired("role")
	activateEntraEligibleCmd.Flags().StringP("justification", "j", "", "justification for activating the role")
	activateEntraEligibleCmd.MarkFlagRequired("justification")
	activateEntraEligibleCmd.Flags().IntP("duration", "d", 240, "duration for the role activation in minutes")
}

var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

var activateAzureCmd = &cobra.Command{
	Use:   "az",
	Short: "activate azure resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: print help
		fmt.Printf("Usage: pim-cli activate az [eligible]\n")
	},
}

var activateEntraCmd = &cobra.Command{
	Use:   "entra",
	Short: "activate Entra ID roles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: print help
		fmt.Printf("Usage: pim-cli activate entra [eligible]\n")
	},
}

var activateAzureEligibleCmd = &cobra.Command{
	Use:   "eligible",
	Short: "Activate eligibility for azure resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// tenantId := util.GetTenantId()
		// accessToken := util.GetAccessToken(tenantId, util.AZ_MGMT_SCOPE)
		// fmt.Printf("token: %s\n", accessToken)

	},
}

var activateEntraEligibleCmd = &cobra.Command{
	Use:   "eligible",
	Short: "Activate eligible Entra ID role",
	Long: `
pim activate entra eligible --role "58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d" `,

	Run: func(cmd *cobra.Command, args []string) {
		roleid, err := cmd.Flags().GetString("role")
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		justification, err := cmd.Flags().GetString("justification")
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		//"PT4H"
		durationMinutes, err := cmd.Flags().GetInt("duration")
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		duration := fmt.Sprintf("PT%dH%dM", durationMinutes/60, durationMinutes%60)
		roleActivationRequestResponse, err := util.ActivateEntraIdEligibleRoles(roleid, justification, duration)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		// TODO
		fmt.Printf("Action: %s\n", roleActivationRequestResponse.Action)
		fmt.Printf("PrincipalId: %s\n", roleActivationRequestResponse.PrincipalId)
		fmt.Printf("RoleDefinitionId: %s\n", roleActivationRequestResponse.RoleDefinitionId)
		fmt.Printf("DirectoryScopeId: %s\n", roleActivationRequestResponse.DirectoryScopeId)
		fmt.Printf("IsValidationOnly: %t\n", roleActivationRequestResponse.IsValidationOnly)
		fmt.Printf("Justification: %s\n", roleActivationRequestResponse.Justification)
		fmt.Printf("ScheduleInfo: StartDateTime=%s, Expiration=%v\n", roleActivationRequestResponse.ScheduleInfo.StartDateTime, roleActivationRequestResponse.ScheduleInfo.Expiration)

	},
}
