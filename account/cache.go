package account

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// ===== Response caching ===============================================
//
// ResponseCache is a simple TTL cache for API responses. It is useful
// for endpoints that return rarely-changing data (e.g. getChats,
// getUserInfo) to reduce the number of API calls. The cache is safe for
// concurrent use.

// CacheEntry holds a cached response and its expiry time.
type CacheEntry struct {
	data      *events.Event
	expiresAt time.Time
}

// ResponseCache is a TTL-based in-memory cache for API responses.
type ResponseCache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
	ttl     time.Duration
	stop    chan struct{}
}

// NewResponseCache creates a ResponseCache with the given default TTL.
func NewResponseCache(ttl time.Duration) *ResponseCache {
	c := &ResponseCache{
		entries: make(map[string]CacheEntry),
		ttl:     ttl,
		stop:    make(chan struct{}),
	}
	if ttl > 0 {
		go c.gcLoop(ttl)
	}
	return c
}

// Close stops the background GC goroutine. Safe to call multiple times.
func (c *ResponseCache) Close() {
	select {
	case <-c.stop:
		// already closed
	default:
		close(c.stop)
	}
}

// Get returns the cached response for key, or nil if not present/expired.
func (c *ResponseCache) Get(key string) *events.Event {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || time.Now().After(e.expiresAt) {
		return nil
	}
	return e.data
}

// Set stores a response with the default TTL.
func (c *ResponseCache) Set(key string, data *events.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{
		data:      data,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// SetWithTTL stores a response with a custom TTL.
func (c *ResponseCache) SetWithTTL(key string, data *events.Event, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{
		data:      data,
		expiresAt: time.Now().Add(ttl),
	}
}

// Delete removes a cached entry.
func (c *ResponseCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Clear removes all cached entries.
func (c *ResponseCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]CacheEntry)
}

// Size returns the number of cached entries.
func (c *ResponseCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// gcLoop periodically evicts expired entries. Exits when Close() is
// called or the stop channel is closed.
func (c *ResponseCache) gcLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for k, e := range c.entries {
				if now.After(e.expiresAt) {
					delete(c.entries, k)
				}
			}
			c.mu.Unlock()
		}
	}
}

// CachedRequest is a helper that wraps a.Request with a cache lookup.
// If the cache has a valid entry for the given key, it is returned
// directly. Otherwise the request is made, the result is cached, and
// then returned.
func (a *Account) CachedRequest(ctx context.Context, cache *ResponseCache, key string, opts RequestOptions) (*events.Event, error) {
	if cache != nil {
		if cached := cache.Get(key); cached != nil {
			return cached, nil
		}
	}
	resp, err := a.Request(ctx, opts)
	if err != nil {
		return nil, err
	}
	if cache != nil {
		cache.Set(key, resp)
	}
	return resp, nil
}

// ===== Request de-duplication =========================================
//
// IdempotencyManager tracks request IDs so that retried requests are not
// executed twice. This is useful for non-idempotent operations like
// sendMessage where a network timeout might cause the caller to retry.

// IdempotencyManager tracks processed request IDs for deduplication.
type IdempotencyManager struct {
	mu   sync.Mutex
	seen map[string]time.Time
	ttl  time.Duration
	stop chan struct{}
}

// NewIdempotencyManager creates an IdempotencyManager with the given TTL.
// Processed request IDs are forgotten after the TTL expires.
func NewIdempotencyManager(ttl time.Duration) *IdempotencyManager {
	m := &IdempotencyManager{
		seen: make(map[string]time.Time),
		ttl:  ttl,
		stop: make(chan struct{}),
	}
	if ttl > 0 {
		go m.gcLoop(ttl)
	}
	return m
}

// Close stops the background GC goroutine. Safe to call multiple times.
func (m *IdempotencyManager) Close() {
	select {
	case <-m.stop:
	default:
		close(m.stop)
	}
}

// IsUnique returns true if the requestID has not been seen before (and
// marks it as seen). Honours the TTL inline — an entry whose stored
// timestamp is older than the TTL is treated as unseen.
func (m *IdempotencyManager) IsUnique(requestID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t, seen := m.seen[requestID]; seen {
		if m.ttl > 0 && time.Since(t) > m.ttl {
			// expired — treat as unseen and re-mark
			m.seen[requestID] = time.Now()
			return true
		}
		return false
	}
	m.seen[requestID] = time.Now()
	return true
}

// MarkAsProcessed explicitly marks a request ID as processed.
func (m *IdempotencyManager) MarkAsProcessed(requestID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seen[requestID] = time.Now()
}

// Forget removes a request ID from the manager.
func (m *IdempotencyManager) Forget(requestID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.seen, requestID)
}

