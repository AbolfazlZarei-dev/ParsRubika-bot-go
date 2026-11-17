package bot

import (
	"context"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/rules"
)

// ===== FSM / State Manager =====================
//
// StateManager is a per-user in-memory finite-state machine and the
// simplest way to build multi-step bot flows (forms, wizards, menus).
//
// The manager is safe for concurrent use. States are scoped per userID
// (which is the sender_id for private chats, or chat_id+sender_id for
// groups). A state value is an opaque string — the bot author decides
// the vocabulary (e.g. "awaiting_phone", "menu:settings").

// StateManager holds per-user state values with optional TTL.
type StateManager struct {
	mu     sync.RWMutex
	states map[string]stateEntry
	ttl    time.Duration // 0 = no expiry
	done   chan struct{}
	once   sync.Once
}

type stateEntry struct {
	value     string
	updatedAt time.Time
	data      map[string]string // optional per-state scratch data
}

// NewStateManager creates a StateManager. If ttl > 0, state entries
// expire after that duration of inactivity.
func NewStateManager(ttl time.Duration) *StateManager {
	sm := &StateManager{
		states: make(map[string]stateEntry),
		ttl:    ttl,
		done:   make(chan struct{}),
	}
	if ttl > 0 {
		go sm.gcLoop(ttl)
	}
	return sm
}

// Set sets the state for userID.
func (sm *StateManager) Set(userID, state string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.states[userID] = stateEntry{
		value:     state,
		updatedAt: time.Now(),
		data:      make(map[string]string),
	}
}

// Get returns the current state for userID ("" if none).
func (sm *StateManager) Get(userID string) string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	e, ok := sm.states[userID]
	if !ok {
		return ""
	}
	if sm.ttl > 0 && time.Since(e.updatedAt) > sm.ttl {
		return ""
	}
	return e.value
}

// Clear removes the state for userID.
func (sm *StateManager) Clear(userID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.states, userID)
}

// SetData stores a key/value pair in the per-user state scratch area.
// The state is created if it does not exist.
func (sm *StateManager) SetData(userID, key, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	e, ok := sm.states[userID]
	if !ok {
		e = stateEntry{updatedAt: time.Now(), data: make(map[string]string)}
	}
	e.data[key] = value
	e.updatedAt = time.Now()
	sm.states[userID] = e
}

// GetData retrieves a key from the per-user state scratch area.
func (sm *StateManager) GetData(userID, key string) string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	e, ok := sm.states[userID]
	if !ok {
		return ""
	}
	return e.data[key]
}

// HasState returns true if the user has any state set (and not expired).
func (sm *StateManager) HasState(userID string) bool {
	return sm.Get(userID) != ""
}

// ListAll returns a snapshot of all active states (for debugging).
func (sm *StateManager) ListAll() map[string]string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	out := make(map[string]string, len(sm.states))
	for uid, e := range sm.states {
		if sm.ttl > 0 && time.Since(e.updatedAt) > sm.ttl {
			continue
		}
		out[uid] = e.value
	}
	return out
}

// Close stops the background expiry goroutine. It is safe to call more
// than once.
func (sm *StateManager) Close() { sm.once.Do(func() { close(sm.done) }) }

// gcLoop periodically evicts expired entries.
func (sm *StateManager) gcLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-sm.done:
			return
		case <-ticker.C:
		}
		sm.mu.Lock()
		now := time.Now()
		for uid, e := range sm.states {
			if sm.ttl > 0 && now.Sub(e.updatedAt) > sm.ttl {
				delete(sm.states, uid)
			}
		}
		sm.mu.Unlock()
	}
}

// StateRule is a Rule that matches when the user is in (or not in) a
// given state.
type StateRule struct {
	manager *StateManager
	states  []string
	invert  bool
	scope   StateScope
}

// StateScope determines whether the state key is the sender_id (User)
// or the chat_id (Chat) or a combination (Both).
type StateScope int

const (
	StateScopeUser StateScope = iota
	StateScopeChat
	StateScopeBoth
)

// NewStateRule creates a rule that matches when the user is in one of
// the given states. Pass invert=true to match when NOT in any of them.
func NewStateRule(manager *StateManager, scope StateScope, invert bool, states ...string) *StateRule {
	return &StateRule{manager: manager, states: states, invert: invert, scope: scope}
}

// Check implements rules.Rule.
func (s *StateRule) Check(ctx context.Context, upd *models.Update) bool {
	key := s.keyFor(upd)
	current := s.manager.Get(key)
	matched := false
	for _, st := range s.states {
		if current == st {
			matched = true
			break
		}
	}
	if s.invert {
		return !matched
	}
	return matched
}

