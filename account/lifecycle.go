package account

import (
	"context"
	cryptoRand "crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/security"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/storage"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport"
)

// Start performs the connection + (optional) sign-in handshake. If a
// session has been saved previously it is reused; otherwise the
// `phoneNumber` (or the one passed in) is used to send a code and
// sign in interactively.
func (a *Account) Start(ctx context.Context, phone string) error {
	if a.transport == nil {
		return coreerrors.New(coreerrors.KindConfig, "transport is nil")
	}
	// 1) Discover DCs.
	cfg, err := transport.DiscoverDCs(ctx, a.transport)
	if err != nil {
		return err
	}
	a.dcConfig = cfg

	// 2) Load any saved credentials.
	rec, err := a.store.Load()
	if err != nil {
		return err
	}
	if rec != nil {
		if rec.Auth != "" {
			a.auth = rec.Auth
		}
		if rec.GUID != "" {
			a.guid = rec.GUID
		}
		if rec.UserAgent != "" {
			a.userAgent = rec.UserAgent
		}
		if rec.PrivateKey != "" {
			a.privateKey = rec.PrivateKey
		}
	}

	// 3) Decode auth + parse the RSA key.
	if a.auth != "" {
		a.decodedAuth = security.DecodeAuth(a.auth)
	}
	if a.privateKey != "" {
		pk, err := security.ParsePrivateKey(a.privateKey)
		if err != nil {
			return coreerrors.Wrap(coreerrors.KindAuth, "parse private key", err)
		}
		a.rsaKey = pk
	}

	// 4) Track whether we had a real auth before this point.
	hadAuth := a.auth != ""

	// Generate a tmp_session secret only if we don't have a real auth.
	if a.auth == "" {
		secret, err := security.RandomSecret(core.AuthLength)
		if err != nil {
			return err
		}
		a.auth = secret
	}
	if a.symmetricKey == "" {
		key, err := security.DeriveKey(a.auth)
		if err != nil {
			return err
		}
		a.symmetricKey = key
	}

	// 5) If we had no real auth, go straight to sign-in (skip
	// the getUserInfo probe — it will always fail with a
	// random tmp_session).
	if !hadAuth {
		if phone == "" {
			phone = a.phoneNumber
		}
		if phone == "" {
			return coreerrors.New(coreerrors.KindAuth, "not registered and no phone number provided")
		}
		if err := a.signInInteractive(ctx, phone); err != nil {
			return err
		}
	} else {
		// We have a real auth — verify it works. A non-NotRegistered
		// error (network, invalid auth, server) must NOT be swallowed:
		// proceeding would leave the account running against an invalid
		// session. Only the NotRegistered status means "this auth is
		// fresh, go sign in".
		if _, err := a.GetUserInfo(ctx, ""); err != nil {
			if !coreerrors.IsStatus(err, "NotRegistered") {
				return coreerrors.Wrap(coreerrors.KindAuth, "verify session", err)
			}
			if phone == "" {
				phone = a.phoneNumber
			}
			if phone == "" {
				return coreerrors.New(coreerrors.KindAuth, "not registered and no phone number provided")
			}
			if err := a.signInInteractive(ctx, phone); err != nil {
				return err
			}
		}
	}

	// 6) Persist the session. A save failure (disk full, permission
	// denied, ...) is surfaced to the caller rather than silently logged,
	// so a session that did not actually reach disk is never reported as
	// successfully started.
	if err := a.store.Save(storage.Record{
		Phone:      a.phoneNumber,
		Auth:       a.auth,
		GUID:       a.guid,
		UserAgent:  a.userAgent,
		PrivateKey: a.privateKey,
	}); err != nil {
		return coreerrors.Wrap(coreerrors.KindSession, "persist session", err)
	}
	return nil
}

