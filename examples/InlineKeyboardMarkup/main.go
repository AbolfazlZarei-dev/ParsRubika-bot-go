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
	// ایجاد پوشه افست‌ها برای ذخیره آخرین وضعیت پولینگ
	if err := os.MkdirAll("offsets", 0755); err != nil {
		log.Fatalf("خطا در ایجاد پوشه offsets: %v", err)
	}
	log.Println("✅ بررسی دسترسی نوشتن با موفقیت انجام شد.")

	log.Println("در حال ساخت کلاینت ربات...")
	// استفاده از حالت سوییچر برای پایداری بیشتر (مانند فایل اول)
	client := ParsRubika.NewClient(token, ParsRubika.WithConnectionMode(ParsRubika.SwitcherMode))

	// میدلور برای نمایش لاگ‌های دریافتی (برای دیباگ بهتر)
	client.Middleware(func(ctx context.Context, update *ParsRubika.Update, next ParsRubika.HandlerFunc) error {
		if update.NewMessage != nil {
			log.Printf("🔔 پیام جدید از کاربر %s در چت %s: %s", update.NewMessage.SenderID, update.ChatID, update.NewMessage.Text)
		} else if update.Type == ParsRubika.CallbackQuery {
			log.Printf("🔔 کلیک دکمه دریافت شد در چت %s", update.ChatID)
		}
		return next(ctx, update)
	})

	// ===============================
	// هندلر دستور /test (از فایل دوم)
	// ===============================
	// این دستور دکمه‌های اینلاین را ارسال می‌کند
	client.OnCommand("test", func(ctx context.Context, update *ParsRubika.Update) error {
		log.Printf("دستور /test دریافت شد. ارسال دکمه‌ها به چت %s...", update.ChatID)

		// ساخت دکمه اول
		btnOption1 := &ParsRubika.InlineKeyboardButton{
			ID:         "opt_1",
			Type:       "Simple",
			ButtonText: "گزینه اول 🔴",
		}

		// ساخت دکمه دوم
		btnOption2 := &ParsRubika.InlineKeyboardButton{
			ID:         "opt_2",
			Type:       "Simple",
			ButtonText: "گزینه دوم 🟢",
		}

		// ساخت دکمه سوم (لینک)
		btnLink := &ParsRubika.InlineKeyboardButton{
			ID:         "link_btn",
			Type:       "Link",
			ButtonText: "گیت‌هاب من",
			URL:        "https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go",
		}

		// چیدن دکمه‌ها در ردیف‌ها
		row1 := ParsRubika.InlineKeyboardRow{
			Buttons: []ParsRubika.InlineKeyboardButton{*btnOption1, *btnOption2},
		}

		row2 := ParsRubika.InlineKeyboardRow{
			Buttons: []ParsRubika.InlineKeyboardButton{*btnLink},
		}

		// ساخت کیبورد کامل
		keypad := &ParsRubika.InlineKeyboardMarkup{
			Rows: []ParsRubika.InlineKeyboardRow{row1, row2},
		}

		// ارسال پیام
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID:       update.ChatID,
			Text:         "لطفاً یکی از گزینه‌ها را انتخاب کنید (تست دکمه‌ها):",
			InlineKeypad: keypad,
		})

		return err
	})

	// ===============================
	// هندلر کلیک دکمه (Callback)
	// ===============================
	client.OnCallbackQuery(func(ctx context.Context, update *ParsRubika.Update) error {
		if update.NewMessage != nil && update.NewMessage.AuxData != nil && update.NewMessage.AuxData.ButtonID != nil {
			clickedBtnID := *update.NewMessage.AuxData.ButtonID
			originalText := update.NewMessage.Text

			var responseText string
			switch clickedBtnID {
			case "opt_1":
				responseText = "شما روی **گزینه اول** کلیک کردید! ✅"
			case "opt_2":
				responseText = "شما روی **گزینه دوم** کلیک کردید! ✅"
			case "link_btn":
				responseText = "شما روی لینک کلیک کردید."
			default:
				responseText = "دکمه ناشناخته!"
			}

			log.Printf("دکمه '%s' فشرده شد. در حال ویرایش پیام...", clickedBtnID)

			// ویرایش پیام برای نمایش نتیجه
			err := client.EditMessageText(ctx, &ParsRubika.EditMessageTextRequest{
				ChatID:    update.ChatID,
				MessageID: update.NewMessage.MessageID,
				Text:      fmt.Sprintf("%s\n\n%s", originalText, responseText),
			})

			return err
		}
		return fmt.Errorf("داده‌های دکمه یافت نشد")
	})

	// ===============================
	// شروع بات با پشتیبانی از آفست
	// ===============================
	ctx := context.Background()
	log.Println("ربات در حال راه‌اندازی و شروع پولینگ با پشتیبانی از آفست...")
	log.Println("لطفاً دستور /test را ارسال کنید.")

	if err := client.Run(ctx, ParsRubika.PollingOptions{
		Limit:        100,
		PollInterval: 2 * time.Second,
	}); err != nil {
		log.Fatalf("❌ خطا در اجرای ربات: %v", err)
	}
}
