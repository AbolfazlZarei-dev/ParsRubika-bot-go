package account

import (
	"context"
	"fmt"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/files"
	"path/filepath"
	"sync"
)

// Own profile, privacy, appearance and profile-photo shortcuts.

// GetNotificationSetting returns the user's notification settings.
func (a *Account) GetNotificationSetting(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getNotificationSetting", Input: map[string]any{}, Encrypt: true})
}

// SetNotification sets notification preferences.
func (a *Account) SetNotification(ctx context.Context, settings map[string]any, updateParameters []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setSetting", Input: map[string]any{
		"settings":          settings,
		"update_parameters": updateParameters,
	}, Encrypt: true})
}

// GetWallpapers returns the available wallpapers.
func (a *Account) GetWallpapers(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getWallpapers", Input: map[string]any{}, Encrypt: true})
}

// ResetWallpapers resets the user's wallpaper to default.
func (a *Account) ResetWallpapers(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "resetWallpapers", Input: map[string]any{}, Encrypt: true})
}

// UpdateMyName updates the current user's display name.
func (a *Account) UpdateMyName(ctx context.Context, firstName, lastName string) (*events.Event, error) {
	return a.UpdateProfile(ctx, firstName, lastName, "", "")
}

// UpdateMyBio updates the current user's bio.
func (a *Account) UpdateMyBio(ctx context.Context, bio string) (*events.Event, error) {
	return a.UpdateProfile(ctx, "", "", bio, "")
}

// UpdateMyBirthday updates the current user's birthday (format: YYYY-MM-DD).
func (a *Account) UpdateMyBirthday(ctx context.Context, birthday string) (*events.Event, error) {
	return a.UpdateProfile(ctx, "", "", "", birthday)
}

// EditShowLastOnline sets the privacy for "last online" visibility.
func (a *Account) EditShowLastOnline(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error) {
	return a.SetPrivacySetting(ctx, "show_my_last_online", value, exceptions)
}

// EditShowPhoneNumber sets the privacy for phone number visibility.
func (a *Account) EditShowPhoneNumber(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error) {
	return a.SetPrivacySetting(ctx, "show_my_phone_number", value, exceptions)
}

// EditShowProfilePhoto sets the privacy for profile photo visibility.
func (a *Account) EditShowProfilePhoto(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error) {
	return a.SetPrivacySetting(ctx, "show_my_profile_photo", value, exceptions)
}

// EditShowBirthday sets the privacy for birthday visibility.
func (a *Account) EditShowBirthday(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error) {
	return a.SetPrivacySetting(ctx, "show_my_birth_date", value, exceptions)
}

// EditCanCalledBy sets the privacy for who can call the user.
func (a *Account) EditCanCalledBy(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error) {
	return a.SetPrivacySetting(ctx, "can_called_by", value, exceptions)
}

// EditCanJoinChatBy sets the privacy for who can add the user to chats.
func (a *Account) EditCanJoinChatBy(ctx context.Context, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error) {
	return a.SetPrivacySetting(ctx, "can_join_chat_by", value, exceptions)
}

// UnblockAllUsers unblocks every blocked user. The blocked-users API
// returns members under "in_chat_members" with each user's GUID in the
// "user_guid" field (not "object_guid").
func (a *Account) UnblockAllUsers(ctx context.Context) error {
	resp, err := a.GetBlockedUsers(ctx)
	if err != nil {
		return err
	}
	members := resp.Slice("in_chat_members")
	for _, m := range members {
		if mm, ok := m.(map[string]any); ok {
			// Blocked-user objects use "user_guid" for the user's GUID.
			guid, _ := mm["user_guid"].(string)
			if guid == "" {
				continue
			}
			if _, err := a.UnblockUser(ctx, guid); err != nil {
				return err
			}
		}
	}
	return nil
}

// ProfilePhotoResult holds the outcome of a profile-photo upload.
type ProfilePhotoResult struct {
	AvatarID string
	Avatar   any
}

