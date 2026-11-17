# `rubino` — Methods profiles

_Source: [`methods_profiles.go`](../../ParsRubika/rubino/methods_profiles.go)_

## Symbols

### `GetProfileList`

```go
func (s *Rubino) GetProfileList(ctx context.Context, limit int, sort, equal string) (*Event, error)
```

GetProfileList returns the user's pages/profiles.

### `GetMyProfileInfo`

```go
func (s *Rubino) GetMyProfileInfo(ctx context.Context, profileID string) (*Event, error)
```

GetMyProfileInfo returns the current user's full profile.

### `CreatePage`

```go
func (s *Rubino) CreatePage(ctx context.Context, username, name, bio, phone, email, website string) (*Event, error)
```

CreatePage creates a new Rubino page.

### `UpdateProfile`

```go
func (s *Rubino) UpdateProfile(ctx context.Context, profileID, username, name, bio, phone, email, website, profileStatus string, isMessageAllowed *bool, isMute *bool) (*Event, error)
```

UpdateProfile updates the current user's profile metadata. Pass
empty strings for fields you don't want to change.

### `IsExistUsername`

```go
func (s *Rubino) IsExistUsername(ctx context.Context, username string) (*Event, error)
```

IsExistUsername checks whether a username is available.

### `GetProfileInfo`

```go
func (s *Rubino) GetProfileInfo(ctx context.Context, profileID, targetProfileID string) (*Event, error)
```

GetProfileInfo returns info about a target profile.

### `BlockProfile`

```go
func (s *Rubino) BlockProfile(ctx context.Context, profileID, blockedID string) (*Event, error)
```

BlockProfile blocks a profile.

### `UnblockProfile`

```go
func (s *Rubino) UnblockProfile(ctx context.Context, profileID, blockedID string) (*Event, error)
```

UnblockProfile unblocks a profile.

### `RemovePage`

```go
func (s *Rubino) RemovePage(ctx context.Context, profileID, recordID string) (*Event, error)
```

RemovePage removes a page. The removeRecord API expects a
`record_id` + `model` pair (model="Profile").

### `SetProfilePublic`

```go
func (s *Rubino) SetProfilePublic(ctx context.Context, profileID string) (*Event, error)
```

SetProfilePublic sets a profile public.

### `SetProfilePrivate`

```go
func (s *Rubino) SetProfilePrivate(ctx context.Context, profileID string) (*Event, error)
```

SetProfilePrivate sets a profile private.

### `EnablePrivateMessages`

```go
func (s *Rubino) EnablePrivateMessages(ctx context.Context, profileID string) (*Event, error)
```

EnablePrivateMessages allows others to message the profile.

### `DisablePrivateMessages`

```go
func (s *Rubino) DisablePrivateMessages(ctx context.Context, profileID string) (*Event, error)
```

DisablePrivateMessages disallows others from messaging the profile.
