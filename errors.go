package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

import (
	"errors"
	"fmt"
)

// APIError ساختاری برای خطاهای بازگشتی از API روبیکا
type APIError struct {
	StatusCode int    `json:"-"` // وضعیت HTTP
	Message    string `json:"message"`
	Code       string `json:"code,omitempty"` // کد خطای اختصاصی API (در صورت وجود)
}

// Error متد رابط خطا برای نمایش پیام خطا
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("خطای API روبیکا (کد: %s، وضعیت %d): %s", e.Code, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("خطای API روبیکا (وضعیت %d): %s", e.StatusCode, e.Message)
}

// Is بررسی اینکه آیا خطای داده شده با خطای API فعلی مطابقت دارد یا خیر
func (e *APIError) Is(target error) bool {
	var apiErr *APIError
	if errors.As(target, &apiErr) {
		return e.StatusCode == apiErr.StatusCode && e.Code == apiErr.Code
	}
	return false
}

// ==================== خطاهای عمومی و پشتیبانی نشده ====================

// ErrUnsupportedMethod خطایی برای متدهای پشتیبانی نشده
var ErrUnsupportedMethod = fmt.Errorf("این متد توسط API رسمی روبیکا پشتیبانی نمی‌شود")

// ==================== خطاهای احراز هویت و دسترسی ====================

// ErrInvalidToken خطایی برای توکن نامعتبر
var ErrInvalidToken = fmt.Errorf("توکن بات نامعتبر است")

// ErrBotBlocked خطایی برای بات مسدود شده
var ErrBotBlocked = fmt.Errorf("بات توسط روبیکا مسدود شده است")

// ErrUserBlocked خطایی برای کاربر مسدود شده
var ErrUserBlocked = fmt.Errorf("کاربر مورد نظر مسدود شده است")

// ErrPermissionDenied خطایی برای عدم دسترسی
var ErrPermissionDenied = fmt.Errorf("شما دسترسی لازم برای انجام این عملیات را ندارید")

// ==================== خطاهای مربوط به منابع (یافت نشد) ====================

// ErrChatNotFound خطایی برای چت یافت نشده
var ErrChatNotFound = fmt.Errorf("چت مورد نظر یافت نشد")

// ErrMessageNotFound خطایی برای پیام یافت نشده
var ErrMessageNotFound = fmt.Errorf("پیام مورد نظر یافت نشد")

// ErrFileNotFound خطایی برای فایل یافت نشده
var ErrFileNotFound = fmt.Errorf("فایل مورد نظر یافت نشد")

// ErrStateNotFound خطایی برای وضعیت (state) یافت نشده
var ErrStateNotFound = fmt.Errorf("وضعیت مورد نظر برای کاربر یافت نشد")

// ==================== خطاهای ورودی و اعتبارسنجی ====================

// ErrInvalidParameter خطایی برای پارامتر نامعتبر
var ErrInvalidParameter = fmt.Errorf("پارامتر ورودی نامعتبر است")

// ErrInvalidFile خطایی برای فایل نامعتبر
var ErrInvalidFile = fmt.Errorf("فایل ورودی نامعتبر است")

// ErrFileTooLarge خطایی برای فایل بسیار بزرگ
var ErrFileTooLarge = fmt.Errorf("حجم فایل بیش از حد مجاز است")

// ErrInvalidMediaType خطایی برای نوع رسانه نامعتبر
var ErrInvalidMediaType = fmt.Errorf("نوع رسانه (Media Type) نامعتبر است")

// ErrInvalidChatType خطایی برای نوع چت نامعتبر
var ErrInvalidChatType = fmt.Errorf("نوع چت نامعتبر است")

// ==================== خطاهای چرخه حیات بات ====================

// ErrBotAlreadyRunning خطایی برای تلاش در راه‌اندازی مجدد بات در حال اجرا
var ErrBotAlreadyRunning = fmt.Errorf("بات در حال حاضر در حال اجراست")

// ErrBotNotRunning خطایی برای تلاش در متوقف کردن باتی که در حال اجرا نیست
var ErrBotNotRunning = fmt.Errorf("بات در حال حاضر متوقف است")

// ==================== خطاهای شبکه و سرور ====================

// ErrTooManyRequests خطایی برای درخواست‌های زیاد
var ErrTooManyRequests = fmt.Errorf("تعداد درخواست‌ها بیش از حد مجاز است. لطفاً کمی صبر کنید")

// ErrInternalServer خطایی برای خطای داخلی سرور
var ErrInternalServer = fmt.Errorf("خطای داخلی سرور روبیکا رخ داده است")

// ErrNetwork خطایی برای خطای شبکه
var ErrNetwork = fmt.Errorf("خطا در اتصال به شبکه")

// ErrTimeout خطایی برای اتمام زمان انتظار (Timeout)
var ErrTimeout = fmt.Errorf("زمان انتظار برای پاسخ سرور به پایان رسید")

// ==================== خطاهای عملیاتی خاص ====================

// ErrUploadFailed خطایی برای شکست در آپلود فایل
var ErrUploadFailed = fmt.Errorf("آپلود فایل با شکست مواجه شد")

// ErrDownloadFailed خطایی برای شکست در دانلود فایل
var ErrDownloadFailed = fmt.Errorf("دانلود فایل با شکست مواجه شد")

// ErrHandlerFailed خطایی برای شکست در اجرای یک هندلر
var ErrHandlerFailed = fmt.Errorf("اجرای هندلر با شکست مواجه شد")

// ErrWebhookSetupFailed خطایی برای شکست در تنظیم وب‌هوک
var ErrWebhookSetupFailed = fmt.Errorf("تنظیم وب‌هوک با شکست مواجه شد")

// ErrPollingFailed خطایی برای شکست در فرآیند پولینگ
var ErrPollingFailed = fmt.Errorf("فرآیند دریافت آپدیت‌ها (پولینگ) با شکست مواجه شد")

// ErrUserBlockedByAntiSpam خطایی برای مسدود شدن کاربر توسط سیستم ضد اسپم
var ErrUserBlockedByAntiSpam = fmt.Errorf("کاربر به دلیل فعالیت مشکوک توسط سیستم ضد اسپم مسدود شد")

// ErrHotReloadFailed خطایی برای شکست در عملیات Hot-Reload
var ErrHotReloadFailed = fmt.Errorf("عملیات بارگذاری مجدد (Hot-Reload) با شکست مواجه شد")

// ErrAPISwitchFailed خطایی برای شکست در سوییچ API
var ErrAPISwitchFailed = fmt.Errorf("عملیات سوییچ API با شکست مواجه شد")

// ErrHostReloadFailed خطایی برای شکست در عملیات Host-Reload
var ErrHostReloadFailed = fmt.Errorf("عملیات بارگذاری مجدد هاست (Host-Reload) با شکست مواجه شد")
