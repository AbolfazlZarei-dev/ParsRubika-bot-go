// Package account implements the ParsRubika user/account client. It
// speaks the encrypted Rubika API (AES-CBC + RSA-PKCS1v15 signing) and
// listens for live updates over a WebSocket connection.
//
// The Account struct is intentionally split across several files so
// each concern lives in its own unit:
//
// - account.go : the Account struct + constructor + options
// - lifecycle.go : Start / Stop / Run / connect
// - dispatch.go : the update dispatcher + listener registration
// - request.go : the encrypted request builder
// - methods_*.go : API method groups (auth, users, chats, ...)
package account

import (
	"context"
	"crypto/rsa"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/enums"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/files"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/security"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/storage"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/text"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport"
)

// Account is the high-level user/account client.
type Account struct {
	// Configuration.
	name        string
	auth        string
	privateKey  string
	phoneNumber string
	userAgent   string
	timeout     time.Duration
	langCode    string
	parseMode   enums.ParseMode
	proxy       string
	platform    string
	maxRetries  int
	sequential  bool
	logger      core.Logger

	// Identity sent inside the `client` field of every request body.
	identity core.Identity

	// Formatter + session storage + transport.
	formatter  text.Formatter
	store      storage.Store
	transport  *transport.Transport
	uploader   *files.Uploader
	downloader *files.Downloader

	// Crypto state — populated during Start().
	decodedAuth  string
	symmetricKey string
	// symmetricKeyBytes is the []byte form of symmetricKey, cached so the
	// hot Request path is allocation-free for the key material
	// (Decrypt/EncryptJSON consume this []byte directly).
	symmetricKeyBytes []byte
	rsaKey            *rsa.PrivateKey
	guid              string

	// Discovered DCs.
	dcConfig *transport.DCConfig

	// initMu guards the lazy initialisation of dcConfig, auth and
	// symmetricKey inside Request (and the WS keep-alive path). It
	// prevents two concurrent Request calls from racing on those
	// fields and double-discovering DCs.
	initMu sync.Mutex

	// debug, when true, enables verbose request/response logging.
	// Off by default so plaintext payloads are never logged.
	debug bool

	// codeProvider, when non-nil, is used by signInInteractive to
	// obtain the SMS code / 2FA password. When nil, the interactive
	// stdin/stdout prompt is used instead. Setting a CodeProvider lets
	// the account client run headlessly (servers, Docker, systemd).
	codeProvider CodeProvider

	// WebSocket connection.
	ws       *transport.WebSocket
	wsCancel context.CancelFunc
	// wsMu guards the ws and cancel fields. Run/Listen write them while
	// Stop/Close read them from another goroutine; without the lock that
	// is a data race even under steady-state operation.
	wsMu sync.Mutex

	// Update dispatching.
	listeners      map[int]listenerEntry
	listenersMu    sync.RWMutex
	nextListenerID int

	// Lifecycle.
	ctx    context.Context
	cancel context.CancelFunc
	// stopped is set by Stop/Close so the reconnect loop in Listen can
	// tell a user-initiated shutdown apart from a dropped connection and
	// avoid redialling after the caller asked to stop.
	stopped atomic.Bool
}

// Listener is the signature of a registered update callback.
type Listener func(ctx context.Context, e *events.Event) error

// CodeProvider is the signature of the callback used by the interactive
// sign-in flow to obtain the SMS code (or 2FA password) from the user.
// The `hint` argument is a short human-readable description of what is
// being requested (e.g. "Code:" or "Password [hint] > "). The provider
// should return the user's response, trimmed of trailing whitespace.
//
// When no CodeProvider is set, sign-in falls back to reading from stdin
// (interactive mode). Setting a CodeProvider makes the account usable
// in headless environments (servers, Docker, systemd).
type CodeProvider func(ctx context.Context, hint string) (string, error)

// Matcher is the contract implemented by every account-side matcher
// (text, command, chat-type, ...). The matchers live in the
// `parsrubika/account/matchers` sub-package.
type Matcher interface {
	// Kind returns the CamelCase update-name this matcher accepts
	// (e.g. "MessageUpdates", "ChatUpdates", "ShowActivities",
	// "ShowNotifications", "RemoveNotifications"). These names are
	// the `capitalize()` of the snake_case keys the server sends in
	// the decrypted `data_enc` payload.
	Kind() string
	// Match returns true if the listener should fire for the supplied
	// raw update payload.
	Match(ctx context.Context, raw map[string]any) bool
}

// listenerEntry pairs a user callback with its matcher.
type listenerEntry struct {
	fn      Listener
	matcher Matcher
}

// Option configures an Account at construction time.
type Option func(*Account) error

// ===== Options ===========================================================

