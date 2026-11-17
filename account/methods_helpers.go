package account

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

// High-level convenience helpers built on the raw API methods.

// String returns a human-readable description of the account.
func (a *Account) String() string {
	return fmt.Sprintf("ParsRubika.Account(guid=%s, auth=%s)", a.guid, maskAuth(a.auth))
}

// maskAuth returns a masked version of the auth string for safe logging.
func maskAuth(auth string) string {
	if len(auth) <= 8 {
		return strings.Repeat("*", len(auth))
	}
	return auth[:4] + strings.Repeat("*", len(auth)-8) + auth[len(auth)-4:]
}

// BanMembers bans multiple members from a group or channel. The GUID
// prefix (g0/c0) selects the correct underlying API.
func (a *Account) BanMembers(ctx context.Context, objectGUID string, memberGUIDs []string) error {
	for _, mg := range memberGUIDs {
		if _, err := a.BanMember(ctx, objectGUID, mg); err != nil {
			return fmt.Errorf("ban %s: %w", mg, err)
		}
	}
	return nil
}

// UnbanMembers unbans multiple members from a group or channel.
func (a *Account) UnbanMembers(ctx context.Context, objectGUID string, memberGUIDs []string) error {
	for _, mg := range memberGUIDs {
		if _, err := a.UnbanMember(ctx, objectGUID, mg); err != nil {
			return fmt.Errorf("unban %s: %w", mg, err)
		}
	}
	return nil
}

// UnbanAllMembers fetches the banned-members list and unbans everyone.
// excludeGUIDs (if non-nil) lists members that should remain banned.
// The banned-members API returns members under "in_chat_members" with
// each member's GUID in the "user_guid" field (not "object_guid").
func (a *Account) UnbanAllMembers(ctx context.Context, objectGUID string, excludeGUIDs []string) error {
	exclude := make(map[string]struct{}, len(excludeGUIDs))
	for _, g := range excludeGUIDs {
		exclude[g] = struct{}{}
	}
	startID := ""
	for {
		page, err := a.GetBannedGroupMembers(ctx, objectGUID, startID)
		if err != nil {
			return err
		}
		members := page.Slice("in_chat_members")
		for _, m := range members {
			if mm, ok := m.(map[string]any); ok {
				// Banned-member objects use "user_guid" for the member's GUID.
				guid, _ := mm["user_guid"].(string)
				if guid == "" {
					continue
				}
				if _, skip := exclude[guid]; skip {
					continue
				}
				if _, err := a.UnbanMember(ctx, objectGUID, guid); err != nil {
					return err
				}
			}
		}
		ns := page.AnyStr("next_start_id")
		if ns == "" || ns == startID {
			break
		}
		startID = ns
	}
	return nil
}

// ClearBlackList is an alias for UnbanAllMembers.
func (a *Account) ClearBlackList(ctx context.Context, objectGUID string, excludeGUIDs []string) error {
	return a.UnbanAllMembers(ctx, objectGUID, excludeGUIDs)
}

// AcceptAllJoinRequests accepts every pending join request for a group
// or channel. excludeGUIDs (if non-nil) lists users to skip.
func (a *Account) AcceptAllJoinRequests(ctx context.Context, objectGUID string, excludeGUIDs []string) error {
	exclude := make(map[string]struct{}, len(excludeGUIDs))
	for _, g := range excludeGUIDs {
		exclude[g] = struct{}{}
	}
	cursor := ""
	for {
		page, err := a.GetJoinRequests(ctx, objectGUID, cursor)
		if err != nil {
			return err
		}
		reqs := page.Slice("join_requests")
		for _, r := range reqs {
			if mm, ok := r.(map[string]any); ok {
				guid, _ := mm["user_guid"].(string)
				if guid == "" {
					continue
				}
				if _, skip := exclude[guid]; skip {
					continue
				}
				if _, err := a.AcceptJoinRequest(ctx, objectGUID, guid); err != nil {
					return err
				}
			}
		}
		ns := page.AnyStr("next_start_id")
		if ns == "" || ns == cursor {
			break
		}
		cursor = ns
	}
	return nil
}

