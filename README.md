# ParsRubika

<p align="center">
  <img src="https://img.parsrubika.ir/banner.png" alt="ParsRubika">
</p>

> **A blazing-fast, modular Go library for the Rubika messaging platform.**

| | |
|---|---|
| **Module**  | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4` |
| **Author**  | Abolfazl Zarei |
| **Version** | 4.0.0 (codename: Nova) |
| **License** | LGPL v3 |
| **Go**      | ≥ 1.20 |

ParsRubika is a from-scratch, production-grade Go implementation of the
Rubika API. It is built around five independent, fully-modular clients
that share one fast, secure foundation. Every line of code was written
specifically for ParsRubika — it is a standalone project with no
dependency on any other Rubika library. The core module builds with
**Go 1.20 and newer** and has only four small dependencies.

---

## Why ParsRubika?

- ⚡ **Fast**: HTTP/2 transport with connection pooling, retries with
  jittered exponential backoff, per-chat ordered parallel dispatch.
  Optional HTTP/3 (QUIC) module.
- 🔒 **Secure**: AES-256-CBC + RSA-PKCS1v15-SHA256, constant-time key
  operations, secure random generation. Encrypted `StringSession`
  (AES-256-GCM + PBKDF2-SHA256, 200k iterations).
- 🧩 **Modular**: 40+ importable packages — every concern (crypto,
  transport, storage, formatting, events, rules, plugins) lives in
  its own package.
- 🎯 **Simple API**: `bot.New(token).Start(ctx)` — that's it.
- 📦 **Pluggable storage**: SQLite, file, in-memory, or string session
  (with optional AES passphrase). Cross-backend converters.
- 🔌 **Plugin framework**: register, enable, disable, reload with
  dependency resolution, cycle detection, rollback on failure, shared
  services and automatic clean-up of everything a plugin registered.
- 🧱 **16 ready-made plugins**: access control, rate limiting, anti-spam,
  anti-link, forced join, help, welcome, broadcast, moderation, stats,
  analytics, i18n, feature flags, key/value store, …
- 🎨 **Logger**: optional coloured, levelled, structured logger with
  JSON mode and rotating files.
- 🗓️ **Scheduler + cron**: one-shot, recurring, and cron-style
  scheduling with a fast field-by-field `Next()` algorithm.
- 📺 **TV client**: a standalone first-class module (just like the
  Rubino, Bot and Account clients), a standalone first-class module.
- 🎨 **Rubino client** as a separate first-class module.
- 🤖 **AI helper**: a multi-file `ai` package with a `Provider`
  interface and **8 ready-made provider packages** — GapGPT, OpenAI,
  DeepSeek, Anthropic (Claude), z.ai (GLM), Google Gemini, Groq and
  OpenRouter — each in its own importable package. No model catalogue
  is hard-coded: every provider declares its own default and accepts
  any model id you pass, so new models work the moment their vendor
  ships them.
- 🧪 **Tested**: unit tests for every package, including the bot core,
  plugins and logger; the suite passes under the race detector.

---

## Architecture

```
(repo root)                # go.mod lives here — the repo root IS the Go module
├── core/                   # Version, author, constants, typed errors, platforms
│   └── errors/             # Sentinel + typed errors
├── security/               # AES, RSA, auth decode, key generation, padding, random
├── text/                   # Markdown / HTML → metadata formatter (UTF-16 aware)
├── enums/                  # ParseMode, ChatType, ChatAction, MediaType, ...
├── storage/                # Store interface + SQLite / Memory / File / StringSession
├── events/                 # Dynamic event wrapper (recursive key lookup)
├── transport/              # HTTP/2 + WebSocket + DC discovery + retry + ratelimit
├── files/                  # Chunked uploader / downloader with progress
├── account/                # User / account client
│   ├── matchers/           # Text, Command, Chat, Media, Author + And/Or/Not + Operator
│   ├── listeners/          # Messages, Chats, Activities, Notifications, Removals
│   └── methods_*.go        # auth, users, chats, messages, groups, channels, ...
├── bot/                    # Bot API client
│   ├── models/             # Chat, Message, Update, Inline, Keypad, Button, ...
│   ├── rules/              # Text, Command, ChatType, ChatID, Media, Button, ...
│   ├── widgets/            # KeypadBuilder / InlineBuilder DSL
│   └── api_*.go            # Messaging, Chats, Members, Files, Keypads, Webhooks
├── rubino/                 # Rubino social-network client
├── tv/                     # Rubika TV / VOD client (standalone, first-class)
├── rubiai/                 # RubiAI assistant client (sessions, chat history, SSE streaming)
├── extras/ai/              # AI plugin (Provider interface, registry, /ask handler)
│   └── providers/          # 8 provider packages, each importable on its own:
│       ├── openai/         #   OpenAI Chat Completions
│       ├── deepseek/       #   DeepSeek (OpenAI-compatible)
│       ├── groq/           #   Groq LPU (OpenAI-compatible)
│       ├── zai/            #   z.ai / Zhipu GLM (OpenAI-compatible)
│       ├── openrouter/     #   OpenRouter (multi-vendor router)
│       ├── anthropic/      #   Anthropic Messages API (Claude)
│       ├── gemini/         #   Google Gemini generateContent
│       └── gapgpt/         #   GapGPT Responses API (Iran-hosted gateway)
├── netutil/                # SSRF-safe URL fetching (private-IP blocking + body caps)
├── logger/                 # Optional coloured / JSON / file logger (no dependencies)
├── plugins/                # Plugin framework: Host, Manager, dependencies, services
│   ├── plugintest/         # Test harness: a bot wired to a fake Rubika API
│   └── builtin/            # access, ratelimit, antispam, antilink, forcejoin, help,
│                           # welcome, users, broadcast, moderation, stats, analytics,
│                           # requestlog, i18n, flags, kv
├── extras/                 # Scheduler, cron, AI helper, forms, QR, tasks, voicekit, handy
└── examples/               # echo, pro_bot, advanced_bot, quick_start, reminder, panel, ...
```

ParsRubika is a **Rubika-only** library. The five supported platforms
(PWA, Android, RubX, Rubikids, Rubino) all belong to the Rubika
ecosystem. Non-Rubika platforms are intentionally not supported.

---

## Installation

```bash
go get github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4
```

The core module needs **Go 1.20+**. SQLite sessions work out of the box
through the in-module `storage/sqlitedriver`; no separate `go get` is
required.

---

## Quick start — Everything in one import

Everything lives in the library itself: the root package builds all five
clients (account, bot, rubino, tv, rubiai) in one call, and three
one-line helpers cover the most common jobs. No long setup code.

```go
package main

