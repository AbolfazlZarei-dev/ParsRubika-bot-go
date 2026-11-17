package account

import (
	"context"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account/listeners"
)

// AddListener registers a callback under the supplied matcher. It is
// defined here (in the account package) so callers can use either
// the listeners sub-package types or their own matcher.
//
// The function signature accepts a `listeners.MatcherChain` which both
// `listeners.Messages`, `listeners.Chats`, ... and every matcher from
// the `matchers` sub-package implement.
func AddListener(a *Account, fn Listener, l listeners.MatcherChain) int {
	return a.AddListener(fn, listenerChain{l})
}

// listenerChain adapts the listeners.MatcherChain interface to the
// account.Matcher interface.
type listenerChain struct{ inner listeners.MatcherChain }

func (lc listenerChain) Kind() string { return lc.inner.Kind() }
func (lc listenerChain) Match(ctx context.Context, raw map[string]any) bool {
	return lc.inner.Match(ctx, raw)
}
