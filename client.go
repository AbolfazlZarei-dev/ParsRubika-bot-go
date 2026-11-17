package parsrubika

// The everything-in-one entry point. ParsRubika's five clients (account,
// bot, rubino, tv, rubiai) each have their own constructor and options;
// this file wires them together behind a single App so an application
// needs a handful of lines instead of a page of setup code — and three
// one-liners (SendText, Hangi, QuickBot) cover the most common jobs
// with a single call.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/rubiai"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/rubino"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/tv"
)

// Config describes what New should build. Every field is optional:
// leave Phone empty to skip the account client, BotToken empty to skip
// the bot client. Session defaults to "session".
type Config struct {
	// Phone is the account phone number in international format
	// (e.g. "+989123456789"). Empty = no account client.
	Phone string
	// Session is the session file name; "<Session>.bin" holds the login
	// so the SMS code is only asked for once. Default "session".
	Session string
	// Passphrase enables AES-256-GCM at-rest encryption of the session
	// file. Empty = plain text.
	Passphrase string
	// CodeProvider receives the login code / 2FA password. Nil = the
	// built-in interactive terminal prompt.
	CodeProvider account.CodeProvider
	// BotToken is a Rubika bot token; empty = no bot client.
	BotToken string
	// RubinoAuth overrides the Rubino auth key. Empty = the account
	// auth is copied over after sign-in.
	RubinoAuth string
	// AIAuth overrides the RubiAI auth key. Empty = the account auth
	// is copied over after sign-in.
	AIAuth string
	// Proxy URL (http://, https://, socks5://) shared by every client.
	Proxy string
	// Debug turns on verbose account logging.
	Debug bool
}

// App bundles the five ParsRubika clients. Nil members were not
// requested in the Config (for example Bot when no BotToken was given).
type App struct {
	// Account is the user-account client (messages, groups, channels,
	// files, avatars, ...). It shares nothing with Bot: Rubika's bot
	// API and the account API are two different services.
	Account *account.Account
	// Bot is the official bot API client.
	Bot *bot.Bot
	// Rubino is the Rubino social-network client. Its auth is wired
	// from the account after Login when Config.RubinoAuth was empty.
	Rubino *rubino.Rubino
	// TV is the Rubika TV / VOD client (no auth needed).
	TV *tv.TV
	// Rubiai is the RubiAI chat client. Its auth is wired from the
	// account after Login when Config.AIAuth was empty.
	Rubiai *rubiai.RubiaAI

	phone     string
	rubinoSet bool
	aiSet     bool
}

// New builds the App described by cfg. It does NOT sign in — call
// Login (one-shot work) or Run (sign in + listen for updates) first.
// A nil cfg means "everything by default": account on session "session".
func New(cfg *Config) (*App, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	app := &App{phone: cfg.Phone, rubinoSet: cfg.RubinoAuth != "", aiSet: cfg.AIAuth != ""}

	if cfg.Phone != "" {
		session := cfg.Session
		if session == "" {
			session = "session"
		}
		opts := []account.Option{
			account.WithPhoneNumber(cfg.Phone),
			account.WithProxy(cfg.Proxy),
			account.WithDebug(cfg.Debug),
		}
		// A nil CodeProvider means "use the built-in interactive stdin
		// prompt" (documented on Config), so only install the option
		// when the caller supplied one — WithCodeProvider(nil) is a
		// validation error by design.
		if cfg.CodeProvider != nil {
			opts = append(opts, account.WithCodeProvider(cfg.CodeProvider))
		}
		if cfg.Passphrase != "" {
			opts = append(opts, account.WithFileSessionPassphrase(session, cfg.Passphrase))
		} else {
			opts = append(opts, account.WithFileSession(session))
		}
		acc, err := account.New(opts...)
		if err != nil {
			return nil, fmt.Errorf("parsrubika: account: %w", err)
		}
		app.Account = acc
	}

	if cfg.BotToken != "" {
		app.Bot = bot.New(cfg.BotToken)
	}

	// Rubino, TV and Rubiai are cheap and auth-light, so they are
	// always constructed; Rubino/Rubiai receive the account auth after
	// Login unless an explicit override was configured.
	app.Rubino = rubino.New(rubino.WithAuth(cfg.RubinoAuth), rubino.WithProxy(cfg.Proxy))
	app.TV = tv.New(tv.WithProxy(cfg.Proxy))
	app.Rubiai = rubiai.New(rubiai.WithProxy(cfg.Proxy))
	if cfg.AIAuth != "" {
		app.Rubiai.SetAuth(cfg.AIAuth)
	}
	return app, nil
}

