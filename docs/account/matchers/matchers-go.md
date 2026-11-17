# `matchers` — Package matchers provides the composable matchers used by the

_Source: [`matchers.go`](../../../ParsRubika/account/matchers/matchers.go)_

## Symbols

### `Matcher`

```go
type Matcher interface
```

Matcher is the contract implemented by every account-side matcher.
It is identical to the account.Matcher interface; it is re-declared
here to avoid an import cycle.
