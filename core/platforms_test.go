package core

import (
	"testing"
)

func TestIdentityForPWA(t *testing.T) {
	id := IdentityFor(PlatformPWAFull, "en")
	if id.Platform != "PWA" {
		t.Errorf("Platform = %q", id.Platform)
	}
	if id.LangCode != "en" {
		t.Errorf("LangCode = %q", id.LangCode)
	}
	if id.Package != "m.rubika.ir" {
		t.Errorf("Package = %q", id.Package)
	}
}

func TestIdentityForAndroid(t *testing.T) {
	id := IdentityFor(PlatformAndroidFull, "fa")
	if id.Platform != "Android" {
		t.Errorf("Platform = %q", id.Platform)
	}
	if id.Package != "app.rbmain.a" {
		t.Errorf("Package = %q", id.Package)
	}
}

func TestIdentityForRubX(t *testing.T) {
	id := IdentityFor(PlatformRubX, "fa")
	if id.Platform != "Android" {
		t.Errorf("Platform = %q", id.Platform)
	}
	if id.Package != "ir.rubx.bapp" {
		t.Errorf("Package = %q", id.Package)
	}
}

func TestIdentityForRubikids(t *testing.T) {
	id := IdentityFor(PlatformRubikids, "fa")
	if id.Package != "app.rbmain.kids" {
		t.Errorf("Package = %q", id.Package)
	}
}

func TestIdentityForRubino(t *testing.T) {
	id := IdentityFor(PlatformRubino, "fa")
	if id.Package != "app.rubino.main" {
		t.Errorf("Package = %q", id.Package)
	}
}

func TestUserAgentForAndroid(t *testing.T) {
	if ua := UserAgentFor(PlatformAndroidFull); ua != "okhttp/3.12.1" {
		t.Errorf("UA = %q", ua)
	}
	if ua := UserAgentFor(PlatformPWAFull); ua == "okhttp/3.12.1" {
		t.Errorf("UA should not be okhttp for PWA")
	}
}

func TestDCDiscoveryURLForPWA(t *testing.T) {
	if u := DCDiscoveryURLFor(PlatformPWAFull); u != "https://getdcmess.iranlms.ir/" {
		t.Errorf("DC URL for PWA = %q", u)
	}
}

func TestStoragePrefixForPWA(t *testing.T) {
	if p := StoragePrefixFor(PlatformPWAFull); p != "messenger" {
		t.Errorf("Storage prefix for PWA = %q", p)
	}
}

func TestAllPlatforms(t *testing.T) {
	// Only Rubika platforms remain: PWA, Android,
	// RubX, Rubikids, Rubino (5 entries).
	if len(AllPlatforms) != 5 {
		t.Errorf("AllPlatforms has %d entries, want 5", len(AllPlatforms))
	}
}

func TestNoNonRubikaPlatforms(t *testing.T) {
	// Regression guard: ParsRubika is Rubika-only. Every entry in
	// AllPlatforms must be one of the known Rubika platform kinds —
	// no unrelated platforms may sneak in.
	for _, p := range AllPlatforms {
		switch p {
		case PlatformPWAFull, PlatformAndroidFull, PlatformRubX, PlatformRubikids, PlatformRubino:
			// OK — a known Rubika platform.
		default:
			t.Errorf("unknown/non-Rubika platform %q present in AllPlatforms", p)
		}
	}
}
