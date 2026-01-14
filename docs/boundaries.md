# Boundaries

View network boundaries configured in Armis.

## List All Boundaries

```go
boundaries, err := client.GetBoundaries(ctx)
if err != nil {
    log.Fatal(err)
}

for _, b := range boundaries {
    fmt.Printf("%s (ID: %d)\n", b.Name, b.ID)
    fmt.Printf("  Affected Sites: %s\n", b.AffectedSites)
}
```

## Get a Single Boundary

```go
boundary, err := client.GetBoundaryByID(ctx, "123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", boundary.Name)
fmt.Printf("Affected Sites: %s\n", boundary.AffectedSites)
fmt.Printf("Rules (AND): %v\n", boundary.RuleAQL.And)
fmt.Printf("Rules (OR): %v\n", boundary.RuleAQL.Or)
```

## Boundary Fields

| Field | Description |
|-------|-------------|
| `ID` | Unique identifier |
| `Name` | Boundary name |
| `AffectedSites` | Sites affected by this boundary |
| `RuleAQL` | AQL rules defining the boundary |
