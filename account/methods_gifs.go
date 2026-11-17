package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// GIF sets.

// GetMyGifSet returns the user's saved gif set.
func (a *Account) GetMyGifSet(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getMyGifSet", Input: map[string]any{}, Encrypt: true})
}

// AddToMyGifSet adds a gif to the user's gif set.
func (a *Account) AddToMyGifSet(ctx context.Context, objectGUID, messageID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "addToMyGifSet",
		Input:   map[string]any{"object_guid": objectGUID, "message_id": messageID},
		Encrypt: true,
	})
}

// RemoveFromMyGifSet removes a gif from the user's gif set. The API
// uses `file_id`, not `object_guid + message_id`.
func (a *Account) RemoveFromMyGifSet(ctx context.Context, fileID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "removeFromMyGifSet",
		Input:   map[string]any{"file_id": fileID},
		Encrypt: true,
	})
}

// DeleteAllGifs removes all gifs from the user's gif set. The file_id
// may come back as either a string or a JSON number depending on the
// server, so we use the event's AnyStr accessor (via a temporary Event
// wrapper) to coerce it.
func (a *Account) DeleteAllGifs(ctx context.Context) error {
	res, err := a.GetMyGifSet(ctx)
	if err != nil {
		return err
	}
	for _, item := range res.Slice("gifs") {
		if m, ok := item.(map[string]any); ok {
			// Wrap in a temporary Event so we can use AnyStr for
			// type-coercion (file_id may be string or number).
			ev := events.New(m)
			fileID := ev.AnyStr("file_id")
			if fileID != "" {
				if _, err := a.RemoveFromMyGifSet(ctx, fileID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
