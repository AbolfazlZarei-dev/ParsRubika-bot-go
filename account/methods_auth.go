package account

import (
	"context"
	"fmt"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"regexp"
	"strings"
)

// Session, device registration and account lifecycle.

// SendCode triggers the SMS / Internal verification code for `phone`.
func (a *Account) SendCode(ctx context.Context, phone, passKey, sendType string) (*events.Event, error) {
	if sendType == "" {
		sendType = "SMS"
	}
	if sendType != "SMS" && sendType != "Internal" {
		return nil, fmt.Errorf("account: send_type must be 'SMS' or 'Internal'")
	}
	in := map[string]any{
		"phone_number": phone,
		"send_type":    sendType,
	}
	// pass_key must NOT be included when empty (server returns INVALID_INPUT).
	if passKey != "" {
		in["pass_key"] = passKey
	}
	return a.Request(ctx, RequestOptions{
		Method:     "sendCode",
		Input:      in,
		TmpSession: true,
		Encrypt:    true,
	})
}

// SignIn completes the registration handshake.
func (a *Account) SignIn(ctx context.Context, phoneCode, phone, phoneCodeHash, publicKey string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method: "signIn",
		Input: map[string]any{
			"phone_code":      phoneCode,
			"phone_number":    phone,
			"phone_code_hash": phoneCodeHash,
			"public_key":      publicKey,
		},
		TmpSession: true,
		Encrypt:    true,
	})
}

// RegisterDevice registers the device model with the server.
func (a *Account) RegisterDevice(ctx context.Context, deviceModel string) (*events.Event, error) {
	var input map[string]any
	switch a.platform {
	case "Android":
		input = androidDevice(a.userAgent, a.langCode, a.identity.AppVersion, deviceModel)
	default:
		input = browserDevice(a.userAgent, a.langCode, a.identity.AppVersion, deviceModel)
	}
	return a.Request(ctx, RequestOptions{Method: "registerDevice", Input: input, Encrypt: true})
}

// UnregisterDevice removes the current device registration from the
// server. Use this for a clean logout (without it the device lingers
// in the server's device list and may keep receiving push tokens).
// For a full logout, follow this with Logout (in methods_auth2fa.go).
func (a *Account) UnregisterDevice(ctx context.Context) (*events.Event, error) {
	input := map[string]any{
		"lang_code":   a.langCode,
		"app_version": a.identity.AppVersion,
	}
	return a.Request(ctx, RequestOptions{Method: "unregisterDevice", Input: input, Encrypt: true})
}

var browserRE = regexp.MustCompile(`(?i)(opera|chrome|safari|firefox|msie|trident)/(\d+)`)

var digitsRE = regexp.MustCompile(`\d+`)

func browserDevice(userAgent, langCode, appVersion, deviceModel string) map[string]any {
	dm := parseBrowserModel(userAgent)
	if deviceModel != "" {
		dm = deviceModel
	}
	return map[string]any{
		"token":          "",
		"lang_code":      langCode,
		"token_type":     "Firebase",
		"app_version":    "PW_" + appVersion,
		"system_version": parseSystemVersion(userAgent),
		"device_model":   dm,
		"device_hash":    "2" + strings.Join(digitsRE.FindAllString(userAgent, -1), ""),
	}
}

func androidDevice(userAgent, langCode, appVersion, _ string) map[string]any {
	return map[string]any{
		"token":          "",
		"lang_code":      langCode,
		"token_type":     "Firebase",
		"app_version":    "MA_" + appVersion,
		"system_version": "SDK 22",
		"device_model":   "samsungSM-G925F",
		"device_hash":    "2" + strings.Join(digitsRE.FindAllString(userAgent, -1), ""),
	}
}

func parseBrowserModel(ua string) string {
	m := browserRE.FindStringSubmatch(ua)
	if m == nil {
		return "Unknown"
	}
	return m[1] + " " + m[2]
}

func parseSystemVersion(ua string) string {
	mapping := map[string]string{
		"Windows NT 10.0": "Windows 10",
		"Windows NT 6.2":  "Windows 8",
		"Windows NT 6.1":  "Windows 7",
		"Windows NT 6.0":  "Windows Vista",
		"Windows NT 5.1":  "Windows XP",
		"Windows NT 5.0":  "Windows 2000",
		"Mac":             "Mac/iOS",
		"X11":             "UNIX",
		"Linux":           "Linux",
	}
	for k, v := range mapping {
		if strings.Contains(ua, k) {
			return v
		}
	}
	return "Unknown"
}

// Logout logs out the current session.
func (a *Account) Logout(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "logout", Input: map[string]any{}, Encrypt: true})
}

// DeleteAccount requests account deletion.
func (a *Account) DeleteAccount(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "requestDeleteAccount", Input: map[string]any{}, Encrypt: true})
}

// UpgradeToAPI6 upgrades an API5 auth to API6 by uploading a new RSA
// public key. Returns the private key (PEM) that the caller should
// store alongside the auth.
func (a *Account) UpgradeToAPI6(ctx context.Context) (privateKeyPEM string, err error) {
	pub, priv, err := generateKeyPair()
	if err != nil {
		return "", err
	}
	_, err = a.Request(ctx, RequestOptions{
		Method:  "setPublicKey",
		Input:   map[string]any{"public_key": pub},
		Encrypt: true,
	})
	if err != nil {
		return "", err
	}
	return priv, nil
}

// GetUnconfirmedSessions returns the unconfirmed session list.
func (a *Account) GetUnconfirmedSessions(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getUnconfirmedSessions", Input: map[string]any{}, Encrypt: true})
}

// ActionOnUnconfirmedSession confirms or denies an unconfirmed session.
func (a *Account) ActionOnUnconfirmedSession(ctx context.Context, action, unconfirmedSessionKey string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "actionOnUnconfirmedSession", Input: map[string]any{
		"action":                  action, // "Confirm" or "Deny"
		"unconfirmed_session_key": unconfirmedSessionKey,
	}, Encrypt: true})
}

// ConfirmUnconfirmedSession is a convenience wrapper (action = "Confirm").
func (a *Account) ConfirmUnconfirmedSession(ctx context.Context, key string) (*events.Event, error) {
	return a.ActionOnUnconfirmedSession(ctx, "Confirm", key)
}

// DenyUnconfirmedSession is a convenience wrapper (action = "Deny").
func (a *Account) DenyUnconfirmedSession(ctx context.Context, key string) (*events.Event, error) {
	return a.ActionOnUnconfirmedSession(ctx, "Deny", key)
}

// generateKeyPair is a local alias to avoid importing security in this file.
func generateKeyPair() (pub, priv string, err error) {
	return generateKeyPairImpl()
}

// TerminateSession terminates a single active session by id.
func (a *Account) TerminateSession(ctx context.Context, sessionID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "terminateSession", Input: map[string]any{"session_key": sessionID}, Encrypt: true})
}

// TerminateOtherSessions terminates every other active session.
func (a *Account) TerminateOtherSessions(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "terminateOtherSessions", Input: map[string]any{}, Encrypt: true})
}

// GetMySessions returns the active sessions list.
func (a *Account) GetMySessions(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getMySessions", Input: map[string]any{}, Encrypt: true})
}