import pr "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4"

func main() {
    // One line, one job — the dual (hangi) profile on a channel:
    pr.Hangi("session", "+989123456789",
        "c0DBedk0525ec07fb54a35b01baf2303", "1.jpg", "2.jpg")
}
```

`parsrubika.New` builds the `App`:

```go
app, err := pr.New(&pr.Config{
    Phone:   "+989123456789", // account client (SMS code is asked once)
    Session: "session",       // login is reused from <session>.bin
    // Passphrase: "...",     // optional AES-256-GCM session encryption
    BotToken: "YOUR_BOT_TOKEN", // optional bot client
})
defer app.Close()

app.Account // full account client: messages, groups, files, avatars...
app.Bot     // bot API client (nil when no token)
app.Rubino  // Rubino (auth is copied from the account after Login)
app.TV      // Rubika TV / VOD
app.Rubiai  // RubiAI (auth is copied from the account after Login)

// one-shot work          -> app.Login(ctx)
// stay online + updates  -> app.Run(ctx)   (blocks)
```

And the one-liners:

```go
// send one text message (signs in, sends, closes):
pr.SendText("session", "+989123456789", "gAAAA...targetGUID", "سلام")

// the dual (hangi) profile — پروفایل هنگی — two pictures, ONE call:
pr.Hangi("session", "+989123456789", "", "1.jpg", "2.jpg")

// a running bot in one line:
pr.QuickBot("TOKEN", func(c *bot.Ctx) { c.Reply("سلام!") })
```

---


## Quick start — Bot API

```go
package main

import (
    "context"
    "os"
    "os/signal"

    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
)

