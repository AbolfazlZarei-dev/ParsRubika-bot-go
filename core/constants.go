package core

// Identity holds the client-side identity sent inside the `client` field
// of every request body. The exact values are tuned to mimic the
// official web / Android clients so the server returns the correct
// feature flags.
type Identity struct {
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	Platform   string `json:"platform"`
	Package    string `json:"package"`
	LangCode   string `json:"lang_code"`
}

// WebIdentity returns the identity used by PWA (web) clients.
func WebIdentity(lang string) Identity {
	return Identity{
		AppName:    "Main",
		AppVersion: "2.4.6",
		Platform:   "PWA",
		Package:    "m.rubika.ir",
		LangCode:   lang,
	}
}

// AndroidIdentity returns the identity used by Android clients.
func AndroidIdentity(lang string) Identity {
	return Identity{
		AppName:    "Main",
		AppVersion: "3.8.2",
		Platform:   "Android",
		Package:    "app.rbmain.a",
		LangCode:   lang,
	}
}

// Platform is the top-level selector used by the account client.
type Platform string

const (
	PlatformWeb     Platform = "PWA"
	PlatformAndroid Platform = "Android"
)

// AllDCDiscoveryURLs is the list of endpoints tried, in order, to
// discover the active Rubika data-centers. The transport package is the canonical
// source (core.AllDCDiscoveryURLs is kept for backwards compatibility).
var AllDCDiscoveryURLs = []string{
	"https://getdcmess.iranlms.ir/",
	"https://messengerget.iranlms.ir/",
}

// DefaultWebSocketPingInterval controls how often the keep-alive ping
// is sent over the WebSocket connection.
const DefaultWebSocketPingInterval = 10 // seconds

// MaxMessageTextLength is the server-imposed maximum size of a single
// text message; longer texts are chunked into multiple messages.
const MaxMessageTextLength = 4200

// RSAKeyBits is the size (in bits) of the RSA key generated at sign-in
// time. The Rubika server accepts 2048-bit keys (the official MaxRubika
// client uses RSA-2048); 2048 is the minimum modern baseline and is
// resistant to feasible factorisation.
//
// Existing on-disk sessions keep using whatever key they were created
// with (the private key is loaded verbatim); only NEW sign-ins generate
// a 2048-bit key. See MIGRATION.md for details.
const RSAKeyBits = 2048

// AuthLength is the length of the auth key string (32 chars).
const AuthLength = 32

// PasskeyLength is the length of the derived AES key (32 chars → AES-256).
const PasskeyLength = 32
