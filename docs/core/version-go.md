# `core` — Package core holds the version, author and shared constants for

_Source: [`version.go`](../../ParsRubika/core/version.go)_

## Symbols

### `Author`

```go
const Author = "Abolfazl Zarei"
```

Author is the maintainer of ParsRubika.

### `Version`

```go
const Version = "4.0.0"
```

Version is the current ParsRubika release.

### `Codename`

```go
const Codename = "Nova"
```

Codename is the marketing name for this release line.

### `License`

```go
const License = "LGPL-3.0-only"
```

License is the SPDX-style identifier of the distribution license.

### `UserAgent`

```go
const UserAgent = "ParsRubika/" + Version + " (go; +" + Author + ")"
```

UserAgent is the default User-Agent string used by every ParsRubika
HTTP client. It already includes the version so that server-side
telemetry can distinguish ParsRubika traffic.

### `APIVersion`

```go
const APIVersion = "6"
```

APIVersion is the wire-format API version sent in every request body.

### `DefaultLangCode`

```go
const DefaultLangCode = "fa"
```

DefaultLangCode is the default language code sent to the server.

### `DefaultTimeout`

```go
const DefaultTimeout = 30
```

DefaultTimeout is the default HTTP timeout, in seconds.

### `DefaultMaxRetries`

```go
const DefaultMaxRetries = 5
```

DefaultMaxRetries is the default number of retries for failed requests.

### `DefaultChunkSize`

```go
const DefaultChunkSize = 1 << 20
```

DefaultChunkSize is the default upload/download chunk size (1 MiB).

### `DefaultDownloadChunkSize`

```go
const DefaultDownloadChunkSize = 128 * 1024
```

DefaultDownloadChunkSize is the default download chunk size (128 KiB).

### `WelcomeBanner`

```go
func WelcomeBanner() string
```

WelcomeBanner returns the standard library welcome banner.
