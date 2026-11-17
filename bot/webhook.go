package bot

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// parseUpdate converts a raw update payload into a typed *models.Update.
func parseUpdate(raw map[string]any, b *Bot) (*models.Update, error) {
	t, _ := raw["type"].(string)
	chatID, _ := raw["chat_id"].(string)
	upd := &models.Update{Type: models.UpdateType(t), ChatID: chatID}
	if v, ok := raw["update_time"]; ok {
		upd.UpdateTime = toString(v)
	}
	switch models.UpdateType(t) {
	case models.UpdateTypeRemovedMessage:
		upd.RemovedMessageID = toString(raw["removed_message_id"])
		return upd, nil
	case models.UpdateTypeStartedBot, models.UpdateTypeStoppedBot:
		return upd, nil
	case models.UpdateTypeUpdatedPayment:
		if pr, ok := raw["updated_payment"].(map[string]any); ok {
			jb, _ := json.Marshal(pr)
			var ps models.PaymentStatus
			if err := json.Unmarshal(jb, &ps); err == nil {
				upd.UpdatedPayment = &ps
			}
		}
		return upd, nil
	case models.UpdateTypeNewMessage, models.UpdateTypeUpdatedMessage:
		key := "new_message"
		if models.UpdateType(t) == models.UpdateTypeUpdatedMessage {
			key = "updated_message"
		}
		msgRaw, _ := raw[key].(map[string]any)
		if msgRaw == nil {
			return nil, &APIError{Status: "InvalidResponse", DevMessage: "missing message payload"}
		}
		rawJSON, err := json.Marshal(msgRaw)
		if err != nil {
			return nil, err
		}
		var msg models.Message
		if err := json.Unmarshal(rawJSON, &msg); err != nil {
			return nil, err
		}
		if models.UpdateType(t) == models.UpdateTypeNewMessage {
			upd.NewMessage = &msg
		} else {
			upd.UpdatedMessage = &msg
		}
		return upd, nil
	}
	return nil, &APIError{Status: "InvalidResponse", DevMessage: "unsupported update type: " + t}
}

// parseInlineMessage converts a raw inline-message payload into a
// typed *models.InlineMessage.
func parseInlineMessage(raw map[string]any, b *Bot) *models.InlineMessage {
	im := &models.InlineMessage{
		SenderID:  toString(raw["sender_id"]),
		Text:      parsedString(raw, "text"),
		MessageID: toString(raw["message_id"]),
		ChatID:    parsedString(raw, "chat_id"),
	}
	decode := func(key string, dst any) bool {
		sub, ok := raw[key].(map[string]any)
		if !ok {
			return false
		}
		jb, _ := json.Marshal(sub)
		return json.Unmarshal(jb, dst) == nil
	}
	var f models.File
	if decode("file", &f) {
		im.File = &f
	}
	var l models.Location
	if decode("location", &l) {
		im.Location = &l
	}
	var a models.AuxData
	if decode("aux_data", &a) {
		im.AuxData = &a
	}
	return im
}

// buildRawUpdateKey identifies a raw (not yet parsed) update.
func buildRawUpdateKey(item map[string]any) string {
	t, _ := item["type"].(string)
	if t == "" {
		return ""
	}
	switch models.UpdateType(t) {
	case models.UpdateTypeRemovedMessage:
		if id := toString(item["removed_message_id"]); id != "" {
			if chat := toString(item["chat_id"]); chat != "" {
				return t + ":" + id + ":" + chat
			}
			return t + ":" + id
		}
	case models.UpdateTypeNewMessage, models.UpdateTypeUpdatedMessage:
		key := "new_message"
		if models.UpdateType(t) == models.UpdateTypeUpdatedMessage {
			key = "updated_message"
		}
		msg, _ := item[key].(map[string]any)
		if msg == nil {
			return ""
		}
		id := toString(msg["message_id"])
		if id == "" {
			return ""
		}
		updateTime := ""
		if v, ok := item["update_time"]; ok {
			updateTime = toString(v)
		} else if v, ok := msg["time"]; ok {
			updateTime = toString(v)
		}
		return t + ":" + id + ":" + updateTime
	}
	return ""
}

// handleWebhook receives updates pushed by the platform.
func (b *Bot) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if sec := b.opts.WebhookSecret; sec != "" {
		got := r.URL.Query().Get("s")
		if subtle.ConstantTimeCompare([]byte(got), []byte(sec)) != 1 {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}
	if max := b.opts.MaxBodyBytes; max > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, max)
	}
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"OK"}`))

	ctx := b.lifecycleContext()
	if inlineRaw, ok := payload["inline_message"].(map[string]any); ok {
		inline := parseInlineMessage(inlineRaw, b)
		b.submit(ctx, wrapInlineAsUpdate(inline))
		return
	}
	if updRaw, ok := payload["update"].(map[string]any); ok {
		if upd, err := parseUpdate(updRaw, b); err == nil && upd != nil {
			b.submit(ctx, upd)
		}
	}
}

// wrapInlineAsUpdate wraps an InlineMessage into a *models.Update so
// handlers registered via On(...) can process it through the same
// signature.
func wrapInlineAsUpdate(i *models.InlineMessage) *models.Update {
	return &models.Update{
		Type:   models.UpdateTypeInlineMessage,
		ChatID: i.ChatID,
		NewMessage: &models.Message{
			MessageID: i.MessageID,
			Text:      i.Text,
			SenderID:  i.SenderID,
			AuxData:   i.AuxData,
			File:      i.File,
			Location:  i.Location,
		},
	}
}

// WebhookHandler returns the http.Handler that serves every webhook
// endpoint, so the bot can be mounted on an existing HTTP server.
func (b *Bot) WebhookHandler() http.Handler {
	mux := http.NewServeMux()
	for _, p := range []string{
		"/webhook", "/webhook/receiveUpdate", "/webhook/receiveInlineMessage",
		"/webhook/receiveQuery", "/webhook/getSelectionItem", "/webhook/searchSelectionItems",
	} {
		mux.HandleFunc(p, b.handleWebhook)
	}
	return mux
}

// runWebhook registers the endpoints and serves them until ctx ends.
func (b *Bot) runWebhook(ctx context.Context, webhookURL, host string, port int) error {
	if err := b.onStart(ctx); err != nil {
		return err
	}
	defer b.onStop(context.Background())
	defer b.drain()

	base := strings.TrimRight(webhookURL, "/") + "/webhook"
	for _, t := range []models.UpdateEndpointType{
		models.EndpointReceiveUpdate,
		models.EndpointReceiveInlineMessage,
		models.EndpointReceiveQuery,
		models.EndpointGetSelectionItem,
		models.EndpointSearchSelectionItems,
	} {
		target := base + "/" + string(t)
		if sec := b.opts.WebhookSecret; sec != "" {
			target += "?s=" + url.QueryEscape(sec)
		}
		if _, err := b.UpdateBotEndpoints(ctx, target, string(t)); err != nil {
			b.warnf("register webhook %s: %v", t, err)
		}
	}
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		Handler:           b.WebhookHandler(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
	}
	go func() {
		<-ctx.Done()
		sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(sctx)
	}()
	b.infof("webhook server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
