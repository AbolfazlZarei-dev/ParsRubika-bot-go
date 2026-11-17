// Package cron is a tiny cron-style scheduler built on top of the
// `scheduler` package. It supports the standard 5-field cron syntax
// (minute, hour, day-of-month, month, day-of-week) and a small set
// of macros: @yearly, @monthly, @weekly, @daily, @hourly.
package cron

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/scheduler"
)

// Job is a single cron entry: a schedule + a callback.
type Job struct {
	Schedule string
	Callback func(context.Context)
}

// Cronner wraps a scheduler and a list of parsed jobs. Use New to
// construct one, Add to register jobs, then Start to begin firing.
type Cronner struct {
	sched   *scheduler.Scheduler
	jobs    []*parsedJob
	running bool
	mu      sync.RWMutex
}

// New constructs a fresh Cronner bound to the supplied context.
func New(parent context.Context) *Cronner {
	return &Cronner{sched: scheduler.New(parent)}
}

// Add parses `spec` and registers `cb` for it. Supported spec forms:
//
// - "@yearly" / "@annually" — once a year, Jan 1 00:00
// - "@monthly" — once a month, day 1 00:00
// - "@weekly" — once a week, Sunday 00:00
// - "@daily" / "@midnight" — once a day, 00:00
// - "@hourly" — once an hour, at minute 0
// - 5-field cron: "M H DOM MON DOW"
// Each field can be a number, "*", "*/N" (every N units), or a
// comma-separated list of numbers.
func (c *Cronner) Add(spec string, cb func(context.Context)) error {
	j, err := parseSpec(spec)
	if err != nil {
		return err
	}
	j.callback = cb
	c.mu.Lock()
	c.jobs = append(c.jobs, j)
	c.mu.Unlock()
	return nil
}

// Start begins firing all registered jobs. Returns immediately.
func (c *Cronner) Start() {
	c.mu.Lock()
	c.running = true
	jobs := append([]*parsedJob(nil), c.jobs...)
	c.mu.Unlock()
	for _, j := range jobs {
		c.scheduleJob(j)
	}
}

// Stop cancels every scheduled job.
func (c *Cronner) Stop() {
	c.sched.Stop()
	c.mu.Lock()
	c.running = false
	c.mu.Unlock()
}

// Running reports whether the cronner is currently active.
func (c *Cronner) Running() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}

// scheduleJob computes the next fire time and registers a one-shot
// task; when the callback runs it reschedules itself.
func (c *Cronner) scheduleJob(j *parsedJob) {
	now := time.Now()
	next := j.Next(now)
	if next.IsZero() {
		return
	}
	delay := next.Sub(now)
	c.sched.After(delay, func(ctx context.Context) {
		j.callback(ctx)
		// Reschedule.
		c.scheduleJob(j)
	})
}

// parsedJob is the internal representation of a parsed cron spec.
type parsedJob struct {
	minute   field
	hour     field
	day      field
	month    field
	weekday  field
	callback func(context.Context)
}

// field is a set of allowed values for one cron field.
type field struct {
	wildcard bool
	values   map[int]struct{}
	step     int
}

// Next returns the next time after `from` that matches the schedule.
// It uses the standard cron field-by-field forward-scanning algorithm
// (the same one used by Vixie cron and robfig/cron) instead of a
// minute-by-minute brute force, so even `@yearly` resolves in O(1)
// iterations rather than ~5 million.
//
// Returns the zero time when no match exists in the next 5 years.
func (j *parsedJob) Next(from time.Time) time.Time {
	// Start at the top of the next minute after `from`.
	t := from.Truncate(time.Minute).Add(time.Minute)

	// Cap the search at 5 years to avoid infinite loops on impossible
	// specs (e.g. "0 0 31 2 *", which can never match).
	deadline := from.AddDate(5, 0, 0)

	for {
		if t.After(deadline) {
			return time.Time{}
		}

		// Month: if the current month doesn't match, advance to the
		// first second of the next matching month and restart.
		if !j.month.matches(int(t.Month())) {
			t = firstOfNextMatchingMonth(t, j.month)
			continue
		}

		// Day: check both day-of-month and day-of-week (cron uses OR
		// semantics when both are restricted, AND when one is `*`).
		if !dayMatches(j, t) {
			// Advance to the start of the next day.
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1)
			continue
		}

		// Hour.
		if !j.hour.matches(t.Hour()) {
			t = t.Add(time.Hour).Truncate(time.Hour)
			continue
		}

		// Minute.
		if !j.minute.matches(t.Minute()) {
			t = t.Add(time.Minute)
			continue
		}

		return t
	}
}

