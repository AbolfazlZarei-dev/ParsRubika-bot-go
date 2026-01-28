package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v2"
)

func main() {
	// توکن ربات خود را اینجا قرار دهید
	token := "YOUR_BOT_TOKEN_HERE"

	log.Println("در حال بررسی دسترسی نوشتن در پوشه 'offsets'...")
	if err := os.MkdirAll("offsets", 0755); err != nil {
		log.Fatalf("خطا در ایجاد پوشه offsets: %v", err)
	}
	log.Println("✅ بررسی دسترسی نوشتن با موفقیت انجام شد.")

	log.Println("در حال ساخت کلاینت ربات...")
	// استفاده از حالت سوییچر برای پایداری بیشتر
	client := ParsRubika.NewClient(token, ParsRubika.WithConnectionMode(ParsRubika.SwitcherMode))

	// میدلور برای نمایش لاگ‌های دریافتی
	client.Middleware(func(ctx context.Context, update *ParsRubika.Update, next ParsRubika.HandlerFunc) error {
		if update.NewMessage != nil {
			log.Printf("🔔 پیام جدید از کاربر %s در چت %s: %s", update.NewMessage.SenderID, update.ChatID, update.NewMessage.Text)
		}
		return next(ctx, update)
	})

	// هندلر برای دستور /start
	client.OnCommand("start", func(ctx context.Context, update *ParsRubika.Update) error {
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   "سلام! به ربات نمونه ParsRubika خوش آمدید.\n\nبرای دیدن تمام فرمت‌های پشتیبانی شده، دستور `/formats` را ارسال کنید.",
		})
		return err
	})

	// هندلر برای دستور /formats
	client.OnCommand("formats", func(ctx context.Context, update *ParsRubika.Update) error {
		senderID := update.NewMessage.SenderID
		log.Printf("دستور /formats دریافت شد از طرف %s. در حال ارسال پیام نمونه...", senderID)
		return sendAllFormatsExample(ctx, client, update.ChatID, senderID)
	})

	ctx := context.Background()
	log.Println("ربات در حال راه‌اندازی و شروع پولینگ...")

	if err := client.Run(ctx, ParsRubika.PollingOptions{
		Limit:        100,
		PollInterval: 2 * time.Second,
	}); err != nil {
		log.Fatalf("❌ خطا در اجرای ربات: %v", err)
	}
}

// sendAllFormatsExample پیامی حاوی تمام فرمت‌های پشتیبانی شده را ارسال می‌کند
func sendAllFormatsExample(ctx context.Context, client *ParsRubika.BotClient, chatID, senderID string) error {
	// ساخت متن با تمام فرمت‌های ممکن
	formattedText := fmt.Sprintf(`این یک پیام تست برای متادیتاهای ParsRubika است:

**این متن برجسته (Bold) است.**

__این متن زیرخط‌دار (Underline) است.__

~~این متن خط‌خورده (Strikethrough) است.~~

`+"`"+`این متن تک‌خطی (Inline Code) است.`+"`"+`

||این متن اسپویلر (Spoiler) است.||

××این یک متن نقل قول (Quote) است.××

[این یک لینک به گیت‌هاب ParsRubika](https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go)

(%s)[شما] این یک منشن است.
`, senderID)

	// ارسال پیام با فرمت مارک‌داون
	_, err := client.SendMessageWithMarkdown(ctx, chatID, formattedText)
	if err != nil {
		return fmt.Errorf("خطا در ارسال پیام فرمت‌دهی شده: %w", err)
	}

	return nil
}
