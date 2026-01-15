# Users

Manage Armis user accounts.

## List All Users

```go
users, err := client.GetUsers(ctx)
if err != nil {
    log.Fatal(err)
}

for _, user := range users {
    fmt.Printf("%s (%s)\n", user.Name, user.Email)
}
```

## Get a Single User

You can fetch a user by their ID or email:

```go
// By ID
user, err := client.GetUser(ctx, "123")

// By email
user, err := client.GetUser(ctx, "jane@example.com")

fmt.Printf("Name: %s, Role: %s\n", user.Name, user.Role)
```

## Create a New User

```go
newUser := armis.UserSettings{
    Name:  "Jane Doe",
    Email: "jane@example.com",
    Role:  "Admin",
}

created, err := client.CreateUser(ctx, newUser)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created user with ID: %d\n", created.ID)
```

## Update a User

```go
updates := armis.UserSettings{
    Name:  "Jane Smith",
    Email: "jane@example.com",
    Title: "Security Analyst",
}

updated, err := client.UpdateUser(ctx, updates, "123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated user: %s\n", updated.Name)
```

## Delete a User

```go
success, err := client.DeleteUser(ctx, "123")
if err != nil {
    log.Fatal(err)
}

if success {
    fmt.Println("User deleted")
}
```

## User Fields

| Field | Description |
|-------|-------------|
| `ID` | Unique identifier |
| `Name` | Full name |
| `Email` | Email address |
| `Username` | Login username |
| `Role` | Assigned role |
| `IsActive` | Whether the user is active |
| `Title` | Job title |
| `Phone` | Phone number |
| `Location` | Location |
| `TwoFactorAuthentication` | 2FA enabled |
| `LastLoginTime` | Last login timestamp |
