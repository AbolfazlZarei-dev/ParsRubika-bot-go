package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

import "time"

// ==================== مدل‌های اصلی و پایه ====================

// Bot اطلاعات بات
type Bot struct {
	BotID        string `json:"bot_id"`
	BotTitle     string `json:"bot_title"`
	Avatar       *File  `json:"avatar"`
	Description  string `json:"description"`
	Username     string `json:"username"`
	StartMessage string `json:"start_message"`
	ShareURL     string `json:"share_url"`
}

// User اطلاعات کاربر
type User struct {
	UserID     string `json:"user_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Username   string `json:"username"`
	Bio        string `json:"bio"`
	IsVerified bool   `json:"is_verified"`
	IsPrivate  bool   `json:"is_private"`
	Avatar     *File  `json:"avatar"`
	LastSeen   string `json:"last_seen"`
}

// Chat اطلاعات چت
type Chat struct {
	ChatID    string       `json:"chat_id"`
	ChatType  ChatTypeEnum `json:"chat_type"`
	UserID    string       `json:"user_id"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Title     string       `json:"title"`
	Username  string       `json:"username"`
}

// File اطلاعات فایل
type File struct {
	FileID   string       `json:"file_id"`
	FileName string       `json:"file_name"`
	Size     int64        `json:"size"`
	Type     FileTypeEnum `json:"type,omitempty"`
}

// AuxData داده‌های کمکی
type AuxData struct {
	StartID  *string `json:"start_id"`
	ButtonID *string `json:"button_id"`
}

// ForwardedFrom اطلاعات فوروارد پیام
type ForwardedFrom struct {
	TypeFrom     ForwardedFromEnum `json:"type_from"`
	MessageID    string            `json:"message_id"`
	FromChatID   string            `json:"from_chat_id"`
	FromSenderID string            `json:"from_sender_id"`
}

// PaymentStatus وضعیت پرداخت
type PaymentStatus struct {
	PaymentID string            `json:"payment_id"`
	Status    PaymentStatusEnum `json:"status"`
}

// ==================== مدل‌های پیام (Message) و محتوای آن ====================

// Message پیام
type Message struct {
	MessageID        string                 `json:"message_id"`
	Text             string                 `json:"text"`
	Time             string                 `json:"time"`
	IsEdited         bool                   `json:"is_edited"`
	SenderType       MessageSenderEnum      `json:"sender_type"`
	SenderID         string                 `json:"sender_id"`
	AuxData          *AuxData               `json:"aux_data"`
	File             *File                  `json:"file"`
	ReplyToMessageID string                 `json:"reply_to_message_id"`
	ForwardedFrom    *ForwardedFrom         `json:"forwarded_from"`
	ForwardedNoLink  map[string]interface{} `json:"forwarded_no_link"` // <--- این خط قبلا اصلاح شد
	Location         *Location              `json:"location"`
	Sticker          *Sticker               `json:"sticker"`
	ContactMessage   *ContactMessage        `json:"contact_message"`
	Poll             *Poll                  `json:"poll"`
	Payment          *PaymentStatus         `json:"payment"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	IsMarkdown       bool                   `json:"is_markdown,omitempty"`
	Font             string                 `json:"font,omitempty"`
}

// Sticker استیکر
type Sticker struct {
	StickerID      string `json:"sticker_id"`
	File           *File  `json:"file"`
	EmojiCharacter string `json:"emoji_character"`
}

// ContactMessage پیام مخاطب
type ContactMessage struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

// Poll نظرسنجی
type Poll struct {
	Question   string      `json:"question"`
	Options    []string    `json:"options"`
	PollStatus *PollStatus `json:"poll_status"`
}

// PollStatus وضعیت نظرسنجی
type PollStatus struct {
	State              PollStatusEnum `json:"state"`
	SelectionIndex     int            `json:"selection_index"`
	PercentVoteOptions []int          `json:"percent_vote_options"`
	TotalVote          int            `json:"total_vote"`
	ShowTotalVotes     bool           `json:"show_total_votes"`
}

// Location موقعیت مکانی
type Location struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

// ==================== مدل‌های دکمه و کیبورد ====================

// LoginUrl اطلاعات برای دکمه لاگین
type LoginUrl struct {
	URL                string  `json:"url"`
	ForwardText        *string `json:"forward_text,omitempty"`
	BotUsername        *string `json:"bot_username,omitempty"`
	RequestWriteAccess *bool   `json:"request_write_access,omitempty"`
}

// WebAppInfo اطلاعات برای دکمه وب اپ
type WebAppInfo struct {
	URL string `json:"url"`
}

// InlineKeyboardButton دکمه کیبورد اینلاین
type InlineKeyboardButton struct {
	Text                         string    `json:"text"`
	URL                          string    `json:"url,omitempty"`
	LoginURL                     *LoginUrl `json:"login_url,omitempty"`
	CallbackData                 string    `json:"callback_data,omitempty"`
	SwitchInlineQuery            string    `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string    `json:"switch_inline_query_current_chat,omitempty"`
}

