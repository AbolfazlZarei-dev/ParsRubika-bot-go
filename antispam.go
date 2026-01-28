package ParsRubika

// نسخه: 2.5.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go

import (
	"sync"
	"time"
)

// AntiSpam برای جلوگیری از اسپم و درخواست‌های مکرر کاربران
type AntiSpam struct {
	requests map[string]time.Time // نگهداری زمان آخرین درخواست هر کاربر
	mu       sync.RWMutex         // قفل برای دسترسی امن همزمان
	cooldown time.Duration        // زمان کول‌داون بین درخواست‌ها
	stopChan chan struct{}        // کانال برای توقف روتین پاکسازی
}

// NewAntiSpam ایجاد یک نمونه جدید از AntiSpam
func NewAntiSpam() *AntiSpam {
	as := &AntiSpam{
		requests: make(map[string]time.Time),
		cooldown: 3 * time.Second, // به طور پیش‌فرض 3 ثانیه بین هر درخواست برای یک کاربر
		stopChan: make(chan struct{}),
	}

	// شروع گاربج کالر برای جلوگیری از نشت حافظه (Memory Leak)
	go as.cleanupRoutine()

	return as
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

// Stop تابعی برای توقف گاربج کالر (بهتر است در متد Stop کلاینت بات فراخوانی شود)
func (as *AntiSpam) Stop() {
	close(as.stopChan)
}

// cleanupRoutine روتین دوره‌ای برای پاکسازی کاربرانی که مدتیست فعالیت نداشتند
func (as *AntiSpam) cleanupRoutine() {
	// هر 5 دقیقه یکبار پاکسازی انجام می‌شود
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-as.stopChan:
			return
		case <-ticker.C:
			as.cleanupOldEntries()
		}
	}
}

// cleanupOldEntries حذف ورودی‌های قدیمی از مپ برای آزادسازی حافظه
func (as *AntiSpam) cleanupOldEntries() {
	as.mu.Lock()
	defer as.mu.Unlock()

	now := time.Now()
	// حذف کاربرانی که در 30 دقیقه گذشته درخواستی نداشته‌اند
	// این زمان را می‌توان بسته به نیاز کم یا زیاد کرد
	for userID, lastRequest := range as.requests {
		if now.Sub(lastRequest) > 30*time.Minute {
			delete(as.requests, userID)
		}
	}
}
