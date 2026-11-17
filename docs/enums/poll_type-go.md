# `enums` — Poll type

_Source: [`poll_type.go`](../../ParsRubika/enums/poll_type.go)_

## Symbols

### `PollType`

```go
type PollType string
```

PollType is the type of a poll (Regular vs Quiz).

### `PollTypeRegular`

```go
const PollTypeRegular PollType = "Regular"
```

### `PollTypeQuiz`

```go
const PollTypeQuiz PollType = "Quiz"
```

### `String`

```go
func (p PollType) String() string { return string(p) }
```

String implements fmt.Stringer.