// GetMessageInfo fetches a single message by ID (delegates to
// GetMessagesByID with a single-element slice).
func (a *Account) GetMessageInfo(ctx context.Context, objectGUID, messageID string) (*events.Event, error) {
	return a.GetMessagesByID(ctx, objectGUID, []string{messageID})
}

// GetLastMessage returns the last message in a chat. Rubika chat-info
// responses nest the last message under the "last_message" key (an
// object containing "message_id"), not as a top-level "last_message_id"
// string. We extract the message_id from that nested object.
func (a *Account) GetLastMessage(ctx context.Context, objectGUID string) (*events.Event, error) {
	info, err := a.GetInfo(ctx, objectGUID)
	if err != nil {
		return nil, err
	}
	// Rubika nests the last message as an object under "last_message".
	lastMsg := info.EventOf("last_message")
	if lastMsg != nil {
		lastID := lastMsg.AnyStr("message_id")
		if lastID != "" {
			return a.GetMessageInfo(ctx, objectGUID, lastID)
		}
	}
	// Fallback: some chat types return last_message_id as a top-level field.
	lastID := info.AnyStr("last_message_id")
	if lastID == "" {
		return nil, fmt.Errorf("no last_message in chat info")
	}
	return a.GetMessageInfo(ctx, objectGUID, lastID)
}

// GetFirstMessage returns the first message in a chat by fetching the
// messages interval and picking the lowest message_id. Message IDs are
// compared numerically (not lexicographically) to handle IDs with
// different digit counts correctly.
func (a *Account) GetFirstMessage(ctx context.Context, objectGUID string) (*events.Event, error) {
	interval, err := a.GetMessagesInterval(ctx, objectGUID, "")
	if err != nil {
		return nil, err
	}
	msgs := interval.Slice("messages")
	if len(msgs) == 0 {
		return nil, fmt.Errorf("no messages in interval")
	}
	var firstID int64 = -1
	for _, m := range msgs {
		if mm, ok := m.(map[string]any); ok {
			midStr := fmt.Sprintf("%v", mm["message_id"])
			midVal, _ := strconv.ParseInt(midStr, 10, 64)
			if firstID < 0 || (midVal > 0 && midVal < firstID) {
				firstID = midVal
			}
		}
	}
	if firstID < 0 {
		return nil, fmt.Errorf("could not determine first message id")
	}
	return a.GetMessageInfo(ctx, objectGUID, strconv.FormatInt(firstID, 10))
}

// SearchMessages searches messages in a chat by text (client-side filter
// over getMessages). It fetches up to maxPages pages of 50 messages each and returns those
// whose text contains the query (case-insensitive). The cursor advances
// to the OLDEST message in each page (msgs[len-1]) so each page fetches
// a distinct window of messages.
func (a *Account) SearchMessages(ctx context.Context, objectGUID, query string, maxPages int) ([]map[string]any, error) {
	if maxPages <= 0 {
		maxPages = 5
	}
	query = strings.ToLower(query)
	var out []map[string]any
	maxID := ""
	for i := 0; i < maxPages; i++ {
		page, err := a.GetMessages(ctx, objectGUID, maxID, 50, "FromMax")
		if err != nil {
			return out, err
		}
		msgs := page.Slice("messages")
		if len(msgs) == 0 {
			break
		}
		for _, m := range msgs {
			mm, ok := m.(map[string]any)
			if !ok {
				continue
			}
			text, _ := mm["text"].(string)
			if strings.Contains(strings.ToLower(text), query) {
				out = append(out, mm)
			}
		}
		// Advance cursor to the OLDEST message in this page so the next
		// call fetches an earlier window. With sort=FromMax, msgs[0] is
		// the newest and msgs[len-1] is the oldest.
		if last, ok := msgs[len(msgs)-1].(map[string]any); ok {
			maxID = fmt.Sprintf("%v", last["message_id"])
		}
	}
	return out, nil
}

