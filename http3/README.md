# http3 (optional module)

HTTP/3 (QUIC) client for ParsRubika. It is a separate module so the core
library keeps building with Go 1.20; `quic-go` needs Go 1.23.

```bash
cd your-project
go get github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/http3
```

```go
import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/http3"

client := http3.NewHTTP3Client(30*time.Second, false)
b := bot.New(token, bot.WithHTTPClient(client))
```

When working from a checkout of this repository run `go mod tidy` once
inside this directory to resolve the indirect dependencies.
