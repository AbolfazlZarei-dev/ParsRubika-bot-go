# `enums` — Extra

_Source: [`extra.go`](../../ParsRubika/enums/extra.go)_

## Symbols

### `VoiceChatState`

```go
type VoiceChatState string
```

VoiceChatState is the speaking state of a voice-chat participant.

### `VoiceChatStateMute`

```go
const VoiceChatStateMute VoiceChatState = "Mute"
```

### `VoiceChatStateUnmute`

```go
const VoiceChatStateUnmute VoiceChatState = "Unmute"
```

### `VoiceChatStateRequestSpeak`

```go
const VoiceChatStateRequestSpeak VoiceChatState = "RequestSpeak"
```

### `VoiceChatStateCancelRequestSpeak`

```go
const VoiceChatStateCancelRequestSpeak VoiceChatState = "CancelRequestSpeak"
```

### `String`

```go
func (v VoiceChatState) String() string { return string(v) }
```

String implements fmt.Stringer.

### `VoiceChatActivity`

```go
type VoiceChatActivity string
```

VoiceChatActivity is the activity indicator broadcast in voice chats.

### `VoiceChatActivitySpeaking`

```go
const VoiceChatActivitySpeaking VoiceChatActivity = "Speaking"
```

### `String`

```go
func (v VoiceChatActivity) String() string { return string(v) }
```

String implements fmt.Stringer.

### `GroupAdminAccess`

```go
type GroupAdminAccess string
```

GroupAdminAccess enumerates the permission keys accepted by the
setGroupAdmin API. Use these when building the access_list argument.

### `GroupAdminAccessSetAdmin`

```go
const GroupAdminAccessSetAdmin GroupAdminAccess = "SetAdmin"
```

### `GroupAdminAccessBanMember`

```go
const GroupAdminAccessBanMember GroupAdminAccess = "BanMember"
```

### `GroupAdminAccessChangeInfo`

```go
const GroupAdminAccessChangeInfo GroupAdminAccess = "ChangeInfo"
```

### `GroupAdminAccessPinMessages`

```go
const GroupAdminAccessPinMessages GroupAdminAccess = "PinMessages"
```

### `GroupAdminAccessSetJoinLink`

```go
const GroupAdminAccessSetJoinLink GroupAdminAccess = "SetJoinLink"
```

### `GroupAdminAccessSetMemberAccess`

```go
const GroupAdminAccessSetMemberAccess GroupAdminAccess = "SetMemberAccess"
```

### `GroupAdminAccessDeleteGlobalAllMessages`

```go
const GroupAdminAccessDeleteGlobalAllMessages GroupAdminAccess = "DeleteGlobalAllMessages"
```

### `GroupAdminAccessDeleteMessages`

```go
const GroupAdminAccessDeleteMessages GroupAdminAccess = "DeleteMessages"
```

### `GroupAdminAccessEditMessages`

```go
const GroupAdminAccessEditMessages GroupAdminAccess = "EditMessages"
```

### `ChannelAdminAccess`

```go
type ChannelAdminAccess string
```

ChannelAdminAccess enumerates the permission keys accepted by the
setChannelAdmin API.

### `ChannelAdminAccessSetAdmin`

```go
const ChannelAdminAccessSetAdmin ChannelAdminAccess = "SetAdmin"
```

### `ChannelAdminAccessAddMember`

```go
const ChannelAdminAccessAddMember ChannelAdminAccess = "AddMember"
```

### `ChannelAdminAccessViewAdmins`

```go
const ChannelAdminAccessViewAdmins ChannelAdminAccess = "ViewAdmins"
```

### `ChannelAdminAccessChangeInfo`

```go
const ChannelAdminAccessChangeInfo ChannelAdminAccess = "ChangeInfo"
```

### `ChannelAdminAccessPinMessages`

```go
const ChannelAdminAccessPinMessages ChannelAdminAccess = "PinMessages"
```

### `ChannelAdminAccessViewMembers`

```go
const ChannelAdminAccessViewMembers ChannelAdminAccess = "ViewMembers"
```

### `ChannelAdminAccessSetJoinLink`

```go
const ChannelAdminAccessSetJoinLink ChannelAdminAccess = "SetJoinLink"
```

### `ChannelAdminAccessSendMessage`

```go
const ChannelAdminAccessSendMessage ChannelAdminAccess = "SendMessages"
```

### `ChannelAdminAccessEditMessage`

```go
const ChannelAdminAccessEditMessage ChannelAdminAccess = "EditMessage"
```

### `ChannelAdminAccessDeleteMessage`

```go
const ChannelAdminAccessDeleteMessage ChannelAdminAccess = "DeleteMessage"
```

### `ChannelAdminAccessDeleteGlobalAllMessages`

```go
const ChannelAdminAccessDeleteGlobalAllMessages ChannelAdminAccess = "DeleteGlobalAllMessages"
```

### `GroupDefaultAccess`

```go
type GroupDefaultAccess string
```

GroupDefaultAccess enumerates the default-permission keys for groups.

### `GroupDefaultAccessSendMessages`

```go
const GroupDefaultAccessSendMessages GroupDefaultAccess = "SendMessages"
```

### `GroupDefaultAccessAddMember`

```go
const GroupDefaultAccessAddMember GroupDefaultAccess = "AddMember"
```

### `GroupDefaultAccessViewAdmins`

```go
const GroupDefaultAccessViewAdmins GroupDefaultAccess = "ViewAdmins"
```

### `GroupDefaultAccessViewMembers`

```go
const GroupDefaultAccessViewMembers GroupDefaultAccess = "ViewMembers"
```

