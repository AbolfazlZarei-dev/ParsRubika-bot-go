package account

import (
	"context"
	"fmt"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/enums"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"strings"
)

// Group and channel voice chats.

// CreateChannelVoiceChat creates a voice chat inside a channel.
func (a *Account) CreateChannelVoiceChat(ctx context.Context, channelGUID, title string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "createChannelVoiceChat",
		Input:   map[string]any{"channel_guid": channelGUID, "title": title},
		Encrypt: true,
	})
}

// DiscardChannelVoiceChat closes a channel voice chat.
func (a *Account) DiscardChannelVoiceChat(ctx context.Context, channelGUID, voiceChatID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "discardChannelVoiceChat",
		Input:   map[string]any{"channel_guid": channelGUID, "voice_chat_id": voiceChatID},
		Encrypt: true,
	})
}

// SetChannelVoiceChatSetting updates voice-chat settings. Only the
// non-empty fields are sent; `updated_parameters` is built automatically.
func (a *Account) SetChannelVoiceChatSetting(ctx context.Context, channelGUID, voiceChatID string, settings map[string]any) (*events.Event, error) {
	in := map[string]any{
		"channel_guid":  channelGUID,
		"voice_chat_id": voiceChatID,
	}
	updated := []string{}
	for k, v := range settings {
		in[k] = v
		updated = append(updated, k)
	}
	in["updated_parameters"] = updated
	return a.Request(ctx, RequestOptions{Method: "setChannelVoiceChatSetting", Input: in, Encrypt: true})
}

// GetGroupVoiceChat returns info about a group voice chat.
func (a *Account) GetGroupVoiceChat(ctx context.Context, groupGUID, voiceChatID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getGroupVoiceChat", Input: map[string]any{
		"group_guid":    groupGUID,
		"voice_chat_id": voiceChatID,
	}, Encrypt: true})
}

// GetGroupVoiceChatParticipants returns the participants of a voice chat.
func (a *Account) GetGroupVoiceChatParticipants(ctx context.Context, groupGUID, voiceChatID, startID string) (*events.Event, error) {
	in := map[string]any{
		"group_guid":    groupGUID,
		"voice_chat_id": voiceChatID,
	}
	if startID != "" {
		in["start_id"] = startID
	}
	return a.Request(ctx, RequestOptions{Method: "getGroupVoiceChatParticipants", Input: in, Encrypt: true})
}

// GetChatVoiceChat returns the voice chat ID for a chat. Only groups
// and channels support voice chats; user/bot/service GUIDs return an
// explicit error rather than being mis-dispatched to GetGroupInfo.
func (a *Account) GetChatVoiceChat(ctx context.Context, objectGUID string) (*events.Event, error) {
	switch {
	case strings.HasPrefix(objectGUID, "c0"):
		return a.GetChannelInfo(ctx, objectGUID)
	case strings.HasPrefix(objectGUID, "g0"):
		return a.GetGroupInfo(ctx, objectGUID)
	default:
		return nil, fmt.Errorf("voice chats are only available for groups and channels, got GUID: %s", objectGUID)
	}
}

// CreateGroupVoiceChat creates a voice chat inside a group. The
// server expects the uniqueness token `rnd` and the protocol
// `version` alongside the GUID; `title` is optional.
func (a *Account) CreateGroupVoiceChat(ctx context.Context, groupGUID, title string) (*events.Event, error) {
	in := map[string]any{
		"group_guid": groupGUID,
		"rnd":        randInt63(),
		"version":    1,
	}
	if title != "" {
		in["title"] = title
	}
	return a.Request(ctx, RequestOptions{
		Method:  "createGroupVoiceChat",
		Input:   in,
		Encrypt: true,
	})
}

// LeaveGroupVoiceChat leaves a group voice chat.
func (a *Account) LeaveGroupVoiceChat(ctx context.Context, groupGUID, voiceChatID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "leaveGroupVoiceChat",
		Input:   map[string]any{"group_guid": groupGUID, "voice_chat_id": voiceChatID},
		Encrypt: true,
	})
}

// GetGroupVoiceChatUpdates fetches the latest voice-chat updates.
func (a *Account) GetGroupVoiceChatUpdates(ctx context.Context, groupGUID, voiceChatID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "getGroupVoiceChatUpdates",
		Input:   map[string]any{"group_guid": groupGUID, "voice_chat_id": voiceChatID},
		Encrypt: true,
	})
}

