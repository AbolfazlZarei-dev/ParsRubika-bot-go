# `rubino` — Methods missing

_Source: [`methods_missing.go`](../../ParsRubika/rubino/methods_missing.go)_

## Symbols

### `GetShareLink`

```go
func (s *Rubino) GetShareLink(ctx context.Context, profileID, postID, postProfileID string) (*events.Event, error)
```

GetShareLink returns the share link for a Rubino post.

### `UploadProfilePhoto`

```go
func (s *Rubino) UploadProfilePhoto(ctx context.Context, profileID, fileID string) (*events.Event, error)
```

UploadProfilePhoto uploads a new profile photo.