// UploadProfilePhoto uploads an image file and sets it as the current
// account's profile picture (avatar).
func (a *Account) UploadProfilePhoto(ctx context.Context, path string) (*ProfilePhotoResult, error) {
	if path == "" {
		return nil, fmt.Errorf("account: UploadProfilePhoto: path is empty")
	}

	// 1) Upload the image file to Rubika's file storage.
	mime := mimeFromExt(filepath.Ext(path))
	src := files.Source{Path: path}
	uploaded, err := a.Upload(ctx, src, mime, filepath.Base(path))
	if err != nil {
		return nil, fmt.Errorf("upload photo file: %w", err)
	}

	fileID := uploaded.AnyStr("file_id")
	if fileID == "" {
		return nil, fmt.Errorf("upload photo: server did not return file_id")
	}

	// 2) Bind the uploaded file to the current account as an avatar.
	guid := a.GUID()
	if guid == "" {
		return nil, fmt.Errorf("upload photo: account GUID is empty (not signed in)")
	}

	resp, err := a.UploadAvatar(ctx, guid, fileID, fileID)
	if err != nil {
		return nil, fmt.Errorf("bind avatar: %w", err)
	}

	avatar, _ := resp.Get("avatar")
	result := &ProfilePhotoResult{Avatar: avatar}

	if av, ok := avatar.(map[string]any); ok {
		if id, ok := av["avatar_id"].(string); ok && id != "" {
			result.AvatarID = id
		}
	}
	if result.AvatarID == "" {
		if id := resp.AnyStr("avatar_id"); id != "" {
			result.AvatarID = id
		}
	}
	return result, nil
}

// DeleteProfilePhoto deletes a profile picture (avatar) by its id.
// If avatarID is empty, the most recent avatar is fetched and deleted.
func (a *Account) DeleteProfilePhoto(ctx context.Context, avatarID string) error {
	guid := a.GUID()
	if guid == "" {
		return fmt.Errorf("delete photo: account GUID is empty (not signed in)")
	}

	if avatarID == "" {
		resp, err := a.GetAvatars(ctx, guid)
		if err != nil {
			return fmt.Errorf("fetch avatars: %w", err)
		}
		avatars := resp.Slice("avatars")
		if len(avatars) == 0 {
			return fmt.Errorf("delete photo: no avatars to delete")
		}
		av, _ := avatars[0].(map[string]any)
		id, _ := av["avatar_id"].(string)
		if id == "" {
			return fmt.Errorf("delete photo: avatar has no id")
		}
		avatarID = id
	}

	if _, err := a.DeleteAvatar(ctx, guid, avatarID); err != nil {
		return fmt.Errorf("delete avatar: %w", err)
	}
	return nil
}

// HangProfileResult holds the outcomes of the two concurrent
// updateProfile calls fired by HangProfile.
type HangProfileResult struct {
	First     *events.Event
	FirstErr  error
	Second    *events.Event
	SecondErr error
}

// UpdateProfile updates the current user's profile. Pass empty strings
// for fields you don't want to change. The `updated_parameters` list
// is built automatically from the non-empty fields.
func (a *Account) UpdateProfile(ctx context.Context, firstName, lastName, bio, birthDate string) (*events.Event, error) {
	in := map[string]any{}
	updated := []string{}
	if firstName != "" {
		in["first_name"] = firstName
		updated = append(updated, "first_name")
	}
	if lastName != "" {
		in["last_name"] = lastName
		updated = append(updated, "last_name")
	}
	if bio != "" {
		in["bio"] = bio
		updated = append(updated, "bio")
	}
	if birthDate != "" {
		in["birth_date"] = birthDate
		updated = append(updated, "birth_date")
	}
	in["updated_parameters"] = updated
	return a.Request(ctx, RequestOptions{Method: "updateProfile", Input: in, Encrypt: true})
}

// buildProfileInput mirrors UpdateProfile's input construction so the
// two hang-shots use exactly the same payload shape.
func (a *Account) buildProfileInput(firstName, lastName, bio, birthDate string) map[string]any {
	in := map[string]any{}
	updated := []string{}
	if firstName != "" {
		in["first_name"] = firstName
		updated = append(updated, "first_name")
	}
	if lastName != "" {
		in["last_name"] = lastName
		updated = append(updated, "last_name")
	}
	if bio != "" {
		in["bio"] = bio
		updated = append(updated, "bio")
	}
	if birthDate != "" {
		in["birth_date"] = birthDate
		updated = append(updated, "birth_date")
	}
	in["updated_parameters"] = updated
	return in
}

