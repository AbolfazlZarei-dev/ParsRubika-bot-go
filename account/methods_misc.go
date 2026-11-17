package account

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/netutil"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Reports, stories, Rubino bridges, link scraping and other API methods.

// GetAvailableReactions returns the reactions available to the user.
func (a *Account) GetAvailableReactions(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getAvailableReactions", Input: map[string]any{}, Encrypt: true})
}

// GetWalletTransactionMessage returns info about a wallet transaction.
func (a *Account) GetWalletTransactionMessage(ctx context.Context, transferID, accessTransfer string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getWalletTransactionMessage", Input: map[string]any{
		"transfer_id":     transferID,
		"access_transfer": accessTransfer,
	}, Encrypt: true})
}

// SearchGlobalObjects performs a global object search.
func (a *Account) SearchGlobalObjects(ctx context.Context, query string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "searchGlobalObjects", Input: map[string]any{"search_text": query}, Encrypt: true})
}

// GetProfileLinkItems returns the items shown on a profile link. The
// API uses `object_guid`, not `profile_link`.
func (a *Account) GetProfileLinkItems(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:     "getProfileLinkItems",
		Input:      map[string]any{"object_guid": objectGUID},
		TmpSession: true,
		Encrypt:    true,
	})
}

// GetTime returns the server time.
func (a *Account) GetTime(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getTime", Input: map[string]any{}, TmpSession: true, Encrypt: true})
}

// GetRelatedObjects returns the related objects of a GUID.
func (a *Account) GetRelatedObjects(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:     "getRelatedObjects",
		Input:      map[string]any{"object_guid": objectGUID},
		TmpSession: true,
		Encrypt:    true,
	})
}

// ReportObject reports an object for abuse. The API uses
// `{object_guid, report_type, report_description, report_type_object,
// message_id?}`.
func (a *Account) ReportObject(ctx context.Context, objectGUID, reportType, reportDescription, reportTypeObject, messageID string) (*events.Event, error) {
	in := map[string]any{
		"object_guid":        objectGUID,
		"report_type":        reportType,
		"report_description": reportDescription,
		"report_type_object": reportTypeObject,
	}
	if messageID != "" {
		in["message_id"] = messageID
	}
	return a.Request(ctx, RequestOptions{Method: "reportObject", Input: in, Encrypt: true})
}

// ReportMessage is a convenience wrapper for reporting a single
// message.
func (a *Account) ReportMessage(ctx context.Context, objectGUID, messageID, reportType, reportDescription string) (*events.Event, error) {
	return a.ReportObject(ctx, objectGUID, reportType, reportDescription, "Message", messageID)
}

// ReportChat is a convenience wrapper for reporting a chat.
func (a *Account) ReportChat(ctx context.Context, objectGUID, reportType, reportDescription string) (*events.Event, error) {
	return a.ReportObject(ctx, objectGUID, reportType, reportDescription, "Object", "")
}

// stripAt removes a leading `@` from a username.
func stripAt(s string) string {
	return strings.TrimPrefix(s, "@")
}

// channelPostLinkRegex matches rubika.ir/post/<hash> URLs.
var channelPostLinkRegex = regexp.MustCompile(`https?://rubika\.ir/post/([A-Za-z0-9]+)`)

// GetChannelPostByLink resolves a channel post from a web URL. It calls
// getLinkFromAppUrl to convert the web link into a deep link, then parses
// the result to extract channel GUID and message ID.
func (a *Account) GetChannelPostByLink(ctx context.Context, postURL string) (*events.Event, error) {
	resp, err := a.GetLinkFromAppURL(ctx, postURL)
	if err != nil {
		return nil, err
	}
	// The response contains a deep link like:
	// rubika://c.rubika.ir/<channel_guid>/<message_id>
	link := resp.AnyStr("link")
	if link == "" {
		// Try alternate field names
		link = resp.AnyStr("app_url", "url", "deep_link")
	}
	if link == "" {
		return resp, nil
	}
	// Parse the deep link
	parts := strings.Split(strings.TrimPrefix(link, "rubika://c.rubika.ir/"), "/")
	if len(parts) >= 2 {
		channelGUID := parts[0]
		messageID := parts[1]
		// Fetch the actual message
		return a.GetMessagesByID(ctx, channelGUID, []string{messageID})
	}
	return resp, nil
}

