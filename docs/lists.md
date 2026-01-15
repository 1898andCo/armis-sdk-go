# Lists

View custom lists in Armis.

## List All Lists

```go
lists, err := client.GetLists(ctx)
if err != nil {
    log.Fatal(err)
}

for _, list := range lists {
    fmt.Printf("%s (%s)\n", list.ListName, list.ListType)
}
```

## Working with List Data

```go
lists, err := client.GetLists(ctx)
if err != nil {
    log.Fatal(err)
}

for _, list := range lists {
    fmt.Printf("List: %s\n", list.ListName)
    fmt.Printf("  ID: %d\n", list.ListID)
    fmt.Printf("  Type: %s\n", list.ListType)
    fmt.Printf("  Description: %s\n", list.Description)
    fmt.Printf("  Created by: %s at %s\n", list.CreatedBy, list.CreationTime)
    fmt.Printf("  Last updated by: %s at %s\n", list.LastUpdatedBy, list.LastUpdateTime)
}
```

## List Fields

| Field | Description |
|-------|-------------|
| `ListID` | Unique identifier |
| `ListName` | Name of the list |
| `ListType` | Type of list |
| `Description` | List description |
| `CreatedBy` | Who created the list |
| `CreationTime` | When it was created |
| `LastUpdatedBy` | Who last modified it |
| `LastUpdateTime` | When it was last modified |
