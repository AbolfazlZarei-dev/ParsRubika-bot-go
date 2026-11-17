# `rubino` — Methods stories extra

_Source: [`methods_stories_extra.go`](../../ParsRubika/rubino/methods_stories_extra.go)_

## Symbols

### `AddStory`

```go
func (s *Rubino) AddStory(ctx context.Context, profileID, filePath, caption, fileType string) (*Event, error)
```

AddStory uploads and publishes a story. `fileType` is "Picture" or
"Video". The file is uploaded via RequestUploadFile + UploadFile.
