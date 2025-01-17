package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	// "io"
	"log"
	"os/exec"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

/*
functions that utilize the az cli

BEWARE: it can be different users logged in to the az cli and mgc cli
*/

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

func AzLogin(tenantId string) {
	cmd := exec.Command("az", "login", "--tenant", tenantId, "--allow-no-subscriptions")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to login: %v\n", err)
	}
}

func GetAccessToken(tenantId string, scope string) (string, error) {

	credential := GetCredential()

	tokenOpts := policy.TokenRequestOptions{
		Scopes: []string{
			scope,
		},
		TenantID: tenantId,
	}
	token, err := credential.GetToken(context.Background(), tokenOpts)
	if err != nil {
		// The refresh token has expired due to inactivity. need to reauthenticate
		if strings.Contains(err.Error(), "AADSTS700082") {
			err = errors.New(`token has expired. Please reauthenticate by running:

pim login "<TENANT_ID>"`)
		}
		return "", err
	}
	// log.Printf("got token, expires on:\n%v\n", token.ExpiresOn)

	return token.Token, nil
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

func GetTenantId() (string, error) {
	accounts, error := AzAccount("show")
	if error != nil {
		return "", error
	}
	account := (*accounts)[0]
	return account.GetTenantId(), nil
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
	tenantId, err := GetTenantId()
	if err != nil {
		return nil, err
	}
	// get access token for the currently logged in user
	token, err := GetAccessToken(tenantId, MS_GRAPH_SCOPE)
	if err != nil {
		return nil, err
	}

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

func AzWhoami() (*UserProfile, error) {
	url := "https://graph.microsoft.com/v1.0/me"
	resp, err := MsGraphRequest(url)

	if err != nil {
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
