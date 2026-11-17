package rubino

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// This file extends the Rubino client with the profile, admin, search
// and feed methods: profile settings & admin management, close friends,
// the search/explore/feed group, post update/share/collaborator/feedback,
// follower removal, the metrics/analytics group, message-owner actions,
// and a local profile_created_at decoder.
//
// Every method mirrors the exact Rubino `method` string and input shape
// used by MaxRubika so the two libraries are behaviourally equivalent.

// ===== Profile settings & admins ======================================

// GetProfileSetting returns the active profile's settings.
func (s *Rubino) GetProfileSetting(ctx context.Context, profileID string) (*events.Event, error) {
	return s.request(ctx, "getProfileSetting", map[string]any{"profile_id": s.profileOr(profileID)})
}

// SetProfileSetting updates the active profile's settings. `settings` is
// the raw settings map (key/value pairs the Rubino API accepts).
func (s *Rubino) SetProfileSetting(ctx context.Context, profileID string, settings map[string]any) (*events.Event, error) {
	in := map[string]any{"profile_id": s.profileOr(profileID)}
	for k, v := range settings {
		in[k] = v
	}
	return s.request(ctx, "setProfileSetting", in)
}

// GetSettingAdminPermissions returns the available admin permission set.
func (s *Rubino) GetSettingAdminPermissions(ctx context.Context) (*events.Event, error) {
	return s.request(ctx, "getSettingAdminPermissions", map[string]any{})
}

// GetProfileAdmins lists the admins of a profile.
func (s *Rubino) GetProfileAdmins(ctx context.Context, profileID string) (*events.Event, error) {
	return s.request(ctx, "getProfileAdmins", map[string]any{"profile_id": s.profileOr(profileID)})
}

// SetProfileAdmin promotes a profile admin with the given access list.
// `accessList` is the set of permission names (e.g. ["EditProfile",
// "AddPost"]); `action` is typically "SetAdmin".
func (s *Rubino) SetProfileAdmin(ctx context.Context, profileID, targetProfileID, action string, accessList []string) (*events.Event, error) {
	if action == "" {
		action = "SetAdmin"
	}
	return s.request(ctx, "setProfileAdmin", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"access_list":       accessList,
		"action":            action,
	})
}

// UnsetProfileAdmin demotes a profile admin.
func (s *Rubino) UnsetProfileAdmin(ctx context.Context, profileID, targetProfileID string) (*events.Event, error) {
	return s.request(ctx, "setProfileAdmin", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"access_list":       []string{},
		"action":            "UnsetAdmin",
	})
}

// AcceptProfileAdmin accepts an incoming admin invitation for a profile.
func (s *Rubino) AcceptProfileAdmin(ctx context.Context, profileID, targetProfileID string) (*events.Event, error) {
	return s.request(ctx, "profileAdminAction", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"action":            "Accept",
	})
}

// DenyProfileAdmin denies an incoming admin invitation.
func (s *Rubino) DenyProfileAdmin(ctx context.Context, profileID, targetProfileID string) (*events.Event, error) {
	return s.request(ctx, "profileAdminAction", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"action":            "Deny",
	})
}

// RemoveProfileAdmin removes a profile admin (used by the profile owner).
func (s *Rubino) RemoveProfileAdmin(ctx context.Context, profileID, targetProfileID string) (*events.Event, error) {
	return s.request(ctx, "profileAdminAction", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"action":            "Remove",
	})
}

// ===== Close friends ==================================================

// GetMyCloseFriends returns the caller's close-friends list.
func (s *Rubino) GetMyCloseFriends(ctx context.Context) (*events.Event, error) {
	return s.request(ctx, "getMyCloseFriends", map[string]any{})
}

// SetCloseFriends sets the close-friends list for a profile.
func (s *Rubino) SetCloseFriends(ctx context.Context, profileID string, targetProfileIDs []string) (*events.Event, error) {
	return s.request(ctx, "closeFriendAction", map[string]any{
		"profile_id":         s.profileOr(profileID),
		"target_profile_ids": targetProfileIDs,
		"action":             "Set",
	})
}

// UnsetCloseFriends removes the given profiles from the close-friends list.
func (s *Rubino) UnsetCloseFriends(ctx context.Context, profileID string, targetProfileIDs []string) (*events.Event, error) {
	return s.request(ctx, "closeFriendAction", map[string]any{
		"profile_id":         s.profileOr(profileID),
		"target_profile_ids": targetProfileIDs,
		"action":             "Unset",
	})
}

