package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

// ==================== enumهای اصلی ====================

type ChatTypeEnum string

const (
	UserChat    ChatTypeEnum = "User"
	BotChat     ChatTypeEnum = "Bot"
	GroupChat   ChatTypeEnum = "Group"
	ChannelChat ChatTypeEnum = "Channel"
)

type FileTypeEnum string

const (
	FileType    FileTypeEnum = "File"
	ImageType   FileTypeEnum = "Image"
	VoiceType   FileTypeEnum = "Voice"
	VideoType   FileTypeEnum = "Video"
	MusicType   FileTypeEnum = "Music"
	GifType     FileTypeEnum = "Gif"
	PhotoType   FileTypeEnum = "Photo"
	StickerType FileTypeEnum = "Sticker"
	AudioType   FileTypeEnum = "Audio"
)

type UpdateTypeEnum string

const (
	UpdatedMessage UpdateTypeEnum = "UpdatedMessage"
	NewMessage     UpdateTypeEnum = "NewMessage"
	RemovedMessage UpdateTypeEnum = "RemovedMessage"
	StartedBot     UpdateTypeEnum = "StartedBot"
	StoppedBot     UpdateTypeEnum = "StoppedBot"
	UpdatedPayment UpdateTypeEnum = "UpdatedPayment"
	CallbackQuery  UpdateTypeEnum = "CallbackQuery"
	InlineQuery    UpdateTypeEnum = "InlineQuery"
)

type ButtonTypeEnum string

const (
	ButtonTypeSimple           ButtonTypeEnum = "Simple"
	ButtonTypeSelection        ButtonTypeEnum = "Selection"
	ButtonTypeCalendar         ButtonTypeEnum = "Calendar"
	ButtonTypeNumberPicker     ButtonTypeEnum = "NumberPicker"
	ButtonTypeStringPicker     ButtonTypeEnum = "StringPicker"
	ButtonTypeLocation         ButtonTypeEnum = "Location"
	ButtonTypePayment          ButtonTypeEnum = "Payment"
	ButtonTypeCameraImage      ButtonTypeEnum = "CameraImage"
	ButtonTypeCameraVideo      ButtonTypeEnum = "CameraVideo"
	ButtonTypeGalleryImage     ButtonTypeEnum = "GalleryImage"
	ButtonTypeGalleryVideo     ButtonTypeEnum = "GalleryVideo"
	ButtonTypeFile             ButtonTypeEnum = "File"
	ButtonTypeAudio            ButtonTypeEnum = "Audio"
	ButtonTypeRecordAudio      ButtonTypeEnum = "RecordAudio"
	ButtonTypeMyPhoneNumber    ButtonTypeEnum = "MyPhoneNumber"
	ButtonTypeMyLocation       ButtonTypeEnum = "MyLocation"
	ButtonTypeTextbox          ButtonTypeEnum = "Textbox"
	ButtonTypeLink             ButtonTypeEnum = "Link"
	ButtonTypeAskMyPhoneNumber ButtonTypeEnum = "AskMyPhoneNumber"
	ButtonTypeAskMyLocation    ButtonTypeEnum = "AskMyLocation"
	ButtonTypeBarcode          ButtonTypeEnum = "Barcode"
	ButtonTypeWebApp           ButtonTypeEnum = "WebApp"
	ButtonTypeLogin            ButtonTypeEnum = "Login"
	ButtonTypeGame             ButtonTypeEnum = "Game"
)

type ChatKeypadTypeEnum string

const (
	NoneKeypad   ChatKeypadTypeEnum = "None"
	NewKeypad    ChatKeypadTypeEnum = "New"
	RemoveKeypad ChatKeypadTypeEnum = "Remove"
)

type UpdateEndpointTypeEnum string

const (
	ReceiveUpdate        UpdateEndpointTypeEnum = "ReceiveUpdate"
	ReceiveInlineMessage UpdateEndpointTypeEnum = "ReceiveInlineMessage"
	ReceiveQuery         UpdateEndpointTypeEnum = "ReceiveQuery"
	GetSelectionItem     UpdateEndpointTypeEnum = "GetSelectionItem"
	SearchSelectionItems UpdateEndpointTypeEnum = "SearchSelectionItems"
)

type MessageSenderEnum string

const (
	UserSender MessageSenderEnum = "User"
	BotSender  MessageSenderEnum = "Bot"
)

type ChatMemberStatusEnum string

const (
	MemberCreator       ChatMemberStatusEnum = "Creator"
	MemberAdministrator ChatMemberStatusEnum = "Administrator"
	MemberMember        ChatMemberStatusEnum = "Member"
	MemberRestricted    ChatMemberStatusEnum = "Restricted"
	MemberLeft          ChatMemberStatusEnum = "Left"
	MemberKicked        ChatMemberStatusEnum = "Kicked"
)

type PollStatusEnum string

const (
	OpenPoll   PollStatusEnum = "Open"
	ClosedPoll PollStatusEnum = "Closed"
)

// ==================== enumهای جدید (نسخه 2.0.0) ====================

// APIType نوع API برای اتصال
type APIType string

