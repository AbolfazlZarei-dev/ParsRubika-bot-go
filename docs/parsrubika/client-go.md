# `parsrubika` — Everything-in-one client (App) and one-line helpers

_Source: [`client.go`](../../ParsRubika/client.go)_

## Symbols

### `Config`

```go
type Config struct
```

Config describes what New should build. Every field is optional:
leave Phone empty to skip the account client, BotToken empty to skip
the bot client. Session defaults to "session".

### `App`

```go
type App struct
```

App bundles the five ParsRubika clients. Nil members were not
requested in the Config (for example Bot when no BotToken was given).

### `New`

```go
func New(cfg *Config) (*App, error)
```

New builds the App described by cfg. It does NOT sign in — call
Login (one-shot work) or Run (sign in + listen for updates) first.
A nil cfg means "everything by default": account on session "session".

### `Login`

```go
func (app *App) Login(ctx context.Context) error
```

Login signs the account in (sendCode → code → signIn, prompting on
the terminal when the session file has no stored login) and wires the
fresh auth key into the Rubino and Rubiai clients. It returns as soon
as the sign-in completes; use Run to also listen for updates.

### `Run`

```go
func (app *App) Run(ctx context.Context) error
```

Run signs in (Login) and then blocks in the account's WebSocket
update loop until ctx is cancelled. It is the "stay online" mode:
attach account listeners before calling Run to react to messages.

### `Close`

```go
func (app *App) Close() error
```

Close releases every client. Safe to call more than once.

### `SendText`

```go
func SendText(session, phone, target, text string) error
```

SendText is the one-line message sender: it signs in with the stored
session (asking for the SMS code only the first time), sends one text
message to target (a user/group/channel GUID or username) and closes.

err := parsrubika.SendText("session", "+989123456789", "gAAAA...", "سلام")

### `Hangi`

```go
func Hangi(session, phone, guid, photo1, photo2 string, opts ...account.HangiOption) (*account.HangiProfileResult, error)
```

Hangi is the one-line dual (hangi) profile — پروفایل هنگی. It uploads
both images, fires ONE simultaneous dual bind and cleans up, leaving
exactly the two photos on the object. Pass guid = "" to hang the pair
on the signed-in account itself, or a group/channel GUID to hang it
on that chat.

res, err := parsrubika.Hangi("session", "+989123456789",
"c0DBedk0525ec07fb54a35b01baf2303", "1.jpg", "2.jpg")

Extra control (rounds, clean, progress, ...) goes through the same
account.HangiOption values used by Account.HangiProfile.

### `QuickBot`

```go
func QuickBot(token string, onMessage func(c *bot.Ctx)) error
```

QuickBot is the one-line bot runner: it builds a bot from token,
registers onMessage as the catch-all handler (it receives every
update) and blocks in the polling loop until the process is killed.

parsrubika.QuickBot(token, func(c *bot.Ctx) {
if c.Update.Text() != "" { c.Reply("سلام!") }
})
