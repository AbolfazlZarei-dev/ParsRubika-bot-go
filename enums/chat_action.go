package enums

// ChatAction is the activity type reported during typing / uploading /
// recording.
type ChatAction string

const (
	ChatActionTyping    ChatAction = "Typing"
	ChatActionUploading ChatAction = "Uploading"
	ChatActionRecording ChatAction = "Recording"
)

// String implements fmt.Stringer.
func (c ChatAction) String() string { return string(c) }
