# Roles

Manage user roles and permissions.

## List All Roles

```go
roles, err := client.GetRoles(ctx)
if err != nil {
    log.Fatal(err)
}

for _, role := range roles {
    fmt.Printf("%s (ID: %d)\n", role.Name, role.ID)
}
```

## Find a Role by Name

```go
role, err := client.GetRoleByName(ctx, "Admin")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found role: %s\n", role.Name)
```

## Find a Role by ID

```go
role, err := client.GetRoleByID(ctx, "5")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found role: %s\n", role.Name)
```

## Create a New Role

```go
newRole := armis.RoleSettings{
    Name: "Security Analyst",
    Permissions: armis.Permissions{
        Device: armis.Device{
            Read: armis.Permission{All: true},
        },
        Alert: armis.Alert{
            Read: armis.Permission{All: true},
        },
    },
}

success, err := client.CreateRole(ctx, newRole)
if err != nil {
    log.Fatal(err)
}

if success {
    fmt.Println("Role created")
}
```

## Update a Role

```go
updates := armis.RoleSettings{
    Name: "Senior Security Analyst",
    Permissions: armis.Permissions{
        Device: armis.Device{
            All: true,
        },
    },
}

updated, err := client.UpdateRole(ctx, updates, "5")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated role: %s\n", updated.Name)
```

## Delete a Role

```go
success, err := client.DeleteRole(ctx, "5")
if err != nil {
    log.Fatal(err)
}

if success {
    fmt.Println("Role deleted")
}
```

## Permission Categories

Roles have permissions across these categories:

- **Device** - View and manage devices
- **Alert** - View and resolve alerts
- **Policy** - View and manage policies
- **Report** - View, export, and manage reports
- **Vulnerability** - View and manage vulnerabilities
- **Settings** - Access to various system settings
- **User** - Manage users
- **RiskFactor** - View and manage risk factors
