# `tv` — Methods extra

_Source: [`methods_extra.go`](../../ParsRubika/tv/methods_extra.go)_

## Symbols

### `GetMediaByLink`

```go
func (t *TV) GetMediaByLink(ctx context.Context, link, trackID string) (*events.Event, error)
```

GetMediaByLink retrieves media details by its share link (e.g.
"https://rubika.ir/vod/294119"). The numeric ID is extracted from the
last path segment; non-numeric links are rejected.

### `GetNextEpisode`

```go
func (t *TV) GetNextEpisode(ctx context.Context, episodeID string) (*events.Event, error)
```

GetNextEpisode retrieves the next episode after the given episode.

### `GetNextMediaInfo`

```go
func (t *TV) GetNextMediaInfo(ctx context.Context, mediaID, trackID string) (*events.Event, error)
```

GetNextMediaInfo retrieves the next media info after the given media.

### `PlayType`

```go
type PlayType string
```

PlayType selects whether GetPlayInfo returns the main stream or the
trailer.

### `PlayTypeMain`

```go
const PlayTypeMain PlayType = "Main"
```

### `PlayTypeTrailer`

```go
const PlayTypeTrailer PlayType = "Trailer"
```

### `GetPlayInfo`

```go
func (t *TV) GetPlayInfo(ctx context.Context, mediaID string, playType PlayType, version int, trackID string) (*events.Event, error)
```

GetPlayInfo retrieves playback information for a media. `playType`
accepts "Main" or "Trailer" (case-insensitive, defaults to Trailer);
`version` defaults to 1.

### `GetPlaySuggestedEpisode`

```go
func (t *TV) GetPlaySuggestedEpisode(ctx context.Context, seriesID, trackID string) (*events.Event, error)
```

GetPlaySuggestedEpisode retrieves the suggested next episode for a
series.

### `SetCurrentPlayTime`

```go
func (t *TV) SetCurrentPlayTime(ctx context.Context, mediaID string, playTime int) (*events.Event, error)
```

SetCurrentPlayTime records the user's current play position (in
seconds) for a media. Negative values are rejected.

### `SetSessionPlayDuration`

```go
func (t *TV) SetSessionPlayDuration(ctx context.Context, sessionID string, playDuration int) (*events.Event, error)
```

SetSessionPlayDuration records the total play duration (in seconds)
for a session.

### `ViewType`

```go
type ViewType string
```

ViewType is the kind of commercial view event registered by
ViewCommercial.

### `ViewTypeFull`

```go
const ViewTypeFull ViewType = "Full"
```

### `ViewTypeHalf`

```go
const ViewTypeHalf ViewType = "Half"
```

### `ViewTypeSkip`

```go
const ViewTypeSkip ViewType = "Skip"
```

### `ViewCommercial`

```go
func (t *TV) ViewCommercial(ctx context.Context, viewType ViewType) (*events.Event, error)
```

ViewCommercial registers a commercial view event. `viewType` accepts
Full/Half/Skip (case-insensitive, defaults to Full).

### `DislikeMedia`

```go
func (t *TV) DislikeMedia(ctx context.Context, mediaID, trackID string) (*events.Event, error)
```

DislikeMedia adds a dislike to a media.

### `UndoLikeOrDislikeMedia`

```go
func (t *TV) UndoLikeOrDislikeMedia(ctx context.Context, mediaID, trackID string) (*events.Event, error)
```

UndoLikeOrDislikeMedia undoes a previous like or dislike.

### `FollowMedia`

```go
func (t *TV) FollowMedia(ctx context.Context, mediaID, trackID string) (*events.Event, error)
```

FollowMedia follows a media (e.g. a series).

### `UnfollowMedia`

```go
func (t *TV) UnfollowMedia(ctx context.Context, mediaID, trackID string) (*events.Event, error)
```

UnfollowMedia unfollows a media.
