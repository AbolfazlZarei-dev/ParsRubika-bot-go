// Package rubino implements the Rubino social-network client. It
// talks to the Rubino API at `https://rubino1.iranlms.ir/` with
// `api_version: "0"`. The client is fully modular: every method
// group lives in its own file (posts, stories, follow, comments,
// profiles, uploads).
package rubino

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport"
)

// RubinoAPIURL is the base URL for the Rubino API.
const RubinoAPIURL = "https://rubino1.iranlms.ir/"

// clientBlock is the JSON-encoded `client` object sent with every
// request. It is constant, so it is marshalled once at startup and
// embedded verbatim as a json.RawMessage — the encoder writes the
// pre-computed bytes straight to the wire without re-traversing or
// re-sorting the map on every call.
var (
	clientBlockOnce sync.Once
	clientBlock     json.RawMessage
)

func clientBlockJSON() json.RawMessage {
	clientBlockOnce.Do(func() {
		b, err := json.Marshal(map[string]any{
			"app_name":    "Main",
			"app_version": "3.7.3",
			"platform":    "Android",
			"package":     "app.rbmain.a",
			"lang_code":   core.DefaultLangCode,
			"store":       "Myket",
			"temp_code":   "34",
		})
		if err != nil {
			// A static, well-formed map never fails to encode; fall back
			// to an empty object so the request still goes out shaped
			// correctly.
			clientBlock = json.RawMessage("{}")
			return
		}
		clientBlock = b
	})
	return clientBlock
}

// requestPayload is the typed envelope for a Rubino request. Using a
// struct instead of a map[string]any lets the JSON encoder stream
// fields in declaration order with no map-key sorting and no
// intermediate map allocation on the hot path.
type requestPayload struct {
	APIVersion string          `json:"api_version"`
	Auth       string          `json:"auth"`
	Client     json.RawMessage `json:"client"`
	Data       map[string]any  `json:"data"`
	Method     string          `json:"method"`
}

// Rubino is the Rubino social-network client.
type Rubino struct {
	auth      string
	userAgent string
	proxy     string
	timeout   int
	logger    *log.Logger
	transport *transport.Transport

	mu      sync.Mutex
	apiURL  string
	profile string // the active profile_id, used by some methods
}

// Option configures a Rubino at construction time.
type Option func(*Rubino)

// WithAuth sets the auth key.
func WithAuth(auth string) Option {
	return func(s *Rubino) { s.auth = auth }
}

// WithUserAgent overrides the User-Agent string.
func WithUserAgent(ua string) Option {
	return func(s *Rubino) {
		if ua != "" {
			s.userAgent = ua
		}
	}
}

// WithProxy sets a proxy URL.
func WithProxy(proxy string) Option {
	return func(s *Rubino) { s.proxy = proxy }
}

// WithTimeout sets the HTTP timeout in seconds.
func WithTimeout(seconds int) Option {
	return func(s *Rubino) {
		if seconds > 0 {
			s.timeout = seconds
		}
	}
}

// WithLogger injects a custom logger.
func WithLogger(l *log.Logger) Option {
	return func(s *Rubino) {
		if l != nil {
			s.logger = l
		}
	}
}

// WithProfile sets the active profile ID. Some methods use this
// default when their `profileID` parameter is empty.
func WithProfile(profileID string) Option {
	return func(s *Rubino) { s.profile = profileID }
}

// New constructs a Rubino client.
func New(opts ...Option) *Rubino {
	s := &Rubino{
		userAgent: core.UserAgent,
		timeout:   core.DefaultTimeout,
		logger:    log.New(os.Stderr, "[parsrubika-rubino] ", log.LstdFlags|log.Lmsgprefix),
		apiURL:    RubinoAPIURL,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	t := transport.New(
		transport.WithUserAgent(s.userAgent),
		transport.WithProxy(s.proxy),
		transport.WithMaxRetries(core.DefaultMaxRetries),
		// Wire the configured timeout through to the transport. Previously
		// s.timeout was set by WithTimeout but never read, so the option
		// was a silent no-op and the transport kept its 30 s default.
		transport.WithTimeout(time.Duration(s.timeout)*time.Second),
	)
	s.transport = t
	return s
}

// request sends a POST request to the Rubino API. The body shape is:
//
//	{
//	 "api_version": "0",
//	 "auth": "<auth>",
//	 "client": {app_name, app_version, platform, package, lang_code},
//	 "data": <input>,
//	 "method": <method>
//	}
func (s *Rubino) request(ctx context.Context, method string, input map[string]any) (*events.Event, error) {
	if s.auth == "" {
		return nil, coreerrors.New(coreerrors.KindAuth, "rubino: auth is required")
	}
	payload := requestPayload{
		APIVersion: "0",
		Auth:       s.auth,
		Client:     clientBlockJSON(),
		Data:       input,
		Method:     method,
	}
	body, status, err := s.transport.Post(ctx, s.apiURL, payload)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, coreerrors.New(coreerrors.KindAPI, fmt.Sprintf("rubino: http %d: %s", status, string(body)))
	}
	var raw map[string]any
	if err := jsonUnmarshal(body, &raw); err != nil {
		return nil, coreerrors.Wrap(coreerrors.KindAPI, "rubino: decode", err)
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
	return nil, coreerrors.FromStatus(stDet, fmt.Sprintf("rubino: %s failed", method))
}

// profileOr returns the supplied profileID, falling back to the
// Rubino's default profile when empty.
func (s *Rubino) profileOr(profileID string) string {
	if profileID != "" {
		return profileID
	}
	return s.profile
}

// rnd returns a random int used by some Rubino methods.
func rnd() int64 { return rand.Int63n(1_000_000_000) }

// Close releases the underlying transport.
func (s *Rubino) Close() error { return nil }

// SetHTTPClient lets callers swap in an HTTP/3 client.
func (s *Rubino) SetHTTPClient(c *http.Client) {
	if c != nil {
		s.transport.SetHTTPClient(c)
	}
}

// SetAuth updates the auth key. Used by callers that want to reuse
// the account's auth key (after the account signs in) instead of a
// separate Rubino login — Rubino accepts the account's auth key.
func (s *Rubino) SetAuth(auth string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.auth = auth
}
