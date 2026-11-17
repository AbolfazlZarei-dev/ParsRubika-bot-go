// Package tv implements the ParsRubika TV client. It talks to the
// Rubika VOD service at `https://rbvod{1-4}.iranlms.ir/` with
// `api_version: "1"`. The client sends the TV-specific headers the
// service expects.
//
// The TV API supports browsing media (movies, series, clips),
// searching, liking / disliking, managing a wish-list, viewing cast
// details, and fetching related / seasonal content.
package tv

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport"
)

// TVServerURLs is the list of TV API servers. A random one is picked
// per request to spread the load.
var TVServerURLs = []string{
	"https://rbvod1.iranlms.ir/",
	"https://rbvod2.iranlms.ir/",
	"https://rbvod3.iranlms.ir/",
	"https://rbvod4.iranlms.ir/",
}

// MetaURL is the URL of the version-meta JSON.
const MetaURL = "https://tv.rubika.ir/meta.json"

// DefaultVersion is the fallback app version when MetaURL is
// unreachable.
const DefaultVersion = "3.2.3"

// TVHeaders are the TV-specific HTTP headers (origin / referer set to
// tv.rubika.ir).
var TVHeaders = map[string]string{
	"accept":          "*/*",
	"accept-encoding": "gzip, deflate, br, zstd",
	"accept-language": "fa",
	"content-type":    "application/json;charset=UTF-8",
	"origin":          "https://tv.rubika.ir",
	"referer":         "https://tv.rubika.ir/",
	"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
}

// TV is the ParsRubika TV client.
type TV struct {
	auth      string
	proxy     string
	logger    *log.Logger
	transport *transport.Transport
	version   string
	versionMu sync.Mutex
	versionOK bool
	// versionOnce ensures the MetaURL fetch runs at most once, even
	// when many goroutines call getVersion concurrently. Without it,
	// every concurrent request would block on versionMu until the
	// first HTTP call to MetaURL finished (which can take seconds on
	// a slow network).
	versionOnce sync.Once
}

// Option configures a TV client at construction time.
type Option func(*TV)

// WithAuth sets the auth key.
func WithAuth(auth string) Option {
	return func(t *TV) { t.auth = auth }
}

// WithProxy sets a proxy URL.
func WithProxy(proxy string) Option {
	return func(t *TV) { t.proxy = proxy }
}

// WithLogger injects a custom logger.
func WithLogger(l *log.Logger) Option {
	return func(t *TV) {
		if l != nil {
			t.logger = l
		}
	}
}

// WithVersion overrides the app version (normally auto-detected
// from MetaURL).
func WithVersion(v string) Option {
	return func(t *TV) {
		if v != "" {
			t.version = v
			t.versionOK = true
		}
	}
}

// New constructs a TV client.
func New(opts ...Option) *TV {
	t := &TV{
		logger:  log.New(os.Stderr, "[parsrubika-tv] ", log.LstdFlags|log.Lmsgprefix),
		version: DefaultVersion,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(t)
		}
	}
	t.transport = transport.New(
		transport.WithProxy(t.proxy),
		transport.WithMaxRetries(core.DefaultMaxRetries),
	)
	// Override the default transport headers with TV-specific ones.
	for k, v := range TVHeaders {
		t.transport.SetHeader(k, v)
	}
	return t
}

// getVersion returns the app version, auto-detecting it from MetaURL
// on the first call. Subsequent calls return the cached value without
// any network I/O. The MetaURL fetch runs at most once thanks to
// sync.Once — concurrent callers wait cheaply instead of all blocking
// on the mutex during the HTTP round-trip.
func (t *TV) getVersion(ctx context.Context) string {
	t.versionOnce.Do(func() {
		// Fast path: the caller already supplied a version via
		// WithVersion, so versionOK is true and there's nothing to
		// fetch. sync.Once still runs this Do once, but we no-op.
		t.versionMu.Lock()
		alreadySet := t.versionOK
		t.versionMu.Unlock()
		if alreadySet {
			return
		}
		// Fetch the version from MetaURL. We do NOT hold versionMu
		// during the HTTP call — only when writing the result.
		body, status, err := t.transport.Get(ctx, MetaURL)
		if err != nil || status >= 400 {
			return
		}
		var meta map[string]any
		if json.Unmarshal(body, &meta) != nil {
			return
		}
		v, _ := meta["version"].(string)
		if v == "" {
			return
		}
		t.versionMu.Lock()
		t.version = v
		t.versionOK = true
		t.versionMu.Unlock()
	})
	t.versionMu.Lock()
	defer t.versionMu.Unlock()
	return t.version
}

// request sends a POST request to a random TV server.
func (t *TV) request(ctx context.Context, method string, input map[string]any) (*events.Event, error) {
	if t.auth == "" {
		return nil, coreerrors.New(coreerrors.KindAuth, "tv: auth is required")
	}
	server := TVServerURLs[rand.Intn(len(TVServerURLs))]
	payload := map[string]any{
		"api_version": "1",
		"auth":        t.auth,
		"client": map[string]any{
			"app_name":    "Main",
			"app_version": t.getVersion(ctx),
			"platform":    "TVWeb",
			"package":     "tv.rubika.ir",
			"lang_code":   core.DefaultLangCode,
		},
		"data":   input,
		"method": method,
	}
	body, status, err := t.transport.Post(ctx, server, payload)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, coreerrors.New(coreerrors.KindAPI, fmt.Sprintf("tv: http %d: %s", status, string(body)))
	}
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, coreerrors.Wrap(coreerrors.KindAPI, "tv: decode", err)
	}
	st, _ := raw["status"].(string)
	stDet, _ := raw["status_det"].(string)
	if st == "OK" && stDet == "OK" {
		data, _ := raw["data"].(map[string]any)
		if data == nil {
			data = map[string]any{}
		}
		return events.New(data), nil
	}
	return nil, coreerrors.FromStatus(stDet, fmt.Sprintf("tv: %s failed", method))
}

