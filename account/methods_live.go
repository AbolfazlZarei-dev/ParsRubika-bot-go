package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Live streams and live location.

// SendLive starts a live stream in a chat.
func (a *Account) SendLive(ctx context.Context, objectGUID, title string, thumbInline any) (*events.Event, error) {
	in := map[string]any{
		"object_guid": objectGUID,
		"device_type": "Software",
		"rnd":         randInt63(),
	}
	if title != "" {
		in["title"] = title
	}
	if thumbInline != nil {
		in["thumb_inline"] = thumbInline
	}
	return a.Request(ctx, RequestOptions{Method: "sendLive", Input: in, Encrypt: true})
}

// GetLiveStatus returns the status of a live stream.
func (a *Account) GetLiveStatus(ctx context.Context, liveID, accessToken string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getLiveStatus", Input: map[string]any{
		"live_id":      liveID,
		"access_token": accessToken,
	}, Encrypt: true})
}

// GetLivePlayURL returns the playback URL of a live stream.
func (a *Account) GetLivePlayURL(ctx context.Context, accessToken, liveID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getLivePlayUrl", Input: map[string]any{
		"access_token": accessToken,
		"live_id":      liveID,
	}, Encrypt: true})
}

// GetLiveComments returns the comments of a live stream.
func (a *Account) GetLiveComments(ctx context.Context, liveID, accessToken string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getLiveComments", Input: map[string]any{
		"live_id":      liveID,
		"access_token": accessToken,
	}, Encrypt: true})
}

// AddLiveComment adds a comment to a live stream.
func (a *Account) AddLiveComment(ctx context.Context, accessToken, liveID, text string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "addLiveComment", Input: map[string]any{
		"access_token": accessToken,
		"live_id":      liveID,
		"text":         text,
	}, Encrypt: true})
}

// SetCurrentLiveLocation shares a live location with map tiles.
func (a *Account) SetCurrentLiveLocation(ctx context.Context, objectGUID string, lat, lon float64, tileSideCount int, tileURLs []string, xLoc, yLoc float64) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setCurrentLiveLocation", Input: map[string]any{
		"object_guid": objectGUID,
		"location": map[string]any{
			"latitude":  lat,
			"longitude": lon,
			"map_view": map[string]any{
				"tile_side_count": tileSideCount,
				"tile_urls":       tileURLs,
				"x_loc":           xLoc,
				"y_loc":           yLoc,
			},
		},
	}, Encrypt: true})
}

// GetMapView returns a map view for the given coordinates.
func (a *Account) GetMapView(ctx context.Context, lat, lon float64) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getMapView", Input: map[string]any{
		"location": map[string]any{"latitude": lat, "longitude": lon},
	}, Encrypt: true})
}
