# `bot` — Webhook

_Source: [`webhook.go`](../../ParsRubika/bot/webhook.go)_

## Symbols

### `WebhookHandler`

```go
func (b *Bot) WebhookHandler() http.Handler
```

WebhookHandler returns the http.Handler that serves every webhook
endpoint, so the bot can be mounted on an existing HTTP server.
