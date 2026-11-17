package enums

// ChatType enumerates the kinds of chat the server can return.
type ChatType string

const (
	ChatTypeUser    ChatType = "User"
	ChatTypeBot     ChatType = "Bot"
	ChatTypeGroup   ChatType = "Group"
	ChatTypeChannel ChatType = "Channel"
	ChatTypeService ChatType = "Service"
)

// String implements fmt.Stringer.
func (c ChatType) String() string { return string(c) }

// FromGUID infers the chat type from the prefix of the supplied GUID:
//
// - "c0..." → Channel
// - "g0..." → Group
// - anything else (including "u0...", "b0...") → User
func ChatTypeFromGUID(guid string) ChatType {
	switch {
	case len(guid) >= 2 && guid[:2] == "c0":
		return ChatTypeChannel
	case len(guid) >= 2 && guid[:2] == "g0":
		return ChatTypeGroup
	default:
		return ChatTypeUser
	}
}
