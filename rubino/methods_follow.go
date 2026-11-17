package rubino

import "context"

// ===== Follow / unfollow ============================================

// Follow requests to follow a profile.
func (s *Rubino) Follow(ctx context.Context, profileID, followeeID string) (*Event, error) {
	return s.request(ctx, "requestFollow", map[string]any{
		"profile_id":  s.profileOr(profileID),
		"followee_id": followeeID,
		"f_type":      "Follow",
	})
}

// Unfollow requests to unfollow a profile.
func (s *Rubino) Unfollow(ctx context.Context, profileID, followeeID string) (*Event, error) {
	return s.request(ctx, "requestFollow", map[string]any{
		"profile_id":  s.profileOr(profileID),
		"followee_id": followeeID,
		"f_type":      "Unfollow",
	})
}

// AcceptFollowRequest accepts a pending follow request.
func (s *Rubino) AcceptFollowRequest(ctx context.Context, profileID, requestID string) (*Event, error) {
	return s.request(ctx, "actionOnRequest", map[string]any{
		"profile_id": s.profileOr(profileID),
		"request_id": requestID,
		"action":     "Accept",
	})
}

// DeclineFollowRequest declines a pending follow request.
func (s *Rubino) DeclineFollowRequest(ctx context.Context, profileID, requestID string) (*Event, error) {
	return s.request(ctx, "actionOnRequest", map[string]any{
		"profile_id": s.profileOr(profileID),
		"request_id": requestID,
		"action":     "Decline",
	})
}

// GetNewFollowRequests returns the pending follow requests.
func (s *Rubino) GetNewFollowRequests(ctx context.Context, profileID string, limit int, sort string) (*Event, error) {
	return s.request(ctx, "getNewFollowRequests", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
	})
}

// GetProfileFollowers returns the followers of a target profile.
func (s *Rubino) GetProfileFollowers(ctx context.Context, profileID, targetProfileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getProfileFollowers", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"limit":             limit,
		"sort":              ifEmpty(sort, "latest"),
		"equal":             equal,
		"f_type":            "Follower",
	})
}

// GetProfileFollowings returns the followings of a target profile.
func (s *Rubino) GetProfileFollowings(ctx context.Context, profileID, targetProfileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getProfileFollowers", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"limit":             limit,
		"sort":              ifEmpty(sort, "latest"),
		"equal":             equal,
		"f_type":            "Following",
	})
}

// SearchFollower searches the followers of a target profile.
func (s *Rubino) SearchFollower(ctx context.Context, profileID, targetProfileID, username string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getProfileFollowers", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"username":          username,
		"limit":             limit,
		"sort":              ifEmpty(sort, "latest"),
		"equal":             equal,
		"f_type":            "Follower",
		"search_type":       "Follower",
	})
}

// SearchFollowing searches the followings of a target profile.
func (s *Rubino) SearchFollowing(ctx context.Context, profileID, targetProfileID, username string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getProfileFollowers", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"username":          username,
		"limit":             limit,
		"sort":              ifEmpty(sort, "latest"),
		"equal":             equal,
		"f_type":            "Following",
		"search_type":       "Following",
	})
}
