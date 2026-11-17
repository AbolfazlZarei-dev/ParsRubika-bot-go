package account

import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"

// Levelled logging helpers. Verbose diagnostics (decrypted payloads,
// access hashes, upload internals) are Debug only so they never reach
// the console unless the caller opts in.
func (a *Account) debugf(format string, v ...any) {
	core.Logf(a.logger, core.LogDebug, "account: "+format, v...)
}
func (a *Account) infof(format string, v ...any) {
	core.Logf(a.logger, core.LogInfo, "account: "+format, v...)
}
func (a *Account) warnf(format string, v ...any) {
	core.Logf(a.logger, core.LogWarn, "account: "+format, v...)
}
func (a *Account) errorf(format string, v ...any) {
	core.Logf(a.logger, core.LogError, "account: "+format, v...)
}
