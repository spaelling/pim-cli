# PIM CLI

PIM CLI written in Golang, aiming to trivialize PIM role listing and activation, especially for multi-tenant and multi-user scenarios.

## Requirements

The PIM CLI utilizes [Az CLI](https://learn.microsoft.com/en-us/cli/azure/) and the less known [mgc CLI](https://learn.microsoft.com/en-us/graph/cli/installation?tabs=macos).
Both are developed and supported by Microsoft. The reason that we use this is due to how authentication works.
The Az CLI has commands to expose an authentication token but has a limited scope in Microsoft Graph.
The Microsoft Graph CLI (mgc) does not expose a token (yet), but it is fairly convoluted and looks to be auto-generated based on the Graph REST APIs.

## Installation

TODO

## Commands

run `pim --help` for a list of command. The help sections are mostly unfinished. My focus is on documenting the most useful commands.

### pim list

Commands for listing different areas. This is in the context of the currently logged in user.

#### pim list tenants

List tenants that the currently signed in users is a member of

```cli
pim list tenants

Name                                          ID                                           
Acme Corporation                            123e4567-e89b-12d3-a456-426614174000         
Globex Corporation                          234e5678-f90c-23d4-b567-526715274001         
Soylent Corp                                345f6789-0a1b-34e5-c678-627816374002         
Initech                                     456g7890-1b2c-45f6-d789-728917474003         
Umbrella Corporation                        567h8901-2c3d-56g7-e890-829018574004         
Hooli                                       678i9012-3d4e-67h8-f901-930119674005
```

#### pim list az eligible

`not implemented yet`

#### pim list entra roles

List all active Entra ID roles

```cli
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
```

##### pim list entra eligible

List eligible Entra ID roles for currently logged in user.

```cli
pim list entra eligible

Principal ID: 65cfef21-f882-40a8-acc4-e00eeb156088

Role:   Application Developer
Id:     cf1c38e5-3621-4004-a7cb-879624dced7c

Role:   Attribute Assignment Administrator
Id:     58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d
```

### pim activate

Activate roles in Entra ID or Azure.

#### pim activate az eligible

`not implemented yet`

#### pim activate entra eligible

```cli
pim activate entra eligible --role "58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d" --justification "Enable admin role using pim cli" --duration 42 --validation true

Role activation is validated!
Action:                 selfActivate
Principal Id:           65cfef21-f882-40a8-acc4-e00eeb156088
RoleDefinition Id:      58a13ea3-c632-46ae-9ee0-9c0d43cd7f3d
Is validation:          true
Justification:          Enable admin role using pim cli
Activation start:       2025-01-17T19:59:30Z
Expires after:          PT42M
```

### pim account

Commands to interact with the currently signed in user

## Planned Features

### Set activation

Activate multiple roles, both Entra ID, Azure, and even group memberships in one command.

Example: Directory Reader and Reader (Azure) on the pseudo-root management group for 8 hours, Owner on Subscription X, Key Vault Administrator on specific Key Vault and Storage Blob Data Contributor on 3 different Storage Accounts for 4 hours, using different justifications.

This could be the daily driver for both working with some low privileged broad access and access for working on a specific project. Activating those roles in the Azure Portal would take 10 minutes or more, each day.

## TODO

- debug and verbose output
- auto complete
- help
- release pipeline
- this README
 