# `account` — Hangi (dual) profile + tolerant image upload

_Source: [`methods_hangi.go`](../../ParsRubika/account/methods_hangi.go)_

## Symbols

### `AvatarFileRef`

```go
type AvatarFileRef struct
```

AvatarFileRef is an image that is already in Rubika's file storage
and is ready to be bound as an avatar.

### `HangiOptions`

```go
type HangiOptions struct
```

HangiOptions tunes HangiProfile. The zero value is what you want in
almost every case: ONE application that leaves exactly the two photos.

### `HangiOption`

```go
type HangiOption func(*HangiOptions)
```

HangiOption mutates HangiOptions.

### `WithHangiRounds`

```go
func WithHangiRounds(n int) HangiOption
```

WithHangiRounds overrides the number of dual-bind rounds (default 1).

### `WithHangiOffset`

```go
func WithHangiOffset(d time.Duration) HangiOption
```

WithHangiOffset overrides the stagger before the second bind.

### `WithHangiGap`

```go
func WithHangiGap(d time.Duration) HangiOption
```

WithHangiGap overrides the pause between rounds.

### `WithHangiClean`

```go
func WithHangiClean() HangiOption
```

WithHangiClean deletes all existing avatars before binding the pair.

### `WithHangiKeepOld`

```go
func WithHangiKeepOld() HangiOption
```

WithHangiKeepOld keeps avatars that this run did not create (skips
the clean-up pass).

### `WithHangiProgress`

```go
func WithHangiProgress(fn func(string)) HangiOption
```

WithHangiProgress installs a progress sink.

### `HangiProfileResult`

```go
type HangiProfileResult struct
```

HangiProfileResult reports what the server ended up with.

### `UploadImage`

```go
func (a *Account) UploadImage(ctx context.Context, path string) (*AvatarFileRef, error)
```

UploadImage uploads a local image file (jpg/png/gif/webp) to Rubika's
upload service and returns the reference needed to bind it as an
avatar.

It talks to the upload service through the same tolerant chunk parser
the hangprofile example used: success is an access_hash_rec anywhere
in the response — string or number, top level or nested — because the
service has been seen to answer without the status/status_det fields
the stock files.Uploader insists on.

### `HangiProfile`

```go
func (a *Account) HangiProfile(ctx context.Context, objectGUID, photo1Path, photo2Path string, opts ...HangiOption) (*HangiProfileResult, error)
```

HangiProfile sets the dual (hangi) profile — پروفایل هنگی — on
objectGUID with two local image files. Pass an empty objectGUID to
target the signed-in account itself; pass a group/channel GUID to
hang the pair on that chat instead.

What it does, in order:
1. optionally deletes every existing avatar (WithHangiClean);
2. uploads both images (tolerant upload, see UploadImage);
3. fires ONE round of two simultaneous uploadAvatar binds — photo1
main + photo2 thumb, and the swapped pair — so the object ends up
with two different pictures (offset staggers the second bind);
4. deletes every avatar this run did NOT create, so rerunning the
helper never piles up entries: the object ends with exactly the
two new photos. Disable with WithHangiKeepOld.

The default is a single application: Rounds=1. More rounds only add
avatar entries server-side; the hang comes from the two writes
colliding, not from repetition.
