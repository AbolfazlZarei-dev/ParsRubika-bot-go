package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"strings"
)

// Two-step verification (2FA), recovery and phone-number change.

// TurnOffTwoStep disables 2FA with the current password.
func (a *Account) TurnOffTwoStep(ctx context.Context, password string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "turnOffTwoStep", Input: map[string]any{"password": password}, Encrypt: true})
}

// ChangePassword changes the 2FA password.
func (a *Account) ChangePassword(ctx context.Context, oldPassword, newPassword, hint string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "changePassword", Input: map[string]any{
		"password":     oldPassword,
		"new_password": newPassword,
		"new_hint":     hint,
	}, Encrypt: true})
}

// AbortTwoStepSetup aborts an in-progress 2FA setup.
func (a *Account) AbortTwoStepSetup(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "abortTwoStepSetup", Input: map[string]any{}, Encrypt: true})
}

// CheckTwoStepPasscode verifies the 2FA password without changing it.
func (a *Account) CheckTwoStepPasscode(ctx context.Context, password string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "checkTwoStepPasscode", Input: map[string]any{"password": password}, Encrypt: true})
}

// RequestRecoveryEmail sets or changes the 2FA recovery email.
func (a *Account) RequestRecoveryEmail(ctx context.Context, password, recoveryEmail string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "requestRecoveryEmail", Input: map[string]any{
		"password":       password,
		"recovery_email": recoveryEmail,
	}, Encrypt: true})
}

// VerifyRecoveryEmail verifies the recovery email with a code.
func (a *Account) VerifyRecoveryEmail(ctx context.Context, password, code string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "verifyRecoveryEmail", Input: map[string]any{
		"password": password,
		"code":     code,
	}, Encrypt: true})
}

// RequestForgetPassword initiates the 2FA password recovery flow.
func (a *Account) RequestForgetPassword(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "requestForgetPassword", Input: map[string]any{}, Encrypt: true})
}

// LoginTwoStepForgetPassword recovers 2FA via phone number.
func (a *Account) LoginTwoStepForgetPassword(ctx context.Context, phoneNumber string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "loginTwoStepForgetPassword", Input: map[string]any{
		"phone_number": normalizePhone(phoneNumber),
	}, TmpSession: true, Encrypt: true})
}

// LoginDisableTwoStep disables 2FA via recovery flow.
func (a *Account) LoginDisableTwoStep(ctx context.Context, phoneNumber, emailCode, forgetPasswordCodeHash string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "loginDisableTwoStep", Input: map[string]any{
		"phone_number":              normalizePhone(phoneNumber),
		"email_code":                emailCode,
		"forget_password_code_hash": forgetPasswordCodeHash,
	}, TmpSession: true, Encrypt: true})
}

// RequestChangePhoneNumber requests a phone number change.
func (a *Account) RequestChangePhoneNumber(ctx context.Context, newPhoneNumber string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "requestChangePhoneNumber", Input: map[string]any{
		"new_phone_number": normalizePhone(newPhoneNumber),
	}, Encrypt: true})
}

// VerifyChangePhoneNumber verifies the phone number change with a code.
func (a *Account) VerifyChangePhoneNumber(ctx context.Context, code, hash string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "verifyChangePhoneNumber", Input: map[string]any{
		"code": code,
		"hash": hash,
	}, Encrypt: true})
}

// stripLeadingAt removes a leading @ from a username.
func stripLeadingAt(s string) string { return strings.TrimPrefix(s, "@") }

// SetupTwoStepVerification configures 2FA. `recoveryEmail` is optional.
func (a *Account) SetupTwoStepVerification(ctx context.Context, password, hint, recoveryEmail string) (*events.Event, error) {
	in := map[string]any{"password": password, "hint": hint}
	if recoveryEmail != "" {
		in["recovery_email"] = recoveryEmail
	}
	return a.Request(ctx, RequestOptions{Method: "setupTwoStepVerification", Input: in, Encrypt: true})
}

// GetTwoPasscodeStatus returns the 2FA status.
func (a *Account) GetTwoPasscodeStatus(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getTwoPasscodeStatus", Input: map[string]any{}, Encrypt: true})
}
