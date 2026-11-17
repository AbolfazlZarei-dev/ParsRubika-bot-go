package plugins

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
)

type quiet struct{}

func (quiet) Printf(string, ...any) {}

func newBot() *bot.Bot {
	return bot.New("tok", bot.WithLogger(quiet{}), bot.WithPersistOffset(false))
}

// probe is a configurable test plugin.
type probe struct {
	Base
	setup    func(ctx context.Context, h *Host) error
	torn     *atomic.Int32
	defaults Config
}

func newProbe(name string, req ...string) *probe {
	return &probe{Base: NewBase(Meta{Name: name, Version: "1", Requires: req}), torn: new(atomic.Int32)}
}
func (p *probe) Setup(ctx context.Context, h *Host) error {
	if p.setup != nil {
		return p.setup(ctx, h)
	}
	return nil
}
func (p *probe) Teardown(context.Context) error { p.torn.Add(1); return nil }
func (p *probe) DefaultConfig() Config          { return p.defaults }

func TestEnableDisableRemovesEverythingFromBot(t *testing.T) {
	b := newBot()
	m := New(b)
	p := newProbe("a")
	p.setup = func(ctx context.Context, h *Host) error {
		h.Command("x", func(c *bot.Ctx) error { return nil })
		h.Use(func(ctx context.Context, b *bot.Bot, u *bot.Update, next func() error) error { return next() })
		h.OnStart(func(context.Context, *bot.Bot) error { return nil })
		h.Help("/x", "does x")
		return h.Provide("svc.a", 42)
	}
	m.MustRegister(p)
	if err := m.Enable(context.Background(), "a"); err != nil {
		t.Fatal(err)
	}
	if !m.IsEnabled("a") || len(m.HelpEntries()) != 1 {
		t.Fatal("plugin not active")
	}
	if v, ok := Service[int](m, "svc.a"); !ok || v != 42 {
		t.Fatal("service missing")
	}
	if m.List()[0].Handlers != 1 {
		t.Fatalf("handlers = %d", m.List()[0].Handlers)
	}
	if err := m.Disable(context.Background(), "a"); err != nil {
		t.Fatal(err)
	}
	if m.IsEnabled("a") || p.torn.Load() != 1 {
		t.Fatal("plugin still active or teardown not called")
	}
	if _, ok := m.Service("svc.a"); ok || len(m.HelpEntries()) != 0 {
		t.Fatal("service/help not withdrawn")
	}
	// the bot must have no handler left: Off of any former ID fails
	for id := 1; id < 8; id++ {
		if b.Off(id) || b.RemoveMiddleware(id) || b.RemoveHook(id) {
			t.Fatalf("leftover registration with id %d", id)
		}
	}
}

func TestDependenciesEnabledFirstAndDisabledLast(t *testing.T) {
	m := New(newBot())
	var order []string
	mk := func(name string, req ...string) *probe {
		p := newProbe(name, req...)
		p.setup = func(context.Context, *Host) error { order = append(order, "+"+name); return nil }
		return p
	}
	m.MustRegister(mk("app", "db", "cache"))
	m.MustRegister(mk("db"))
	m.MustRegister(mk("cache", "db"))
	if err := m.Enable(context.Background(), "app"); err != nil {
		t.Fatal(err)
	}
	if strings.Join(order, ",") != "+db,+cache,+app" {
		t.Fatalf("order = %v", order)
	}
	// disabling db must take its dependents down first
	if err := m.Disable(context.Background(), "db"); err != nil {
		t.Fatal(err)
	}
	if m.IsEnabled("app") || m.IsEnabled("cache") || m.IsEnabled("db") {
		t.Fatalf("enabled: %v", m.Enabled())
	}
}

func TestCycleAndUnknownDependency(t *testing.T) {
	m := New(newBot())
	m.MustRegister(newProbe("a", "b"))
	m.MustRegister(newProbe("b", "a"))
	m.MustRegister(newProbe("c", "ghost"))
	if err := m.Enable(context.Background(), "a"); !errors.Is(err, ErrDependency) {
		t.Fatalf("cycle: %v", err)
	}
	if err := m.Enable(context.Background(), "c"); !errors.Is(err, ErrDependency) {
		t.Fatalf("unknown dep: %v", err)
	}
	if err := m.Enable(context.Background(), "nope"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("not found: %v", err)
	}
}

