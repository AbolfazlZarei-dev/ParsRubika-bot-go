package ParsRubika

// نسخه: 2.5.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go

import (
	"context"
	"fmt"
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

// StartPolling شروع دریافت آپدیت‌ها با پولینگ (بدون لاگ و بهینه)
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

	if err := c.init(ctx); err != nil {
		return fmt.Errorf("could not start polling without bot ID: %w", err)
	}

	// دریافت آفست ذخیره شده
	var offset *string
	savedOffset := c.GetOffset()
	if savedOffset != "" {
		offset = &savedOffset
	} else {
		// فقط در اولین اجرا، صف را پاک کرده و اولین آفست را دریافت می‌کنیم
		c.clearUpdateQueue(ctx, &offset)
	}

	pollingTicker := time.NewTicker(opts.PollInterval)
	defer pollingTicker.Stop()

	var retryCount int
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.stopChan:
			return nil
		case <-pollingTicker.C:
			updates, err := c.GetUpdates(ctx, offset, opts.Limit)
			if err != nil {
				if c.networkStabilityManager.IsRetryableError(err) {
					delay := c.networkStabilityManager.CalculateBackoffDelay(retryCount)
					retryCount++
					time.Sleep(delay)
					continue
				} else {
					return err
				}
			}

			retryCount = 0

			if len(updates.Updates) > 0 {
				for _, update := range updates.Updates {
					if err := opts.Handler(ctx, &update); err != nil {
					}
				}

				if updates.NextOffsetID != "" {
					offset = &updates.NextOffsetID
					// ذخیره سریع آفست جدید برای هر دسته از آپدیت‌ها
					if err := c.SetOffset(updates.NextOffsetID); err != nil {
					}
				}
			}
		}
	}
}

// clearUpdateQueue پاک‌سازی صف آپدیت‌های قدیمی (بدون لاگ)
func (c *BotClient) clearUpdateQueue(ctx context.Context, offset **string) error {
	for {
		discardUpdates, err := c.GetUpdates(ctx, *offset, 100)
		if err != nil {
			return fmt.Errorf("خطا در پاک‌سازی صف: %w", err)
		}

		if len(discardUpdates.Updates) == 0 {
			break
		}

		if discardUpdates.NextOffsetID != "" {
			*offset = &discardUpdates.NextOffsetID
			// ذخیره آخرین آفست در حین پاک‌سازی
			if err := c.SetOffset(discardUpdates.NextOffsetID); err != nil {
			}
		} else {
			break
		}
	}
	return nil
}
