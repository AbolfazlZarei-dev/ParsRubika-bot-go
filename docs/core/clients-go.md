# `core` — Clients

_Source: [`clients.go`](../../ParsRubika/core/clients.go)_

## Symbols

### `ClientKind`

```go
type ClientKind string
```

ClientKind enumerates the first-class ParsRubika clients. Each kind
maps to a standalone, fully-modular client package that shares one
fast, secure foundation (core / security / transport / events).

The five clients are registered here as metadata only — core must not
import any of them (that would create an import cycle, since every
client imports core). The Package field carries the import path so
tooling and docs can point at the right module without a hard
dependency.

Mirrors the platform registry in platforms.go: a single, authoritative
place where every first-class client is listed, in priority order.

### `ClientBot`

```go
const ClientBot ClientKind = "bot"
```

ClientBot is the official Rubika Bot HTTP API client.

### `ClientAccount`

```go
const ClientAccount ClientKind = "account"
```

ClientAccount is the Rubika user (account) client — the encrypted
messenger envelope plus the live WebSocket update stream.

### `ClientRubino`

```go
const ClientRubino ClientKind = "rubino"
```

ClientRubino is the Rubino social-network client.

### `ClientTV`

```go
const ClientTV ClientKind = "tv"
```

ClientTV is the Rubika TV / VOD client.

### `ClientRubiAI`

```go
const ClientRubiAI ClientKind = "rubiai"
```

ClientRubiAI is the RubiAI assistant client (chat, sessions, SSE
streaming, speech-to-text).

### `ClientInfo`

```go
type ClientInfo struct
```

ClientInfo is the metadata for one first-class client. It is the
client-side analogue of Identity: name, import path, a one-line
description and the service's base URL. All fields are static strings
so core stays dependency-free.

### `AllClients`

```go
var AllClients = []ClientInfo{
```

AllClients lists every first-class ParsRubika client, in priority
order: the messaging clients first (bot, account), then the
Rubika-ecosystem vertical clients (rubino, tv, rubiai). It mirrors
AllPlatforms so the two registries stay in lock-step.

### `ClientFor`

```go
func ClientFor(k ClientKind) ClientInfo
```

ClientFor returns the ClientInfo for the given kind, or the zero
ClientInfo (with an empty Kind) when the kind is unknown. Callers
should check the returned Kind against the empty string.

### `ClientNames`

```go
func ClientNames() []string
```

ClientNames returns the Name of every first-class client, in
AllClients order. Handy for dashboards, /help output and docs.

### `ClientKinds`

```go
func ClientKinds() []ClientKind
```

ClientKinds returns every ClientKind, in AllClients order.