func main() {
    b := bot.New("YOUR_BOT_TOKEN")

    b.Command("start", func(c *bot.Ctx) error {
        _, err := c.Reply("Hello! 👋")
        return err
    })
    b.Message(func(c *bot.Ctx) error { // echo everything else
        _, err := c.Reply(c.Text())
        return err
    })

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()
    _ = b.Start(ctx) // returns when ctx is cancelled
}
```

`*bot.Ctx` carries the update and offers `Reply`, `Send`, `Edit`,
`Delete`, `Forward`, `CopyTo`, `Args`, `Command`, `ButtonID`, `Chat`,
`Ban` and `Set`/`Get` for passing values from middleware to handlers.
The lower-level `b.On(rules, func(ctx, b, upd) error)` stays available.

---

## Quick start — Account client (interactive)

```go
package main

import (
    "context"
    "log"

    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account"
    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account/listeners"
    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account/matchers"
    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

func main() {
    acc, err := account.New(
        account.WithPhoneNumber("+989123456789"),
        account.WithSQLiteSession("mysession"),
    )
    if err != nil { log.Fatal(err) }
    defer acc.Close()

    account.AddListener(acc, func(ctx context.Context, e *events.Event) error {
        _, err := acc.SendMessage(ctx, e.ObjectGUID(), "pong! 🏓")
        return err
    }, listeners.NewMessages(matchers.NewIsText("MessageUpdates")))

    log.Fatal(acc.Run(context.Background(), ""))
}
```

## Quick start — Account client (headless / server)

For servers, Docker, or systemd where stdin is not available, supply a
`CodeProvider` callback. The sign-in flow calls it to obtain the SMS
code or 2FA password instead of prompting on the terminal:

```go
acc, err := account.New(
    account.WithPhoneNumber("+989123456789"),
    account.WithCodeProvider(func(ctx context.Context, hint string) (string, error) {
        // hint is "Code:" or "Password [hint] > "
        // Read the code from your HTTP UI, a Redis queue, a Telegram
        // control bot, etc.
        return myUI.PromptUser(ctx, hint)
    }),
    account.WithSQLiteSession("mysession"),
)
```

### Account "hang" / dual profile

The account client also ships `HangProfile`, the well-known Rubika
"dual profile" trick: it fires two `updateProfile` requests
back-to-back so the two writes race on the server, which can leave the
persisted profile in an inconsistent (stacked / dual) state.

```go
res, err := acc.HangProfile(ctx, "Foo", "Bar", "bio-text", "2000-01-01")
// res.First  / res.FirstErr   — outcome of shot 1
// res.Second / res.SecondErr  — outcome of shot 2
```

Both shots use the **same** payload; the "hang" effect comes from the
race itself, not from conflicting values. The method honours context
cancellation and never panics if one shot fails — its error lands on the
returned result.

### The dual-avatar hangi profile (پروفایل هنگی)

`HangiProfile` sets a **dual profile picture**: two different images on
one object at the same time — main photo of the first bind is photo1
with photo2 as its thumbnail, and the second bind fires the swapped
pair, so the object shows two different pictures instead of one stacked
photo. It targets the signed-in account, or any group/channel GUID you
administer:

```go
// objectGUID "" = the account itself; pass a group/channel GUID for chats.
res, err := acc.HangiProfile(ctx, "c0DBedk0525ec07fb54a35b01baf2303",
    "1.jpg", "2.jpg",
    account.WithHangiProgress(func(msg string) { fmt.Println(msg) }))
// res.AvatarIDs  — avatars created by this run
// res.Deleted    — leftover avatars removed by the clean-up pass
// res.FinalCount — what the server lists at the end
```

One application is the correct usage: the default is a single round of
two simultaneous `uploadAvatar` binds, and the clean-up pass deletes
every avatar the run did not create — so rerunning never piles up
photos and the object ends with **exactly** the two new pictures.
Options: `WithHangiRounds`, `WithHangiOffset`, `WithHangiGap`,
`WithHangiClean`, `WithHangiKeepOld`, `WithHangiProgress`. Or skip the
client plumbing entirely with the one-line `parsrubika.Hangi(...)` (see
[Quick start — Everything in one import](#quick-start--everything-in-one-import)).
A runnable example lives in `examples/hangprofile`.

## Quick start — Bot with ready-made plugins

```go
b := bot.New("YOUR_BOT_TOKEN", bot.WithLogger(logger.New()))

m := plugins.New(b)
_ = builtin.RegisterAll(m, map[string]plugins.Config{
    "access":    {"admins": []string{"b0YourUserID"}},
    "ratelimit": {"rate": 2.0, "burst": 6},
    "forcejoin": {"channels": []string{"c0Channel"}, "links": []string{"rubika.ir/mychannel"}},
})
m.Install() // plugins start and stop together with the bot

b.Command("ping", func(c *bot.Ctx) error { _, err := c.Reply("pong"); return err })
_ = b.Start(ctx)
```

See [Plugins](#plugins) for the full list and for writing your own, and
`examples/pro_bot` for a complete program.

---

## Quick start — Rubino client

```go
package main

import (
    "context"
    "log"

    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/rubino"
)

func main() {
    s := rubino.New(rubino.WithAuth("YOUR_AUTH"))
    defer s.Close()

    res, err := s.GetProfileList(context.Background(), 50, "latest", "")
    if err != nil { log.Fatal(err) }
    log.Println(res.JSON(""))
}
```

## Quick start — TV client

The TV client is a **standalone, first-class module** (just like the
Rubino, Bot and Account clients) and lives at `parsrubika/v4/tv` —
not under `extras/`.

```go
package main

import (
    "context"
    "log"

    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/tv"
)

func main() {
    t := tv.New(tv.WithAuth("YOUR_AUTH"))
    defer t.Close()

    res, err := t.Search(context.Background(), "news", "0")
    if err != nil { log.Fatal(err) }
    log.Println(res.JSON(""))
}
```

## Quick start — RubiAI client

The RubiAI client talks to the RubiAI assistant API and supports SSE
streaming for token-by-token output. Reuse the account auth key, or
supply a dedicated one.

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/rubiai"
)

func main() {
    ai := rubiai.New(rubiai.WithAuth("YOUR_AUTH"))
    defer ai.Close()

    // Blocking call — assembles the full reply.
    res, err := ai.SendMessage(context.Background(), rubiai.SendMessageOptions{
        Message: "Summarise the Rubika API in one line.",
    })
    if err != nil { log.Fatal(err) }
    fmt.Println(res.TextContent)

    // Streaming call — invoked per SSE event.
    _, _ = ai.SendMessageStream(context.Background(),
        rubiai.SendMessageOptions{Message: "Tell me a story."},
        func(event string, data map[string]any) {
            if event == "updateMessage" {
                fmt.Print(".")
            }
        })
    fmt.Println()
}
```

## AI providers

The `extras/ai` package is split across several small files
(`provider.go`, `generator.go`, `plugin.go`, `prompt.go`, `doc.go`) and
exposes a single `Provider` interface. Every supported LLM service is
its **own package** under `extras/ai/providers/<name>` — ParsRubika
never ships a model catalogue, so each provider declares its own default
model and accepts any model id you pass through `ai.Request.Model`.

| Provider package | Service | Wire format | Default model |
|---|---|---|---|
| `providers/openai` | OpenAI | Chat Completions | `gpt-4o-mini` |
| `providers/deepseek` | DeepSeek | OpenAI-compatible | `deepseek-chat` |
| `providers/groq` | Groq (LPU) | OpenAI-compatible | `llama-3.3-70b-versatile` |
| `providers/zai` | z.ai / Zhipu GLM | OpenAI-compatible | `glm-4.6` |
| `providers/openrouter` | OpenRouter | OpenAI-compatible router | `openai/gpt-4o-mini` |
| `providers/anthropic` | Anthropic (Claude) | Messages API | `claude-opus-4-1` |
| `providers/gemini` | Google Gemini | `generateContent` | `gemini-2.0-flash` |
| `providers/gapgpt` | GapGPT (Iran-hosted) | Responses API | `gapgpt-4o-mini` |

All eight implement the same interface and accept the same options
(`WithAPIKey`, `WithModel`, `WithHTTPClient`, `WithTimeout`, …), so they
are drop-in replacements for each other.

```go
package main

import (
    "context"
    "os"
    "os/signal"

    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai"
    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/openai"
)

func main() {
    b := bot.New(os.Getenv("BOT_TOKEN"))

    // One provider per service — swap the import + constructor to
    // change the backend. The model is overridable per call.
    p := openai.New(openai.WithAPIKey(os.Getenv("OPENAI_API_KEY")))

    // Drive a provider directly ...
    res, _ := p.Generate(context.Background(), ai.Request{
        Model:  "gpt-4o",      // optional; falls back to the default
        Prompt: "Summarise the Rubika API in one line.",
    })
    _ = res.Text

    // ...or register it with the AI plugin so /ask answers in chat.
    ai.Setup(b, ai.Config{Provider: p})

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()
    _ = b.Start(ctx)
}
```

## Security

ParsRubika ships with defence-in-depth for the surfaces it exposes on
the user's behalf:

- **Session at-rest encryption** — `FileStore` and `SQLiteStore` can
  seal the `auth` key and RSA `private_key` with AES-256-GCM
  (PBKDF2-SHA256, 200k iters). Use `WithFileSessionPassphrase` /
  `WithSQLiteSessionPassphrase`. Legacy plain-text sessions auto-migrate
  on the next Save.
- **SSRF protection** — the `netutil` package blocks requests to
  private/loopback/link-local/cloud-metadata IPs at both pre-flight and
  the dial layer (defeating DNS rebinding), and caps response bodies.
  `UploadFromURL`, `FetchJSON` and `ScrapeChannelPostLink` use it.
- **Decompression-bomb caps** — `StringSession` and `bot.ImportTokenString`
  cap decrypted/decompressed payloads (64 KiB).
- **Debug-log redaction** — `WithDebug(true)` logs a non-sensitive
  summary and masks `auth`, `private_key`, `sign`, `data_enc`, `file_id`,
  `access_hash_*`, `text`, `phone_number`, `bio`, `password`.
- **RSA 2048** — new sign-ins generate 2048-bit keys (matching the
  official client); existing 1024-bit sessions keep working.

See `docs/security/` for the session encryption reference.

## Quick start — Scheduler + cron

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/cron"
    "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/scheduler"
)

