# `account` — Methods profile

_Source: [`methods_profile.go`](../../ParsRubika/account/methods_profile.go)_

## Symbols

### `GetNotificationSetting`

```go
func (a *Account) GetNotificationSetting(ctx context.Context) (*events.Event, error)
```

GetNotificationSetting returns the user's notification settings.

### `SetNotification`

```go
func (a *Account) SetNotification(ctx context.Context, settings map[string]any, updateParameters []string) (*events.Event, error)
```

SetNotification sets notification preferences.

### `GetWallpapers`

```go
func (a *Account) GetWallpapers(ctx context.Context) (*events.Event, error)
```

GetWallpapers returns the available wallpapers.

### `ResetWallpapers`

```go
func (a *Account) ResetWallpapers(ctx context.Context) (*events.Event, error)
```

ResetWallpapers resets the user's wallpaper to default.

### `UpdateMyName`

```go
func (a *Account) UpdateMyName(ctx context.Context, firstName, lastName string) (*events.Event, error)
```

UpdateMyName updates the current user's display name.

### `UpdateMyBio`

```go
func (a *Account) UpdateMyBio(ctx context.Context, bio string) (*events.Event, error)
```

UpdateMyBio updates the current user's bio.

### `UpdateMyBirthday`

```go
func (a *Account) UpdateMyBirthday(ctx context.Context, birthday string) (*events.Event, error)
```

UpdateMyBirthday updates the current user's birthday (format: YYYY-MM-DD).

### `EditShowLastOnline`

```go
func (a *Account) EditShowLastOnline(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error)
```

EditShowLastOnline sets the privacy for "last online" visibility.

### `EditShowPhoneNumber`

```go
func (a *Account) EditShowPhoneNumber(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error)
```

EditShowPhoneNumber sets the privacy for phone number visibility.

### `EditShowProfilePhoto`

```go
func (a *Account) EditShowProfilePhoto(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error)
```

EditShowProfilePhoto sets the privacy for profile photo visibility.

### `EditShowBirthday`

```go
func (a *Account) EditShowBirthday(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error)
```

EditShowBirthday sets the privacy for birthday visibility.

### `EditCanCalledBy`

```go
func (a *Account) EditCanCalledBy(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error)
```

EditCanCalledBy sets the privacy for who can call the user.

### `EditCanJoinChatBy`

```go
func (a *Account) EditCanJoinChatBy(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error)
```

EditCanJoinChatBy sets the privacy for who can add the user to chats.

### `UnblockAllUsers`

```go
func (a *Account) UnblockAllUsers(ctx context.Context) error
```

UnblockAllUsers unblocks every blocked user. The blocked-users API
returns members under "in_chat_members" with each user's GUID in the
"user_guid" field (not "object_guid").

### `ProfilePhotoResult`

```go
type ProfilePhotoResult struct
```

ProfilePhotoResult holds the outcome of a profile-photo upload.

### `UploadProfilePhoto`

```go
func (a *Account) UploadProfilePhoto(ctx context.Context, path string) (*ProfilePhotoResult, error)
```

UploadProfilePhoto uploads an image file and sets it as the current
account's profile picture (avatar).

### `DeleteProfilePhoto`

```go
func (a *Account) DeleteProfilePhoto(ctx context.Context, avatarID string) error
```

DeleteProfilePhoto deletes a profile picture (avatar) by its id.
If avatarID is empty, the most recent avatar is fetched and deleted.

### `HangProfileResult`

```go
type HangProfileResult struct
```

HangProfileResult holds the outcomes of the two concurrent
updateProfile calls fired by HangProfile.

### `UpdateProfile`

```go
func (a *Account) UpdateProfile(ctx context.Context, firstName, lastName, bio, birthDate string) (*events.Event, error)
```

UpdateProfile updates the current user's profile. Pass empty strings
for fields you don't want to change. The `updated_parameters` list
is built automatically from the non-empty fields.

### `HangProfile`

```go
func (a *Account) HangProfile(ctx context.Context, firstName, lastName, bio, birthDate string) (*HangProfileResult, error)
```

HangProfile (a.k.a "profile hang" / "dual profile") fires two
updateProfile requests back-to-back against the current account so
the two writes race on the server. This is the well-known Rubika
trick used to register a "dual"/"hang" profile: when two
updateProfile calls land near-simultaneously, the server may persist
an inconsistent state and the profile then renders as two stacked
identities in clients.

The two shots use the SAME payload (firstName/lastName/bio/birthDate)
so both writes target the same fields — the "hang" effect comes from
the race itself, not from conflicting values. If you want the two
shots to carry different fields, build them with UpdateProfile
directly.

Both requests are issued concurrently (two goroutines) and the call
blocks until both return. The context passed to each shot is the
caller's context; cancelling it aborts whichever shot is still
in-flight. The method returns the result of the two shots; it does
not raise on a single shot failing — inspect FirstErr / SecondErr.

### `UpdateUsername`

```go
func (a *Account) UpdateUsername(ctx context.Context, username string) (*events.Event, error)
```

UpdateUsername changes the user's username (strips leading `@`).

### `PrivacyValue`

```go
type PrivacyValue string
```

PrivacyValue is one of "Everybody", "MyContacts", "Nobody".

### `PrivacyEverybody`

```go
const PrivacyEverybody PrivacyValue = "Everybody"
```

### `PrivacyMyContacts`

```go
const PrivacyMyContacts PrivacyValue = "MyContacts"
```

### `PrivacyNobody`

```go
const PrivacyNobody PrivacyValue = "Nobody"
```

### `PrivacyExceptions`

```go
type PrivacyExceptions struct
```

PrivacyExceptions is the include/exclude list for a privacy setting.

### `SetPrivacySetting`

```go
func (a *Account) SetPrivacySetting(ctx context.Context, settingType string, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error)
```

SetPrivacySetting sets a single privacy setting. `settingType` is one
of `ShowLastOnline`, `ShowPhone`, `ShowPhoto`, `CanJoinChat`,
`ShowBirthDate`, `CanCalled`. The `exceptions` parameter is optional.

### `SetSetting`

```go
func (a *Account) SetSetting(ctx context.Context, settings map[string]any, updateParameters []string) (*events.Event, error)
```

SetSetting is a low-level escape hatch that lets the caller build
the `settings` dict manually.

### `GetPrivacySetting`

```go
func (a *Account) GetPrivacySetting(ctx context.Context) (*events.Event, error)
```

GetPrivacySetting returns the user's privacy settings.

### `GetBlockedUsers`

```go
func (a *Account) GetBlockedUsers(ctx context.Context, startID ...string) (*events.Event, error)
```

GetBlockedUsers returns the list of blocked users; pass startID to
continue a previous page.
