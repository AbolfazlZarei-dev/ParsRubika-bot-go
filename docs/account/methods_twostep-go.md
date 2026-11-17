# `account` — Methods twostep

_Source: [`methods_twostep.go`](../../ParsRubika/account/methods_twostep.go)_

## Symbols

### `TurnOffTwoStep`

```go
func (a *Account) TurnOffTwoStep(ctx context.Context, password string) (*events.Event, error)
```

TurnOffTwoStep disables 2FA with the current password.

### `ChangePassword`

```go
func (a *Account) ChangePassword(ctx context.Context, oldPassword, newPassword, hint string) (*events.Event, error)
```

ChangePassword changes the 2FA password.

### `AbortTwoStepSetup`

```go
func (a *Account) AbortTwoStepSetup(ctx context.Context) (*events.Event, error)
```

AbortTwoStepSetup aborts an in-progress 2FA setup.

### `CheckTwoStepPasscode`

```go
func (a *Account) CheckTwoStepPasscode(ctx context.Context, password string) (*events.Event, error)
```

CheckTwoStepPasscode verifies the 2FA password without changing it.

### `RequestRecoveryEmail`

```go
func (a *Account) RequestRecoveryEmail(ctx context.Context, password, recoveryEmail string) (*events.Event, error)
```

RequestRecoveryEmail sets or changes the 2FA recovery email.

### `VerifyRecoveryEmail`

```go
func (a *Account) VerifyRecoveryEmail(ctx context.Context, password, code string) (*events.Event, error)
```

VerifyRecoveryEmail verifies the recovery email with a code.

### `RequestForgetPassword`

```go
func (a *Account) RequestForgetPassword(ctx context.Context) (*events.Event, error)
```

RequestForgetPassword initiates the 2FA password recovery flow.

### `LoginTwoStepForgetPassword`

```go
func (a *Account) LoginTwoStepForgetPassword(ctx context.Context, phoneNumber string) (*events.Event, error)
```

LoginTwoStepForgetPassword recovers 2FA via phone number.

### `LoginDisableTwoStep`

```go
func (a *Account) LoginDisableTwoStep(ctx context.Context, phoneNumber, emailCode, forgetPasswordCodeHash string) (*events.Event, error)
```

LoginDisableTwoStep disables 2FA via recovery flow.

### `RequestChangePhoneNumber`

```go
func (a *Account) RequestChangePhoneNumber(ctx context.Context, newPhoneNumber string) (*events.Event, error)
```

RequestChangePhoneNumber requests a phone number change.

### `VerifyChangePhoneNumber`

```go
func (a *Account) VerifyChangePhoneNumber(ctx context.Context, code, hash string) (*events.Event, error)
```

VerifyChangePhoneNumber verifies the phone number change with a code.

### `SetupTwoStepVerification`

```go
func (a *Account) SetupTwoStepVerification(ctx context.Context, password, hint, recoveryEmail string) (*events.Event, error)
```

SetupTwoStepVerification configures 2FA. `recoveryEmail` is optional.

### `GetTwoPasscodeStatus`

```go
func (a *Account) GetTwoPasscodeStatus(ctx context.Context) (*events.Event, error)
```

GetTwoPasscodeStatus returns the 2FA status.
