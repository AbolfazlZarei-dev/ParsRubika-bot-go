# `enums` — Report

_Source: [`report.go`](../../ParsRubika/enums/report.go)_

## Symbols

### `ReportReason`

```go
type ReportReason string
```

ReportReason is the reason code sent when reporting an object.

### `ReportReasonPorn`

```go
const ReportReasonPorn ReportReason = "Porn"
```

### `ReportReasonViolence`

```go
const ReportReasonViolence ReportReason = "Violence"
```

### `ReportReasonSpam`

```go
const ReportReasonSpam ReportReason = "Spam"
```

### `ReportReasonAbuse`

```go
const ReportReasonAbuse ReportReason = "Abuse"
```

### `ReportReasonOther`

```go
const ReportReasonOther ReportReason = "Other"
```

### `String`

```go
func (r ReportReason) String() string { return string(r) }
```

String implements fmt.Stringer.
