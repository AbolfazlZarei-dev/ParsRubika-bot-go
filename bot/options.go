package bot

import (
	"net/http"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
)

// Options carries the optional configuration of a Bot.
//
// An Options value can be passed directly to New, or individual
// settings can be supplied through the With* functional options.
// Both styles can be mixed freely.
type Options struct {
	// BaseURL is the API endpoint template; "%s" is replaced by the token.
	BaseURL string
	// RateLimit is the minimum spacing between outgoing API calls.
	// Long-poll requests are not subject to it.
	RateLimit time.Duration
	// UseWebhook is kept for source compatibility; Run selects the mode
	// from its webhookURL argument.
	UseWebhook bool
	// Timeout is the per-request HTTP timeout.
	Timeout time.Duration
	// PollInterval is the pause between two polling rounds.
	PollInterval time.Duration
	// LongPollTimeout is the server-side wait passed to getUpdates.
	LongPollTimeout time.Duration
	// MaxRetries is the number of retries (not attempts) per request.
	MaxRetries int
	// BackoffFactor is the base delay of the exponential backoff.
	BackoffFactor time.Duration
	// MaxBackoff caps a single backoff sleep.
	MaxBackoff time.Duration
	// RetryStatuses lists the HTTP statuses that are retried.
	RetryStatuses map[int]struct{}
	// UserAgent overrides the default User-Agent header.
	UserAgent string
	// ParseMode is the default text parse mode.
	ParseMode ParseMode
	// PersistOffset stores the polling offset on disk between runs.
	PersistOffset bool
	// OffsetDir is the directory of the offset file (default: ".").
	OffsetDir string
	// HTTPClient replaces the default HTTP client.
	HTTPClient *http.Client

	// Concurrency is the number of updates processed in parallel.
	// Updates of the same chat are always processed in order.
	// Zero selects DefaultConcurrency; 1 gives fully sequential handling.
	Concurrency int
	// SkipBacklog drops the first batch of pending updates when the bot
	// starts without a stored offset (so old messages are not replayed).
	SkipBacklog bool
	// DedupSize is the number of recent update keys remembered to
	// suppress duplicates.
	DedupSize int
	// ShutdownTimeout is how long Stop/cancellation waits for in-flight
	// handlers to finish.
	ShutdownTimeout time.Duration
	// WebhookSecret, when set, must be echoed by the caller in the "s"
	// query parameter of every webhook request.
	WebhookSecret string
	// MaxBodyBytes limits the size of a webhook request body.
	MaxBodyBytes int64
	// ErrorHandler receives every handler/middleware error and panic.
	ErrorHandler func(ctx *ErrorContext)
	// Logger replaces the default logger sink.
	Logger Logger
}

// Defaults shared across the package.
const (
	DefaultConcurrency = 8
	defaultDedupSize   = 2048
)

// DefaultOptions returns sensible defaults for a Bot.
func DefaultOptions() Options {
	return Options{
		BaseURL:         BaseURL,
		RateLimit:       100 * time.Millisecond,
		Timeout:         35 * time.Second,
		PollInterval:    300 * time.Millisecond,
		LongPollTimeout: 20 * time.Second,
		MaxRetries:      3,
		BackoffFactor:   400 * time.Millisecond,
		MaxBackoff:      15 * time.Second,
		RetryStatuses:   map[int]struct{}{408: {}, 425: {}, 429: {}, 500: {}, 502: {}, 503: {}, 504: {}},
		UserAgent:       core.UserAgent,
		ParseMode:       ParseMode("markdown"),
		PersistOffset:   true,
		Concurrency:     DefaultConcurrency,
		SkipBacklog:     true,
		DedupSize:       defaultDedupSize,
		ShutdownTimeout: 10 * time.Second,
		MaxBodyBytes:    1 << 20,
	}
}

// Option mutates the configuration of a Bot at construction time.
type Option interface{ apply(*Options) }

type optionFunc func(*Options)

func (f optionFunc) apply(o *Options) { f(o) }