// InlineKeyboardMarkup کیبورد اینلاین
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]*InlineKeyboardButton `json:"inline_keyboard"`
}

// KeyboardButton دکمه کیبورد پاسخ
type KeyboardButton struct {
	Text            string                  `json:"text"`
	RequestContact  bool                    `json:"request_contact,omitempty"`
	RequestLocation bool                    `json:"request_location,omitempty"`
	RequestPoll     *KeyboardButtonPollType `json:"request_poll,omitempty"`
	WebApp          *WebAppInfo             `json:"web_app,omitempty"`
}

// KeyboardButtonPollType نوع نظرسنجی برای دکمه کیبورد
type KeyboardButtonPollType struct {
	Type string `json:"type"`
}

// ReplyKeyboardMarkup کیبورد پاسخ
type ReplyKeyboardMarkup struct {
	Keyboard              [][]*KeyboardButton `json:"keyboard"`
	IsPersistent          bool                `json:"is_persistent,omitempty"`
	ResizeKeyboard        bool                `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool                `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string              `json:"input_field_placeholder,omitempty"`
	Selective             bool                `json:"selective,omitempty"`
}

// ReplyKeyboardRemove حذف کیبورد پاسخ
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective,omitempty"`
}

// ==================== مدل‌های آپدیت (Update) ====================

// Update آپدیت
type Update struct {
	Type             UpdateTypeEnum `json:"type"`
	ChatID           string         `json:"chat_id"`
	RemovedMessageID *string        `json:"removed_message_id"`
	NewMessage       *Message       `json:"new_message"`
	UpdatedMessage   *Message       `json:"updated_message"`
	UpdatedPayment   *PaymentStatus `json:"updated_payment"`
}

// InlineMessage پیام اینلاین
type InlineMessage struct {
	SenderID  string    `json:"sender_id"`
	Text      string    `json:"text"`
	File      *File     `json:"file"`
	Location  *Location `json:"location"`
	AuxData   *AuxData  `json:"aux_data"`
	MessageID string    `json:"message_id"`
	ChatID    string    `json:"chat_id"`
}

// ==================== مدل‌های درخواست (Request) ====================

// BaseRequest مدل پایه برای درخواست‌ها
type BaseRequest struct{}

// EditMessageKeypadRequest درخواست ویرایش کیبورد پیام (اینلاین)
type EditMessageKeypadRequest struct {
	ChatID         string                `json:"chat_id"`
	MessageID      string                `json:"message_id"`
	InlineKeyboard *InlineKeyboardMarkup `json:"inline_keyboard"`
}

// EditChatKeypadRequest درخواست ویرایش کیبورد چت (پاسخ)
type EditChatKeypadRequest struct {
	ChatID         string               `json:"chat_id"`
	ReplyKeyboard  *ReplyKeyboardMarkup `json:"reply_keyboard,omitempty"`
	RemoveKeyboard *ReplyKeyboardRemove `json:"remove_keyboard,omitempty"`
	ChatKeypadType ChatKeypadTypeEnum   `json:"chat_keypad_type"`
}

// SendMessageRequest درخواست ارسال پیام
type SendMessageRequest struct {
	ChatID               string                `json:"chat_id"`
	Text                 string                `json:"text"`
	DisableNotification  bool                  `json:"disable_notification,omitempty"`
	ReplyToMessageID     string                `json:"reply_to_message_id,omitempty"`
	InlineKeyboardMarkup *InlineKeyboardMarkup `json:"inline_keyboard,omitempty"`
	ReplyKeyboardMarkup  *ReplyKeyboardMarkup  `json:"reply_keyboard,omitempty"`
}

