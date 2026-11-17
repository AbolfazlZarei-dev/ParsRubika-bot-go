// Package netutil provides SSRF defenses for outbound HTTP requests the
// library issues on behalf of the user (URL-based uploads, JSON fetches,
// link scraping). It combines two layers:
//
//  1. ValidateURL / AssertURL — a fast syntactic + DNS-resolution check
//     run before the request is sent. Rejects non-http(s) schemes,
//     localhost, and any host that resolves to a private / loopback /
//     link-local / metadata IP (including the cloud metadata endpoint
//     169.254.169.254).
//
//  2. NewSafeTransport — an *http.Transport whose Dialer.Control hook
//     inspects the resolved IP right before connect and aborts the dial
//     if it is private. This defeats DNS-rebinding: even if the resolver
//     returned a public IP at pre-flight time, the actual connect is
//     blocked when the address is private.
//
// Both layers are used together so a single miss at one layer is caught
// by the other.
package netutil

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"time"
)

// ErrBlocked is returned when a URL or address is refused by the SSRF
// policy (private / loopback / link-local / metadata IP, or a forbidden
// scheme).
var ErrBlocked = errors.New("netutil: address blocked by SSRF policy")

// AllowedSchemes are the URL schemes permitted for outbound fetches.
var AllowedSchemes = map[string]bool{"http": true, https: true}

const https = "https"

// IsBlockedIP reports whether ip points at a host the library must never
// contact on the user's behalf: loopback, private RFC-1918 ranges,
// link-local (incl. the cloud metadata endpoint 169.254.169.254),
// unspecified, multicast, and other non-global-unicast addresses.
func IsBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsMulticast() || ip.IsUnspecified() {
		return true
	}
	// Cloud metadata endpoints: AWS/GCP/Azure live on 169.254.169.254.
	// IsLinkLocalUnicast already covers it, but be explicit for clarity.
	if ip.Equal(net.IPv4(169, 254, 169, 254)) {
		return true
	}
	// Carrier-grade NAT (100.64.0.0/10) and other reserved ranges are
	// not IsPrivate in older Go; reject anything that is not global
	// unicast to be safe.
	if !ip.IsGlobalUnicast() {
		return true
	}
	return false
}

// ValidateURL parses rawURL and applies the syntactic SSRF policy:
// http/https scheme only, non-empty host, no localhost, and an IP-literal
// host (if any) that is not private. It does NOT perform DNS resolution,
// so a hostname that resolves to a private IP is only caught by
// AssertURL or NewSafeTransport.
func ValidateURL(rawURL string) (*url.URL, error) {
	if rawURL == "" {
		return nil, fmt.Errorf("%w: empty URL", ErrBlocked)
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("netutil: parse URL: %w", err)
	}
	if !AllowedSchemes[u.Scheme] {
		return nil, fmt.Errorf("%w: scheme %q not allowed", ErrBlocked, u.Scheme)
	}
	host := u.Hostname()
	if host == "" {
		return nil, fmt.Errorf("%w: empty host", ErrBlocked)
	}
	// Reject the literal name "localhost" and "*.localhost".
	if host == "localhost" || hasSuffixDotLocalhost(host) {
		return nil, fmt.Errorf("%w: localhost host", ErrBlocked)
	}
	// If the host is already an IP literal we can decide now.
	if ip := net.ParseIP(host); ip != nil {
		if IsBlockedIP(ip) {
			return nil, fmt.Errorf("%w: private IP literal %s", ErrBlocked, ip)
		}
	}
	return u, nil
}

func hasSuffixDotLocalhost(host string) bool {
	// "foo.localhost"
	end := len(host) - len(".localhost")
	return end > 0 && host[end:] == ".localhost"
}

// AssertURL runs ValidateURL and then resolves the host name (when it is
// not already an IP literal) and verifies that NONE of the resolved
// addresses are private. The check uses the default resolver with the
// supplied context so it honours cancellation / timeouts.
//
// AssertURL is a pre-flight check. For full defense against DNS
// rebinding, also route the request through NewSafeTransport.
func AssertURL(ctx context.Context, rawURL string) (*url.URL, error) {
	u, err := ValidateURL(rawURL)
	if err != nil {
		return nil, err
	}
	host := u.Hostname()
	if net.ParseIP(host) != nil {
		return u, nil // already validated by ValidateURL
	}
	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil {
		return nil, fmt.Errorf("netutil: resolve %s: %w", host, err)
	}
	for _, ip := range ips {
		if IsBlockedIP(ip) {
			return nil, fmt.Errorf("%w: %s resolves to private IP %s", ErrBlocked, host, ip)
		}
	}
	return u, nil
}

// ControlSSRF is the Dialer.Control hook: it receives the resolved
// IP:port just before connect and refuses private addresses. This runs
// AFTER name resolution, so it sees the real address being contacted.
//
// It is exported so the main API/WebSocket transport can install the
// same SSRF defense used by NewSafeTransport, without having to replace
// its whole http.Transport.
func ControlSSRF(network, address string, _ syscall.RawConn) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		// No port? Treat the whole string as the host.
		host = address
	}
	if ip := net.ParseIP(host); ip != nil && IsBlockedIP(ip) {
		return fmt.Errorf("%w: dial to %s refused", ErrBlocked, ip)
	}
	return nil
}

// controlSSRF is retained as an unexported alias for internal callers.
func controlSSRF(network, address string, c syscall.RawConn) error {
	return ControlSSRF(network, address, c)
}

// NewSafeTransport returns an *http.Transport that blocks dials to
// private / loopback / link-local / metadata IPs via a Control hook.
// Use it as the transport of any http.Client that fetches user-supplied
// URLs (uploads, scrapes, JSON fetches).
func NewSafeTransport(timeout time.Duration) *http.Transport {
	dialer := &net.Dialer{
		Timeout:  timeout,
		Control:  controlSSRF,
		Resolver: net.DefaultResolver,
	}
	return &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// NewSafeClient returns an http.Client whose transport blocks SSRF
// targets and whose overall timeout is bounded.
func NewSafeClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &http.Client{
		Transport: NewSafeTransport(timeout),
		Timeout:   timeout,
	}
}

// MaxFetchBytes caps how many bytes a user-URL fetch may read into memory
// (UploadFromURL, FetchJSON, ScrapeChannelPostLink). It prevents a
// malicious remote server from exhausting RAM by streaming an enormous
// response.
const MaxFetchBytes = 50 << 20 // 50 MiB
