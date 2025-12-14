package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

import (
	"sync"
	"time"
)

// AntiSpam برای جلوگیری از اسپم و درخواست‌های مکرر کاربران
type AntiSpam struct {
	requests map[string]time.Time // نگهداری زمان آخرین درخواست هر کاربر
	mu       sync.RWMutex         // قفل برای دسترسی امن همزمان
	cooldown time.Duration        // زمان کول‌داون بین درخواست‌ها
}

// NewAntiSpam ایجاد یک نمونه جدید از AntiSpam
func NewAntiSpam() *AntiSpam {
	return &AntiSpam{
		requests: make(map[string]time.Time),
		cooldown: 3 * time.Second, // به طور پیش‌فرض 3 ثانیه بین هر درخواست برای یک کاربر
	}
}

// SetCooldown تنظیم زمان کول‌داون
func (as *AntiSpam) SetCooldown(duration time.Duration) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.cooldown = duration
}

// Check بررسی اینکه آیا کاربر اسپم می‌کند یا خیر
// اگر کاربر اسپم نکند، true برمی‌گرداند. در غیر این صورت false برمی‌گرداند.
func (as *AntiSpam) Check(userID string) bool {
	as.mu.Lock()
	defer as.mu.Unlock()

	now := time.Now()
	lastRequest, exists := as.requests[userID]

	// اگر کاربر قبلاً درخواست نداده یا زمان کول‌داون گذشته باشد، اجازه بده
	if !exists || now.Sub(lastRequest) >= as.cooldown {
		as.requests[userID] = now
		return true
	}

	// کاربر در حال اسپم کردن است
	return false
}

// Reset بازنشانی وضعیت ضد اسپم برای یک کاربر خاص
func (as *AntiSpam) Reset(userID string) {
	as.mu.Lock()
	defer as.mu.Unlock()
	delete(as.requests, userID)
}
