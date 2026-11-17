// Package errors defines the typed and sentinel errors used across
// ParsRubika. Every error returned by a public ParsRubika function is
// either one of the sentinels below, a *TypedError, or wraps one of
// them. Callers can therefore use `errors.Is` and `errors.As` for
// precise, non-string-based matching.
package errors

import (
	"errors"
	"fmt"
)

// Kind classifies an error so that callers can react without parsing
// the message text.
type Kind int

const (
	KindUnknown Kind = iota
	KindNetwork
	KindTimeout
	KindAuth
	KindValidation
	KindAPI
	KindUpload
	KindDownload
	KindSession
	KindConfig
	KindNotImplemented
	KindCancelled
	KindStopDispatch
)

// TypedError is the structured error returned by ParsRubika. It carries
// a Kind, a human-readable message, and an optional server-side status
// string (e.g. "InvalidAuth", "NotRegistered").
type TypedError struct {
	Kind    Kind
	Status  string
	Msg     string
	Wrapped error
}

// Error implements the error interface.
func (e *TypedError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("parsrubika[%s]: %s: %v", e.Kind, e.Msg, e.Wrapped)
	}
	if e.Status != "" {
		return fmt.Sprintf("parsrubika[%s]: %s: %s", e.Kind, e.Status, e.Msg)
	}
	return fmt.Sprintf("parsrubika[%s]: %s", e.Kind, e.Msg)
}

// Unwrap supports errors.Is / errors.As.
func (e *TypedError) Unwrap() error { return e.Wrapped }

// String implements fmt.Stringer for Kind.
func (k Kind) String() string {
	switch k {
	case KindNetwork:
		return "network"
	case KindTimeout:
		return "timeout"
	case KindAuth:
		return "auth"
	case KindValidation:
		return "validation"
	case KindAPI:
		return "api"
	case KindUpload:
		return "upload"
	case KindDownload:
		return "download"
	case KindSession:
		return "session"
	case KindConfig:
		return "config"
	case KindNotImplemented:
		return "not-implemented"
	case KindCancelled:
		return "cancelled"
	case KindStopDispatch:
		return "stop-dispatch"
	default:
		return "unknown"
	}
}

// New constructs a fresh TypedError.
func New(kind Kind, msg string) *TypedError {
	return &TypedError{Kind: kind, Msg: msg}
}

// Wrap returns a TypedError that wraps `inner` while preserving its
// kind for `errors.Is` matching.
func Wrap(kind Kind, msg string, inner error) *TypedError {
	return &TypedError{Kind: kind, Msg: msg, Wrapped: inner}
}

// WithStatus attaches a server-side status string to the error.
func WithStatus(kind Kind, status, msg string) *TypedError {
	return &TypedError{Kind: kind, Status: status, Msg: msg}
}

// ===== Sentinel errors ===================================================

// Common sentinel errors for use with `errors.Is`.
var (
	ErrNetwork         = New(KindNetwork, "network error")
	ErrTimeout         = New(KindTimeout, "request timed out")
	ErrInvalidAuth     = WithStatus(KindAuth, "InvalidAuth", "invalid auth key")
	ErrNotRegistered   = WithStatus(KindAuth, "NotRegistered", "user not registered")
	ErrCodeUsed        = WithStatus(KindAuth, "CodeIsUsed", "verification code already used")
	ErrCodeExpired     = WithStatus(KindAuth, "CodeIsExpired", "verification code expired")
	ErrTooManyRequests = WithStatus(KindAPI, "TooRequests", "too many requests")
	ErrServerError     = WithStatus(KindAPI, "ServerError", "server error")
	ErrURLNotFound     = WithStatus(KindAPI, "UrlNotFound", "url not found")
	ErrInvalidInput    = WithStatus(KindValidation, "InvalidInput", "invalid input")
	ErrNoConnection    = WithStatus(KindNetwork, "NoConnection", "no connection")
	ErrUnsupportedAPI  = WithStatus(KindAPI, "NotSupportedApiVersion", "unsupported API version")
	ErrUploadFailed    = New(KindUpload, "upload failed")
	ErrDownloadFailed  = New(KindDownload, "download failed")
	ErrSessionEmpty    = New(KindSession, "session is empty")
	ErrSessionCorrupt  = New(KindSession, "session data is corrupt")
	ErrConfigInvalid   = New(KindConfig, "configuration is invalid")
	ErrCancelled       = New(KindCancelled, "operation cancelled")
	ErrStopDispatch    = New(KindStopDispatch, "stop dispatching")
	ErrNotImplemented  = New(KindNotImplemented, "not implemented")
)

// FromStatus maps a server-side `status_det` string to the matching
// sentinel error. Unknown statuses produce a generic *TypedError with
// the original status preserved.
func FromStatus(status, msg string) error {
	if msg == "" {
		msg = status
	}
	switch status {
	case "InvalidAuth":
		return ErrInvalidAuth
	case "NotRegistered":
		return ErrNotRegistered
	case "CodeIsUsed":
		return ErrCodeUsed
	case "CodeIsExpired":
		return ErrCodeExpired
	case "TooRequests":
		return ErrTooManyRequests
	case "ServerError":
		return ErrServerError
	case "UrlNotFound":
		return ErrURLNotFound
	case "InvalidInput":
		return ErrInvalidInput
	case "NoConnection":
		return ErrNoConnection
	case "NotSupportedApiVersion":
		return ErrUnsupportedAPI
	default:
		return WithStatus(KindAPI, status, msg)
	}
}

// IsKind reports whether err is a *TypedError of the given Kind.
func IsKind(err error, k Kind) bool {
	var te *TypedError
	if errors.As(err, &te) {
		return te.Kind == k
	}
	return false
}

// IsStatus reports whether err is a *TypedError with the given status.
func IsStatus(err error, status string) bool {
	var te *TypedError
	if errors.As(err, &te) {
		return te.Status == status
	}
	return false
}