func main() {
    s := scheduler.New(context.Background())
    s.After(5*time.Second, func(ctx context.Context) {
        fmt.Println("5 seconds elapsed!")
    })
    s.Every(1*time.Minute, func(ctx context.Context) {
        fmt.Println("Tick.")
    })

    c := cron.New(context.Background())
    _ = c.Add("@daily", func(ctx context.Context) {
        fmt.Println("Daily report goes here.")
    })
    c.Start()

    select {}
}
```

---

## Sending media

The Bot API exposes typed helpers for every Rubika media type. Each one
is a thin wrapper around `SendFile` with the correct `FileType`:

| Method | FileType | Use |
|---|---|---|
| `SendImage` | `Image` | Photos (server generates thumbnail) |
| `SendVideo` | `Video` | Video files |
| `SendGif` | `Gif` | Animated GIFs |
| `SendAudio` / `SendMusic` | `Music` | Audio with a music player |
| `SendVoice` | `Voice` | Voice notes (`.ogg`) |
| `SendFile` | `File` | Generic file (default) |
| `SendSticker` | — | Sticker by ID |
| `SendPoll` | — | Polls |
| `SendLocation` | — | Location attachment |
| `SendContact` | — | Contact attachment |

Example — send a photo with a caption:

```go
_, err := b.SendImage(ctx, chatID, "/path/to/photo.jpg", "", "Look at this! 📸")
```

Example — send a video by a previously-uploaded file ID:

```go
_, err := b.SendVideo(ctx, chatID, "", "fileID-here", "My video", "")
```

---

## Plugins

A plugin is a self-contained feature that plugs into the bot through a
`Host`. Everything it registers (handlers, middleware, hooks, timers,
services) is tracked and removed when the plugin is disabled.

```go
type Hello struct{ plugins.Base }

