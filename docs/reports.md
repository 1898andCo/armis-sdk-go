# Reports

Create, access, and manage scheduled reports in Armis.

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

## Create a Report

```go
report := armis.CreateReportRequest{
    ReportName: "Critical Vulnerabilities Report",
    ASQ:        "in:vulnerabilities severity:critical",
    ExportConfiguration: armis.ExportConfiguration{
        Columns: armis.ExportColumns{
            Vulnerabilities: []string{"name", "severity", "cveId", "affectedDevices"},
        },
    },
    Schedule: armis.CreateSchedule{
        Email:            []string{"security@example.com"},
        RepeatAmount:     "1",
        RepeatUnit:       "week",
        TimeOfDay:        "09:00",
        Timezone:         "America/New_York",
        Weekdays:         []string{"monday"},
        ReportFileFormat: "csv",
    },
}

created, err := client.CreateReport(ctx, report)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created report: %s (ID: %d)\n", created.ReportName, created.ID)
```

## Update a Report

```go
update := armis.UpdateReportRequest{
    ReportName:   "Updated Report Name",
    ASQ:          "in:devices timeFrame:\"7 Days\"",
    EmailSubject: "Weekly Device Report",
    Schedule: armis.CreateSchedule{
        Email:            []string{"team@example.com"},
        RepeatAmount:     "1",
        RepeatUnit:       "week",
        TimeOfDay:        "08:00",
        Timezone:         "America/New_York",
        Weekdays:         []string{"monday"},
        ReportFileFormat: "pdf",
    },
}

updated, err := client.UpdateReport(ctx, "123", update)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated report: %s\n", updated.ReportName)
```

### Partial Updates

Since UpdateReport uses PATCH, you can update only specific fields:

```go
// Update only the report name
update := armis.UpdateReportRequest{
    ReportName: "New Report Name",
}

updated, err := client.UpdateReport(ctx, "123", update)
if err != nil {
    log.Fatal(err)
}
```

## Delete a Report

```go
success, err := client.DeleteReport(ctx, "123")
if err != nil {
    log.Fatal(err)
}

if success {
    fmt.Println("Report deleted")
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

## Create Report Request Fields

| Field | Description |
|-------|-------------|
| `ReportName` | Name of the report (required) |
| `ASQ` | AQL query for the report (required) |
| `EmailSubject` | Custom email subject for scheduled reports |
| `ExportConfiguration` | Columns to include in exports |
| `Schedule` | Schedule configuration for recurring reports |

## Update Report Request Fields

| Field | Description |
|-------|-------------|
| `ReportName` | Name of the report (optional) |
| `ASQ` | AQL query for the report (optional) |
| `EmailSubject` | Custom email subject for scheduled reports (optional) |
| `ExportConfiguration` | Columns to include in exports (optional) |
| `Schedule` | Schedule configuration for recurring reports (optional) |

> **Note:** All fields are optional for update requests since PATCH supports partial updates.

## Export Configuration

| Field | Description |
|-------|-------------|
| `Columns.Devices` | Device columns to export |
| `Columns.Vulnerabilities` | Vulnerability columns to export |
| `Columns.Activities` | Activity columns to export |

## Validation Rules

- ReportName is required
- ASQ (query) is required