// WithAuth sets the auth key.
func WithAuth(auth string) Option {
	return func(a *Account) error {
		if auth == "" {
			return coreerrors.New(coreerrors.KindValidation, "auth cannot be empty")
		}
		a.auth = auth
		return nil
	}
}

// WithPrivateKey sets the RSA private key (PEM-encoded).
func WithPrivateKey(pk string) Option {
	return func(a *Account) error {
		a.privateKey = security.NormalizePrivateKey(pk)
		return nil
	}
}

// WithPhoneNumber sets the phone number used during the registration
// handshake.
func WithPhoneNumber(phone string) Option {
	return func(a *Account) error { a.phoneNumber = phone; return nil }
}

// WithUserAgent overrides the default User-Agent string.
func WithUserAgent(ua string) Option {
	return func(a *Account) error {
		if ua != "" {
			a.userAgent = ua
		}
		return nil
	}
}

// WithTimeout sets the HTTP timeout.
func WithTimeout(d time.Duration) Option {
	return func(a *Account) error {
		if d <= 0 {
			return coreerrors.New(coreerrors.KindValidation, "timeout must be positive")
		}
		a.timeout = d
		return nil
	}
}

// WithLangCode sets the language code (default "fa").
func WithLangCode(code string) Option {
	return func(a *Account) error {
		if code == "" {
			code = core.DefaultLangCode
		}
		a.langCode = code
		return nil
	}
}

// WithParseMode sets the default parse mode.
func WithParseMode(mode enums.ParseMode) Option {
	return func(a *Account) error { a.parseMode = mode; return nil }
}

// WithProxy sets an HTTP(S) proxy URL.
func WithProxy(proxy string) Option {
	return func(a *Account) error { a.proxy = proxy; return nil }
}

// WithPlatform selects "PWA" (default) or "Android".
func WithPlatform(p string) Option {
	return func(a *Account) error {
		switch p {
		case "PWA", "Android":
			a.platform = p
		default:
			return coreerrors.New(coreerrors.KindValidation, "unsupported platform: "+p)
		}
		return nil
	}
}

// WithMaxRetries overrides the default max-retries value.
func WithMaxRetries(n int) Option {
	return func(a *Account) error {
		if n <= 0 {
			n = core.DefaultMaxRetries
		}
		a.maxRetries = n
		return nil
	}
}

// WithSequentialListeners, when true, makes the dispatcher run
// listeners sequentially and stop after the first match.
func WithSequentialListeners(b bool) Option {
	return func(a *Account) error { a.sequential = b; return nil }
}

// WithLogger injects a custom logger.
func WithLogger(l core.Logger) Option {
	return func(a *Account) error {
		if l != nil {
			a.logger = l
		}
		return nil
	}
}

// WithDebug enables verbose request/response/WebSocket logging for
// troubleshooting. Off by default — when disabled, plaintext payloads
// are never logged.
func WithDebug(b bool) Option {
	return func(a *Account) error { a.debug = b; return nil }
}

// WithCodeProvider installs a callback used by the interactive sign-in
// flow to obtain the SMS code (or 2FA password) from the caller. When
// set, the account client can run headlessly (servers, Docker, systemd)
// instead of prompting on stdin/stdout.
//
// Example:
//
//	acc, _ := account.New(
//	 account.WithPhoneNumber("+989123456789"),
//	 account.WithCodeProvider(func(ctx context.Context, hint string) (string, error) {
//	 // hint is "Code:" or "Password [hint] > "
//	 return myHTTPUI.PromptUser(ctx, hint)
//	 }),
//	)
func WithCodeProvider(p CodeProvider) Option {
	return func(a *Account) error {
		if p == nil {
			return coreerrors.New(coreerrors.KindValidation, "code provider cannot be nil")
		}
		a.codeProvider = p
		return nil
	}
}

// WithStore lets the caller supply a custom Store implementation.
func WithStore(s storage.Store) Option {
	return func(a *Account) error { a.store = s; return nil }
}

// WithSQLiteSession uses a SQLite session file with the given name.
//
// The session is stored in plain text. For at-rest encryption of the
// auth key and RSA private key use WithSQLiteSessionPassphrase instead.
func WithSQLiteSession(name string) Option {
	return func(a *Account) error {
		s, err := storage.NewSQLiteStore(name)
		if err != nil {
			return err
		}
		a.store = s
		a.name = name
		return nil
	}
}

// WithSQLiteSessionPassphrase uses a SQLite session file with the given
// name and enables AES-256-GCM at-rest encryption for the auth and
// private_key columns. The passphrase is NOT stored alongside the
// session file; losing it makes the session unrecoverable.
func WithSQLiteSessionPassphrase(name, passphrase string) Option {
	return func(a *Account) error {
		s, err := storage.NewSQLiteStoreWithPassphrase(name, passphrase)
		if err != nil {
			return err
		}
		a.store = s
		a.name = name
		return nil
	}
}

