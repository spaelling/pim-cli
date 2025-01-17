package util

import (
	"encoding/json"
	"fmt"
	"time"

	// "io"
	"log"
	"os/exec"
)

/*
functions that utilize the mgc cli
*/

func MgcWhoami() (*UserProfile, error) {
	cmd := exec.Command("mgc", "users", "get", "--user-id", "me")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to login: %v\n", err)
	}

	var user UserProfile
	err = json.Unmarshal(output, &user)
	if err != nil {
		log.Fatalf("failed to unmarshal output: %v\n", err)
	}

	return &user, nil
}

func GetRoleDefinitionByID(roleDefinitionID string) (*DirectoryRole, error) {
	roles, err := ListEntraIdRoleDefinitions()
	if err != nil {
		log.Fatalf("failed to get roles: %v\n", err)
		return nil, err
	}

	for _, role := range *roles {
		if role.RoleTemplateID == roleDefinitionID {
			return &role, nil
		}
	}
	return nil, fmt.Errorf("role with RoleDefinitionID %s not found", roleDefinitionID)
}

func ListEntraIdRoleDefinitions() (*[]DirectoryRole, error) {
	cmd := exec.Command("mgc", "directory-roles", "list")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to login: %v\n", err)
	}

	var directoryRolesResponse DirectoryRolesResponse
	err = json.Unmarshal(output, &directoryRolesResponse)
	if err != nil {
		log.Fatalf("failed to unmarshal output: %v\n", err)
	}

	return &directoryRolesResponse.Roles, nil
}

func ListEntraIdEligibleRoles() (*[]RoleEligibilityScheduleInstance, error) {
	cmd := exec.Command("mgc", "role-management", "directory", "role-eligibility-schedule-instances", "filter-by-current-user-with-on", "get", "--on", "principal")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to login: %v\n", err)
	}

	var roleEligibilityResponse RoleEligibilityScheduleInstancesResponse
	err = json.Unmarshal(output, &roleEligibilityResponse)
	if err != nil {
		log.Fatalf("failed to unmarshal output: %v\n", err)
	}

	// for _, role := range roleEligibilityResponse.RoleEligibilityScheduleInstances {
	// 	fmt.Print(role.String())
	// }

	return &roleEligibilityResponse.RoleEligibilityScheduleInstances, nil
}

var roleActivationRequestDefaults = RoleActivationRequest{
	Action:           "selfActivate",
	DirectoryScopeId: "/",
	IsValidationOnly: false,
	ScheduleInfo: struct {
		StartDateTime string `json:"StartDateTime"`
		Expiration    struct {
			Type     string `json:"Type"`
			Duration string `json:"Duration"`
		} `json:"Expiration"`
	}{
		Expiration: struct {
			Type     string `json:"Type"`
			Duration string `json:"Duration"`
		}{
			Type: "AfterDuration",
		},
	},
}

func ActivateEntraIdEligibleRoles(roleDefinitionId string, justification string, duration string, isValidation bool) (*RoleActivationRequest, error) {
	roleActivationRequest := roleActivationRequestDefaults

	// Get the current user
	me, err := MgcWhoami()
	if err != nil {
		return nil, err
	}

	principalId := me.ID
	startDateTime := time.Now().UTC().Format(time.RFC3339) // Set to current UTC date-time

	roleActivationRequest.PrincipalId = principalId
	roleActivationRequest.RoleDefinitionId = roleDefinitionId
	roleActivationRequest.ScheduleInfo.StartDateTime = startDateTime
	roleActivationRequest.Justification = justification
	roleActivationRequest.ScheduleInfo.Expiration.Duration = duration
	roleActivationRequest.IsValidationOnly = isValidation

	// fmt.Print(roleActivationRequest.String())

	roleActivationRequestBytes, err := json.Marshal(roleActivationRequest)
	if err != nil {
		log.Fatalf("failed to marshal body: %v\n", err)
	}
	bodyString := string(roleActivationRequestBytes)
	cmd := exec.Command("mgc", "role-management", "directory", "role-assignment-schedule-requests", "create", "--body", bodyString)
	output, err := cmd.Output()
	if err != nil {
		// TODO: exit status 255 - role is already activated
		log.Fatalf("failed to activate: %v\n", err)
	}

	var roleActivationRequestResponse RoleActivationRequest
	err = json.Unmarshal(output, &roleActivationRequestResponse)
	if err != nil {
		log.Fatalf("failed to unmarshal output: %v\n", err)
	}

	return &roleActivationRequestResponse, nil
}
