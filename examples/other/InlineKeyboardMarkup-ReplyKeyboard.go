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
	// استفاده از حالت سوییچر برای پایداری بیشتر
	client := ParsRubika.NewClient(token, ParsRubika.WithConnectionMode(ParsRubika.SwitcherMode))

	// میدلور برای نمایش لاگ‌های دریافتی (برای دیباگ بهتر)
	client.Middleware(func(ctx context.Context, update *ParsRubika.Update, next ParsRubika.HandlerFunc) error {
		if update.NewMessage != nil {
			log.Printf("🔔 پیام جدید از کاربر %s در چت %s: %s", update.NewMessage.SenderID, update.ChatID, update.NewMessage.Text)
		} else if update.Type == ParsRubika.CallbackQuery {
			log.Printf("🔔 کلیک دکمه اینلاین دریافت شد در چت %s", update.ChatID)
		}
		return next(ctx, update)
	})

	// ===============================
	// هندلر دستور /start (اصلاح شده)
	// ===============================
	client.OnCommand("start", func(ctx context.Context, update *ParsRubika.Update) error {
		log.Printf("دستور /start دریافت شد. ارسال پیام خوشامدگویی و کیبورد شیشه‌ای به چت %s...", update.ChatID)

		// ساخت دکمه‌ها برای کیبورد شیشه‌ای که با استارت ارسال می‌شود
		btnHelp := &ParsRubika.KeyboardButton{ID: "start_help", Type: "Simple", ButtonText: "❓ راهنما"}
		btnAbout := &ParsRubika.KeyboardButton{ID: "start_about", Type: "Simple", ButtonText: "ℹ️ درباره ربات"}

		// چیدن دکمه‌ها در یک ردیف
		row := ParsRubika.ReplyKeyboardRow{Buttons: []ParsRubika.KeyboardButton{*btnHelp, *btnAbout}}

		// ساخت کیبورد کامل
		startKeypad := &ParsRubika.ReplyKeyboardMarkup{
			Rows:            []ParsRubika.ReplyKeyboardRow{row},
			ResizeKeyboard:  true,  // کوچک کردن کیبورد برای ظاهر بهتر
			OneTimeKeyboard: false, // کیبورد پس از استفاده مخفی نشود
		}

		// ارسال پیام خوشامدگویی به همراه کیبورد شیشه‌ای
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID:         update.ChatID,
			Text:           "سلام! به ربات نمونه ParsRubika خوش آمدید.\n\nاز دکمه‌های زیر برای شروع استفاده کنید:",
			ChatKeypad:     startKeypad,
			ChatKeypadType: ParsRubika.NewKeypad, // این فیلد کیبورد را نمایش می‌دهد
		})
		return err
	})

	// ===============================
	// هندلر دستور /test (دکمه‌های اینلاین)
	// ===============================
	client.OnCommand("test", func(ctx context.Context, update *ParsRubika.Update) error {
		log.Printf("دستور /test دریافت شد. ارسال دکمه‌های اینلاین به چت %s...", update.ChatID)

		btnOption1 := &ParsRubika.InlineKeyboardButton{ID: "opt_1", Type: "Simple", ButtonText: "گزینه اول 🔴"}
		btnOption2 := &ParsRubika.InlineKeyboardButton{ID: "opt_2", Type: "Simple", ButtonText: "گزینه دوم 🟢"}
		btnLink := &ParsRubika.InlineKeyboardButton{ID: "link_btn", Type: "Link", ButtonText: "گیت‌هاب من", URL: "https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go"}

		row1 := ParsRubika.InlineKeyboardRow{Buttons: []ParsRubika.InlineKeyboardButton{*btnOption1, *btnOption2}}
		row2 := ParsRubika.InlineKeyboardRow{Buttons: []ParsRubika.InlineKeyboardButton{*btnLink}}

		keypad := &ParsRubika.InlineKeyboardMarkup{Rows: []ParsRubika.InlineKeyboardRow{row1, row2}}

		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID:       update.ChatID,
			Text:         "این دکمه‌های اینلاین هستند (به پیام متصل‌اند):",
			InlineKeypad: keypad,
		})
		return err
	})

	// ===============================
	// هندلر کلیک دکمه اینلاین (Callback)
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

			log.Printf("دکمه اینلاین '%s' فشرده شد. در حال ویرایش پیام...", clickedBtnID)

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
	// هندلر یکپارچه برای تمام دکمه‌های شیشه‌ای
	// ===============================
	// این هندلر پیام‌های ارسالی از تمام دکمه‌های شیشه‌ای را مدیریت می‌کند
	client.AddHandler(func(ctx context.Context, update *ParsRubika.Update) error {
		if update.NewMessage == nil || update.NewMessage.Text == "" {
			return nil
		}

		var responseText string
		switch update.NewMessage.Text {
		// دکمه‌هایی که با دستور /start ارسال می‌شوند
		case "❓ راهنما":
			responseText = "این ربات نمونه قابلیت‌های زیر را دارد:\n\n/test : ارسال دکمه‌های اینلاین\n/glass : ارسال دکمه شیشه‌ای دیگر\n/removeglass : حذف کیبورد فعلی"
		case "ℹ️ درباره ربات":
			responseText = "این یک ربات نمونه برای کتابخانه ParsRubika نسخه 2 است.\nسازنده: ابوالفضل زارعی\nلینک گیت‌هاب: https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go"

		// دکمه‌هایی که با دستور /glass ارسال می‌شوند
		case "گزینه اول شیشه‌ای 🔴":
			responseText = "شما روی گزینه اول شیشه‌ای کلیک کردید! 🎉"
		case "گزینه دوم شیشه‌ای 🟢":
			responseText = "شما روی گزینه دوم شیشه‌ای کلیک کردید! 🎉"

		default:
			return nil // اگر پیام مربوط به دکمه‌های ما نبود، از هندلر رد شو
		}

		log.Printf("دکمه شیشه‌ای فشرده شد: %s در چت %s", update.NewMessage.Text, update.ChatID)
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   responseText,
		})
		return err
	}, nil)

	// ===============================
	// هندلر دستور /glass (برای ارسال یک دکمه شیشه‌ای دیگر)
	// ===============================
	client.OnCommand("glass", func(ctx context.Context, update *ParsRubika.Update) error {
		log.Printf("دستور /glass دریافت شد. ارسال کیبورد شیشه‌ای به چت %s...", update.ChatID)

		btn1 := &ParsRubika.KeyboardButton{ID: "glass_opt_1", Type: "Simple", ButtonText: "گزینه اول شیشه‌ای 🔴"}
		btn2 := &ParsRubika.KeyboardButton{ID: "glass_opt_2", Type: "Simple", ButtonText: "گزینه دوم شیشه‌ای 🟢"}
		btn3 := &ParsRubika.KeyboardButton{ID: "glass_loc", Type: "AskMyLocation", ButtonText: "ارسال موقعیت من 📍"}
		btn4 := &ParsRubika.KeyboardButton{ID: "glass_contact", Type: "AskMyPhoneNumber", ButtonText: "ارسال شماره تماس 📞"}

		row1 := ParsRubika.ReplyKeyboardRow{Buttons: []ParsRubika.KeyboardButton{*btn1, *btn2}}
		row2 := ParsRubika.ReplyKeyboardRow{Buttons: []ParsRubika.KeyboardButton{*btn3, *btn4}}

		chatKeypad := &ParsRubika.ReplyKeyboardMarkup{
			Rows:            []ParsRubika.ReplyKeyboardRow{row1, row2},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}

		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID:         update.ChatID,
			Text:           "این یک کیبورد شیشه‌ای دیگر است. روی دکمه‌ها کلیک کنید:",
			ChatKeypad:     chatKeypad,
			ChatKeypadType: ParsRubika.NewKeypad,
		})
		return err
	})

	// ===============================
	// هندلر برای دریافت موقعیت و تماس (این‌ها به طور خودکار کار می‌کنند)
	// ===============================
	client.OnLocation(func(ctx context.Context, update *ParsRubika.Update) error {
		log.Printf("موقعیت کاربر %s دریافت شد: %s, %s", update.NewMessage.SenderID, update.NewMessage.Location.Latitude, update.NewMessage.Location.Longitude)
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   fmt.Sprintf("موقعیت شما با موفقیت دریافت شد!\nطول جغرافیایی: %s\nعرض جغرافیایی: %s", update.NewMessage.Location.Longitude, update.NewMessage.Location.Latitude),
		})
		return err
	})

	client.OnContact(func(ctx context.Context, update *ParsRubika.Update) error {
		contact := update.NewMessage.ContactMessage
		log.Printf("شماره تماس کاربر %s دریافت شد: %s", update.NewMessage.SenderID, contact.PhoneNumber)
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   fmt.Sprintf("شماره تماس شما با موفقیت دریافت شد:\n%s %s", contact.FirstName, contact.PhoneNumber),
		})
		return err
	})

	// ===============================
	// هندلر دستور /removeglass (برای حذف دکمه شیشه‌ای)
	// ===============================
	client.OnCommand("removeglass", func(ctx context.Context, update *ParsRubika.Update) error {
		log.Printf("دستور /removeglass دریافت شد. در حال حذف کیبورد شیشه‌ای از چت %s...", update.ChatID)

		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID:         update.ChatID,
			Text:           "کیبورد شیشه‌ای با موفقیت حذف شد.",
			ChatKeypadType: ParsRubika.RemoveKeypad, // این فیلد کیبورد را از چت حذف می‌کند
		})
		return err
	})

	// ===============================
	// شروع بات با پشتیبانی از آفست
	// ===============================
	ctx := context.Background()
	log.Println("ربات در حال راه‌اندازی و شروع پولینگ با پشتیبانی از آفست...")
	log.Println("لطفاً دستور /start را ارسال کنید.")

	if err := client.Run(ctx, ParsRubika.PollingOptions{
		Limit:        100,
		PollInterval: 2 * time.Second,
	}); err != nil {
		log.Fatalf("❌ خطا در اجرای ربات: %v", err)
	}
}
