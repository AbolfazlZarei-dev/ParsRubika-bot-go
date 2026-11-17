package flags

import (
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestGlobalAndPerUserFlags(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New(), plugins.Config{"flags": map[string]any{"beta": false, "new_ui": true}})
	env.Enable("flags")
	f, _ := plugins.Service[Flags](env.Manager, plugins.ServiceFlags)
	if f.IsEnabled("beta") || !f.IsEnabled("new_ui") || f.IsEnabled("missing") {
		t.Fatal("initial values wrong")
	}
	f.SetFor("beta", "vip", true)
	f.SetFor("new_ui", "grumpy", false)
	if !f.IsEnabledFor("beta", "vip") || f.IsEnabledFor("beta", "other") || f.IsEnabledFor("new_ui", "grumpy") || !f.IsEnabledFor("new_ui", "other") {
		t.Fatal("per-user override wrong")
	}
}
