# Getting Started

## Installation

```bash
go get github.com/1898andCo/armis-sdk-go/v2
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/1898andCo/armis-sdk-go/v2/armis"
)

func main() {
    client, err := armis.NewClient(
        os.Getenv("ARMIS_API_KEY"),
        os.Getenv("ARMIS_API_URL"), // e.g., "https://your-instance.armis.com"
    )
    if err != nil {
        log.Fatal(err)
    }

    users, err := client.GetUsers(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d users\n", len(users))
}
```

> **Note:** Set `ARMIS_API_URL` to your organization's Armis instance URL (e.g., `https://your-instance.armis.com`).

## Configuration Options

```go
client, err := armis.NewClient(
    os.Getenv("ARMIS_API_KEY"),                           // Required: API key
    os.Getenv("ARMIS_API_URL"),                           // Required: Your Armis instance URL
    armis.WithAPIVersion("v1"),                           // Optional: API version (default: v1)
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
