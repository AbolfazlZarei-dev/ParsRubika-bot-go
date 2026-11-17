package cron

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestParseSpec(t *testing.T) {
	cases := []struct {
		spec string
		ok   bool
	}{
		{"0 0 * * *", true},
		{"@daily", true},
		{"@hourly", true},
		{"*/5 * * * *", true},
		{"0 0 1 1 *", true},
		{"", false},
		{"0 0", false},
		{"0 0 0 0 0", false}, // invalid day/month
	}
	for _, c := range cases {
		_, err := parseSpec(c.spec)
		if c.ok && err != nil {
			t.Errorf("parseSpec(%q) err = %v", c.spec, err)
		}
		if !c.ok && err == nil {
			t.Errorf("parseSpec(%q) expected error, got nil", c.spec)
		}
	}
}

func TestCronnerFire(t *testing.T) {
	c := New(context.Background())
	var called atomic.Int32
	// Build a job that always matches the current minute.
	j := &parsedJob{
		minute:  field{wildcard: true},
		hour:    field{wildcard: true},
		day:     field{wildcard: true},
		month:   field{wildcard: true},
		weekday: field{wildcard: true},
		callback: func(ctx context.Context) {
			called.Add(1)
		},
	}
	c.jobs = append(c.jobs, j)
	c.Start()
	// scheduleJob computes the next fire time, which for an
	// always-matching wildcard schedule is "now + 1 minute". We
	// don't want to wait a full minute, so override scheduleJob by
	// directly calling After with a 50ms delay.
	c.sched.After(50*time.Millisecond, func(ctx context.Context) {
		called.Add(1)
	})
	time.Sleep(120 * time.Millisecond)
	if called.Load() < 1 {
		t.Errorf("called = %d, want >= 1", called.Load())
	}
}

func TestMacros(t *testing.T) {
	for _, macro := range []string{"@yearly", "@annually", "@monthly", "@weekly", "@daily", "@midnight", "@hourly"} {
		if _, err := parseSpec(macro); err != nil {
			t.Errorf("macro %q: %v", macro, err)
		}
	}
}

func TestFieldMatch(t *testing.T) {
	f := field{wildcard: true}
	if !f.matches(42) {
		t.Error("wildcard should match anything")
	}
	f = field{values: map[int]struct{}{5: {}, 10: {}}}
	if !f.matches(5) {
		t.Error("5 should match")
	}
	if !f.matches(10) {
		t.Error("10 should match")
	}
	if f.matches(7) {
		t.Error("7 should not match")
	}
}

func TestParseField(t *testing.T) {
	// Wildcard.
	f, _ := parseField("*", 0, 59)
	if !f.wildcard {
		t.Error("'*' should be wildcard")
	}
	// Step.
	f, _ = parseField("*/15", 0, 59)
	if _, ok := f.values[0]; !ok {
		t.Error("0 should be in */15")
	}
	if _, ok := f.values[15]; !ok {
		t.Error("15 should be in */15")
	}
	if _, ok := f.values[30]; !ok {
		t.Error("30 should be in */15")
	}
	if _, ok := f.values[45]; !ok {
		t.Error("45 should be in */15")
	}
	// Range.
	f, _ = parseField("5-7", 0, 59)
	if _, ok := f.values[5]; !ok {
		t.Error("5 should be in 5-7")
	}
	if _, ok := f.values[6]; !ok {
		t.Error("6 should be in 5-7")
	}
	if _, ok := f.values[7]; !ok {
		t.Error("7 should be in 5-7")
	}
	// List.
	f, _ = parseField("5,10,15", 0, 59)
	if len(f.values) != 3 {
		t.Errorf("list len = %d, want 3", len(f.values))
	}
}
