package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Saved-music playlist.

// GetSavedMusicPlaylist returns the user's saved music playlist.
func (a *Account) GetSavedMusicPlaylist(ctx context.Context, userGUID string) (*events.Event, error) {
	guid := userGUID
	if guid == "" {
		guid = a.guid
	}
	return a.Request(ctx, RequestOptions{Method: "getSavedMusicPlaylist", Input: map[string]any{
		"object_guid": guid,
	}, Encrypt: true})
}

// SetSavedMusicPlaylist replaces the entire saved music playlist.
func (a *Account) SetSavedMusicPlaylist(ctx context.Context, playlistTracks []map[string]any) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setSavedMusicPlaylist", Input: map[string]any{
		"object_guid":     a.guid,
		"playlist_tracks": playlistTracks,
	}, Encrypt: true})
}

// AddSavedMusicTrack adds a single track to the saved music playlist.
func (a *Account) AddSavedMusicTrack(ctx context.Context, track map[string]any) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "addSavedMusicTrack", Input: map[string]any{
		"object_guid": a.guid,
		"added_track": track,
	}, Encrypt: true})
}

// DeleteSavedMusicPlaylist deletes the entire playlist (or a single position).
func (a *Account) DeleteSavedMusicPlaylist(ctx context.Context, position *int) (*events.Event, error) {
	if position == nil {
		return a.SetSavedMusicPlaylist(ctx, []map[string]any{})
	}
	res, err := a.GetSavedMusicPlaylist(ctx, "")
	if err != nil {
		return nil, err
	}
	tracks := res.Slice("playlist_tracks")
	var newTracks []map[string]any
	for i, t := range tracks {
		if i == *position {
			continue
		}
		if m, ok := t.(map[string]any); ok {
			newTracks = append(newTracks, m)
		}
	}
	return a.SetSavedMusicPlaylist(ctx, newTracks)
}
