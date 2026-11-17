// Package parsrubika is a high-performance, modular Go library for the
// Rubika messaging platform, written from scratch with a focus on speed,
// security and ergonomics.
//
// ParsRubika is organised around five independent clients that share the
// same foundation:
//
// - parsrubika/account : the user / account client (encrypted API + WebSocket)
// - parsrubika/bot : the official bot HTTP API client
// - parsrubika/rubino : the Rubino social-network client
// - parsrubika/rubiai : the RubiAI streaming (SSE) chat client
// - parsrubika/tv : the Rubika TV / VOD client
// - parsrubika/plugins : a plugin framework shared by all clients
//
// The foundation is split into small, single-purpose packages so each
// piece can be used in isolation:
//
// - parsrubika/core : version, author, constants and typed errors
// - parsrubika/security : AES-CBC, RSA-PKCS1v15, key generation
// - parsrubika/text : markdown / html → metadata formatter
// - parsrubika/enums : shared enumerations (72 reactions, voice-chat states, etc.)
// - parsrubika/storage : pluggable session storage (sqlite / memory / file / string)
// - parsrubika/events : dynamic event wrapper used by the account client
// - parsrubika/transport: HTTP/2, HTTP/3, WebSocket, DC discovery, sync.Pool
// - parsrubika/files : chunked uploader / downloader with progress + thumbnails
// - parsrubika/extras : scheduler, cron, forms, qr, tasks, voicekit, ai, handy
//
// It covers the complete Rubika surface: the Bot API, the user-account
// client, Rubino and TV, with a colourful optional logger and a
// plugin framework on top.
//
// # Everything in one import
//
// This root package is also the everything-in-one entry point. New
// builds all five clients together behind a single App, and three
// one-line helpers cover the most common jobs — no long setup code:
//
//	app, _ := parsrubika.New(&parsrubika.Config{Phone: "+98...", Session: "session"})
//	defer app.Close()
//	app.Account.SendMessage(ctx, guid, "سلام")   // full account client
//	app.Rubino.MyProfile()                       // Rubino
//	parsrubika.SendText("session", "+98...", guid, "سلام") // one message
//	parsrubika.Hangi("session", "+98...", "", "1.jpg", "2.jpg") // dual-avatar hangi profile
//	parsrubika.QuickBot(token, func(c *bot.Ctx) { c.Reply("سلام!") })     // one-line bot
//
// Author : Abolfazl Zarei
// Version: 4.0.0
// License: LGPL v3
package parsrubika
