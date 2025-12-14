package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

import (
	"log"
	"sync"
)

// StateManager برای مدیریت وضعیت‌های کاربران
type StateManager struct {
	states map[string]map[string]interface{} // نگهداری وضعیت هر کاربر به صورت map درون map
	mu     sync.RWMutex                      // قفل برای دسترسی امن همزمان (خواندن/نوشتن)
}

// NewStateManager یک نمونه جدید از StateManager ایجاد می‌کند
func NewStateManager() *StateManager {
	return &StateManager{
		states: make(map[string]map[string]interface{}),
	}
}

// SetState یک مقدار را برای کاربر مشخص ذخیره می‌کند
func (sm *StateManager) SetState(userID, key string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// اگر کاربر برای اولین بار است، یک map جدید برای او ایجاد کن
	if _, ok := sm.states[userID]; !ok {
		sm.states[userID] = make(map[string]interface{})
	}
	sm.states[userID][key] = value
	log.Printf("State: Set for user %s, key '%s' with value %v", userID, key, value)
}

// GetState یک مقدار را برای کاربر و کلید مشخص برمی‌گرداند
func (sm *StateManager) GetState(userID, key string) (interface{}, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if userState, ok := sm.states[userID]; ok {
		if value, ok := userState[key]; ok {
			log.Printf("State: Get for user %s, key '%s' found", userID, key)
			return value, true
		}
	}
	log.Printf("State: Get for user %s, key '%s' not found", userID, key)
	return nil, false
}

// DeleteState یک کلید را برای کاربر مشخص حذف می‌کند
func (sm *StateManager) DeleteState(userID, key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if userState, ok := sm.states[userID]; ok {
		delete(userState, key)
		log.Printf("State: Deleted for user %s, key '%s'", userID, key)
	}
}

// DeleteUserState تمام وضعیت‌های یک کاربر را حذف می‌کند
func (sm *StateManager) DeleteUserState(userID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.states, userID)
	log.Printf("State: All states deleted for user %s", userID)
}

// GetAllUserStates تمام وضعیت‌های یک کاربر را برمی‌گرداند
func (sm *StateManager) GetAllUserStates(userID string) (map[string]interface{}, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	userState, ok := sm.states[userID]
	if !ok {
		return nil, false
	}

	// یک کپی از وضعیت‌ها برمی‌گردانیم تا از تغییرات خارجی جلوگیری شود
	copyState := make(map[string]interface{})
	for k, v := range userState {
		copyState[k] = v
	}
	return copyState, true
}