// dayMatches implements cron's day-of-month / day-of-week matching
// rules. When both fields are restricted (not wildcard), the day
// matches if EITHER matches (OR semantics). When one is a wildcard,
// both must match (AND semantics).
func dayMatches(j *parsedJob, t time.Time) bool {
	domRestricted := !j.day.wildcard
	dowRestricted := !j.weekday.wildcard
	if domRestricted && dowRestricted {
		return j.day.matches(t.Day()) || j.weekday.matches(int(t.Weekday()))
	}
	return j.day.matches(t.Day()) && j.weekday.matches(int(t.Weekday()))
}

// firstOfNextMatchingMonth returns the first day of the next month
// (after `t`) whose value matches `f`. If the current month matches,
// it returns the first day of the *current* month (the caller will
// then re-check day/hour/minute).
func firstOfNextMatchingMonth(t time.Time, f field) time.Time {
	// Walk forward month-by-month until one matches.
	year, month, _ := t.Date()
	for i := 0; i < 12; i++ {
		if f.matches(int(month)) {
			return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
		}
		// Advance to the next month.
		if month == time.December {
			year++
			month = time.January
		} else {
			month++
		}
	}
	// No match in the next 12 months — return far future so the
	// deadline check terminates the loop.
	return t.AddDate(1, 0, 0)
}

func (j *parsedJob) match(t time.Time) bool {
	if !j.minute.matches(int(t.Minute())) {
		return false
	}
	if !j.hour.matches(int(t.Hour())) {
		return false
	}
	if !j.day.matches(int(t.Day())) {
		return false
	}
	if !j.month.matches(int(t.Month())) {
		return false
	}
	if !j.weekday.matches(int(t.Weekday())) {
		return false
	}
	return true
}

func (f *field) matches(v int) bool {
	if f.wildcard {
		return true
	}
	_, ok := f.values[v]
	return ok
}

// parseSpec parses a cron spec into a parsedJob.
func parseSpec(spec string) (*parsedJob, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, fmt.Errorf("cron: empty spec")
	}
	// Macros.
	switch strings.ToLower(spec) {
	case "@yearly", "@annually":
		spec = "0 0 1 1 *"
	case "@monthly":
		spec = "0 0 1 * *"
	case "@weekly":
		spec = "0 0 * * 0"
	case "@daily", "@midnight":
		spec = "0 0 * * *"
	case "@hourly":
		spec = "0 * * * *"
	}
	parts := strings.Fields(spec)
	if len(parts) != 5 {
		return nil, fmt.Errorf("cron: spec must have 5 fields, got %d", len(parts))
	}
	j := &parsedJob{}
	var err error
	if j.minute, err = parseField(parts[0], 0, 59); err != nil {
		return nil, fmt.Errorf("cron: minute field: %w", err)
	}
	if j.hour, err = parseField(parts[1], 0, 23); err != nil {
		return nil, fmt.Errorf("cron: hour field: %w", err)
	}
	if j.day, err = parseField(parts[2], 1, 31); err != nil {
		return nil, fmt.Errorf("cron: day-of-month field: %w", err)
	}
	if j.month, err = parseField(parts[3], 1, 12); err != nil {
		return nil, fmt.Errorf("cron: month field: %w", err)
	}
	if j.weekday, err = parseField(parts[4], 0, 6); err != nil {
		return nil, fmt.Errorf("cron: day-of-week field: %w", err)
	}
	return j, nil
}

// parseField parses a single cron field. Supports "*", "*/N", "N",
// "N-M" (range), and comma-separated lists.
func parseField(s string, min, max int) (field, error) {
	f := field{values: map[int]struct{}{}}
	if s == "*" {
		f.wildcard = true
		return f, nil
	}
	for _, part := range strings.Split(s, ",") {
		if strings.HasPrefix(part, "*/") {
			stepStr := part[2:]
			step, err := strconv.Atoi(stepStr)
			if err != nil || step <= 0 {
				return f, fmt.Errorf("invalid step %q", stepStr)
			}
			for i := min; i <= max; i += step {
				f.values[i] = struct{}{}
			}
			continue
		}
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, err := strconv.Atoi(bounds[0])
			if err != nil {
				return f, fmt.Errorf("invalid range %q", part)
			}
			hi, err := strconv.Atoi(bounds[1])
			if err != nil {
				return f, fmt.Errorf("invalid range %q", part)
			}
			if lo < min || hi > max || lo > hi {
				return f, fmt.Errorf("range %q out of bounds [%d, %d]", part, min, max)
			}
			for i := lo; i <= hi; i++ {
				f.values[i] = struct{}{}
			}
			continue
		}
		v, err := strconv.Atoi(part)
		if err != nil {
			return f, fmt.Errorf("invalid value %q", part)
		}
		if v < min || v > max {
			return f, fmt.Errorf("value %d out of bounds [%d, %d]", v, min, max)
		}
		f.values[v] = struct{}{}
	}
	return f, nil
}
