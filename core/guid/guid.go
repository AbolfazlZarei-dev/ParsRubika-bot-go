// Package guid provides utility functions for working with Rubika
// GUIDs — detecting their type, validating their format, and
// normalising usernames/links.
package guid

import "strings"

// Kind is the kind of object a GUID refers to.
type Kind string

const (
	KindUser    Kind = "User"
	KindBot     Kind = "Bot"
	KindGroup   Kind = "Group"
	KindChannel Kind = "Channel"
	KindService Kind = "Service"
	KindUnknown Kind = "Unknown"
)

// Detect infers the kind of object from the GUID prefix.
func Detect(g string) Kind {
	switch {
	case strings.HasPrefix(g, "c0"):
		return KindChannel
	case strings.HasPrefix(g, "g0"):
		return KindGroup
	case strings.HasPrefix(g, "b0"):
		return KindBot
	case strings.HasPrefix(g, "u0"):
		return KindUser
	case strings.HasPrefix(g, "s0"):
		return KindService
	default:
		return KindUnknown
	}
}

// IsUser reports whether `g` is a user GUID.
func IsUser(g string) bool { return Detect(g) == KindUser }

// IsBot reports whether `g` is a bot GUID.
func IsBot(g string) bool { return Detect(g) == KindBot }

// IsGroup reports whether `g` is a group GUID.
func IsGroup(g string) bool { return Detect(g) == KindGroup }

// IsChannel reports whether `g` is a channel GUID.
func IsChannel(g string) bool { return Detect(g) == KindChannel }

// IsService reports whether `g` is a service GUID.
func IsService(g string) bool { return Detect(g) == KindService }

// IsValid reports whether `g` looks like a valid GUID (32 chars
// after the 2-char prefix).
func IsValid(g string) bool {
	if len(g) < 34 {
		return false
	}
	return Detect(g) != KindUnknown
}

// StripAt removes a leading `@` from a username.
func StripAt(s string) string { return strings.TrimPrefix(s, "@") }

// TrailingSlug returns the trailing path segment of a URL or the
// input itself when it does not contain a slash.
func TrailingSlug(link string) string {
	if link == "" {
		return ""
	}
	if i := strings.LastIndex(link, "/"); i >= 0 {
		return link[i+1:]
	}
	return link
}
