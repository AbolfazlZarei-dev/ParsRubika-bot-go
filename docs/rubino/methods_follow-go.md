# `rubino` — Methods follow

_Source: [`methods_follow.go`](../../ParsRubika/rubino/methods_follow.go)_

## Symbols

### `Follow`

```go
func (s *Rubino) Follow(ctx context.Context, profileID, followeeID string) (*Event, error)
```

Follow requests to follow a profile.

### `Unfollow`

```go
func (s *Rubino) Unfollow(ctx context.Context, profileID, followeeID string) (*Event, error)
```

Unfollow requests to unfollow a profile.

### `AcceptFollowRequest`

```go
func (s *Rubino) AcceptFollowRequest(ctx context.Context, profileID, requestID string) (*Event, error)
```

AcceptFollowRequest accepts a pending follow request.

### `DeclineFollowRequest`

```go
func (s *Rubino) DeclineFollowRequest(ctx context.Context, profileID, requestID string) (*Event, error)
```

DeclineFollowRequest declines a pending follow request.

### `GetNewFollowRequests`

```go
func (s *Rubino) GetNewFollowRequests(ctx context.Context, profileID string, limit int, sort string) (*Event, error)
```

GetNewFollowRequests returns the pending follow requests.

### `GetProfileFollowers`

```go
func (s *Rubino) GetProfileFollowers(ctx context.Context, profileID, targetProfileID string, limit int, sort, equal string) (*Event, error)
```

GetProfileFollowers returns the followers of a target profile.

### `GetProfileFollowings`

```go
func (s *Rubino) GetProfileFollowings(ctx context.Context, profileID, targetProfileID string, limit int, sort, equal string) (*Event, error)
```

GetProfileFollowings returns the followings of a target profile.

### `SearchFollower`

```go
func (s *Rubino) SearchFollower(ctx context.Context, profileID, targetProfileID, username string, limit int, sort, equal string) (*Event, error)
```

SearchFollower searches the followers of a target profile.

### `SearchFollowing`

```go
func (s *Rubino) SearchFollowing(ctx context.Context, profileID, targetProfileID, username string, limit int, sort, equal string) (*Event, error)
```

SearchFollowing searches the followings of a target profile.
