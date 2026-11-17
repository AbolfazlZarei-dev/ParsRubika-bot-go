package broadcast

import (
	"strings"
	"testing"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/access"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/users"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestBroadcastReachesEveryKnownChatAndIsAdminOnly(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(access.New(), plugins.Config{"admins": []string{"root"}})
	env.Manager.MustRegister(users.New())
	env.Manager.MustRegister(New(), plugins.Config{"delay": "1ms"})
	env.Enable("broadcast") // pulls in access + users through Requires
	_ = env.Say("b0a", "a", "hi")
	_ = env.Say("b0b", "b", "hi")
	env.Reset()

	_ = env.Say("b0a", "a", "/broadcast nope") // not an admin
	if len(env.Calls("sendMessage")) != 0 {
		t.Fatal("non-admin could broadcast")
	}
	_ = env.Say("b0root", "root", "/broadcast hello all")
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		n := 0
		for _, s := range env.Sent() {
			if s == "hello all" {
				n++
			}
		}
		if n == 3 { // a, b and the admin's own chat
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	n, summary := 0, false
	for _, s := range env.Sent() {
		if s == "hello all" {
			n++
		}
		if strings.HasPrefix(s, "✅") {
			summary = true
		}
	}
	if n != 3 || !summary {
		t.Fatalf("delivered=%d summary=%v sent=%v", n, summary, env.Sent())
	}
}