func New() *Hello {
    return &Hello{plugins.NewBase(plugins.Meta{
        Name: "hello", Version: "1.0.0", Description: "Says hello",
        Requires: []string{"kv"},          // enabled first, disabled last
    })}
}

func (p *Hello) Setup(ctx context.Context, h *plugins.Host) error {
    n := h.Config().Int("times", 1)        // defaults merged with user config
    h.Help("/hello", "Say hello")          // shows up in /help
    h.Command("hello", func(c *bot.Ctx) error {
        _, err := c.Reply(strings.Repeat("Hello! ", n))
        return err
    })
    h.Every(time.Minute, func(ctx context.Context) { h.Debugf("alive") })
    return nil
}
```

Optional interfaces: `Defaulter` (default config), `Validator`
(reject bad config before start), `Teardowner` (release your own
resources). `Manager` offers `Register`, `Enable`, `Disable`, `Reload`,
`EnableAll`/`DisableAll`, `Install`, `List`, `Configure`, services
(`h.Provide`, `plugins.Service[T](m, key)`) and loading configuration
from JSON (`plugins.LoadConfigFile`, `WithConfigs`, manifests in a
directory via `ApplyManifests`).

**Built-in plugins** (`plugins/builtin`, register all with `builtin.RegisterAll`):

| Plugin | What it does |
|---|---|
| `access` | admins, allow-list, block-list; publishes the `access` service |
| `ratelimit` | per-user token bucket |
| `antispam` | flood / repeated-text detection with temporary mutes |
| `antilink` | removes links of non-admins in groups |
| `forcejoin` | forced channel subscription (fails open on API errors) |
| `help` | automatic `/help` from the plugins' documented commands |
| `welcome` | greeting on `StartedBot` and `/start` |
| `users` | remembers chats (persistent), skips users who stopped the bot |
| `broadcast` | admin `/broadcast <text>` to all known chats |
| `moderation` | `/ban`, `/unban`, `/warn`, `/warnings`, `/resetwarns` |
| `stats` | admin `/stats`: uptime, counters, memory |
| `analytics` | counters for update types and commands |
| `requestlog` | per-update timing logs, slow-handler warnings |
| `i18n` | per-user language, `/lang`, translator service |
| `flags` | feature flags with per-user overrides |
| `kv` | key/value store with TTL and JSON persistence |

Test your own plugins with `plugins/plugintest`: it starts a bot against
an in-process fake Rubika API, so `env.Say("b0user", "u1", "/hello")`
and `env.Sent()` need no network.

---

## Events, handlers and rules

Handlers are registered with `On` (low level) or `Handle`, `Command`,
`Text`, `Regex`, `Button`, `Message` and `Fallback` (with `*bot.Ctx`).
Each handler is gated by a slice of `Rule`s — **all** rules must pass.
The first matching handler wins; a handler can return `bot.ErrContinue`
to let the next matching handler run as well.

- Order: higher **priority** first, then registration order. Plugin
  handlers use `bot.PriorityPlugin` so an application catch-all cannot
  swallow their commands; `Fallback` always runs last.
- `Off(id)` removes a handler, `RemoveMiddleware(id)` a middleware.
- **Middleware** wraps the handlers: `func(ctx, b, upd, next) error`.
  Not calling `next()` stops the update; calling it and inspecting the
  result lets you time, log or recover handlers.
- A panic in a handler or middleware never kills the bot: it is
  recovered, counted and sent to `bot.WithErrorHandler` (or logged).
- Duplicate updates (same message ID and time) are dropped.

```go
// Only in private chats, only "/ask <question>".
b.Handle([]bot.Rule{rules.PrivateChat{}, bot.CommandRule("ask")},
    func(c *bot.Ctx) error {
        _, err := c.Reply("You asked: " + c.ArgsText())
        return err
    })
