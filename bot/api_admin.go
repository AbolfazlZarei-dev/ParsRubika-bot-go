package bot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// SendQuiz is a convenience wrapper around SendPoll that sets type=Quiz
// and requires a correct option index.
func (b *Bot) SendQuiz(ctx context.Context, chatID, question string, options []string, correctOptionIdx int, explanation string, replyTo string, disableNotif bool, keypad *models.Keypad, keypadType models.ChatKeypadType) (*models.MessageID, error) {
	idx := correctOptionIdx
	return b.SendPoll(ctx, chatID, question, options, "Quiz", false, true, &idx, explanation, replyTo, disableNotif, keypad, keypadType)
}

// ClosePoll closes/stops a poll so no further votes are accepted. The
// Rubika bot API exposes this as a dedicated endpoint in newer versions;
// when not available, use the account-side setPollAction.
func (b *Bot) ClosePoll(ctx context.Context, chatID, messageID string) error {
	_, err := b.makeRequest(ctx, "closePoll", map[string]any{
		"chat_id":    chatID,
		"message_id": messageID,
	})
	return err
}

// GetChatAdmins is an alias for GetChatAdministrators (common naming).
func (b *Bot) GetChatAdmins(ctx context.Context, chatID string) ([]map[string]any, error) {
	return b.GetChatAdministrators(ctx, chatID)
}

// GetChatMembers returns a page of chat members. startID is the
// pagination cursor (empty for the first page).
func (b *Bot) GetChatMembers(ctx context.Context, chatID, startID string) (map[string]any, error) {
	payload := map[string]any{"chat_id": chatID}
	if startID != "" {
		payload["start_id"] = startID
	}
	return b.makeRequest(ctx, "getChatMembers", payload)
}

// GetChatInfo is an alias for GetChat.
func (b *Bot) GetChatInfo(ctx context.Context, chatID string) (*models.Chat, error) {
	return b.GetChat(ctx, chatID)
}

// SetChatTitle updates a group/channel title.
func (b *Bot) SetChatTitle(ctx context.Context, chatID, title string) error {
	_, err := b.makeRequest(ctx, "editChatTitle", map[string]any{
		"chat_id": chatID,
		"title":   title,
	})
	return err
}

// SetChatDescription updates a group/channel description.
func (b *Bot) SetChatDescription(ctx context.Context, chatID, description string) error {
	_, err := b.makeRequest(ctx, "editChatDescription", map[string]any{
		"chat_id":     chatID,
		"description": description,
	})
	return err
}

// SetChatPhoto sets the chat avatar from an uploaded file_id.
func (b *Bot) SetChatPhoto(ctx context.Context, chatID, fileID string) error {
	_, err := b.makeRequest(ctx, "editChatPhoto", map[string]any{
		"chat_id": chatID,
		"file_id": fileID,
	})
	return err
}

// RemoveChatPhoto removes the chat avatar.
func (b *Bot) RemoveChatPhoto(ctx context.Context, chatID string) error {
	_, err := b.makeRequest(ctx, "editChatPhoto", map[string]any{
		"chat_id": chatID,
		"file_id": "Remove",
	})
	return err
}

// AddChatMembers adds multiple users to a group/channel.
func (b *Bot) AddChatMembers(ctx context.Context, chatID string, userIDs []string) error {
	_, err := b.makeRequest(ctx, "addChatMembers", map[string]any{
		"chat_id":  chatID,
		"user_ids": userIDs,
	})
	return err
}

// BanChatMembers bans multiple members sequentially. Returns the first
// error encountered (if any). Successfully banned members before the
// error are not rolled back.
func (b *Bot) BanChatMembers(ctx context.Context, chatID string, userIDs []string) error {
	for _, uid := range userIDs {
		if err := b.BanChatMember(ctx, chatID, uid); err != nil {
			return fmt.Errorf("ban %s: %w", uid, err)
		}
	}
	return nil
}

// UnbanChatMembers unbans multiple members sequentially.
func (b *Bot) UnbanChatMembers(ctx context.Context, chatID string, userIDs []string) error {
	for _, uid := range userIDs {
		if err := b.UnbanChatMember(ctx, chatID, uid); err != nil {
			return fmt.Errorf("unban %s: %w", uid, err)
		}
	}
	return nil
}