// ===== Search =========================================================

// SearchPost searches posts by text. `usage` is "Search" by default;
// "Hashtag" switches to hashtag search semantics.
func (s *Rubino) SearchPost(ctx context.Context, searchText, usage, startID string) (*events.Event, error) {
	if usage == "" {
		usage = "Search"
	}
	in := map[string]any{"search_text": searchText, "usage": usage}
	if startID != "" {
		in["start_id"] = startID
	}
	return s.request(ctx, "searchPost", in)
}

// SearchHashTag searches posts by hashtag.
func (s *Rubino) SearchHashTag(ctx context.Context, searchText, startID string) (*events.Event, error) {
	in := map[string]any{"search_text": searchText}
	if startID != "" {
		in["start_id"] = startID
	}
	return s.request(ctx, "searchHashTag", in)
}

// SearchLocation searches posts near a geographic point.
func (s *Rubino) SearchLocation(ctx context.Context, searchText string, lat, lng float64, startID string) (*events.Event, error) {
	in := map[string]any{
		"search_text": searchText,
		"lat":         lat,
		"lng":         lng,
	}
	if startID != "" {
		in["start_id"] = startID
	}
	return s.request(ctx, "searchLocation", in)
}

// ===== Posts: update / share / collaborator / feedback ================

// GetPostsByIDs fetches a batch of posts by ID.
func (s *Rubino) GetPostsByIDs(ctx context.Context, profileID string, postIDs []string) (*events.Event, error) {
	return s.request(ctx, "getPostsByIds", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_ids":   postIDs,
	})
}

// UpdatePost edits a post's caption / settings. `updateParameters` lists
// the fields being changed (e.g. ["caption"]).
func (s *Rubino) UpdatePost(ctx context.Context, profileID, postID string, fields map[string]any, updateParameters []string) (*events.Event, error) {
	in := map[string]any{
		"profile_id":         s.profileOr(profileID),
		"post_id":            postID,
		"updated_parameters": updateParameters,
	}
	for k, v := range fields {
		in[k] = v
	}
	return s.request(ctx, "updatePost", in)
}

// UpdatePostShareCount records a share event for a post.
func (s *Rubino) UpdatePostShareCount(ctx context.Context, profileID, postID string) (*events.Event, error) {
	return s.request(ctx, "updatePostShareCount", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_id":    postID,
	})
}

// PostCollaboratorAction adds/removes a collaborator on a post.
// `action` is "AddCollaborator" or "RemoveCollaborator".
func (s *Rubino) PostCollaboratorAction(ctx context.Context, profileID, postID, targetProfileID, action string) (*events.Event, error) {
	if action == "" {
		action = "AddCollaborator"
	}
	return s.request(ctx, "postCollaboratorAction", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"post_id":           postID,
		"target_profile_id": targetProfileID,
		"action":            action,
	})
}

// SetPostFeedback submits feedback (e.g. not-interested) on a post.
func (s *Rubino) SetPostFeedback(ctx context.Context, profileID, postID, feedback string) (*events.Event, error) {
	return s.request(ctx, "setPostFeedback", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_id":    postID,
		"feedback":   feedback,
	})
}

// ===== Feed ===========================================================

// GetFeedData returns the caller's feed.
func (s *Rubino) GetFeedData(ctx context.Context, profileID, startID string, limit int) (*events.Event, error) {
	in := map[string]any{"profile_id": s.profileOr(profileID)}
	if startID != "" {
		in["start_id"] = startID
	}
	if limit > 0 {
		in["limit"] = limit
	}
	return s.request(ctx, "getFeedData", in)
}

// ViewFeedData records a feed-view event (used for analytics).
func (s *Rubino) ViewFeedData(ctx context.Context, profileID string, feedRowIDs []string) (*events.Event, error) {
	return s.request(ctx, "viewFeedData", map[string]any{
		"profile_id":   s.profileOr(profileID),
		"feed_row_ids": feedRowIDs,
	})
}

// ===== Followers / feedback ==========================================

// RemoveFollower removes a follower from a profile.
func (s *Rubino) RemoveFollower(ctx context.Context, profileID, targetProfileID string) (*events.Event, error) {
	return s.request(ctx, "profileAction", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"action":            "RemoveFollower",
	})
}

