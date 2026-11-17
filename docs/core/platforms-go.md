# `core` — Platforms

_Source: [`platforms.go`](../../ParsRubika/core/platforms.go)_

## Symbols

### `PlatformKind`

```go
type PlatformKind string
```

PlatformKind enumerates the supported Rubika client platforms. Each
platform maps to a distinct identity (app_name, app_version, package,
headers) sent to the Rubika servers.

ParsRubika is a Rubika-only library: every platform below belongs to
the Rubika ecosystem. Non-Rubika platforms are intentionally not
supported.

### `PlatformPWAFull`

```go
const PlatformPWAFull PlatformKind = "PWA"
```

PlatformPWAFull is the default PWA / web platform.

### `PlatformAndroidFull`

```go
const PlatformAndroidFull PlatformKind = "Android"
```

PlatformAndroidFull is the Android app platform.

### `PlatformRubX`

```go
const PlatformRubX PlatformKind = "RubX"
```

PlatformRubX is the RubX Android variant.

### `PlatformRubikids`

```go
const PlatformRubikids PlatformKind = "Rubikids"
```

PlatformRubikids is the Rubikids Android variant.

### `PlatformRubino`

```go
const PlatformRubino PlatformKind = "Rubino"
```

PlatformRubino is the Rubino Android variant.

### `AllPlatforms`

```go
var AllPlatforms = []PlatformKind{
```

AllPlatforms lists every supported Rubika platform, in priority order.

### `IdentityFor`

```go
func IdentityFor(p PlatformKind, lang string) Identity
```

IdentityFor returns the per-platform Identity. Pass an empty
language code to use the default ("fa").

### `UserAgentFor`

```go
func UserAgentFor(p PlatformKind) string
```

UserAgentFor returns the User-Agent string for the given platform.
Android platforms use `okhttp/3.12.1`; everything else uses the
default ParsRubika User-Agent.

### `DCDiscoveryURLFor`

```go
func DCDiscoveryURLFor(p PlatformKind) string
```

DCDiscoveryURL returns the Rubika DC-discovery endpoint. All Rubika
platforms share the same endpoint.

### `StoragePrefixFor`

```go
func StoragePrefixFor(p PlatformKind) string
```

StoragePrefixFor returns the file-server URL prefix for the given
platform. All Rubika platforms use the `messenger` prefix.
