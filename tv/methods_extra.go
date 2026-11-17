package tv

import (
	"context"
	"strings"

	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// This file adds the playback, discovery and engagement TV methods:
// getMediaByShareString, getNextEpisode,
// getNextMediaInfo, getPlayInfo, getPlaySuggestedEpisode,
// setCurrentPlayTime, setSessionPlayDuration, viewCommercial, and the
// track_id-aware Dislike / Undo / Follow / UnFollow media actions.
//
// Each method mirrors the exact Rubika `method` string and input shape
// used by MaxRubika so the two libraries are behaviourally equivalent
// (excluding the Shad platform, which is out of scope).

// ===== Media lookup ====================================================

// GetMediaByLink retrieves media details by its share link (e.g.
// "https://rubika.ir/vod/294119"). The numeric ID is extracted from the
// last path segment; non-numeric links are rejected.
func (t *TV) GetMediaByLink(ctx context.Context, link, trackID string) (*events.Event, error) {
	// Strip trailing slashes and take the last path segment.
	stripped := strings.TrimRight(link, "/")
	if i := strings.LastIndex(stripped, "/"); i >= 0 {
		stripped = stripped[i+1:]
	}
	if stripped == "" || !isAllDigits(stripped) {
		return nil, coreerrors.New(coreerrors.KindValidation,
			"tv: could not extract a numeric media ID from "+link)
	}
	in := map[string]any{"share_string": stripped}
	if trackID != "" {
		in["track_id"] = trackID
	}
	return t.request(ctx, "getMediaByShareString", in)
}

// GetNextEpisode retrieves the next episode after the given episode.
func (t *TV) GetNextEpisode(ctx context.Context, episodeID string) (*events.Event, error) {
	return t.request(ctx, "getNextEpisode", map[string]any{"episode_id": episodeID})
}

// GetNextMediaInfo retrieves the next media info after the given media.
func (t *TV) GetNextMediaInfo(ctx context.Context, mediaID, trackID string) (*events.Event, error) {
	in := map[string]any{"media_id": mediaID}
	if trackID != "" {
		in["track_id"] = trackID
	}
	return t.request(ctx, "getNextMediaInfo", in)
}

// PlayType selects whether GetPlayInfo returns the main stream or the
// trailer.
type PlayType string

const (
	PlayTypeMain    PlayType = "Main"
	PlayTypeTrailer PlayType = "Trailer"
)

// normalizePlayType maps loose input to Main/Trailer (case-insensitive).
func normalizePlayType(in string) PlayType {
	switch strings.ToLower(strings.TrimSpace(in)) {
	case "main":
		return PlayTypeMain
	case "trailer":
		return PlayTypeTrailer
	}
	return ""
}

// GetPlayInfo retrieves playback information for a media. `playType`
// accepts "Main" or "Trailer" (case-insensitive, defaults to Trailer);
// `version` defaults to 1.
func (t *TV) GetPlayInfo(ctx context.Context, mediaID string, playType PlayType, version int, trackID string) (*events.Event, error) {
	pt := normalizePlayType(string(playType))
	if pt == "" {
		pt = PlayTypeTrailer
	}
	if version <= 0 {
		version = 1
	}
	in := map[string]any{
		"media_id": mediaID,
		"type":     string(pt),
		"version":  version,
	}
	if trackID != "" {
		in["track_id"] = trackID
	}
	return t.request(ctx, "getPlayInfo", in)
}

// GetPlaySuggestedEpisode retrieves the suggested next episode for a
// series.
func (t *TV) GetPlaySuggestedEpisode(ctx context.Context, seriesID, trackID string) (*events.Event, error) {
	in := map[string]any{"series_id": seriesID}
	if trackID != "" {
		in["track_id"] = trackID
	}
	return t.request(ctx, "getPlaySuggestedEpisode", in)
}

// ===== Playback tracking ===============================================

// SetCurrentPlayTime records the user's current play position (in
// seconds) for a media. Negative values are rejected.
func (t *TV) SetCurrentPlayTime(ctx context.Context, mediaID string, playTime int) (*events.Event, error) {
	if playTime < 0 {
		return nil, coreerrors.New(coreerrors.KindValidation,
			"tv: play_time must be >= 0")
	}
	return t.request(ctx, "setCurrentPlayTime", map[string]any{
		"media_id":  mediaID,
		"play_time": playTime,
	})
}

// SetSessionPlayDuration records the total play duration (in seconds)
// for a session.
func (t *TV) SetSessionPlayDuration(ctx context.Context, sessionID string, playDuration int) (*events.Event, error) {
	return t.request(ctx, "setSessionPlayDuration", map[string]any{
		"session_id":    sessionID,
		"play_duration": playDuration,
	})
}

// ViewType is the kind of commercial view event registered by
// ViewCommercial.
type ViewType string

const (
	ViewTypeFull ViewType = "Full"
	ViewTypeHalf ViewType = "Half"
	ViewTypeSkip ViewType = "Skip"
)

func normalizeViewType(in string) ViewType {
	switch strings.ToLower(strings.TrimSpace(in)) {
	case "full":
		return ViewTypeFull
	case "half":
		return ViewTypeHalf
	case "skip":
		return ViewTypeSkip
	}
	return ""
}

// ViewCommercial registers a commercial view event. `viewType` accepts
// Full/Half/Skip (case-insensitive, defaults to Full).
func (t *TV) ViewCommercial(ctx context.Context, viewType ViewType) (*events.Event, error) {
	vt := normalizeViewType(string(viewType))
	if vt == "" {
		vt = ViewTypeFull
	}
	return t.request(ctx, "viewCommercial", map[string]any{"view_type": string(vt)})
}

// ===== track_id-aware media actions ====================================

// actionOnMediaTrack is the track_id-aware variant of actionOnMedia.
// The track_id is the session hash string the TV web client uses for
// analytics; it is optional.
func (t *TV) actionOnMediaTrack(ctx context.Context, action, mediaID, trackID string) (*events.Event, error) {
	in := map[string]any{"action": action, "media_id": mediaID}
	if trackID != "" {
		in["track_id"] = trackID
	}
	return t.request(ctx, "actionOnMedia", in)
}

// DislikeMedia adds a dislike to a media.
func (t *TV) DislikeMedia(ctx context.Context, mediaID, trackID string) (*events.Event, error) {
	return t.actionOnMediaTrack(ctx, "Dislike", mediaID, trackID)
}

// UndoLikeOrDislikeMedia undoes a previous like or dislike.
func (t *TV) UndoLikeOrDislikeMedia(ctx context.Context, mediaID, trackID string) (*events.Event, error) {
	return t.actionOnMediaTrack(ctx, "UndoLikeOrDislike", mediaID, trackID)
}

// FollowMedia follows a media (e.g. a series).
func (t *TV) FollowMedia(ctx context.Context, mediaID, trackID string) (*events.Event, error) {
	return t.actionOnMediaTrack(ctx, "Follow", mediaID, trackID)
}

// UnfollowMedia unfollows a media.
func (t *TV) UnfollowMedia(ctx context.Context, mediaID, trackID string) (*events.Event, error) {
	return t.actionOnMediaTrack(ctx, "UnFollow", mediaID, trackID)
}

// isAllDigits reports whether s is non-empty and all ASCII digits.
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
