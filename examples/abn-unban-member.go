package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v2"
)

// --- اطلاعات و تنظیمات اصلی ربات ---
const (
	// توکن ربات خود را اینجا قرار دهید
	BotToken = "token"

	// آیدی چتی که می‌خواهید پیام‌های آزمایشی و اطلاع‌رسانی به آن ارسال شود
	TargetChatID = "chat-id"
)

// --- متغیر ادمین ---
// این متغیر فقط یک بار پر می‌شود و دیگر تغییر نمی‌کند
var AdminID = ""

// کش ساده برای نگهداری ارتباط بین MessageID و UserID کاربر
var messageCache = struct {
	sync.RWMutex
	data map[string]string // Key: MessageID, Value: UserID
}{
	data: make(map[string]string),
}

// تابع کمکی برای دریافت ورودی از کاربر در ترمینال
func promptUser(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func main() {
	log.Println("🚀 راه‌اندازی ربات ParsRubika نسخه 2.0.0")

	// --- مرحله 1: انتخاب حالت اتصال API ---
	var connectionMode ParsRubika.ConnectionMode
	for {
		fmt.Println("\nلطفاً حالت اتصال API را انتخاب کنید:")
		fmt.Println("1. Bot API (فقط API اصلی روبیکا)")
		fmt.Println("2. Messenger API (فقط API جدید مسنجر)")
		fmt.Println("3. Smart Switcher (سوییچ هوشمند بین دو API)")

		choice := promptUser("گزینه مورد نظر (1-3): ")
		switch choice {
		case "1":
			connectionMode = ParsRubika.BotAPIMode
			log.Println("✅ حالت اتصال: Bot API انتخاب شد.")
			goto API_CHOSEN
		case "2":
			connectionMode = ParsRubika.MessengerAPIMode
			log.Println("✅ حالت اتصال: Messenger API انتخاب شد.")
			goto API_CHOSEN
		case "3":
			connectionMode = ParsRubika.SwitcherMode
			log.Println("✅ حالت اتصال: Smart Switcher انتخاب شد.")
			goto API_CHOSEN
		default:
			fmt.Println("❌ گزینه نامعتبر. لطفاً دوباره تلاش کنید.")
		}
	}
API_CHOSEN:

	// --- مرحله 2: انتخاب نحوه اجرای ربات ---
	var runningMode string
	for {
		fmt.Println("\nنحوه اجرای ربات را انتخاب کنید:")
		fmt.Println("1. Polling (دریافت آپدیت با پولینگ)")
		fmt.Println("2. Webhook (دریافت آپدیت از طریق وب‌هوک)")
		fmt.Println("3. Host-Reload + Polling (توسعه با بارگذاری مجدد خودکار)")

		choice := promptUser("گزینه مورد نظر (1-3): ")
		switch choice {
		case "1":
			runningMode = "polling"
			log.Println("✅ حالت اجرا: Polling انتخاب شد.")
			goto RUN_MODE_CHOSEN
		case "2":
			runningMode = "webhook"
			log.Println("✅ حالت اجرا: Webhook انتخاب شد.")
			goto RUN_MODE_CHOSEN
		case "3":
			runningMode = "host_reload"
			log.Println("✅ حالت اجرا: Host-Reload + Polling انتخاب شد.")
			goto RUN_MODE_CHOSEN
		default:
			fmt.Println("❌ گزینه نامعتبر. لطفاً دوباره تلاش کنید.")
		}
	}
RUN_MODE_CHOSEN:

	// --- مرحله 3: ساخت و پیکربندی کلاینت ربات ---
	client := ParsRubika.NewClient(BotToken,
		ParsRubika.WithConnectionMode(connectionMode),
		ParsRubika.WithHotReload(runningMode == "host_reload"),
		ParsRubika.WithFastConnection(),
		ParsRubika.WithNotificationOptions(ParsRubika.NotificationOptions{
			Enabled:     true,
			SendToOwner: true,
			ChatID:      TargetChatID,
			Message:     "🔔 اطلاعیه سوییچ API",
		}),
	)

	// --- مرحله 3.5: هندلر عمومی (فعال‌سازی یک‌بار مصرف و کشینگ) ---
	client.OnMessageUpdates(func(ctx context.Context, update *ParsRubika.Update) error {
		if update.NewMessage != nil {
			// --- منطق فعال‌سازی یک‌بار مصرف ---
			if update.NewMessage.Text == "فعال" {
				senderID := update.NewMessage.SenderID

				if AdminID == "" {
					// ادمین هنوز تعیین نشده است
					AdminID = senderID
					log.Printf("✅ ادمین اصلی تایید شد: %s", AdminID)

					_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
						ChatID: update.ChatID,
						Text:   fmt.Sprintf("✅ شما به عنوان **ادمین اصلی** ثبت شدید.\nآیدی شما: `%s`", AdminID),
					})
					return err
				} else {
					// ادمین قبلاً انتخاب شده است
					_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
						ChatID: update.ChatID,
						Text:   "⛔ ادغین اصلی قبلاً انتخاب شده است. شما نمی‌توانید ادمین شوید.",
					})
					return err
				}
			}

			// --- کشینگ پیام‌ها برای ریپلای ---
			messageCache.Lock()
			messageCache.data[update.NewMessage.MessageID] = update.NewMessage.SenderID
			messageCache.Unlock()
		}
		return nil
	})

	// --- مرحله 4: ثبت هندلرها ---

	// هندلر دستور /start
	client.OnCommand("start", func(ctx context.Context, update *ParsRubika.Update) error {
		_, err := client.GetUserInfo(ctx, update.NewMessage.SenderID)
		if err != nil {
			log.Printf("خطا در دریافت اطلاعات کاربر %s: %v", update.NewMessage.SenderID, err)
		}

		var statusMsg string
		if AdminID == "" {
			statusMsg = "\n\n⚠️ **ادمین هنوز انتخاب نشده است!**\nاگر صاحب ربات هستید، پیام **«فعال ربات من»** را ارسال کنید."
		} else {
			statusMsg = ""
		}

		welcomeText := fmt.Sprintf(
			"👋 به ربات خوش آمدید!\n\n"+
				"حالت اتصال فعلی: `%s`"+statusMsg,
			client.GetAPIManager().GetCurrentAPI(),
		)
		_, err = client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   welcomeText,
		})
		return err
	})

	// --- هندلر بن کردن (فقط برای ادمین اصلی) ---
	client.OnCommand("ban", func(ctx context.Context, update *ParsRubika.Update) error {
		senderID := update.NewMessage.SenderID

		// بررسی دسترسی: اگر ادمین خالی است یا فرستنده ادمین نیست
		if AdminID == "" || senderID != AdminID {
			client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   "⛔ شما اجازه دسترسی به این دستور را ندارید.\n(فقط ادمین اصلی می‌تواند این کار را انجام دهد)",
			})
			return nil
		}

		if update.NewMessage.ReplyToMessageID == "" {
			client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   "❌ لطفاً روی پیام کاربری که می‌خواهید بن کنید، ریپلای کنید.",
			})
			return nil
		}

		messageCache.RLock()
		targetUserID, found := messageCache.data[update.NewMessage.ReplyToMessageID]
		messageCache.RUnlock()

		if !found {
			client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   "❌ پیام در حافظه یافت نشد.",
			})
			return nil
		}

		err := client.BanChatMember(ctx, &ParsRubika.BanChatMemberRequest{
			ChatID: update.ChatID,
			UserID: targetUserID,
		})

		if err != nil {
			client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   fmt.Sprintf("❌ خطا: %v", err),
			})
			return err
		}

		client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   fmt.Sprintf("✅ کاربر `%s` توسط ادمین بن شد.", targetUserID),
		})
		return nil
	})

	// --- هندلر انبن کردن (فقط برای ادمین اصلی) ---
	client.OnCommand("unban", func(ctx context.Context, update *ParsRubika.Update) error {
		senderID := update.NewMessage.SenderID

		if AdminID == "" || senderID != AdminID {
			client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   "⛔ شما اجازه دسترسی ندارید.",
			})
			return nil
		}

		if update.NewMessage.ReplyToMessageID == "" {
			client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   "❌ لطفاً روی پیام کاربری که می‌خواهید انبن کنید، ریپلای کنید.",
			})
			return nil
		}

		messageCache.RLock()
		targetUserID, found := messageCache.data[update.NewMessage.ReplyToMessageID]
		messageCache.RUnlock()

		if !found {
			client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   "❌ پیام در حافظه یافت نشد.",
			})
			return nil
		}

		err := client.UnbanChatMember(ctx, &ParsRubika.UnbanChatMemberRequest{
			ChatID: update.ChatID,
			UserID: targetUserID,
		})

		if err != nil {
			client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
				ChatID: update.ChatID,
				Text:   fmt.Sprintf("❌ خطا: %v", err),
			})
			return err
		}

		client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   fmt.Sprintf("✅ کاربر `%s` انبن شد.", targetUserID),
		})
		return nil
	})

	// --- مرحله 5: اجرای ربات ---
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Printf("بات با موفقیت راه‌اندازی شد. در حال اتصال به API: %s", client.GetAPIManager().GetCurrentAPI())

	switch runningMode {
	case "polling":
		log.Println("🔄 ربات در حال اجرا با Polling...")
		err := client.Run(ctx, ParsRubika.PollingOptions{})
		if err != nil {
			log.Fatalf("خطا در اجرای ربات: %v", err)
		}

	case "webhook":
		portStr := promptUser("پورت وب‌هوک را وارد کنید (مثال: 8080): ")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("پورت نامعتبر است: %v", err)
		}
		log.Printf("🌐 ربات در حال اجرا با Webhook روی پورت %d...", port)
		err = client.StartWebhook(ctx, ParsRubika.WebhookOptions{
			Port:    port,
			Path:    "/webhook",
			Handler: client.ProcessUpdate,
		})
		if err != nil {
			log.Fatalf("خطا در اجرای وب‌هوک: %v", err)
		}

	case "host_reload":
		log.Println("♻️ ربات در حال اجرا با Polling و قابلیت Host-Reload...")
		log.Println("هر تغییر در فایل‌های .go باعث کامپایل مجدد و ری‌استارت ربات خواهد شد.")
		err := client.Run(ctx, ParsRubika.PollingOptions{})
		if err != nil {
			log.Fatalf("خطا در اجرای ربات: %v", err)
		}
	}

	log.Println("✅ ربات با موفقیت متوقف شد.")
}
