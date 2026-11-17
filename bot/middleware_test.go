package bot

import "testing"

// TestChatTypeOf_PrefixSemantics locks in the chat-type classification
// used by ChatTypeMiddleware. A user (private) chat is identified by
// the "u0" prefix; "b0" is a bot GUID and must NOT classify as
// "private". This previously regressed (b0 -> "private"), causing
// ChatTypeMiddleware("private") to admit bot GUIDs and reject real
// user chats.
func TestChatTypeOf_PrefixSemantics(t *testing.T) {
	cases := []struct {
		name   string
		chatID string
		want   string
	}{
		{"user chat is private", "u0abcdef", "private"},
		{"bot GUID is NOT private", "b0abcdef", ""},
		{"group", "g0abcdef", "group"},
		{"channel", "c0abcdef", "channel"},
		{"service is unclassified", "s0abcdef", ""},
		{"empty", "", ""},
		{"nil update", "", ""},
		{"unknown prefix", "x0abcdef", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			upd := &Update{ChatID: c.chatID}
			if got := chatTypeOf(upd); got != c.want {
				t.Errorf("chatTypeOf(chatID=%q) = %q, want %q", c.chatID, got, c.want)
			}
		})
	}
}