func TestFailedSetupRollsBack(t *testing.T) {
	b := newBot()
	m := New(b)
	boom := errors.New("boom")
	p := newProbe("bad")
	p.setup = func(ctx context.Context, h *Host) error {
		h.Command("x", func(c *bot.Ctx) error { return nil })
		return boom
	}
	m.MustRegister(p)
	err := m.Enable(context.Background(), "bad")
	var se *SetupError
	if !errors.As(err, &se) || !errors.Is(err, boom) {
		t.Fatalf("err = %v", err)
	}
	if m.IsEnabled("bad") || b.Off(1) {
		t.Fatal("handler of failed plugin was not removed")
	}
	if m.List()[0].Error == "" {
		t.Fatal("error not reported in Info")
	}
}

func TestPanicInSetupIsContained(t *testing.T) {
	m := New(newBot())
	p := newProbe("p")
	p.setup = func(context.Context, *Host) error { panic("oops") }
	m.MustRegister(p)
	if err := m.Enable(context.Background(), "p"); err == nil || !strings.Contains(err.Error(), "oops") {
		t.Fatalf("err = %v", err)
	}
}

func TestConfigDefaultsMergeAndTypedGetters(t *testing.T) {
	m := New(newBot(), WithConfigs(map[string]Config{"p": {"n": 7.0}}))
	p := newProbe("p")
	p.defaults = Config{"n": 1, "name": "x", "d": "2s"}
	var got Config
	p.setup = func(_ context.Context, h *Host) error { got = h.Config(); return nil }
	m.MustRegister(p, Config{"name": "y"})
	if err := m.Enable(context.Background(), "p"); err != nil {
		t.Fatal(err)
	}
	if got.Int("n", 0) != 7 || got.String("name", "") != "y" || got.Duration("d", 0) != 2*time.Second {
		t.Fatalf("cfg = %v", got)
	}
}

func TestEveryStopsOnDisable(t *testing.T) {
	m := New(newBot())
	var ticks atomic.Int32
	p := newProbe("tick")
	p.setup = func(_ context.Context, h *Host) error {
		h.Every(5*time.Millisecond, func(context.Context) { ticks.Add(1) })
		return nil
	}
	m.MustRegister(p)
	_ = m.Enable(context.Background(), "tick")
	time.Sleep(40 * time.Millisecond)
	_ = m.Disable(context.Background(), "tick")
	n := ticks.Load()
	time.Sleep(30 * time.Millisecond)
	if n == 0 || ticks.Load() != n {
		t.Fatalf("ticks before=%d after=%d", n, ticks.Load())
	}
}

func TestEnableAllHonoursPriorityAndDisabledFlag(t *testing.T) {
	m := New(newBot())
	var order []string
	mk := func(name string, prio int, off bool) *probe {
		p := &probe{Base: NewBase(Meta{Name: name, Priority: prio, Disabled: off}), torn: new(atomic.Int32)}
		p.setup = func(context.Context, *Host) error { order = append(order, name); return nil }
		return p
	}
	m.MustRegister(mk("late", 200, false))
	m.MustRegister(mk("early", 10, false))
	m.MustRegister(mk("off", 1, true))
	if err := m.EnableAll(context.Background()); err != nil {
		t.Fatal(err)
	}
	if strings.Join(order, ",") != "early,late" {
		t.Fatalf("order = %v", order)
	}
	_ = m.DisableAll(context.Background())
	if len(m.Enabled()) != 0 {
		t.Fatal("DisableAll left plugins active")
	}
}

func TestRegisterValidation(t *testing.T) {
	m := New(newBot())
	if err := m.Register(newProbe("Bad Name")); !errors.Is(err, ErrInvalidName) {
		t.Fatalf("invalid name: %v", err)
	}
	m.MustRegister(newProbe("ok"))
	if err := m.Register(newProbe("OK")); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate: %v", err)
	}
}

func TestReloadAppliesNewConfig(t *testing.T) {
	m := New(newBot())
	p := newProbe("r")
	var v int
	p.setup = func(_ context.Context, h *Host) error { v = h.Config().Int("n", 0); return nil }
	m.MustRegister(p, Config{"n": 1})
	_ = m.Enable(context.Background(), "r")
	if err := m.Reload(context.Background(), "r", Config{"n": 5}); err != nil {
		t.Fatal(err)
	}
	if v != 5 || p.torn.Load() != 1 {
		t.Fatalf("v=%d torn=%d", v, p.torn.Load())
	}
}

func TestParseConfigJSON(t *testing.T) {
	c, err := ParseConfigJSON([]byte(`{"plugins":{"AntiSpam":{"max_messages":5,"words":["a","b"]}}}`))
	if err != nil {
		t.Fatal(err)
	}
	cfg := c["antispam"]
	if cfg.Int("max_messages", 0) != 5 || len(cfg.Strings("words")) != 2 {
		t.Fatalf("cfg = %v", cfg)
	}
}
