package ParsRubika

// سازنده ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

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

// BotCommand دستورات بات
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
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
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
	Size     string `json:"size"`
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

// MessageTextUpdate آپدیت متن پیام
type MessageTextUpdate struct {
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
}

// MessageKeypadUpdate آپدیت کیبورد پیام
type MessageKeypadUpdate struct {
	MessageID    string  `json:"message_id"`
	InlineKeypad *Keypad `json:"inline_keypad"`
}

// --- شروع تعاریف جدید برای انواع دکمه ---

// ButtonSelectionItem آیتم در دکمه انتخاب
type ButtonSelectionItem struct {
	Text     string                  `json:"text"`
	ImageUrl string                  `json:"image_url"`
	Type     ButtonSelectionTypeEnum `json:"type"`
}

// ButtonSelection دکمه انتخاب
type ButtonSelection struct {
	SelectionID      string                    `json:"selection_id"`
	SearchType       ButtonSelectionSearchEnum `json:"search_type"`
	GetType          ButtonSelectionGetEnum    `json:"get_type"`
	IsMultiSelection bool                      `json:"is_multi_selection"`
	ColumnsCount     string                    `json:"columns_count"`
	Title            string                    `json:"title"`
	Items            []ButtonSelectionItem     `json:"items"`
}

// ButtonCalendar دکمه تقویم
type ButtonCalendar struct {
	DefaultValue *string                `json:"default_value"`
	Type         ButtonCalendarTypeEnum `json:"type"`
	MinYear      string                 `json:"min_year"`
	MaxYear      string                 `json:"max_year"`
	Title        string                 `json:"title"`
}

// ButtonNumberPicker دکمه انتخاب عدد
type ButtonNumberPicker struct {
	MinValue     string  `json:"min_value"`
	MaxValue     string  `json:"max_value"`
	DefaultValue *string `json:"default_value"`
	Title        string  `json:"title"`
}

// ButtonStringPicker دکمه انتخاب رشته
type ButtonStringPicker struct {
	Items        []string `json:"items"`
	DefaultValue *string  `json:"default_value"`
	Title        *string  `json:"title"`
}

// ButtonLocation دکمه موقعیت
type ButtonLocation struct {
	DefaultPointerLocation *Location              `json:"default_pointer_location"`
	DefaultMapLocation     *Location              `json:"default_map_location"`
	Type                   ButtonLocationTypeEnum `json:"type"`
	Title                  *string                `json:"title"`
	LocationImageUrl       string                 `json:"location_image_url"`
}

// ButtonTextbox دکمه کادر متنی
type ButtonTextbox struct {
	TypeLine     ButtonTextboxTypeLineEnum   `json:"type_line"`
	TypeKeypad   ButtonTextboxTypeKeypadEnum `json:"type_keypad"`
	PlaceHolder  *string                     `json:"place_holder"`
	Title        *string                     `json:"title"`
	DefaultValue *string                     `json:"default_value"`
}

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

// --- پایان تعاریف جدید ---

// Button دکمه
type Button struct {
	ID         string         `json:"id"`
	Type       ButtonTypeEnum `json:"type"`
	ButtonText string         `json:"button_text"`
	// فیلدهای اختصاصی برای انواع دکمه
	Url                *string             `json:"url,omitempty"`
	ButtonSelection    *ButtonSelection    `json:"button_selection,omitempty"`
	ButtonCalendar     *ButtonCalendar     `json:"button_calendar,omitempty"`
	ButtonNumberPicker *ButtonNumberPicker `json:"button_number_picker,omitempty"`
	ButtonStringPicker *ButtonStringPicker `json:"button_string_picker,omitempty"`
	ButtonLocation     *ButtonLocation     `json:"button_location,omitempty"`
	ButtonTextbox      *ButtonTextbox      `json:"button_textbox,omitempty"`
	// --- فیلدهای جدید برای دکمه‌های پیشرفته ---
	SwitchInlineQuery *string     `json:"switch_inline_query,omitempty"`
	LoginUrl          *LoginUrl   `json:"login_url,omitempty"`
	CallbackGame      *string     `json:"callback_game,omitempty"`
	Pay               *bool       `json:"pay,omitempty"`
	WebApp            *WebAppInfo `json:"web_app,omitempty"`
}

// KeypadRow ردیف کیبورد
type KeypadRow struct {
	Buttons []Button `json:"buttons"`
}

