package core

import (
	"strings"
	"testing"

	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/enums"
)

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if Author == "" {
		t.Error("Author should not be empty")
	}
}

func TestWelcomeBanner(t *testing.T) {
	b := WelcomeBanner()
	if !strings.Contains(b, Version) {
		t.Errorf("WelcomeBanner should contain version: %s", b)
	}
	if !strings.Contains(b, Author) {
		t.Errorf("WelcomeBanner should contain author: %s", b)
	}
}

func TestWebIdentity(t *testing.T) {
	id := WebIdentity("fa")
	if id.Platform != "PWA" {
		t.Errorf("Platform = %q, want PWA", id.Platform)
	}
	if id.LangCode != "fa" {
		t.Errorf("LangCode = %q, want fa", id.LangCode)
	}
}

func TestAndroidIdentity(t *testing.T) {
	id := AndroidIdentity("fa")
	if id.Platform != "Android" {
		t.Errorf("Platform = %q, want Android", id.Platform)
	}
}

func TestChatTypeFromGUID(t *testing.T) {
	cases := map[string]enums.ChatType{
		"c0CHANNEL": enums.ChatTypeChannel,
		"g0GROUP":   enums.ChatTypeGroup,
		"u0USER":    enums.ChatTypeUser,
	}
	for in, want := range cases {
		if got := enums.ChatTypeFromGUID(in); got != want {
			t.Errorf("ChatTypeFromGUID(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTypedError(t *testing.T) {
	e := coreerrors.New(coreerrors.KindNetwork, "boom")
	if e.Error() == "" {
		t.Error("Error() should not be empty")
	}
	if !coreerrors.IsKind(e, coreerrors.KindNetwork) {
		t.Error("IsKind should return true")
	}
	if coreerrors.IsKind(e, coreerrors.KindAuth) {
		t.Error("IsKind should return false for non-matching kind")
	}
}

func TestFromStatus(t *testing.T) {
	e := coreerrors.FromStatus("InvalidAuth", "bad key")
	if !coreerrors.IsStatus(e, "InvalidAuth") {
		t.Error("IsStatus should return true for InvalidAuth")
	}
	e2 := coreerrors.FromStatus("SomeUnknownStatus", "msg")
	if !coreerrors.IsStatus(e2, "SomeUnknownStatus") {
		t.Error("IsStatus should return true for unknown status")
	}
}

func TestSentinels(t *testing.T) {
	if coreerrors.ErrInvalidAuth == nil {
		t.Error("ErrInvalidAuth should not be nil")
	}
	if !coreerrors.IsStatus(coreerrors.ErrInvalidAuth, "InvalidAuth") {
		t.Error("ErrInvalidAuth should be status InvalidAuth")
	}
	if !coreerrors.IsStatus(coreerrors.ErrNotRegistered, "NotRegistered") {
		t.Error("ErrNotRegistered should be status NotRegistered")
	}
}
