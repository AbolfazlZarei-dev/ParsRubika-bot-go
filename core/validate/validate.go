// Package validate provides input-validation helpers for Rubika
// API parameters — chat IDs, usernames, phone numbers, etc.
package validate

import (
	"regexp"
	"strings"
)

// chatIDRE matches valid Rubika chat IDs: @username, or a 32-char
// GUID prefixed with b0/g0/c0/u0/s0.
var chatIDRE = regexp.MustCompile(`^(@[a-zA-Z0-9_]{3,32}|(b0|g0|c0|u0|s0)[a-zA-Z0-9]{30,})$`)

// userGUIDRE matches user GUIDs only (u0 prefix).
var userGUIDRE = regexp.MustCompile(`^u0[a-zA-Z0-9]{30,}$`)

// botGUIDRE matches bot GUIDs only (b0 prefix).
var botGUIDRE = regexp.MustCompile(`^b0[a-zA-Z0-9]{30,}$`)

// usernameRE matches valid Rubika usernames: 3-32 chars,
// alphanumeric + underscore, must contain at least one letter,
// must not start/end with `_` or a digit.
var usernameRE = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{1,30}[a-zA-Z0-9]$`)

// tokenRE matches valid Rubika bot tokens: 5 uppercase letters
// followed by 59 alphanumerics (64 total).
var tokenRE = regexp.MustCompile(`^[A-Z]{5}[0-9A-Z]{59}$`)

// ChatID reports whether `id` is a valid Rubika chat ID.
func ChatID(id string) bool { return chatIDRE.MatchString(id) }

// UserGUID reports whether `g` is a valid user GUID.
func UserGUID(g string) bool { return userGUIDRE.MatchString(g) }

// BotGUID reports whether `g` is a valid bot GUID.
func BotGUID(g string) bool { return botGUIDRE.MatchString(g) }

// Username reports whether `u` is a valid Rubika username.
func Username(u string) bool {
	u = strings.TrimPrefix(u, "@")
	return usernameRE.MatchString(u)
}

// Token reports whether `t` is a valid Rubika bot token.
func Token(t string) bool { return tokenRE.MatchString(t) }

// PhoneNumber normalises a Persian/Iranian phone number to the
// international format (98...). Returns "" on invalid input.
//
// The order of operations matters: the "+" prefix must be stripped
// BEFORE the "00"/"0" checks, otherwise inputs like "+098..." are
// mis-processed (the "0" check fires on "+0..." and converts the local
// "0" to "98", leaving "+98..." which then has its "+" stripped to
// "98..." but the leading local "0" was never removed).
//
// Supported input forms:
//
// 09981530922 → 989981530922 (local, leading 0)
// 9981530922 → 989981530922 (bare mobile, no prefix)
// +989981530922 → 989981530922 (international with +)
// 00989981530922 → 989981530922 (international with 00)
// 989981530922 → 989981530922 (already normalized)
func PhoneNumber(p string) string {
	// Convert Persian digits.
	p = strings.NewReplacer(
		"۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4",
		"۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9",
	).Replace(p)
	// Strip separators.
	p = strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(p)
	// Strip "+" FIRST so the "00"/"0" checks below see only digits.
	if strings.HasPrefix(p, "+") {
		p = p[1:]
	}
	if strings.HasPrefix(p, "00") {
		p = p[2:]
	}
	if strings.HasPrefix(p, "0") {
		p = "98" + p[1:]
	}
	// Bare Iranian mobile: 10 digits starting with "9" (e.g. 9981530922)
	// → prepend "98". Iranian mobiles without country code are exactly
	// 10 digits and always start with 9.
	if len(p) == 10 && strings.HasPrefix(p, "9") {
		p = "98" + p
	}
	if len(p) < 7 || len(p) > 15 {
		return ""
	}
	for _, r := range p {
		if r < '0' || r > '9' {
			return ""
		}
	}
	return p
}

// MessageText reports whether `t` is a non-empty message text.
func MessageText(t string) bool { return strings.TrimSpace(t) != "" }
