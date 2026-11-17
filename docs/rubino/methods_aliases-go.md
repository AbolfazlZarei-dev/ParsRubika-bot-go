# `rubino` — Methods aliases

_Source: [`methods_aliases.go`](../../ParsRubika/rubino/methods_aliases.go)_

## Symbols

### `GetPostLikesByIDs`

```go
func (s *Rubino) GetPostLikesByIDs(ctx context.Context, profileID, postID, postProfileID string, limit int, sort, equal string) (*events.Event, error)
```

GetPostLikesByIDs returns the users who liked a post. It is the
explicit-parameter form of GetPostLikes (profile, post and
post-owner IDs are passed separately).

### `GetCommentRepliesByIDs`

```go
func (s *Rubino) GetCommentRepliesByIDs(ctx context.Context, profileID, postID, commentID string, limit int, sort string) (*events.Event, error)
```

GetCommentRepliesByIDs returns replies to a comment. It is the
explicit-parameter form of GetCommentReplies.

### `GetLikedCommentProfilesByIDs`

```go
func (s *Rubino) GetLikedCommentProfilesByIDs(ctx context.Context, profileID, postID, commentID string, limit int, sort string) (*events.Event, error)
```

GetLikedCommentProfilesByIDs returns profiles that liked a comment.
(GetLikedCommentProfiles already exists; this is an alias.)

### `LikeCommentAction`

```go
func (s *Rubino) LikeCommentAction(ctx context.Context, profileID, commentID, postID, actionType string) (*events.Event, error)
```

LikeCommentAction performs a like/unlike action on a comment.

### `SearchProfileByIDs`

```go
func (s *Rubino) SearchProfileByIDs(ctx context.Context, profileID, username string, limit int, sort, equal string) (*events.Event, error)
```

SearchProfileByIDs searches for profiles by username. It is the
explicit-parameter form of SearchProfile.

### `AddPostViewCountByIDs`

```go
func (s *Rubino) AddPostViewCountByIDs(ctx context.Context, profileID, postID, postProfileID string) (*events.Event, error)
```

AddPostViewCountByIDs registers a view on a post. (AddPostViewCount
already exists; this is an alias with explicit profileID.)

### `GetProfilesStoriesByIDs`

```go
func (s *Rubino) GetProfilesStoriesByIDs(ctx context.Context, profileID string, limit int) (*events.Event, error)
```

GetProfilesStoriesByIDs returns stories for the given profile.
(GetProfilesStories already exists; this is an alias with limit.)

### `RequestFollow`

```go
func (s *Rubino) RequestFollow(ctx context.Context, profileID, followeeID, fType string) (*events.Event, error)
```

RequestFollow performs a follow/unfollow action. Follow and
Unfollow wrap it; this exposes the underlying requestFollow endpoint
directly.

### `GetProfileFollowersByType`

```go
func (s *Rubino) GetProfileFollowersByType(ctx context.Context, profileID, targetProfileID, fType string, limit int, sort, equal string) (*events.Event, error)
```

GetProfileFollowersByType returns followers or followings based on
fType ("Follower" or "Following"). (GetProfileFollowers/Followings
already exist; this exposes the underlying endpoint directly.)
