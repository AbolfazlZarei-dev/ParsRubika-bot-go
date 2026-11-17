# `account` — Avatar management (list, upload, set main, delete)

_Source: [`methods_avatar.go`](../../ParsRubika/account/methods_avatar.go)_

## Symbols

### `GetAvatars`

```go
func (a *Account) GetAvatars(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetAvatars returns the avatar list of an object.

### `UploadAvatar`

```go
func (a *Account) UploadAvatar(ctx context.Context, objectGUID, mainFileID, thumbnailFileID string) (*events.Event, error)
```

UploadAvatar uploads a new avatar for the supplied object. mainFileID
is the full-size photo and thumbnailFileID the small one; passing the
same id for both reuses one upload as its own thumbnail. (Wire format
verified against maxrubika's working upload_avatar: the server wants
object_guid + main_file_id + thumbnail_file_id; file_id/dc_id/
access_hash_rec are rejected with INVALID_INPUT.)

### `SetMainAvatar`

```go
func (a *Account) SetMainAvatar(ctx context.Context, objectGUID, avatarID string) (*events.Event, error)
```

SetMainAvatar promotes one of the object's existing avatars (by
avatar_id) to be the current/main one. Wire method: updateAvatar.

### `DeleteAvatar`

```go
func (a *Account) DeleteAvatar(ctx context.Context, objectGUID, avatarID string) (*events.Event, error)
```

DeleteAvatar deletes an avatar.

### `DeleteAllAvatars`

```go
func (a *Account) DeleteAllAvatars(ctx context.Context, objectGUID string) error
```

DeleteAllAvatars deletes every avatar of an object by looping over
getAvatars and calling deleteAvatar for each.
