# Reports

Access scheduled reports in Armis.

## List All Reports

```go
reports, err := client.GetReports(ctx)
if err != nil {
    log.Fatal(err)
}

for _, r := range reports {
    scheduled := "not scheduled"
    if r.IsScheduled {
        scheduled = "scheduled"
    }
    fmt.Printf("%s (%s) - %s\n", r.ReportName, r.ReportType, scheduled)
}
```

## Get a Single Report

```go
report, err := client.GetReportByID(ctx, "123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", report.ReportName)
fmt.Printf("Type: %s\n", report.ReportType)
fmt.Printf("Query: %s\n", report.Asq)
fmt.Printf("Created: %s\n", report.CreationTime)

if report.IsScheduled {
    fmt.Printf("Schedule:\n")
    fmt.Printf("  Every %v %s\n", report.Schedule.RepeatAmount, report.Schedule.RepeatUnit)
    fmt.Printf("  At: %s\n", report.Schedule.TimeOfDay)
    fmt.Printf("  Format: %s\n", report.Schedule.ReportFileFormat)
    fmt.Printf("  Recipients: %v\n", report.Schedule.Email)
}
```

## Report Fields

| Field | Description |
|-------|-------------|
| `ID` | Unique identifier |
| `ReportName` | Name of the report |
| `ReportType` | Type of report |
| `Asq` | AQL query for the report |
| `IsScheduled` | Whether the report runs on a schedule |
| `CreationTime` | When the report was created |

## Schedule Fields

| Field | Description |
|-------|-------------|
| `RepeatAmount` | How often to repeat |
| `RepeatUnit` | Unit (days, weeks, etc.) |
| `TimeOfDay` | Time to run |
| `Timezone` | Timezone for scheduling |
| `Weekdays` | Days of week to run |
| `Email` | Email recipients |
| `ReportFileFormat` | Output format |
