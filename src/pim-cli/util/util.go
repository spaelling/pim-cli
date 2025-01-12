package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	// "io"
	"log"
	"net/http"
	"os/exec"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

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

type AzureAccounts []AzureAccount

func (a AzureAccount) String() string {
	return fmt.Sprintf("EnvironmentName: %s\nHomeTenantId: %s\nID: %s\nIsDefault: %t\nManagedByTenants: %v\nName: %s\nState: %s\nTenantId: %s\nUser: %v\n",
		a.EnvironmentName, a.HomeTenantId, a.ID, a.IsDefault, a.ManagedByTenants, a.Name, a.State, a.TenantId, a.User.Name)
}

func (a AzureAccounts) String() string {
	var str string
	for _, account := range a {
		str += account.String()
		str += "\n"
	}
	return str
}

func GetCredential() *azidentity.AzureCLICredential {
	// Create a new Azure CLI credential
	credentialOpts := azidentity.AzureCLICredentialOptions{
		AdditionallyAllowedTenants: []string{"*"},
	}
	credential, err := azidentity.NewAzureCLICredential(&credentialOpts)
	if err != nil {
		log.Fatalf("failed to get credential: %v\n", err)
	}
	return credential
}

func GetAccessToken(tenantId string, scope string) string {

	credential := GetCredential()

	tokenOpts := policy.TokenRequestOptions{
		Scopes: []string{
			scope,
		},
		TenantID: tenantId,
	}
	token, err := credential.GetToken(context.Background(), tokenOpts)
	if err != nil {
		log.Fatalf("failed to get token: %v\n", err)
	}
	log.Printf("got token, expires on:\n%v\n", token.ExpiresOn)

	return token.Token
}

// This does much the same as GetTenants(), but it is for the currently logged in user and also shows the tenant displayname
func ListTenants() ([]armsubscriptions.TenantIDDescription, error) {

	var tenants []armsubscriptions.TenantIDDescription
	credential := GetCredential()

	client, err := armsubscriptions.NewTenantsClient(credential, nil)
	if err != nil {
		// log.Fatalf("failed to create tenants client: %v\n", err)
		return tenants, err
	}

	pager := client.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			log.Fatalf("failed to get next page of tenants: %v\n", err)
		}
		for _, tenant := range page.Value {
			tenants = append(tenants, *tenant)
		}
	}

	return tenants, nil
}

func AzLogin(tenantId string) {
	cmd := exec.Command("az", "login", "--tenant", tenantId, "--allow-no-subscriptions")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to login: %v\n", err)
	}
}

func GetTenantId() string {
	accounts, _ := AzAccount("show")
	account := (*accounts)[0]
	return account.GetTenantId()
}

func (a AzureAccount) GetTenantId() string {
	return a.TenantId
}

func (a AzureAccounts) GetTenants() *[]string {
	// Create a set to store unique tenant ids
	tenantSet := make(map[string]struct{})

	// Iterate over each Azure account and add the tenant id to the set
	for _, account := range a {
		tenantSet[account.GetTenantId()] = struct{}{}
	}

	// Initialize a slice to store all unique tenant ids
	var tenants []string

	// Iterate over the set and append each tenant id to the slice
	for tenant := range tenantSet {
		tenants = append(tenants, tenant)
	}

	// Return the slice of unique tenant ids
	return &tenants
}

func (a AzureAccount) GetUser() string {
	return a.User.Name
}

func (a AzureAccounts) GetUsers() *[]string {
	// Create a set to store unique user names
	userSet := make(map[string]struct{})

	// Iterate over each Azure account and add the user name to the set
	for _, account := range a {
		userSet[account.GetUser()] = struct{}{}
	}

	// Initialize a slice to store all unique user names
	var users []string

	// Iterate over the set and append each user name to the slice
	for user := range userSet {
		users = append(users, user)
	}

	// Return the slice of unique user names
	return &users
}

func AzAccountShow() (*AzureAccounts, error) {
	var account AzureAccount = AzureAccount{}
	var accounts AzureAccounts = AzureAccounts{}

	cmd := exec.Command("az", "account", "show")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("failed to run az account: %v\n", err)
		return &accounts, err
	}

	err = json.Unmarshal(output, &account)
	accounts = append(accounts, account)

	// log.Printf("%v\n", account.String())
	if err != nil {
		log.Printf("failed to unmarshal output: %v\n", err)
		return &accounts, err
	}

	return &accounts, nil
}

func AzAccountList() (*AzureAccounts, error) {
	var account AzureAccounts

	cmd := exec.Command("az", "account", "list")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("failed to run az account: %v\n", err)
		return &account, err
	}

	err = json.Unmarshal(output, &account)
	if err != nil {
		log.Printf("failed to unmarshal output: %v\n", err)
		return &account, err
	}

	return &account, nil
}

func AzAccount(flag string) (*AzureAccounts, error) {
	switch flag {
	case "show":
		return AzAccountShow()
	case "list":
		return AzAccountList()
	default:
		text := fmt.Sprintf("unknown flag %s\n", flag)
		return nil, errors.New(text)
	}
}

func MsGraphRequest(url string) (*http.Response, error) {
	// get tenant id for the currently logged in user
	tenantId := GetTenantId()
	// get access token for the currently logged in user
	token := GetAccessToken(tenantId, AZ_MGMT_SCOPE)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

func Whoami() (*UserProfile, error) {
	url := "https://graph.microsoft.com/v1.0/me"
	resp, err := MsGraphRequest(url)

	if err != nil {
		log.Fatalf("failed to get user: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to get user: %v\n", resp.Status)
		return nil, err
	}

	var profile UserProfile = UserProfile{}
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		log.Fatalf("failed to decode response: %v\n", err)
		return nil, err
	}
	return &profile, nil
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

type RoleEligibilityScheduleInstancesResponse struct {
	OdataContext                     string                            `json:"@odata.context"`
	RoleEligibilityScheduleInstances []RoleEligibilityScheduleInstance `json:"value"`
}

func (r RoleEligibilityScheduleInstance) String() string {
	return fmt.Sprintf("ID: %s\nPrincipalID: %s\nRoleDefinitionID: %s\n",
		r.ID, r.PrincipalID, r.RoleDefinitionID)
}

type DirectoryRole struct {
	ID              string      `json:"id"`
	DeletedDateTime interface{} `json:"deletedDateTime"`
	Description     string      `json:"description"`
	DisplayName     string      `json:"displayName"`
	RoleTemplateID  string      `json:"roleTemplateId"`
}

type DirectoryRolesResponse struct {
	OdataContext string          `json:"@odata.context"`
	Roles        []DirectoryRole `json:"value"`
}

func (r DirectoryRole) String() string {
	return fmt.Sprintf("ID: %s\nDeletedDateTime: %v\nDescription: %s\nDisplayName: %s\nRoleTemplateID: %s\n",
		r.ID, r.DeletedDateTime, r.Description, r.DisplayName, r.RoleTemplateID)
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
