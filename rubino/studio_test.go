package rubino

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	s := New(WithAuth("test-auth"))
	if s == nil {
		t.Fatal("expected non-nil Rubino")
	}
	defer s.Close()
	if s.auth != "test-auth" {
		t.Errorf("auth = %q", s.auth)
	}
}

func TestProfileOr(t *testing.T) {
	s := New(WithAuth("a"), WithProfile("default-profile"))
	if s.profileOr("") != "default-profile" {
		t.Error("profileOr(\"\") should return the default profile")
	}
	if s.profileOr("explicit") != "explicit" {
		t.Error("profileOr(\"explicit\") should return the explicit value")
	}
}

func TestRequestNoAuth(t *testing.T) {
	s := New()
	defer s.Close()
	_, err := s.request(context.Background(), "getProfileList", map[string]any{})
	if err == nil {
		t.Error("expected error when auth is missing")
	}
}

func TestRnd(t *testing.T) {
	a := rnd()
	b := rnd()
	// rnd() should always return a positive int64. Consecutive
	// calls should usually differ.
	if a < 0 || b < 0 {
		t.Error("rnd should be positive")
	}
}
