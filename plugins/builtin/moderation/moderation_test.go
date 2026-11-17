package moderation

import (
	"strings"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/access"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func rig(t *testing.T) *plugintest.Env {
	env := plugintest.New(t)
	env.Manager.MustRegister(access.New(), plugins.Config{"admins": []string{"root"}})
	env.Manager.MustRegister(New(), plugins.Config{"max_warns": 2})
	env.Enable("access", "moderation")
	env.Respond("getChatMember", func(map[string]any) map[string]any { return map[string]any{"type": "Member"} })
	return env
}

func TestOnlyAdminsCanModerate(t *testing.T) {
	env := rig(t)
	_ = env.Say("g0grp", "nobody", "/ban victim")
	if len(env.Calls("banChatMember")) != 0 {
		t.Fatal("non-admin banned someone")
	}
	_ = env.Say("g0grp", "root", "/ban victim")
	if c := env.Calls("banChatMember"); len(c) != 1 || c[0].Payload["user_id"] != "victim" {
		t.Fatalf("calls = %v", c)
	}
}

func TestWarnEscalatesToBan(t *testing.T) {
	env := rig(t)
	_ = env.Say("g0grp", "root", "/warn bob spamming")
	_ = env.Say("g0grp", "root", "/warnings bob")
	_ = env.Say("g0grp", "root", "/warn bob")
	sent := strings.Join(env.Sent(), "|")
	if !strings.Contains(sent, "اخطار 1 از 2") || !strings.Contains(sent, "spamming") {
		t.Fatalf("first warning text: %s", sent)
	}
	if len(env.Calls("banChatMember")) != 1 {
		t.Fatal("second warning must ban")
	}
	_ = env.Say("g0grp", "root", "/warnings bob")
	last := env.Sent()[len(env.Sent())-1]
	if !strings.Contains(last, "0") {
		t.Fatalf("warnings must reset after the ban: %s", last)
	}
}
