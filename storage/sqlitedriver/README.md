# sqlitedriver (optional module)

Registers a pure-Go (cgo-free) SQLite driver named `sqlite`, which
`storage.NewSQLiteStore` and `account.WithSQLiteSession` need. It is a
separate module so the core library keeps building with Go 1.20; the
driver itself needs Go 1.21.

```go
import _ "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/storage/sqlitedriver"
```

When working from a checkout of this repository run `go mod tidy` once
inside this directory to resolve the indirect dependencies.
