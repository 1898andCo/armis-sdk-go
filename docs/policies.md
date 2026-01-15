# Policies

Create and manage security policies.

## List All Policies

```go
policies, err := client.GetAllPolicies(ctx)
if err != nil {
    log.Fatal(err)
}

for _, policy := range policies {
    status := "disabled"
    if policy.IsEnabled {
        status = "enabled"
    }
    fmt.Printf("%s (%s) - %s\n", policy.Name, policy.RuleType, status)
}
```

## Get a Single Policy

```go
policy, err := client.GetPolicy(ctx, "123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Policy: %s\n", policy.Name)
fmt.Printf("Description: %s\n", policy.Description)
fmt.Printf("Type: %s\n", policy.RuleType)
```

## Create a Policy

```go
policy := armis.PolicySettings{
    Name:        "Detect Unauthorized Devices",
    Description: "Alert when unmanaged devices connect to the network",
    RuleType:    "DEVICE",
    IsEnabled:   true,
    Rules: armis.Rules{
        And: []any{
            map[string]any{
                "field":    "managedBy",
                "operator": "equals",
                "value":    "none",
            },
        },
    },
    Actions: []armis.Action{
        {
            Type: "alert",
            Params: armis.Params{
                Severity: "high",
                Title:    "Unmanaged Device Detected",
            },
        },
    },
}

result, err := client.CreatePolicy(ctx, policy)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created policy with ID: %d\n", result.ID)
```

## Update a Policy

```go
updates := armis.PolicySettings{
    Name:        "Detect Unauthorized Devices (Updated)",
    Description: "Updated description",
    RuleType:    "DEVICE",
    IsEnabled:   false,
    Rules: armis.Rules{
        And: []any{
            map[string]any{
                "field":    "managedBy",
                "operator": "equals",
                "value":    "none",
            },
        },
    },
}

updated, err := client.UpdatePolicy(ctx, updates, "123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated: %s\n", updated.Name)
```

## Delete a Policy

```go
success, err := client.DeletePolicy(ctx, "123")
if err != nil {
    log.Fatal(err)
}

if success {
    fmt.Println("Policy deleted")
}
```

## Policy Rule Types

| Type | Description |
|------|-------------|
| `DEVICE` | Triggers based on device attributes |
| `ACTIVITY` | Triggers based on device activity |
| `IP_CONNECTION` | Triggers based on network connections |
| `VULNERABILITY` | Triggers based on vulnerabilities |

## Validation Rules

- Name is required
- Rules cannot be empty (must have `And` or `Or` conditions)
- Description must be under 500 characters
- RuleType must be one of the valid types above
