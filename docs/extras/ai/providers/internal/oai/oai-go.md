# `oai` — Package oai is the shared transport for OpenAI-compatible chat

_Source: [`oai.go`](../../../../../../ParsRubika/extras/ai/providers/internal/oai/oai.go)_

## Symbols

### `Config`

```go
type Config struct
```

Config is the per-provider wiring passed to Chat. Every field except
HTTPClient (which defaults to http.DefaultClient) must be set by the
provider package.

### `Chat`

```go
func Chat(ctx context.Context, c Config, name string, req ai.Request) (ai.Response, error)
```

Chat sends an OpenAI-compatible chat-completions request and returns
the assembled reply. The provider's Name is stamped onto the
resulting ai.Response.Provider.

### `ParseInt`

```go
func ParseInt(s string) int { v, _ := strconv.Atoi(strings.TrimSpace(s)); return v }
```

ParseInt is a small helper for providers that read numeric headers.
Kept here so every provider imports oai for one reason only.
