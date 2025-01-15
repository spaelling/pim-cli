package util

import (
	"fmt"
)

type DirectoryRole struct {
	ID              string      `json:"id"`
	DeletedDateTime interface{} `json:"deletedDateTime"`
	Description     string      `json:"description"`
	DisplayName     string      `json:"displayName"`
	RoleTemplateID  string      `json:"roleTemplateId"`
}

func (r DirectoryRole) String() string {
	return fmt.Sprintf("ID: %s\nDeletedDateTime: %v\nDescription: %s\nDisplayName: %s\nRoleTemplateID: %s\n",
		r.ID, r.DeletedDateTime, r.Description, r.DisplayName, r.RoleTemplateID)
}

type DirectoryRolesResponse struct {
	OdataContext string          `json:"@odata.context"`
	Roles        []DirectoryRole `json:"value"`
}

type ScheduleInfo struct {
	StartDateTime string `json:"StartDateTime"`
	Expiration    struct {
		Type     string `json:"Type"`
		Duration string `json:"Duration"`
	} `json:"Expiration"`
}

type RoleActivationRequest struct {
	Action           string       `json:"Action"`
	PrincipalId      string       `json:"PrincipalId"`
	RoleDefinitionId string       `json:"RoleDefinitionId"`
	DirectoryScopeId string       `json:"DirectoryScopeId"`
	IsValidationOnly bool         `json:"isValidationOnly"`
	Justification    string       `json:"Justification"`
	ScheduleInfo     ScheduleInfo `json:"ScheduleInfo"`
}

func (r RoleActivationRequest) String() string {
	return fmt.Sprintf("Action: %s\nPrincipalId: %s\nRoleDefinitionId: %s\nDirectoryScopeId: %s\nIsValidationOnly: %t\nJustification: %s\nScheduleInfo: StartDateTime=%s, Expiration=%v\n",
		r.Action, r.PrincipalId, r.RoleDefinitionId, r.DirectoryScopeId, r.IsValidationOnly, r.Justification, r.ScheduleInfo.StartDateTime, r.ScheduleInfo.Expiration)
}

type User struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type AzureAccount struct {
	EnvironmentName  string        `json:"environmentName"`
	HomeTenantId     string        `json:"homeTenantId"`
	ID               string        `json:"id"`
	IsDefault        bool          `json:"isDefault"`
	ManagedByTenants []interface{} `json:"managedByTenants"`
	Name             string        `json:"name"`
	State            string        `json:"state"`
	TenantId         string        `json:"tenantId"`
	User             User          `json:"user"`
}

func (a AzureAccount) String() string {
	return fmt.Sprintf("EnvironmentName: %s\nHomeTenantId: %s\nID: %s\nIsDefault: %t\nManagedByTenants: %v\nName: %s\nState: %s\nTenantId: %s\nUser: %v\n",
		a.EnvironmentName, a.HomeTenantId, a.ID, a.IsDefault, a.ManagedByTenants, a.Name, a.State, a.TenantId, a.User.Name)
}

type AzureAccounts []AzureAccount

func (a AzureAccounts) String() string {
	var str string
	for _, account := range a {
		str += account.String()
		str += "\n"
	}
	return str
}

type UserProfile struct {
	OdataContext      string   `json:"@odata.context"`
	BusinessPhones    []string `json:"businessPhones"`
	DisplayName       string   `json:"displayName"`
	GivenName         string   `json:"givenName"`
	ID                string   `json:"id"`
	JobTitle          string   `json:"jobTitle"`
	Mail              string   `json:"mail"`
	MobilePhone       string   `json:"mobilePhone"`
	OfficeLocation    string   `json:"officeLocation"`
	PreferredLanguage string   `json:"preferredLanguage"`
	Surname           string   `json:"surname"`
	UserPrincipalName string   `json:"userPrincipalName"`
}

func (u UserProfile) String() string {

	return fmt.Sprintf(`
DisplayName: %s
GivenName: %s
ID: %s
JobTitle: %s
Mail: %s
MobilePhone: %s
OfficeLocation: %s
PreferredLanguage: %s
Surname: %s
UserPrincipalName: %s`,
		u.DisplayName, u.GivenName, u.ID, u.JobTitle, u.Mail, u.MobilePhone, u.OfficeLocation, u.PreferredLanguage, u.Surname, u.UserPrincipalName)
}

func (u UserProfile) Table() string {
	return fmt.Sprintf("%-45s %-45s %-45s %-45s\n%-45s %-45s %-45s %-45s\n", "DisplayName", "ID", "Mail", "UserPrincipalName", u.DisplayName, u.ID, u.Mail, u.UserPrincipalName)
}

type RoleEligibilityScheduleInstance struct {
	ID                        string      `json:"id"`
	PrincipalID               string      `json:"principalId"`
	RoleDefinitionID          string      `json:"roleDefinitionId"`
	DirectoryScopeID          string      `json:"directoryScopeId"`
	AppScopeID                interface{} `json:"appScopeId"`
	StartDateTime             string      `json:"startDateTime"`
	EndDateTime               interface{} `json:"endDateTime"`
	MemberType                string      `json:"memberType"`
	RoleEligibilityScheduleID string      `json:"roleEligibilityScheduleId"`
}

func (r RoleEligibilityScheduleInstance) String() string {
	return fmt.Sprintf("ID: %s\nPrincipalID: %s\nRoleDefinitionID: %s\n",
		r.ID, r.PrincipalID, r.RoleDefinitionID)
}

type RoleEligibilityScheduleInstancesResponse struct {
	OdataContext                     string                            `json:"@odata.context"`
	RoleEligibilityScheduleInstances []RoleEligibilityScheduleInstance `json:"value"`
}
