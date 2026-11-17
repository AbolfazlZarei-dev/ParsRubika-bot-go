# `rubino` — Methods stories

_Source: [`methods_stories.go`](../../ParsRubika/rubino/methods_stories.go)_

## Symbols

### `GetProfilesStories`

```go
func (s *Rubino) GetProfilesStories(ctx context.Context, profileID string, limit int) (*Event, error)
```

GetProfilesStories returns the stories from the user's contacts.

### `GetProfilesStoryList`

```go
func (s *Rubino) GetProfilesStoryList(ctx context.Context, profileID string, profileStoryIDs []map[string]any) (*Event, error)
```

GetProfilesStoryList returns multiple stories for several profiles.
`profileStoryIDs` is a slice of `{profile_id, story_ids}` maps.

### `GetStory`

```go
func (s *Rubino) GetStory(ctx context.Context, profileID, storyProfileID string, storyIDs []string) (*Event, error)
```

GetStory returns a single story.

### `GetStoryViewers`

```go
func (s *Rubino) GetStoryViewers(ctx context.Context, profileID, storyID string, limit int) (*Event, error)
```

GetStoryViewers returns the viewers of a story.

### `AddViewStory`

```go
func (s *Rubino) AddViewStory(ctx context.Context, profileID, storyProfileID string, storyIDs []string) (*Event, error)
```

AddViewStory adds a view to a story.

### `GetMyArchiveStories`

```go
func (s *Rubino) GetMyArchiveStories(ctx context.Context, profileID string, limit int, sort, equal string) (*Event, error)
```

GetMyArchiveStories returns the user's archived stories.

### `DeleteStory`

```go
func (s *Rubino) DeleteStory(ctx context.Context, profileID string, storyIDs []string) (*Event, error)
```

DeleteStory deletes one or more stories. The Rubino API expects the
plural `story_ids` key for batch deletes.

### `GetProfileHighlights`

```go
func (s *Rubino) GetProfileHighlights(ctx context.Context, profileID, targetProfileID string, limit int, sort, equal string) (*Event, error)
```

GetProfileHighlights returns the highlights of a target profile.

### `CreateHighlight`

```go
func (s *Rubino) CreateHighlight(ctx context.Context, profileID, highlightName string, storyIDs []string, highlightCoverPicture any) (*Event, error)
```

CreateHighlight creates a highlight.

### `HighlightStory`

```go
func (s *Rubino) HighlightStory(ctx context.Context, profileID, highlightID, storyID string) (*Event, error)
```

HighlightStory adds a story to an existing highlight.

### `RemoveStoryFromHighlight`

```go
func (s *Rubino) RemoveStoryFromHighlight(ctx context.Context, profileID, highlightID string, removeStoryIDs []string) (*Event, error)
```

RemoveStoryFromHighlight removes a story from a highlight.

### `GetHighlightStoryIDs`

```go
func (s *Rubino) GetHighlightStoryIDs(ctx context.Context, profileID, highlightID, targetProfileID string) (*Event, error)
```

GetHighlightStoryIDs returns the story IDs in a highlight.

### `GetHighlightStories`

```go
func (s *Rubino) GetHighlightStories(ctx context.Context, profileID, highlightID, targetProfileID string, storyIDs []string) (*Event, error)
```

GetHighlightStories returns the stories inside a highlight.

### `GetPostsByHashTag`

```go
func (s *Rubino) GetPostsByHashTag(ctx context.Context, profileID, hashtag, startID string) (*Event, error)
```

GetPostsByHashTag returns posts tagged with the supplied hashtag.

### `GetHashTagTrend`

```go
func (s *Rubino) GetHashTagTrend(ctx context.Context, profileID string, limit int, sort, equal string) (*Event, error)
```

GetHashTagTrend returns trending hashtags.

### `SearchProfile`

```go
func (s *Rubino) SearchProfile(ctx context.Context, profileID, username string, limit int, sort, equal string) (*Event, error)
```

SearchProfile searches for profiles by username.

### `SetReportPost`

```go
func (s *Rubino) SetReportPost(ctx context.Context, profileID, postID, reason string) (*Event, error)
```

SetReportPost reports a post. `reason` is a short reason string.

### `SetReportPage`

```go
func (s *Rubino) SetReportPage(ctx context.Context, profileID, reason string) (*Event, error)
```

SetReportPage reports a page.

### `GetSuggested`

```go
func (s *Rubino) GetSuggested(ctx context.Context, profileID string, limit int, sort string) (*Event, error)
```

GetSuggested returns suggested profiles to follow.

### `GetMyProfilePosts`

```go
func (s *Rubino) GetMyProfilePosts(ctx context.Context, profileID string, limit int, sort string) (*Event, error)
```

GetMyProfilePosts returns the posts of the user's own profile.
