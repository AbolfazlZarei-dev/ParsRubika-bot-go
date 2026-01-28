package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v2"
)

func main() {
	token := "YOUR_BOT_TOKEN_HERE"

	// ایجاد پوشه آفست‌ها
	if err := os.MkdirAll("offsets", 0755); err != nil {
		log.Fatalf("خطا در ایجاد پوشه offsets: %v", err)
	}

	log.Println("در حال ساخت کلاینت ربات...")
	client := ParsRubika.NewClient(token)

	// میدلور برای لاگ کردن پیام‌ها
	client.Middleware(func(ctx context.Context, update *ParsRubika.Update, next ParsRubika.HandlerFunc) error {
		if update.NewMessage != nil {
			log.Printf("📩 پیام: %s", update.NewMessage.Text)
		}
		return next(ctx, update)
	})

	// ===============================
	// هندلر دستور /start
	// ===============================
	client.OnCommand("start", func(ctx context.Context, update *ParsRubika.Update) error {
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   "سلام! به ربات ساده خوش آمدید 🌹\n\nبرای مشاهده راهنما دستور /help را ارسال کنید.",
		})
		return err
	})

	// ===============================
	// هندلر دستور /help
	// ===============================
	client.OnCommand("help", func(ctx context.Context, update *ParsRubika.Update) error {
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   "📚 راهنمای استفاده از ربات:\n\n/start: شروع مجدد و پیام خوش‌آمدگویی\n/help: مشاهده همین راهنما",
		})
		return err
	})

	ctx := context.Background()
	log.Println("ربات شروع شد. دستورات /start و /help را تست کنید.")

	if err := client.Run(ctx, ParsRubika.PollingOptions{
		Limit:        100,
		PollInterval: 2 * time.Second,
	}); err != nil {
		log.Fatalf("❌ خطا: %v", err)
	}
}
