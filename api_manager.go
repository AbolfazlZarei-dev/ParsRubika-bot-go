package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// APIManager مدیریت APIها و سوییچ هوشمند بین آن‌ها
type APIManager struct {
	currentAPI     APIType
	healthStatus   map[APIType]*HealthStatus
	mu             sync.RWMutex
	client         *BotClient
	switchHistory  []APIEvent
	healthTicker   *time.Ticker
	stopHealthChan chan struct{}
}

// NewAPIManager ایجاد مدیر API
func NewAPIManager(client *BotClient) *APIManager {
	return &APIManager{
		currentAPI: BotAPI, // به طور پیش‌فرض از Bot API شروع می‌کنیم
		healthStatus: map[APIType]*HealthStatus{
			BotAPI: {
				APIType:      BotAPI,
				IsHealthy:    true,
				ResponseTime: 0,
				LastChecked:  time.Now(),
				ErrorCount:   0,
			},
			MessengerAPI: {
				APIType:      MessengerAPI,
				IsHealthy:    true,
				ResponseTime: 0,
				LastChecked:  time.Now(),
				ErrorCount:   0,
			},
		},
		client:         client,
		switchHistory:  make([]APIEvent, 0),
		stopHealthChan: make(chan struct{}),
	}
}

// GetCurrentAPI دریافت API فعلی
func (am *APIManager) GetCurrentAPI() APIType {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.currentAPI
}

// SwitchAPI سوییچ کردن API (منطق جدید و غیرمسدودکننده)
func (am *APIManager) SwitchAPI(newAPI APIType, reason string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.currentAPI == newAPI {
		return nil
	}

	event := APIEvent{
		Timestamp: time.Now(),
		FromAPI:   am.currentAPI,
		ToAPI:     newAPI,
		Reason:    reason,
	}
	am.switchHistory = append(am.switchHistory, event)

	log.Printf("🔄 در حال سوییچ از %s به %s به دلیل: %s", am.currentAPI, newAPI, reason)

	// 1. ابتدا وضعیت داخلی API را فوری تغییر می‌دهیم
	am.currentAPI = newAPI

	log.Printf("✅ وضعیت داخلی با موفقیت به %s سوییچ شد", newAPI)

	// 2. سپس اطلاع‌رسانی‌ها را در یک goroutine جداگانه ارسال می‌کنیم
	// تا اگر با Timeout مواجه شدند، عملکرد اصلی بات متوقف نشود.
	if am.client.notificationOpts != nil && am.client.notificationOpts.Enabled {
		go func() {
			// ارسال اطلاع‌رسانی "در حال سوییچ"
			am.client.sendAPIChangeNotification(
				fmt.Sprintf("🔔 در حال سوییچ از %s به %s به دلیل: %s", event.FromAPI, event.ToAPI, event.Reason),
			)
			// ارسال اطلاع‌رسانی "سوییچ موفق"
			am.client.sendAPIChangeNotification(
				fmt.Sprintf("✅ با موفقیت به %s سوییچ شد", event.ToAPI),
			)
		}()
	}

	return nil
}

// CheckHealth بررسی سلامت API
func (am *APIManager) CheckHealth(ctx context.Context, apiType APIType) error {
	start := time.Now()

	var url string
	if apiType == BotAPI {
		url = fmt.Sprintf("%s/%s/getMe", am.client.baseURL, am.client.token)
	} else {
		url = fmt.Sprintf("%s/%s/getMe", am.client.messengerURL, am.client.token)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		am.updateHealthStatus(apiType, time.Since(start), false, err)
		return err
	}

	resp, err := am.client.httpClient.Do(req)
	if err != nil {
		am.updateHealthStatus(apiType, time.Since(start), false, err)
		return err
	}
	defer resp.Body.Close()

	responseTime := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		apiErr := fmt.Errorf("HTTP status: %d", resp.StatusCode)
		am.updateHealthStatus(apiType, responseTime, false, apiErr)
		return apiErr
	}

	// اگر کد وضعیت 200 OK باشد، این یک نشانه خوب برای سلامت API است.
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// خطای پارس کردن را برای دیباگ لاگ می‌کنیم، اما سلامت API را رد نمی‌کنیم.
		log.Printf("اخطار: خطا در پارس کردن JSON از %s (وضعیت HTTP: %d): %v. API به عنوان سالم در نظر گرفته می‌شود.", apiType, resp.StatusCode, err)
		am.updateHealthStatus(apiType, responseTime, true, nil)
		return nil
	}

	if status, ok := result["status"].(string); ok && status != "OK" {
		apiErr := fmt.Errorf("non-OK status: %s", status)
		am.updateHealthStatus(apiType, responseTime, false, apiErr)
		return apiErr
	}

	am.updateHealthStatus(apiType, responseTime, true, nil)
	return nil
}

