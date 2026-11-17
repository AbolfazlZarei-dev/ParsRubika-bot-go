package rubino

import (
	"context"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// ===== Missing Rubino methods =====================
//
// Endpoints that complete the studio package.

// GetBlockedProfiles returns the list of profiles the current user has
// blocked.
func (s *Rubino) GetBlockedProfiles(ctx context.Context, profileID string, limit int, sort, equal string) (*events.Event, error) {
	if limit <= 0 {
		limit = 50
	}
	if sort == "" {
		sort = "FromMax"
	}
	return s.request(ctx, "getBlockedProfiles", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
		"sort":       sort,
		"equal":      equal != "",
	})
}

// GetNewEvents returns recent Rubino events (notifications).
func (s *Rubino) GetNewEvents(ctx context.Context, profileID string, limit int, sort string) (*events.Event, error) {
	if limit <= 0 {
		limit = 50
	}
	if sort == "" {
		sort = "FromMax"
	}
	return s.request(ctx, "getNewEvents", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
		"sort":       sort,
	})
}

// EnableProfileNotifications turns on notifications for a profile.
func (s *Rubino) EnableProfileNotifications(ctx context.Context, profileID string) (*events.Event, error) {
	mute := false
	return s.request(ctx, "updateProfile", map[string]any{
		"profile_id": s.profileOr(profileID),
		"is_mute":    &mute,
	})
}

// MuteProfileNotifications turns off notifications for a profile.
func (s *Rubino) MuteProfileNotifications(ctx context.Context, profileID string) (*events.Event, error) {
	mute := true
	return s.request(ctx, "updateProfile", map[string]any{
		"profile_id": s.profileOr(profileID),
		"is_mute":    &mute,
	})
}

// DeleteStory deletes one or more stories by ID.
func (s *Rubino) DeleteStoryByID(ctx context.Context, profileID string, storyIDs []string) (*events.Event, error) {
	return s.request(ctx, "deleteStory", map[string]any{
		"profile_id": s.profileOr(profileID),
		"story_ids":  storyIDs,
	})
}

// SetProfileStatus updates the profile's status (Public/Private).
func (s *Rubino) SetProfileStatus(ctx context.Context, profileID, status string) (*events.Event, error) {
	return s.request(ctx, "updateProfile", map[string]any{
		"profile_id":     s.profileOr(profileID),
		"profile_status": status,
	})
}

// EnablePrivateMessagesByProfile allows other profiles to message this one.
func (s *Rubino) EnablePrivateMessagesByProfile(ctx context.Context, profileID string) (*events.Event, error) {
	allowed := true
	return s.request(ctx, "updateProfile", map[string]any{
		"profile_id":         s.profileOr(profileID),
		"is_message_allowed": &allowed,
	})
}

// DisablePrivateMessagesByProfile blocks other profiles from messaging.
func (s *Rubino) DisablePrivateMessagesByProfile(ctx context.Context, profileID string) (*events.Event, error) {
	allowed := false
	return s.request(ctx, "updateProfile", map[string]any{
		"profile_id":         s.profileOr(profileID),
		"is_message_allowed": &allowed,
	})
}

// RemoveStoryFromHighlightByID removes specific stories from a highlight.
func (s *Rubino) RemoveStoryFromHighlightByID(ctx context.Context, profileID, highlightID string, removeStoryIDs []string) (*events.Event, error) {
	return s.request(ctx, "editHighlight", map[string]any{
		"profile_id":       s.profileOr(profileID),
		"highlight_id":     highlightID,
		"remove_story_ids": removeStoryIDs,
	})
}

// DeleteCommentByIDs deletes a comment from a post.
func (s *Rubino) DeleteCommentByIDs(ctx context.Context, profileID, postID, commentID string) (*events.Event, error) {
	return s.request(ctx, "removeRecord", map[string]any{
		"profile_id": s.profileOr(profileID),
		"model":      "Comment",
		"record_id":  commentID,
		"post_id":    postID,
	})
}

// DeletePage removes a Rubino page/record.
func (s *Rubino) DeletePage(ctx context.Context, profileID, recordID string) (*events.Event, error) {
	return s.request(ctx, "removeRecord", map[string]any{
		"profile_id": s.profileOr(profileID),
		"model":      "Profile",
		"record_id":  recordID,
	})
}

// SetReportRecord reports a post or profile. reason is a string
// description.
func (s *Rubino) SetReportRecord(ctx context.Context, profileID, reason, postID, postProfileID string) (*events.Event, error) {
	input := map[string]any{
		"profile_id": s.profileOr(profileID),
		"reason":     reason,
	}
	if postID != "" {
		input["post_id"] = postID
	}
	if postProfileID != "" {
		input["post_profile_id"] = postProfileID
	}
	return s.request(ctx, "setReportRecord", input)
}