// signInInteractive performs the sendCode → signIn handshake.
func (a *Account) signInInteractive(ctx context.Context, phone string) error {
	phone = normalizePhone(phone)
	if phone == "" {
		return coreerrors.New(coreerrors.KindValidation, "invalid phone number")
	}

	res, err := a.SendCode(ctx, phone, "", "SMS")
	if err != nil {
		return err
	}
	phoneCodeHash := res.Str("phone_code_hash")
	hint := res.Str("hint_pass_key")
	status := res.Str("status")
	// 2FA pass-key loop. Bounded to 5 attempts to prevent an
	// infinite loop on persistent wrong input.
	for attempts := 0; status == "SendPassKey" && attempts < 5; attempts++ {
		passKey, err := a.readCode(ctx, "Password ["+hint+"] > ")
		if err != nil {
			return err
		}
		res, err = a.SendCode(ctx, phone, passKey, "SMS")
		if err != nil {
			return err
		}
		status = res.Str("status")
		// Re-read phone_code_hash in case the server rotated it.
		if h := res.Str("phone_code_hash"); h != "" {
			phoneCodeHash = h
		}
	}
	if status == "SendPassKey" {
		return coreerrors.New(coreerrors.KindAuth, "too many wrong 2FA attempts")
	}

	publicKey, privateKeyPEM, err := security.GenerateKeyPair(core.RSAKeyBits)
	if err != nil {
		return err
	}
	a.privateKey = privateKeyPEM
	// Parse the newly generated private key so we can decrypt the
	// auth returned by signIn.
	newKey, err := security.ParsePrivateKey(a.privateKey)
	if err != nil {
		return coreerrors.Wrap(coreerrors.KindAuth, "parse new private key", err)
	}
	a.rsaKey = newKey
	// Code loop. Bounded to 5 attempts.
	for attempts := 0; attempts < 5; attempts++ {
		code, err := a.readCode(ctx, "Code: ")
		if err != nil {
			return err
		}
		res, err := a.SignIn(ctx, code, phone, phoneCodeHash, publicKey)
		if err != nil {
			// Allow retry on code-already-used / wrong-code by
			// re-prompting; other errors abort.
			if isRetryableCodeErr(err) && attempts < 4 {
				continue
			}
			return err
		}
		// Request() returns nil error only when status=="OK" &&
		// status_det=="OK", so sign-in succeeded. The response
		// data does not carry a nested "status" key — extract
		// the auth directly.
		encAuth := res.Str("auth")
		auth, err := security.DecryptOAEP(a.rsaKey, encAuth)
		if err != nil {
			return err
		}
		a.auth = auth
		a.symmetricKey, err = security.DeriveKey(auth)
		if err != nil {
			return err
		}
		a.decodedAuth = security.DecodeAuth(a.auth)
		if a.privateKey != "" {
			pk, err := security.ParsePrivateKey(a.privateKey)
			if err != nil {
				return err
			}
			a.rsaKey = pk
		}
		if user := res.EventOf("user"); user != nil {
			a.guid = user.Str("user_guid")
		}
		if _, err := a.RegisterDevice(ctx, a.name); err != nil {
			return err
		}
		return nil
	}
	return coreerrors.New(coreerrors.KindAuth, "too many wrong code attempts")
}

// isRetryableCodeErr reports whether a signIn error is worth retrying
// (wrong code / code already used) versus a hard failure.
func isRetryableCodeErr(err error) bool {
	msg := err.Error()
	return containsFold(msg, "CodeIsUsed") ||
		containsFold(msg, "CodeIsExpired") ||
		containsFold(msg, "InvalidInput")
}

// containsFold reports whether s contains substr case-insensitively.
func containsFold(s, substr string) bool {
	ls := strings.ToLower(s)
	return strings.Contains(ls, strings.ToLower(substr))
}

// Stop tears down the WebSocket and the transport. It is idempotent and
// signals the reconnect loop not to redial afterwards.
func (a *Account) Stop() {
	a.stopped.Store(true)
	a.wsMu.Lock()
	if a.cancel != nil {
		a.cancel()
	}
	ws := a.ws
	a.wsMu.Unlock()
	if ws != nil {
		_ = ws.Close()
	}
}

// Run is the convenience method that calls Start + Listen. It blocks
// until `ctx` is cancelled or Stop/Close is called.
func (a *Account) Run(ctx context.Context, phone string) error {
	if err := a.Start(ctx, phone); err != nil {
		return err
	}
	runCtx, cancel := context.WithCancel(ctx)
	a.wsMu.Lock()
	a.ctx, a.cancel = runCtx, cancel
	a.wsMu.Unlock()
	return a.Listen(runCtx)
}