// SendPollRequest درخواست ارسال نظرسنجی
type SendPollRequest struct {
	ChatID   string   `json:"chat_id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

// GetUpdatesRequest درخواست دریافت آپدیت‌ها
type GetUpdatesRequest struct {
	OffsetID string `json:"offset_id,omitempty"`
	Limit    int    `json:"limit"`
}

// GetChatRequest درخواست دریافت اطلاعات چت
type GetChatRequest struct {
	ChatID string `json:"chat_id"`
}

// SendLocationRequest درخواست ارسال موقعیت مکانی
type SendLocationRequest struct {
	ChatID               string                `json:"chat_id"`
	Latitude             string                `json:"latitude"`
	Longitude            string                `json:"longitude"`
	DisableNotification  bool                  `json:"disable_notification,omitempty"`
	ReplyToMessageID     string                `json:"reply_to_message_id,omitempty"`
	InlineKeyboardMarkup *InlineKeyboardMarkup `json:"inline_keyboard,omitempty"`
	ReplyKeyboardMarkup  *ReplyKeyboardMarkup  `json:"reply_keyboard,omitempty"`
}

// SendContactRequest درخواست ارسال مخاطب
type SendContactRequest struct {
	ChatID               string                `json:"chat_id"`
	FirstName            string                `json:"first_name"`
	LastName             string                `json:"last_name"`
	PhoneNumber          string                `json:"phone_number"`
	DisableNotification  bool                  `json:"disable_notification,omitempty"`
	ReplyToMessageID     string                `json:"reply_to_message_id,omitempty"`
	InlineKeyboardMarkup *InlineKeyboardMarkup `json:"inline_keyboard,omitempty"`
	ReplyKeyboardMarkup  *ReplyKeyboardMarkup  `json:"reply_keyboard,omitempty"`
}

// ForwardMessageRequest درخواست فوروارد پیام
type ForwardMessageRequest struct {
	FromChatID          string `json:"from_chat_id"`
	MessageID           string `json:"message_id"`
	ToChatID            string `json:"to_chat_id"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
}

// EditMessageTextRequest درخواست ویرایش متن پیام
type EditMessageTextRequest struct {
	ChatID               string                `json:"chat_id"`
	MessageID            string                `json:"message_id"`
	Text                 string                `json:"text"`
	InlineKeyboardMarkup *InlineKeyboardMarkup `json:"inline_keyboard,omitempty"`
}

// DeleteMessageRequest درخواست حذف پیام
type DeleteMessageRequest struct {
	ChatID    string `json:"chat_id"`
	MessageID string `json:"message_id"`
}

// SetCommandsRequest درخواست تنظیم دستورات بات
type SetCommandsRequest struct {
	BotCommands []BotCommand `json:"bot_commands"`
}

// BotCommand دستورات بات
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// GetFileRequest درخواست دریافت اطلاعات فایل
type GetFileRequest struct {
	FileID string `json:"file_id"`
}

// SendFileRequest درخواست ارسال فایل
type SendFileRequest struct {
	ChatID               string                `json:"chat_id"`
	FileID               string                `json:"file_id"`
	Text                 string                `json:"text,omitempty"`
	ReplyToMessageID     string                `json:"reply_to_message_id,omitempty"`
	DisableNotification  bool                  `json:"disable_notification,omitempty"`
	InlineKeyboardMarkup *InlineKeyboardMarkup `json:"inline_keyboard,omitempty"`
	ReplyKeyboardMarkup  *ReplyKeyboardMarkup  `json:"reply_keyboard,omitempty"`
}

// RequestSendFileRequest درخواست آدرس آپلود
type RequestSendFileRequest struct {
	Type FileTypeEnum `json:"type"`
}

// BanChatMemberRequest درخواست مسدود کردن عضو چت
type BanChatMemberRequest struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

// UnbanChatMemberRequest درخواست رفع مسدودیت عضو چت
type UnbanChatMemberRequest struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