// DeleteMyMessages deletes all of the current user's messages in a chat
// by fetching and deleting them in batches. The author is identified by
// the "author_object_guid" field (not "author_guid"). The cursor advances
// to the OLDEST message in each page so the next call fetches an earlier
// window — using the newest would re-fetch nearly the same page and
// loop ~forever.
func (a *Account) DeleteMyMessages(ctx context.Context, objectGUID string) error {
	maxID := ""
	for {
		page, err := a.GetMessages(ctx, objectGUID, maxID, 50, "FromMax")
		if err != nil {
			return err
		}
		msgs := page.Slice("messages")
		if len(msgs) == 0 {
			break
		}
		var ids []string
		for _, m := range msgs {
			if mm, ok := m.(map[string]any); ok {
				// Rubika uses "author_object_guid" for the message author.
				if author, _ := mm["author_object_guid"].(string); author == a.guid {
					ids = append(ids, fmt.Sprintf("%v", mm["message_id"]))
				}
			}
		}
		if len(ids) > 0 {
			if _, err := a.DeleteMessages(ctx, objectGUID, ids); err != nil {
				return err
			}
		}
		// Advance cursor to the OLDEST message in this page.
		if last, ok := msgs[len(msgs)-1].(map[string]any); ok {
			maxID = fmt.Sprintf("%v", last["message_id"])
		}
	}
	return nil
}

// DeleteAllMessages deletes every message in a chat (regardless of
// author). Use with caution. The cursor advances to the OLDEST message
// in each page (see DeleteMyMessages for rationale).
func (a *Account) DeleteAllMessages(ctx context.Context, objectGUID string) error {
	maxID := ""
	for {
		page, err := a.GetMessages(ctx, objectGUID, maxID, 50, "FromMax")
		if err != nil {
			return err
		}
		msgs := page.Slice("messages")
		if len(msgs) == 0 {
			break
		}
		var ids []string
		for _, m := range msgs {
			if mm, ok := m.(map[string]any); ok {
				ids = append(ids, fmt.Sprintf("%v", mm["message_id"]))
			}
		}
		if len(ids) > 0 {
			if _, err := a.DeleteMessages(ctx, objectGUID, ids); err != nil {
				return err
			}
		}
		// Advance cursor to the OLDEST message in this page.
		if last, ok := msgs[len(msgs)-1].(map[string]any); ok {
			maxID = fmt.Sprintf("%v", last["message_id"])
		}
	}
	return nil
}

// GetGUID resolves a username, join-link, or GUID to a GUID. If the
// input is already a GUID (starts with u0/g0/c0/b0/s0) it is returned
// as-is. Join-links (rubika.ir/joing/... or /joinc/...) are resolved
// via groupPreviewByJoinLink / channelPreviewByJoinLink.
func (a *Account) GetGUID(ctx context.Context, input string) (string, error) {
	if isGUID(input) {
		return input, nil
	}
	if strings.HasPrefix(input, "@") {
		obj, err := a.GetObjectByUsername(ctx, strings.TrimPrefix(input, "@"))
		if err != nil {
			return "", err
		}
		for _, key := range []string{"user", "channel", "group", "bot"} {
			if objMap := obj.Map(key); objMap != nil {
				if guid, _ := objMap["object_guid"].(string); guid != "" {
					return guid, nil
				}
			}
		}
		return "", fmt.Errorf("could not resolve username %q", input)
	}
	// Join-link detection: /joing/ → group, /joinc/ → channel.
	if strings.Contains(input, "/joing/") {
		preview, err := a.GroupPreviewByJoinLink(ctx, input)
		if err != nil {
			return "", err
		}
		guid := preview.AnyStr("group_guid", "object_guid")
		if guid == "" {
			return "", fmt.Errorf("could not resolve group link %q", input)
		}
		return guid, nil
	}
	if strings.Contains(input, "/joinc/") {
		preview, err := a.ChannelPreviewByJoinLink(ctx, input)
		if err != nil {
			return "", err
		}
		guid := preview.AnyStr("channel_guid", "object_guid")
		if guid == "" {
			return "", fmt.Errorf("could not resolve channel link %q", input)
		}
		return guid, nil
	}
	// Fallback: try treating it as a GUID-less chat info.
	preview, err := a.GetInfo(ctx, input)
	if err != nil {
		return "", err
	}
	guid := preview.AnyStr("object_guid")
	if guid == "" {
		return "", fmt.Errorf("could not resolve link %q", input)
	}
	return guid, nil
}

