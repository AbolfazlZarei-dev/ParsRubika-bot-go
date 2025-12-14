# ParsRubika-bot-go

# 🤖 ParsRubika | کتابخانه روبیکا برای Go (Golang) – ساخت ربات روبیکا با Go

**ParsRubika** یک **کتابخانه کامل، فارسی و متن‌باز برای ساخت ربات روبیکا با زبان Go (Golang)** است. اگر به دنبال یک **کتابخانه روبیکا برای Go** هستید که سریع، امن، حرفه‌ای و مناسب پروژه‌های واقعی باشد، **ParsRubika بهترین انتخاب شماست**.




<div align="center">

![Rubika Bot](https://img.shields.io/badge/Rubika-Bot%20API-red?style=for-the-badge&logo=telegram&logoColor=white)
![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Stars](https://img.shields.io/github/stars/Abolfazl-Zarei/ParsRubika-bot-go?style=for-the-badge&color=gold)

<br />

<img src="https://s6.uupload.ir/files/img_20251127_142046_212_c3oj.jpg" width="150"/>
<img src="https://img.icons8.com/color/96/000000/golang.png" width="80"/>
<img src="https://img.icons8.com/color/96/000000/robot-2.png" width="80"/>
<img src="https://img.icons8.com/color/96/000000/api.png" width="80"/>

**🔗 مخزن گیت‌هاب:**  
[[https://github.com/Abolfazl-Zarei/ParsRubika-bot-go](https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go)](https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go)

**👤 توسعه‌دهنده:** ابوالفضل زارعی  
**📧 ایمیل:** `ninjacode.ir@gmail.com`  
**🆔 روبیکا:** `NinjaCode`  
**📢 چنل روبیکا:** `Ninja_code`

</div>

## 🎯 معرفی

**ParsRubika** یک کتابخانه **فارسی** و **قدرتمند** برای ساخت ربات‌های پیام‌رسان **روبیکا** با زبان **Go** است. این کتابخانه با معماری مدرن و امکانات پیشرفته، توسعه ربات‌های حرفه‌ای را برای شما ساده می‌کند.

### 🌟 ویژگی‌های کلیدی

| ویژگی | توضیح |
|-------|-----------|
| **✅ پشتیبانی کامل API** | پوشش کامل تمام متدهای رسمی روبیکا |
| **🔄 Polling & Webhook** | دو روش مدرن برای دریافت آپدیت‌ها |
| **🎛 مدیریت وضعیت** | State Management داخلی برای مدیریت مکالمات کاربران |
| **📁 مدیریت فایل‌ها** | آپلود و دانلود انواع فایل‌ها با سادگی |
| **⌨ کیبوردهای پویا** | پشتیبانی کامل از کیبوردهای Reply و Inline |
| **🛡 خطایابی هوشمند** | مدیریت خودکار خطاها و تلاش مجدد هوشمند |
| **🔄 API Switcher** | سوییچ خودکار بین APIهای روبیکا برای پایداری بالا |
| **🔄 Hot-Reload** | بارگذاری مجدد کد بدون نیاز به ری‌استارت دستی |
| **🛡️ Anti-Spam** | سیستم داخلی برای جلوگیری از اسپم کاربران |

با این کتابخانه می‌توانید:

* 🤖 **ربات روبیکا با Go بسازید**
* 📦 از یک **کتابخانه روبیکا با زبان Go**
* 🚀 ربات‌های **سریع، پایدار و مقیاس‌پذیر**
* 🔌 اتصال مستقیم و پایدار به **API روبیکا**
* ⚡ توسعه سریع بات‌های حرفه‌ای روبیکا
| **📚 مستندات کامل** | راهنمای جامع و مثال‌های کاربردی |

---

## ⚡ نصب و راه‌اندازی

### 📋 پیش‌نیازها

1.  **نصب Go (نسخه 1.21 یا بالاتر):**
    ```bash
    # اوبونتی/دبیان
    sudo apt update && sudo apt install golang-go

    # مک
    brew install go

    # ویندوز
    # از سایت https://golang.org/dl دانلود کنید
    ```

2.**دریافت توکن ربات از `@BotFather` در روبیکا.**

### 📥 نصب کتابخانه

```bash
go get github.com/Abolfazl-Zarei/ParsRubika-bot-go/v2
```

---

## 🚀 شروع سریع

یک ربات ساده "سلام دنیا" در چند دقیقه:

```go
// main.go
package main

import (
    "context"
    "log"
    "os"
    "time"

    ParsRubika "github.com/Abolfazl-Zarei/ParsRubika-bot-go/v2"
)

func main() {
    // 🔑 دریافت توکن از متغیر محیطی
    botToken := os.Getenv("RUBIKA_BOT_TOKEN")
    if botToken == "" {
        log.Fatal("❌ توکن ربات یافت نشد! متغیر محیطی RUBIKA_BOT_TOKEN را تنظیم کنید.")
    }

    // 🤖 ایجاد نمونه ربات
    bot := ParsRubika.NewClient(botToken,
        ParsRubika.WithRateLimitDelay(1*time.Second),
        ParsRubika.WithConnectionMode(ParsRubika.SwitcherMode), // حالت هوشمند سوییچر
    )

    // 🎯 تنظیم هندلر برای پیام‌های متنی
    bot.OnMessageUpdates(func(ctx context.Context, update *ParsRubika.Update) error {
        if update.NewMessage != nil && update.NewMessage.Text != "" {
            // 📨 پاسخ به پیام کاربر
            _, err := bot.SendMessage(ctx, &ParsRubika.SendMessageRequest{
                ChatID: update.ChatID,
                Text:   "👋 سلام! من با ParsRubika ساخته شده‌ام! 🚀",
            })
            return err
        }
        return nil
    })

    // 🚀 اجرای ربات
    ctx := context.Background()
    log.Println("🤖 ربات با موفقیت راه‌اندازی شد...")
    if err := bot.Run(ctx); err != nil {
        log.Fatal("💥 خطا در اجرای ربات:", err)
    }
}
```

**اجرای ربات:**

```bash
# تنظیم توکن
export RUBIKA_BOT_TOKEN="your_bot_token_here"

# اجرا
go run main.go
```

---

## 📡 API Reference

### 💬 ارسال پیام

#### ارسال پیام متنی
```go
_, err := bot.SendMessage(ctx, &ParsRubika.SendMessageRequest{
    ChatID: "CHAT_ID",
    Text:   "این یک پیام تست است!",
})
```

#### ارسال پیام با کیبورد اینلاین (Inline)
```go
keyboard := &ParsRubika.InlineKeyboardMarkup{
    InlineKeyboard: [][]*ParsRubika.InlineKeyboardButton{
        {
            {Text: "👍", CallbackData: "like"},
            {Text: "👎", CallbackData: "dislike"},
        },
    },
}

_, err := bot.SendMessage(ctx, &ParsRubika.SendMessageRequest{
    ChatID:              "CHAT_ID",
    Text:                "این پیام را دوست داشتید؟",
    InlineKeyboardMarkup: keyboard,
})
```

#### ارسال پیام با کیبورد پاسخ (Reply)
```go
keyboard := &ParsRubika.ReplyKeyboardMarkup{
    Keyboard: [][]*ParsRubika.KeyboardButton{
        {{Text: "📚 راهنما"}, {Text: "👤 اطلاعات"}},
        {{Text: "🔊 تکرار متن"}, {Text: "💾 وضعیت"}},
    },
    ResizeKeyboard: true,
}

_, err := bot.SendMessage(ctx, &ParsRubika.SendMessageRequest{
    ChatID:             "CHAT_ID",
    Text:               "منوی اصلی:",
    ReplyKeyboardMarkup: keyboard,
})
```

### 🖼️ ارسال مدیا

```go
// ارسال عکس
_, err := bot.SendPhoto(ctx, "CHAT_ID", "path/to/image.jpg", "این یک عکس است.")

// ارسال ویدیو
_, err := bot.SendVideo(ctx, "CHAT_ID", "path/to/video.mp4", "این یک ویدیو است.")

// ارسال فایل
_, err := bot.SendDocument(ctx, "CHAT_ID", "path/to/file.pdf", "این یک فایل است.")
```

### ✏️ ویرایش و مدیریت پیام‌ها

```go
// ویرایش متن پیام
err := bot.EditMessageText(ctx, &ParsRubika.EditMessageTextRequest{
    ChatID:    "CHAT_ID",
    MessageID: "MESSAGE_ID",
    Text:      "متن جدید ویرایش شده.",
})

// حذف پیام
err := bot.DeleteMessage(ctx, &ParsRubika.DeleteMessageRequest{
    ChatID:    "CHAT_ID",
    MessageID: "MESSAGE_ID",
})

// پین کردن پیام
err := bot.SetPin(ctx, "CHAT_ID", "MESSAGE_ID")
```

### 🎯 مدیریت هندلر‌ها

هندلرها قلب منطق ربات شما هستند. ParsRubika روش‌های ساده‌ای برای ثبت هندلر بر اساس نوع آپدیت فراهم می‌کند.

```go
// هندلر برای تمام پیام‌های جدید
bot.OnMessageUpdates(func(ctx context.Context, update *ParsRubika.Update) error {
    // منطق پردازش پیام
    return nil
})

// هندلر برای یک دستور خاص
bot.OnCommand("start", func(ctx context.Context, update *ParsRubika.Update) error {
    _, err := bot.SendMessage(ctx, &ParsRubika.SendMessageRequest{
        ChatID: update.ChatID,
        Text:   "به ربات خوش آمدید!",
    })
    return err
})

// هندلر برای عکس‌ها
bot.OnPhoto(func(ctx context.Context, update *ParsRubika.Update) error {
    _, err := bot.SendMessage(ctx, &ParsRubika.SendMessageRequest{
        ChatID: update.ChatID,
        Text:   "عکس زیبایی است!",
    })
    return err
})

// هندلر برای دکمه‌های اینلاین
bot.OnCallbackQuery(func(ctx context.Context, update *ParsRubika.Update) error {
    if update.NewMessage != nil && update.NewMessage.AuxData != nil {
        buttonID := *update.NewMessage.AuxData.ButtonID
        // منطق بر اساس buttonID
    }
    return nil
})
```

### 📊 مدیریت وضعیت (State Management)

برای مدیریت مکالمات چند مرحله‌ای، می‌توانید از State Manager داخلی استفاده کنید.

```go
// ذخیره وضعیت کاربر
bot.SetState(userID, "step", "awaiting_name")

// بازیابی وضعیت کاربر
if step, exists := bot.GetState(userID, "step"); exists {
    if step == "awaiting_name" {
        // منطق مرحله دریافت نام
    }
}

// حذف وضعیت کاربر
bot.DeleteState(userID, "step")
```

### 🛡️ سیستم ضد اسپم

```go
// بررسی اینکه آیا کاربر اسپم می‌کند یا خیر
if bot.CheckAntiSpam(userID) {
    // کاربر اسپم نمی‌کند، پردازش را ادامه بده
    processUserMessage()
} else {
    // کاربر در حال اسپم کردن است
    bot.SendMessage(ctx, &ParsRubika.SendMessageRequest{
        ChatID: chatID,
        Text:   "⚠️ لطفاً آرام‌تر باشید و درخواست‌ها را با فاصله ارسال کنید.",
    })
}
```

---

## 🌐 دریافت آپدیت‌ها

### 📡 Polling (روش پیش‌فرض)

در این روش، ربات به طور مداوم از سرور روبیکا آپدیت‌های جدید را درخواست می‌کند.

```go
// این روش به صورت خودکار در bot.Run() فراخوانی می‌شود
// اما می‌توانید تنظیمات آن را تغییر دهید:
ctx := context.Background()
err := bot.StartPolling(ctx, ParsRubika.PollingOptions{
    Limit:        100,             // حداکثر تعداد آپدیت در هر درخواست
    PollInterval: 2 * time.Second, // فاصله بین درخواست‌ها
})
```

### 🌐 Webhook (برای عملکرد بالا)

در این روش، سرور روبیکا آپدیت‌ها را برای شما ارسال می‌کند. این روش برای ربات‌های پرترافیک بهتر است.

```go
webhookOpts := ParsRubika.WebhookOptions{
    Port:    8080,
    Path:    "/webhook",
    Secret:  "your_strong_secret", // برای امنیت
}

err := bot.StartWebhook(ctx, webhookOpts)
```

| ویژگی | 📡 Polling | 🌐 Webhook |
|--------|------------|------------|
| **سادگی** | ✅ بسیار ساده | ⚠️ نیاز به سرور عمومی |
| **Performance** | ⚠️ متوسط | ✅ بسیار بالا |
| **Real-time** | ❌ تأخیر دارد | ✅ فوری |

---

## 🚀 مثال پیشرفته: ربات نظرسنجی

```go
package main

import (
    // ... ایمپورت‌ها
)

func main() {
    // ... مقداردهی اولیه ربات

    bot.OnMessageUpdates(func(ctx context.Context, update *ParsRubika.Update) error {
        msg := update.NewMessage
        if msg == nil {
            return nil
        }

        switch {
        case strings.HasPrefix(msg.Text, "/create_poll"):
            // پارس کردن دستور: /create_poll سوال | گزینه۱ | گزینه۲
            parts := strings.Split(strings.TrimPrefix(msg.Text, "/create_poll "), "|")
            if len(parts) < 3 {
                // ارسال پیام خطا
                return nil
            }
            question := strings.TrimSpace(parts[0])
            options := []string{strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])}
            
            // ایجاد نظرسنجی
            _, err := bot.CreatePoll(ctx, update.ChatID, question, options)
            return err
            
        case strings.HasPrefix(msg.Text, "/poll_stats"):
            // دریافت و نمایش آمار نظرسنجی
            // این بخش نیاز به ذخیره اطلاعات نظرسنجی دارد
            // ...
        }
        return nil
    })

    // ... اجرای ربات
}
```

---

## ☁️ استقرار و دیپلوی

### 🐳 استفاده از داکر (Docker)

برای استقرار آسان، می‌توانید از داکر استفاده کنید.

**Dockerfile:**
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o rubika-bot .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/rubika-bot .
COPY --from=builder /app/config.yaml .
EXPOSE 8080
CMD ["./rubika-bot"]
```

**docker-compose.yml:**
```yaml
version: '3.8'
services:
  rubika-bot:
    build: .
    ports:
      - "8080:8080"
    environment:
      - RUBIKA_BOT_TOKEN=${RUBIKA_BOT_TOKEN}
    restart: unless-stopped
```

**اجرای پروژه:**
```bash
docker-compose up --build
```

---

## 🔧 عیب‌یابی

### 🚨 خطاهای رایج

| خطا | دلیل احتمالی | راه‌حل |
|------|-------------|--------|
| `401 Unauthorized` | توکن نامعتبر یا اشتباه | توکن را از `@BotFather` مجدداً دریافت و بررسی کنید. |
| `connection refused` | مشکل در اتصال به اینترنت یا فایروال | اتصال به اینترنت و تنظیمات فایروال را بررسی کنید. |
| `too many requests` | ارسال درخواست‌های بیش از حد مجاز (Rate Limit) | از `WithRateLimitDelay` برای افزایش فاصله بین درخواست‌ها استفاده کنید. |

### 📊 لاگینگ پیشرفته

برای دیباگ بهتر، می‌توانید یک میدلور لاگینگ اضافه کنید:

```go
bot.Middleware(func(ctx context.Context, update *ParsRubika.Update, next ParsRubika.HandlerFunc) error {
    start := time.Now()
    log.Printf("📥 آپدیت دریافت شد: %s از چت %s", update.Type, update.ChatID)
    
    err := next(ctx, update)
    
    log.Printf("✅ پردازش آپدیت در %v با خطا: %v", time.Since(start), err)
    return err
})
```

---

## 📞 پشتیبانی

- **👤 ایدی روبیکا:** `NinjaCode`
- **📢 چنل روبیکا:** `Ninja_code`
- **📧 ایمیل:** `ninjacode.ir@gmail.com`
- **🐙 گیت‌هاب:** [Abolfazl-Zarei](https://github.com/Abolfazl-Zarei)

برای گزارش باگ یا درخواست ویژگی جدید، لطفاً یک [Issue](https://github.com/Abolfazl-Zarei/ParsRubika-bot-go/issues) در گیت‌هاب ایجاد کنید.

---

## 🎉 نتیجه‌گیری

**ParsRubika** با ارائه یک API تمیز، مستندات کامل و قابلیت‌های پیشرفته مانند مدیریت وضعیت، سوییچ هوشمند API و Hot-Reload، فرآیند ساخت ربات‌های قدرتمند و پایدار برای روبیکا را به شدت ساده می‌کند. با استفاده از این کتابخانه، می‌توانید بر روی منطق اصلی ربات خود تمرکز کنید و از نگرانی در مورد جزئیات فنی پیاده‌سازی رها شوید.

### 🚀 شروع نهایی

```go
package main

// ... ایمپورت‌ها

func main() {
    bot := ParsRubika.NewClient(os.Getenv("RUBIKA_BOT_TOKEN"),
        ParsRubika.WithConnectionMode(ParsRubika.SwitcherMode),
        ParsRubika.WithHotReload(true),
    )
    
    bot.OnMessageUpdates(yourHandler)
    
    ctx := context.Background()
    log.Println("🎉 ربات شروع به کار کرد!")
    bot.Run(ctx)
}
```


