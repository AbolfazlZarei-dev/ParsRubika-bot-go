package core

// ClientKind enumerates the first-class ParsRubika clients. Each kind
// maps to a standalone, fully-modular client package that shares one
// fast, secure foundation (core / security / transport / events).
//
// The five clients are registered here as metadata only — core must not
// import any of them (that would create an import cycle, since every
// client imports core). The Package field carries the import path so
// tooling and docs can point at the right module without a hard
// dependency.
//
// Mirrors the platform registry in platforms.go: a single, authoritative
// place where every first-class client is listed, in priority order.
type ClientKind string

const (
	// ClientBot is the official Rubika Bot HTTP API client.
	ClientBot ClientKind = "bot"
	// ClientAccount is the Rubika user (account) client — the encrypted
	// messenger envelope plus the live WebSocket update stream.
	ClientAccount ClientKind = "account"
	// ClientRubino is the Rubino social-network client.
	ClientRubino ClientKind = "rubino"
	// ClientTV is the Rubika TV / VOD client.
	ClientTV ClientKind = "tv"
	// ClientRubiAI is the RubiAI assistant client (chat, sessions, SSE
	// streaming, speech-to-text).
	ClientRubiAI ClientKind = "rubiai"
)

// ClientInfo is the metadata for one first-class client. It is the
// client-side analogue of Identity: name, import path, a one-line
// description and the service's base URL. All fields are static strings
// so core stays dependency-free.
type ClientInfo struct {
	// Kind is the ClientKind this info describes.
	Kind ClientKind
	// Name is the human-readable client name.
	Name string
	// Package is the full Go import path of the client package.
	Package string
	// Description is a one-line summary of what the client does.
	Description string
	// BaseURL is the primary service endpoint the client talks to. For
	// clients with multiple servers (TV, account) this is the first /
	// discovery endpoint.
	BaseURL string
}

// AllClients lists every first-class ParsRubika client, in priority
// order: the messaging clients first (bot, account), then the
// Rubika-ecosystem vertical clients (rubino, tv, rubiai). It mirrors
// AllPlatforms so the two registries stay in lock-step.
var AllClients = []ClientInfo{
	{
		Kind:        ClientBot,
		Name:        "Bot API",
		Package:     "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot",
		Description: "Official Rubika Bot HTTP API client with rules, middleware, polling and webhook.",
		BaseURL:     "https://botapi.rubika.ir/v3/",
	},
	{
		Kind:        ClientAccount,
		Name:        "Account",
		Package:     "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account",
		Description: "Rubika user/account client: encrypted envelope, live WebSocket updates, full method surface.",
		BaseURL:     "https://getdcmess.iranlms.ir/",
	},
	{
		Kind:        ClientRubino,
		Name:        "Rubino",
		Package:     "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/rubino",
		Description: "Rubino social-network client: posts, stories, follow, comments, profiles, uploads.",
		BaseURL:     "https://rubino1.iranlms.ir/",
	},
	{
		Kind:        ClientTV,
		Name:        "TV",
		Package:     "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/tv",
		Description: "Rubika TV / VOD client: browse, search, like, wish-list, cast, series (standalone first-class module).",
		BaseURL:     "https://rbvod1.iranlms.ir/",
	},
	{
		Kind:        ClientRubiAI,
		Name:        "RubiAI",
		Package:     "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/rubiai",
		Description: "RubiAI assistant client: chat history, sessions, SSE streaming and speech-to-text.",
		BaseURL:     "https://rubichat1.iranlms.ir/",
	},
}

// ClientFor returns the ClientInfo for the given kind, or the zero
// ClientInfo (with an empty Kind) when the kind is unknown. Callers
// should check the returned Kind against the empty string.
func ClientFor(k ClientKind) ClientInfo {
	for _, c := range AllClients {
		if c.Kind == k {
			return c
		}
	}
	return ClientInfo{}
}

// ClientNames returns the Name of every first-class client, in
// AllClients order. Handy for dashboards, /help output and docs.
func ClientNames() []string {
	out := make([]string, 0, len(AllClients))
	for _, c := range AllClients {
		out = append(out, c.Name)
	}
	return out
}

// ClientKinds returns every ClientKind, in AllClients order.
func ClientKinds() []ClientKind {
	out := make([]ClientKind, 0, len(AllClients))
	for _, c := range AllClients {
		out = append(out, c.Kind)
	}
	return out
}