// RestrictChatMember restricts a member in a chat until the given time
// (Unix seconds; 0 or negative means forever).
func (b *Bot) RestrictChatMember(ctx context.Context, chatID, userID string, until int64) error {
	_, err := b.makeRequest(ctx, "restrictChatMember", map[string]any{
		"chat_id":    chatID,
		"user_id":    userID,
		"until_date": until,
	})
	return err
}

// PromoteChatMember elevates a member to admin with the supplied rights.
// rights is a list of permission strings (e.g. "SendMessages", "PinMessages").
func (b *Bot) PromoteChatMember(ctx context.Context, chatID, userID string, rights []string) error {
	_, err := b.makeRequest(ctx, "promoteChatMember", map[string]any{
		"chat_id":  chatID,
		"user_id":  userID,
		"rights":   rights,
		"is_admin": true,
	})
	return err
}

// DemoteChatMember removes admin rights from a member.
func (b *Bot) DemoteChatMember(ctx context.Context, chatID, userID string) error {
	_, err := b.makeRequest(ctx, "promoteChatMember", map[string]any{
		"chat_id":  chatID,
		"user_id":  userID,
		"rights":   []string{},
		"is_admin": false,
	})
	return err
}

// ExportChatInviteLink generates (or re-exports) the primary invite link.
func (b *Bot) ExportChatInviteLink(ctx context.Context, chatID string) (string, error) {
	data, err := b.makeRequest(ctx, "exportChatInviteLink", map[string]any{"chat_id": chatID})
	if err != nil {
		return "", err
	}
	link, _ := data["invite_link"].(string)
	return link, nil
}

// RevokeChatInviteLink revokes an existing invite link.
func (b *Bot) RevokeChatInviteLink(ctx context.Context, chatID, link string) error {
	_, err := b.makeRequest(ctx, "revokeChatInviteLink", map[string]any{
		"chat_id": chatID,
		"link":    link,
	})
	return err
}

// CreateGroup creates a new group with the given title and initial members.
func (b *Bot) CreateGroup(ctx context.Context, title string, memberUserIDs []string) (map[string]any, error) {
	return b.makeRequest(ctx, "createGroup", map[string]any{
		"title":           title,
		"member_user_ids": memberUserIDs,
	})
}

// CreateChannel creates a new channel.
func (b *Bot) CreateChannel(ctx context.Context, title, description string) (map[string]any, error) {
	payload := map[string]any{"title": title}
	if description != "" {
		payload["description"] = description
	}
	return b.makeRequest(ctx, "createChannel", payload)
}

// LeaveChat leaves a group or channel.
func (b *Bot) LeaveChat(ctx context.Context, chatID string) error {
	_, err := b.makeRequest(ctx, "leaveChat", map[string]any{"chat_id": chatID})
	return err
}

// RemoveChatKeypad removes the reply keypad for a chat (alias for
// EditChatKeypad with type=Remove).
func (b *Bot) RemoveChatKeypad(ctx context.Context, chatID string) error {
	return b.EditChatKeypad(ctx, chatID, models.ChatKeypadRemove, nil)
}

// RegisterAllEndpoints registers all 5 webhook endpoint types in one
// call.
func (b *Bot) RegisterAllEndpoints(ctx context.Context, url string) error {
	endpoints := []string{
		"ReceiveUpdate",
		"ReceiveInlineMessage",
		"ReceiveQuery",
		"GetSelectionItem",
		"SearchSelectionItems",
	}
	for _, ep := range endpoints {
		if _, err := b.UpdateBotEndpoints(ctx, url, ep); err != nil {
			return fmt.Errorf("register %s: %w", ep, err)
		}
	}
	return nil
}

// ===== Bot-side convenience helpers ==================
//
// Membership helpers (IsMember, IsAdmin, ...). Each one issues a
// getChatMember call and inspects the returned role and permissions.

// MemberIsAdmin returns true if the user is an admin or the chat owner.
func (b *Bot) MemberIsAdmin(ctx context.Context, chatID, userID string) (bool, error) {
	m, err := b.GetChatMember(ctx, chatID, userID)
	if err != nil {
		return false, err
	}
	role, _ := m["type"].(string)
	return role == "Admin" || role == "Creator" || role == "Owner", nil
}

// MemberIsOwner returns true if the user is the chat owner/creator.
func (b *Bot) MemberIsOwner(ctx context.Context, chatID, userID string) (bool, error) {
	m, err := b.GetChatMember(ctx, chatID, userID)
	if err != nil {
		return false, err
	}
	role, _ := m["type"].(string)
	return role == "Creator" || role == "Owner", nil
}

