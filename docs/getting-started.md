# Getting Started

## Installation

```bash
go get github.com/1898andCo/armis-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/1898andCo/armis-sdk-go/armis"
)

func main() {
    client, err := armis.NewClient(os.Getenv("ARMIS_API_KEY"))
    if err != nil {
        panic(err)
    }

    // You're ready to go!
    users, _ := client.GetUsers(context.Background())
    fmt.Printf("Found %d users\n", len(users))
}
```

## Using a Custom Armis Instance

```go
client, err := armis.NewClient(
    os.Getenv("ARMIS_API_KEY"),
    armis.WithAPIURL("https://your-company.armis.com"),
)
```

## Next Steps

- [Users](users.md) - Manage user accounts
- [Roles](roles.md) - Configure role permissions
- [Policies](policies.md) - Create and manage security policies
- [Collectors](collectors.md) - Manage data collectors
- [Search](search.md) - Query devices and alerts with AQL
- [Boundaries](boundaries.md) - View network boundaries
- [Reports](reports.md) - Access scheduled reports
- [Sites](sites.md) - View site information
- [Lists](lists.md) - View custom lists
