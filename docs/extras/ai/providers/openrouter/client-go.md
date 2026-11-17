# `openrouter` — Everything-in-one client (App) and one-line helpers

_Source: [`client.go`](../../../../../ParsRubika/extras/ai/providers/openrouter/client.go)_

## Symbols

### `Name`

```go
const Name = "openrouter"
```

Name is the provider identifier.

### `DefaultBaseURL`

```go
const DefaultBaseURL = "https://openrouter.ai/api/v1"
```

DefaultBaseURL is the OpenRouter API root.

### `DefaultModel`

```go
const DefaultModel = "openai/gpt-4o-mini"
```

DefaultModel is the fallback model when ai.Request.Model is empty.
OpenRouter fronts hundreds of models — always pass the exact
`vendor/model` slug you want.

### `Client`

```go
type Client struct
```

Client is the OpenRouter provider. It implements ai.Provider.

### `Option`

```go
type Option func(*Client)
```

Option configures a Client at construction time.

### `WithAPIKey`

```go
func WithAPIKey(key string) Option
```

WithAPIKey sets the OpenRouter API key (required).

### `WithBaseURL`

```go
func WithBaseURL(url string) Option
```

WithBaseURL overrides the default API root.

### `WithModel`

```go
func WithModel(model string) Option
```

WithModel overrides the default model slug for every Generate call.

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

### `WithRankings`

```go
func WithRankings(siteURL, title string) Option
```

WithRankings sets the optional HTTP-Referer (your site URL) and
X-Title (your app name) headers so your app shows up on OpenRouter's
leaderboards. Both are optional per OpenRouter's API.

### `New`

```go
func New(opts ...Option) *Client
```

New constructs an OpenRouter Client.

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

Generate sends a chat-completions request to OpenRouter.

### `Generator`

```go
func Generator(opts ...Option) ai.Generator { return ai.NewGenerator(New(opts...)) }
```

Generator adapts an OpenRouter Client to the simple ai.Generator
signature.
