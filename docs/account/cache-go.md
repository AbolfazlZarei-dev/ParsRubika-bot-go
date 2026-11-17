# `account` — Cache

_Source: [`cache.go`](../../ParsRubika/account/cache.go)_

## Symbols

### `CacheEntry`

```go
type CacheEntry struct
```

CacheEntry holds a cached response and its expiry time.

### `ResponseCache`

```go
type ResponseCache struct
```

ResponseCache is a TTL-based in-memory cache for API responses.

### `NewResponseCache`

```go
func NewResponseCache(ttl time.Duration) *ResponseCache
```

NewResponseCache creates a ResponseCache with the given default TTL.

### `Close`

```go
func (c *ResponseCache) Close()
```

Close stops the background GC goroutine. Safe to call multiple times.

### `Get`

```go
func (c *ResponseCache) Get(key string) *events.Event
```

Get returns the cached response for key, or nil if not present/expired.

### `Set`

```go
func (c *ResponseCache) Set(key string, data *events.Event)
```

Set stores a response with the default TTL.

### `SetWithTTL`

```go
func (c *ResponseCache) SetWithTTL(key string, data *events.Event, ttl time.Duration)
```

SetWithTTL stores a response with a custom TTL.

### `Delete`

```go
func (c *ResponseCache) Delete(key string)
```

Delete removes a cached entry.

### `Clear`

```go
func (c *ResponseCache) Clear()
```

Clear removes all cached entries.

### `Size`

```go
func (c *ResponseCache) Size() int
```

Size returns the number of cached entries.

### `CachedRequest`

```go
func (a *Account) CachedRequest(ctx context.Context, cache *ResponseCache, key string, opts RequestOptions) (*events.Event, error)
```

CachedRequest is a helper that wraps a.Request with a cache lookup.
If the cache has a valid entry for the given key, it is returned
directly. Otherwise the request is made, the result is cached, and
then returned.

### `IdempotencyManager`

```go
type IdempotencyManager struct
```

IdempotencyManager tracks processed request IDs for deduplication.

### `NewIdempotencyManager`

```go
func NewIdempotencyManager(ttl time.Duration) *IdempotencyManager
```

NewIdempotencyManager creates an IdempotencyManager with the given TTL.
Processed request IDs are forgotten after the TTL expires.

### `Close`

```go
func (m *IdempotencyManager) Close()
```

Close stops the background GC goroutine. Safe to call multiple times.

### `IsUnique`

```go
func (m *IdempotencyManager) IsUnique(requestID string) bool
```

IsUnique returns true if the requestID has not been seen before (and
marks it as seen). Honours the TTL inline — an entry whose stored
timestamp is older than the TTL is treated as unseen.

### `MarkAsProcessed`

```go
func (m *IdempotencyManager) MarkAsProcessed(requestID string)
```

MarkAsProcessed explicitly marks a request ID as processed.

### `Forget`

```go
func (m *IdempotencyManager) Forget(requestID string)
```

Forget removes a request ID from the manager.

### `QueuedMessage`

```go
type QueuedMessage struct
```

QueuedMessage holds a message destined for a chat.

### `MessageQueue`

```go
type MessageQueue struct
```

MessageQueue is a bounded async queue for outgoing messages.

### `NewMessageQueue`

```go
func NewMessageQueue(a *Account, bufSize int, concurrency int) *MessageQueue
```

NewMessageQueue creates a MessageQueue with the given buffer size and
starts a worker goroutine that drains the queue.

### `ErrQueueClosed`

```go
var ErrQueueClosed = fmt.Errorf("account: message queue is closed")
```

ErrQueueClosed is returned by Submit when the queue has been stopped.

### `ErrQueueFull`

```go
var ErrQueueFull = fmt.Errorf("account: message queue is full")
```

ErrQueueFull is returned by Submit when the queue buffer is full.

### `Submit`

```go
func (mq *MessageQueue) Submit(objectGUID, text string, opts ...SendOptions) <-chan error
```

Submit enqueues a message for sending. Returns a channel that will
receive the error (or nil on success). The channel is buffered so it
does not block the worker.

Submit is safe to call concurrently with Stop: the closed flag is
checked under mu, and the send is performed while holding mu, so Stop
can never close(mq.ch) between Submit's check and its send.

### `Stop`

```go
func (mq *MessageQueue) Stop()
```

Stop shuts down the queue, drains any buffered messages (sending them
an error so callers don't block forever), and waits for in-flight
workers to exit. Safe to call multiple times.

### `Pending`

```go
func (mq *MessageQueue) Pending() int { return len(mq.ch) }
```

Pending returns the number of messages waiting in the queue.
