package account

import (
	"context"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account/matchers"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

func TestAccountNew(t *testing.T) {
	a, err := New(WithAuth("abcdefghijklmnopqrstuvwxyz012345"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil Account")
	}
	if a.Auth() == "" {
		t.Error("Auth should not be empty")
	}
	defer a.Close()
}

func TestAccountInvalidAuth(t *testing.T) {
	if _, err := New(WithAuth("")); err == nil {
		t.Error("expected error on empty auth")
	}
}

func TestAccountOptions(t *testing.T) {
	a, err := New(
		WithAuth("abcdefghijklmnopqrstuvwxyz012345"),
		WithMaxRetries(10),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer a.Close()
	if a.maxRetries != 10 {
		t.Errorf("maxRetries = %d, want 10", a.maxRetries)
	}
}

func TestAddRemoveListener(t *testing.T) {
	a, _ := New(WithAuth("abcdefghijklmnopqrstuvwxyz012345"))
	defer a.Close()
	// Use sequential mode from the start so the test is deterministic.
	a.sequential = true
	called := 0
	id := AddListener(a, func(ctx context.Context, e *events.Event) error {
		called++
		return nil
	}, matchers.NewIsText("MessageUpdates"))
	if id < 1 {
		t.Fatalf("expected positive ID, got %d", id)
	}
	// Dispatch a fake update using the snake_case name the server
	// actually sends (`message_updates` -> capitalize -> "MessageUpdates").
	a.dispatch(context.Background(), "message_updates", map[string]any{
		"message": map[string]any{"text": "hello"},
	})
	if called == 0 {
		t.Error("listener was not called")
	}
	if !a.RemoveListener(id) {
		t.Error("RemoveListener should return true for existing ID")
	}
	if a.RemoveListener(id) {
		t.Error("RemoveListener should return false after removal")
	}
}

func TestCapitalize(t *testing.T) {
	cases := map[string]string{
		"messages":             "Messages",
		"message_updates":      "MessageUpdates",
		"chat_updates":         "ChatUpdates",
		"show_activities":      "ShowActivities",
		"show_notifications":   "ShowNotifications",
		"remove_notifications": "RemoveNotifications",
	}
	for in, want := range cases {
		if got := capitalize(in); got != want {
			t.Errorf("capitalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNormalizeSelf(t *testing.T) {
	cases := map[string]string{
		"me":      "u0SELF",
		"cloud":   "u0SELF",
		"self":    "u0SELF",
		"g0OTHER": "g0OTHER",
	}
	for in, want := range cases {
		if got := normalizeSelf(in, "u0SELF"); got != want {
			t.Errorf("normalizeSelf(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNormalizePhone(t *testing.T) {
	cases := map[string]string{
		"989121234567":   "989121234567",
		"09121234567":    "989121234567",
		"00989121234567": "989121234567",
		"":               "",
	}
	for in, want := range cases {
		if got := normalizePhone(in); got != want {
			t.Errorf("normalizePhone(%q) = %q, want %q", in, got, want)
		}
	}
}
