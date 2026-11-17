package rubino

import "context"

// ===== Stories / highlights ==========================================

// GetProfilesStories returns the stories from the user's contacts.
func (s *Rubino) GetProfilesStories(ctx context.Context, profileID string, limit int) (*Event, error) {
	return s.request(ctx, "getProfilesStories", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
	})
}

// GetProfilesStoryList returns multiple stories for several profiles.
// `profileStoryIDs` is a slice of `{profile_id, story_ids}` maps.
func (s *Rubino) GetProfilesStoryList(ctx context.Context, profileID string, profileStoryIDs []map[string]any) (*Event, error) {
	return s.request(ctx, "getProfilesStoryList", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"profile_story_ids": profileStoryIDs,
	})
}

// GetStory returns a single story.
func (s *Rubino) GetStory(ctx context.Context, profileID, storyProfileID string, storyIDs []string) (*Event, error) {
	return s.request(ctx, "getStory", map[string]any{
		"profile_id":       s.profileOr(profileID),
		"story_profile_id": storyProfileID,
		"story_ids":        storyIDs,
	})
}

// GetStoryViewers returns the viewers of a story.
func (s *Rubino) GetStoryViewers(ctx context.Context, profileID, storyID string, limit int) (*Event, error) {
	return s.request(ctx, "getStoryViewers", map[string]any{
		"profile_id": s.profileOr(profileID),
		"story_id":   storyID,
		"limit":      limit,
	})
}

// AddViewStory adds a view to a story.
func (s *Rubino) AddViewStory(ctx context.Context, profileID, storyProfileID string, storyIDs []string) (*Event, error) {
	return s.request(ctx, "addViewStory", map[string]any{
		"profile_id":       s.profileOr(profileID),
		"story_profile_id": storyProfileID,
		"story_ids":        storyIDs,
	})
}

// GetMyArchiveStories returns the user's archived stories.
func (s *Rubino) GetMyArchiveStories(ctx context.Context, profileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getMyArchiveStories", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
		"equal":      equal,
	})
}

// DeleteStory deletes one or more stories. The Rubino API expects the
// plural `story_ids` key for batch deletes.
func (s *Rubino) DeleteStory(ctx context.Context, profileID string, storyIDs []string) (*Event, error) {
	return s.request(ctx, "deleteStory", map[string]any{
		"profile_id": s.profileOr(profileID),
		"story_ids":  storyIDs,
	})
}

// GetProfileHighlights returns the highlights of a target profile.
func (s *Rubino) GetProfileHighlights(ctx context.Context, profileID, targetProfileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getProfileHighlights", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"limit":             limit,
		"sort":              ifEmpty(sort, "latest"),
		"equal":             equal,
	})
}

// CreateHighlight creates a highlight.
func (s *Rubino) CreateHighlight(ctx context.Context, profileID, highlightName string, storyIDs []string, highlightCoverPicture any) (*Event, error) {
	return s.request(ctx, "addHighlight", map[string]any{
		"profile_id":              s.profileOr(profileID),
		"highlight_name":          highlightName,
		"story_ids":               storyIDs,
		"highlight_cover_picture": highlightCoverPicture,
	})
}

// HighlightStory adds a story to an existing highlight.
func (s *Rubino) HighlightStory(ctx context.Context, profileID, highlightID, storyID string) (*Event, error) {
	return s.request(ctx, "highlightStory", map[string]any{
		"profile_id":   s.profileOr(profileID),
		"highlight_id": highlightID,
		"story_id":     storyID,
	})
}

// RemoveStoryFromHighlight removes a story from a highlight.
func (s *Rubino) RemoveStoryFromHighlight(ctx context.Context, profileID, highlightID string, removeStoryIDs []string) (*Event, error) {
	return s.request(ctx, "editHighlight", map[string]any{
		"profile_id":       s.profileOr(profileID),
		"highlight_id":     highlightID,
		"remove_story_ids": removeStoryIDs,
	})
}

// GetHighlightStoryIDs returns the story IDs in a highlight.
func (s *Rubino) GetHighlightStoryIDs(ctx context.Context, profileID, highlightID, targetProfileID string) (*Event, error) {
	return s.request(ctx, "getHighlightStoryIds", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"highlight_id":      highlightID,
		"target_profile_id": targetProfileID,
	})
}

// GetHighlightStories returns the stories inside a highlight.
func (s *Rubino) GetHighlightStories(ctx context.Context, profileID, highlightID, targetProfileID string, storyIDs []string) (*Event, error) {
	return s.request(ctx, "getHighlightStories", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"highlight_id":      highlightID,
		"target_profile_id": targetProfileID,
		"story_ids":         storyIDs,
	})
}

// GetPostsByHashTag returns posts tagged with the supplied hashtag.
func (s *Rubino) GetPostsByHashTag(ctx context.Context, profileID, hashtag, startID string) (*Event, error) {
	return s.request(ctx, "getPostsByHashTag", map[string]any{
		"profile_id": s.profileOr(profileID),
		"hashtag":    hashtag,
		"start_id":   startID,
	})
}

// GetHashTagTrend returns trending hashtags.
func (s *Rubino) GetHashTagTrend(ctx context.Context, profileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getHashTagTrend", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
		"equal":      equal,
	})
}

// SearchProfile searches for profiles by username.
func (s *Rubino) SearchProfile(ctx context.Context, profileID, username string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "searchProfile", map[string]any{
		"profile_id": s.profileOr(profileID),
		"username":   username,
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
		"equal":      equal,
	})
}

// SetReportPost reports a post. `reason` is a short reason string.
func (s *Rubino) SetReportPost(ctx context.Context, profileID, postID, reason string) (*Event, error) {
	return s.request(ctx, "setReportRecord", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_id":    postID,
		"reason":     reason,
	})
}

// SetReportPage reports a page.
func (s *Rubino) SetReportPage(ctx context.Context, profileID, reason string) (*Event, error) {
	return s.request(ctx, "setReportRecord", map[string]any{
		"profile_id": s.profileOr(profileID),
		"reason":     reason,
	})
}

// GetSuggested returns suggested profiles to follow.
func (s *Rubino) GetSuggested(ctx context.Context, profileID string, limit int, sort string) (*Event, error) {
	return s.request(ctx, "getSuggested", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
	})
}

// GetMyProfilePosts returns the posts of the user's own profile.
func (s *Rubino) GetMyProfilePosts(ctx context.Context, profileID string, limit int, sort string) (*Event, error) {
	return s.request(ctx, "getMyProfilePosts", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
	})
}
