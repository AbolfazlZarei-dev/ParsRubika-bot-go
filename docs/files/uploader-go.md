# `files` — Uploader

_Source: [`uploader.go`](../../ParsRubika/files/uploader.go)_

## Symbols

### `MaxUploadSize`

```go
const MaxUploadSize = 2 << 30 // 2 GiB
```

MaxUploadSize is the client-side sanity cap on a single uploaded file.
The Rubika server enforces its own limit; this guard only rejects
absurd sizes before any byte is read.

### `Source`

```go
type Source struct
```

Source is the union of an on-disk path or an in-memory blob for
upload. Exactly one of Path / Bytes must be set.

### `Uploader`

```go
type Uploader struct
```

Uploader uploads a single file using the chunked protocol. It is
safe for concurrent use after construction (each upload runs on its
own context).

### `NewUploader`

```go
func NewUploader(client *http.Client, chunkSize int64, maxRetries int, backoff time.Duration, headers map[string]string) *Uploader
```

NewUploader constructs an Uploader.

### `UploadResult`

```go
type UploadResult struct
```

UploadResult is the metadata returned by the server after a
successful upload.

### `RequestSendFile`

```go
type RequestSendFile func(ctx context.Context, fileName string, size int64, mime string) (fileID string, dcID int, uploadURL string, accessHashSend string, err error)
```

RequestSendFile is the signature of the function used to obtain the
initial (file_id, dc_id, upload_url, access_hash_send) tuple.

### `Upload`

```go
func (u *Uploader) Upload(ctx context.Context, src Source, mime, fileName string, requestSendFile RequestSendFile, cb ProgressCallback) (*UploadResult, error)
```

Upload uploads `src` to the server. It calls `requestSendFile` to
obtain the upload metadata, then chunks the file and uploads each
chunk to `uploadURL`. When the server returns `ERROR_TRY_AGAIN`, the
upload is restarted.

`cb` (if non-nil) is called after every chunk with the total size
and the uploaded-so-far byte count.

### `MaxParallelDownloads`

```go
const MaxParallelDownloads = 8
```

MaxParallelDownloads bounds the number of concurrent chunk
downloads when `gather` is enabled.

### `Result`

```go
type Result struct
```

Result is the union type returned by Downloader.Download. Exactly
one of Path / Bytes will be populated depending on `saveAs`.

### `Downloader`

```go
type Downloader struct
```

Downloader fetches a file by (dc_id, file_id, access_hash) from the
Rubika file servers using chunked HTTP requests.

### `NewDownloader`

```go
func NewDownloader(client *http.Client, chunkSize int64, maxRetries int, backoff time.Duration, headers map[string]string) *Downloader
```

NewDownloader constructs a Downloader.

### `Download`

```go
func (d *Downloader) Download(ctx context.Context, dcID int, fileID, accessHash string, size int64, saveAs string, gather bool, cb ProgressCallback) (*Result, error)
```

Download fetches the file. When `saveAs` is non-empty the bytes are
written to disk and the Result.Path is populated. When `gather` is
true and `saveAs` is empty, chunks are fetched in parallel.
