# `rubiai` — Methods

_Source: [`methods.go`](../../ParsRubika/rubiai/methods.go)_

## Symbols

### `Feedback`

```go
type Feedback string
```

Feedback is the user-sent opinion on an AI message.

### `FeedbackLike`

```go
const FeedbackLike Feedback = "Like"
```

### `FeedbackDislike`

```go
const FeedbackDislike Feedback = "Dislike"
```

### `FeedbackNone`

```go
const FeedbackNone Feedback = "None"
```

### `WelcomeLocation`

```go
type WelcomeLocation string
```

WelcomeLocation is the application area that getWelcome describes.

### `WelcomeVOD`

```go
const WelcomeVOD WelcomeLocation = "VOD"
```

### `WelcomeHome`

```go
const WelcomeHome WelcomeLocation = "Home"
```

### `WelcomeMessenger`

```go
const WelcomeMessenger WelcomeLocation = "Messenger"
```

### `WelcomeSocial`

```go
const WelcomeSocial WelcomeLocation = "Social"
```

### `GetWelcome`

```go
func (r *RubiaAI) GetWelcome(ctx context.Context, location string, locationID string) (*events.Event, error)
```

GetWelcome retrieves welcome information for an application location.
`location` accepts VOD/Home/Messenger/Social (case-insensitive).

### `GetServiceAccess`

```go
func (r *RubiaAI) GetServiceAccess(ctx context.Context) (*events.Event, error)
```

GetServiceAccess retrieves the current AI service access information.

### `GetFeedbackIssues`

```go
func (r *RubiaAI) GetFeedbackIssues(ctx context.Context) (*events.Event, error)
```

GetFeedbackIssues retrieves the available feedback issue types.

### `GetSessions`

```go
func (r *RubiaAI) GetSessions(ctx context.Context, allPages bool, limit int, startID string) (*events.Event, error)
```

GetSessions retrieves AI chat sessions. When allPages is true (the
default) the method auto-paginates following has_continue/next_start_id
and returns a single event wrapping {"sessions":[...]}. Otherwise it
returns the raw page for the supplied start_id/limit.

### `GetChatHistory`

```go
func (r *RubiaAI) GetChatHistory(ctx context.Context, sessionID string, allPages bool, limit int, startID string) (*events.Event, error)
```

GetChatHistory retrieves messages from an AI session. When allPages is
true (the default) the method auto-paginates and returns a single event
wrapping {"messages":[...]}. Otherwise it returns the raw page.

### `DeleteSession`

```go
func (r *RubiaAI) DeleteSession(ctx context.Context, sessionID string) (*events.Event, error)
```

DeleteSession deletes an AI chat session.

### `EditSession`

```go
func (r *RubiaAI) EditSession(ctx context.Context, sessionID string, title string, isPinned *bool) (*events.Event, error)
```

EditSession updates the pin state and/or title of a session. Pass the
zero value for a field you do not want to change (title="" / isPinned
left as its sentinel via the pointers).

### `SendFeedback`

```go
func (r *RubiaAI) SendFeedback(ctx context.Context, sessionID, messageID, feedback string) (*events.Event, error)
```

SendFeedback submits feedback (Like/Dislike/None) for an AI message.

### `SetFeedbackDetail`

```go
func (r *RubiaAI) SetFeedbackDetail(ctx context.Context, sessionID, messageID, issueType, issueDesc string) (*events.Event, error)
```

SetFeedbackDetail attaches a detailed issue type (and optional
description) to previously-submitted feedback. The issue types come
from GetFeedbackIssues.

### `StopSessionStream`

```go
func (r *RubiaAI) StopSessionStream(ctx context.Context, sessionID, messageID, feedback string) (*events.Event, error)
```

StopSessionStream stops an active AI session stream and submits the
final feedback (Like/Dislike; Good/Bad are accepted as aliases).

### `SpeechToText`

```go
func (r *RubiaAI) SpeechToText(ctx context.Context, voiceData []byte) (*events.Event, error)
```

SpeechToText converts raw voice bytes to text via the AI speech-to-text
service. The bytes are base64-encoded before being sent.

### `SpeechToTextFromFile`

```go
func (r *RubiaAI) SpeechToTextFromFile(ctx context.Context, path string) (*events.Event, error)
```

SpeechToTextFromFile reads the voice file at `path` and passes its
contents to SpeechToText.

### `DownloadImage`

```go
func (r *RubiaAI) DownloadImage(ctx context.Context, taskID, saveAs string, maxWait, interval time.Duration) (*events.Event, error)
```

DownloadImage retrieves an AI-generated image, polling until it is
ready. When `saveAs` is non-empty the decoded image bytes are written
there. Returns the final API response event (status Success/Error).

The polling loop honours ctx cancellation; maxWait and interval
default to 60 s / 2 s when zero.
