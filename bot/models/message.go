package models

// ForwardedFrom carries the forward metadata.
type ForwardedFrom struct {
	TypeFrom     string `json:"type_from"`
	MessageID    string `json:"message_id,omitempty"`
	FromChatID   string `json:"from_chat_id,omitempty"`
	FromSenderID string `json:"from_sender_id,omitempty"`
}

// PaymentStatus reflects the payment state attached to a message.
type PaymentStatus struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}

// PollStatus holds the live state of a poll.
type PollStatus struct {
	State              string `json:"state,omitempty"`
	SelectionIndex     int    `json:"selection_index,omitempty"`
	PercentVoteOptions []int  `json:"percent_vote_options,omitempty"`
	TotalVote          int    `json:"total_vote,omitempty"`
	ShowTotalVotes     *bool  `json:"show_total_votes,omitempty"`
}

// Poll is a poll attachment.
type Poll struct {
	Question   string      `json:"question"`
	Options    []string    `json:"options"`
	PollStatus *PollStatus `json:"poll_status,omitempty"`
}

// Location is a (lat, lon) pair.
type Location struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

// LiveLocation carries a live-location update.
type LiveLocation struct {
	StartTime       string    `json:"start_time"`
	LivePeriod      int       `json:"live_period"`
	CurrentLocation *Location `json:"current_location,omitempty"`
	UserID          string    `json:"user_id,omitempty"`
	Status          string    `json:"status,omitempty"`
	LastUpdateTime  string    `json:"last_update_time,omitempty"`
}

// ContactMessage carries a contact vCard payload.
type ContactMessage struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
}

// Sticker is the sticker metadata.
type Sticker struct {
	StickerID      string `json:"sticker_id"`
	File           *File  `json:"file,omitempty"`
	EmojiCharacter string `json:"emoji_character,omitempty"`
}

// AuxData carries the aux_data attached to a button click.
type AuxData struct {
	StartID  string `json:"start_id,omitempty"`
	ButtonID string `json:"button_id,omitempty"`
}

// JoinChannelData is the payload attached to a "join channel" button.
type JoinChannelData struct {
	Username string `json:"username"`
	AskJoin  bool   `json:"ask_join,omitempty"`
}

// OpenChatData carries the open-chat payload attached to a button.
type OpenChatData struct {
	ObjectGUID string `json:"object_guid"`
	ObjectType string `json:"object_type,omitempty"`
}

// Message is the parsed form of a `new_message` / `updated_message` payload.
type Message struct {
	MessageID        string          `json:"message_id"`
	Time             string          `json:"time,omitempty"`
	Text             string          `json:"text,omitempty"`
	IsEdited         bool            `json:"is_edited,omitempty"`
	SenderType       string          `json:"sender_type,omitempty"`
	SenderID         string          `json:"sender_id,omitempty"`
	AuxData          *AuxData        `json:"aux_data,omitempty"`
	File             *File           `json:"file,omitempty"`
	ReplyToMessageID string          `json:"reply_to_message_id,omitempty"`
	ForwardedFrom    *ForwardedFrom  `json:"forwarded_from,omitempty"`
	ForwardedNoLink  string          `json:"forwarded_no_link,omitempty"`
	Location         *Location       `json:"location,omitempty"`
	Sticker          *Sticker        `json:"sticker,omitempty"`
	ContactMessage   *ContactMessage `json:"contact_message,omitempty"`
	Poll             *Poll           `json:"poll,omitempty"`
	LiveLocation     *LiveLocation   `json:"live_location,omitempty"`
	Metadata         *Metadata       `json:"metadata,omitempty"`
}

// MessageID is the response of every send* call.
type MessageID struct {
	MessageID    string `json:"message_id,omitempty"`
	NewMessageID string `json:"new_message_id,omitempty"`
	FileID       string `json:"file_id,omitempty"`
	ChatID       string `json:"chat_id,omitempty"`
}

// MetaDataParts is one slice of the metadata attached to a text message.
type MetaDataParts struct {
	FromIndex         int    `json:"from_index,omitempty"`
	Length            int    `json:"length,omitempty"`
	Type              string `json:"type,omitempty"`
	LinkURL           string `json:"link_url,omitempty"`
	MentionTextUserID string `json:"mention_text_user_id,omitempty"`
}

// Metadata is the structured-form metadata block.
type Metadata struct {
	MetaDataParts []MetaDataParts `json:"meta_data_parts,omitempty"`
}