// GetPostByAskLink resolves a post from an ask link. It calls
// getLinkFromAppUrl and then fetches the message.
func (a *Account) GetPostByAskLink(ctx context.Context, askLink string) (*events.Event, error) {
	return a.GetChannelPostByLink(ctx, askLink)
}

// ScrapeChannelPostLink fetches the rubika.ir web page for a post and
// extracts the deep link. This is a fallback when getLinkFromAppUrl
// does not resolve.
//
// Security: only http(s) URLs that match the rubika.ir post-link pattern
// are accepted; the dial is routed through netutil's SSRF-safe transport
// so private/loopback/link-local IPs cannot be reached even via DNS
// rebinding, and the response is capped at netutil.MaxFetchBytes.
func (a *Account) ScrapeChannelPostLink(ctx context.Context, postURL string) (string, error) {
	matches := channelPostLinkRegex.FindStringSubmatch(postURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid channel post URL: %s", postURL)
	}
	if _, err := netutil.AssertURL(ctx, postURL); err != nil {
		return "", err
	}
	client := netutil.NewSafeClient(a.timeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, postURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, netutil.MaxFetchBytes+1))
	if err != nil {
		return "", err
	}
	if len(body) > netutil.MaxFetchBytes {
		return "", fmt.Errorf("scrape %s: response exceeds %d bytes", postURL, netutil.MaxFetchBytes)
	}
	// Look for a deep link or meta refresh in the HTML.
	// The page typically contains: <meta property="og:url" content="rubika://...">
	html := string(body)
	idx := strings.Index(html, "rubika://")
	if idx >= 0 {
		end := strings.IndexAny(html[idx:], "\"'<>\n ")
		if end > 0 {
			return html[idx : idx+end], nil
		}
	}
	return "", fmt.Errorf("could not extract deep link from %s", postURL)
}

// FetchJSON is a low-level helper that performs an HTTP GET and decodes
// the response as JSON. It is useful for calling third-party APIs
// (e.g. translation services, Rubino download helpers).
//
// Security: the URL is validated through the netutil SSRF policy and the
// dial goes through a SSRF-safe transport; the response is capped at
// netutil.MaxFetchBytes.
func (a *Account) FetchJSON(ctx context.Context, url string, out any) error {
	if _, err := netutil.AssertURL(ctx, url); err != nil {
		return err
	}
	client := netutil.NewSafeClient(a.timeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, netutil.MaxFetchBytes+1))
	if err != nil {
		return err
	}
	if len(body) > netutil.MaxFetchBytes {
		return fmt.Errorf("fetch %s: response exceeds %d bytes", url, netutil.MaxFetchBytes)
	}
	return json.Unmarshal(body, out)
}

// GetContactsStories returns the stories of the user's contacts.
//
// (GetContactsLastOnline is already defined in methods_users.go.)
func (a *Account) GetContactsStories(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getStories", Input: map[string]any{}, Encrypt: true})
}

// SendRubinoPost shares a Rubino post to a chat.
func (a *Account) SendRubinoPost(ctx context.Context, objectGUID, postID, postProfileID string, isMute bool) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "sendRubinoPost", Input: map[string]any{
		"object_guid":     objectGUID,
		"post_id":         postID,
		"post_profile_id": postProfileID,
		"rnd":             randInt63(),
		"is_mute":         isMute,
	}, Encrypt: true})
}

