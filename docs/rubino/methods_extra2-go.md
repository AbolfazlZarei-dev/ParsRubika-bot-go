# `rubino` — Methods extra2

_Source: [`methods_extra2.go`](../../ParsRubika/rubino/methods_extra2.go)_

## Symbols

### `GetProfileSetting`

```go
func (s *Rubino) GetProfileSetting(ctx context.Context, profileID string) (*events.Event, error)
```

GetProfileSetting returns the active profile's settings.

### `SetProfileSetting`

```go
func (s *Rubino) SetProfileSetting(ctx context.Context, profileID string, settings map[string]any) (*events.Event, error)
```

SetProfileSetting updates the active profile's settings. `settings` is
the raw settings map (key/value pairs the Rubino API accepts).

### `GetSettingAdminPermissions`

```go
func (s *Rubino) GetSettingAdminPermissions(ctx context.Context) (*events.Event, error)
```

GetSettingAdminPermissions returns the available admin permission set.

### `GetProfileAdmins`

```go
func (s *Rubino) GetProfileAdmins(ctx context.Context, profileID string) (*events.Event, error)
```

GetProfileAdmins lists the admins of a profile.

### `SetProfileAdmin`

```go
func (s *Rubino) SetProfileAdmin(ctx context.Context, profileID, targetProfileID, action string, accessList []string) (*events.Event, error)
```

SetProfileAdmin promotes a profile admin with the given access list.
`accessList` is the set of permission names (e.g. ["EditProfile",
"AddPost"]); `action` is typically "SetAdmin".

### `UnsetProfileAdmin`

```go
func (s *Rubino) UnsetProfileAdmin(ctx context.Context, profileID, targetProfileID string) (*events.Event, error)
```

UnsetProfileAdmin demotes a profile admin.

### `AcceptProfileAdmin`

```go
func (s *Rubino) AcceptProfileAdmin(ctx context.Context, profileID, targetProfileID string) (*events.Event, error)
```

AcceptProfileAdmin accepts an incoming admin invitation for a profile.

### `DenyProfileAdmin`

```go
func (s *Rubino) DenyProfileAdmin(ctx context.Context, profileID, targetProfileID string) (*events.Event, error)
```

DenyProfileAdmin denies an incoming admin invitation.

### `RemoveProfileAdmin`

```go
func (s *Rubino) RemoveProfileAdmin(ctx context.Context, profileID, targetProfileID string) (*events.Event, error)
```

RemoveProfileAdmin removes a profile admin (used by the profile owner).

### `GetMyCloseFriends`

```go
func (s *Rubino) GetMyCloseFriends(ctx context.Context) (*events.Event, error)
```

GetMyCloseFriends returns the caller's close-friends list.

### `SetCloseFriends`

```go
func (s *Rubino) SetCloseFriends(ctx context.Context, profileID string, targetProfileIDs []string) (*events.Event, error)
```

SetCloseFriends sets the close-friends list for a profile.

### `UnsetCloseFriends`

```go
func (s *Rubino) UnsetCloseFriends(ctx context.Context, profileID string, targetProfileIDs []string) (*events.Event, error)
```

UnsetCloseFriends removes the given profiles from the close-friends list.

### `SearchPost`

```go
func (s *Rubino) SearchPost(ctx context.Context, searchText, usage, startID string) (*events.Event, error)
```

SearchPost searches posts by text. `usage` is "Search" by default;
"Hashtag" switches to hashtag search semantics.

### `SearchHashTag`

```go
func (s *Rubino) SearchHashTag(ctx context.Context, searchText, startID string) (*events.Event, error)
```

SearchHashTag searches posts by hashtag.

### `SearchLocation`

```go
func (s *Rubino) SearchLocation(ctx context.Context, searchText string, lat, lng float64, startID string) (*events.Event, error)
```

SearchLocation searches posts near a geographic point.

### `GetPostsByIDs`

```go
func (s *Rubino) GetPostsByIDs(ctx context.Context, profileID string, postIDs []string) (*events.Event, error)
```

GetPostsByIDs fetches a batch of posts by ID.

### `UpdatePost`