// NotInterestedProfile marks a profile as "not interested".
func (s *Rubino) NotInterestedProfile(ctx context.Context, profileID, targetProfileID string) (*events.Event, error) {
	return s.request(ctx, "setProfileFeedback", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"feedback":          "NotInterested",
	})
}

// ===== Metrics / analytics ===========================================

// GetBoostURL returns the profile's boost (promotion) URL.
func (s *Rubino) GetBoostURL(ctx context.Context, profileID string) (*events.Event, error) {
	return s.request(ctx, "getBoostUrl", map[string]any{"profile_id": s.profileOr(profileID)})
}

// GetMetricsConfig returns the metrics configuration.
func (s *Rubino) GetMetricsConfig(ctx context.Context) (*events.Event, error) {
	return s.request(ctx, "getMetricsConfig", map[string]any{})
}

// TrackClientMetrics submits a client-metrics payload.
func (s *Rubino) TrackClientMetrics(ctx context.Context, payload map[string]any) (*events.Event, error) {
	return s.request(ctx, "trackClientMetrics", payload)
}

// AddImpressionIDs records impression IDs for analytics.
func (s *Rubino) AddImpressionIDs(ctx context.Context, profileID string, impressionIDs []string) (*events.Event, error) {
	return s.request(ctx, "addImpressionIds", map[string]any{
		"profile_id":     s.profileOr(profileID),
		"impression_ids": impressionIDs,
	})
}

// AddUsageTime records usage time for a profile.
func (s *Rubino) AddUsageTime(ctx context.Context, profileID string, seconds int) (*events.Event, error) {
	return s.request(ctx, "addUsageTime", map[string]any{
		"profile_id": s.profileOr(profileID),
		"usage_time": seconds,
	})
}

// ===== Message owner =================================================

// AcceptMessageOwner accepts a message-owner request for a post.
func (s *Rubino) AcceptMessageOwner(ctx context.Context, profileID, postID string) (*events.Event, error) {
	return s.request(ctx, "settingMessageOwnerAction", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_id":    postID,
		"action":     "Accept",
	})
}

// DenyMessageOwner denies a message-owner request.
func (s *Rubino) DenyMessageOwner(ctx context.Context, profileID, postID string) (*events.Event, error) {
	return s.request(ctx, "settingMessageOwnerAction", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_id":    postID,
		"action":     "Deny",
	})
}

// RemoveMessageOwner removes the message owner of a post.
func (s *Rubino) RemoveMessageOwner(ctx context.Context, profileID, postID string) (*events.Event, error) {
	return s.request(ctx, "settingMessageOwnerAction", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_id":    postID,
		"action":     "Remove",
	})
}

// ===== Local helpers (no API call) ===================================

// ProfileCreatedAt decodes the creation timestamp encoded in the first 8
// hex characters of a Rubino profile_id. Rubino profile IDs are 24-hex
// strings whose first 4 bytes are a big-endian Unix timestamp. It makes
// no API call.
//
// Returns an error when the ID is too short or not valid hex.
func ProfileCreatedAt(profileID string) (time.Time, error) {
	id := strings.TrimSpace(profileID)
	if len(id) < 8 {
		return time.Time{}, coreerrors.New(coreerrors.KindValidation,
			"rubino: profile_id too short (need ≥8 hex chars)")
	}
	b, err := hex.DecodeString(id[:8])
	if err != nil {
		return time.Time{}, coreerrors.Wrap(coreerrors.KindValidation,
			"rubino: profile_id hex decode", err)
	}
	if len(b) < 4 {
		return time.Time{}, coreerrors.New(coreerrors.KindValidation,
			"rubino: profile_id hex too short")
	}
	ts := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return time.Unix(int64(ts), 0).UTC(), nil
}

// ProfileCreatedAtString is a convenience wrapper returning the RFC-3339
// creation timestamp of a profile_id.
func ProfileCreatedAtString(profileID string) (string, error) {
	t, err := ProfileCreatedAt(profileID)
	if err != nil {
		return "", err
	}
	return t.Format(time.RFC3339), nil
}

// formatHexID renders a byte slice as a lowercase hex string, used by
// helpers that render profile IDs.
func formatHexID(b []byte) string { return hex.EncodeToString(b) }

// strconvInt is retained so future input validators can use it without
// re-importing strconv in this file.
func strconvInt(s string) (int, error) { return strconv.Atoi(s) }

// fmtErr is a tiny wrapper to keep error construction concise.
func fmtErr(format string, a ...any) error { return fmt.Errorf(format, a...) }
