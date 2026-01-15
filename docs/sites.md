# Sites

View site information in Armis.

## List All Sites

```go
sites, err := client.GetSites(ctx)
if err != nil {
    log.Fatal(err)
}

for _, site := range sites {
    fmt.Printf("%s - %s\n", site.Name, site.Location)
}
```

## Working with Site Data

```go
sites, err := client.GetSites(ctx)
if err != nil {
    log.Fatal(err)
}

for _, site := range sites {
    fmt.Printf("Site: %s\n", site.Name)
    fmt.Printf("  ID: %s\n", site.ID)
    fmt.Printf("  Location: %s\n", site.Location)
    fmt.Printf("  Tier: %s\n", site.Tier)

    if site.Lat != 0 && site.Lng != 0 {
        fmt.Printf("  Coordinates: %.4f, %.4f\n", site.Lat, site.Lng)
    }
}
```

## Site Fields

| Field | Description |
|-------|-------------|
| `ID` | Unique identifier |
| `Name` | Site name |
| `Location` | Physical location description |
| `Tier` | Site tier/priority level |
| `Lat` | Latitude coordinate |
| `Lng` | Longitude coordinate |
| `User` | Associated user |
