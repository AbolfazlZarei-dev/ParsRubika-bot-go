package rubiai

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Feedback is the user-sent opinion on an AI message.
type Feedback string

const (
	FeedbackLike    Feedback = "Like"
	FeedbackDislike Feedback = "Dislike"
	FeedbackNone    Feedback = "None"
)

// normalizeFeedback maps loose user input (case-insensitive, plus the
// Good/Bad aliases used by stop_session_stream) to the canonical value.
// Returns "" when the input is not recognised.
func normalizeFeedback(in string, allowNone bool) Feedback {
	switch strings.ToLower(strings.TrimSpace(in)) {
	case "like", "good":
		return FeedbackLike
	case "dislike", "bad":
		return FeedbackDislike
	case "none":
		if allowNone {
			return FeedbackNone
		}
	}
	return ""
}

// WelcomeLocation is the application area that getWelcome describes.
type WelcomeLocation string

const (
	WelcomeVOD       WelcomeLocation = "VOD"
	WelcomeHome      WelcomeLocation = "Home"
	WelcomeMessenger WelcomeLocation = "Messenger"
	WelcomeSocial    WelcomeLocation = "Social"
)

// normalizeLocation maps loose input to one of the four accepted values.
func normalizeLocation(in string) WelcomeLocation {
	switch strings.ToLower(strings.TrimSpace(in)) {
	case "vod":
		return WelcomeVOD
	case "home":
		return WelcomeHome
	case "messenger":
		return WelcomeMessenger
	case "social":
		return WelcomeSocial
	}
	return ""
}

// GetWelcome retrieves welcome information for an application location.
// `location` accepts VOD/Home/Messenger/Social (case-insensitive).
func (r *RubiaAI) GetWelcome(ctx context.Context, location string, locationID string) (*events.Event, error) {
	loc := normalizeLocation(location)
	if loc == "" {
		return nil, coreerrors.New(coreerrors.KindValidation,
			"rubiai: location must be one of VOD/Home/Messenger/Social")
	}
	return r.request(ctx, "getWelcome", map[string]any{
		"location":    string(loc),
		"location_id": locationID,
	})
}

// GetServiceAccess retrieves the current AI service access information.
func (r *RubiaAI) GetServiceAccess(ctx context.Context) (*events.Event, error) {
	return r.request(ctx, "getServiceAccess", nil)
}

// GetFeedbackIssues retrieves the available feedback issue types.
func (r *RubiaAI) GetFeedbackIssues(ctx context.Context) (*events.Event, error) {
	return r.request(ctx, "getFeedbackIssues", nil)
}

// GetSessions retrieves AI chat sessions. When allPages is true (the
// default) the method auto-paginates following has_continue/next_start_id
// and returns a single event wrapping {"sessions":[...]}. Otherwise it
// returns the raw page for the supplied start_id/limit.
func (r *RubiaAI) GetSessions(ctx context.Context, allPages bool, limit int, startID string) (*events.Event, error) {
	if limit <= 0 {
		limit = 10
	}
	if !allPages {
		return r.request(ctx, "getSessions", map[string]any{
			"start_id": startID,
			"limit":    limit,
		})
	}
	all := []any{}
	cursor := startID
	for {
		ev, err := r.request(ctx, "getSessions", map[string]any{
			"start_id": cursor,
			"limit":    limit,
		})
		if err != nil {
			return nil, err
		}
		if sess := findSlice(ev, "sessions"); sess != nil {
			all = append(all, sess...)
		}
		if !ev.Bool("has_continue") {
			break
		}
		next := ev.AnyStr("next_start_id")
		if next == "" || next == cursor {
			break
		}
		cursor = next
	}
	return events.New(map[string]any{"sessions": all}), nil
}