// HangProfile (a.k.a "profile hang" / "dual profile") fires two
// updateProfile requests back-to-back against the current account so
// the two writes race on the server. This is the well-known Rubika
// trick used to register a "dual"/"hang" profile: when two
// updateProfile calls land near-simultaneously, the server may persist
// an inconsistent state and the profile then renders as two stacked
// identities in clients.
//
// The two shots use the SAME payload (firstName/lastName/bio/birthDate)
// so both writes target the same fields — the "hang" effect comes from
// the race itself, not from conflicting values. If you want the two
// shots to carry different fields, build them with UpdateProfile
// directly.
//
// Both requests are issued concurrently (two goroutines) and the call
// blocks until both return. The context passed to each shot is the
// caller's context; cancelling it aborts whichever shot is still
// in-flight. The method returns the result of the two shots; it does
// not raise on a single shot failing — inspect FirstErr / SecondErr.
func (a *Account) HangProfile(ctx context.Context, firstName, lastName, bio, birthDate string) (*HangProfileResult, error) {
	in := a.buildProfileInput(firstName, lastName, bio, birthDate)

	var wg sync.WaitGroup
	res := &HangProfileResult{}

	fire := func(ev **events.Event, errp *error) {
		defer wg.Done()
		// A fresh map each shot: the transport / encryption layer may
		// mutate the input and we do not want the two concurrent calls
		// to share one backing map.
		shot := make(map[string]any, len(in))
		for k, v := range in {
			shot[k] = v
		}
		e, err := a.Request(ctx, RequestOptions{Method: "updateProfile", Input: shot, Encrypt: true})
		*ev = e
		*errp = err
	}

	wg.Add(2)
	go fire(&res.First, &res.FirstErr)
	go fire(&res.Second, &res.SecondErr)
	wg.Wait()

	return res, nil
}

// UpdateUsername changes the user's username (strips leading `@`).
func (a *Account) UpdateUsername(ctx context.Context, username string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "updateUsername", Input: map[string]any{"username": stripAt(username)}, Encrypt: true})
}

// PrivacyValue is one of "Everybody", "MyContacts", "Nobody".
type PrivacyValue string

const (
	PrivacyEverybody  PrivacyValue = "Everybody"
	PrivacyMyContacts PrivacyValue = "MyContacts"
	PrivacyNobody     PrivacyValue = "Nobody"
)

// PrivacyExceptions is the include/exclude list for a privacy setting.
type PrivacyExceptions struct {
	IncludeUsers []string
	ExcludeUsers []string
}

// SetPrivacySetting sets a single privacy setting. `settingType` is one
// of `ShowLastOnline`, `ShowPhone`, `ShowPhoto`, `CanJoinChat`,
// `ShowBirthDate`, `CanCalled`. The `exceptions` parameter is optional.
func (a *Account) SetPrivacySetting(ctx context.Context, settingType string, value PrivacyValue, exceptions *PrivacyExceptions) (*events.Event, error) {
	settings := map[string]any{settingType: string(value)}
	updateParameters := []string{settingType}
	if exceptions != nil {
		settings[settingType+"_exceptions"] = map[string]any{
			"setting_type":  settingType,
			"include_users": exceptions.IncludeUsers,
			"exclude_users": exceptions.ExcludeUsers,
		}
	}
	return a.Request(ctx, RequestOptions{
		Method: "setSetting",
		Input: map[string]any{
			"settings":          settings,
			"update_parameters": updateParameters,
		},
		Encrypt: true,
	})
}

// SetSetting is a low-level escape hatch that lets the caller build
// the `settings` dict manually.
func (a *Account) SetSetting(ctx context.Context, settings map[string]any, updateParameters []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method: "setSetting",
		Input: map[string]any{
			"settings":          settings,
			"update_parameters": updateParameters,
		},
		Encrypt: true,
	})
}

// GetPrivacySetting returns the user's privacy settings.
func (a *Account) GetPrivacySetting(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getPrivacySetting", Input: map[string]any{}, Encrypt: true})
}

// GetBlockedUsers returns the list of blocked users; pass startID to
// continue a previous page.
func (a *Account) GetBlockedUsers(ctx context.Context, startID ...string) (*events.Event, error) {
	in := map[string]any{}
	if len(startID) > 0 && startID[0] != "" {
		in["start_id"] = startID[0]
	}
	return a.Request(ctx, RequestOptions{Method: "getBlockedUsers", Input: in, Encrypt: true})
}
