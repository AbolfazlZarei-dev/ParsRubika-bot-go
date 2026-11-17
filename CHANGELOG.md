# Changelog

## 4.0.0-r2 — "Everything in one import" (2026-11-14)

### Added
- **parsrubika (root)**: the everything-in-one facade. `New(&Config{...})`
  builds the account, bot, rubino, tv and rubiai clients together behind
  a single `App`; `Login` signs in and wires the account auth into
  Rubino/Rubiai automatically, `Run` stays online in the WebSocket update
  loop. Three one-line helpers cover the common jobs: `SendText` (one
  message), `Hangi` (the dual-avatar profile), `QuickBot` (a running bot
  with a catch-all handler). Users no longer need pages of setup code —
  the library does the wiring.
- **account**: `HangiProfile` — the dual-avatar پروفایل هنگی in ONE
  call. It uploads both images with a new tolerant image upload
  (`Account.UploadImage`: accepts every upload-service response shape
  that yields an access_hash_rec, with bounded ERROR_TRY_AGAIN
  restarts), fires ONE round of two simultaneous `uploadAvatar` binds
  (photo1 main + photo2 thumbnail, and the swapped pair), then deletes
  every avatar the run did not create so the object ends with exactly
  the two pictures — reruns never pile up avatar entries. Options:
  `WithHangiRounds` (default 1), `WithHangiOffset` (300ms),
  `WithHangiGap`, `WithHangiClean`, `WithHangiKeepOld`,
  `WithHangiProgress`. Works on the account itself and on any
  group/channel GUID (`object_guid`).
- **docs**: the whole `docs/` tree regenerated from the current sources
  (67 packages, including the new root facade and the hangi methods),
  plus README/CHANGELOG updates for everything above.

### Changed
- **examples/hangprofile**: reduced from ~500 lines to a thin flag
  wrapper around one library call. New `--guid` flag targets a
  group/channel (e.g. `--guid c0DBedk0525ec07fb54a35b01baf2303`);
  `--rounds` now defaults to **1** — the previous 5-round loop was why
  the server kept several avatar entries; the "hang" comes from the two
  writes colliding, not from repetition.

### Fixed
- **parsrubika (root facade)**: `New` no longer fails with
  "code provider cannot be nil" when `Config.CodeProvider` was left
  nil. A nil provider means "use the built-in interactive stdin
  prompt" (documented on `Config`), so the facade now only installs
  the option when the caller supplied one. This broke the very first
  run of `SendText`/`Hangi` (no stored session yet → SMS login
  required → immediate validation error). Regression tests:
  `client_test.go`.
- The hangprofile flow no longer accumulates avatars: after the single
  dual bind, leftover avatars (previous runs) are pruned, so the object
  ends with exactly the two requested photos.

---

## 4.0.0 — "Nova"

A major-version release that lands the findings of the full
security/architecture audit and adds the account-side "dual profile"
helper. Every change is backed by a regression/race test. No public
API is broken; all additions are opt-in (`WithPassphrase`,
`ToSQLiteWithPassphrase`, `WebSocket.SetProxy`, `PostRetryPolicy`,
`netutil.ControlSSRF`, `HangProfile`).

### Added
- **account**: new `HangProfile` method (the "dual profile" / "profile
  hang" trick). It fires two `updateProfile` requests back-to-back
  against the current account so the two writes race on the server;
  the persisted state can then render as a stacked/"dual" profile in
  Rubika clients.
  - Signature:
    `func (a *Account) HangProfile(ctx context.Context, firstName, lastName, bio, birthDate string) (*HangProfileResult, error)`
  - Both shots use the same payload (mirroring `UpdateProfile`); the
    "hang" effect comes from the race itself, not conflicting values.
  - Returns a `HangProfileResult{First, FirstErr, Second, SecondErr}`
    so the caller can inspect each shot independently; a single
    failed shot does not raise — its error goes on the result.
  - Honours context cancellation: cancelling the context aborts any
    in-flight shot.
  - Backed by regression + race tests
    (`TestHangProfileFiresTwoShots`, `TestHangProfileSurfacesErrors`,
    `TestHangProfileRespectsContextCancellation`).

## [4.0.0] — 2026-11-13 (security & concurrency audit release)

This release lands the findings of the full security/architecture
audit. Every change is backed by a regression/race test. No public
API is broken; all additions are opt-in (`WithPassphrase`,
`ToSQLiteWithPassphrase`, `WebSocket.SetProxy`, `PostRetryPolicy`,
`netutil.ControlSSRF`).

### Security
- **security**: removed modulo bias from `RandomSecret` (now uses
  `crypto/rand.Int`); `pkcs7Pad` no longer mutates the caller's
  backing array.
- **storage**: SQLite session file forced to mode `0600` with
  `journal_mode=DELETE` (no world-readable WAL sidecars); session
  selection is deterministic (`ORDER BY rowid DESC`); FileStore
  writes are now atomic+fsync'd via a random temp file (no
  symlink/TOCTOU race, no fixed `.tmp` path); base64 decode of
  sealed/string-session fields is size-capped before allocation;
  `StringSession.ToSQLiteWithPassphrase` preserves encryption on
  export (the old `ToSQLite` silently downgraded to plain text).