// GetChatHistory retrieves messages from an AI session. When allPages is
// true (the default) the method auto-paginates and returns a single event
// wrapping {"messages":[...]}. Otherwise it returns the raw page.
func (r *RubiaAI) GetChatHistory(ctx context.Context, sessionID string, allPages bool, limit int, startID string) (*events.Event, error) {
	if sessionID == "" {
		return nil, coreerrors.New(coreerrors.KindValidation, "rubiai: session_id is required")
	}
	if limit <= 0 {
		limit = 100
	}
	if !allPages {
		return r.request(ctx, "getChatHistory", map[string]any{
			"session_id": sessionID,
			"limit":      limit,
			"start_id":   startID,
		})
	}
	all := []any{}
	cursor := startID
	for {
		ev, err := r.request(ctx, "getChatHistory", map[string]any{
			"session_id": sessionID,
			"limit":      limit,
			"start_id":   cursor,
		})
		if err != nil {
			return nil, err
		}
		if msgs := findSlice(ev, "messages"); msgs != nil {
			all = append(all, msgs...)
		}
		if !ev.Bool("has_continue") {
			break
		}
		next := ev.AnyStr("next_start_id")
		if next == "" || next == cursor {
			break
		}
		cursor = next
	}
	return events.New(map[string]any{"messages": all}), nil
}

// DeleteSession deletes an AI chat session.
func (r *RubiaAI) DeleteSession(ctx context.Context, sessionID string) (*events.Event, error) {
	if sessionID == "" {
		return nil, coreerrors.New(coreerrors.KindValidation, "rubiai: session_id is required")
	}
	return r.request(ctx, "deleteSession", map[string]any{"session_id": sessionID})
}

// EditSession updates the pin state and/or title of a session. Pass the
// zero value for a field you do not want to change (title="" / isPinned
// left as its sentinel via the pointers).
func (r *RubiaAI) EditSession(ctx context.Context, sessionID string, title string, isPinned *bool) (*events.Event, error) {
	if sessionID == "" {
		return nil, coreerrors.New(coreerrors.KindValidation, "rubiai: session_id is required")
	}
	input := map[string]any{
		"session_id":        sessionID,
		"update_parameters": []string{},
	}
	if isPinned != nil {
		input["is_pinned"] = *isPinned
		input["update_parameters"] = append(input["update_parameters"].([]string), "is_pinned")
	}
	if title != "" {
		input["title"] = title
		input["update_parameters"] = append(input["update_parameters"].([]string), "title")
	}
	return r.request(ctx, "editSession", input)
}

// SendFeedback submits feedback (Like/Dislike/None) for an AI message.
func (r *RubiaAI) SendFeedback(ctx context.Context, sessionID, messageID, feedback string) (*events.Event, error) {
	fb := normalizeFeedback(feedback, true)
	if fb == "" {
		return nil, coreerrors.New(coreerrors.KindValidation,
			fmt.Sprintf("rubiai: invalid feedback %q (Like/Dislike/None)", feedback))
	}
	return r.request(ctx, "sendFeedback", map[string]any{
		"session_id": sessionID,
		"message_id": messageID,
		"feedback":   string(fb),
	})
}

// SetFeedbackDetail attaches a detailed issue type (and optional
// description) to previously-submitted feedback. The issue types come
// from GetFeedbackIssues.
func (r *RubiaAI) SetFeedbackDetail(ctx context.Context, sessionID, messageID, issueType, issueDesc string) (*events.Event, error) {
	return r.request(ctx, "setFeedbackDetail", map[string]any{
		"session_id":          sessionID,
		"message_id":          messageID,
		"feedback_issue_type": issueType,
		"feedback_issue_desc": issueDesc,
	})
}

// StopSessionStream stops an active AI session stream and submits the
// final feedback (Like/Dislike; Good/Bad are accepted as aliases).
func (r *RubiaAI) StopSessionStream(ctx context.Context, sessionID, messageID, feedback string) (*events.Event, error) {
	fb := normalizeFeedback(feedback, false)
	if fb == "" {
		return nil, coreerrors.New(coreerrors.KindValidation,
			fmt.Sprintf("rubiai: invalid feedback %q (Like/Dislike)", feedback))
	}
	return r.request(ctx, "stopSessionStream", map[string]any{
		"session_id": sessionID,
		"message_id": messageID,
		"feedback":   string(fb),
	})
}

