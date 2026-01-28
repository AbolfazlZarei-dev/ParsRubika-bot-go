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

	// میدلور برای نمایش پیام‌های دریافتی
	client.Middleware(func(ctx context.Context, update *ParsRubika.Update, next ParsRubika.HandlerFunc) error {
		if update.NewMessage != nil {
			log.Printf("📩 پیام دریافت شد: %s", update.NewMessage.Text)
		}
		return next(ctx, update)
	})

	// ===============================
	// دستور /start: نمایش کیبورد
	// ===============================
	client.OnCommand("start", func(ctx context.Context, update *ParsRubika.Update) error {
		log.Println("دستور start دریافت شد. نمایش کیبورد...")

		// ساخت دکمه‌ها
		btn1 := &ParsRubika.KeyboardButton{
			ID:         "k_red",
			Type:       "Simple",
			ButtonText: "قرمز 🔴",
		}

		btn2 := &ParsRubika.KeyboardButton{
			ID:         "k_blue",
			Type:       "Simple",
			ButtonText: "آبی 🔵",
		}

		btn3 := &ParsRubika.KeyboardButton{
			ID:         "k_green",
			Type:       "Simple",
			ButtonText: "سبز 🟢",
		}

		// چیدن دکمه‌ها در ردیف‌ها
		row1 := ParsRubika.ReplyKeyboardRow{Buttons: []ParsRubika.KeyboardButton{*btn1, *btn2}}
		row2 := ParsRubika.ReplyKeyboardRow{Buttons: []ParsRubika.KeyboardButton{*btn3}}

		// ساخت مارک‌آپ
		keypad := &ParsRubika.ReplyKeyboardMarkup{
			Rows:           []ParsRubika.ReplyKeyboardRow{row1, row2},
			ResizeKeyboard: true, // سایز کیبورد را کوچک نگه می‌دارد
		}

		// ارسال پیام همراه با کیبورد در یک درخواست
		// نکته: اگر ChatKeypadType را به "New" تنظیم نکنید، کیبورد نمایش داده نمی‌شود
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID:         update.ChatID,
			Text:           "سلام! این یک کیبورد ساده پایین صفحه است.\nروی دکمه‌ها کلیک کنید یا /hide را بزنید.",
			ChatKeypad:     keypad,
			ChatKeypadType: "New", // <-- حیاتی است: ایجاد کیبورد جدید
		})

		return err
	})

	// ===============================
	// هندلر کلیک دکمه‌ها (بر اساس متن)
	// ===============================
	// دکمه‌های Reply Keyboard متن ارسال می‌کنند، پس باید متن را چک کنیم
	client.OnMessageUpdates(func(ctx context.Context, update *ParsRubika.Update) error {
		if update.NewMessage == nil {
			return nil
		}

		text := update.NewMessage.Text

		switch text {
		case "قرمز 🔴":
			log.Println("دکمه قرمز کلیک شد")
			_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   "شما روی دکمه **قرمز** کلیک کردید!",
			})
			return err

		case "آبی 🔵":
			log.Println("دکمه آبی کلیک شد")
			_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   "شما روی دکمه **آبی** کلیک کردید!",
			})
			return err

		case "سبز 🟢":
			log.Println("دکمه سبز کلیک شد")
			_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   "شما روی دکمه **سبز** کلیک کردید!",
			})
			return err
		}

		return nil
	})

	// ===============================
	// دستور /hide: مخفی کردن کیبورد
	// ===============================
	client.OnCommand("hide", func(ctx context.Context, update *ParsRubika.Update) error {
		log.Println("در حال حذف کیبورد...")

		err := client.EditChatKeypad(ctx, &ParsRubika.EditChatKeypadRequest{
			ChatID:         update.ChatID,
			ChatKeypadType: "Remove", // حذف کیبورد
		})

		if err == nil {
			_, _ = client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   "کیبورد مخفی شد.",
			})
		}
		return err
	})

	ctx := context.Background()
	log.Println("ربات کیبورد ساده شروع شد. دستور /start را بزنید.")

	if err := client.Run(ctx, ParsRubika.PollingOptions{
		Limit:        100,
		PollInterval: 2 * time.Second,
	}); err != nil {
		log.Fatalf("❌ خطا: %v", err)
	}
}
