# Collectors

Manage Armis data collectors.

## List All Collectors

```go
collectors, err := client.GetCollectors(ctx)
if err != nil {
    log.Fatal(err)
}

for _, c := range collectors {
    fmt.Printf("%s - %s (%s)\n", c.Name, c.Status, c.IPAddress)
}
```

## Get a Single Collector

```go
collector, err := client.GetCollectorByID(ctx, "123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", collector.Name)
fmt.Printf("Status: %s\n", collector.Status)
fmt.Printf("IP: %s\n", collector.IPAddress)
fmt.Printf("Last Seen: %s\n", collector.LastSeen)
```

## Create a Collector

```go
newCollector := armis.CreateCollectorSettings{
    Name:           "Office Collector",
    DeploymentType: "physical",
}

created, err := client.CreateCollector(ctx, newCollector)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Collector ID: %d\n", created.CollectorID)
fmt.Printf("License Key: %s\n", created.LicenseKey)
fmt.Printf("Username: %s\n", created.User)
fmt.Printf("Password: %s\n", created.Password)
```

## Update a Collector

```go
updates := armis.UpdateCollectorSettings{
    Name:           "Main Office Collector",
    DeploymentType: "physical",
}

updated, err := client.UpdateCollector(ctx, "123", updates)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated: %s\n", updated.Name)
```

## Delete a Collector

```go
success, err := client.DeleteCollector(ctx, "123")
if err != nil {
    log.Fatal(err)
}

if success {
    fmt.Println("Collector deleted")
}
```

## Collector Fields

| Field | Description |
|-------|-------------|
| `Name` | Collector name |
| `Status` | Current status |
| `IPAddress` | IP address |
| `MacAddress` | MAC address |
| `LastSeen` | Last communication time |
| `BootTime` | Last boot time |
| `Type` | Collector type |
| `CollectorNumber` | Assigned number |
| `ClusterID` | Cluster assignment |
