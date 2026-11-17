// Package bot implements the official Rubika Bot HTTP API client.
//
// It is split across several small files so each concern lives in its
// own unit:
//
// - bot.go : the Bot struct + constructor + options
// - lifecycle.go : Start / Stop / Run (polling + webhook)
// - polling.go : the long-poll loop
// - webhook.go : the HTTP webhook server
// - request.go : the encrypted request dispatcher
// - ratelimit.go : the per-token rate limiter
// - models/ : Chat, Message, Update, Inline, Keypad, Button, ...
// - rules/ : Text, Command, ChatType, ChatID, Media, ...
// - api/ : Messaging, Chats, Members, Files, Keypads, Webhooks
package bot

import (
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/enums"
)

// Author / Version / Codename re-exports.
const (
	Author   = core.Author
	Version  = core.Version
	Codename = core.Codename
)

// BaseURL is the bot API endpoint template. The token is interpolated
// in the final URL: `<BaseURL>/<token>/`.
const BaseURL = "https://botapi.rubika.ir/v3/%s/"

// ParseMode is re-exported here for ergonomic imports.
type ParseMode = enums.ParseMode
