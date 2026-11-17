package kv

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestTTLAndPrefix(t *testing.T) {
	env := plugintest.New(t)
	p := New()
	clock := time.Unix(100, 0)
	p.now = func() time.Time { return clock }
	env.Manager.MustRegister(p)
	env.Enable("kv")
	s, ok := plugins.Service[Store](env.Manager, plugins.ServiceKV)
	if !ok {
		t.Fatal("service missing")
	}
	s.Set("user:1", "a")
	s.SetTTL("user:2", "b", time.Minute)
	s.Set("other", 1)
	if got := s.Keys("user:"); len(got) != 2 {
		t.Fatalf("keys = %v", got)
	}
	clock = clock.Add(2 * time.Minute)
	if _, ok := s.Get("user:2"); ok {
		t.Fatal("expired key still readable")
	}
	if got := s.Keys("user:"); len(got) != 1 {
		t.Fatalf("keys after expiry = %v", got)
	}
	s.Delete("user:1")
	if _, ok := s.Get("user:1"); ok {
		t.Fatal("delete failed")
	}
}

func TestPersistenceAcrossRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kv.json")
	env := plugintest.New(t)
	env.Manager.MustRegister(New(), plugins.Config{"path": path})
	env.Enable("kv")
	s, _ := plugins.Service[Store](env.Manager, plugins.ServiceKV)
	s.Set("k", "v")
	if err := env.Manager.Disable(nil2(), "kv"); err != nil { // Teardown flushes
		t.Fatal(err)
	}
	env2 := plugintest.New(t)
	env2.Manager.MustRegister(New(), plugins.Config{"path": path})
	env2.Enable("kv")
	s2, _ := plugins.Service[Store](env2.Manager, plugins.ServiceKV)
	if v, ok := s2.Get("k"); !ok || v != "v" {
		t.Fatalf("value lost across restart: %v %v", v, ok)
	}
}