// SetHTTPClient lets callers swap in an HTTP/3 client.
func (t *TV) SetHTTPClient(c *http.Client) {
	if c != nil {
		t.transport.SetHTTPClient(c)
	}
}

// SetAuth updates the auth key. Used by callers that want to reuse
// the account's auth key (after the account signs in).
func (t *TV) SetAuth(auth string) {
	t.auth = auth
}

// ===== Browsing / listing ==============================================

// GetListItems returns the items of a list (e.g. a category).
// `listType` defaults to "Media"; `startID` defaults to "0".
func (t *TV) GetListItems(ctx context.Context, listID, listType, startID string) (*events.Event, error) {
	if listType == "" {
		listType = "Media"
	}
	if startID == "" {
		startID = "0"
	}
	return t.request(ctx, "getListItems", map[string]any{
		"list_id":  listID,
		"type":     listType,
		"start_id": startID,
	})
}

// GetListing returns the items of a listing page. `listingID`
// defaults to "home".
func (t *TV) GetListing(ctx context.Context, listingID string) (*events.Event, error) {
	if listingID == "" {
		listingID = "home"
	}
	return t.request(ctx, "getListing", map[string]any{"listing_id": listingID})
}

// GetCustomMenuItems returns the user's custom menu.
func (t *TV) GetCustomMenuItems(ctx context.Context) (*events.Event, error) {
	return t.request(ctx, "getCustomMenuItems", map[string]any{})
}

// GetPropertyTypes returns the available property types (genres,
// categories, etc.).
func (t *TV) GetPropertyTypes(ctx context.Context) (*events.Event, error) {
	return t.request(ctx, "getPropertyTypes", map[string]any{})
}

// GetMedia returns info about a single media.
func (t *TV) GetMedia(ctx context.Context, mediaID, trackID string) (*events.Event, error) {
	in := map[string]any{"media_id": mediaID}
	if trackID != "" {
		in["track_id"] = trackID
	}
	return t.request(ctx, "getMedia", in)
}

// Search searches for media by text.
func (t *TV) Search(ctx context.Context, text, startID string) (*events.Event, error) {
	in := map[string]any{"search_text": text}
	if startID != "" {
		in["start_id"] = startID
	}
	return t.request(ctx, "search", in)
}

// GetRelated returns the related medias.
func (t *TV) GetRelated(ctx context.Context, mediaID, startID string) (*events.Event, error) {
	if startID == "" {
		startID = "0"
	}
	return t.request(ctx, "getRelated", map[string]any{
		"media_id": mediaID,
		"start_id": startID,
	})
}

// ===== Cast / series ===================================================

// GetCastMedias returns the medias for a cast (actor/director).
func (t *TV) GetCastMedias(ctx context.Context, castID, startID string) (*events.Event, error) {
	if startID == "" {
		startID = "0"
	}
	return t.request(ctx, "getCastMedias", map[string]any{
		"cast_id":  castID,
		"start_id": startID,
	})
}

// GetCastDetails returns info about a cast member.
func (t *TV) GetCastDetails(ctx context.Context, castID string) (*events.Event, error) {
	return t.request(ctx, "getCastDetails", map[string]any{"cast_id": castID})
}

// GetSeasonEpisodes returns the episodes of a season.
func (t *TV) GetSeasonEpisodes(ctx context.Context, seasonID, seriesID string) (*events.Event, error) {
	return t.request(ctx, "getSeasonEpisodes", map[string]any{
		"season_id": seasonID,
		"series_id": seriesID,
	})
}

// ===== Account / wish-list =============================================

// GetAccountInfo returns the user's TV account info.
func (t *TV) GetAccountInfo(ctx context.Context) (*events.Event, error) {
	return t.request(ctx, "getAccountInfo", map[string]any{})
}

// GetWishList returns the user's wish-list.
func (t *TV) GetWishList(ctx context.Context, startID string) (*events.Event, error) {
	in := map[string]any{}
	if startID != "" {
		in["start_id"] = startID
	}
	return t.request(ctx, "getWishList", in)
}

// ===== Media actions (like / dislike / wish-list) =====================

// actionOnMedia is the low-level helper for media actions.
func (t *TV) actionOnMedia(ctx context.Context, action, mediaID string) (*events.Event, error) {
	return t.request(ctx, "actionOnMedia", map[string]any{
		"action":   action,
		"media_id": mediaID,
	})
}

// LikeMedia adds a like to a media.
func (t *TV) LikeMedia(ctx context.Context, mediaID string) (*events.Event, error) {
	return t.actionOnMedia(ctx, "Like", mediaID)
}

// UnlikeMedia removes a like from a media. The action verb is
// "Unlike" (matching the Rubika like/unlike vocabulary); "Dislike" is
// a separate thumbs-down action, not the inverse of Like.
func (t *TV) UnlikeMedia(ctx context.Context, mediaID string) (*events.Event, error) {
	return t.actionOnMedia(ctx, "Unlike", mediaID)
}

// AddWishMedia adds a media to the wish-list.
func (t *TV) AddWishMedia(ctx context.Context, mediaID string) (*events.Event, error) {
	return t.actionOnMedia(ctx, "AddWishList", mediaID)
}

// RemoveWishMedia removes a media from the wish-list.
func (t *TV) RemoveWishMedia(ctx context.Context, mediaID string) (*events.Event, error) {
	return t.actionOnMedia(ctx, "RemoveWishList", mediaID)
}

// Close releases the underlying transport.
func (t *TV) Close() error { return nil }
