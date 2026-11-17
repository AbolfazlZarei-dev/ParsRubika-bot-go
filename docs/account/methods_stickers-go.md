# `account` — Methods stickers

_Source: [`methods_stickers.go`](../../ParsRubika/account/methods_stickers.go)_

## Symbols

### `GetStickerSetting`

```go
func (a *Account) GetStickerSetting(ctx context.Context) (*events.Event, error)
```

GetStickerSetting returns the user's sticker settings.

### `GetMyArchivedStickerSets`

```go
func (a *Account) GetMyArchivedStickerSets(ctx context.Context, startID string) (*events.Event, error)
```

GetMyArchivedStickerSets returns the user's archived sticker sets.

### `GetMyStickerSets`

```go
func (a *Account) GetMyStickerSets(ctx context.Context) (*events.Event, error)
```

GetMyStickerSets returns the user's installed sticker sets.

### `GetStickerSetByID`

```go
func (a *Account) GetStickerSetByID(ctx context.Context, setID string) (*events.Event, error)
```

GetStickerSetByID fetches a sticker set by id.

### `GetStickersByEmoji`

```go
func (a *Account) GetStickersByEmoji(ctx context.Context, emoji, suggestBy string) (*events.Event, error)
```

GetStickersByEmoji returns stickers matching an emoji. `suggestBy`
defaults to "All" when empty.

### `GetStickersBySetIDs`

```go
func (a *Account) GetStickersBySetIDs(ctx context.Context, setIDs []string) (*events.Event, error)
```

GetStickersBySetIDs returns stickers from multiple sets.

### `SearchStickers`

```go
func (a *Account) SearchStickers(ctx context.Context, query, startID string) (*events.Event, error)
```

SearchStickers searches for stickers by keyword (paged).

### `GetTrendStickerSets`

```go
func (a *Account) GetTrendStickerSets(ctx context.Context, startID string) (*events.Event, error)
```

GetTrendStickerSets returns the trending sticker sets (paged).

### `ActionOnStickerSet`

```go
func (a *Account) ActionOnStickerSet(ctx context.Context, setID, action string) (*events.Event, error)
```

ActionOnStickerSet installs/uninstalls/archives a sticker set.
`action` is "Add", "Remove", or "Archive".

### `AddStickerSet`

```go
func (a *Account) AddStickerSet(ctx context.Context, setID string) (*events.Event, error)
```

AddStickerSet installs a sticker set (action = "Add").

### `DeleteStickerSet`

```go
func (a *Account) DeleteStickerSet(ctx context.Context, setID string) (*events.Event, error)
```

DeleteStickerSet uninstalls a sticker set (action = "Remove").

### `ArchiveStickerSet`

```go
func (a *Account) ArchiveStickerSet(ctx context.Context, setID string) (*events.Event, error)
```

ArchiveStickerSet archives a sticker set (action = "Archive").

### `ReorderStickerSets`

```go
func (a *Account) ReorderStickerSets(ctx context.Context, orderedSetIDs []string) (*events.Event, error)
```

ReorderStickerSets reorders the user's installed sticker sets.
