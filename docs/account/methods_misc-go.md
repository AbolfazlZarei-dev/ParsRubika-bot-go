# `account` — Methods misc

_Source: [`methods_misc.go`](../../ParsRubika/account/methods_misc.go)_

## Symbols

### `GetAvailableReactions`

```go
func (a *Account) GetAvailableReactions(ctx context.Context) (*events.Event, error)
```

GetAvailableReactions returns the reactions available to the user.

### `GetWalletTransactionMessage`

```go
func (a *Account) GetWalletTransactionMessage(ctx context.Context, transferID, accessTransfer string) (*events.Event, error)
```

GetWalletTransactionMessage returns info about a wallet transaction.

### `SearchGlobalObjects`

```go
func (a *Account) SearchGlobalObjects(ctx context.Context, query string) (*events.Event, error)
```

SearchGlobalObjects performs a global object search.

### `GetProfileLinkItems`

```go
func (a *Account) GetProfileLinkItems(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetProfileLinkItems returns the items shown on a profile link. The
API uses `object_guid`, not `profile_link`.

### `GetTime`

```go
func (a *Account) GetTime(ctx context.Context) (*events.Event, error)
```

GetTime returns the server time.

### `GetRelatedObjects`

```go
func (a *Account) GetRelatedObjects(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetRelatedObjects returns the related objects of a GUID.

### `ReportObject`

```go
func (a *Account) ReportObject(ctx context.Context, objectGUID, reportType, reportDescription, reportTypeObject, messageID string) (*events.Event, error)
```

ReportObject reports an object for abuse. The API uses
`{object_guid, report_type, report_description, report_type_object,
message_id?}`.

### `ReportMessage`

```go
func (a *Account) ReportMessage(ctx context.Context, objectGUID, messageID, reportType, reportDescription string) (*events.Event, error)
```

ReportMessage is a convenience wrapper for reporting a single
message.

### `ReportChat`

```go
func (a *Account) ReportChat(ctx context.Context, objectGUID, reportType, reportDescription string) (*events.Event, error)
```

ReportChat is a convenience wrapper for reporting a chat.

### `GetChannelPostByLink`

```go
func (a *Account) GetChannelPostByLink(ctx context.Context, postURL string) (*events.Event, error)
```

GetChannelPostByLink resolves a channel post from a web URL. It calls
getLinkFromAppUrl to convert the web link into a deep link, then parses
the result to extract channel GUID and message ID.

### `GetPostByAskLink`

```go
func (a *Account) GetPostByAskLink(ctx context.Context, askLink string) (*events.Event, error)
```

GetPostByAskLink resolves a post from an ask link. It calls
getLinkFromAppUrl and then fetches the message.

### `ScrapeChannelPostLink`

```go
func (a *Account) ScrapeChannelPostLink(ctx context.Context, postURL string) (string, error)
```

ScrapeChannelPostLink fetches the rubika.ir web page for a post and
extracts the deep link. This is a fallback when getLinkFromAppUrl
does not resolve.

Security: only http(s) URLs that match the rubika.ir post-link pattern
are accepted; the dial is routed through netutil's SSRF-safe transport
so private/loopback/link-local IPs cannot be reached even via DNS
rebinding, and the response is capped at netutil.MaxFetchBytes.

### `FetchJSON`

```go
func (a *Account) FetchJSON(ctx context.Context, url string, out any) error
```

FetchJSON is a low-level helper that performs an HTTP GET and decodes
the response as JSON. It is useful for calling third-party APIs
(e.g. translation services, Rubino download helpers).

Security: the URL is validated through the netutil SSRF policy and the
dial goes through a SSRF-safe transport; the response is capped at
netutil.MaxFetchBytes.

### `GetContactsStories`

```go
func (a *Account) GetContactsStories(ctx context.Context) (*events.Event, error)
```

GetContactsStories returns the stories of the user's contacts.

(GetContactsLastOnline is already defined in methods_users.go.)

### `SendRubinoPost`

```go
func (a *Account) SendRubinoPost(ctx context.Context, objectGUID, postID, postProfileID string, isMute bool) (*events.Event, error)
```

SendRubinoPost shares a Rubino post to a chat.

### `SendRubinoProduct`

```go
func (a *Account) SendRubinoProduct(ctx context.Context, objectGUID, storeID, productID string, productVarieties []map[string]any) (*events.Event, error)
```

SendRubinoProduct shares a Rubino product to a chat.

### `SendRubinoStory`

```go
func (a *Account) SendRubinoStory(ctx context.Context, objectGUID, storyID, storyProfileID, replyText, storyType string, isMute bool) (*events.Event, error)
```

SendRubinoStory replies to or shares a Rubino story.

### `IsRubikaLink`

```go
func IsRubikaLink(s string) bool { return rubikaLinkRegex.MatchString(s) }
```

IsRubikaLink reports whether s contains a rubika.ir URL.

### `IsGroupLink`

```go
func IsGroupLink(s string) bool { return groupLinkRegex.MatchString(s) }
```

IsGroupLink reports whether s is a rubika group join link.

### `IsChannelLink`

```go
func IsChannelLink(s string) bool { return channelLinkRegex.MatchString(s) }
```

IsChannelLink reports whether s is a rubika channel join link.

### `IsPostLink`

```go
func IsPostLink(s string) bool { return postLinkRegex.MatchString(s) }
```

IsPostLink reports whether s is a rubika post link.

### `IsUsername`

```go
func IsUsername(s string) bool { return usernameRegex.MatchString(s) }
```

IsUsername reports whether s is a @username (4-32 alphanumeric/underscore).

### `GetRubikaLinks`

```go
func GetRubikaLinks(s string) []string
```

GetRubikaLinks extracts all rubika.ir URLs from s.

### `GetGroupLinks`

```go
func GetGroupLinks(s string) []string
```

GetGroupLinks extracts all group join links from s.

### `GetChannelLinks`

```go
func GetChannelLinks(s string) []string
```

GetChannelLinks extracts all channel join links from s.

### `GetPostLinks`

```go
func GetPostLinks(s string) []string
```

GetPostLinks extracts all post links from s.

### `GetUsernames`

```go
func GetUsernames(s string) []string
```

GetUsernames extracts all @usernames from s.

### `GetServerTime`

```go
func (a *Account) GetServerTime(ctx context.Context) (time.Time, error)
```

GetServerTime fetches the current server time. The server answers with a
timestamp (as a string or number); it is parsed robustly.
