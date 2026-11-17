# `bot` — Options

_Source: [`options.go`](../../ParsRubika/bot/options.go)_

## Symbols

### `Options`

```go
type Options struct
```

Options carries the optional configuration of a Bot.

An Options value can be passed directly to New, or individual
settings can be supplied through the With* functional options.
Both styles can be mixed freely.

### `DefaultConcurrency`

```go
const DefaultConcurrency = 8
```

Defaults shared across the package.

### `DefaultOptions`

```go
func DefaultOptions() Options
```

DefaultOptions returns sensible defaults for a Bot.

### `Option`

```go
type Option interface
```

Option mutates the configuration of a Bot at construction time.

### `WithBaseURL`

```go
func WithBaseURL(u string) Option { return optionFunc(func(o *Options) { o.BaseURL = u }) }
```

WithBaseURL sets the API endpoint template.

### `WithRateLimit`

```go
func WithRateLimit(d time.Duration) Option { return optionFunc(func(o *Options) { o.RateLimit = d }) }
```

WithRateLimit sets the minimum spacing between API calls (0 disables it).

### `WithTimeout`

```go
func WithTimeout(d time.Duration) Option { return optionFunc(func(o *Options) { o.Timeout = d }) }
```

WithTimeout sets the per-request HTTP timeout.

### `WithLongPollTimeout`

```go
func WithLongPollTimeout(d time.Duration) Option
```

WithLongPollTimeout sets the server-side wait passed to getUpdates.
Lower values make the bot more responsive to newly sent messages but
increase the number of HTTP round-trips. Default is 20s.

### `WithPollInterval`

```go
func WithPollInterval(d time.Duration) Option
```

WithPollInterval sets the pause between polling rounds.

### `WithRetries`

```go
func WithRetries(n int, backoff time.Duration) Option
```

WithRetries sets the number of retries (0 disables retrying).

### `WithConcurrency`

```go
func WithConcurrency(n int) Option { return optionFunc(func(o *Options) { o.Concurrency = n }) }
```

WithConcurrency sets how many updates are handled in parallel.

### `WithParseMode`

```go
func WithParseMode(m ParseMode) Option { return optionFunc(func(o *Options) { o.ParseMode = m }) }
```

WithParseMode sets the default parse mode of outgoing text.

### `WithPersistOffset`

```go
func WithPersistOffset(on bool) Option { return optionFunc(func(o *Options) { o.PersistOffset = on }) }
```

WithPersistOffset toggles on-disk offset persistence.

### `WithOffsetDir`

```go
func WithOffsetDir(dir string) Option { return optionFunc(func(o *Options) { o.OffsetDir = dir }) }
```

WithOffsetDir sets the directory used for the offset file.

### `WithSkipBacklog`

```go
func WithSkipBacklog(on bool) Option { return optionFunc(func(o *Options) { o.SkipBacklog = on }) }
```

WithSkipBacklog toggles dropping of the first pending batch.

### `WithHTTPClient`

```go
func WithHTTPClient(c *http.Client) Option { return optionFunc(func(o *Options) { o.HTTPClient = c }) }
```

WithHTTPClient installs a custom HTTP client.

### `WithUserAgent`

```go
func WithUserAgent(ua string) Option { return optionFunc(func(o *Options) { o.UserAgent = ua }) }
```

WithUserAgent overrides the User-Agent header.

### `WithWebhookSecret`

```go
func WithWebhookSecret(s string) Option { return optionFunc(func(o *Options) { o.WebhookSecret = s }) }
```

WithWebhookSecret protects the webhook endpoints with a shared secret.

### `WithErrorHandler`

```go
func WithErrorHandler(fn func(*ErrorContext)) Option
```

WithErrorHandler installs a global error/panic handler.

### `WithLogger`

```go
func WithLogger(l Logger) Option { return optionFunc(func(o *Options) { o.Logger = l }) }
```

WithLogger installs a custom logger.
