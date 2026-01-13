<a href="https://www.armis.com">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://www.armis.com/wp-content/uploads/2024/02/armis-logo-white.svg">
    <source media="(prefers-color-scheme: light)" srcset="https://www.armis.com/wp-content/uploads/2020/03/Armis-Logo-1.svg">
    <img src="https://www.armis.com/wp-content/uploads/2020/03/Armis-Logo-1.svg" alt="Armis logo" title="Armis" align="right" height="50">
  </picture>
</a>

# Armis Go SDK

A Go client library for the [Armis Centrix](https://www.armis.com/) API. This SDK provides a simple, idiomatic, and thread-safe way to interact with the Armis platform.

## Requirements

- [Go](https://golang.org/doc/install) >= 1.22

## Installation

```sh
go get github.com/1898andCo/armis-sdk-go
```

## Authentication

The SDK uses API key authentication. To obtain your API key from the Armis console:

1. Go to **Settings > API Management**
2. Click **Create** to create a new API key (if one doesn't exist)
3. Click **Show** to access the secret key and copy it

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/1898andCo/armis-sdk-go/armis"
)

func main() {
    client, err := armis.NewClient(
        os.Getenv("ARMIS_API_KEY"),
        armis.WithAPIURL("https://your-instance.armis.com"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    results, err := client.GetSearch(ctx, "in:devices", true, true)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d devices\n", results.Total)
}
```

## Configuration Options

```go
client, err := armis.NewClient(
    apiKey,
    armis.WithAPIURL("https://your-instance.armis.com"),  // Custom API URL
    armis.WithAPIVersion("v1"),                           // API version (default: v1)
    armis.WithHTTPClient(&http.Client{                    // Custom HTTP client
        Timeout: 60 * time.Second,
    }),
)
```

## Contributing

Check out our [Contributing Guide](./CONTRIBUTING.md) for information on how to contribute to the SDK.

For bug reports and feature requests, please use the [issue tracker](https://github.com/1898andCo/armis-sdk-go/issues).

PRs are welcome! We follow the typical "fork-and-pull" Git workflow:

1. **Fork** the repo on GitHub
2. **Clone** the project to your own machine
3. **Commit** changes to your own branch
4. **Push** your work back up to your fork
5. Submit a **Pull Request** so that we can review your changes

> [!TIP]
> Be sure to merge the latest changes from "upstream" before making a pull request!
