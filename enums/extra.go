package enums

// VoiceChatState is the speaking state of a voice-chat participant.
type VoiceChatState string

const (
	VoiceChatStateMute               VoiceChatState = "Mute"
	VoiceChatStateUnmute             VoiceChatState = "Unmute"
	VoiceChatStateRequestSpeak       VoiceChatState = "RequestSpeak"
	VoiceChatStateCancelRequestSpeak VoiceChatState = "CancelRequestSpeak"
)

// String implements fmt.Stringer.
func (v VoiceChatState) String() string { return string(v) }

// VoiceChatActivity is the activity indicator broadcast in voice chats.
type VoiceChatActivity string

const (
	VoiceChatActivitySpeaking VoiceChatActivity = "Speaking"
)

// String implements fmt.Stringer.
func (v VoiceChatActivity) String() string { return string(v) }

// GroupAdminAccess enumerates the permission keys accepted by the
// setGroupAdmin API. Use these when building the access_list argument.
type GroupAdminAccess string

const (
	GroupAdminAccessSetAdmin                GroupAdminAccess = "SetAdmin"
	GroupAdminAccessBanMember               GroupAdminAccess = "BanMember"
	GroupAdminAccessChangeInfo              GroupAdminAccess = "ChangeInfo"
	GroupAdminAccessPinMessages             GroupAdminAccess = "PinMessages"
	GroupAdminAccessSetJoinLink             GroupAdminAccess = "SetJoinLink"
	GroupAdminAccessSetMemberAccess         GroupAdminAccess = "SetMemberAccess"
	GroupAdminAccessDeleteGlobalAllMessages GroupAdminAccess = "DeleteGlobalAllMessages"
	GroupAdminAccessDeleteMessages          GroupAdminAccess = "DeleteMessages"
	GroupAdminAccessEditMessages            GroupAdminAccess = "EditMessages"
)

// ChannelAdminAccess enumerates the permission keys accepted by the
// setChannelAdmin API.
type ChannelAdminAccess string

const (
	ChannelAdminAccessSetAdmin                ChannelAdminAccess = "SetAdmin"
	ChannelAdminAccessAddMember               ChannelAdminAccess = "AddMember"
	ChannelAdminAccessViewAdmins              ChannelAdminAccess = "ViewAdmins"
	ChannelAdminAccessChangeInfo              ChannelAdminAccess = "ChangeInfo"
	ChannelAdminAccessPinMessages             ChannelAdminAccess = "PinMessages"
	ChannelAdminAccessViewMembers             ChannelAdminAccess = "ViewMembers"
	ChannelAdminAccessSetJoinLink             ChannelAdminAccess = "SetJoinLink"
	ChannelAdminAccessSendMessage             ChannelAdminAccess = "SendMessages"
	ChannelAdminAccessEditMessage             ChannelAdminAccess = "EditMessage"
	ChannelAdminAccessDeleteMessage           ChannelAdminAccess = "DeleteMessage"
	ChannelAdminAccessDeleteGlobalAllMessages ChannelAdminAccess = "DeleteGlobalAllMessages"
)

// GroupDefaultAccess enumerates the default-permission keys for groups.
type GroupDefaultAccess string

const (
	GroupDefaultAccessSendMessages GroupDefaultAccess = "SendMessages"
	GroupDefaultAccessAddMember    GroupDefaultAccess = "AddMember"
	GroupDefaultAccessViewAdmins   GroupDefaultAccess = "ViewAdmins"
	GroupDefaultAccessViewMembers  GroupDefaultAccess = "ViewMembers"
)

// CallDiscardReason is the reason a call was ended.
type CallDiscardReason string

const (
	CallDiscardReasonMissed     CallDiscardReason = "Missed"
	CallDiscardReasonDisconnect CallDiscardReason = "Disconnect"
)

// String implements fmt.Stringer.
func (c CallDiscardReason) String() string { return string(c) }

// CallType distinguishes voice from video calls.
type CallType string

const (
	CallTypeVoice CallType = "Voice"
	CallTypeVideo CallType = "Video"
)

// AskSpamAction is the action taken on an ask-spam prompt.
type AskSpamAction string

const (
	AskSpamCancel         AskSpamAction = "Cancel"
	AskSpamBlockUser      AskSpamAction = "BlockUser"
	AskSpamAddToContact   AskSpamAction = "AddToContact"
	AskSpamReportAndLeave AskSpamAction = "ReportAndLeave"
)

// DeleteMessageType is the scope of a delete operation.
type DeleteMessageType string

const (
	DeleteMessageLocal     DeleteMessageType = "Local"
	DeleteMessageGlobal    DeleteMessageType = "Global"
	DeleteMessageScheduled DeleteMessageType = "Scheduled"
)

// ChatAdsAction is the action applied to a chat advertisement.
type ChatAdsAction string

const (
	ChatAdsActionView     ChatAdsAction = "View"
	ChatAdsActionClick    ChatAdsAction = "Click"
	ChatAdsActionDontShow ChatAdsAction = "DontShow"
	ChatAdsActionReport   ChatAdsAction = "Report"
)

// JoinRequestAction is the action applied to a pending join request.
type JoinRequestAction string

const (
	JoinRequestAccept JoinRequestAction = "Accept"
	JoinRequestReject JoinRequestAction = "Reject"
)

// OwnershipAction is the action applied to an ownership-transfer request.
type OwnershipAction string

const (
	OwnershipAccept OwnershipAction = "Accept"
	OwnershipReject OwnershipAction = "Reject"
)

// StickerSetAction is the action applied to a sticker set.
type StickerSetAction string

const (
	StickerSetAdd     StickerSetAction = "Add"
	StickerSetRemove  StickerSetAction = "Remove"
	StickerSetArchive StickerSetAction = "Archive"
	StickerSetDelete  StickerSetAction = "Delete"
)

// ChannelType is the visibility/type of a channel.
type ChannelType string

const (
	ChannelTypePublic  ChannelType = "Public"
	ChannelTypePrivate ChannelType = "Private"
)

// PollAction is the action applied to a poll.
type PollAction string

const (
	PollActionStop    PollAction = "Stop"
	PollActionRetract PollAction = "Retract"
)

// ReactionAction is the add/remove action for message reactions.
type ReactionAction string

const (
	ReactionAdd    ReactionAction = "Add"
	ReactionRemove ReactionAction = "Remove"
)

// RubinoStoryType is how a Rubino story is forwarded.
type RubinoStoryType string

const (
	RubinoStoryReply  RubinoStoryType = "Reply"
	RubinoStoryDirect RubinoStoryType = "Direct"
)

// SortOrder is the pagination direction used by Rubino/TV endpoints.
type SortOrder string

const (
	SortFromMax SortOrder = "FromMax"
	SortFromMin SortOrder = "FromMin"
)

// SearchType is the type of message search.
type SearchType string

const (
	SearchText    SearchType = "Text"
	SearchHashtag SearchType = "Hashtag"
)

// SendType is the channel for delivering a login code.
type SendType string

const (
	SendTypeSMS      SendType = "SMS"
	SendTypeInternal SendType = "Internal"
)