// GetMyGUID returns the current user's GUID.
func (a *Account) GetMyGUID() string { return a.guid }

// CheckJoin verifies that a user has joined a group/channel by fetching
// the member list and scanning for the user's GUID.
func (a *Account) CheckJoin(ctx context.Context, objectGUID, userGUID string) (bool, error) {
	startID := ""
	for i := 0; i < 50; i++ {
		var page *events.Event
		var err error
		if strings.HasPrefix(objectGUID, "c0") {
			page, err = a.GetChannelAllMembers(ctx, objectGUID, startID, "", 200)
		} else {
			page, err = a.GetGroupAllMembers(ctx, objectGUID, startID, "")
		}
		if err != nil {
			return false, err
		}
		// Regular member-list APIs return results under "in_chat"
		// (not "in_chat_members", which is for banned-members).
		members := page.Slice("in_chat")
		for _, m := range members {
			if mm, ok := m.(map[string]any); ok {
				// Member objects use "user_guid" for the member's GUID.
				if guid, _ := mm["user_guid"].(string); guid == userGUID {
					return true, nil
				}
			}
		}
		ns := page.AnyStr("next_start_id")
		if ns == "" || ns == startID {
			break
		}
		startID = ns
	}
	return false, nil
}

// BulkSendResult is the outcome of a single bulk send target.
type BulkSendResult struct {
	ObjectGUID string
	Err        error
	Success    bool
}

// SendBulkMessage sends the same text to multiple chats concurrently.
// concurrency limits in-flight sends.
func (a *Account) SendBulkMessage(ctx context.Context, objectGUIDs []string, text string, concurrency int, delayBetween time.Duration, opts ...SendOptions) []BulkSendResult {
	if concurrency <= 0 {
		concurrency = 5
	}
	results := make([]BulkSendResult, len(objectGUIDs))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for i, oid := range objectGUIDs {
		select {
		case <-ctx.Done():
			results[i] = BulkSendResult{ObjectGUID: oid, Err: ctx.Err()}
			continue
		case sem <- struct{}{}:
		}
		wg.Add(1)
		go func(idx int, guid string) {
			defer wg.Done()
			defer func() { <-sem }()
			_, err := a.SendMessage(ctx, guid, text, opts...)
			results[idx] = BulkSendResult{ObjectGUID: guid, Err: err, Success: err == nil}
		}(i, oid)
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

// isGUID reports whether s looks like a Rubika GUID (starts with one of
// the known prefixes and is at least 34 chars long).
func isGUID(s string) bool {
	for _, p := range []string{"u0", "g0", "c0", "b0", "s0"} {
		if strings.HasPrefix(s, p) && len(s) >= 34 {
			return true
		}
	}
	return false
}

func randInt63() int64 { return rand.Int63n(1_000_000) + 1 }

// anyToStr coerces a scalar value (string, float64, int, json.Number)
// to its string form. Returns "" for unsupported types.
func anyToStr(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatInt(int64(x), 10)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case json.Number:
		return x.String()
	default:
		return ""
	}
}
