# `files` — Thumbnail

_Source: [`thumbnail.go`](../../ParsRubika/files/thumbnail.go)_

## Symbols

### `ThumbnailResult`

```go
type ThumbnailResult struct
```

ThumbnailResult holds the base64-encoded thumbnail and the original
image dimensions. The server expects `thumb_inline` (base64 string),
and the REAL width/height of the source image (not arbitrary values
like 200x200). Without `thumb_inline` the server rejects Image-type
sendMessage requests with INVALID_INPUT.

### `GenerateThumbnail`

```go
func GenerateThumbnail(src Source) (*ThumbnailResult, error)
```

GenerateThumbnail reads `src` (path or bytes), decodes the image,
resizes it to a small thumbnail (max 100px on the longest side,
preserving aspect ratio), re-encodes it as JPEG (quality 85 — much
smaller than PNG for photos), and returns a ThumbnailResult.

Returns nil if the image cannot be decoded — callers should fall
back to type="File" in that case.
