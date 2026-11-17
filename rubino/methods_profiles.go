package rubino

import "context"

// GetProfileList returns the user's pages/profiles.
func (s *Rubino) GetProfileList(ctx context.Context, limit int, sort, equal string) (*Event, error) {
	in := map[string]any{
		"limit": limit,
		"sort":  ifEmpty(sort, "latest"),
		"equal": equal,
	}
	return s.request(ctx, "getProfileList", in)
}

// GetMyProfileInfo returns the current user's full profile.
func (s *Rubino) GetMyProfileInfo(ctx context.Context, profileID string) (*Event, error) {
	return s.request(ctx, "getMyProfileInfo", map[string]any{
		"profile_id": s.profileOr(profileID),
	})
}

// CreatePage creates a new Rubino page.
func (s *Rubino) CreatePage(ctx context.Context, username, name, bio, phone, email, website string) (*Event, error) {
	return s.request(ctx, "createPage", map[string]any{
		"username": username,
		"name":     name,
		"bio":      bio,
		"phone":    phone,
		"email":    email,
		"website":  website,
	})
}

// UpdateProfile updates the current user's profile metadata. Pass
// empty strings for fields you don't want to change.
func (s *Rubino) UpdateProfile(ctx context.Context, profileID, username, name, bio, phone, email, website, profileStatus string, isMessageAllowed *bool, isMute *bool) (*Event, error) {
	in := map[string]any{
		"profile_id": s.profileOr(profileID),
	}
	if username != "" {
		in["username"] = username
	}
	if name != "" {
		in["name"] = name
	}
	if bio != "" {
		in["bio"] = bio
	}
	if phone != "" {
		in["phone"] = phone
	}
	if email != "" {
		in["email"] = email
	}
	if website != "" {
		in["website"] = website
	}
	if profileStatus != "" {
		in["profile_status"] = profileStatus
	}
	if isMessageAllowed != nil {
		in["is_message_allowed"] = *isMessageAllowed
	}
	if isMute != nil {
		in["is_mute"] = *isMute
	}
	return s.request(ctx, "updateProfile", in)
}

// IsExistUsername checks whether a username is available.
func (s *Rubino) IsExistUsername(ctx context.Context, username string) (*Event, error) {
	return s.request(ctx, "isExistUsername", map[string]any{"username": username})
}

// GetProfileInfo returns info about a target profile.
func (s *Rubino) GetProfileInfo(ctx context.Context, profileID, targetProfileID string) (*Event, error) {
	return s.request(ctx, "getProfileInfo", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
	})
}

// BlockProfile blocks a profile.
func (s *Rubino) BlockProfile(ctx context.Context, profileID, blockedID string) (*Event, error) {
	return s.request(ctx, "blockProfile", map[string]any{
		"profile_id":  s.profileOr(profileID),
		"blocked_id":  blockedID,
		"action_type": "Block",
	})
}

// UnblockProfile unblocks a profile.
func (s *Rubino) UnblockProfile(ctx context.Context, profileID, blockedID string) (*Event, error) {
	return s.request(ctx, "blockProfile", map[string]any{
		"profile_id":  s.profileOr(profileID),
		"blocked_id":  blockedID,
		"action_type": "Unblock",
	})
}

// RemovePage removes a page. The removeRecord API expects a
// `record_id` + `model` pair (model="Profile").
func (s *Rubino) RemovePage(ctx context.Context, profileID, recordID string) (*Event, error) {
	return s.request(ctx, "removeRecord", map[string]any{
		"profile_id": s.profileOr(profileID),
		"record_id":  recordID,
		"model":      "Profile",
	})
}

// SetProfilePublic sets a profile public.
func (s *Rubino) SetProfilePublic(ctx context.Context, profileID string) (*Event, error) {
	return s.request(ctx, "updateProfile", map[string]any{
		"profile_id":     s.profileOr(profileID),
		"profile_status": "Public",
	})
}

// SetProfilePrivate sets a profile private.
func (s *Rubino) SetProfilePrivate(ctx context.Context, profileID string) (*Event, error) {
	return s.request(ctx, "updateProfile", map[string]any{
		"profile_id":     s.profileOr(profileID),
		"profile_status": "Private",
	})
}

// EnablePrivateMessages allows others to message the profile.
func (s *Rubino) EnablePrivateMessages(ctx context.Context, profileID string) (*Event, error) {
	return s.request(ctx, "updateProfile", map[string]any{
		"profile_id":         s.profileOr(profileID),
		"is_message_allowed": true,
	})
}

// DisablePrivateMessages disallows others from messaging the profile.
func (s *Rubino) DisablePrivateMessages(ctx context.Context, profileID string) (*Event, error) {
	return s.request(ctx, "updateProfile", map[string]any{
		"profile_id":         s.profileOr(profileID),
		"is_message_allowed": false,
	})
}

func ifEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
