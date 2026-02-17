# Armis Go SDK

<div align="center">

[![Armis SDK Tests](https://github.com/1898andCo/armis-sdk-go/actions/workflows/sdk-tests.yml/badge.svg)](https://github.com/1898andCo/armis-sdk-go/actions/workflows/sdk-tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/1898andCo/armis-sdk-go/v2)](https://goreportcard.com/report/github.com/1898andCo/armis-sdk-go/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/1898andCo/armis-sdk-go/v2.svg)](https://pkg.go.dev/github.com/1898andCo/armis-sdk-go/v2)
[![Go Version](https://img.shields.io/github/go-mod/go-version/1898andCo/armis-sdk-go)](go.mod)
[![License](https://img.shields.io/github/license/1898andCo/armis-sdk-go)](LICENSE)

</div>

A Go client library for the [Armis Centrix](https://www.armis.com/) API. This SDK provides a simple, idiomatic, and thread-safe way to interact with the Armis platform.

## Requirements

- [Go](https://golang.org/doc/install) >= 1.25

## Installation

```sh
go get github.com/1898andCo/armis-sdk-go/v2
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
import (
    "net/http"
    "time"

    "github.com/1898andCo/armis-sdk-go/v2/armis"
)

client, err := armis.NewClient(
    apiKey,                                               // Required: API key
    "https://your-instance.armis.com",                    // Required: API URL
    armis.WithAPIVersion("v1"),                           // Optional: API version (default: v1)
    armis.WithHTTPClient(&http.Client{                    // Optional: Custom HTTP client
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

## LLM Support

This project includes an [llms.txt](./llms.txt) file following the [llms.txt specification](https://llmstxt.org) to help Large Language Models understand and work with this codebase more effectively.