// Keypad کیبورد
type Keypad struct {
	Rows           []KeypadRow `json:"rows"`
	ResizeKeyboard bool        `json:"resize_keyboard"`
	OnTimeKeyboard bool        `json:"on_time_keyboard"`
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

// PollStatus وضعیت نظرسنجی
type PollStatus struct {
	State              PollStatusEnum `json:"state"`
	SelectionIndex     int            `json:"selection_index"`
	PercentVoteOptions []int          `json:"percent_vote_options"`
	TotalVote          int            `json:"total_vote"`
	ShowTotalVotes     bool           `json:"show_total_votes"`
}

// Poll نظرسنجی
type Poll struct {
	Question   string      `json:"question"`
	Options    []string    `json:"options"`
	PollStatus *PollStatus `json:"poll_status"`
}

// Location موقعیت مکانی
type Location struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

// Message پیام
type Message struct {
	// نکته: در پاسخ‌های مربوط به getUpdates و receiveUpdate، این فیلد به صورت عدد (int64) برگردانده می‌شود.
	MessageID        int64             `json:"message_id"`
	Text             string            `json:"text"`
	Time             string            `json:"time"`
	IsEdited         bool              `json:"is_edited"`
	SenderType       MessageSenderEnum `json:"sender_type"`
	SenderID         string            `json:"sender_id"`
	AuxData          *AuxData          `json:"aux_data"`
	File             *File             `json:"file"`
	ReplyToMessageID string            `json:"reply_to_message_id"`
	ForwardedFrom    *ForwardedFrom    `json:"forwarded_from"`
	ForwardedNoLink  string            `json:"forwarded_no_link"`
	Location         *Location         `json:"location"`
	Sticker          *Sticker          `json:"sticker"`
	ContactMessage   *ContactMessage   `json:"contact_message"`
	Poll             *Poll             `json:"poll"`
	Payment          *PaymentStatus    `json:"payment"`
}

// Update آپدیت
type Update struct {
	Type   UpdateTypeEnum `json:"type"`
	ChatID string         `json:"chat_id"`
	// نکته: در API این فیلد به صورت string برگردانده می‌شود.
	RemovedMessageID *string        `json:"removed_message_id"`
	NewMessage       *Message       `json:"new_message"`
	UpdatedMessage   *Message       `json:"updated_message"`
	UpdatedPayment   *PaymentStatus `json:"updated_payment"`
}

// InlineMessage پیام اینلاین
type InlineMessage struct {
	SenderID string    `json:"sender_id"`
	Text     string    `json:"text"`
	File     *File     `json:"file"`
	Location *Location `json:"location"`
	AuxData  *AuxData  `json:"aux_data"`
	// نکته: در پاسخ‌های مربوط به receiveInlineMessage، این فیلد به صورت رشته (string) برگردانده می‌شود.
	MessageID string `json:"message_id"`
	ChatID    string `json:"chat_id"`
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

// --- مدل‌های درخواست و پاسخ ---

type SendMessageRequest struct {
	ChatID              string             `json:"chat_id"`
	Text                string             `json:"text"`
	ChatKeypad          *Keypad            `json:"chat_keypad,omitempty"`
	DisableNotification bool               `json:"disable_notification,omitempty"`
	InlineKeypad        *Keypad            `json:"inline_keypad,omitempty"`
	ReplyToMessageID    string             `json:"reply_to_message_id,omitempty"`
	ChatKeypadType      ChatKeypadTypeEnum `json:"chat_keypad_type,omitempty"`
}

type SendPollRequest struct {
	ChatID   string   `json:"chat_id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type GetUpdatesRequest struct {
	OffsetID string `json:"offset_id,omitempty"`
	Limit    int    `json:"limit"`
}

type GetChatRequest struct {
	ChatID string `json:"chat_id"`
}

type SendLocationRequest struct {
	ChatID              string             `json:"chat_id"`
	Latitude            string             `json:"latitude"`
	Longitude           string             `json:"longitude"`
	ChatKeypad          *Keypad            `json:"chat_keypad,omitempty"`
	DisableNotification bool               `json:"disable_notification,omitempty"`
	InlineKeypad        *Keypad            `json:"inline_keypad,omitempty"`
	ReplyToMessageID    string             `json:"reply_to_message_id,omitempty"`
	ChatKeypadType      ChatKeypadTypeEnum `json:"chat_keypad_type,omitempty"`
}

type SendContactRequest struct {
	ChatID              string             `json:"chat_id"`
	FirstName           string             `json:"first_name"`
	LastName            string             `json:"last_name"`
	PhoneNumber         string             `json:"phone_number"`
	ChatKeypad          *Keypad            `json:"chat_keypad,omitempty"`
	DisableNotification bool               `json:"disable_notification,omitempty"`
	InlineKeypad        *Keypad            `json:"inline_keypad,omitempty"`
	ReplyToMessageID    string             `json:"reply_to_message_id,omitempty"`
	ChatKeypadType      ChatKeypadTypeEnum `json:"chat_keypad_type,omitempty"`
}

type ForwardMessageRequest struct {
	FromChatID          string `json:"from_chat_id"`
	MessageID           string `json:"message_id"`
	ToChatID            string `json:"to_chat_id"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
}

type EditMessageTextRequest struct {
	ChatID    string `json:"chat_id"`
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
}

type EditMessageKeypadRequest struct {
	ChatID       string  `json:"chat_id"`
	MessageID    string  `json:"message_id"`
	InlineKeypad *Keypad `json:"inline_keypad"`
}

type DeleteMessageRequest struct {
	ChatID    string `json:"chat_id"`
	MessageID string `json:"message_id"`
}

type SetCommandsRequest struct {
	BotCommands []BotCommand `json:"bot_commands"`
}

type EditChatKeypadRequest struct {
	ChatID         string             `json:"chat_id"`
	ChatKeypad     *Keypad            `json:"chat_keypad,omitempty"`
	ChatKeypadType ChatKeypadTypeEnum `json:"chat_keypad_type"`
}

type GetFileRequest struct {
	FileID string `json:"file_id"`
}

type SendFileRequest struct {
	ChatID              string             `json:"chat_id"`
	FileID              string             `json:"file_id"`
	Text                string             `json:"text,omitempty"`
	ReplyToMessageID    string             `json:"reply_to_message_id,omitempty"`
	DisableNotification bool               `json:"disable_notification,omitempty"`
	ChatKeypad          *Keypad            `json:"chat_keypad,omitempty"`
	InlineKeypad        *Keypad            `json:"inline_keypad,omitempty"`
	ChatKeypadType      ChatKeypadTypeEnum `json:"chat_keypad_type,omitempty"`
}

type RequestSendFileRequest struct {
	Type FileTypeEnum `json:"type"`
}

type BanChatMemberRequest struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

type UnbanChatMemberRequest struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

type PinChatMessageRequest struct {
	ChatID              string `json:"chat_id"`
	MessageID           string `json:"message_id"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
}

type UnpinChatMessageRequest struct {
	ChatID    string `json:"chat_id"`
	MessageID string `json:"message_id"`
}

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

type GetChatMemberRequest struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

type ChatMember struct {
	User     Chat   `json:"user"`
	Status   string `json:"status"`
	JoinDate string `json:"join_date"`
}

type GetChatMemberResponse struct {
	Member ChatMember `json:"member"`
}

type GetChatAdministratorsRequest struct {
	ChatID string `json:"chat_id"`
}

type GetChatAdministratorsResponse struct {
	Administrators []ChatMember `json:"administrators"`
}

type GetChatMemberCountRequest struct {
	ChatID string `json:"chat_id"`
}

type GetChatMemberCountResponse struct {
	Count int `json:"count"`
}

// --- مدل‌های پاسخ ---

type MessageIDResponse struct {
	MessageID string `json:"message_id"`
}

type GetUpdatesResponse struct {
	Updates      []Update `json:"updates"`
	NextOffsetID string   `json:"next_offset_id"`
}

type ForwardMessageResponse struct {
	NewMessageID string `json:"new_message_id"`
}

type GetFileResponse struct {
	DownloadURL string `json:"download_url"`
	FileName    string `json:"file_name"`
	Size        string `json:"size"`
}

type RequestSendFileResponse struct {
	UploadURL string `json:"upload_url"`
}

type FileUploadResponse struct {
	FileID string `json:"file_id"`
}

// --- مدل‌های درخواست و پاسخ برای متدهای جدید ---

// GetSelectionItemRequest درخواست دریافت آیتم‌های انتخابی
type GetSelectionItemRequest struct {
	SelectionID string `json:"selection_id"`
	Query       string `json:"query,omitempty"`
}

// GetSelectionItemResponse پاسخ دریافت آیتم‌های انتخابی
type GetSelectionItemResponse struct {
	Items []ButtonSelectionItem `json:"items"`
}

// SearchSelectionItemsRequest درخواست جستجوی آیتم‌های انتخابی
type SearchSelectionItemsRequest struct {
	SelectionID string `json:"selection_id"`
	Query       string `json:"query"`
}

// SearchSelectionItemsResponse پاسخ جستجوی آیتم‌های انتخابی
type SearchSelectionItemsResponse struct {
	Items []ButtonSelectionItem `json:"items"`
}

// ReceiveUpdateRequest درخواست دریافت آپدیت‌ها
type ReceiveUpdateRequest struct {
	Update *Update `json:"update"`
}

// ReceiveInlineMessageRequest درخواست دریافت پیام اینلاین
type ReceiveInlineMessageRequest struct {
	InlineMessage *InlineMessage `json:"inline_message"`
}

// ReceiveQueryRequest درخواست دریافت کوئری
type ReceiveQueryRequest struct {
	// نکته: این فیلد در وب‌هوک به صورت string برگردانده می‌شود.
	QueryID  string `json:"query_id"`
	ChatID   string `json:"chat_id"`
	SenderID string `json:"sender_id"`
	Data     string `json:"data"`
}
