package core

// PlatformKind enumerates the supported Rubika client platforms. Each
// platform maps to a distinct identity (app_name, app_version, package,
// headers) sent to the Rubika servers.
//
// ParsRubika is a Rubika-only library: every platform below belongs to
// the Rubika ecosystem. Non-Rubika platforms are intentionally not
// supported.
type PlatformKind string

const (
	// PlatformPWAFull is the default PWA / web platform.
	PlatformPWAFull PlatformKind = "PWA"
	// PlatformAndroidFull is the Android app platform.
	PlatformAndroidFull PlatformKind = "Android"
	// PlatformRubX is the RubX Android variant.
	PlatformRubX PlatformKind = "RubX"
	// PlatformRubikids is the Rubikids Android variant.
	PlatformRubikids PlatformKind = "Rubikids"
	// PlatformRubino is the Rubino Android variant.
	PlatformRubino PlatformKind = "Rubino"
)

// AllPlatforms lists every supported Rubika platform, in priority order.
var AllPlatforms = []PlatformKind{
	PlatformPWAFull,
	PlatformAndroidFull,
	PlatformRubX,
	PlatformRubikids,
	PlatformRubino,
}

// IdentityFor returns the per-platform Identity. Pass an empty
// language code to use the default ("fa").
func IdentityFor(p PlatformKind, lang string) Identity {
	if lang == "" {
		lang = DefaultLangCode
	}
	switch p {
	case PlatformAndroidFull:
		return AndroidIdentity(lang)
	case PlatformRubX:
		return Identity{
			AppName:    "Main",
			AppVersion: "3.8.2",
			Platform:   "Android",
			Package:    "ir.rubx.bapp",
			LangCode:   lang,
		}
	case PlatformRubikids:
		return Identity{
			AppName:    "Main",
			AppVersion: "3.8.2",
			Platform:   "Android",
			Package:    "app.rbmain.kids",
			LangCode:   lang,
		}
	case PlatformRubino:
		return Identity{
			AppName:    "Main",
			AppVersion: "3.8.2",
			Platform:   "Android",
			Package:    "app.rubino.main",
			LangCode:   lang,
		}
	default: // PWA
		return WebIdentity(lang)
	}
}

// UserAgentFor returns the User-Agent string for the given platform.
// Android platforms use `okhttp/3.12.1`; everything else uses the
// default ParsRubika User-Agent.
func UserAgentFor(p PlatformKind) string {
	switch p {
	case PlatformAndroidFull, PlatformRubX, PlatformRubikids, PlatformRubino:
		return "okhttp/3.12.1"
	default:
		return UserAgent
	}
}

// DCDiscoveryURL returns the Rubika DC-discovery endpoint. All Rubika
// platforms share the same endpoint.
func DCDiscoveryURLFor(p PlatformKind) string {
	_ = p
	return "https://getdcmess.iranlms.ir/"
}

// StoragePrefixFor returns the file-server URL prefix for the given
// platform. All Rubika platforms use the `messenger` prefix.
func StoragePrefixFor(p PlatformKind) string {
	_ = p
	return "messenger"
}
