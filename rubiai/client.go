// Package rubiai implements the RubiAI assistant client. It talks to the
// RubiAI API at `https://rubichat1.iranlms.ir/` with `api_version: "1"`,
// plain JSON over HTTPS (no envelope encryption), and Server-Sent-Events
// streaming for sendMessage.
//
// The package mirrors the conventions of the rubino and tv packages:
// every method returns `*events.Event` (the dynamic map wrapper), options
// configure the client, and an httptest-friendly `request` path makes the
// whole surface unit-testable without a live RubiAI account.
package rubiai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport"
)

// BaseURL is the default RubiAI API endpoint.
const BaseURL = "https://rubichat1.iranlms.ir/"

// ClientInfo is the `client` block sent inside every RubiAI request. It
// identifies the caller as the official Android Rubika app so the server
// accepts the request.
var ClientInfo = map[string]any{
	"app_name":    "Main",
	"app_version": "4.1.0",
	"package":     "app.rbmain.a",
	"platform":    "Android",
}

// clientBlock is the JSON-encoded form of ClientInfo. It is built once
// (lazily, on first use) and embedded verbatim as a json.RawMessage so
// the encoder never re-traverses or re-sorts the map on the hot path.
var (
	clientBlockOnce sync.Once
	clientBlock     json.RawMessage
)

func clientBlockJSON() json.RawMessage {
	clientBlockOnce.Do(func() {
		b, err := json.Marshal(ClientInfo)
		if err != nil {
			clientBlock = json.RawMessage("{}")
			return
		}
		clientBlock = b
	})
	return clientBlock
}

// requestPayload is the typed envelope for a RubiAI request. A struct
// streams fields in declaration order with no map-key sorting and no
// per-call map allocation, which keeps the encode step cheap.
type requestPayload struct {
	APIVersion string             `json:"api_version"`
	Auth       string             `json:"auth"`
	Data       requestPayloadData `json:"data"`
}

type requestPayloadData struct {
	Input  map[string]any  `json:"input"`
	Method string          `json:"method"`
	Client json.RawMessage `json:"client"`
}

// RubiaAI is the RubiAI assistant client.
type RubiaAI struct {
	auth       string
	userAgent  string
	proxy      string
	timeout    time.Duration
	maxRetries int
	logger     *log.Logger
	transport  *transport.Transport

	apiURL string
}

// Option configures a RubiaAI client at construction time.
type Option func(*RubiaAI)

// WithAuth sets the RubiAI auth key. Required for every request.
func WithAuth(auth string) Option {
	return func(r *RubiaAI) {
		if auth != "" {
			r.auth = auth
		}
	}
}

// WithUserAgent overrides the default User-Agent string.
func WithUserAgent(ua string) Option {
	return func(r *RubiaAI) {
		if ua != "" {
			r.userAgent = ua
		}
	}
}

// WithProxy sets a proxy URL (http://, https://, socks5://, socks5h://).
func WithProxy(proxy string) Option {
	return func(r *RubiaAI) { r.proxy = proxy }
}

// WithTimeout sets the HTTP timeout.
func WithTimeout(d time.Duration) Option {
	return func(r *RubiaAI) {
		if d > 0 {
			r.timeout = d
		}
	}
}

// WithMaxRetries overrides the default retry count for transient errors.
func WithMaxRetries(n int) Option {
	return func(r *RubiaAI) {
		if n > 0 {
			r.maxRetries = n
		}
	}
}

// WithLogger injects a custom logger.
func WithLogger(l *log.Logger) Option {
	return func(r *RubiaAI) {
		if l != nil {
			r.logger = l
		}
	}
}

// New constructs a RubiaAI client.
func New(opts ...Option) *RubiaAI {
	r := &RubiaAI{
		userAgent:  core.UserAgent,
		timeout:    time.Duration(core.DefaultTimeout) * time.Second,
		maxRetries: core.DefaultMaxRetries,
		logger:     log.New(os.Stderr, "[parsrubika-rubiai] ", log.LstdFlags|log.Lmsgprefix),
		apiURL:     BaseURL,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}
	t := transport.New(
		transport.WithUserAgent(r.userAgent),
		transport.WithProxy(r.proxy),
		transport.WithTimeout(r.timeout),
		transport.WithMaxRetries(r.maxRetries),
	)
	r.transport = t
	return r
}

// SetHTTPClient lets callers swap in a pre-configured *http.Client
// (e.g. an HTTP/3 client or one bound to a specific proxy).
func (r *RubiaAI) SetHTTPClient(c *http.Client) {
	if c != nil {
		r.transport.SetHTTPClient(c)
	}
}

// SetAuth updates the auth key. Useful when reusing the account client's
// auth key after the account signs in (RubiAI accepts the account auth).
func (r *RubiaAI) SetAuth(auth string) {
	if auth != "" {
		r.auth = auth
	}
}

// Close releases the underlying transport.
func (r *RubiaAI) Close() error { return nil }

// request sends a normal (non-streaming) RubiAI request. The body shape is:
//
//	{
//	  "api_version": "1",
//	  "auth": "<auth>",
//	  "data": { "input": <input>, "method": <method>, "client": <ClientInfo> }
//	}
//
// On success (data.status == "OK") the data.output field is returned
// (falling back to the whole data block when output is absent). On
// failure the status_det is mapped to a typed error via coreerrors.
func (r *RubiaAI) request(ctx context.Context, method string, input map[string]any) (*events.Event, error) {
	if r.auth == "" {
		return nil, coreerrors.New(coreerrors.KindAuth, "rubiai: auth is required")
	}
	if input == nil {
		input = map[string]any{}
	}
	payload := requestPayload{
		APIVersion: "1",
		Auth:       r.auth,
		Data: requestPayloadData{
			Input:  input,
			Method: method,
			Client: clientBlockJSON(),
		},
	}
	body, status, err := r.transport.Post(ctx, r.apiURL, payload)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, coreerrors.New(coreerrors.KindAPI, fmt.Sprintf("rubiai: http %d: %s", status, string(body)))
	}
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, coreerrors.Wrap(coreerrors.KindAPI, "rubiai: decode", err)
	}
	data, _ := raw["data"].(map[string]any)
	if data == nil {
		return nil, coreerrors.New(coreerrors.KindAPI, "rubiai: response missing data block")
	}
	st, _ := data["status"].(string)
	if st != "OK" {
		stDet, _ := data["status_det"].(string)
		// Normalise INVALID_AUTH → INVALID_ACCESS to match the Python
		// exception taxonomy and the rest of ParsRubika.
		if stDet == "INVALID_AUTH" {
			stDet = "INVALID_ACCESS"
		}
		return nil, coreerrors.FromStatus(stDet, fmt.Sprintf("rubiai: %s failed", method))
	}
	// output is the method-specific payload; when absent, return the
	// whole data block so callers always get an *events.Event.
	if out, ok := data["output"].(map[string]any); ok && out != nil {
		return events.New(out), nil
	}
	return events.New(data), nil
}
