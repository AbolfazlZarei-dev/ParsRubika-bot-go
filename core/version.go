// Package core holds the version, author and shared constants for
// ParsRubika. It is intentionally tiny — every other package can import
// it without pulling in any heavy dependencies.
package core

// Author is the maintainer of ParsRubika.
const Author = "Abolfazl Zarei"

// Version is the current ParsRubika release.
const Version = "4.0.0"

// Codename is the marketing name for this release line.
const Codename = "Nova"

// License is the SPDX-style identifier of the distribution license.
const License = "LGPL-3.0-only"

// UserAgent is the default User-Agent string used by every ParsRubika
// HTTP client. It already includes the version so that server-side
// telemetry can distinguish ParsRubika traffic.
const UserAgent = "ParsRubika/" + Version + " (go; +" + Author + ")"

// APIVersion is the wire-format API version sent in every request body.
const APIVersion = "6"

// DefaultLangCode is the default language code sent to the server.
const DefaultLangCode = "fa"

// DefaultTimeout is the default HTTP timeout, in seconds.
const DefaultTimeout = 30

// DefaultMaxRetries is the default number of retries for failed requests.
const DefaultMaxRetries = 5

// DefaultChunkSize is the default upload/download chunk size (1 MiB).
const DefaultChunkSize = 1 << 20

// DefaultDownloadChunkSize is the default download chunk size (128 KiB).
const DefaultDownloadChunkSize = 128 * 1024

// WelcomeBanner returns the standard library welcome banner.
func WelcomeBanner() string {
	return "Welcome to ParsRubika " + Version + " (" + Codename + ")\n" +
		"Authored by " + Author + ". Licensed under " + License + ".\n" +
		"Absolutely no warranty. Use responsibly.\n"
}
