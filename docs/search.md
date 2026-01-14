# Search

Query Armis data using AQL (Armis Query Language).

## Basic Search

```go
results, err := client.GetSearch(ctx, "in:devices", true, true)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d results (total: %d)\n", results.Count, results.Total)
```

## Search Parameters

```go
// GetSearch(ctx, aql, includeSample, includeTotal)
// - aql: The AQL query string
// - includeSample: Include sample results in response
// - includeTotal: Include total count in response

results, err := client.GetSearch(ctx, "in:devices", true, true)
```

## Search for Devices

```go
// All devices
results, _ := client.GetSearch(ctx, "in:devices", true, true)

// Windows devices
results, _ := client.GetSearch(ctx, "in:devices os:windows", true, true)

// High risk devices
results, _ := client.GetSearch(ctx, "in:devices riskLevel:>7", true, true)

// Devices seen in last 24 hours
results, _ := client.GetSearch(ctx, "in:devices lastSeen:>now-1d", true, true)
```

## Search for Alerts

```go
// All alerts
results, _ := client.GetSearch(ctx, "in:alerts", true, true)

// High severity alerts
results, _ := client.GetSearch(ctx, "in:alerts severity:high", true, true)

// Unresolved alerts
results, _ := client.GetSearch(ctx, "in:alerts status:unresolved", true, true)
```

## Search for Vulnerabilities

```go
// All vulnerabilities
results, _ := client.GetSearch(ctx, "in:vulnerabilities", true, true)

// Critical vulnerabilities
results, _ := client.GetSearch(ctx, "in:vulnerabilities severity:critical", true, true)
```

## Working with Results

```go
results, err := client.GetSearch(ctx, "in:alerts severity:high", true, true)
if err != nil {
    log.Fatal(err)
}

for _, item := range results.Results {
    fmt.Printf("Alert: %s\n", item.Title)
    fmt.Printf("  Severity: %s\n", item.Severity)
    fmt.Printf("  Status: %s\n", item.Status)
    fmt.Printf("  Time: %s\n", item.Time)
}
```

## Pagination

Results include pagination info:

```go
results, _ := client.GetSearch(ctx, "in:devices", true, true)

fmt.Printf("Count: %d\n", results.Count)    // Results in this page
fmt.Printf("Total: %d\n", results.Total)    // Total matching results

if results.Next != nil {
    fmt.Printf("Next page starts at: %d\n", *results.Next)
}
```

## Common AQL Patterns

| Query | Description |
|-------|-------------|
| `in:devices` | All devices |
| `in:alerts` | All alerts |
| `in:vulnerabilities` | All vulnerabilities |
| `in:activities` | All activities |
| `in:devices type:laptop` | Laptops only |
| `in:devices manufacturer:Apple` | Apple devices |
| `in:alerts status:unresolved` | Open alerts |
| `in:devices riskLevel:>5` | Medium+ risk devices |
