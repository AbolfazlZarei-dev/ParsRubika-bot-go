# `account` — Methods files

_Source: [`methods_files.go`](../../ParsRubika/account/methods_files.go)_

## Symbols

### `UploadFromPath`

```go
func (a *Account) UploadFromPath(ctx context.Context, path string) (*UploadResult, error)
```

UploadFromPath uploads a local file by path and returns the upload
result event. The MIME type is auto-detected from the file extension.

### `UploadFromBytes`

```go
func (a *Account) UploadFromBytes(ctx context.Context, data []byte, fileName, mime string) (*UploadResult, error)
```

UploadFromBytes uploads raw bytes as a file.

### `UploadResult`

```go
type UploadResult struct
```

UploadResult wraps an upload response with metadata.

### `DownloadToPath`

```go
func (a *Account) DownloadToPath(ctx context.Context, fileInline *events.Event, saveAs string) error
```

DownloadToPath downloads a file inline to the given path.

### `DownloadToBytes`

```go
func (a *Account) DownloadToBytes(ctx context.Context, fileInline *events.Event) ([]byte, error)
```

DownloadToBytes downloads a file inline and returns the raw bytes.

### `UploadFromURL`

```go
func (a *Account) UploadFromURL(ctx context.Context, url, fileName string) (*UploadResult, error)
```

UploadFromURL downloads a file from a URL and uploads it to Rubika.
This is useful for forwarding media from other sources.

Security: the URL is validated and resolved through the netutil SSRF
policy (http/https only; private, loopback, link-local and cloud
metadata IPs are refused both at pre-flight and at the dial layer),
and the response body is capped at netutil.MaxFetchBytes so a
malicious remote cannot exhaust memory.

### `FileExists`

```go
func FileExists(path string) bool
```

FileExists reports whether a file exists at the given path.

### `FileSize`

```go
func FileSize(path string) int64
```

FileSize returns the size of a file in bytes, or 0 on error.

### `DownloadProfilePicture`

```go
func (a *Account) DownloadProfilePicture(ctx context.Context, objectGUID string, saveAs string) (any, error)
```

DownloadProfilePicture downloads the avatar of an object. If saveAs is
empty the raw bytes are returned; otherwise the file path is returned.

### `RequestSendFile`

```go
func (a *Account) RequestSendFile(ctx context.Context, fileName string, size int64, mime string) (*events.Event, error)
```

RequestSendFile asks the server to allocate an upload slot.

### `Upload`

```go
func (a *Account) Upload(ctx context.Context, src files.Source, mime, fileName string) (*events.Event, error)
```

Upload uploads a file and returns the metadata wrapped in an
*events.Event.

### `Download`

```go
func (a *Account) Download(ctx context.Context, fileInline *events.Event, saveAs string, chunkSize int64, cb files.ProgressCallback) (*files.Result, error)
```

Download downloads a file by its file_inline metadata. When `saveAs`
is empty the raw bytes are returned.