// Login signs the account in (sendCode → code → signIn, prompting on
// the terminal when the session file has no stored login) and wires the
// fresh auth key into the Rubino and Rubiai clients. It returns as soon
// as the sign-in completes; use Run to also listen for updates.
func (app *App) Login(ctx context.Context) error {
	if app.Account == nil {
		return nil
	}
	if err := app.Account.Start(ctx, app.phone); err != nil {
		return fmt.Errorf("parsrubika: login: %w", err)
	}
	if auth := app.Account.Auth(); auth != "" {
		if !app.rubinoSet {
			app.Rubino.SetAuth(auth)
		}
		if !app.aiSet {
			app.Rubiai.SetAuth(auth)
		}
	}
	return nil
}

// Run signs in (Login) and then blocks in the account's WebSocket
// update loop until ctx is cancelled. It is the "stay online" mode:
// attach account listeners before calling Run to react to messages.
func (app *App) Run(ctx context.Context) error {
	if app.Account == nil {
		return errors.New("parsrubika: Run needs Config.Phone")
	}
	if err := app.Login(ctx); err != nil {
		return err
	}
	return app.Account.Listen(ctx)
}

// Close releases every client. Safe to call more than once.
func (app *App) Close() error {
	if app.Account != nil {
		_ = app.Account.Close()
	}
	return nil
}

// SendText is the one-line message sender: it signs in with the stored
// session (asking for the SMS code only the first time), sends one text
// message to target (a user/group/channel GUID or username) and closes.
//
//	err := parsrubika.SendText("session", "+989123456789", "gAAAA...", "سلام")
func SendText(session, phone, target, text string) error {
	app, err := New(&Config{Phone: phone, Session: session})
	if err != nil {
		return err
	}
	defer app.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := app.Login(ctx); err != nil {
		return err
	}
	_, err = app.Account.SendMessage(ctx, target, text)
	return err
}

// Hangi is the one-line dual (hangi) profile — پروفایل هنگی. It uploads
// both images, fires ONE simultaneous dual bind and cleans up, leaving
// exactly the two photos on the object. Pass guid = "" to hang the pair
// on the signed-in account itself, or a group/channel GUID to hang it
// on that chat.
//
//	res, err := parsrubika.Hangi("session", "+989123456789",
//	    "c0DBedk0525ec07fb54a35b01baf2303", "1.jpg", "2.jpg")
//
// Extra control (rounds, clean, progress, ...) goes through the same
// account.HangiOption values used by Account.HangiProfile.
func Hangi(session, phone, guid, photo1, photo2 string, opts ...account.HangiOption) (*account.HangiProfileResult, error) {
	app, err := New(&Config{Phone: phone, Session: session})
	if err != nil {
		return nil, err
	}
	defer app.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	if err := app.Login(ctx); err != nil {
		return nil, err
	}
	if guid == "" {
		guid = app.Account.GUID()
	}
	return app.Account.HangiProfile(ctx, guid, photo1, photo2, opts...)
}

// QuickBot is the one-line bot runner: it builds a bot from token,
// registers onMessage as the catch-all handler (it receives every
// update) and blocks in the polling loop until the process is killed.
//
//	parsrubika.QuickBot(token, func(c *bot.Ctx) {
//	    if c.Update.Text() != "" { c.Reply("سلام!") }
//	})
func QuickBot(token string, onMessage func(c *bot.Ctx)) error {
	if token == "" {
		return errors.New("parsrubika: QuickBot needs a bot token")
	}
	if onMessage == nil {
		return errors.New("parsrubika: QuickBot needs a handler")
	}
	b := bot.New(token)
	b.Handle(nil, func(c *bot.Ctx) error {
		onMessage(c)
		return nil
	})
	return b.Run(context.Background(), "")
}