// PinChatMessageRequest درخواست پین کردن پیام
type PinChatMessageRequest struct {
	ChatID              string `json:"chat_id"`
	MessageID           string `json:"message_id"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
}

// UnpinChatMessageRequest درخواست آنپین کردن پیام
type UnpinChatMessageRequest struct {
	ChatID    string `json:"chat_id"`
	MessageID string `json:"message_id"`
}

// PromoteChatMemberRequest درخواست ارتقای عضو چت
type PromoteChatMemberRequest struct {
	ChatID              string `json:"chat_id"`
	UserID              string `json:"user_id"`
	IsAdministrator     *bool  `json:"is_administrator,omitempty"`
	CanChangeInfo       *bool  `json:"can_change_info,omitempty"`
	CanPostMessages     *bool  `json:"can_post_messages,omitempty"`
	CanEditMessages     *bool  `json:"can_edit_messages,omitempty"`
	CanDeleteMessages   *bool  `json:"can_delete_messages,omitempty"`
	CanInviteUsers      *bool  `json:"can_invite_users,omitempty"`
	CanPinMessages      *bool  `json:"can_pin_messages,omitempty"`
	CanManageVideoChats *bool  `json:"can_manage_video_chats,omitempty"`
	CanManageChat       *bool  `json:"can_manage_chat,omitempty"`
	CanManageTopics     *bool  `json:"can_manage_topics,omitempty"`
}

// GetChatMemberRequest درخواست دریافت اطلاعات عضو چت
type GetChatMemberRequest struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

// ChatMember عضو چت
type ChatMember struct {
	User     Chat   `json:"user"`
	Status   string `json:"status"`
	JoinDate string `json:"join_date"`
}

// GetChatAdministratorsRequest درخواست دریافت مدیران چت
type GetChatAdministratorsRequest struct {
	ChatID string `json:"chat_id"`
}

// GetChatMemberCountRequest درخواست دریافت تعداد اعضای چت
type GetChatMemberCountRequest struct {
	ChatID string `json:"chat_id"`
}

// UpdateBotEndpointsRequest درخواست به‌روزرسانی endpointهای بات
type UpdateBotEndpointsRequest struct {
	URL  string                 `json:"url"`
	Type UpdateEndpointTypeEnum `json:"type"`
}

// GetSelectionItemRequest درخواست دریافت آیتم‌های انتخابی
type GetSelectionItemRequest struct {
	SelectionID string `json:"selection_id"`
	Query       string `json:"query,omitempty"`
}

// SearchSelectionItemsRequest درخواست جستجوی آیتم‌های انتخابی
type SearchSelectionItemsRequest struct {
	SelectionID string `json:"selection_id"`
	Query       string `json:"query"`
}

// ReceiveUpdateRequest درخواست دریافت آپدیت‌ها (برای وب‌هوک)
type ReceiveUpdateRequest struct {
	Update *Update `json:"update"`
}

// ReceiveInlineMessageRequest درخواست دریافت پیام اینلاین (برای وب‌هوک)
type ReceiveInlineMessageRequest struct {
	InlineMessage *InlineMessage `json:"inline_message"`
}

// ReceiveQueryRequest درخواست دریافت کوئری (برای وب‌هوک)
type ReceiveQueryRequest struct {
	QueryID  string `json:"query_id"`
	ChatID   string `json:"chat_id"`
	SenderID string `json:"sender_id"`
	Data     string `json:"data"`
}

// ==================== مدل‌های پاسخ (Response) ====================

// BaseResponse مدل پایه برای پاسخ‌ها
type BaseResponse struct {
	OK          bool                `json:"ok"`
	ErrorCode   int                 `json:"error_code,omitempty"`
	Description string              `json:"description,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

// ResponseParameters پارامترهای پاسخ
type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`
	RetryAfter      int   `json:"retry_after,omitempty"`
}

// MessageIDResponse پاسخ حاوی آیدی پیام
type MessageIDResponse struct {
	MessageID string `json:"message_id"`
}

// GetUpdatesResponse پاسخ دریافت آپدیت‌ها
type GetUpdatesResponse struct {
	Updates      []Update `json:"updates"`
	NextOffsetID string   `json:"next_offset_id"`
}

// ForwardMessageResponse پاسخ فوروارد پیام
type ForwardMessageResponse struct {
	NewMessageID string `json:"new_message_id"`
}

// GetFileResponse پاسخ دریافت اطلاعات فایل
type GetFileResponse struct {
	DownloadURL string `json:"download_url"`
	FileName    string `json:"file_name"`
	Size        int64  `json:"size"`
}

// RequestSendFileResponse پاسخ حاوی آدرس آپلود
type RequestSendFileResponse struct {
	UploadURL string `json:"upload_url"`
}

// FileUploadResponse پاسخ آپلود فایل
type FileUploadResponse struct {
	FileID string `json:"file_id"`
}

// GetChatMemberResponse پاسخ دریافت اطلاعات عضو چت
type GetChatMemberResponse struct {
	Member ChatMember `json:"member"`
}

// GetChatAdministratorsResponse پاسخ دریافت مدیران چت
type GetChatAdministratorsResponse struct {
	Administrators []ChatMember `json:"administrators"`
}

// GetChatMemberCountResponse پاسخ دریافت تعداد اعضای چت
type GetChatMemberCountResponse struct {
	Count int `json:"count"`
}

// GetSelectionItemResponse پاسخ دریافت آیتم‌های انتخابی
type GetSelectionItemResponse struct {
	Items []ButtonSelectionItem `json:"items"`
}

// SearchSelectionItemsResponse پاسخ جستجوی آیتم‌های انتخابی
type SearchSelectionItemsResponse struct {
	Items []ButtonSelectionItem `json:"items"`
}

// ButtonSelectionItem آیتم در دکمه انتخاب
type ButtonSelectionItem struct {
	Text     string                  `json:"text"`
	ImageUrl string                  `json:"image_url"`
	Type     ButtonSelectionTypeEnum `json:"type"`
}

// ==================== مدل‌های جدید (نسخه 2.0.0) ====================

// HealthStatus وضعیت سلامت یک API
type HealthStatus struct {
	APIType      APIType       `json:"api_type"`
	IsHealthy    bool          `json:"is_healthy"`
	ResponseTime time.Duration `json:"response_time"`
	LastChecked  time.Time     `json:"last_checked"`
	ErrorCount   int           `json:"error_count"`
	LastError    error         `json:"last_error,omitempty"`
}

// APIEvent رویداد تغییر API
type APIEvent struct {
	Timestamp time.Time `json:"timestamp"`
	FromAPI   APIType   `json:"from_api"`
	ToAPI     APIType   `json:"to_api"`
	Reason    string    `json:"reason"`
}

// NotificationOptions گزینه‌های اطلاع‌رسانی سوییچ API
type NotificationOptions struct {
	Enabled     bool   `json:"enabled"`
	SendToOwner bool   `json:"send_to_owner"`
	SendToAll   bool   `json:"send_to_all"`
	Message     string `json:"message"`
	ChatID      string `json:"chat_id"` // برای ارسال به چت خاص (مالک)
}

// NetworkEvent رویداد شبکه برای تحلیل هوشمند
type NetworkEvent struct {
	Timestamp    time.Time     `json:"timestamp"`
	APIType      APIType       `json:"api_type"`
	ResponseTime time.Duration `json:"response_time"`
	Success      bool          `json:"success"`
	Error        error         `json:"error,omitempty"`
}

// NetworkSample نمونه آماری برای مدل پیش‌بینی
type NetworkSample struct {
	TimeOfDay    float64       `json:"time_of_day"`
	DayOfWeek    int           `json:"day_of_week"`
	ResponseTime time.Duration `json:"response_time"`
	Success      bool          `json:"success"`
}

// UserRecord رکورد کامل از یک کاربر برای ذخیره در فایل
type UserRecord struct {
	UserID        string    `json:"user_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Username      string    `json:"username"`
	Bio           string    `json:"bio"`
	IsVerified    bool      `json:"is_verified"`
	IsPrivate     bool      `json:"is_private"`
	Avatar        *File     `json:"avatar"`
	LastSeen      string    `json:"last_seen"`
	StartedAt     time.Time `json:"started_at"`      // زمانی که کاربر بات را استارت زد
	StartChatID   string    `json:"start_chat_id"`   // آیدی چتی که استارت از آن ارسال شده
	StartChatType string    `json:"start_chat_type"` // نوع چت (Private, Group, Channel)
}
