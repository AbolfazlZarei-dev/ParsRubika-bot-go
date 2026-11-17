# `account` — Methods music

_Source: [`methods_music.go`](../../ParsRubika/account/methods_music.go)_

## Symbols

### `GetSavedMusicPlaylist`

```go
func (a *Account) GetSavedMusicPlaylist(ctx context.Context, userGUID string) (*events.Event, error)
```

GetSavedMusicPlaylist returns the user's saved music playlist.

### `SetSavedMusicPlaylist`

```go
func (a *Account) SetSavedMusicPlaylist(ctx context.Context, playlistTracks []map[string]any) (*events.Event, error)
```

SetSavedMusicPlaylist replaces the entire saved music playlist.

### `AddSavedMusicTrack`

```go
func (a *Account) AddSavedMusicTrack(ctx context.Context, track map[string]any) (*events.Event, error)
```

AddSavedMusicTrack adds a single track to the saved music playlist.

### `DeleteSavedMusicPlaylist`

```go
func (a *Account) DeleteSavedMusicPlaylist(ctx context.Context, position *int) (*events.Event, error)
```

DeleteSavedMusicPlaylist deletes the entire playlist (or a single position).
