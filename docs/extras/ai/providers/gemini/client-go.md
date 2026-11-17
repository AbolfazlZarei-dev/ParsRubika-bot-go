# `gemini` — Everything-in-one client (App) and one-line helpers

_Source: [`client.go`](../../../../../ParsRubika/extras/ai/providers/gemini/client.go)_

## Symbols

### `Name`

```go
const Name = "gemini"
```

Name is the provider identifier.

### `DefaultBaseURL`

```go
const DefaultBaseURL = "https://generativelanguage.googleapis.com/v1beta"
```

DefaultBaseURL is the Gemini API root (without the model segment).

### `DefaultModel`

```go
const DefaultModel = "gemini-3.5-flash-lite"
```

DefaultModel is the model used when ai.Request.Model is empty.

### `Client`

```go
type Client struct
```

Client is the Gemini provider. It implements ai.Provider.

### `Option`

```go
type Option func(*Client)
```

Option configures a Client at construction time.

### `WithAPIKey`

```go
func WithAPIKey(key string) Option
```

WithAPIKey sets the Gemini API key (sent as the ?key= query param).

### `WithBaseURL`

```go
func WithBaseURL(url string) Option
```

WithBaseURL overrides the default API root.

### `WithModel`

```go
func WithModel(model string) Option
```

WithModel overrides the default Gemini model for every Generate call.

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

New constructs a Gemini Client.

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

Generate sends a generateContent request to Gemini.

### `Generator`

```go
func Generator(opts ...Option) ai.Generator { return ai.NewGenerator(New(opts...)) }
```

Generator adapts a Gemini Client to the simple ai.Generator signature.