// SpeechToText converts raw voice bytes to text via the AI speech-to-text
// service. The bytes are base64-encoded before being sent.
func (r *RubiaAI) SpeechToText(ctx context.Context, voiceData []byte) (*events.Event, error) {
	if len(voiceData) == 0 {
		return nil, coreerrors.New(coreerrors.KindValidation, "rubiai: voice_data is empty")
	}
	return r.request(ctx, "speechToText", map[string]any{
		"voice_data": base64.StdEncoding.EncodeToString(voiceData),
	})
}

// SpeechToTextFromFile reads the voice file at `path` and passes its
// contents to SpeechToText.
func (r *RubiaAI) SpeechToTextFromFile(ctx context.Context, path string) (*events.Event, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, coreerrors.Wrap(coreerrors.KindValidation, "rubiai: read voice file", err)
	}
	return r.SpeechToText(ctx, data)
}

// DownloadImage retrieves an AI-generated image, polling until it is
// ready. When `saveAs` is non-empty the decoded image bytes are written
// there. Returns the final API response event (status Success/Error).
//
// The polling loop honours ctx cancellation; maxWait and interval
// default to 60 s / 2 s when zero.
func (r *RubiaAI) DownloadImage(ctx context.Context, taskID, saveAs string, maxWait, interval time.Duration) (*events.Event, error) {
	if taskID == "" {
		return nil, coreerrors.New(coreerrors.KindValidation, "rubiai: task_id is required")
	}
	if maxWait <= 0 {
		maxWait = 60 * time.Second
	}
	if interval <= 0 {
		interval = 2 * time.Second
	}
	deadline := time.Now().Add(maxWait)
	var last *events.Event
	for {
		ev, err := r.request(ctx, "downloadImage", map[string]any{"task_id": taskID})
		if err != nil {
			return nil, err
		}
		last = ev
		switch ev.AnyStr("status") {
		case "Success":
			if saveAs != "" {
				imgB64 := ev.AnyStr("image")
				if imgB64 == "" {
					return nil, coreerrors.New(coreerrors.KindAPI, "rubiai: image data missing in response")
				}
				img, err := base64.StdEncoding.DecodeString(imgB64)
				if err != nil {
					return nil, coreerrors.Wrap(coreerrors.KindAPI, "rubiai: decode image base64", err)
				}
				if err := os.WriteFile(saveAs, img, 0o600); err != nil {
					return nil, coreerrors.Wrap(coreerrors.KindAPI, "rubiai: write image file", err)
				}
			}
			return ev, nil
		case "Error":
			return ev, nil
		}
		if time.Now().After(deadline) {
			return last, nil
		}
		select {
		case <-ctx.Done():
			return last, ctx.Err()
		case <-time.After(interval):
		}
	}
}

// findSlice returns the slice found under `key` (recursively, mirroring
// the Python Data.find_keys helper) or nil.
func findSlice(ev *events.Event, key string) []any {
	if ev == nil {
		return nil
	}
	if s := ev.Slice(key); s != nil {
		return s
	}
	// events.Event.Slice only looks at top-level; fall back to a recursive
	// search over the raw map for nested shapes the server sometimes uses.
	return findSliceRecursive(ev.Map(), key)
}

func findSliceRecursive(m map[string]any, key string) []any {
	if m == nil {
		return nil
	}
	if v, ok := m[key]; ok {
		if s, ok := v.([]any); ok {
			return s
		}
	}
	for _, v := range m {
		switch x := v.(type) {
		case map[string]any:
			if r := findSliceRecursive(x, key); r != nil {
				return r
			}
		case []any:
			for _, item := range x {
				if mm, ok := item.(map[string]any); ok {
					if r := findSliceRecursive(mm, key); r != nil {
						return r
					}
				}
			}
		}
	}
	return nil
}
