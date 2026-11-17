package models

// Bot is the bot identity returned by `getMe`.
type Bot struct {
	BotID        string `json:"bot_id"`
	BotTitle     string `json:"bot_title"`
	Avatar       *File  `json:"avatar,omitempty"`
	Description  string `json:"description,omitempty"`
	Username     string `json:"username,omitempty"`
	StartMessage string `json:"start_message,omitempty"`
	ShareURL     string `json:"share_url,omitempty"`
}

// BotCommand is a single entry in the bot's command list.
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// UpdateType enumerates the update types returned by `getUpdates`.
type UpdateType string

const (
	UpdateTypeNewMessage     UpdateType = "NewMessage"
	UpdateTypeUpdatedMessage UpdateType = "UpdatedMessage"
	UpdateTypeRemovedMessage UpdateType = "RemovedMessage"
	UpdateTypeStartedBot     UpdateType = "StartedBot"
	UpdateTypeStoppedBot     UpdateType = "StoppedBot"
	UpdateTypeUpdatedPayment UpdateType = "UpdatedPayment"
	UpdateTypeInlineMessage  UpdateType = "InlineMessage"
)

// ChatKeypadType selects how a chat keypad is applied (None, New, Remove).
type ChatKeypadType string

const (
	ChatKeypadNone   ChatKeypadType = "None"
	ChatKeypadNew    ChatKeypadType = "New"
	ChatKeypadRemove ChatKeypadType = "Remove"
)

// FileType enumerates the supported file upload types.
type FileType string

const (
	FileTypeFile  FileType = "File"
	FileTypeImage FileType = "Image"
	FileTypeVoice FileType = "Voice"
	FileTypeMusic FileType = "Music"
	FileTypeGif   FileType = "Gif"
	FileTypeVideo FileType = "Video"
)

// ButtonType enumerates the supported button types.
type ButtonType string

const (
	ButtonTypeSimple       ButtonType = "Simple"
	ButtonTypeSelection    ButtonType = "Selection"
	ButtonTypeCalendar     ButtonType = "Calendar"
	ButtonTypeNumberPicker ButtonType = "NumberPicker"
	ButtonTypeStringPicker ButtonType = "StringPicker"
	ButtonTypeLocation     ButtonType = "Location"
	ButtonTypePayment      ButtonType = "Payment"
	ButtonTypeCameraImage  ButtonType = "CameraImage"
	ButtonTypeCameraVideo  ButtonType = "CameraVideo"
	ButtonTypeGalleryImage ButtonType = "GalleryImage"
	ButtonTypeGalleryVideo ButtonType = "GalleryVideo"
	ButtonTypeFile         ButtonType = "File"
	ButtonTypeAudio        ButtonType = "Audio"
	ButtonTypeRecordAudio  ButtonType = "RecordAudio"
	ButtonTypeMyPhone      ButtonType = "MyPhoneNumber"
	ButtonTypeMyLocation   ButtonType = "MyLocation"
	ButtonTypeTextbox      ButtonType = "Textbox"
	ButtonTypeLink         ButtonType = "Link"
	ButtonTypeAskMyPhone   ButtonType = "AskMyPhoneNumber"
	ButtonTypeAskLocation  ButtonType = "AskLocation"
	ButtonTypeBarcode      ButtonType = "Barcode"
)

// UpdateEndpointType names the webhook endpoint kinds accepted by
// `updateBotEndpoints`.
type UpdateEndpointType string

const (
	EndpointReceiveUpdate        UpdateEndpointType = "ReceiveUpdate"
	EndpointReceiveInlineMessage UpdateEndpointType = "ReceiveInlineMessage"
	EndpointReceiveQuery         UpdateEndpointType = "ReceiveQuery"
	EndpointGetSelectionItem     UpdateEndpointType = "GetSelectionItem"
	EndpointSearchSelectionItems UpdateEndpointType = "SearchSelectionItems"
)
