package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Object ownership transfer.

// RequestChangeOwner requests ownership transfer to a new user.
func (a *Account) RequestChangeOwner(ctx context.Context, objectGUID, newOwnerGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "requestChangeObjectOwner", Input: map[string]any{
		"object_guid":         objectGUID,
		"new_owner_user_guid": newOwnerGUID,
	}, Encrypt: true})
}

// AcceptOwnershipRequest accepts a pending ownership transfer.
func (a *Account) AcceptOwnershipRequest(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "replyRequestObjectOwner", Input: map[string]any{
		"object_guid": objectGUID,
		"action":      "Accept",
	}, Encrypt: true})
}

// RejectOwnershipRequest rejects a pending ownership transfer.
func (a *Account) RejectOwnershipRequest(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "replyRequestObjectOwner", Input: map[string]any{
		"object_guid": objectGUID,
		"action":      "Reject",
	}, Encrypt: true})
}

// CancelChangeOwner cancels a pending ownership transfer.
func (a *Account) CancelChangeOwner(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "cancelChangeObjectOwner", Input: map[string]any{
		"object_guid": objectGUID,
	}, Encrypt: true})
}

// GetPendingOwner returns the pending owner of a chat.
func (a *Account) GetPendingOwner(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getPendingObjectOwner", Input: map[string]any{
		"object_guid": objectGUID,
	}, Encrypt: true})
}