- **transport**: POST is no longer retried on ambiguous transport
  errors (prevents duplicate messages) — only on explicit
  server-signalled statuses (429/5xx/408/425); response bodies are
  capped at 16 MiB (`MaxResponseBody`); proxy setup is fail-closed
  (`ErrProxyUnavailable`, never falls back to a direct connection);
  HTTP redirects are disabled on the API client so the auth-bearing
  JSON body is never re-sent to a different host; DC API URLs are
  validated (https-only, no private/loopback literals) and WSS URLs
  validated (wss-only, no private literals); error strings are
  truncated and scrubbed of PEM/secret material; `MaxConnsPerHost`
  bounds live connections; the `headers` map is now mutex-guarded.
- **websocket**: `SetReadLimit(8 MiB)` + a 120 s read deadline; proxy
  support (`SetProxy`) for HTTP/HTTPS/SOCKS5.
- **account**: public `Uploader`/`Downloader` are no longer built
  with a nil HTTP client (was a nil-panic); `RequestOptions.Auth` is
  honoured on the v6 path; `Start` now returns (not swallows) the
  `GetUserInfo` and session-`Save` errors.
- **files**: download writes to a temp file (mode `0600`) and renames
  atomically; chunk length is verified (truncated chunks rejected);
  final size is verified; partial files are cleaned up on error;
  `file_name` sent to the server is `filepath.Base` only (local path
  no longer leaked); both `Path`+`Bytes` set is an error; chunk size
  is capped at 8 MiB; in-memory download paths refuse files > 50 MiB;
  `Content-Disposition` filenames are sanitised (path-traversal /
  control-char safe); the MIME→extension rename can no longer escape
  the target directory.
- **http3**: response body capped; redirects disabled; errors
  scrubbed.
- **plugins/antilink**: allowlist match is now exact-host or
  real-subdomain (`host == d || HasSuffix(host, "."+d)`), closing the
  `strings.Contains` bypass (`example.com.evil.com`).

### Concurrency / lifecycle
- **account**: `ws`/`cancel` fields are mutex-guarded; a `stopped`
  flag prevents reconnect after `Stop`/`Close` (also on the direct
  `Listen` path); a ctx-watcher closes the socket so a blocking
  `Read` unblocks immediately on cancellation; `MessageQueue.Submit`
  can no longer panic on send-on-closed-channel (mutex + `closed`
  flag).
- **transport**: `SetHeader`/`applyHeaders` race fixed; retry
  backoff now has ±25 % jitter.
- **core/guard**: `Retry` returns a real error for `max < 1` /
  `nil fn` instead of a silent `nil`.

### Bug fixes
- **bot**: private-chat detection used the wrong GUID prefix.
  `models.Update.IsPrivateChat()` and `bot.chatTypeOf` (used by
  `ChatTypeMiddleware`) matched the `"b0"` prefix, but `"b0"` is the
  **bot** GUID prefix — real one-to-one (user-to-bot) chats use
  `"u0"`, as `core/guid.Detect`, `rules.Private`, `rules.ChatTypeIs`
  and `bot.IsPrivateUpdate` all already assumed. The mismatch left
  `Ctx.IsPrivate()`, `rules.PrivateChat{}` and the `access`,
  `antilink`, `forcejoin` plugins never matching real private chats
  while misclassifying bot GUIDs as private. Both sites now check
  `"u0"`; the misleading doc comments are corrected and a regression
  test locks the prefix semantics (`u0`→private, `b0`→not).

### Tests / CI
- New regression + race tests for every fix above (security/,
  storage/, transport/, account/, files/, sqlitedriver/, bot/).
- The `sqlitedriver` sub-module now actually runs its tests (links
  `modernc.org/sqlite`); CI gained a `go test -race` step for it.

### Docs
- Root package doc lists `rubiai` (the RubiAI SSE client) among the
  five clients.