### `CallDiscardReason`

```go
type CallDiscardReason string
```

CallDiscardReason is the reason a call was ended.

### `CallDiscardReasonMissed`

```go
const CallDiscardReasonMissed CallDiscardReason = "Missed"
```

### `CallDiscardReasonDisconnect`

```go
const CallDiscardReasonDisconnect CallDiscardReason = "Disconnect"
```

### `String`

```go
func (c CallDiscardReason) String() string { return string(c) }
```

String implements fmt.Stringer.

### `CallType`

```go
type CallType string
```

CallType distinguishes voice from video calls.

### `CallTypeVoice`

```go
const CallTypeVoice CallType = "Voice"
```

### `CallTypeVideo`

```go
const CallTypeVideo CallType = "Video"
```

### `AskSpamAction`

```go
type AskSpamAction string
```

AskSpamAction is the action taken on an ask-spam prompt.

### `AskSpamCancel`

```go
const AskSpamCancel AskSpamAction = "Cancel"
```

### `AskSpamBlockUser`

```go
const AskSpamBlockUser AskSpamAction = "BlockUser"
```

### `AskSpamAddToContact`

```go
const AskSpamAddToContact AskSpamAction = "AddToContact"
```

### `AskSpamReportAndLeave`

```go
const AskSpamReportAndLeave AskSpamAction = "ReportAndLeave"
```

### `DeleteMessageType`

```go
type DeleteMessageType string
```

DeleteMessageType is the scope of a delete operation.

### `DeleteMessageLocal`

```go
const DeleteMessageLocal DeleteMessageType = "Local"
```

### `DeleteMessageGlobal`

```go
const DeleteMessageGlobal DeleteMessageType = "Global"
```

### `DeleteMessageScheduled`

```go
const DeleteMessageScheduled DeleteMessageType = "Scheduled"
```

### `ChatAdsAction`

```go
type ChatAdsAction string
```

ChatAdsAction is the action applied to a chat advertisement.

### `ChatAdsActionView`

```go
const ChatAdsActionView ChatAdsAction = "View"
```

### `ChatAdsActionClick`

```go
const ChatAdsActionClick ChatAdsAction = "Click"
```

### `ChatAdsActionDontShow`

```go
const ChatAdsActionDontShow ChatAdsAction = "DontShow"
```

### `ChatAdsActionReport`

```go
const ChatAdsActionReport ChatAdsAction = "Report"
```

### `JoinRequestAction`

```go
type JoinRequestAction string
```

JoinRequestAction is the action applied to a pending join request.

### `JoinRequestAccept`

```go
const JoinRequestAccept JoinRequestAction = "Accept"
```

### `JoinRequestReject`

```go
const JoinRequestReject JoinRequestAction = "Reject"
```

### `OwnershipAction`

```go
type OwnershipAction string
```

OwnershipAction is the action applied to an ownership-transfer request.

### `OwnershipAccept`

```go
const OwnershipAccept OwnershipAction = "Accept"
```

### `OwnershipReject`

```go
const OwnershipReject OwnershipAction = "Reject"
```

### `StickerSetAction`

```go
type StickerSetAction string
```

StickerSetAction is the action applied to a sticker set.

### `StickerSetAdd`

```go
const StickerSetAdd StickerSetAction = "Add"
```

### `StickerSetRemove`

```go
const StickerSetRemove StickerSetAction = "Remove"
```

### `StickerSetArchive`

```go
const StickerSetArchive StickerSetAction = "Archive"
```

### `StickerSetDelete`

```go
const StickerSetDelete StickerSetAction = "Delete"
```

### `ChannelType`

```go
type ChannelType string
```

ChannelType is the visibility/type of a channel.

### `ChannelTypePublic`

```go
const ChannelTypePublic ChannelType = "Public"
```

### `ChannelTypePrivate`

```go
const ChannelTypePrivate ChannelType = "Private"
```

### `PollAction`

```go
type PollAction string
```

PollAction is the action applied to a poll.

### `PollActionStop`

```go
const PollActionStop PollAction = "Stop"
```

### `PollActionRetract`

```go
const PollActionRetract PollAction = "Retract"
```

### `ReactionAction`

```go
type ReactionAction string
```

ReactionAction is the add/remove action for message reactions.

### `ReactionAdd`

```go
const ReactionAdd ReactionAction = "Add"
```

### `ReactionRemove`

```go
const ReactionRemove ReactionAction = "Remove"
```

### `RubinoStoryType`

```go
type RubinoStoryType string
```

RubinoStoryType is how a Rubino story is forwarded.

### `RubinoStoryReply`

```go
const RubinoStoryReply RubinoStoryType = "Reply"
```

### `RubinoStoryDirect`

```go
const RubinoStoryDirect RubinoStoryType = "Direct"
```

### `SortOrder`

```go
type SortOrder string
```

SortOrder is the pagination direction used by Rubino/TV endpoints.

### `SortFromMax`

```go
const SortFromMax SortOrder = "FromMax"
```

### `SortFromMin`

```go
const SortFromMin SortOrder = "FromMin"
```

### `SearchType`

```go
type SearchType string
```

SearchType is the type of message search.

### `SearchText`

```go
const SearchText SearchType = "Text"
```

### `SearchHashtag`

```go
const SearchHashtag SearchType = "Hashtag"
```

### `SendType`

```go
type SendType string
```

SendType is the channel for delivering a login code.

### `SendTypeSMS`

```go
const SendTypeSMS SendType = "SMS"
```

### `SendTypeInternal`

```go
const SendTypeInternal SendType = "Internal"
```
