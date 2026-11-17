package rubino

import (
	"context"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// ===== Additional Rubino methods =====================================
//
// Explicit-parameter variants and a few endpoints not covered elsewhere.

// GetPostLikesByIDs returns the users who liked a post. It is the
// explicit-parameter form of GetPostLikes (profile, post and
// post-owner IDs are passed separately).
func (s *Rubino) GetPostLikesByIDs(ctx context.Context, profileID, postID, postProfileID string, limit int, sort, equal string) (*events.Event, error) {
	if limit <= 0 {
		limit = 50
	}
	if sort == "" {
		sort = "FromMax"
	}
	return s.request(ctx, "getPostLikes", map[string]any{
		"profile_id":      s.profileOr(profileID),
		"post_id":         postID,
		"post_profile_id": postProfileID,
		"limit":           limit,
		"sort":            sort,
		"equal":           equal != "",
	})
}

// GetCommentRepliesByIDs returns replies to a comment. It is the
// explicit-parameter form of GetCommentReplies.
func (s *Rubino) GetCommentRepliesByIDs(ctx context.Context, profileID, postID, commentID string, limit int, sort string) (*events.Event, error) {
	if limit <= 0 {
		limit = 50
	}
	if sort == "" {
		sort = "FromMax"
	}
	return s.request(ctx, "getCommentReplies", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_id":    postID,
		"comment_id": commentID,
		"limit":      limit,
		"sort":       sort,
	})
}

// GetLikedCommentProfilesByIDs returns profiles that liked a comment.
// (GetLikedCommentProfiles already exists; this is an alias.)
func (s *Rubino) GetLikedCommentProfilesByIDs(ctx context.Context, profileID, postID, commentID string, limit int, sort string) (*events.Event, error) {
	if limit <= 0 {
		limit = 50
	}
	if sort == "" {
		sort = "FromMax"
	}
	return s.request(ctx, "getLikedCommentProfiles", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_id":    postID,
		"comment_id": commentID,
		"limit":      limit,
		"sort":       sort,
	})
}

// LikeCommentAction performs a like/unlike action on a comment.
func (s *Rubino) LikeCommentAction(ctx context.Context, profileID, commentID, postID, actionType string) (*events.Event, error) {
	return s.request(ctx, "likeCommentAction", map[string]any{
		"profile_id":  s.profileOr(profileID),
		"comment_id":  commentID,
		"post_id":     postID,
		"action_type": actionType,
	})
}

// SearchProfileByIDs searches for profiles by username. It is the
// explicit-parameter form of SearchProfile.
func (s *Rubino) SearchProfileByIDs(ctx context.Context, profileID, username string, limit int, sort, equal string) (*events.Event, error) {
	if limit <= 0 {
		limit = 50
	}
	if sort == "" {
		sort = "FromMax"
	}
	return s.request(ctx, "searchProfile", map[string]any{
		"profile_id": s.profileOr(profileID),
		"username":   username,
		"limit":      limit,
		"sort":       sort,
		"equal":      equal != "",
	})
}

// AddPostViewCountByIDs registers a view on a post. (AddPostViewCount
// already exists; this is an alias with explicit profileID.)
func (s *Rubino) AddPostViewCountByIDs(ctx context.Context, profileID, postID, postProfileID string) (*events.Event, error) {
	return s.request(ctx, "addPostViewCount", map[string]any{
		"profile_id":      s.profileOr(profileID),
		"post_id":         postID,
		"post_profile_id": postProfileID,
	})
}

// GetProfilesStoriesByIDs returns stories for the given profile.
// (GetProfilesStories already exists; this is an alias with limit.)
func (s *Rubino) GetProfilesStoriesByIDs(ctx context.Context, profileID string, limit int) (*events.Event, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.request(ctx, "getProfilesStories", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
	})
}

// RequestFollow performs a follow/unfollow action. Follow and
// Unfollow wrap it; this exposes the underlying requestFollow endpoint
// directly.
func (s *Rubino) RequestFollow(ctx context.Context, profileID, followeeID, fType string) (*events.Event, error) {
	return s.request(ctx, "requestFollow", map[string]any{
		"profile_id":  s.profileOr(profileID),
		"followee_id": followeeID,
		"f_type":      fType,
	})
}

// GetProfileFollowersByType returns followers or followings based on
// fType ("Follower" or "Following"). (GetProfileFollowers/Followings
// already exist; this exposes the underlying endpoint directly.)
func (s *Rubino) GetProfileFollowersByType(ctx context.Context, profileID, targetProfileID, fType string, limit int, sort, equal string) (*events.Event, error) {
	if limit <= 0 {
		limit = 50
	}
	if sort == "" {
		sort = "FromMax"
	}
	return s.request(ctx, "getProfileFollowers", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"f_type":            fType,
		"limit":             limit,
		"sort":              sort,
		"equal":             equal != "",
	})
}
