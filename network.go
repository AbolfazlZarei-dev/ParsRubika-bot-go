package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

import (
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

// NetworkStabilityManager مدیریت پایداری شبکه به صورت هوشمند
type NetworkStabilityManager struct {
	client          *BotClient
	history         []NetworkEvent
	mu              sync.Mutex
	predictionModel *PredictionModel
}

// PredictionModel مدل پیش‌بینی برای شبکه
type PredictionModel struct {
	samples []NetworkSample
	mu      sync.Mutex
}

// NewNetworkStabilityManager ایجاد یک نمونه جدید
func NewNetworkStabilityManager(client *BotClient) *NetworkStabilityManager {
	return &NetworkStabilityManager{
		client:          client,
		history:         make([]NetworkEvent, 0),
		predictionModel: NewPredictionModel(),
	}
}

// NewPredictionModel ایجاد مدل پیش‌بینی جدید
func NewPredictionModel() *PredictionModel {
	return &PredictionModel{
		samples: make([]NetworkSample, 0),
	}
}

// RecordEvent ثبت رویداد شبکه
func (nsm *NetworkStabilityManager) RecordEvent(apiType APIType, responseTime time.Duration, success bool, err error) {
	nsm.mu.Lock()
	defer nsm.mu.Unlock()

	event := NetworkEvent{
		Timestamp:    time.Now(),
		APIType:      apiType,
		ResponseTime: responseTime,
		Success:      success,
		Error:        err,
	}

	nsm.history = append(nsm.history, event)

	// نگهداری فقط 1000 رویداد آخر
	if len(nsm.history) > 1000 {
		nsm.history = nsm.history[1:]
	}

	// افزودن نمونه به مدل پیش‌بینی
	now := time.Now()
	hour := float64(now.Hour()) + float64(now.Minute())/60.0
	dayOfWeek := int(now.Weekday())

	sample := NetworkSample{
		TimeOfDay:    hour,
		DayOfWeek:    dayOfWeek,
		ResponseTime: responseTime,
		Success:      success,
	}

	nsm.predictionModel.AddSample(sample)
}

// AddSample افزودن نمونه به مدل پیش‌بینی
func (pm *PredictionModel) AddSample(sample NetworkSample) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.samples = append(pm.samples, sample)

	// نگهداری فقط 10000 نمونه آخر
	if len(pm.samples) > 10000 {
		pm.samples = pm.samples[1:]
	}
}

// PredictPerformance پیش‌بینی عملکرد شبکه
func (pm *PredictionModel) PredictPerformance() (expectedResponseTime time.Duration, successProbability float64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if len(pm.samples) < 10 {
		return 0, 0.5 // مقدار پیش‌فرض اگر نمونه کافی نباشد
	}

	now := time.Now()
	hour := float64(now.Hour()) + float64(now.Minute())/60.0
	dayOfWeek := int(now.Weekday())

	// پیدا کردن نمونه‌های مشابه
	var similarSamples []NetworkSample
	for _, sample := range pm.samples {
		hourDiff := sample.TimeOfDay - hour
		if hourDiff < 0 {
			hourDiff = -hourDiff
		}

		dayDiff := sample.DayOfWeek - dayOfWeek
		if dayDiff < 0 {
			dayDiff = -dayDiff
		}

		if hourDiff <= 2.0 && dayDiff <= 1 {
			similarSamples = append(similarSamples, sample)
		}
	}

	if len(similarSamples) < 5 {
		similarSamples = pm.samples
	}

	var totalResponseTime time.Duration
	successCount := 0

	for _, sample := range similarSamples {
		totalResponseTime += sample.ResponseTime
		if sample.Success {
			successCount++
		}
	}

	// *** این قسمت اصلاح شده است ***
	// به جای := از = استفاده می‌کنیم چون متغیرها در امضای تابع تعریف شده‌اند
	expectedResponseTime = totalResponseTime / time.Duration(len(similarSamples))
	successProbability = float64(successCount) / float64(len(similarSamples))

	return expectedResponseTime, successProbability
}

// CalculateBackoffDelay محاسبه تأخیر برای تلاش مجدد با الگوریتم نمایی و Jitter هوشمند
func (nsm *NetworkStabilityManager) CalculateBackoffDelay(retryCount int) time.Duration {
	// پیش‌بینی عملکرد شبکه
	_, successProbability := nsm.predictionModel.PredictPerformance()

	// اگر احتمال موفقیت کم است، تأخیر بیشتری اعمال کن
	if successProbability < 0.5 {
		retryCount += 2
	}

	// الگوریتم نمایی: baseDelay * (2 ^ retryCount)
	baseDelay := 1 * time.Second
	maxDelay := 30 * time.Second

	exponentialDelay := time.Duration(1<<uint(retryCount)) * baseDelay
	if exponentialDelay > maxDelay {
		exponentialDelay = maxDelay
	}

	// Jitter تا 25% از تأخیر اصلی
	jitter := time.Duration(rand.Int63n(int64(exponentialDelay / 4)))
	return exponentialDelay + jitter
}

// IsRetryableError بررسی اینکه آیا یک خطا قابل تلاش مجدد است یا خیر
func (nsm *NetworkStabilityManager) IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	if nsm.client.ignoreTimeout && isTimeoutError(err) {
		return true
	}

	if isNetworkError(err) {
		return true
	}

	return false
}

// isTimeoutError بررسی اینکه آیا خطا از نوع timeout است
func isTimeoutError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "deadline exceeded") ||
		strings.Contains(errStr, "context deadline exceeded")
}

// isNetworkError بررسی اینکه آیا خطا از نوع شبکه است
func isNetworkError(err error) bool {
	if netErr, ok := err.(net.Error); ok {
		return netErr.Temporary() || netErr.Timeout()
	}
	errStr := err.Error()
	return strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "no such host")
}

// GetNetworkStatistics دریافت آمار شبکه
func (nsm *NetworkStabilityManager) GetNetworkStatistics(apiType APIType) (avgResponseTime time.Duration, successRate float64, requestCount int) {
	nsm.mu.Lock()
	defer nsm.mu.Unlock()

	var totalResponseTime time.Duration
	successCount := 0
	requestCount = 0

	for _, event := range nsm.history {
		if event.APIType == apiType {
			totalResponseTime += event.ResponseTime
			requestCount++
			if event.Success {
				successCount++
			}
		}
	}

	if requestCount == 0 {
		return 0, 0, 0
	}

	avgResponseTime = totalResponseTime / time.Duration(requestCount)
	successRate = float64(successCount) / float64(requestCount)

	return avgResponseTime, successRate, requestCount
}