// Listen starts the WebSocket update loop. It blocks until `ctx` is
// cancelled or a fatal error occurs.
//
// Reconnect strategy: exponential backoff with a cap and full jitter,
// reset to zero after a successful handshake. The previous fixed 3 s
// sleep could hammer the server (and burn CPU) during a persistent
// outage; the new strategy slows down (up to wsMaxBackoff) and recovers
// immediately once the connection is healthy again.
func (a *Account) Listen(ctx context.Context) error {
	if a.dcConfig == nil || a.dcConfig.WSSURL == "" {
		cfg, err := transport.DiscoverDCs(ctx, a.transport)
		if err != nil {
			return err
		}
		a.dcConfig = cfg
	}
	a.debugf("WebSocket URL = %s", a.dcConfig.WSSURL)
	ws := transport.NewWebSocket(a.dcConfig.WSSURL, nil)
	// Wire the configured proxy through to the WebSocket dialer so WS
	// traffic honours the same proxy as the HTTP path (audit: WS proxy
	// support).
	if a.proxy != "" {
		ws.SetProxy(a.proxy)
	}
	a.wsMu.Lock()
	a.ws = ws
	a.wsMu.Unlock()

	// Cancellation watcher: a blocking ws.Read() cannot be interrupted by
	// ctx alone, so on ctx cancellation (or Stop, which cancels ctx) we
	// close the socket to unblock the read loop immediately. Without this
	// the read goroutine would hang until the read deadline or a peer FIN.
	go func() {
		<-ctx.Done()
		a.wsMu.Lock()
		c := a.ws
		a.wsMu.Unlock()
		if c != nil {
			_ = c.Close()
		}
	}()

	// wsSem bounds the number of concurrently in-flight handleWSMessage
	// goroutines so a message burst cannot spawn an unbounded number of
	// goroutines (each of which decrypts + dispatches). The dispatcher
	// already bounds listener fan-out at dispatchConcurrency.
	wsSem := make(chan struct{}, wsConcurrentHandlers)

	backoff := wsInitialBackoff
	for {
		// A user-initiated Stop/Close sets `stopped`; never redial after
		// that, even on the direct-Listen path where a.cancel is nil.
		if a.stopped.Load() {
			return nil
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := ws.Dial(ctx); err != nil {
			if a.stopped.Load() {
				return nil
			}
			a.warnf("ws dial failed: %v (backoff %v)", err, backoff)
			if sleepCtxDur(ctx, backoff+jitter()) {
				return ctx.Err()
			}
			backoff = nextBackoff(backoff)
			continue
		}
		a.infof("ws connected, sending handshake...")
		if err := ws.Handshake(a.auth); err != nil {
			a.warnf("ws handshake failed: %v (backoff %v)", err, backoff)
			_ = ws.Close()
			if sleepCtxDur(ctx, backoff+jitter()) {
				return ctx.Err()
			}
			backoff = nextBackoff(backoff)
			continue
		}
		// Healthy connection: reset the backoff so the next drop does
		// not penalise a transient blip with a long sleep.
		backoff = wsInitialBackoff
		a.infof("ws handshake OK, listening for updates...")
		// Per-connection keep-alive context. Each Dial gets its own
		// cancellable context so the keep-alive goroutine exits as
		// soon as this specific WS connection drops, rather than
		// accumulating across redials (which would leak goroutines
		// and produce duplicate pings on the new connection).
		keepAliveCtx, keepAliveCancel := context.WithCancel(ctx)
		go func() {
			if err := ws.KeepAlive(keepAliveCtx, pingInterval(), func() error {
				_, err := a.GetChatsUpdates(keepAliveCtx)
				return err
			}); err != nil {
				a.debugf("ws keepalive: %v", err)
			}
		}()
		a.readWSLoop(ctx, ws, wsSem)
		keepAliveCancel() // stop the keep-alive immediately
		_ = ws.Close()
		if a.stopped.Load() {
			return nil
		}
		// After a dropped connection that had succeeded, start the
		// backoff from the initial value (transient blip) but still
		// grow it if the next dial also fails immediately.
		if sleepCtxDur(ctx, backoff+jitter()) {
			return ctx.Err()
		}
		backoff = nextBackoff(backoff)
	}
}

// readWSLoop reads messages from the WebSocket until it is closed. Each
// decoded message is handled in its own goroutine, bounded by wsSem so a
// burst of frames cannot spawn an unbounded number of goroutines.
func (a *Account) readWSLoop(ctx context.Context, ws *transport.WebSocket, wsSem chan struct{}) {
	for {
		if err := ctx.Err(); err != nil {
			return
		}
		raw, err := ws.Read()
		if err != nil {
			a.warnf("ws read error: %v", err)
			return
		}
		if a.debug {
			a.infof("ws received %d bytes", len(raw))
		}
		msg, err := transport.DecodeMessage(raw)
		if err != nil {
			a.warnf("ws decode error: %v", err)
			continue
		}
		if a.debug {
			a.infof("ws message keys: %v", mapKeys(msg))
		}
		// Acquire the semaphore before spawning so the wait happens on
		// the read loop (cheap) instead of piling up goroutines. On
		// ctx cancellation we drop the message rather than block forever.
		select {
		case wsSem <- struct{}{}:
		case <-ctx.Done():
			return
		}
		go func(m map[string]any) {
			defer func() { <-wsSem }()
			a.handleWSMessage(ctx, m)
		}(msg)
	}
}

// handleWSMessage is in handlews.go

// ===== small helpers ====================================================

// normalizePhone normalises a Persian/Iranian phone number to the
// international format (98...). It handles all common input forms:
//
// 09981530922 → 989981530922 (local, leading 0)
// 9981530922 → 989981530922 (bare mobile, no prefix)
// +989981530922 → 989981530922 (international with +)
// 00989981530922 → 989981530922 (international with 00)
// 989981530922 → 989981530922 (already normalized)
//
// Persian digits and separators (space, dash, parens) are also handled.
func normalizePhone(p string) string {
	// Convert Persian digits to ASCII.
	p = strings.NewReplacer("۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4",
		"۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9").Replace(p)
	// Strip separators.
	p = strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(p)
	// Strip leading "+".
	if strings.HasPrefix(p, "+") {
		p = p[1:]
	}
	// Strip leading "00" (international prefix).
	if strings.HasPrefix(p, "00") {
		p = p[2:]
	}
	// Local format: leading "0" → replace with "98".
	if strings.HasPrefix(p, "0") {
		p = "98" + p[1:]
	}
	// Bare Iranian mobile: 10 digits starting with "9" (e.g. 9981530922)
	// → prepend "98". Iranian mobiles without country code are exactly
	// 10 digits and always start with 9.
	if len(p) == 10 && strings.HasPrefix(p, "9") {
		p = "98" + p
	}
	// Validate length and digits.
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

// readCode obtains a code/password from the configured CodeProvider,
// or falls back to the interactive stdin prompt when no provider is
// set. The `hint` argument is a short human-readable prompt shown to
// the user (e.g. "Code:" or "Password [hint] > ").
func (a *Account) readCode(ctx context.Context, hint string) (string, error) {
	if a.codeProvider != nil {
		return a.codeProvider(ctx, hint)
	}
	return prompt(hint), nil
}

// prompt reads a single line from stdin (interactive mode). It uses
// bufio.ReadString so codes/passwords containing spaces are preserved
// (fmt.Fscanln would stop at the first whitespace and truncate them).
func prompt(msg string) string {
	fmt.Print(msg)
	s, _ := stdin().ReadString('\n')
	return strings.TrimRight(s, "\r\n")
}

func pingInterval() time.Duration {
	return durationOfSeconds(core.DefaultWebSocketPingInterval)
}

// WebSocket reconnect backoff parameters. The previous implementation
// slept a fixed 3 s on every drop, which both hammered the server during
// a persistent outage and felt sluggish on a transient blip. The values
// below produce a 1 s → 2 s → 4 s → 8 s → 16 s → 32 s → 60 s sequence
// (capped) with up to ±25 % jitter, reset to 1 s after a successful
// handshake.
const (
	wsInitialBackoff     = 1 * time.Second
	wsMaxBackoff         = 60 * time.Second
	wsConcurrentHandlers = 64 // bounds in-flight handleWSMessage goroutines
)

// nextBackoff doubles d up to wsMaxBackoff.
func nextBackoff(d time.Duration) time.Duration {
	d *= 2
	if d > wsMaxBackoff {
		d = wsMaxBackoff
	}
	return d
}

// jitter returns a random offset in [-25 %, +25 %] of the base backoff
// so a fleet of clients reconnecting after a network blip do not all
// retry simultaneously (thundering herd).
func jitter() time.Duration {
	// crypto/rand-based ±25 % of one second → [-250ms, +250ms].
	n, err := randIntn(500)
	if err != nil {
		return 0
	}
	return time.Duration(n-250) * time.Millisecond
}

func sleepCtx(ctx context.Context, seconds int) bool {
	return sleepCtxDur(ctx, time.Duration(seconds)*time.Second)
}

// sleepCtxDur sleeps for d, returning true if ctx was cancelled during
// the sleep (so the caller can bail out).
func sleepCtxDur(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		d = wsInitialBackoff
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return true
	case <-t.C:
		return false
	}
}

// randIntn returns a uniformly-distributed non-negative pseudo-random
// int in [0, n) using crypto/rand so reconnect jitter is not guessable
// from the system PRNG seed.
func randIntn(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("randIntn: n must be positive")
	}
	var b [8]byte
	if _, err := cryptoRand.Read(b[:]); err != nil {
		return 0, err
	}
	// Reduce 64-bit value mod n (slight modulo bias is acceptable for jitter).
	v := uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
	return int(v % uint64(n)), nil
}

// mapKeys returns the keys of a map for debug logging.
func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
