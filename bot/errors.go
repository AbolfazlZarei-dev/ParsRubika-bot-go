package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
)

// APIError is returned when the Bot API answers with a non-OK status.
type APIError struct {
	Status     string
	DevMessage string
	// HTTPStatus is the transport level status code (0 when unknown).
	HTTPStatus int
	// RetryAfter is the server suggested delay, when provided.
	RetryAfter time.Duration
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.DevMessage != "" {
		return "bot: " + e.Status + ": " + e.DevMessage
	}
	return "bot: " + e.Status
}

// Temporary reports whether retrying the call might succeed.
func (e *APIError) Temporary() bool {
	if e.HTTPStatus == 429 || e.HTTPStatus >= 500 {
		return true
	}
	return isRetryableAPIStatus(e.Status)
}

// IsStatus reports whether err (or any error it wraps) is an *APIError
// carrying the given status string.
func IsStatus(err error, status string) bool {
	var ae *APIError
	return errors.As(err, &ae) && ae.Status == status
}

// IsRateLimited reports whether err is a "too many requests" failure.
func IsRateLimited(err error) bool {
	var ae *APIError
	if !errors.As(err, &ae) {
		return false
	}
	return ae.HTTPStatus == 429 || ae.Status == "TooRequests" || ae.Status == "429"
}

// NewAPIError is a convenience constructor for callers that want to
// synthesise an error from a custom status string.
func NewAPIError(status, devMessage string) *APIError {
	return &APIError{Status: status, DevMessage: devMessage}
}

// ErrContinue can be returned by a handler to let the dispatcher keep
// looking for the next matching handler instead of stopping at this one.
var ErrContinue = errors.New("bot: continue to next handler")

// ErrStopped is returned by operations attempted after the bot stopped.
var ErrStopped = errors.New("bot: stopped")

// ErrorContext describes a failure raised while processing an update.
type ErrorContext struct {
	Ctx    context.Context
	Bot    *Bot
	Update *Update
	Err    error
	// Panic holds the recovered value when the failure was a panic.
	Panic any
	// Stack is the goroutine stack captured at the time of a panic.
	Stack []byte
}

// Error summarises the failure.
func (e *ErrorContext) Error() string {
	if e.Panic != nil {
		return fmt.Sprintf("panic: %v", e.Panic)
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

// Logger is the logging contract used by the bot (see core.Logger).
type Logger = core.Logger

// stdLogger adapts *log.Logger to Logger (kept for clarity).
type stdLogger struct{ l *log.Logger }

func (s stdLogger) Printf(format string, args ...any) { s.l.Printf(format, args...) }
