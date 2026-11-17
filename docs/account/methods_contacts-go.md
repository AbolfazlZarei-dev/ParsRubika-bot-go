# `account` — Methods contacts

_Source: [`methods_contacts.go`](../../ParsRubika/account/methods_contacts.go)_

## Symbols

### `GetContacts`

```go
func (a *Account) GetContacts(ctx context.Context, startID string) (*events.Event, error)
```

GetContacts returns the user's address book (paged).

### `GetContactsUpdates`

```go
func (a *Account) GetContactsUpdates(ctx context.Context, state ...string) (*events.Event, error)
```

GetContactsUpdates returns the contacts updates after `state`
(omit it to start from the newest state).

### `AddAddressBook`

```go
func (a *Account) AddAddressBook(ctx context.Context, phone, firstName, lastName string) (*events.Event, error)
```

AddAddressBook adds a single phone number to the address book.

### `DeleteContact`

```go
func (a *Account) DeleteContact(ctx context.Context, userGUID string) (*events.Event, error)
```

DeleteContact deletes a contact from the address book.

### `ResetContacts`

```go
func (a *Account) ResetContacts(ctx context.Context) (*events.Event, error)
```

ResetContacts resets the address book.

### `GetContactsLastOnline`

```go
func (a *Account) GetContactsLastOnline(ctx context.Context, userGUIDs []string) (*events.Event, error)
```

GetContactsLastOnline returns the last-online time for a list of
user GUIDs.

### `GetTopUsers`

```go
func (a *Account) GetTopUsers(ctx context.Context) (*events.Event, error)
```

GetTopUsers returns the user's top chat partners.

### `RemoveFromTopUsers`

```go
func (a *Account) RemoveFromTopUsers(ctx context.Context, userGUID string) (*events.Event, error)
```

RemoveFromTopUsers removes a user from the top-chats list.
