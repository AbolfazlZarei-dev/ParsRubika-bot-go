package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Contacts, address book and top-chat lists.

// GetContacts returns the user's address book (paged).
func (a *Account) GetContacts(ctx context.Context, startID string) (*events.Event, error) {
	in := map[string]any{}
	if startID != "" {
		in["start_id"] = startID
	}
	return a.Request(ctx, RequestOptions{Method: "getContacts", Input: in, Encrypt: true})
}

// GetContactsUpdates returns the contacts updates after `state`
// (omit it to start from the newest state).
func (a *Account) GetContactsUpdates(ctx context.Context, state ...string) (*events.Event, error) {
	in := map[string]any{}
	if len(state) > 0 && state[0] != "" {
		in["state"] = state[0]
	}
	return a.Request(ctx, RequestOptions{Method: "getContactsUpdates", Input: in, Encrypt: true})
}

// AddAddressBook adds a single phone number to the address book.
func (a *Account) AddAddressBook(ctx context.Context, phone, firstName, lastName string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method: "addAddressBook",
		Input: map[string]any{
			"phone":      phone,
			"first_name": firstName,
			"last_name":  lastName,
		},
		Encrypt: true,
	})
}

// DeleteContact deletes a contact from the address book.
func (a *Account) DeleteContact(ctx context.Context, userGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "deleteContact", Input: map[string]any{"user_guid": userGUID}, Encrypt: true})
}

// ResetContacts resets the address book.
func (a *Account) ResetContacts(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "resetContacts", Input: map[string]any{}, Encrypt: true})
}

// GetContactsLastOnline returns the last-online time for a list of
// user GUIDs.
func (a *Account) GetContactsLastOnline(ctx context.Context, userGUIDs []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "getContactsLastOnline",
		Input:   map[string]any{"user_guids": userGUIDs},
		Encrypt: true,
	})
}

// GetTopUsers returns the user's top chat partners.
func (a *Account) GetTopUsers(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getTopChatUsers", Input: map[string]any{}, Encrypt: true})
}

// RemoveFromTopUsers removes a user from the top-chats list.
func (a *Account) RemoveFromTopUsers(ctx context.Context, userGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "removeFromTopChatUsers", Input: map[string]any{"user_guid": userGUID}, Encrypt: true})
}
