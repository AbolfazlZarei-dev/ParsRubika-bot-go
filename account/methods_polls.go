package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Poll creation and voting.

// CreatePoll creates and sends a poll message.
func (a *Account) CreatePoll(ctx context.Context, objectGUID, question string, options []string, pollType string, isAnonymous, allowMulti bool, correctIdx *int, explanation, replyToMessageID string) (*events.Event, error) {
	if pollType == "" {
		pollType = "Regular"
	}
	in := map[string]any{
		"object_guid":             objectGUID,
		"question":                question,
		"options":                 options,
		"type":                    pollType,
		"allows_multiple_answers": allowMulti,
		"is_anonymous":            isAnonymous,
		"rnd":                     randInt63(),
	}
	if correctIdx != nil {
		in["correct_option_index"] = *correctIdx
	}
	if explanation != "" {
		in["explanation"] = explanation
	}
	if replyToMessageID != "" {
		in["reply_to_message_id"] = replyToMessageID
	}
	return a.Request(ctx, RequestOptions{Method: "createPoll", Input: in, Encrypt: true})
}

// VotePoll votes on a single poll option. The API expects a single
// `selection_index: int`, not a list.
func (a *Account) VotePoll(ctx context.Context, pollID string, selectionIndex int) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "votePoll",
		Input:   map[string]any{"poll_id": pollID, "selection_index": selectionIndex},
		Encrypt: true,
	})
}

// GetPollStatus fetches the status of a poll. The API only needs
// `{poll_id}` — no `object_guid`.
func (a *Account) GetPollStatus(ctx context.Context, pollID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getPollStatus", Input: map[string]any{"poll_id": pollID}, Encrypt: true})
}

// GetPollOptionVoters returns the voters of a poll option, paginated
// by `startID` (the API uses `start_id`, not `cursor`).
func (a *Account) GetPollOptionVoters(ctx context.Context, pollID string, optionIndex int, startID string) (*events.Event, error) {
	in := map[string]any{"poll_id": pollID, "selection_index": optionIndex}
	if startID != "" {
		in["start_id"] = startID
	}
	return a.Request(ctx, RequestOptions{Method: "getPollOptionVoters", Input: in, Encrypt: true})
}

// SetPollAction stops or retracts a poll.
func (a *Account) SetPollAction(ctx context.Context, pollID, action string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setPollAction", Input: map[string]any{
		"poll_id": pollID,
		"action":  action, // "Stop" or "Retract"
	}, Encrypt: true})
}

// StopPoll stops a poll (action = "Stop").
func (a *Account) StopPoll(ctx context.Context, pollID string) (*events.Event, error) {
	return a.SetPollAction(ctx, pollID, "Stop")
}

// RetractPoll retracts a poll (action = "Retract").
func (a *Account) RetractPoll(ctx context.Context, pollID string) (*events.Event, error) {
	return a.SetPollAction(ctx, pollID, "Retract")
}
