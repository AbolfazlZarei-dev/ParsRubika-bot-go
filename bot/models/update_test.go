package models

import "testing"

// TestIsPrivateChat_PrefixSemantics locks in the Rubika GUID-prefix
// semantics for private-chat detection. A one-to-one (user-to-bot)
// chat is identified by the "u0" prefix; "b0" is a bot GUID and must
// NOT be treated as private. This previously regressed (b0 was used),
// which made IsPrivateChat and everything delegating to it
// (Ctx.IsPrivate, rules.PrivateChat) misclassify bot GUIDs as private
// chats and never match real user chats.
func TestIsPrivateChat_PrefixSemantics(t *testing.T) {
	cases := []struct {
		name   string
		chatID string
		want   bool
	}{
		{"user chat is private", "u0" + repeat("a", 32), true},
		{"bot GUID is NOT private", "b0" + repeat("a", 32), false},
		{"group is not private", "g0" + repeat("a", 32), false},
		{"channel is not private", "c0" + repeat("a", 32), false},
		{"service is not private", "s0" + repeat("a", 32), false},
		{"empty is not private", "", false},
		{"single char is not private", "u", false},
		{"unknown prefix is not private", "x0abcdef", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u := &Update{ChatID: c.chatID}
			if got := u.IsPrivateChat(); got != c.want {
				t.Errorf("IsPrivateChat(chatID=%q) = %v, want %v", c.chatID, got, c.want)
			}
		})
	}
}

func repeat(s string, n int) string {
	out := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		out = append(out, s...)
	}
	return string(out)
}
