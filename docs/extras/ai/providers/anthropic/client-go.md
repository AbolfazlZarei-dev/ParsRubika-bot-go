# `anthropic` — Everything-in-one client (App) and one-line helpers

_Source: [`client.go`](../../../../../ParsRubika/extras/ai/providers/anthropic/client.go)_

## Symbols

### `Name`

```go
const Name = "anthropic"
```

Name is the provider identifier.

### `DefaultBaseURL`

```go
const DefaultBaseURL = "https://api.anthropic.com/v1"
```

DefaultBaseURL is the Anthropic API root.

### `DefaultModel`

```go
const DefaultModel = "claude-opus-5-5"
```

DefaultModel is Anthropic's recommended default for most workloads.
A fixed id with no date suffix; override per call for any other
Claude model.

### `APIVersion`

```go
const APIVersion = "2023-06-01"
```

APIVersion is the anthropic-version header value. The classic
stable version is widely supported across all Claude models.

### `DefaultMaxTokens`

```go
const DefaultMaxTokens = 1024
```

DefaultMaxTokens is used when ai.Request.MaxTokens is 0, because
the Messages API requires a non-zero max_tokens.

### `Client`

```go
type Client struct
```

Client is the Anthropic provider. It implements ai.Provider.

### `Option`

```go
type Option func(*Client)
```

Option configures a Client at construction time.

### `WithAPIKey`

```go
func WithAPIKey(key string) Option
```

WithAPIKey sets the Anthropic API key (sent as x-api-key).

### `WithBaseURL`

```go
func WithBaseURL(url string) Option
```

WithBaseURL overrides the default API root.

### `WithModel`

```go
func WithModel(model string) Option
```

WithModel overrides the default Claude model for every Generate call.

### `WithHTTPClient`

```go
func WithHTTPClient(h *http.Client) Option
```

WithHTTPClient injects a pre-configured *http.Client.

### `WithHeader`

```go
func WithHeader(key, value string) Option
```

WithHeader adds an extra request header.

### `WithUserAgent`

```go
func WithUserAgent(ua string) Option
```

WithUserAgent overrides the default User-Agent header.

### `WithTimeout`

```go
func WithTimeout(d time.Duration) Option { return func(c *Client) { c.timeout = d } }
```

WithTimeout sets the HTTP timeout when no custom *http.Client is set.

### `New`

```go
func New(opts ...Option) *Client
```

New constructs an Anthropic Client.

### `Name`

```go
func (c *Client) Name() string { return Name }
```

Name returns the provider identifier.

### `DefaultModel`

```go
func (c *Client) DefaultModel() string { return c.model }
```

DefaultModel returns the model used when ai.Request.Model is empty.

### `Generate`

```go
func (c *Client) Generate(ctx context.Context, req ai.Request) (ai.Response, error)
```

Generate sends a Messages API request to Anthropic.

### `Generator`

```go
func Generator(opts ...Option) ai.Generator { return ai.NewGenerator(New(opts...)) }
```

Generator adapts an Anthropic Client to the simple ai.Generator
signature.
