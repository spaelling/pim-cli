cls
# alias pim="./pim"
go build -o pim .
./pim activate entra eligible --help
./pim activate entra eligible --role "58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d" --justification "Enable admin role using pim cli" --duration 33 --validation true

"pim list entra roles"
./pim list entra roles

./pim whoami



"pim list entra eligible"
./pim list entra eligible

# ./pim whoami -h


"pim version"
./pim v
"pim account show"
./pim account show
"pim account list users"
./pim account list users
"pim account list tenants"
./pim account list tenants
"pim list tenants"
./pim list tenants
"pim login contoso"
./pim login "1b775964-7849-4f1a-8052-60b8e5c59b96"

az account get-access-token --scope "https://graph.microsoft.com/RoleEligibilitySchedule.Read.Directory"

az login --scope "https://graph.microsoft.com/RoleEligibilitySchedule.Read.Directory"

az account get-access-token --help

mgc login --help

mgc login --strategy InteractiveBrowser --scopes "RoleAssignmentSchedule.ReadWrite.Directory" #"RoleEligibilitySchedule.Read.Directory", 

mgc identity-governance privileged-access get --help
mgc identity-governance privileged-access get --select id --debug

# application developer cf1c38e5-3621-4004-a7cb-879624dced7c
# Attribute Assignment Administrator 58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d
$body = @"
{
    "Action": "selfActivate",
    "PrincipalId": "65cfef21-f882-40a8-acc4-e00eeb156088",
    "RoleDefinitionId": "58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d",
    "DirectoryScopeId": "/",
    "isValidationOnly": false,
    "Justification": "Enable admin role",
    "ScheduleInfo": {
        "StartDateTime": "$(Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ')",
        "Expiration": {
            "Type": "AfterDuration",
            "Duration": "PT4H"
        }
    }
}
"@

$body | ConvertFrom-Json

mgc role-management directory role-assignment-schedule-requests create --body $Body
mgc role-management directory role-assignment-schedule-requests  filter-by-current-user-with-on get --on principal

mgc role-management directory role-eligibility-schedule-requests create --body $Body --debug

mgc  role-management directory role-eligibility-schedule-instances filter-by-current-user-with-on get --on principal

cls; mgc directory-roles list # --directory-role-id 'cf1c38e5-3621-4004-a7cb-879624dced7c'

mgc directory-roles get