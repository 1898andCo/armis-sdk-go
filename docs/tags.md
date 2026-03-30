# Tags

Retrieve all tags configured in Armis. Tags are returned as a flat list of strings with automatic pagination.

## List All Tags

```go
tags, err := client.GetTags(ctx)
if err != nil {
    log.Fatal(err)
}

for _, tag := range tags {
    fmt.Println(tag)
}
```

## Response

`GetTags` returns `([]string, error)`. The client handles pagination automatically, fetching all pages and returning the complete list.

## Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `Tags` | `[]string` | List of tag names |
| `Count` | `int` | Number of tags in the current page |
| `Next` | `*int` | Offset for the next page (`nil` when no more pages) |
| `Prev` | `*int` | Offset for the previous page |
