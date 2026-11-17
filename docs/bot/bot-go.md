# `bot` — Package bot implements the official Rubika Bot HTTP API client

_Source: [`bot.go`](../../ParsRubika/bot/bot.go)_

## Symbols

### `Author`

```go
const Author = core.Author
```

Author / Version / Codename re-exports.

### `Version`

```go
const Version = core.Version
```

Author / Version / Codename re-exports.

### `Codename`

```go
const Codename = core.Codename
```

Author / Version / Codename re-exports.

### `BaseURL`

```go
const BaseURL = "https://botapi.rubika.ir/v3/%s/"
```

BaseURL is the bot API endpoint template. The token is interpolated
in the final URL: `<BaseURL>/<token>/`.