// MemberIsJoined returns true if the user is currently a member of the
// chat (i.e. getChatMember returns a non-empty result).
func (b *Bot) MemberIsJoined(ctx context.Context, chatID, userID string) (bool, error) {
	m, err := b.GetChatMember(ctx, chatID, userID)
	if err != nil {
		return false, err
	}
	_, ok := m["type"]
	return ok, nil
}

// CheckJoin verifies that a user has joined the given channel/group. It
// fetches the member list and scans for userID. This is a best-effort
// check — very large chats may not return all members in one page.
func (b *Bot) CheckJoin(ctx context.Context, chatID, userID string) (bool, error) {
	startID := ""
	for {
		page, err := b.GetChatMembers(ctx, chatID, startID)
		if err != nil {
			return false, err
		}
		members, _ := page["members"].([]any)
		for _, m := range members {
			if mm, ok := m.(map[string]any); ok {
				if uid, _ := mm["user_id"].(string); uid == userID {
					return true, nil
				}
			}
		}
		next, _ := page["next_start_id"].(string)
		if next == "" || next == startID {
			break
		}
		startID = next
	}
	return false, nil
}

// GetAvatarMe downloads the bot's own avatar. Returns the saved path or
// raw bytes depending on saveAs/asBytes (same semantics as DownloadFile).
func (b *Bot) GetAvatarMe(ctx context.Context, saveAs string, asBytes bool, chunkSize int, progress func(downloaded, total int64)) (any, error) {
	me, err := b.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	if me.Avatar == nil || me.Avatar.FileID == "" {
		return nil, &APIError{Status: "NotFound", DevMessage: "bot has no avatar"}
	}
	return b.DownloadFile(ctx, me.Avatar.FileID, saveAs, asBytes, chunkSize, progress)
}

// ===== Auto-delete / scheduled send helpers ============================
//
// Each helper spawns a goroutine that sleeps and then performs the
// action. The returned CancelFunc aborts it.

// CancelFunc aborts a pending background action.
type CancelFunc func()

// AutoDeleteMessage schedules a message for deletion after the given
// delay. The returned CancelFunc aborts the pending delete if called
// before the delay elapses.
func (b *Bot) AutoDeleteMessage(ctx context.Context, chatID, messageID string, delay time.Duration) CancelFunc {
	cctx, cancel := context.WithCancel(ctx)
	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-cctx.Done():
			return
		case <-timer.C:
			_ = b.DeleteMessage(cctx, chatID, messageID)
		}
	}()
	return CancelFunc(func() { cancel() })
}

// DeleteAfter is an alias for AutoDeleteMessage with a clearer name.
func (b *Bot) DeleteAfter(ctx context.Context, chatID, messageID string, delay time.Duration) CancelFunc {
	return b.AutoDeleteMessage(ctx, chatID, messageID, delay)
}

// EditAfter schedules a message edit after the given delay.
func (b *Bot) EditAfter(ctx context.Context, chatID, messageID, newText string, delay time.Duration) CancelFunc {
	cctx, cancel := context.WithCancel(ctx)
	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-cctx.Done():
			return
		case <-timer.C:
			_ = b.EditMessageText(cctx, chatID, messageID, newText)
		}
	}()
	return CancelFunc(func() { cancel() })
}

// DisableInlineKeyboard removes the inline keypad of a message after an
// optional delay (0 = immediately).
func (b *Bot) DisableInlineKeyboard(ctx context.Context, chatID, messageID string, delay time.Duration) CancelFunc {
	do := func() {
		_ = b.EditMessageKeypad(ctx, chatID, messageID, &models.Keypad{})
	}
	if delay <= 0 {
		do()
		return func() {}
	}
	cctx, cancel := context.WithCancel(ctx)
	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-cctx.Done():
			return
		case <-timer.C:
			do()
		}
	}()
	return CancelFunc(func() { cancel() })
}

// SendScheduledMessage sleeps for delay then sends the message. Returns
// the eventual *MessageID via the callback (non-blocking).
func (b *Bot) SendScheduledMessage(ctx context.Context, chatID, text string, delay time.Duration, opts ...SendOptions) CancelFunc {
	cctx, cancel := context.WithCancel(ctx)
	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-cctx.Done():
			return
		case <-timer.C:
			_, _ = b.SendMessage(cctx, chatID, text, opts...)
		}
	}()
	return CancelFunc(func() { cancel() })
}

