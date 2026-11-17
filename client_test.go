package parsrubika

// Regression tests for the facade constructor. The one that matters:
// a Config without a CodeProvider must build — nil means "use the
// built-in interactive stdin prompt" (see Config.CodeProvider), and
// passing WithCodeProvider(nil) used to fail the whole New with
// "code provider cannot be nil", breaking SendText/Hangi on the
// first (no stored session) run.

import (
	"context"
	"path/filepath"
	"testing"
)

func TestNewWithNilCodeProvider(t *testing.T) {
	app, err := New(&Config{
		Phone:   "+989123456789",
		Session: filepath.Join(t.TempDir(), "sess"),
	})
	if err != nil {
		t.Fatalf("New with nil CodeProvider must not fail: %v", err)
	}
	defer app.Close()
	if app.Account == nil {
		t.Fatal("expected an account client for a Config with Phone")
	}
}

func TestNewWithCustomCodeProvider(t *testing.T) {
	app, err := New(&Config{
		Phone:        "+989123456789",
		Session:      filepath.Join(t.TempDir(), "sess"),
		CodeProvider: func(ctx context.Context, hint string) (string, error) { return "12345", nil },
	})
	if err != nil {
		t.Fatalf("New with a custom CodeProvider must not fail: %v", err)
	}
	defer app.Close()
}
