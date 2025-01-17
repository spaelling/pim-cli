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
	activateEntraEligibleCmd.Flags().BoolP("validation", "z", false, "validate the role activation request")
}

// MARK: activate
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

// MARK: activate az
var activateAzureCmd = &cobra.Command{
	Use:   "az",
	Short: "",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

// MARK: activate entra
var activateEntraCmd = &cobra.Command{
	Use:   "entra",
	Short: "",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

// MARK: activate az eligible
var activateAzureEligibleCmd = &cobra.Command{
	Use:   "eligible",
	Short: "Activate eligibility for azure resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Not implemented yet\n")
	},
}

// MARK: activate entra roles
var activateEntraEligibleCmd = &cobra.Command{
	Use:   "eligible",
	Short: "Activate eligible Entra ID role",
	Long: `
usage example:

validate role activation:

pim activate entra eligible --role "58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d" --justification "Enable admin role using pim cli" --duration 42 --validation true

Role activation is validated!
Action:                 selfActivate
Principal Id:           65cfef21-f882-40a8-acc4-e00eeb156088
RoleDefinition Id:      58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d
Is validation:          true
Justification:          Enable admin role using pim cli
Activation start:       2025-01-17T19:59:30Z
Expires after:          PT42M
`,

	Run: func(cmd *cobra.Command, args []string) {
		// get roleid, justification, and duration from flags
		// return if we encounter an error
		roleid, err := cmd.Flags().GetString("role")
		if err != nil {
			fmt.Printf("Error getting role id: %s\n", err)
			return
		}
		justification, err := cmd.Flags().GetString("justification")
		if err != nil {
			fmt.Printf("Error getting justification: %s\n", err)
			return
		}
		durationMinutes, err := cmd.Flags().GetInt("duration")
		if err != nil {
			fmt.Printf("Error getting duration: %s\n", err)
			return
		}
		if durationMinutes > 8*60 {
			fmt.Printf("Error: duration must be less than 8 hours/480 minutes\n")
			return
		}
		// convert duration to PT{hours}H{minutes}M format
		duration := fmt.Sprintf("PT%dH%dM", durationMinutes/60, durationMinutes%60)
		// IsValidationOnly
		isValidation, err := cmd.Flags().GetBool("validation")
		if err != nil {
			fmt.Printf("Error getting validationflag: %s\n", err)
			return
		}

		// activate the role
		roleActivationRequestResponse, err := util.ActivateEntraIdEligibleRoles(roleid, justification, duration, isValidation)
		if err != nil {
			fmt.Printf("Error when activating role: %s\n", err)
		}

		if isValidation {
			fmt.Print("Role activation is validated!\n")
		} else {
			fmt.Print("Role activation is succesfull!\n")
		}

		// TODO: look up name of the role

		// print the response
		fmt.Printf("Action:\t\t\t%s\n", roleActivationRequestResponse.Action)
		fmt.Printf("Principal Id:\t\t%s\n", roleActivationRequestResponse.PrincipalId)
		fmt.Printf("RoleDefinition Id:\t%s\n", roleActivationRequestResponse.RoleDefinitionId)
		fmt.Printf("Is validation:\t\t%t\n", roleActivationRequestResponse.IsValidationOnly)
		fmt.Printf("Justification:\t\t%s\n", roleActivationRequestResponse.Justification)
		fmt.Printf("Activation start:\t%s\n", roleActivationRequestResponse.ScheduleInfo.StartDateTime)
		fmt.Printf("Expires after:\t\t%v\n", roleActivationRequestResponse.ScheduleInfo.Expiration.Duration)

	},
}
