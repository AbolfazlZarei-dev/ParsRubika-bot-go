# `openai` — Everything-in-one client (App) and one-line helpers

_Source: [`client.go`](../../../../../ParsRubika/extras/ai/providers/openai/client.go)_

## Symbols

### `Name`

```go
const Name = "openai"
```

Name is the provider identifier surfaced in ai.Response.Provider.

### `DefaultBaseURL`

```go
const DefaultBaseURL = "https://api.openai.com/v1"
```

DefaultBaseURL is the OpenAI API root.

### `DefaultModel`

```go
const DefaultModel = "gpt-4o-mini"
```

DefaultModel is the model used when ai.Request.Model is empty. It is
a stable, low-cost default; override it per call (or at construction
with WithModel) for any other OpenAI model.

### `Client`

```go
type Client struct
```

Client is the OpenAI provider. It implements ai.Provider.

### `Option`

```go
type Option func(*Client)
```

Option configures a Client at construction time.

### `WithAPIKey`

```go
func WithAPIKey(key string) Option
```

WithAPIKey sets the OpenAI API key (required).

### `WithBaseURL`

```go
func WithBaseURL(url string) Option
```

WithBaseURL overrides the default API root (useful for Azure OpenAI or
self-hosted gateways that mirror the OpenAI shape).

### `WithModel`

```go
func WithModel(model string) Option
```

WithModel overrides the default model for every Generate call.

### `WithHTTPClient`

```go
func WithHTTPClient(h *http.Client) Option
```

WithHTTPClient injects a pre-configured *http.Client (e.g. one bound
to a proxy or an HTTP/3 transport).

### `WithHeader`

```go
func WithHeader(key, value string) Option
```

WithHeader adds an extra request header (merged on top of the
defaults set by oai.Chat).

### `WithUserAgent`

```go
func WithUserAgent(ua string) Option
```

WithUserAgent overrides the default User-Agent header.

### `WithTimeout`

```go
func WithTimeout(d time.Duration) Option
```

WithTimeout sets the HTTP timeout used when no custom *http.Client was
supplied. Zero keeps the default (no timeout — rely on ctx).

### `New`

```go
func New(opts ...Option) *Client
```

New constructs an OpenAI Client.

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

Generate sends a chat-completions request to OpenAI.

### `Generator`

```go
func Generator(opts ...Option) ai.Generator
```

Generator is a convenience that builds a Client from opts and adapts
it to the simple ai.Generator signature. It is the one-liner for
callers that only want `func(ctx, prompt) (string, error)`.
