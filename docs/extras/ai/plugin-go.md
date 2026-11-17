# `ai` — Plugin

_Source: [`plugin.go`](../../../ParsRubika/extras/ai/plugin.go)_

## Symbols

### `Config`

```go
type Config struct
```

Config configures the AI plugin installed by Setup.

A generator can be supplied in three ways, tried in this order:
1. Config.Provider — a full Provider implementation (preferred);
2. Config.Generator — a plain Generator callback;
3. the global registry populated via RegisterProvider/RegisterGenerator.

When Config.Provider is set it is also registered globally so later
Setup calls (and other components) reuse it without re-passing it.

### `Setup`

```go
func Setup(b *bot.Bot, cfg Config)
```

Setup registers an AI handler on the bot. The handler fires only for
messages that are a command matching Config.Command (default
`/ask <prompt>`). Using a real command rule — instead of a catch-all
text rule — ensures the AI plugin never shadows other text handlers.

The chosen generator is resolved once at Setup time from
Config.Provider, Config.Generator or the global registry; if none is
available the handler replies with a friendly "not configured"
message instead of panicking.