```

Rules live in `bot/rules` (about 100): chat type, text and regex,
command (exact — `/warnings` does not match `/warn`), sender and chat
sets, media, polls, locations, forwards, replies, buttons, … plus the
combinators `And`, `Or`, `Not` and `rules.Func(func(ctx, upd) bool)`.

The Account client uses an analogous `Matcher` system
(`account/matchers`) with `And` / `Or` / `Not` combinators and an
operator system (`==`, `!=`, `<`, `<=`, `>`, `>=`).

---

## Concurrency & dispatch

- **Bot**: updates of the **same chat are processed strictly in
  order**, different chats run in parallel (default 8 at once, change it
  with `bot.WithConcurrency(n)`; `1` means fully sequential). On shutdown
  the bot waits (up to `ShutdownTimeout`) for in-flight handlers.
- **Polling** never competes with message sending for the rate limiter,
  backs off exponentially (with jitter, capped) after errors and honours
  `Retry-After`. The first batch of old updates is skipped unless
  `bot.WithSkipBacklog(false)`; the offset is stored atomically on disk.
- **Webhook**: `b.Run(ctx, "https://your.host", ":8080")` registers all
  endpoints, limits the body size and can require a shared secret
  (`bot.WithWebhookSecret`). `b.WebhookHandler()` mounts the bot on your
  own `http.Server`.
- **Account**: listeners run **concurrently** under a bounded semaphore
  (64 per update). Use `account.WithSequentialListeners(true)` for
  strict in-order processing.
- **Tasks pool** (`extras/tasks`): a bounded worker pool with a
  send-on-closed-channel-safe `Submit` and a `Wait()` that does NOT
  stop the pool.
- **Plugins**: enabling and disabling is serialised; a plugin that
  fails to start is rolled back completely.

---

## Storage & sessions

ParsRubika ships four storage backends, all implementing the same
`Store` interface:

| Backend | Use | Construction |
|---|---|---|
| `SQLiteStore` | Persistent session | `account.WithSQLiteSession("name")` |
| `FileStore` | JSON file on disk (base64-encoded) | `account.WithFileSession("name")` |
| `MemoryStore` | Ephemeral, in-process | `storage.NewMemoryStore()` |
| `StringSession` | Portable, AES-256-GCM encrypted string | `storage.NewStringSession(passphrase)` |

The `StringSession` can be exported from one backend and imported into
another, so you can start a session on a dev machine and move it to a
server.

---

## Error handling

ParsRubika errors are typed via `core/errors`:

```go
_, err := b.SendMessage(ctx, chatID, "hi")
if err != nil {
    var apiErr *coreerrors.Error
    if errors.As(err, &apiErr) {
        switch apiErr.Kind {
        case coreerrors.KindNetwork:  // retryable transport error
        case coreerrors.KindAPI:      // server returned non-OK status
        case coreerrors.KindAuth:     // auth/key problem
        case coreerrors.KindValidation: // bad input
        }
        log.Printf("status=%q message=%s", apiErr.Status, apiErr.Message)
    }
}
```

Sentinel errors (`errors.Is`) and status checks
(`coreerrors.IsStatus(err, "NotRegistered")`) are also supported.

---

## Logging

Every component accepts any logger with a `Printf` method (the standard
`*log.Logger` works). For readable, levelled, coloured output use the
bundled `logger` package — it has no dependencies:

```go
log := logger.New(
    logger.WithLevel(logger.LevelDebug),
    logger.WithName("mybot"),
    logger.WithFile("logs/bot.log", logger.FormatText, 10, 3), // rotating, plain text
)
defer log.Close()

