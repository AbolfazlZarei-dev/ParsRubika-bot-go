# `guid` — Package guid provides utility functions for working with Rubika

_Source: [`guid.go`](../../../ParsRubika/core/guid/guid.go)_

## Symbols

### `Kind`

```go
type Kind string
```

Kind is the kind of object a GUID refers to.

### `KindUser`

```go
const KindUser Kind = "User"
```

### `KindBot`

```go
const KindBot Kind = "Bot"
```

### `KindGroup`

```go
const KindGroup Kind = "Group"
```

### `KindChannel`

```go
const KindChannel Kind = "Channel"
```

### `KindService`

```go
const KindService Kind = "Service"
```

### `KindUnknown`

```go
const KindUnknown Kind = "Unknown"
```

### `Detect`

```go
func Detect(g string) Kind
```

Detect infers the kind of object from the GUID prefix.

### `IsUser`

```go
func IsUser(g string) bool { return Detect(g) == KindUser }
```

IsUser reports whether `g` is a user GUID.

### `IsBot`

```go
func IsBot(g string) bool { return Detect(g) == KindBot }
```

IsBot reports whether `g` is a bot GUID.

### `IsGroup`

```go
func IsGroup(g string) bool { return Detect(g) == KindGroup }
```

IsGroup reports whether `g` is a group GUID.

### `IsChannel`

```go
func IsChannel(g string) bool { return Detect(g) == KindChannel }
```

IsChannel reports whether `g` is a channel GUID.

### `IsService`

```go
func IsService(g string) bool { return Detect(g) == KindService }
```

IsService reports whether `g` is a service GUID.

### `IsValid`

```go
func IsValid(g string) bool
```

IsValid reports whether `g` looks like a valid GUID (32 chars
after the 2-char prefix).

### `StripAt`

```go
func StripAt(s string) string { return strings.TrimPrefix(s, "@") }
```

StripAt removes a leading `@` from a username.

### `TrailingSlug`

```go
func TrailingSlug(link string) string
```

TrailingSlug returns the trailing path segment of a URL or the
input itself when it does not contain a slash.
