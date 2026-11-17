# `files` — Package files implements the chunked uploader and downloader used

_Source: [`files.go`](../../ParsRubika/files/files.go)_

## Symbols

### `ProgressCallback`

```go
type ProgressCallback func(total, current int64)
```

ProgressCallback is invoked after each chunk upload/download with
the current (total, current) byte counters.
