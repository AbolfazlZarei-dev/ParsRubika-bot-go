# `rubino` — Methods extra

_Source: [`methods_extra.go`](../../ParsRubika/rubino/methods_extra.go)_

## Symbols

### `GetBlockedProfiles`

```go
func (s *Rubino) GetBlockedProfiles(ctx context.Context, profileID string, limit int, sort, equal string) (*events.Event, error)
```

GetBlockedProfiles returns the list of profiles the current user has
blocked.

### `GetNewEvents`

```go
func (s *Rubino) GetNewEvents(ctx context.Context, profileID string, limit int, sort string) (*events.Event, error)
```

GetNewEvents returns recent Rubino events (notifications).

### `EnableProfileNotifications`

```go
func (s *Rubino) EnableProfileNotifications(ctx context.Context, profileID string) (*events.Event, error)
```

EnableProfileNotifications turns on notifications for a profile.

### `MuteProfileNotifications`

```go
func (s *Rubino) MuteProfileNotifications(ctx context.Context, profileID string) (*events.Event, error)
```

MuteProfileNotifications turns off notifications for a profile.

### `DeleteStoryByID`

```go
func (s *Rubino) DeleteStoryByID(ctx context.Context, profileID string, storyIDs []string) (*events.Event, error)
```

DeleteStory deletes one or more stories by ID.

### `SetProfileStatus`

```go
func (s *Rubino) SetProfileStatus(ctx context.Context, profileID, status string) (*events.Event, error)
```

SetProfileStatus updates the profile's status (Public/Private).

### `EnablePrivateMessagesByProfile`

```go
func (s *Rubino) EnablePrivateMessagesByProfile(ctx context.Context, profileID string) (*events.Event, error)
```

EnablePrivateMessagesByProfile allows other profiles to message this one.

### `DisablePrivateMessagesByProfile`

```go
func (s *Rubino) DisablePrivateMessagesByProfile(ctx context.Context, profileID string) (*events.Event, error)
```

DisablePrivateMessagesByProfile blocks other profiles from messaging.

### `RemoveStoryFromHighlightByID`

```go
func (s *Rubino) RemoveStoryFromHighlightByID(ctx context.Context, profileID, highlightID string, removeStoryIDs []string) (*events.Event, error)
```

RemoveStoryFromHighlightByID removes specific stories from a highlight.

### `DeleteCommentByIDs`

```go
func (s *Rubino) DeleteCommentByIDs(ctx context.Context, profileID, postID, commentID string) (*events.Event, error)
```

DeleteCommentByIDs deletes a comment from a post.

### `DeletePage`

```go
func (s *Rubino) DeletePage(ctx context.Context, profileID, recordID string) (*events.Event, error)
```

DeletePage removes a Rubino page/record.

### `SetReportRecord`

```go
func (s *Rubino) SetReportRecord(ctx context.Context, profileID, reason, postID, postProfileID string) (*events.Event, error)
```

SetReportRecord reports a post or profile. reason is a string
description.
