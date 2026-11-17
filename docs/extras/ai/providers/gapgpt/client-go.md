# `gapgpt` — Everything-in-one client (App) and one-line helpers

_Source: [`client.go`](../../../../../ParsRubika/extras/ai/providers/gapgpt/client.go)_

## Symbols

### `Name`

```go
const Name = "gapgpt"
```

Name is the provider identifier.

### `DefaultBaseURL`

```go
const DefaultBaseURL = "https://api.gapgpt.app/v1"
```

DefaultBaseURL is the GapGPT API root.

### `DefaultModel`

```go
const DefaultModel = "gapgpt-deepseek-v4.1-flash"
```

DefaultModel is GapGPT's hosted DeepSeek flash model.

### `Client`

```go
type Client struct
```

Client is the GapGPT provider. It implements ai.Provider.

### `Option`

```go
type Option func(*Client)
```

Option configures a Client at construction time.

### `WithAPIKey`

```go
func WithAPIKey(key string) Option
```

WithAPIKey sets the GapGPT API key (Bearer token).

### `WithBaseURL`

```go
func WithBaseURL(url string) Option
```

WithBaseURL overrides the default API root.

### `WithModel`

```go
func WithModel(model string) Option
```

WithModel overrides the default model for every Generate call.

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

New constructs a GapGPT Client.

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

Generate sends a Responses request to GapGPT.

### `Generator`

```go
func Generator(opts ...Option) ai.Generator { return ai.NewGenerator(New(opts...)) }
```

Generator adapts a GapGPT Client to the simple ai.Generator signature.