const (
	BotAPI       APIType = "bot_api"       // API اصلی روبیکا
	MessengerAPI APIType = "messenger_api" // API جدید مسنجر
)

// ConnectionMode حالت اتصال ربات
type ConnectionMode string

const (
	BotAPIMode       ConnectionMode = "bot_api_only"       // فقط از طریق Bot API
	MessengerAPIMode ConnectionMode = "messenger_api_only" // فقط از طریق API Messenger
	SwitcherMode     ConnectionMode = "switcher"           // حالت هوشمند سوییچر
)

type ForwardedFromEnum string

const (
	ForwardedFromUser    ForwardedFromEnum = "User"
	ForwardedFromChannel ForwardedFromEnum = "Channel"
	ForwardedFromBot     ForwardedFromEnum = "Bot"
)

type PaymentStatusEnum string

const (
	PaymentPaid    PaymentStatusEnum = "Paid"
	PaymentNotPaid PaymentStatusEnum = "NotPaid"
)

type ButtonSelectionTypeEnum string

const (
	ButtonSelectionTextOnly   ButtonSelectionTypeEnum = "TextOnly"
	ButtonSelectionTextImgThu ButtonSelectionTypeEnum = "TextImgThu"
	ButtonSelectionTextImgBig ButtonSelectionTypeEnum = "TextImgBig"
)

type ButtonSelectionSearchEnum string

const (
	ButtonSelectionSearchNone  ButtonSelectionSearchEnum = "None"
	ButtonSelectionSearchLocal ButtonSelectionSearchEnum = "Local"
	ButtonSelectionSearchApi   ButtonSelectionSearchEnum = "Api"
)

type ButtonSelectionGetEnum string

const (
	ButtonSelectionGetLocal ButtonSelectionGetEnum = "Local"
	ButtonSelectionGetApi   ButtonSelectionGetEnum = "Api"
)

type ButtonCalendarTypeEnum string

const (
	ButtonCalendarDatePersian   ButtonCalendarTypeEnum = "DatePersian"
	ButtonCalendarDateGregorian ButtonCalendarTypeEnum = "DateGregorian"
)

type ButtonTextboxTypeLineEnum string

const (
	ButtonTextboxLineSingle ButtonTextboxTypeLineEnum = "SingleLine"
	ButtonTextboxLineMulti  ButtonTextboxTypeLineEnum = "MultiLine"
)

type ButtonTextboxTypeKeypadEnum string

const (
	ButtonTextboxKeypadString ButtonTextboxTypeKeypadEnum = "String"
	ButtonTextboxKeypadNumber ButtonTextboxTypeKeypadEnum = "Number"
)

type ButtonLocationTypeEnum string

const (
	ButtonLocationPicker ButtonLocationTypeEnum = "Picker"
	ButtonLocationView   ButtonLocationTypeEnum = "View"
)

// ==================== enumهای جدید برای فیلترها ====================

type FilterTypeEnum string

const (
	FilterTypeText     FilterTypeEnum = "Text"
	FilterTypePhoto    FilterTypeEnum = "Photo"
	FilterTypeAudio    FilterTypeEnum = "Audio"
	FilterTypeVideo    FilterTypeEnum = "Video"
	FilterTypeDocument FilterTypeEnum = "Document"
	FilterTypeSticker  FilterTypeEnum = "Sticker"
	FilterTypeLocation FilterTypeEnum = "Location"
	FilterTypeContact  FilterTypeEnum = "Contact"
	FilterTypePoll     FilterTypeEnum = "Poll"
	FilterTypeCommand  FilterTypeEnum = "Command"
	FilterTypeCallback FilterTypeEnum = "Callback"
	FilterTypeInline   FilterTypeEnum = "Inline"
	FilterTypeEdited   FilterTypeEnum = "Edited"
	FilterTypeForward  FilterTypeEnum = "Forward"
	FilterTypeReply    FilterTypeEnum = "Reply"
	FilterTypeGroup    FilterTypeEnum = "Group"
	FilterTypeChannel  FilterTypeEnum = "Channel"
	FilterTypePrivate  FilterTypeEnum = "Private"
	FilterTypeBot      FilterTypeEnum = "Bot"
	FilterTypeUser     FilterTypeEnum = "User"
)

// ==================== enumهای جدید برای فعالیت چت ====================

type ChatActivityEnum string

const (
	ChatActivityTyping          ChatActivityEnum = "typing"
	ChatActivityUploadPhoto     ChatActivityEnum = "upload_photo"
	ChatActivityRecordVideo     ChatActivityEnum = "record_video"
	ChatActivityUploadVideo     ChatActivityEnum = "upload_video"
	ChatActivityRecordAudio     ChatActivityEnum = "record_audio"
	ChatActivityUploadAudio     ChatActivityEnum = "upload_audio"
	ChatActivityUploadDocument  ChatActivityEnum = "upload_document"
	ChatActivityFindLocation    ChatActivityEnum = "find_location"
	ChatActivityRecordVideoNote ChatActivityEnum = "record_video_note"
	ChatActivityUploadVideoNote ChatActivityEnum = "upload_video_note"
	ChatActivityChooseSticker   ChatActivityEnum = "choose_sticker"
)
