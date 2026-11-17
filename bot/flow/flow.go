// Package flow provides a conversation flow manager for the bot
// API. Use it for wizards, forms and any
// dialogue that spans several messages.
//
// A flow is a named multi-step conversation. Each step has a
// callback that receives the user's input and returns the name of
// the next step (or "" to finish). Flows are scoped per-user so
// multiple users can be in different steps simultaneously.
package flow

import (
	"context"
	"sync"
)

// Bot is the minimal interface a flow client needs.
type Bot interface {
	SendMessage(ctx context.Context, chatID, text string, _ ...any) (any, error)
}

// Step is the signature of a single flow step.
type Step func(ctx context.Context, b Bot, chatID, userID, input string) (next string, err error)

// Flow is a named multi-step conversation.
type Flow struct {
	name    string
	steps   map[string]Step
	manager *Manager
}

// Manager tracks all active flows across all users.
type Manager struct {
	mu     sync.Mutex
	flows  map[string]*Flow       // flow name → Flow
	active map[string]*activeFlow // "user@flow" → active state
}

type activeFlow struct {
	flowName string
	stepName string
	data     map[string]string // per-user scratch space
}

// NewManager constructs a fresh flow Manager.
func NewManager() *Manager {
	return &Manager{
		flows:  map[string]*Flow{},
		active: map[string]*activeFlow{},
	}
}

// Register creates a new flow under the given name and returns it
// so the caller can add steps.
func (m *Manager) Register(name string) *Flow {
	f := &Flow{name: name, steps: map[string]Step{}, manager: m}
	m.mu.Lock()
	m.flows[name] = f
	m.mu.Unlock()
	return f
}

// Add registers a step in the flow.
func (f *Flow) Add(stepName string, step Step) {
	f.steps[stepName] = step
}

// Start begins the flow for a user at the "start" step.
func (m *Manager) Start(flowName, userID string) {
	key := userID + "@" + flowName
	m.mu.Lock()
	defer m.mu.Unlock()
	m.active[key] = &activeFlow{
		flowName: flowName,
		stepName: "start",
		data:     map[string]string{},
	}
}

// IsActive reports whether a flow is in progress for a user.
func (m *Manager) IsActive(flowName, userID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.active[userID+"@"+flowName]
	return ok
}

// IsActiveAny reports whether any flow is in progress for a user.
func (m *Manager) IsActiveAny(userID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k := range m.active {
		if len(k) > len(userID)+1 && k[:len(userID)+1] == userID+"@" {
			return true
		}
	}
	return false
}

// Stop ends any active flow for a user.
func (m *Manager) Stop(flowName, userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.active, userID+"@"+flowName)
}

// StopAll ends all flows for a user.
func (m *Manager) StopAll(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k := range m.active {
		if len(k) > len(userID)+1 && k[:len(userID)+1] == userID+"@" {
			delete(m.active, k)
		}
	}
}

// GetData retrieves a value from the user's flow scratch space.
func (m *Manager) GetData(flowName, userID, key string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	af, ok := m.active[userID+"@"+flowName]
	if !ok {
		return ""
	}
	return af.data[key]
}

// SetData stores a value in the user's flow scratch space.
func (m *Manager) SetData(flowName, userID, key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	af, ok := m.active[userID+"@"+flowName]
	if !ok {
		return
	}
	if af.data == nil {
		af.data = map[string]string{}
	}
	af.data[key] = value
}

// Handle processes an incoming message for a user. If the user has
// an active flow, the current step is invoked. Returns true when
// the message was consumed by a flow, false otherwise.
func (m *Manager) Handle(ctx context.Context, b Bot, chatID, userID, input string) (bool, error) {
	m.mu.Lock()
	var af *activeFlow
	for k, v := range m.active {
		if len(k) > len(userID)+1 && k[:len(userID)+1] == userID+"@" {
			af = v
			break
		}
	}
	m.mu.Unlock()
	if af == nil {
		return false, nil
	}
	f, ok := m.flows[af.flowName]
	if !ok {
		m.StopAll(userID)
		return false, nil
	}
	step, ok := f.steps[af.stepName]
	if !ok {
		m.StopAll(userID)
		return false, nil
	}
	next, err := step(ctx, b, chatID, userID, input)
	if err != nil {
		return true, err
	}
	if next == "" {
		m.StopAll(userID)
		return true, nil
	}
	m.mu.Lock()
	af.stepName = next
	m.mu.Unlock()
	return true, nil
}