// SetGroupVoiceChatSetting updates voice-chat settings. Only the
// non-empty fields in `settings` are sent; the `updated_parameters`
// array is built automatically.
func (a *Account) SetGroupVoiceChatSetting(ctx context.Context, groupGUID, voiceChatID string, settings map[string]any) (*events.Event, error) {
	in := map[string]any{
		"group_guid":    groupGUID,
		"voice_chat_id": voiceChatID,
	}
	updated := []string{}
	for k, v := range settings {
		in[k] = v
		updated = append(updated, k)
	}
	in["updated_parameters"] = updated
	return a.Request(ctx, RequestOptions{Method: "setGroupVoiceChatSetting", Input: in, Encrypt: true})
}

// StartVoiceChat creates a voice chat in a group (alias for
// CreateGroupVoiceChat with a more intuitive name).
func (a *Account) StartVoiceChat(ctx context.Context, groupGUID, title string) (*events.Event, error) {
	return a.CreateGroupVoiceChat(ctx, groupGUID, title)
}

// RequestSpeakVoiceChat requests to speak in a voice chat.
func (a *Account) RequestSpeakVoiceChat(ctx context.Context, objectGUID, voiceChatID, participantGUID string) (*events.Event, error) {
	return a.SetVoiceChatState(ctx, objectGUID, voiceChatID, string(enums.VoiceChatStateRequestSpeak), participantGUID)
}

// CancelRequestSpeakVoiceChat cancels a pending speak request.
func (a *Account) CancelRequestSpeakVoiceChat(ctx context.Context, objectGUID, voiceChatID, participantGUID string) (*events.Event, error) {
	return a.SetVoiceChatState(ctx, objectGUID, voiceChatID, string(enums.VoiceChatStateCancelRequestSpeak), participantGUID)
}

// JoinVoiceChat joins a voice chat. The API method is dispatched to
// `joinGroupVoiceChat` or `joinChannelVoiceChat` based on the GUID
// prefix. `sdpOfferData` is the WebRTC SDP offer (may be empty for a
// simple presence-only join).
func (a *Account) JoinVoiceChat(ctx context.Context, objectGUID, voiceChatID, sdpOfferData string) (*events.Event, error) {
	method := "joinGroupVoiceChat"
	key := "group_guid"
	if strings.HasPrefix(objectGUID, "c0") {
		method = "joinChannelVoiceChat"
		key = "channel_guid"
	}
	in := map[string]any{
		key:                objectGUID,
		"voice_chat_id":    voiceChatID,
		"sdp_offer_data":   sdpOfferData,
		"self_object_guid": a.guid,
	}
	return a.Request(ctx, RequestOptions{Method: method, Input: in, Encrypt: true})
}

// DiscardVoiceChat closes a voice chat (dispatches to Group/Channel).
func (a *Account) DiscardVoiceChat(ctx context.Context, objectGUID, voiceChatID string) (*events.Event, error) {
	method := "discardGroupVoiceChat"
	key := "group_guid"
	if strings.HasPrefix(objectGUID, "c0") {
		method = "discardChannelVoiceChat"
		key = "channel_guid"
	}
	return a.Request(ctx, RequestOptions{
		Method:  method,
		Input:   map[string]any{key: objectGUID, "voice_chat_id": voiceChatID},
		Encrypt: true,
	})
}

// SetVoiceChatState updates the voice-chat state. `action` is "Mute"
// or "Unmute". The API method is dispatched to
// `setGroupVoiceChatState` / `setChannelVoiceChatState`.
func (a *Account) SetVoiceChatState(ctx context.Context, objectGUID, voiceChatID, action, participantGUID string) (*events.Event, error) {
	method := "setGroupVoiceChatState"
	key := "group_guid"
	if strings.HasPrefix(objectGUID, "c0") {
		method = "setChannelVoiceChatState"
		key = "channel_guid"
	}
	if participantGUID == "" {
		participantGUID = a.guid
	}
	return a.Request(ctx, RequestOptions{
		Method: method,
		Input: map[string]any{
			key:                       objectGUID,
			"voice_chat_id":           voiceChatID,
			"action":                  action,
			"participant_object_guid": participantGUID,
		},
		Encrypt: true,
	})
}

// SendGroupVoiceChatActivity sends an activity notification in a voice
// chat. `participantGUID` defaults to the current user's GUID.
func (a *Account) SendGroupVoiceChatActivity(ctx context.Context, objectGUID, voiceChatID, action, participantGUID string) (*events.Event, error) {
	if participantGUID == "" {
		participantGUID = a.guid
	}
	return a.Request(ctx, RequestOptions{
		Method: "sendGroupVoiceChatActivity",
		Input: map[string]any{
			"object_guid":             objectGUID,
			"voice_chat_id":           voiceChatID,
			"action":                  action,
			"participant_object_guid": participantGUID,
		},
		Encrypt: true,
	})
}
