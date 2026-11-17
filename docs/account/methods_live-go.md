# `account` — Methods live

_Source: [`methods_live.go`](../../ParsRubika/account/methods_live.go)_

## Symbols

### `SendLive`

```go
func (a *Account) SendLive(ctx context.Context, objectGUID, title string, thumbInline any) (*events.Event, error)
```

SendLive starts a live stream in a chat.

### `GetLiveStatus`

```go
func (a *Account) GetLiveStatus(ctx context.Context, liveID, accessToken string) (*events.Event, error)
```

GetLiveStatus returns the status of a live stream.

### `GetLivePlayURL`

```go
func (a *Account) GetLivePlayURL(ctx context.Context, accessToken, liveID string) (*events.Event, error)
```

GetLivePlayURL returns the playback URL of a live stream.

### `GetLiveComments`

```go
func (a *Account) GetLiveComments(ctx context.Context, liveID, accessToken string) (*events.Event, error)
```

GetLiveComments returns the comments of a live stream.

### `AddLiveComment`

```go
func (a *Account) AddLiveComment(ctx context.Context, accessToken, liveID, text string) (*events.Event, error)
```

AddLiveComment adds a comment to a live stream.

### `SetCurrentLiveLocation`

```go
func (a *Account) SetCurrentLiveLocation(ctx context.Context, objectGUID string, lat, lon float64, tileSideCount int, tileURLs []string, xLoc, yLoc float64) (*events.Event, error)
```

SetCurrentLiveLocation shares a live location with map tiles.

### `GetMapView`

```go
func (a *Account) GetMapView(ctx context.Context, lat, lon float64) (*events.Event, error)
```

GetMapView returns a map view for the given coordinates.