```go
func (s *Rubino) UpdatePost(ctx context.Context, profileID, postID string, fields map[string]any, updateParameters []string) (*events.Event, error)
```

UpdatePost edits a post's caption / settings. `updateParameters` lists
the fields being changed (e.g. ["caption"]).

### `UpdatePostShareCount`

```go
func (s *Rubino) UpdatePostShareCount(ctx context.Context, profileID, postID string) (*events.Event, error)
```

UpdatePostShareCount records a share event for a post.

### `PostCollaboratorAction`

```go
func (s *Rubino) PostCollaboratorAction(ctx context.Context, profileID, postID, targetProfileID, action string) (*events.Event, error)
```

PostCollaboratorAction adds/removes a collaborator on a post.
`action` is "AddCollaborator" or "RemoveCollaborator".

### `SetPostFeedback`

```go
func (s *Rubino) SetPostFeedback(ctx context.Context, profileID, postID, feedback string) (*events.Event, error)
```

SetPostFeedback submits feedback (e.g. not-interested) on a post.

### `GetFeedData`

```go
func (s *Rubino) GetFeedData(ctx context.Context, profileID, startID string, limit int) (*events.Event, error)
```

GetFeedData returns the caller's feed.

### `ViewFeedData`

```go
func (s *Rubino) ViewFeedData(ctx context.Context, profileID string, feedRowIDs []string) (*events.Event, error)
```

ViewFeedData records a feed-view event (used for analytics).

### `RemoveFollower`

```go
func (s *Rubino) RemoveFollower(ctx context.Context, profileID, targetProfileID string) (*events.Event, error)
```

RemoveFollower removes a follower from a profile.

### `NotInterestedProfile`

```go
func (s *Rubino) NotInterestedProfile(ctx context.Context, profileID, targetProfileID string) (*events.Event, error)
```

NotInterestedProfile marks a profile as "not interested".

### `GetBoostURL`

```go
func (s *Rubino) GetBoostURL(ctx context.Context, profileID string) (*events.Event, error)
```

GetBoostURL returns the profile's boost (promotion) URL.

### `GetMetricsConfig`

```go
func (s *Rubino) GetMetricsConfig(ctx context.Context) (*events.Event, error)
```

GetMetricsConfig returns the metrics configuration.

### `TrackClientMetrics`

```go
func (s *Rubino) TrackClientMetrics(ctx context.Context, payload map[string]any) (*events.Event, error)
```

TrackClientMetrics submits a client-metrics payload.

### `AddImpressionIDs`

```go
func (s *Rubino) AddImpressionIDs(ctx context.Context, profileID string, impressionIDs []string) (*events.Event, error)
```

AddImpressionIDs records impression IDs for analytics.

### `AddUsageTime`

```go
func (s *Rubino) AddUsageTime(ctx context.Context, profileID string, seconds int) (*events.Event, error)
```

AddUsageTime records usage time for a profile.

### `AcceptMessageOwner`

```go
func (s *Rubino) AcceptMessageOwner(ctx context.Context, profileID, postID string) (*events.Event, error)
```

AcceptMessageOwner accepts a message-owner request for a post.

### `DenyMessageOwner`

```go
func (s *Rubino) DenyMessageOwner(ctx context.Context, profileID, postID string) (*events.Event, error)
```

DenyMessageOwner denies a message-owner request.

### `RemoveMessageOwner`

```go
func (s *Rubino) RemoveMessageOwner(ctx context.Context, profileID, postID string) (*events.Event, error)
```

RemoveMessageOwner removes the message owner of a post.

### `ProfileCreatedAt`

```go
func ProfileCreatedAt(profileID string) (time.Time, error)
```

ProfileCreatedAt decodes the creation timestamp encoded in the first 8
hex characters of a Rubino profile_id. Rubino profile IDs are 24-hex
strings whose first 4 bytes are a big-endian Unix timestamp. It makes
no API call.

Returns an error when the ID is too short or not valid hex.

### `ProfileCreatedAtString`

```go
func ProfileCreatedAtString(profileID string) (string, error)
```

ProfileCreatedAtString is a convenience wrapper returning the RFC-3339
creation timestamp of a profile_id.