// ===== Bulk / broadcast helpers =======================================
//
// Bulk sending and editing use a bounded goroutine pool so that
// broadcasting to thousands of chats does not exhaust memory.

// BulkResult is the outcome of a single bulk operation target.
type BulkResult struct {
	ChatID  string
	Err     error
	Success bool
}

// SendBulkMessage sends the same text to multiple chats concurrently.
// concurrency limits the number of in-flight sends. delayBetween is the
// pause between launches (0 = no pause). Returns per-chat results.
func (b *Bot) SendBulkMessage(ctx context.Context, chatIDs []string, text string, concurrency int, delayBetween time.Duration, opts ...SendOptions) []BulkResult {
	if concurrency <= 0 {
		concurrency = 5
	}
	results := make([]BulkResult, len(chatIDs))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for i, cid := range chatIDs {
		select {
		case <-ctx.Done():
			results[i] = BulkResult{ChatID: cid, Err: ctx.Err()}
			continue
		case sem <- struct{}{}:
		}
		wg.Add(1)
		go func(idx int, chatID string) {
			defer wg.Done()
			defer func() { <-sem }()
			_, err := b.SendMessage(ctx, chatID, text, opts...)
			results[idx] = BulkResult{ChatID: chatID, Err: err, Success: err == nil}
		}(i, cid)
		if delayBetween > 0 {
			select {
			case <-ctx.Done():
			case <-time.After(delayBetween):
			}
		}
	}
	wg.Wait()
	return results
}

// DeleteBulkMessage deletes multiple messages from the same chat
// concurrently. messageIDs must all belong to chatID.
func (b *Bot) DeleteBulkMessage(ctx context.Context, chatID string, messageIDs []string, concurrency int) []BulkResult {
	if concurrency <= 0 {
		concurrency = 5
	}
	results := make([]BulkResult, len(messageIDs))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for i, mid := range messageIDs {
		select {
		case <-ctx.Done():
			results[i] = BulkResult{ChatID: mid, Err: ctx.Err()}
			continue
		case sem <- struct{}{}:
		}
		wg.Add(1)
		go func(idx int, msgID string) {
			defer wg.Done()
			defer func() { <-sem }()
			err := b.DeleteMessage(ctx, chatID, msgID)
			results[idx] = BulkResult{ChatID: msgID, Err: err, Success: err == nil}
		}(i, mid)
	}
	wg.Wait()
	return results
}

// EditBulkMessage edits multiple messages in the same chat. messages maps
// messageID → new text.
func (b *Bot) EditBulkMessage(ctx context.Context, chatID string, messages map[string]string, concurrency int) []BulkResult {
	if concurrency <= 0 {
		concurrency = 5
	}
	results := make([]BulkResult, 0, len(messages))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for mid, txt := range messages {
		select {
		case <-ctx.Done():
			results = append(results, BulkResult{ChatID: mid, Err: ctx.Err()})
			continue
		case sem <- struct{}{}:
		}
		wg.Add(1)
		go func(messageID, text string) {
			defer wg.Done()
			defer func() { <-sem }()
			err := b.EditMessageText(ctx, chatID, messageID, text)
			mu.Lock()
			results = append(results, BulkResult{ChatID: messageID, Err: err, Success: err == nil})
			mu.Unlock()
		}(mid, txt)
	}
	wg.Wait()
	return results
}

// SendButtonJoin sends a message with a "Join Channel" inline button.
func (b *Bot) SendButtonJoin(ctx context.Context, chatID, title, username, text string, replyTo string) (*models.MessageID, error) {
	kp := &models.Keypad{
		Rows: []models.KeypadRow{{
			Buttons: []models.Button{widgetsJoinChannel("join_"+username, title, username)},
		}},
	}
	o := SendOptions{InlineKeypad: kp, ReplyToMessageID: replyTo}
	return b.SendMessage(ctx, chatID, text, o)
}

// SendButtonURL sends a message with a URL inline button.
func (b *Bot) SendButtonURL(ctx context.Context, chatID, title, url, text string, replyTo string) (*models.MessageID, error) {
	kp := &models.Keypad{
		Rows: []models.KeypadRow{{
			Buttons: []models.Button{widgetsLink("url_"+title, title, url)},
		}},
	}
	o := SendOptions{InlineKeypad: kp, ReplyToMessageID: replyTo}
	return b.SendMessage(ctx, chatID, text, o)
}