b := bot.New(token, bot.WithLogger(log))
acc := account.New(account.WithLogger(log))

log.Info("ready", logger.F("user", 42))
log.Warn("slow upstream", logger.F("took", 2*time.Second), logger.Err(err))
```

- Levels `Trace … Error`; set `PARSRUBIKA_LOG_LEVEL=debug` or call `SetLevel` at runtime.
- Colours are automatic on terminals; `NO_COLOR` disables and `FORCE_COLOR` forces them (Windows 10+ supported).
- `logger.WithJSON()` / `logger.FormatJSON` give one JSON object per line.
- `log.Named("db")`, `log.With(logger.F(...))` build child loggers; `log.StdLogger(level)` adapts to `*log.Logger`.
- Verbose internals (decrypted payloads, access hashes) are logged at **Debug** only.

Components that implement the optional `Debugf/Infof/Warnf/Errorf`
methods get properly levelled records; plain `Printf` loggers still work.

---

## Available packages

| Package | Purpose |
|---|---|
| `parsrubika` (root) | **Everything in one import**: `New`/`App` builds account+bot+rubino+tv+rubiai together, plus the one-liners `SendText`, `Hangi`, `QuickBot` |
| `core` | Version, author, constants, 5 Rubika platforms, 5 first-class client registry (bot/account/rubino/tv/rubiai) |
| `core/errors` | Typed errors with Kind + Status, sentinel errors, `errors.Is/As` support |
| `security` | AES-256-CBC, RSA-PKCS1v15, RSA-OAEP, key gen, decode_auth, key derivation |
| `text` | Markdown / HTML → metadata formatter with UTF-16 offsets |
| `enums` | ParseMode, ChatType, ChatAction, MediaType, PollType, Reaction, Access, ReportReason |
| `storage` | Store interface, SQLiteStore, MemoryStore, FileStore, **StringSession** (encrypted) |
| `events` | Dynamic `Event` wrapper with `Find`, `Str`, `Int`, `Bool`, property accessors |
| `transport` | HTTP/2 + WebSocket + DC discovery + retry + ratelimit + headers |
| `files` | Chunked uploader (bounded ERROR_TRY_AGAIN restarts) + chunked downloader + progress |
| `account` | User/account client: lifecycle, dispatch, encrypted Request builder, **HangProfile** (dual-profile trick), **HangiProfile** (dual-avatar پروفایل هنگی + tolerant image upload) |
| `account/matchers` | 25+ matchers with operator system (==, !=, <, <=, >, >=) + And/Or/Not |
| `account/listeners` | Messages, Chats, Activities, Notifications, Removals |
| `bot` | Bot API client: lifecycle, polling, webhook, request |
| `bot/models` | Chat, Message, Update, Inline, Keypad, Button, Poll, Location, File, ... |
| `bot/rules` | ~100 rules and combinators: Text, Command, Group, Channel, Photo, Button, Func, And/Or/Not, ... |
| `bot/widgets` | KeypadBuilder + InlineBuilder DSL with all 20+ button types |
| `rubino` | Rubino social-network client (full method surface) |
| `tv` | Rubika TV / VOD client — standalone first-class module (browse, search, like, wish-list, cast, series) |
| `logger` | Coloured / JSON / rotating-file logger, no dependencies |
| `plugins` | Plugin framework: Host, Manager, dependencies, services, config, discovery |
| `plugins/builtin/*` | 16 ready-made plugins (see [Plugins](#plugins)) |
| `plugins/plugintest` | Test harness: a bot against an in-process fake Rubika API |
| `extras/scheduler` | One-shot and recurring task scheduler + SendLater/SendRecurring helpers |
| `extras/cron` | 5-field cron syntax + @yearly/@monthly/@weekly/@daily/@hourly macros |
| `extras/ai` | AI plugin helper: `Provider` interface + registry + `/ask` handler (multi-file) |
| `extras/ai/providers/*` | 8 provider packages: `openai`, `deepseek`, `groq`, `zai`, `openrouter`, `anthropic`, `gemini`, `gapgpt` |
| `extras/handy` | Echo, PingReply, CommandHelp, MentionChat, MentionUserID, CopyMessage |
| `extras/tasks` | Bounded goroutine pool + concurrency-limited distributor |

---

## Documentation

Per-package API reference is generated under `docs/`. Each package has an
`index.md` overview plus one file per source file listing every exported
symbol with its signature and documentation. Start at
[`docs/index.md`](docs/index.md).

The reference is regenerated for every release straight from the Go
sources, so it always matches the code you have.

---

## Credits & Contact

**ParsRubika** is developed and maintained by **[Abolfazl Zarei](https://abolfazlzarei.sbs)**.

For documentation, official announcements, support, and project updates, you can find us on the following platforms:

| Platform | Official Link |
|:---:|:---|
| 🌐 **Website** | [abolfazlzarei.sbs](https://abolfazlzarei.sbs) |
| 📚 **Documentation** | [parsrubika.ir](https://parsrubika.ir) |
| 💬 **Rubika** | [@NinjaCode](https://rubika.ir/NinjaCode) |
| ✈️ **Telegram** | [@Abolfazl_PGR](https://t.me/Abolfazl_PGR) |
| 🚀 **ParsRubika Channel** | [@parsrubika_lib](https://rubika.ir/parsrubika_lib) |

Stay connected to receive the latest releases, development updates, announcements, and improvements.

---

## License

ParsRubika is licensed under the **GNU Lesser General Public License v3.0 (LGPL-3.0)**.

For more information, see the [LICENSE](LICENSE) file.

---

<p align="center">
  <strong>ParsRubika v4.0.0</strong>
  <br><br>
  Developed and maintained with dedication by
  <br>
  <a href="https://abolfazlzarei.sbs"><strong>Abolfazl Zarei</strong></a>
  <br><br>
  <a href="https://parsrubika.ir">Documentation</a>
  &nbsp;·&nbsp;
  <a href="https://rubika.ir/parsrubika_lib">Official Channel</a>
  &nbsp;·&nbsp;
  <a href="https://rubika.ir/NinjaCode">Rubika</a>
  &nbsp;·&nbsp;
  <a href="https://t.me/Abolfazl_PGR">Telegram</a>
  &nbsp;·&nbsp;
  <a href="https://abolfazlzarei.sbs">Website</a>
</p>