// SendRubinoProduct shares a Rubino product to a chat.
func (a *Account) SendRubinoProduct(ctx context.Context, objectGUID, storeID, productID string, productVarieties []map[string]any) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "sendRubinoProduct", Input: map[string]any{
		"object_guid":       objectGUID,
		"store_id":          storeID,
		"product_id":        productID,
		"product_varieties": productVarieties,
		"rnd":               randInt63(),
	}, Encrypt: true})
}

// SendRubinoStory replies to or shares a Rubino story.
func (a *Account) SendRubinoStory(ctx context.Context, objectGUID, storyID, storyProfileID, replyText, storyType string, isMute bool) (*events.Event, error) {
	if storyType == "" {
		storyType = "Reply"
	}
	return a.Request(ctx, RequestOptions{Method: "sendRubinoStory", Input: map[string]any{
		"object_guid":      objectGUID,
		"story_id":         storyID,
		"story_profile_id": storyProfileID,
		"reply_text":       replyText,
		"type":             storyType, // "Reply" or "Direct"
		"rnd":              randInt63(),
		"is_mute":          isMute,
	}, Encrypt: true})
}

// rubikaLinkRegex matches any rubika.ir URL.
var rubikaLinkRegex = regexp.MustCompile(`https?://rubika\.ir/\S+`)

// groupLinkRegex matches rubika.ir/joing/<hash> URLs.
var groupLinkRegex = regexp.MustCompile(`https?://rubika\.ir/joing/[A-Za-z0-9]+`)

// channelLinkRegex matches rubika.ir/joinc/<hash> URLs.
var channelLinkRegex = regexp.MustCompile(`https?://rubika\.ir/joinc/[A-Za-z0-9]+`)

// postLinkRegex matches rubika.ir/post/<hash> URLs.
var postLinkRegex = regexp.MustCompile(`https?://rubika\.ir/post/[A-Za-z0-9]+`)

// usernameRegex matches @username patterns (4-32 chars).
var usernameRegex = regexp.MustCompile(`@([A-Za-z0-9_]{4,32})`)

// IsRubikaLink reports whether s contains a rubika.ir URL.
func IsRubikaLink(s string) bool { return rubikaLinkRegex.MatchString(s) }

// IsGroupLink reports whether s is a rubika group join link.
func IsGroupLink(s string) bool { return groupLinkRegex.MatchString(s) }

// IsChannelLink reports whether s is a rubika channel join link.
func IsChannelLink(s string) bool { return channelLinkRegex.MatchString(s) }

// IsPostLink reports whether s is a rubika post link.
func IsPostLink(s string) bool { return postLinkRegex.MatchString(s) }

// IsUsername reports whether s is a @username (4-32 alphanumeric/underscore).
func IsUsername(s string) bool { return usernameRegex.MatchString(s) }

// GetRubikaLinks extracts all rubika.ir URLs from s.
func GetRubikaLinks(s string) []string {
	return rubikaLinkRegex.FindAllString(s, -1)
}

// GetGroupLinks extracts all group join links from s.
func GetGroupLinks(s string) []string {
	return groupLinkRegex.FindAllString(s, -1)
}

// GetChannelLinks extracts all channel join links from s.
func GetChannelLinks(s string) []string {
	return channelLinkRegex.FindAllString(s, -1)
}

// GetPostLinks extracts all post links from s.
func GetPostLinks(s string) []string {
	return postLinkRegex.FindAllString(s, -1)
}

// GetUsernames extracts all @usernames from s.
func GetUsernames(s string) []string {
	return usernameRegex.FindAllString(s, -1)
}

// GetServerTime fetches the current server time. The server answers with a
// timestamp (as a string or number); it is parsed robustly.
func (a *Account) GetServerTime(ctx context.Context) (time.Time, error) {
	resp, err := a.GetTime(ctx)
	if err != nil {
		return time.Time{}, err
	}
	ts := resp.AnyStr("time")
	if ts == "" {
		return time.Time{}, fmt.Errorf("no time field in response")
	}
	t, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time %q: %w", ts, err)
	}
	return time.Unix(t, 0), nil
}
