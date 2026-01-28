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

	// میدلور
	client.Middleware(func(ctx context.Context, update *ParsRubika.Update, next ParsRubika.HandlerFunc) error {
		if update.NewMessage != nil {
			log.Printf("📩 پیام: %s", update.NewMessage.Text)
		}
		return next(ctx, update)
	})

	// ===============================
	// هندلر دستور /start: تنظیم منوی بات
	// ===============================
	client.OnCommand("start", func(ctx context.Context, update *ParsRubika.Update) error {
		log.Println("در حال تنظیم دستورات بات (Bot Commands)...")

		// تعریف لیست دستورات
		commands := []ParsRubika.BotCommand{
			{
				Command:     "start",
				Description: "شروع مجدد ربات 🚀",
			},
			{
				Command:     "info",
				Description: "دریافت اطلاعات بات ℹ️",
			},
			{
				Command:     "help",
				Description: "راهنما و پشتیبانی ❓",
			},
		}

		// ارسال دستورات به سرور برای تنظیم منوی بات
		err := client.SetCommands(ctx, &ParsRubika.SetCommandsRequest{
			BotCommands: commands,
		})
		if err != nil {
			return err
		}

		log.Println("دستورات با موفقیت تنظیم شدند.")

		// ارسال یک پیام ساده بدون هیچ دکمه‌ای
		_, err = client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   "منوی دستورات بات تنظیم شد.\n\n(این پیام دارای هیچ دکمه‌ای نیست)",
		})
		return err
	})

	// ===============================
	// هندلرهای تستی (اختیاری)
	// ===============================
	// فقط برای اطمینان از اینکه دستورات کار می‌کنند
	client.OnCommand("info", func(ctx context.Context, update *ParsRubika.Update) error {
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   "این اطلاعات بات است.",
		})
		return err
	})

	ctx := context.Background()
	log.Println("ربات شروع شد. دستور /start را بزنید.")

	if err := client.Run(ctx, ParsRubika.PollingOptions{
		Limit:        100,
		PollInterval: 2 * time.Second,
	}); err != nil {
		log.Fatalf("❌ خطا: %v", err)
	}
}