func (s *StateRule) keyFor(upd *models.Update) string {
	switch s.scope {
	case StateScopeChat:
		return upd.ChatID
	case StateScopeBoth:
		return upd.ChatID + ":" + senderOf(upd)
	default:
		return senderOf(upd)
	}
}

// senderOf extracts the sender_id from an update.
func senderOf(upd *models.Update) string {
	if upd.NewMessage != nil {
		return upd.NewMessage.SenderID
	}
	if upd.UpdatedMessage != nil {
		return upd.UpdatedMessage.SenderID
	}
	if upd.ExternalData != nil {
		if sid, ok := upd.ExternalData["sender_id"].(string); ok {
			return sid
		}
	}
	return ""
}

// Compile-time assertion that StateRule implements rules.Rule.
var _ rules.Rule = (*StateRule)(nil)

// ===== Conversation manager ===========================================

// Conversation is a multi-step flow tied to a single user. It tracks the
// current step name and optional per-step data.
type Conversation struct {
	userID  string
	manager *StateManager
	steps   map[string]ConversationStep
	current string
}

// ConversationStep is a single step in a conversation flow.
type ConversationStep struct {
	Name     string
	Prompt   string
	Validate func(input string) error
	Process  func(input string) (nextStep string, err error)
}

// ConversationManager owns all active conversations and their step
// definitions.
type ConversationManager struct {
	mu       sync.Mutex
	flows    map[string]map[string]ConversationStep // flowName → steps
	active   map[string]*Conversation               // userID → active
	stateMgr *StateManager
}

// NewConversationManager creates a new ConversationManager.
func NewConversationManager(stateMgr *StateManager) *ConversationManager {
	return &ConversationManager{
		flows:    make(map[string]map[string]ConversationStep),
		active:   make(map[string]*Conversation),
		stateMgr: stateMgr,
	}
}

// RegisterFlow registers a named conversation flow with its steps.
func (cm *ConversationManager) RegisterFlow(name string, steps []ConversationStep) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	m := make(map[string]ConversationStep, len(steps))
	for _, s := range steps {
		m[s.Name] = s
	}
	cm.flows[name] = m
}

// Start begins a conversation flow for userID. The first step's prompt
// is returned.
func (cm *ConversationManager) Start(userID, flowName string) (string, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	flow, ok := cm.flows[flowName]
	if !ok {
		return "", &APIError{Status: "UnknownFlow", DevMessage: flowName}
	}
	// Pick the first step (map iteration order is random; use the
	// first step from the original slice by re-scanning).
	var firstStep ConversationStep
	for _, s := range flow {
		firstStep = s
		break
	}
	conv := &Conversation{
		userID:  userID,
		manager: cm.stateMgr,
		steps:   flow,
		current: firstStep.Name,
	}
	cm.active[userID] = conv
	cm.stateMgr.Set(userID, "conv:"+flowName+":"+firstStep.Name)
	return firstStep.Prompt, nil
}

// Handle processes an input for an active conversation. Returns the next
// step's prompt (or "" if the conversation is complete) and an error.
func (cm *ConversationManager) Handle(userID, input string) (string, error) {
	cm.mu.Lock()
	conv, ok := cm.active[userID]
	cm.mu.Unlock()
	if !ok {
		return "", &APIError{Status: "NoActiveConversation", DevMessage: userID}
	}
	step, ok := conv.steps[conv.current]
	if !ok {
		return "", &APIError{Status: "InvalidStep", DevMessage: conv.current}
	}
	if step.Validate != nil {
		if err := step.Validate(input); err != nil {
			return step.Prompt, err
		}
	}
	next, err := step.Process(input)
	if err != nil {
		return step.Prompt, err
	}
	if next == "" {
		// Conversation complete
		cm.mu.Lock()
		delete(cm.active, userID)
		cm.mu.Unlock()
		cm.stateMgr.Clear(userID)
		return "", nil
	}
	conv.current = next
	cm.stateMgr.Set(userID, "conv:"+next)
	if nextStep, ok := conv.steps[next]; ok {
		return nextStep.Prompt, nil
	}
	return "", nil
}

// IsActive returns true if the user has an active conversation.
func (cm *ConversationManager) IsActive(userID string) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	_, ok := cm.active[userID]
	return ok
}

// Stop forcibly ends a conversation.
func (cm *ConversationManager) Stop(userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.active, userID)
	cm.stateMgr.Clear(userID)
}
