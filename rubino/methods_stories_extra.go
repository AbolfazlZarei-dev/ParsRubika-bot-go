package rubino

import (
	"context"
	"fmt"
	"os"
)

// AddStory uploads and publishes a story. `fileType` is "Picture" or
// "Video". The file is uploaded via RequestUploadFile + UploadFile.
func (s *Rubino) AddStory(ctx context.Context, profileID, filePath, caption, fileType string) (*Event, error) {
	fileName := baseName(filePath)
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("rubino: stat %s: %w", filePath, err)
	}
	res, err := s.RequestUploadFile(ctx, profileID, fileName, stat.Size(), fileType)
	if err != nil {
		return nil, err
	}
	uploadURL := res.Str("upload_url")
	if uploadURL == "" {
		return nil, fmt.Errorf("rubino: no upload_url in response")
	}
	fileID, err := s.UploadFile(ctx, uploadURL, fileName, filePath)
	if err != nil {
		return nil, err
	}
	in := map[string]any{
		"profile_id": s.profileOr(profileID),
		"file_id":    fileID,
		"rnd":        rnd(),
	}
	if caption != "" {
		in["caption"] = caption
	}
	return s.request(ctx, "addStory", in)
}