// updateHealthStatus به‌روزرسانی وضعیت سلامت
func (am *APIManager) updateHealthStatus(apiType APIType, responseTime time.Duration, success bool, err error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	status := am.healthStatus[apiType]
	status.ResponseTime = responseTime
	status.LastChecked = time.Now()
	status.LastError = err

	if success {
		status.IsHealthy = true
		// اگر درخواست موفق بود، شمارنده خطا را ریست می‌کنیم
		status.ErrorCount = 0
	} else {
		status.IsHealthy = false
		status.ErrorCount++
	}

	// ثبت رویداد در شبکه برای تحلیل هوشمند
	am.client.networkStabilityManager.RecordEvent(apiType, responseTime, success, err)
}

// GetHealthStatus دریافت وضعیت سلامت API
func (am *APIManager) GetHealthStatus(apiType APIType) *HealthStatus {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.healthStatus[apiType]
}

// GetSwitchHistory دریافت تاریخچه سوییچ‌ها
func (am *APIManager) GetSwitchHistory() []APIEvent {
	am.mu.RLock()
	defer am.mu.RUnlock()
	history := make([]APIEvent, len(am.switchHistory))
	copy(history, am.switchHistory)
	return history
}

// StartHealthMonitoring شروع نظارت بر سلامت APIها
func (am *APIManager) StartHealthMonitoring(ctx context.Context, interval time.Duration) {
	am.healthTicker = time.NewTicker(interval)
	defer am.healthTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-am.stopHealthChan:
			log.Println("نظارت بر سلامت API متوقف شد.")
			return
		case <-am.healthTicker.C:
			var wg sync.WaitGroup
			// بررسی سلامت هر دو API به صورت همزمان
			for _, apiType := range []APIType{BotAPI, MessengerAPI} {
				wg.Add(1)
				go func(at APIType) {
					defer wg.Done()
					if err := am.CheckHealth(ctx, at); err != nil {
						log.Printf("خطا در بررسی سلامت %s: %v", at, err)
					}
				}(apiType)
			}

			// منتظر می‌مانیم تا هر دو بررسی سلامت تمام شوند
			wg.Wait()

			// حالا با اطلاعات به‌روز، نیاز به سوییچ را بررسی می‌کنیم
			if am.client.connectionMode == SwitcherMode {
				am.checkForSwitch()
			}
		}
	}
}

// StopHealthMonitoring توقف نظارت بر سلامت
func (am *APIManager) StopHealthMonitoring() {
	if am.healthTicker != nil {
		am.healthTicker.Stop()
	}
	close(am.stopHealthChan)
}

// checkForSwitch بررسی نیاز به سوییچ API (منطق جدید و دقیق)
func (am *APIManager) checkForSwitch() {
	botStatus := am.healthStatus[BotAPI]
	messengerStatus := am.healthStatus[MessengerAPI]

	// --- منطق وقتی روی Bot API هستیم ---
	if am.currentAPI == BotAPI {
		// اگر Bot API ناسالم، کند یا دارای خطاهای مکرر است
		if !botStatus.IsHealthy || botStatus.ResponseTime > 5*time.Second || botStatus.ErrorCount > 3 {
			if messengerStatus.IsHealthy {
				log.Printf("🚨 Bot API مشکل دارد (سالم: %v, تأخیر: %v, خطاها: %d). در حال سوییچ به Messenger API...", botStatus.IsHealthy, botStatus.ResponseTime, botStatus.ErrorCount)
				am.SwitchAPI(MessengerAPI, "Bot API ناسالم یا کند بود")
			} else {
				log.Printf("⚠️ هر دو API در دسترس نیستند. روی Bot API باقی می‌مانم.")
			}
		} else {
			log.Printf("✅ Bot API سالم و پایدار است. نیازی به سوییچ نیست.")
		}
		return
	}

	// --- منطق وقتی روی Messenger API هستیم ---
	if am.currentAPI == MessengerAPI {
		// اگر Bot API دوباره کاملاً سالم و سریع شد، به آن بازمی‌گردیم
		if botStatus.IsHealthy && botStatus.ResponseTime < 2*time.Second && botStatus.ErrorCount == 0 {
			log.Printf("🎉 Bot API دوباره سالم و سریع شد (تأخیر: %v). در حال بازگشت به Bot API...", botStatus.ResponseTime)
			am.SwitchAPI(BotAPI, "Bot API دوباره سالم و سریع شد")
		} else {
			log.Printf("📡 هنوز روی Messenger API هستیم. Bot API هنوز آماده بازگشت نیست (سالم: %v, تأخیر: %v, خطاها: %d).", botStatus.IsHealthy, botStatus.ResponseTime, botStatus.ErrorCount)
		}
		return
	}
}
