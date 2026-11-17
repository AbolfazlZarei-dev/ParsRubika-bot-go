package i18n

import (
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestTranslationLanguageSwitchAndFallback(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New(), plugins.Config{
		"default": "fa",
		"messages": map[string]any{
			"fa": map[string]any{"hi": "سلام %s", "lang.changed": "زبان تغییر کرد", "only.fa": "فقط فارسی"},
			"en": map[string]any{"hi": "Hello %s", "lang.changed": "Language changed"},
		},
	})
	env.Enable("i18n")
	tr, _ := plugins.Service[Translator](env.Manager, plugins.ServiceI18n)
	if tr.T("u", "hi", "Ali") != "سلام Ali" {
		t.Fatalf("default language: %q", tr.T("u", "hi", "Ali"))
	}
	_ = env.Say("b0u", "u", "/lang en")
	if tr.Lang("u") != "en" || tr.T("u", "hi", "Ali") != "Hello Ali" {
		t.Fatal("language switch failed")
	}
	if tr.T("u", "only.fa") != "فقط فارسی" {
		t.Fatal("missing key must fall back to the default language")
	}
	if tr.T("u", "unknown.key") != "unknown.key" {
		t.Fatal("unknown key must be returned as is")
	}
	if tr.SetLang("u", "xx") {
		t.Fatal("unknown language accepted")
	}
	if got := env.Sent(); len(got) != 1 || got[0] != "Language changed" {
		t.Fatalf("sent = %v", got)
	}
}
