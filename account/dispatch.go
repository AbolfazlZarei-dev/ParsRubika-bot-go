package account

import (
	"context"
	"errors"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// dispatchConcurrency caps the number of listener goroutines that may
// run concurrently for a single update. Without a cap, a 10k-message
// flood with N listeners would spawn 10k*N goroutines and exhaust
// scheduler/memory. The default is intentionally generous (listener
// callbacks are usually I/O-bound) but bounded.
const dispatchConcurrency = 64

// dispatch invokes every registered listener whose matcher matches
// the incoming update. Listeners run either sequentially (when
// `a.sequential` is true) or concurrently under a bounded semaphore
// (otherwise) so a flood of updates cannot spawn unbounded goroutines.
func (a *Account) dispatch(ctx context.Context, name string, raw map[string]any) {
	a.listenersMu.RLock()
	entries := make([]listenerEntry, 0, len(a.listeners))
	for _, e := range a.listeners {
		entries = append(entries, e)
	}
	a.listenersMu.RUnlock()

	// In sequential mode we just walk the list and stop on dispatchStop.
	if a.sequential {
		ev := events.New(raw)
		for _, e := range entries {
			if e.matcher.Kind() != capitalize(name) {
				continue
			}
			if !e.matcher.Match(ctx, raw) {
				continue
			}
			if err := e.fn(ctx, ev); err != nil {
				if errors.Is(err, dispatchStop) {
					return
				}
				a.infof("listener %s: %v", e.matcher.Kind(), err)
			}
		}
		return
	}

	// Concurrent mode: bound the fan-out with a semaphore so a flood of
	// updates cannot exhaust goroutines. Each matching listener runs
	// in its own goroutine (so a slow listener cannot block the WS
	// read loop or other listeners), but at most `dispatchConcurrency`
	// run at once.
	//
	// IMPORTANT: we deliberately do NOT wg.Wait() here. The WS read
	// loop must return to ws.Read() immediately so it can pick up the
	// next frame; waiting for listeners (which often make their own
	// network calls — SendMessage, GetChats, ...) would serialise the
	// whole pipeline and make the bot appear unresponsive under load.
	// Each listener goroutine gets its own shallow copy of `raw` so
	// that a listener calling ev.Set (which writes to the map under
	// the Event's mutex) cannot race with another listener reading
	// from the same map.
	sem := make(chan struct{}, dispatchConcurrency)
	for _, e := range entries {
		if e.matcher.Kind() != capitalize(name) {
			continue
		}
		if !e.matcher.Match(ctx, raw) {
			continue
		}
		// Shallow-copy raw so each listener has its own map. A
		// shallow copy is sufficient because the top-level values
		// (nested maps/slices) are themselves treated as read-only
		// by the helpers; only the top-level keys are mutated by
		// ev.Set.
		rawCopy := make(map[string]any, len(raw))
		for k, v := range raw {
			rawCopy[k] = v
		}
		ev := events.New(rawCopy)
		sem <- struct{}{}
		go func(fn Listener, m Matcher, ev *events.Event) {
			defer func() { <-sem }()
			defer func() {
				if r := recover(); r != nil {
					a.errorf("listener %s panic: %v", m.Kind(), r)
				}
			}()
			if err := fn(ctx, ev); err != nil {
				if errors.Is(err, dispatchStop) {
					return
				}
				a.infof("listener %s: %v", m.Kind(), err)
			}
		}(e.fn, e.matcher, ev)
	}
}

// capitalize converts snake_case to CamelCase (matches the helper used
// in the WS dispatch path: "message_updates" -> "MessageUpdates").
func capitalize(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

// AddListener registers a callback under the supplied matcher. It
// returns the integer ID of the listener; pass it to RemoveListener to
// deregister.
func (a *Account) AddListener(fn Listener, m Matcher) int {
	if fn == nil || m == nil {
		return -1
	}
	a.listenersMu.Lock()
	defer a.listenersMu.Unlock()
	a.nextListenerID++
	a.listeners[a.nextListenerID] = listenerEntry{fn: fn, matcher: m}
	return a.nextListenerID
}

// RemoveListener removes a previously-registered listener by its ID.
func (a *Account) RemoveListener(id int) bool {
	a.listenersMu.Lock()
	defer a.listenersMu.Unlock()
	if _, ok := a.listeners[id]; ok {
		delete(a.listeners, id)
		return true
	}
	return false
}

// dispatchStop is the sentinel returned by a listener to halt further
// dispatching of an update.
var dispatchStop = errors.New("account: stop dispatch")
