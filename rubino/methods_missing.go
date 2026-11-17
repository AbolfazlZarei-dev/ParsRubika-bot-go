package rubino

import (
	"context"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// GetShareLink returns the share link for a Rubino post.
func (s *Rubino) GetShareLink(ctx context.Context, profileID, postID, postProfileID string) (*events.Event, error) {
	return s.request(ctx, "getShareLink", map[string]any{
		"profile_id":      s.profileOr(profileID),
		"post_id":         postID,
		"post_profile_id": postProfileID,
	})
}

// UploadProfilePhoto uploads a new profile photo.
func (s *Rubino) UploadProfilePhoto(ctx context.Context, profileID, fileID string) (*events.Event, error) {
	return s.request(ctx, "updateProfilePhoto", map[string]any{
		"profile_id": s.profileOr(profileID),
		"file_id":    fileID,
	})
}