// apply lets a plain Options value act as an Option: every non-zero
// field overrides the current configuration. The two boolean switches
// (UseWebhook, PersistOffset, SkipBacklog) are always copied, as before.
func (src Options) apply(dst *Options) {
	if src.BaseURL != "" {
		dst.BaseURL = src.BaseURL
	}
	if src.RateLimit > 0 {
		dst.RateLimit = src.RateLimit
	}
	dst.UseWebhook = src.UseWebhook
	if src.Timeout > 0 {
		dst.Timeout = src.Timeout
	}
	if src.PollInterval > 0 {
		dst.PollInterval = src.PollInterval
	}
	if src.LongPollTimeout > 0 {
		dst.LongPollTimeout = src.LongPollTimeout
	}
	if src.MaxRetries > 0 {
		dst.MaxRetries = src.MaxRetries
	}
	if src.BackoffFactor > 0 {
		dst.BackoffFactor = src.BackoffFactor
	}
	if src.MaxBackoff > 0 {
		dst.MaxBackoff = src.MaxBackoff
	}
	if src.RetryStatuses != nil {
		dst.RetryStatuses = src.RetryStatuses
	}
	if src.UserAgent != "" {
		dst.UserAgent = src.UserAgent
	}
	if src.ParseMode != "" {
		dst.ParseMode = src.ParseMode
	}
	dst.PersistOffset = src.PersistOffset
	if src.OffsetDir != "" {
		dst.OffsetDir = src.OffsetDir
	}
	if src.HTTPClient != nil {
		dst.HTTPClient = src.HTTPClient
	}
	if src.Concurrency > 0 {
		dst.Concurrency = src.Concurrency
	}
	if src.DedupSize > 0 {
		dst.DedupSize = src.DedupSize
	}
	if src.ShutdownTimeout > 0 {
		dst.ShutdownTimeout = src.ShutdownTimeout
	}
	if src.WebhookSecret != "" {
		dst.WebhookSecret = src.WebhookSecret
	}
	if src.MaxBodyBytes > 0 {
		dst.MaxBodyBytes = src.MaxBodyBytes
	}
	if src.ErrorHandler != nil {
		dst.ErrorHandler = src.ErrorHandler
	}
	if src.Logger != nil {
		dst.Logger = src.Logger
	}
	if src.SkipBacklog {
		dst.SkipBacklog = true
	}
}

// WithBaseURL sets the API endpoint template.
func WithBaseURL(u string) Option { return optionFunc(func(o *Options) { o.BaseURL = u }) }

// WithRateLimit sets the minimum spacing between API calls (0 disables it).
func WithRateLimit(d time.Duration) Option { return optionFunc(func(o *Options) { o.RateLimit = d }) }

// WithTimeout sets the per-request HTTP timeout.
func WithTimeout(d time.Duration) Option { return optionFunc(func(o *Options) { o.Timeout = d }) }

// WithLongPollTimeout sets the server-side wait passed to getUpdates.
// Lower values make the bot more responsive to newly sent messages but
// increase the number of HTTP round-trips. Default is 20s.
func WithLongPollTimeout(d time.Duration) Option {
	return optionFunc(func(o *Options) { o.LongPollTimeout = d })
}

// WithPollInterval sets the pause between polling rounds.
func WithPollInterval(d time.Duration) Option {
	return optionFunc(func(o *Options) { o.PollInterval = d })
}

// WithRetries sets the number of retries (0 disables retrying).
func WithRetries(n int, backoff time.Duration) Option {
	return optionFunc(func(o *Options) {
		o.MaxRetries = n
		if backoff > 0 {
			o.BackoffFactor = backoff
		}
	})
}

// WithConcurrency sets how many updates are handled in parallel.
func WithConcurrency(n int) Option { return optionFunc(func(o *Options) { o.Concurrency = n }) }

// WithParseMode sets the default parse mode of outgoing text.
func WithParseMode(m ParseMode) Option { return optionFunc(func(o *Options) { o.ParseMode = m }) }

// WithPersistOffset toggles on-disk offset persistence.
func WithPersistOffset(on bool) Option { return optionFunc(func(o *Options) { o.PersistOffset = on }) }

// WithOffsetDir sets the directory used for the offset file.
func WithOffsetDir(dir string) Option { return optionFunc(func(o *Options) { o.OffsetDir = dir }) }

// WithSkipBacklog toggles dropping of the first pending batch.
func WithSkipBacklog(on bool) Option { return optionFunc(func(o *Options) { o.SkipBacklog = on }) }

// WithHTTPClient installs a custom HTTP client.
func WithHTTPClient(c *http.Client) Option { return optionFunc(func(o *Options) { o.HTTPClient = c }) }

// WithUserAgent overrides the User-Agent header.
func WithUserAgent(ua string) Option { return optionFunc(func(o *Options) { o.UserAgent = ua }) }

// WithWebhookSecret protects the webhook endpoints with a shared secret.
func WithWebhookSecret(s string) Option { return optionFunc(func(o *Options) { o.WebhookSecret = s }) }

// WithErrorHandler installs a global error/panic handler.
func WithErrorHandler(fn func(*ErrorContext)) Option {
	return optionFunc(func(o *Options) { o.ErrorHandler = fn })
}

// WithLogger installs a custom logger.
func WithLogger(l Logger) Option { return optionFunc(func(o *Options) { o.Logger = l }) }