// gcLoop periodically evicts expired entries. Exits when Close() is
// called.
func (m *IdempotencyManager) gcLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-m.stop:
			return
		case <-ticker.C:
			m.mu.Lock()
			now := time.Now()
			for k, t := range m.seen {
				if now.Sub(t) > m.ttl {
					delete(m.seen, k)
				}
			}
			m.mu.Unlock()
		}
	}
}

// ===== Outgoing message queue =========================================
//
// MessageQueue is a bounded async queue for outgoing messages. It
// decouples producers from the network so a burst of sends does not
// overwhelm the rate limiter.

// QueuedMessage holds a message destined for a chat.
type QueuedMessage struct {
	ObjectGUID string
	Text       string
	Opts       []SendOptions
	Result     chan error
}

// MessageQueue is a bounded async queue for outgoing messages.
type MessageQueue struct {
	ch        chan *QueuedMessage
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	closeOnce sync.Once
	done      chan struct{}
	// mu guards the `closed` flag. Submit holds it across the send so
	// that Stop can set closed=true (and then close(mq.ch)) atomically
	// with respect to any in-flight Submit — eliminating the TOCTOU
	// window that would otherwise let Submit send on a closed channel
	// and panic.
	mu     sync.Mutex
	closed bool
}

// NewMessageQueue creates a MessageQueue with the given buffer size and
// starts a worker goroutine that drains the queue.
func NewMessageQueue(a *Account, bufSize int, concurrency int) *MessageQueue {
	if bufSize <= 0 {
		bufSize = 100
	}
	if concurrency <= 0 {
		concurrency = 1
	}
	ctx, cancel := context.WithCancel(context.Background())
	mq := &MessageQueue{
		ch:     make(chan *QueuedMessage, bufSize),
		cancel: cancel,
		done:   make(chan struct{}),
	}
	for i := 0; i < concurrency; i++ {
		mq.wg.Add(1)
		go func() {
			defer mq.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-mq.ch:
					if !ok {
						return
					}
					_, err := a.SendMessage(ctx, msg.ObjectGUID, msg.Text, msg.Opts...)
					msg.Result <- err
				}
			}
		}()
	}
	return mq
}

// ErrQueueClosed is returned by Submit when the queue has been stopped.
var ErrQueueClosed = fmt.Errorf("account: message queue is closed")

// ErrQueueFull is returned by Submit when the queue buffer is full.
var ErrQueueFull = fmt.Errorf("account: message queue is full")

// Submit enqueues a message for sending. Returns a channel that will
// receive the error (or nil on success). The channel is buffered so it
// does not block the worker.
//
// Submit is safe to call concurrently with Stop: the closed flag is
// checked under mu, and the send is performed while holding mu, so Stop
// can never close(mq.ch) between Submit's check and its send.
func (mq *MessageQueue) Submit(objectGUID, text string, opts ...SendOptions) <-chan error {
	result := make(chan error, 1)
	mq.mu.Lock()
	if mq.closed {
		mq.mu.Unlock()
		result <- ErrQueueClosed
		return result
	}
	msg := &QueuedMessage{
		ObjectGUID: objectGUID,
		Text:       text,
		Opts:       opts,
		Result:     result,
	}
	// Non-blocking send (the select has a default), so holding mu here
	// is safe — it never blocks waiting for a worker.
	select {
	case mq.ch <- msg:
	default:
		mq.mu.Unlock()
		result <- ErrQueueFull
		return result
	}
	mq.mu.Unlock()
	return result
}

// Stop shuts down the queue, drains any buffered messages (sending them
// an error so callers don't block forever), and waits for in-flight
// workers to exit. Safe to call multiple times.
func (mq *MessageQueue) Stop() {
	mq.closeOnce.Do(func() {
		// Set closed=true under mu BEFORE closing mq.ch. Any Submit that
		// holds mu after this point sees closed=true and rejects; any
		// Submit that already released mu has already completed its send
		// against an open channel. There is no window in which a send
		// lands on a closed channel.
		mq.mu.Lock()
		mq.closed = true
		close(mq.done)
		mq.mu.Unlock()
		mq.cancel()  // signal workers to stop
		close(mq.ch) // unblock workers waiting on the channel
	})
	// Drain remaining buffered messages so their callers don't block
	// forever on <-result. Workers may have already exited, so we
	// drain here rather than relying on them.
	for msg := range mq.ch {
		msg.Result <- ErrQueueClosed
	}
	mq.wg.Wait()
}

// Pending returns the number of messages waiting in the queue.
func (mq *MessageQueue) Pending() int { return len(mq.ch) }
