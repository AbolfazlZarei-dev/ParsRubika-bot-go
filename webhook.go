package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// WebhookOptions تنظیمات وب‌هوک
type WebhookOptions struct {
	Port    int
	Path    string
	Handler HandlerFunc
	Secret  string // برای احراز هویت (امنیت)
}

// StartWebhook راه‌اندازی سرور وب‌هوک (بهینه شده)
func (c *BotClient) StartWebhook(ctx context.Context, opts WebhookOptions) error {
	if opts.Handler == nil {
		opts.Handler = c.ProcessUpdate
	}

	mux := http.NewServeMux()
	mux.HandleFunc(opts.Path, c.webhookHandler(opts))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", opts.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		<-ctx.Done()
		log.Println("در حال بستن سرور وب‌هوک...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("خطا در بستن سرور وب‌هوک: %v", err)
		}
		log.Println("سرور وب‌هوک با موفقیت بسته شد.")
	}()

	log.Printf("سرور وب‌هوک روی پورت %d و مسیر %s راه‌اندازی شد", opts.Port, opts.Path)

	if err := c.Start(ctx); err != nil {
		return err
	}

	return server.ListenAndServe()
}

// webhookHandler هندلر اصلی وب‌هوک (بهینه شده)
func (c *BotClient) webhookHandler(opts WebhookOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if opts.Secret != "" {
			secret := r.Header.Get("X-Rubika-Secret")
			if secret != opts.Secret {
				log.Printf("تلاش برای دسترسی غیرمجاز به وب‌هوک از آیپی: %s", r.RemoteAddr)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("خطا در خواندن بدنه درخواست: %v", err)
			http.Error(w, "خطا در خواندن بدنه", http.StatusBadRequest)
			return
		}

		update, err := c.parseWebhookBody(body)
		if err != nil {
			log.Printf("خطا در پارس کردن بدنه وب‌هوک: %v", err)
			http.Error(w, "فرمت بدنه نامعتبر", http.StatusBadRequest)
			return
		}

		if err := opts.Handler(r.Context(), update); err != nil {
			log.Printf("خطا در پردازش آپدیت وب‌هوک: %v", err)
			http.Error(w, "خطا در پردازش", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}

// parseWebhookBody پارس کردن بدنه وب‌هوک
func (c *BotClient) parseWebhookBody(body []byte) (*Update, error) {
	var updateWrapper struct {
		Update *Update `json:"update"`
	}
	if err := json.Unmarshal(body, &updateWrapper); err == nil && updateWrapper.Update != nil {
		return updateWrapper.Update, nil
	}

	var inlineWrapper struct {
		InlineMessage *InlineMessage `json:"inline_message"`
	}
	if err := json.Unmarshal(body, &inlineWrapper); err == nil && inlineWrapper.InlineMessage != nil {
		return c.convertInlineToUpdate(inlineWrapper.InlineMessage), nil
	}

	var queryWrapper struct {
		QueryID  string `json:"query_id"`
		ChatID   string `json:"chat_id"`
		SenderID string `json:"sender_id"`
		Data     string `json:"data"`
	}
	if err := json.Unmarshal(body, &queryWrapper); err == nil && queryWrapper.QueryID != "" {
		return c.convertQueryToUpdate(&queryWrapper), nil
	}

	return nil, fmt.Errorf("نوع وب‌هوک شناخته نشد")
}

// convertInlineToUpdate تبدیل InlineMessage به Update
func (c *BotClient) convertInlineToUpdate(inlineMsg *InlineMessage) *Update {
	return &Update{
		Type:   NewMessage,
		ChatID: inlineMsg.ChatID,
		NewMessage: &Message{
			MessageID: inlineMsg.MessageID,
			Text:      inlineMsg.Text,
			Time:      strconv.FormatInt(time.Now().Unix(), 10),
			SenderID:  inlineMsg.SenderID,
			AuxData:   inlineMsg.AuxData,
			File:      inlineMsg.File,
			Location:  inlineMsg.Location,
		},
	}
}

// convertQueryToUpdate تبدیل Query به Update
func (c *BotClient) convertQueryToUpdate(query *struct {
	QueryID  string `json:"query_id"`
	ChatID   string `json:"chat_id"`
	SenderID string `json:"sender_id"`
	Data     string `json:"data"`
}) *Update {
	return &Update{
		Type:   CallbackQuery,
		ChatID: query.ChatID,
		NewMessage: &Message{
			MessageID: query.QueryID,
			Text:      query.Data,
			Time:      strconv.FormatInt(time.Now().Unix(), 10),
			SenderID:  query.SenderID,
			AuxData:   &AuxData{ButtonID: &query.QueryID},
		},
	}
}
