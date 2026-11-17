package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Avatar management (list, upload, set main, delete).

// GetAvatars returns the avatar list of an object.
func (a *Account) GetAvatars(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getAvatars", Input: map[string]any{"object_guid": objectGUID}, Encrypt: true})
}

// UploadAvatar uploads a new avatar for the supplied object. mainFileID
// is the full-size photo and thumbnailFileID the small one; passing the
// same id for both reuses one upload as its own thumbnail. (Wire format
// verified against maxrubika's working upload_avatar: the server wants
// object_guid + main_file_id + thumbnail_file_id; file_id/dc_id/
// access_hash_rec are rejected with INVALID_INPUT.)
func (a *Account) UploadAvatar(ctx context.Context, objectGUID, mainFileID, thumbnailFileID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method: "uploadAvatar",
		Input: map[string]any{
			"object_guid":       objectGUID,
			"main_file_id":      mainFileID,
			"thumbnail_file_id": thumbnailFileID,
		},
		Encrypt: true,
	})
}

// SetMainAvatar promotes one of the object's existing avatars (by
// avatar_id) to be the current/main one. Wire method: updateAvatar.
func (a *Account) SetMainAvatar(ctx context.Context, objectGUID, avatarID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "updateAvatar",
		Input:   map[string]any{"object_guid": objectGUID, "avatar_id": avatarID},
		Encrypt: true,
	})
}

// DeleteAvatar deletes an avatar.
func (a *Account) DeleteAvatar(ctx context.Context, objectGUID, avatarID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "deleteAvatar",
		Input:   map[string]any{"object_guid": objectGUID, "avatar_id": avatarID},
		Encrypt: true,
	})
}

// DeleteAllAvatars deletes every avatar of an object by looping over
// getAvatars and calling deleteAvatar for each.
func (a *Account) DeleteAllAvatars(ctx context.Context, objectGUID string) error {
	resp, err := a.GetAvatars(ctx, objectGUID)
	if err != nil {
		return err
	}
	avatars := resp.Slice("avatars")
	for _, av := range avatars {
		if avm, ok := av.(map[string]any); ok {
			avatarID, _ := avm["avatar_id"].(string)
			if avatarID == "" {
				continue
			}
			if _, err := a.DeleteAvatar(ctx, objectGUID, avatarID); err != nil {
				return err
			}
		}
	}
	return nil
}
