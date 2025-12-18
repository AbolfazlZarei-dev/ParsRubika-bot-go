package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go

import (
	"context"
	"fmt"
	"log"
	"time"
)

// PollingOptions تنظیمات فرآیند Polling
type PollingOptions struct {
	Handler           HandlerFunc
	RetryTimeout      time.Duration // این فیلد دیگر به صورت مستقیم استفاده نمی‌شود و توسط NetworkStabilityManager مدیریت می‌شود
	Limit             int
	AllowEmptyUpdates bool
	PollInterval      time.Duration
	Timeout           time.Duration
}

// StartPolling شروع دریافت آپدیت‌ها با پولینگ (هوشمند شده)
func (c *BotClient) StartPolling(ctx context.Context, opts PollingOptions) error {
	if opts.Handler == nil {
		opts.Handler = c.ProcessUpdate
	}

	if opts.Limit == 0 {
		opts.Limit = 100
	}
	if opts.PollInterval == 0 {
		opts.PollInterval = 2 * time.Second
	}
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	if err := c.getBotID(ctx); err != nil {
		return fmt.Errorf("could not start polling without bot ID: %w", err)
	}

	var offset *string

	log.Println("در حال پاک‌سازی صف آپدیت‌های قدیمی...")
	if err := c.clearUpdateQueue(ctx, &offset); err != nil {
		log.Printf("اخطار در پاک‌سازی صف: %v", err)
	}

	log.Printf("شروع پولینگ هوشمند با API: %s", c.apiManager.GetCurrentAPI())

	pollingTicker := time.NewTicker(opts.PollInterval)
	defer pollingTicker.Stop()

	var retryCount int
	for {
		select {
		case <-ctx.Done():
			log.Println("پولینگ توسط context متوقف شد")
			return ctx.Err()

		case <-c.stopChan:
			log.Println("پولینگ توسط بات متوقف شد")
			return nil

		case <-pollingTicker.C:
			log.Printf("🔄 در حال درخواست آپدیت با آفست: %v", offset)
			updates, err := c.GetUpdates(ctx, offset, opts.Limit)
			if err != nil {
				if c.networkStabilityManager.IsRetryableError(err) {
					delay := c.networkStabilityManager.CalculateBackoffDelay(retryCount)
					log.Printf("خطا در دریافت آپدیت‌ها: %v. تلاش مجدد در %v", err, delay)
					retryCount++
					time.Sleep(delay)
					continue
				} else {
					log.Printf("خطای غیرقابل بازیابی در پولینگ: %v", err)
					return err
				}
			}

			retryCount = 0

			// --- این لاگ‌های جدید کلیدی هستند ---
			log.Printf("📨 پاسخ از سرور دریافت شد. تعداد آپدیت‌ها: %d, آفست بعدی: '%s'", len(updates.Updates), updates.NextOffsetID)
			// --- پایان لاگ‌های جدید ---

			if len(updates.Updates) > 0 {
				log.Printf("📥 شروع پردازش %d آپدیت جدید...", len(updates.Updates))

				for _, update := range updates.Updates {
					log.Printf("🔄 در حال پردازش آپدیت با نوع: %s از چت: %s", update.Type, update.ChatID)
					if err := opts.Handler(ctx, &update); err != nil {
						log.Printf("خطا در پردازش آپدیت: %v", err)
					}
				}

				if updates.NextOffsetID != "" {
					offset = &updates.NextOffsetID
					log.Printf("✅ آفست به‌روزرسانی شد به: %s", *offset)
				}
			} else {
				log.Println("📭 هیچ آپدیت جدیدی در این تیک وجود نداشت.")
			}
		}
	}
}

// clearUpdateQueue پاک‌سازی صف آپدیت‌های قدیمی (بهینه شده)
func (c *BotClient) clearUpdateQueue(ctx context.Context, offset **string) error {
	for {
		discardUpdates, err := c.GetUpdates(ctx, *offset, 100)
		if err != nil {
			return fmt.Errorf("خطا در پاک‌سازی صف: %w", err)
		}

		if len(discardUpdates.Updates) == 0 {
			log.Println("صف آپدیت‌ها با موفقیت پاک شد")
			break
		}

		log.Printf("تعداد %d آپدیت قدیمی دور ریخته شد", len(discardUpdates.Updates))

		if discardUpdates.NextOffsetID != "" {
			*offset = &discardUpdates.NextOffsetID
		} else {
			break
		}
	}
	return nil
}
