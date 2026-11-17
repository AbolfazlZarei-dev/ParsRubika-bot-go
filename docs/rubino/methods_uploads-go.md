# `rubino` — Methods uploads

_Source: [`methods_uploads.go`](../../ParsRubika/rubino/methods_uploads.go)_

## Symbols

### `RequestUploadFile`

```go
func (s *Rubino) RequestUploadFile(ctx context.Context, profileID, fileName string, fileSize int64, fileType string) (*Event, error)
```

RequestUploadFile asks the Rubino API to allocate an upload slot.
`fileType` is "Picture", "Video", "Audio" etc.

### `UploadFile`

```go
func (s *Rubino) UploadFile(ctx context.Context, uploadURL, fileName, filePath string) (string, error)
```

UploadFile uploads a local file to the supplied `uploadURL` using
multipart/form-data. Returns the file_id assigned by the server.

The Rubino upload protocol is similar to the messenger one: the
client requests an upload URL via RequestUploadFile, then POSTs the
file as `multipart/form-data` to the returned URL.

### `UploadAndPost`

```go
func (s *Rubino) UploadAndPost(ctx context.Context, profileID, postType, filePath, caption string, width, height int) (*Event, error)
```

UploadAndPost is a high-level helper that uploads a local file then
publishes a post. `postType` is "Picture" or "Video".

### `FormDataContentType`

```go
func (w *studioMultipartWriter) FormDataContentType() string
```
