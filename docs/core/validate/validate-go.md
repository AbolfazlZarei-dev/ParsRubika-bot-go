# `validate` — Package validate provides input-validation helpers for Rubika

_Source: [`validate.go`](../../../ParsRubika/core/validate/validate.go)_

## Symbols

### `ChatID`

```go
func ChatID(id string) bool { return chatIDRE.MatchString(id) }
```

ChatID reports whether `id` is a valid Rubika chat ID.

### `UserGUID`

```go
func UserGUID(g string) bool { return userGUIDRE.MatchString(g) }
```

UserGUID reports whether `g` is a valid user GUID.

### `BotGUID`

```go
func BotGUID(g string) bool { return botGUIDRE.MatchString(g) }
```

BotGUID reports whether `g` is a valid bot GUID.

### `Username`

```go
func Username(u string) bool
```

Username reports whether `u` is a valid Rubika username.

### `Token`

```go
func Token(t string) bool { return tokenRE.MatchString(t) }
```

Token reports whether `t` is a valid Rubika bot token.

### `PhoneNumber`

```go
func PhoneNumber(p string) string
```

PhoneNumber normalises a Persian/Iranian phone number to the
international format (98...). Returns "" on invalid input.

The order of operations matters: the "+" prefix must be stripped
BEFORE the "00"/"0" checks, otherwise inputs like "+098..." are
mis-processed (the "0" check fires on "+0..." and converts the local
"0" to "98", leaving "+98..." which then has its "+" stripped to
"98..." but the leading local "0" was never removed).

Supported input forms:

09981530922 → 989981530922 (local, leading 0)
9981530922 → 989981530922 (bare mobile, no prefix)
+989981530922 → 989981530922 (international with +)
00989981530922 → 989981530922 (international with 00)
989981530922 → 989981530922 (already normalized)

### `MessageText`

```go
func MessageText(t string) bool { return strings.TrimSpace(t) != "" }
```

MessageText reports whether `t` is a non-empty message text.
