// Package voicekit provides high-level helpers for the Rubika
// voice-chat API. It wraps the low-level voice-chat methods on the
// account client with a more ergonomic API.
package voicekit

import (
	"context"
	"fmt"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// State is the state of a voice-chat participant.
type State string

const (
	StateMute   State = "Mute"
	StateUnmute State = "Unmute"
)

// Activity is the activity reported in a voice chat.
type Activity string

const (
	ActivitySpeaking Activity = "Speaking"
)

// Join joins a voice chat.
func Join(ctx context.Context, acc *account.Account, objectGUID, voiceChatID, sdpOffer string) (*events.Event, error) {
	return acc.JoinVoiceChat(ctx, objectGUID, voiceChatID, sdpOffer)
}

// Leave leaves a voice chat.
func Leave(ctx context.Context, acc *account.Account, objectGUID, voiceChatID string) (*events.Event, error) {
	return acc.DiscardVoiceChat(ctx, objectGUID, voiceChatID)
}

// MuteParticipant mutes a participant in a voice chat.
func MuteParticipant(ctx context.Context, acc *account.Account, objectGUID, voiceChatID, participantGUID string) (*events.Event, error) {
	return acc.SetVoiceChatState(ctx, objectGUID, voiceChatID, string(StateMute), participantGUID)
}

// UnmuteParticipant unmutes a participant in a voice chat.
func UnmuteParticipant(ctx context.Context, acc *account.Account, objectGUID, voiceChatID, participantGUID string) (*events.Event, error) {
	return acc.SetVoiceChatState(ctx, objectGUID, voiceChatID, string(StateUnmute), participantGUID)
}

// ReportSpeaking sends a speaking activity notification.
func ReportSpeaking(ctx context.Context, acc *account.Account, objectGUID, voiceChatID, participantGUID string) (*events.Event, error) {
	return acc.SendGroupVoiceChatActivity(ctx, objectGUID, voiceChatID, string(ActivitySpeaking), participantGUID)
}

// Create creates a voice chat in a group.
func Create(ctx context.Context, acc *account.Account, groupGUID, title string) (*events.Event, error) {
	return acc.CreateGroupVoiceChat(ctx, groupGUID, title)
}

// Close closes a voice chat.
func Close(ctx context.Context, acc *account.Account, objectGUID, voiceChatID string) (*events.Event, error) {
	return acc.DiscardVoiceChat(ctx, objectGUID, voiceChatID)
}

// GetVoiceChatID extracts the voice_chat_id from a chat-info
// event. This is a convenience helper since the ID is nested
// deeply in the response.
func GetVoiceChatID(chatInfo *events.Event) string {
	if chatInfo == nil {
		return ""
	}
	return chatInfo.Str("group_voice_chat_id", "voice_chat_id")
}

// FormatError returns a human-readable description of a voice-chat
// error.
func FormatError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("voice chat error: %v", err)
}