// WithFileSession uses a JSON file session with the given name.
//
// The session is stored in plain text. For at-rest encryption of the
// auth key and RSA private key use WithFileSessionPassphrase instead.
func WithFileSession(name string) Option {
	return func(a *Account) error {
		s, err := storage.NewFileStore(name)
		if err != nil {
			return err
		}
		a.store = s
		a.name = name
		return nil
	}
}

// WithFileSessionPassphrase uses a JSON file session with the given name
// and enables AES-256-GCM at-rest encryption for the auth and
// private_key fields. The passphrase is NOT stored alongside the session
// file; losing it makes the session unrecoverable.
func WithFileSessionPassphrase(name, passphrase string) Option {
	return func(a *Account) error {
		s, err := storage.NewFileStoreWithPassphrase(name, passphrase)
		if err != nil {
			return err
		}
		a.store = s
		a.name = name
		return nil
	}
}

// WithHTTPClient lets the caller supply a pre-configured *http.Client
// (e.g. a HTTP/3 client). When set, the supplied client is used as-is
// and the proxy/timeout options are ignored.
func WithHTTPClient(c *http.Client) Option {
	return func(a *Account) error {
		// Create a transport that uses the caller's HTTP client.
		t := transport.New()
		t.SetHTTPClient(c)
		a.transport = t
		return nil
	}
}

// ===== Constructor ======================================================

// New constructs an Account using the supplied options.
func New(opts ...Option) (*Account, error) {
	a := &Account{
		name:       "ParsRubika",
		userAgent:  core.UserAgent,
		timeout:    time.Duration(core.DefaultTimeout) * time.Second,
		langCode:   core.DefaultLangCode,
		parseMode:  enums.ParseModeMarkdown,
		platform:   "PWA",
		maxRetries: core.DefaultMaxRetries,
		logger:     log.New(os.Stderr, "[parsrubika-account] ", log.LstdFlags|log.Lmsgprefix),
		listeners:  map[int]listenerEntry{},
	}
	a.identity = core.WebIdentity(a.langCode)
	for _, opt := range opts {
		if opt != nil {
			if err := opt(a); err != nil {
				return nil, err
			}
		}
	}
	if a.store == nil {
		a.store = storage.NewMemoryStore()
	}
	if a.transport == nil {
		a.transport = transport.New(
			transport.WithUserAgent(a.userAgent),
			transport.WithProxy(a.proxy),
			transport.WithTimeout(a.timeout),
			transport.WithMaxRetries(a.maxRetries),
		)
	}
	if a.platform == "Android" {
		a.identity = core.AndroidIdentity(a.langCode)
		a.transport.SetHeader("user-agent", "okhttp/3.12.1")
		a.transport.SetHeader("origin", "")
		a.transport.SetHeader("referer", "")
	}
	// The public Uploader()/Downloader() accessors return these, so they
	// must be usable directly: hand them the transport's HTTP client
	// (non-nil) rather than nil, which would nil-panic on Upload/Download.
	httpClient := a.transport.HTTPClient()
	a.uploader = files.NewUploader(httpClient, core.DefaultChunkSize, a.maxRetries, time.Second,
		map[string]string{"auth": a.auth, "user-agent": a.userAgent})
	a.downloader = files.NewDownloader(httpClient, core.DefaultDownloadChunkSize, a.maxRetries, time.Second,
		map[string]string{"auth": a.auth, "user-agent": a.userAgent})
	return a, nil
}

// ===== Public accessors =================================================

// Name returns the session name.
func (a *Account) Name() string { return a.name }

// Auth returns the configured auth key.
func (a *Account) Auth() string { return a.auth }

// GUID returns the GUID of the logged-in user.
func (a *Account) GUID() string { return a.guid }

// ParseMode returns the default parse mode.
func (a *Account) ParseMode() enums.ParseMode { return a.parseMode }

// UserAgent returns the configured User-Agent string.
func (a *Account) UserAgent() string { return a.userAgent }

// Store returns the session store.
func (a *Account) Store() storage.Store { return a.store }

// Formatter returns the text formatter.
func (a *Account) Formatter() text.Formatter { return a.formatter }

// Uploader returns the file uploader.
func (a *Account) Uploader() *files.Uploader { return a.uploader }

// Downloader returns the file downloader.
func (a *Account) Downloader() *files.Downloader { return a.downloader }

// Close releases the WebSocket + transport + store. It is idempotent and
// signals the reconnect loop (if running) not to redial afterwards.
func (a *Account) Close() error {
	a.stopped.Store(true)
	a.wsMu.Lock()
	if a.cancel != nil {
		a.cancel()
	}
	ws := a.ws
	a.wsMu.Unlock()
	if ws != nil {
		_ = ws.Close()
	}
	if a.store != nil {
		_ = a.store.Close()
	}
	return nil
}
