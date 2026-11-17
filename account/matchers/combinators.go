package matchers

import "context"

// And matches when ALL sub-matchers match (short-circuit).
type And struct {
	matchers []Matcher
	kind     string
}

// NewAnd constructs an And matcher.
func NewAnd(kind string, ms ...Matcher) *And {
	if kind == "" {
		kind = firstKind(ms)
	}
	return &And{matchers: ms, kind: kind}
}

// Kind implements Matcher.
func (a *And) Kind() string { return a.kind }

// Match implements Matcher.
func (a *And) Match(ctx context.Context, raw map[string]any) bool {
	for _, m := range a.matchers {
		if !m.Match(ctx, raw) {
			return false
		}
	}
	return true
}

// Or matches when ANY sub-matcher matches.
type Or struct {
	matchers []Matcher
	kind     string
}

// NewOr constructs an Or matcher.
func NewOr(kind string, ms ...Matcher) *Or {
	if kind == "" {
		kind = firstKind(ms)
	}
	return &Or{matchers: ms, kind: kind}
}

// Kind implements Matcher.
func (o *Or) Kind() string { return o.kind }

// Match implements Matcher.
func (o *Or) Match(ctx context.Context, raw map[string]any) bool {
	for _, m := range o.matchers {
		if m.Match(ctx, raw) {
			return true
		}
	}
	return false
}

// Not negates the wrapped matcher.
type Not struct {
	inner Matcher
	kind  string
}

// NewNot constructs a Not matcher.
func NewNot(kind string, m Matcher) *Not {
	if kind == "" {
		kind = m.Kind()
	}
	return &Not{inner: m, kind: kind}
}

// Kind implements Matcher.
func (n *Not) Kind() string { return n.kind }

// Match implements Matcher.
func (n *Not) Match(ctx context.Context, raw map[string]any) bool {
	return !n.inner.Match(ctx, raw)
}

func firstKind(ms []Matcher) string {
	for _, m := range ms {
		if m != nil {
			return m.Kind()
		}
	}
	return "MessageUpdates"
}
