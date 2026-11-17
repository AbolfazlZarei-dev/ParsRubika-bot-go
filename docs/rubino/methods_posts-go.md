# `rubino` — Methods posts

_Source: [`methods_posts.go`](../../ParsRubika/rubino/methods_posts.go)_

## Symbols

### `LikePostAction`

```go
func (s *Rubino) LikePostAction(ctx context.Context, profileID, postID, postProfileID, actionType string) (*Event, error)
```

LikePostAction likes a post.

### `Like`

```go
func (s *Rubino) Like(ctx context.Context, postID, postProfileID string) (*Event, error)
```

Like is a convenience wrapper for LikePostAction with action="Like".

### `Unlike`

```go
func (s *Rubino) Unlike(ctx context.Context, postID, postProfileID string) (*Event, error)
```

Unlike is a convenience wrapper for LikePostAction with action="Unlike".

### `GetComments`

```go
func (s *Rubino) GetComments(ctx context.Context, profileID, postID, postProfileID string, limit int, sort, equal string) (*Event, error)
```

GetComments returns the comments on a post.

### `AddComment`

```go
func (s *Rubino) AddComment(ctx context.Context, profileID, content, postID, postProfileID string) (*Event, error)
```

AddComment comments on a post.

### `AddReplyComment`

```go
func (s *Rubino) AddReplyComment(ctx context.Context, profileID, content, commentID, postID string) (*Event, error)
```

AddReplyComment replies to a comment on a post.

### `LikeComment`

```go
func (s *Rubino) LikeComment(ctx context.Context, profileID, commentID, postID string) (*Event, error)
```

LikeComment likes a comment.

### `UnlikeComment`

```go
func (s *Rubino) UnlikeComment(ctx context.Context, profileID, commentID, postID string) (*Event, error)
```

UnlikeComment unlikes a comment.

### `GetCommentReplies`

```go
func (s *Rubino) GetCommentReplies(ctx context.Context, profileID, postID, commentID string, limit int, sort string) (*Event, error)
```

GetCommentReplies returns the replies to a comment.

### `DeleteComment`

```go
func (s *Rubino) DeleteComment(ctx context.Context, profileID, commentID, postID string) (*Event, error)
```

DeleteComment deletes a comment. The removeRecord API expects a
`record_id` + `model` pair (model="Comment"); the post_id is also
required so the server can locate the comment.

### `GetPostLikes`

```go
func (s *Rubino) GetPostLikes(ctx context.Context, profileID, postID, postProfileID string, limit int, sort, equal string) (*Event, error)
```

GetPostLikes returns the users who liked a post.

### `GetLikedCommentProfiles`

```go
func (s *Rubino) GetLikedCommentProfiles(ctx context.Context, profileID, postID, commentID string, limit int, sort string) (*Event, error)
```

GetLikedCommentProfiles returns the profiles who liked a comment.

### `AddPostViewCount`

```go
func (s *Rubino) AddPostViewCount(ctx context.Context, postID, postProfileID string) (*Event, error)
```

AddPostViewCount registers a view on a post.

### `AddPostViewTime`

```go
func (s *Rubino) AddPostViewTime(ctx context.Context, postID, postProfileID string, duration int) (*Event, error)
```

AddPostViewTime registers a view duration on a post.

### `GetRecentFollowingPosts`

```go
func (s *Rubino) GetRecentFollowingPosts(ctx context.Context, profileID string, limit int, sort, equal string) (*Event, error)
```

GetRecentFollowingPosts returns the feed of recent posts from
followings.

### `GetProfilePosts`

```go
func (s *Rubino) GetProfilePosts(ctx context.Context, profileID, targetProfileID string, limit int, sort, equal string) (*Event, error)
```

GetProfilePosts returns the posts of a specific target profile.

### `GetBookmarkedPosts`

```go
func (s *Rubino) GetBookmarkedPosts(ctx context.Context, profileID string, limit int, sort, equal string) (*Event, error)
```

GetBookmarkedPosts returns the user's saved posts.

### `SavePost`

```go
func (s *Rubino) SavePost(ctx context.Context, profileID, postID, postProfileID string) (*Event, error)
```

SavePost bookmarks a post.

### `UnsavePost`

```go
func (s *Rubino) UnsavePost(ctx context.Context, profileID, postID, postProfileID string) (*Event, error)
```

UnsavePost removes a bookmark.

### `GetExplorePosts`

```go
func (s *Rubino) GetExplorePosts(ctx context.Context, profileID string, limit int, sort, equal string) (*Event, error)
```

GetExplorePosts returns the explore feed.

### `GetRelatedExplorePost`

```go
func (s *Rubino) GetRelatedExplorePost(ctx context.Context, profileID, postID, postProfileID string, limit int) (*Event, error)
```

GetRelatedExplorePost returns related explore posts.

### `GetExplorePostTopics`

```go
func (s *Rubino) GetExplorePostTopics(ctx context.Context, profileID string) (*Event, error)
```

GetExplorePostTopics returns the explore post topics.

### `GetTaggedPosts`

```go
func (s *Rubino) GetTaggedPosts(ctx context.Context, profileID, targetProfileID string, limit int, sort, equal string) (*Event, error)
```

GetTaggedPosts returns the tagged posts of a target profile.

### `GetPostByShareLink`

```go
func (s *Rubino) GetPostByShareLink(ctx context.Context, profileID, shareLink string) (*Event, error)
```

GetPostByShareLink resolves a `rubika.ir/post/...` share link.

### `DeletePost`

```go
func (s *Rubino) DeletePost(ctx context.Context, profileID, postID string) (*Event, error)
```

DeletePost deletes a post. The removeRecord API expects a
`record_id` + `model` pair (model="Post").

### `AddPost`

```go
func (s *Rubino) AddPost(ctx context.Context, profileID, postType, fileID, caption, fileName string, width, height int) (*Event, error)
```

AddPost publishes a new post. `postType` is "Picture" or "Video".
The file must already be uploaded (see RequestUploadFile + UploadFile).
