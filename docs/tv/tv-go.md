# `tv` — Package tv implements the ParsRubika TV client. It talks to the

_Source: [`tv.go`](../../ParsRubika/tv/tv.go)_

## Symbols

### `TVServerURLs`

```go
var TVServerURLs = []string{
```

TVServerURLs is the list of TV API servers. A random one is picked
per request to spread the load.

### `MetaURL`

```go
const MetaURL = "https://tv.rubika.ir/meta.json"
```

MetaURL is the URL of the version-meta JSON.

### `DefaultVersion`

```go
const DefaultVersion = "3.2.3"
```

DefaultVersion is the fallback app version when MetaURL is
unreachable.

### `TVHeaders`

```go
var TVHeaders = map[string]string{
```

TVHeaders are the TV-specific HTTP headers (origin / referer set to
tv.rubika.ir).

### `TV`

```go
type TV struct
```

TV is the ParsRubika TV client.

### `Option`

```go
type Option func(*TV)
```

Option configures a TV client at construction time.

### `WithAuth`

```go
func WithAuth(auth string) Option
```

WithAuth sets the auth key.

### `WithProxy`

```go
func WithProxy(proxy string) Option
```

WithProxy sets a proxy URL.

### `WithLogger`

```go
func WithLogger(l *log.Logger) Option
```

WithLogger injects a custom logger.

### `WithVersion`

```go
func WithVersion(v string) Option
```

WithVersion overrides the app version (normally auto-detected
from MetaURL).

### `New`

```go
func New(opts ...Option) *TV
```

New constructs a TV client.

### `SetHTTPClient`

```go
func (t *TV) SetHTTPClient(c *http.Client)
```

SetHTTPClient lets callers swap in an HTTP/3 client.

### `SetAuth`

```go
func (t *TV) SetAuth(auth string)
```

SetAuth updates the auth key. Used by callers that want to reuse
the account's auth key (after the account signs in).

### `GetListItems`

```go
func (t *TV) GetListItems(ctx context.Context, listID, listType, startID string) (*events.Event, error)
```

GetListItems returns the items of a list (e.g. a category).
`listType` defaults to "Media"; `startID` defaults to "0".

### `GetListing`

```go
func (t *TV) GetListing(ctx context.Context, listingID string) (*events.Event, error)
```

GetListing returns the items of a listing page. `listingID`
defaults to "home".

### `GetCustomMenuItems`

```go
func (t *TV) GetCustomMenuItems(ctx context.Context) (*events.Event, error)
```

GetCustomMenuItems returns the user's custom menu.

### `GetPropertyTypes`

```go
func (t *TV) GetPropertyTypes(ctx context.Context) (*events.Event, error)
```

GetPropertyTypes returns the available property types (genres,
categories, etc.).

### `GetMedia`

```go
func (t *TV) GetMedia(ctx context.Context, mediaID, trackID string) (*events.Event, error)
```

GetMedia returns info about a single media.

### `Search`

```go
func (t *TV) Search(ctx context.Context, text, startID string) (*events.Event, error)
```

Search searches for media by text.

### `GetRelated`

```go
func (t *TV) GetRelated(ctx context.Context, mediaID, startID string) (*events.Event, error)
```

GetRelated returns the related medias.

### `GetCastMedias`

```go
func (t *TV) GetCastMedias(ctx context.Context, castID, startID string) (*events.Event, error)
```

GetCastMedias returns the medias for a cast (actor/director).

### `GetCastDetails`

```go
func (t *TV) GetCastDetails(ctx context.Context, castID string) (*events.Event, error)
```

GetCastDetails returns info about a cast member.

### `GetSeasonEpisodes`

```go
func (t *TV) GetSeasonEpisodes(ctx context.Context, seasonID, seriesID string) (*events.Event, error)
```

GetSeasonEpisodes returns the episodes of a season.

### `GetAccountInfo`

```go
func (t *TV) GetAccountInfo(ctx context.Context) (*events.Event, error)
```

GetAccountInfo returns the user's TV account info.

### `GetWishList`

```go
func (t *TV) GetWishList(ctx context.Context, startID string) (*events.Event, error)
```

GetWishList returns the user's wish-list.

### `LikeMedia`

```go
func (t *TV) LikeMedia(ctx context.Context, mediaID string) (*events.Event, error)
```

LikeMedia adds a like to a media.

### `UnlikeMedia`

```go
func (t *TV) UnlikeMedia(ctx context.Context, mediaID string) (*events.Event, error)
```

UnlikeMedia removes a like from a media. The action verb is
"Unlike" (matching the Rubika like/unlike vocabulary); "Dislike" is
a separate thumbs-down action, not the inverse of Like.

### `AddWishMedia`

```go
func (t *TV) AddWishMedia(ctx context.Context, mediaID string) (*events.Event, error)
```

AddWishMedia adds a media to the wish-list.

### `RemoveWishMedia`

```go
func (t *TV) RemoveWishMedia(ctx context.Context, mediaID string) (*events.Event, error)
```

RemoveWishMedia removes a media from the wish-list.

### `Close`

```go
func (t *TV) Close() error { return nil }
```

Close releases the underlying transport.
