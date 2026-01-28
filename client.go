package ParsRubika

// نسخه: 2.5.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultBaseURL      = "https://botapi.rubika.ir/v3"
	defaultMessengerURL = "https://messengerg2b1.iranlms.ir/v3"
	userAgent           = "ParsRubika/2.5.0"
	version             = "2.5.0"
)

// HandlerFunc برای پردازش آپدیت‌ها
type HandlerFunc func(ctx context.Context, update *Update) error

// MiddlewareFunc برای middleware
type MiddlewareFunc func(ctx context.Context, update *Update, next HandlerFunc) error

// FilterFunc برای فیلتر کردن آپدیت‌ها
type FilterFunc func(ctx context.Context, update *Update) bool

// BotClient کلاینت اصلی بات
type BotClient struct {
	token        string
	baseURL      string
	messengerURL string
	httpClient   *http.Client
	botID        string
	mu           sync.RWMutex
	lastSentTime time.Time

	// مدیریت API
	apiManager     *APIManager
	connectionMode ConnectionMode

	// مدیریت وضعیت
	isRunning bool
	stopChan  chan struct{}

	// هندلرها و میدلورها
	handlers      []Handler
	middlewares   []MiddlewareFunc
	updateFilters []FilterFunc

	// نرخ محدودیت
	rateLimitDuration time.Duration
	maxRetries        int

	// قابلیت‌های جدید
	antiSpam         *AntiSpam
	hotReloadEnabled bool
	stateManager     *StateManager
	ignoreTimeout    bool
	metadata         map[string]interface{}
	reloadManager    *ReloadManager

	// اطلاع‌رسانی
	notificationOpts *NotificationOptions

	// هاست ری‌لود واقعی
	hostReloadWatcher *HostReloadWatcher

	// مدیریت پایداری شبکه
	networkStabilityManager *NetworkStabilityManager

	// برای ذخیره اطلاعات کاربران
	registeredUsers []*User
	fileMu          sync.Mutex // برای دسترسی امن به فایل
	usersFilePath   string     // مسیر فایل ذخیره کاربران

	// برای ذخیره آفست آیدی
	offsetsDir     string // پوشه اصلی برای ذخیره فایل‌های آفست
	offsetFilePath string // مسیر کامل فایل آفست که در init() تنظیم می‌شود
}

// Handler اطلاعات هندلر
type Handler struct {
	handler HandlerFunc
	filter  FilterFunc
	order   int
}

// ClientOption برای تنظیمات اولیه کلاینت
type ClientOption func(*BotClient)

// WithBaseURL برای تغییر آدرس پایه API
func WithBaseURL(url string) ClientOption {
	return func(c *BotClient) {
		c.baseURL = url
	}
}

// WithMessengerURL برای تغییر آدرس پایه API Messenger
func WithMessengerURL(url string) ClientOption {
	return func(c *BotClient) {
		c.messengerURL = url
	}
}

// WithHTTPClient برای تنظیم کلاینت HTTP سفارشی
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *BotClient) {
		c.httpClient = client
	}
}

// WithRateLimitDelay تنظیم تأخیر نرخ محدودیت
func WithRateLimitDelay(delay time.Duration) ClientOption {
	return func(c *BotClient) {
		c.rateLimitDuration = delay
	}
}

// WithMaxRetries تنظیم حداکثر تلاش مجدد
func WithMaxRetries(retries int) ClientOption {
	return func(c *BotClient) {
		c.maxRetries = retries
	}
}

// WithIgnoreTimeout تنظیم نادیده گرفتن خطاهای timeout
func WithIgnoreTimeout(ignore bool) ClientOption {
	return func(c *BotClient) {
		c.ignoreTimeout = ignore
	}
}

// WithHotReload فعال‌سازی قابلیت Hot-Reload
func WithHotReload(enabled bool) ClientOption {
	return func(c *BotClient) {
		c.hotReloadEnabled = enabled
	}
}

// WithConnectionMode تنظیم حالت اتصال (Bot API, Messenger API, Switcher)
func WithConnectionMode(mode ConnectionMode) ClientOption {
	return func(c *BotClient) {
		c.connectionMode = mode
	}
}

// WithNotificationOptions تنظیم گزینه‌های اطلاع‌رسانی سوییچ API
func WithNotificationOptions(opts NotificationOptions) ClientOption {
	return func(c *BotClient) {
		c.notificationOpts = &opts
	}
}

// WithFastConnection فعال‌سازی حالت اتصال سریع (زمان انتظار و پولینگ کوتاه‌تر)
func WithFastConnection() ClientOption {
	return func(c *BotClient) {
		// کاهش زمان انتظار کلی HTTP درخواست‌ها از 30 ثانیه به 10 ثانیه
		c.httpClient.Timeout = 10 * time.Second
		// کاهش فاصله زمانی بین درخواست‌ها برای سرعت بیشتر
		c.rateLimitDuration = 500 * time.Millisecond
	}
}

// WithUltraFastConnection فعال‌سازی حالت اتصال فوق سریع (با ریسک مسدود شدن)
func WithUltraFastConnection() ClientOption {
	return func(c *BotClient) {
		// کاهش زمان انتظار به 5 ثانیه
		c.httpClient.Timeout = 5 * time.Second
		// کاهش فاصله زمانی بین درخواست‌ها به 200 میلی‌ثانیه
		c.rateLimitDuration = 200 * time.Millisecond
	}
}

// NewClient ایجاد یک نمونه جدید از BotClient
func NewClient(token string, opts ...ClientOption) *BotClient {
	client := &BotClient{
		token:             token,
		baseURL:           defaultBaseURL,
		messengerURL:      defaultMessengerURL,
		rateLimitDuration: 1 * time.Second,
		maxRetries:        3,
		stopChan:          make(chan struct{}),
		ignoreTimeout:     true,
		hotReloadEnabled:  false,
		stateManager:      NewStateManager(),
		metadata:          make(map[string]interface{}),
		antiSpam:          NewAntiSpam(),
		connectionMode:    SwitcherMode,
		registeredUsers:   make([]*User, 0),
		usersFilePath:     "users.json",
		offsetsDir:        "offsets", // تعریف پوشه آفست‌ها
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	client.networkStabilityManager = NewNetworkStabilityManager(client)
	client.apiManager = NewAPIManager(client)
	client.reloadManager = NewReloadManager(client)
	client.hostReloadWatcher = NewHostReloadWatcher(client)

	// *** رفع باگ کلیدی: تنظیم API اولیه بر اساس حالت اتصال ***
	switch client.connectionMode {
	case BotAPIMode:
		client.apiManager.currentAPI = BotAPI
	case MessengerAPIMode:
		client.apiManager.currentAPI = MessengerAPI
	case SwitcherMode:
		client.apiManager.currentAPI = BotAPI // پیش‌فرض برای سوییچر
	}

	// بارگذاری کاربران
	client.loadUsersFromFile()

	return client
}

// ==================== Magic / Lifecycle Methods ====================

// init مقداردهی اولیه کلاینت
func (c *BotClient) init(ctx context.Context) error {
	if err := c.getBotID(ctx); err != nil {
		return fmt.Errorf("خطا در مقداردهی اولیه بات: %w", err)
	}

	// پس از دریافت botID، مسیر فایل آفست را تنظیم می‌کنیم
	if err := os.MkdirAll(c.offsetsDir, 0755); err != nil {
		return fmt.Errorf("خطا در ایجاد پوشه آفست‌ها: %w", err)
	}
	c.offsetFilePath = filepath.Join(c.offsetsDir, c.botID+".json")

	// بارگذاری آفست ذخیره شده
	c.GetOffset() // این متد آفست را از فایل می‌خواند و در حافظه نگه می‌دارد

	return nil
}

// Enter برای context manager
func (c *BotClient) Enter(ctx context.Context) (*BotClient, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

// Exit برای context manager
func (c *BotClient) Exit() error {
	return c.close()
}

// ==================== Core Run / Engine / Session Methods ====================

// Start شروع بات (نسخه 2.0.0)
func (c *BotClient) Start(ctx context.Context, opts ...interface{}) error {
	c.mu.Lock()
	if c.isRunning {
		c.mu.Unlock()
		return fmt.Errorf("بات در حال حاضر در حال اجراست")
	}
	c.isRunning = true
	c.mu.Unlock()

	log.Println("بات کار کرد")

	if c.connectionMode == SwitcherMode {
		go c.apiManager.StartHealthMonitoring(ctx, 30*time.Second)
	}

	if c.hotReloadEnabled {
		c.reloadManager.StartWatching()
	}

	if c.hostReloadWatcher != nil {
		if err := c.hostReloadWatcher.StartWatching(ctx); err != nil {
		}
	}

	if err := c.onStart(ctx); err != nil {
		return err
	}

	return nil
}

// Stop توقف بات
func (c *BotClient) Stop(ctx context.Context) error {
	c.mu.Lock()
	if !c.isRunning {
		c.mu.Unlock()
		return nil
	}
	c.isRunning = false
	close(c.stopChan)
	c.stopChan = make(chan struct{})
	c.mu.Unlock()

	c.apiManager.StopHealthMonitoring()

	if err := c.onShutdown(ctx); err != nil {
		return err
	}

	return c.close()
}

// Run اجرای بات با پولینگ
func (c *BotClient) Run(ctx context.Context, pollingOpts ...PollingOptions) error {
	if err := c.Start(ctx); err != nil {
		return err
	}

	if len(pollingOpts) > 0 {
		opts := pollingOpts[0]
		if opts.Handler == nil {
			opts.Handler = c.ProcessUpdate
		}
		return c.StartPolling(ctx, opts)
	}

	<-ctx.Done()
	return c.Stop(ctx)
}

// Connect اتصال به API (همان init)
func (c *BotClient) Connect(ctx context.Context) error {
	return c.init(ctx)
}

// Disconnect قطع اتصال
func (c *BotClient) Disconnect() error {
	return c.close()
}

// Close بستن کلاینت
func (c *BotClient) Close() error {
	return c.close()
}

// ==================== Update Handling / Dispatching Methods ====================

// GetUpdates دریافت آپدیت‌ها
func (c *BotClient) GetUpdates(ctx context.Context, offsetID *string, limit int) (*GetUpdatesResponse, error) {
	var offsetIDStr string
	if offsetID != nil {
		offsetIDStr = *offsetID
	}
	req := &GetUpdatesRequest{OffsetID: offsetIDStr, Limit: limit}
	var resp GetUpdatesResponse
	err := c.makeRequest(ctx, http.MethodPost, "getUpdates", req, &resp)
	return &resp, err
}

// ProcessUpdate پردازش یک آپدیت (متد عمومی)
func (c *BotClient) ProcessUpdate(ctx context.Context, update *Update) error {
	return c.processUpdate(ctx, update)
}

// ==================== Handler Management Methods ====================

// AddHandler اضافه کردن هندلر
func (c *BotClient) AddHandler(handler HandlerFunc, filter FilterFunc, order ...int) {
	handlerOrder := 0
	if len(order) > 0 {
		handlerOrder = order[0]
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, Handler{
		handler: handler,
		filter:  filter,
		order:   handlerOrder,
	})
}

// RemoveHandler حذف هندلر
func (c *BotClient) RemoveHandler(handler HandlerFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, h := range c.handlers {
		if fmt.Sprintf("%p", h.handler) == fmt.Sprintf("%p", handler) {
			c.handlers = append(c.handlers[:i], c.handlers[i+1:]...)
			break
		}
	}
}

// Middleware اضافه کردن میدلور
func (c *BotClient) Middleware(middleware MiddlewareFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.middlewares = append(c.middlewares, middleware)
}

// OnUpdate هندلر برای همه آپدیت‌ها
func (c *BotClient) OnUpdate(handler HandlerFunc) {
	c.AddHandler(handler, nil)
}

// OnStart هندلر برای شروع بات
func (c *BotClient) OnStart(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		return update.Type == StartedBot
	})
}

// OnShutdown هندلر برای خاموشی بات
func (c *BotClient) OnShutdown(handler HandlerFunc) {
	// در پیاده‌سازی کامل‌تر، این رویداد می‌تواند مدیریت شود
}

// OnChatUpdates هندلر برای آپدیت‌های چت
func (c *BotClient) OnChatUpdates(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		return update.ChatID != "" &&
			(update.Type == NewMessage || update.Type == UpdatedMessage)
	})
}

// OnMessageUpdates هندلر برای آپدیت‌های پیام
func (c *BotClient) OnMessageUpdates(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		return (update.NewMessage != nil || update.UpdatedMessage != nil) &&
			update.Type != StartedBot && update.Type != StoppedBot
	})
}

// OnPhoto هندلر برای پیام‌های عکس
func (c *BotClient) OnPhoto(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		if update.NewMessage != nil && update.NewMessage.File != nil {
			return update.NewMessage.File.Type == ImageType || update.NewMessage.File.Type == PhotoType
		}
		if update.UpdatedMessage != nil && update.UpdatedMessage.File != nil {
			return update.UpdatedMessage.File.Type == ImageType || update.UpdatedMessage.File.Type == PhotoType
		}
		return false
	})
}

// OnAudio هندلر برای پیام‌های صوتی
func (c *BotClient) OnAudio(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		if update.NewMessage != nil && update.NewMessage.File != nil {
			return update.NewMessage.File.Type == AudioType || update.NewMessage.File.Type == VoiceType
		}
		if update.UpdatedMessage != nil && update.UpdatedMessage.File != nil {
			return update.UpdatedMessage.File.Type == AudioType || update.UpdatedMessage.File.Type == VoiceType
		}
		return false
	})
}

// OnVideo هندلر برای پیام‌های ویدیویی
func (c *BotClient) OnVideo(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		if update.NewMessage != nil && update.NewMessage.File != nil {
			return update.NewMessage.File.Type == VideoType
		}
		if update.UpdatedMessage != nil && update.UpdatedMessage.File != nil {
			return update.UpdatedMessage.File.Type == VideoType
		}
		return false
	})
}

// OnDocument هندلر برای پیام‌های سند
func (c *BotClient) OnDocument(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		if update.NewMessage != nil && update.NewMessage.File != nil {
			return update.NewMessage.File.Type == FileType
		}
		if update.UpdatedMessage != nil && update.UpdatedMessage.File != nil {
			return update.UpdatedMessage.File.Type == FileType
		}
		return false
	})
}

// OnSticker هندلر برای پیام‌های استیکر
func (c *BotClient) OnSticker(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		if update.NewMessage != nil && update.NewMessage.Sticker != nil {
			return true
		}
		if update.UpdatedMessage != nil && update.UpdatedMessage.Sticker != nil {
			return true
		}
		return false
	})
}

// OnLocation هندلر برای پیام‌های موقعیت مکانی
func (c *BotClient) OnLocation(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		if update.NewMessage != nil && update.NewMessage.Location != nil {
			return true
		}
		if update.UpdatedMessage != nil && update.UpdatedMessage.Location != nil {
			return true
		}
		return false
	})
}

// OnContact هندلر برای پیام‌های مخاطب
func (c *BotClient) OnContact(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		if update.NewMessage != nil && update.NewMessage.ContactMessage != nil {
			return true
		}
		if update.UpdatedMessage != nil && update.UpdatedMessage.ContactMessage != nil {
			return true
		}
		return false
	})
}

// OnPoll هندلر برای پیام‌های نظرسنجی
func (c *BotClient) OnPoll(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		if update.NewMessage != nil && update.NewMessage.Poll != nil {
			return true
		}
		if update.UpdatedMessage != nil && update.UpdatedMessage.Poll != nil {
			return true
		}
		return false
	})
}

// OnCommand هندلر برای دستورات
func (c *BotClient) OnCommand(command string, handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		if update.NewMessage != nil && update.NewMessage.Text != "" {
			return update.NewMessage.Text == "/"+command ||
				(len(update.NewMessage.Text) > len(command)+1 &&
					update.NewMessage.Text[:len(command)+1] == "/"+command+" ")
		}
		if update.UpdatedMessage != nil && update.UpdatedMessage.Text != "" {
			return update.UpdatedMessage.Text == "/"+command ||
				(len(update.UpdatedMessage.Text) > len(command)+1 &&
					update.UpdatedMessage.Text[:len(command)+1] == "/"+command+" ")
		}
		return false
	})
}

// OnCallbackQuery هندلر برای کوئری‌های دکمه
func (c *BotClient) OnCallbackQuery(handler HandlerFunc) {
	c.AddHandler(handler, func(ctx context.Context, update *Update) bool {
		return update.NewMessage != nil && update.NewMessage.AuxData != nil && update.NewMessage.AuxData.ButtonID != nil
	})
}

// OnInlineQuery هندلر برای کوئری‌های اینلاین
func (c *BotClient) OnInlineQuery(handler HandlerFunc) {
	// دربخش مدل پیاده سازی شده
}

// ==================== Sending Messages / Media / Files Methods ====================

// SendMessage ارسال پیام متنی
func (c *BotClient) SendMessage(ctx context.Context, req *SendMessageRequest) (string, error) {
	c.applyRateLimit()
	var resp MessageIDResponse
	err := c.makeRequest(ctx, http.MethodPost, "sendMessage", req, &resp)
	return resp.MessageID, err
}

// EditMessageText ویرایش متن پیام
func (c *BotClient) EditMessageText(ctx context.Context, req *EditMessageTextRequest) error {
	c.applyRateLimit()
	return c.makeRequest(ctx, http.MethodPost, "editMessageText", req, nil)
}

// DeleteMessage حذف پیام
func (c *BotClient) DeleteMessage(ctx context.Context, req *DeleteMessageRequest) error {
	c.applyRateLimit()
	return c.makeRequest(ctx, http.MethodPost, "deleteMessage", req, nil)
}

// ForwardMessage فوروارد پیام
func (c *BotClient) ForwardMessage(ctx context.Context, req *ForwardMessageRequest) (string, error) {
	c.applyRateLimit()
	var resp ForwardMessageResponse
	err := c.makeRequest(ctx, http.MethodPost, "forwardMessage", req, &resp)
	return resp.NewMessageID, err
}

// ForwardMessages فوروارد چندین پیام
func (c *BotClient) ForwardMessages(ctx context.Context, fromChatID string, messageIDs []string, toChatID string) ([]string, error) {
	var newMessageIDs []string
	for _, msgID := range messageIDs {
		req := &ForwardMessageRequest{
			FromChatID: fromChatID,
			MessageID:  msgID,
			ToChatID:   toChatID,
		}
		newMsgID, err := c.ForwardMessage(ctx, req)
		if err != nil {
			return nil, err
		}
		newMessageIDs = append(newMessageIDs, newMsgID)
	}
	return newMessageIDs, nil
}

// SendPhoto ارسال عکس
func (c *BotClient) SendPhoto(ctx context.Context, chatID string, filePath string, caption ...string) (string, error) {
	fileID, err := c.uploadFile(ctx, filePath, ImageType)
	if err != nil {
		return "", err
	}

	req := &SendFileRequest{
		ChatID: chatID,
		FileID: fileID,
	}

	if len(caption) > 0 {
		req.Text = caption[0]
	}

	return c.SendFile(ctx, req)
}

// SendVideo ارسال ویدیو
func (c *BotClient) SendVideo(ctx context.Context, chatID string, filePath string, caption ...string) (string, error) {
	fileID, err := c.uploadFile(ctx, filePath, VideoType)
	if err != nil {
		return "", err
	}

	req := &SendFileRequest{
		ChatID: chatID,
		FileID: fileID,
	}

	if len(caption) > 0 {
		req.Text = caption[0]
	}

	return c.SendFile(ctx, req)
}

// SendMusic ارسال موزیک
func (c *BotClient) SendMusic(ctx context.Context, chatID string, filePath string, caption ...string) (string, error) {
	fileID, err := c.uploadFile(ctx, filePath, MusicType)
	if err != nil {
		return "", err
	}

	req := &SendFileRequest{
		ChatID: chatID,
		FileID: fileID,
	}

	if len(caption) > 0 {
		req.Text = caption[0]
	}

	return c.SendFile(ctx, req)
}

// SendDocument ارسال سند
func (c *BotClient) SendDocument(ctx context.Context, chatID string, filePath string, caption ...string) (string, error) {
	fileID, err := c.uploadFile(ctx, filePath, FileType)
	if err != nil {
		return "", err
	}

	req := &SendFileRequest{
		ChatID: chatID,
		FileID: fileID,
	}

	if len(caption) > 0 {
		req.Text = caption[0]
	}

	return c.SendFile(ctx, req)
}

// SendFile ارسال فایل
func (c *BotClient) SendFile(ctx context.Context, req *SendFileRequest) (string, error) {
	c.applyRateLimit()
	var resp MessageIDResponse
	err := c.makeRequest(ctx, http.MethodPost, "sendFile", req, &resp)
	return resp.MessageID, err
}

// SendVoice ارسال پیام صوتی
func (c *BotClient) SendVoice(ctx context.Context, chatID string, filePath string, caption ...string) (string, error) {
	fileID, err := c.uploadFile(ctx, filePath, VoiceType)
	if err != nil {
		return "", err
	}

	req := &SendFileRequest{
		ChatID: chatID,
		FileID: fileID,
	}

	if len(caption) > 0 {
		req.Text = caption[0]
	}

	return c.SendFile(ctx, req)
}

// SendGif ارسال گیف
func (c *BotClient) SendGif(ctx context.Context, chatID string, filePath string, caption ...string) (string, error) {
	fileID, err := c.uploadFile(ctx, filePath, GifType)
	if err != nil {
		return "", err
	}

	req := &SendFileRequest{
		ChatID: chatID,
		FileID: fileID,
	}

	if len(caption) > 0 {
		req.Text = caption[0]
	}

	return c.SendFile(ctx, req)
}

// SendSticker ارسال استیکر
func (c *BotClient) SendSticker(ctx context.Context, chatID string, filePath string) (string, error) {
	fileID, err := c.uploadFile(ctx, filePath, StickerType)
	if err != nil {
		return "", err
	}

	req := &SendFileRequest{
		ChatID: chatID,
		FileID: fileID,
	}

	return c.SendFile(ctx, req)
}

// SendLocation ارسال موقعیت مکانی
func (c *BotClient) SendLocation(ctx context.Context, req *SendLocationRequest) (string, error) {
	c.applyRateLimit()
	var resp MessageIDResponse
	err := c.makeRequest(ctx, http.MethodPost, "sendLocation", req, &resp)
	return resp.MessageID, err
}

// SendContact ارسال مخاطب
func (c *BotClient) SendContact(ctx context.Context, req *SendContactRequest) (string, error) {
	c.applyRateLimit()
	var resp MessageIDResponse
	err := c.makeRequest(ctx, http.MethodPost, "sendContact", req, &resp)
	return resp.MessageID, err
}

// SendPoll ارسال نظرسنجی
func (c *BotClient) SendPoll(ctx context.Context, req *SendPollRequest) (string, error) {
	c.applyRateLimit()
	var resp MessageIDResponse
	err := c.makeRequest(ctx, http.MethodPost, "sendPoll", req, &resp)
	return resp.MessageID, err
}

// SendChatActivity ارسال فعالیت چت
func (c *BotClient) SendChatActivity(ctx context.Context, chatID string, activity string) error {
	req := map[string]interface{}{
		"chat_id":  chatID,
		"activity": activity,
	}
	return c.makeRequest(ctx, http.MethodPost, "sendChatActivity", req, nil)
}

// ==================== Edit & UI Methods ====================

// EditMessage ویرایش پیام
func (c *BotClient) EditMessage(ctx context.Context, chatID string, messageID string, newText string) error {
	req := &EditMessageTextRequest{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      newText,
	}
	return c.EditMessageText(ctx, req)
}

// EditInlineKeypad ویرایش کیبورد پیام
func (c *BotClient) EditInlineKeypad(ctx context.Context, req *EditMessageKeypadRequest) error {
	c.applyRateLimit()
	return c.makeRequest(ctx, http.MethodPost, "editInlineKeypad", req, nil)
}

// EditChatKeypad ویرایش کیبورد چت
func (c *BotClient) EditChatKeypad(ctx context.Context, req *EditChatKeypadRequest) error {
	c.applyRateLimit()
	return c.makeRequest(ctx, http.MethodPost, "editChatKeypad", req, nil)
}

// SetPin پین کردن پیام
func (c *BotClient) SetPin(ctx context.Context, chatID string, messageID string) error {
	req := &PinChatMessageRequest{
		ChatID:    chatID,
		MessageID: messageID,
	}
	return c.PinChatMessage(ctx, req)
}

// SetUnpin آنپین کردن پیام
func (c *BotClient) SetUnpin(ctx context.Context, chatID string, messageID string) error {
	req := &UnpinChatMessageRequest{
		ChatID:    chatID,
		MessageID: messageID,
	}
	return c.UnpinChatMessage(ctx, req)
}

// ==================== Polls Methods ====================

// CreatePoll ایجاد نظرسنجی
func (c *BotClient) CreatePoll(ctx context.Context, chatID string, question string, options []string) (string, error) {
	req := &SendPollRequest{
		ChatID:   chatID,
		Question: question,
		Options:  options,
	}
	return c.SendPoll(ctx, req)
}

// VotePoll رأی دادن به نظرسنجی
func (c *BotClient) VotePoll(ctx context.Context, chatID string, messageID string, optionIndex int) error {
	req := map[string]interface{}{
		"chat_id":      chatID,
		"message_id":   messageID,
		"option_index": optionIndex,
	}
	return c.makeRequest(ctx, http.MethodPost, "votePoll", req, nil)
}

// GetPollStatus دریافت وضعیت نظرسنجی
func (c *BotClient) GetPollStatus(ctx context.Context, chatID string, messageID string) (*PollStatus, error) {
	var pollStatus PollStatus
	req := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
	}
	err := c.makeRequest(ctx, http.MethodPost, "getPollStatus", req, &pollStatus)
	return &pollStatus, err
}

// ==================== User / Profile Methods ====================

// GetMe دریافت اطلاعات بات
func (c *BotClient) GetMe(ctx context.Context) (*Bot, error) {
	var bot Bot
	err := c.makeRequest(ctx, http.MethodPost, "getMe", nil, &bot)
	return &bot, err
}

// GetBotInfo دریافت اطلاعات کامل بات
func (c *BotClient) GetBotInfo(ctx context.Context) (*Bot, error) {
	return c.GetMe(ctx)
}

// TerminateOtherSessions خاتمه جلسات دیگر
func (c *BotClient) TerminateOtherSessions(ctx context.Context) error {
	return c.makeRequest(ctx, http.MethodPost, "terminateOtherSessions", nil, nil)
}

// GetUserInfo دریافت اطلاعات کاربر
func (c *BotClient) GetUserInfo(ctx context.Context, userID string) (*User, error) {
	var user User
	req := map[string]interface{}{"user_id": userID}
	err := c.makeRequest(ctx, http.MethodPost, "getUserInfo", req, &user)
	return &user, err
}

// GetMembers دریافت اعضا
func (c *BotClient) GetMembers(ctx context.Context, chatID string) ([]ChatMember, error) {
	var members []ChatMember
	req := map[string]interface{}{"chat_id": chatID}
	err := c.makeRequest(ctx, http.MethodPost, "getChatMembers", req, &members)
	return members, err
}

// UpdateProfile به‌روزرسانی پروفایل
func (c *BotClient) UpdateProfile(ctx context.Context, firstName, lastName, bio string) error {
	req := map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
		"bio":        bio,
	}
	return c.makeRequest(ctx, http.MethodPost, "updateProfile", req, nil)
}

// UpdateUsername به‌روزرسانی نام کاربری
func (c *BotClient) UpdateUsername(ctx context.Context, username string) error {
	req := map[string]interface{}{"username": username}
	return c.makeRequest(ctx, http.MethodPost, "updateUsername", req, nil)
}

// UploadAvatar آپلود آواتار
func (c *BotClient) UploadAvatar(ctx context.Context, filePath string) (*File, error) {
	return c.UploadFileDirectly(ctx, filePath, ImageType)
}

// DeleteAvatar حذف آواتار
func (c *BotClient) DeleteAvatar(ctx context.Context) error {
	return c.makeRequest(ctx, http.MethodPost, "deleteAvatar", nil, nil)
}

// ==================== Chats / Groups / Channels (Admin) Methods ====================

// GetChat دریافت اطلاعات چت
func (c *BotClient) GetChat(ctx context.Context, chatID string) (*Chat, error) {
	var chat Chat
	err := c.makeRequest(ctx, http.MethodPost, "getChat", &GetChatRequest{ChatID: chatID}, &chat)
	return &chat, err
}

// GetChatMember دریافت اطلاعات عضو چت
func (c *BotClient) GetChatMember(ctx context.Context, req *GetChatMemberRequest) (*GetChatMemberResponse, error) {
	var resp GetChatMemberResponse
	err := c.makeRequest(ctx, http.MethodPost, "getChatMember", req, &resp)
	return &resp, err
}

// GetChatMemberCount دریافت تعداد اعضای چت
func (c *BotClient) GetChatMemberCount(ctx context.Context, req *GetChatMemberCountRequest) (*GetChatMemberCountResponse, error) {
	var resp GetChatMemberCountResponse
	err := c.makeRequest(ctx, http.MethodPost, "getChatMemberCount", req, &resp)
	return &resp, err
}

// GetChatAdministrators دریافت مدیران چت
func (c *BotClient) GetChatAdministrators(ctx context.Context, req *GetChatAdministratorsRequest) (*GetChatAdministratorsResponse, error) {
	var resp GetChatAdministratorsResponse
	err := c.makeRequest(ctx, http.MethodPost, "getChatAdministrators", req, &resp)
	return &resp, err
}

// SetChatPermissions تنظیم مجوزهای چت
func (c *BotClient) SetChatPermissions(ctx context.Context, chatID string, permissions map[string]bool) error {
	req := map[string]interface{}{
		"chat_id":     chatID,
		"permissions": permissions,
	}
	return c.makeRequest(ctx, http.MethodPost, "setChatPermissions", req, nil)
}

// BanChatMember مسدود کردن عضو چت
func (c *BotClient) BanChatMember(ctx context.Context, req *BanChatMemberRequest) error {
	return c.makeRequest(ctx, http.MethodPost, "banChatMember", req, nil)
}

// UnbanChatMember رفع مسدودیت عضو چت
func (c *BotClient) UnbanChatMember(ctx context.Context, req *UnbanChatMemberRequest) error {
	return c.makeRequest(ctx, http.MethodPost, "unbanChatMember", req, nil)
}

// PinChatMessage پین کردن پیام
func (c *BotClient) PinChatMessage(ctx context.Context, req *PinChatMessageRequest) error {
	return c.makeRequest(ctx, http.MethodPost, "pinChatMessage", req, nil)
}

// UnpinChatMessage آنپین کردن پیام
func (c *BotClient) UnpinChatMessage(ctx context.Context, req *UnpinChatMessageRequest) error {
	return c.makeRequest(ctx, http.MethodPost, "unpinChatMessage", req, nil)
}

// UnpinAllChatMessages آنپین کردن همه پیام‌ها
func (c *BotClient) UnpinAllChatMessages(ctx context.Context, chatID string) error {
	req := map[string]interface{}{"chat_id": chatID}
	return c.makeRequest(ctx, http.MethodPost, "unpinAllChatMessages", req, nil)
}

// PromoteChatMember ارتقای عضو چت
func (c *BotClient) PromoteChatMember(ctx context.Context, req *PromoteChatMemberRequest) error {
	return c.makeRequest(ctx, http.MethodPost, "promoteChatMember", req, nil)
}

// SetCommands تنظیم دستورات بات
func (c *BotClient) SetCommands(ctx context.Context, req *SetCommandsRequest) error {
	return c.makeRequest(ctx, http.MethodPost, "setCommands", req, nil)
}

// UpdateBotEndpoints به‌روزرسانی endpointهای بات
func (c *BotClient) UpdateBotEndpoints(ctx context.Context, url string, endpointType UpdateEndpointTypeEnum) error {
	req := map[string]interface{}{"url": url, "type": endpointType}
	return c.makeRequest(ctx, http.MethodPost, "updateBotEndpoints", req, nil)
}

// GetSelectionItem دریافت آیتم‌های انتخابی
func (c *BotClient) GetSelectionItem(ctx context.Context, req *GetSelectionItemRequest) (*GetSelectionItemResponse, error) {
	var resp GetSelectionItemResponse
	err := c.makeRequest(ctx, http.MethodPost, "getSelectionItem", req, &resp)
	return &resp, err
}

// SearchSelectionItems جستجوی آیتم‌های انتخابی
func (c *BotClient) SearchSelectionItems(ctx context.Context, req *SearchSelectionItemsRequest) (*SearchSelectionItemsResponse, error) {
	var resp SearchSelectionItemsResponse
	err := c.makeRequest(ctx, http.MethodPost, "searchSelectionItems", req, &resp)
	return &resp, err
}

// ==================== File Upload/Download Methods ====================

// RequestSendFile درخواست آدرس آپلود
func (c *BotClient) RequestSendFile(ctx context.Context, fileType FileTypeEnum) (*RequestSendFileResponse, error) {
	var resp RequestSendFileResponse
	err := c.makeRequest(ctx, http.MethodPost, "requestSendFile", &RequestSendFileRequest{Type: fileType}, &resp)
	return &resp, err
}

// UploadFile آپلود فایل
func (c *BotClient) UploadFile(uploadURL, filePath string) (*FileUploadResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("خطا در باز کردن فایل: %w", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("خطا در ساخت فرم فایل: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("خطا در کپی کردن محتوای فایل: %w", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", uploadURL, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("خطا در ایجاد درخواست آپلود: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("خطا در ارسال درخواست آپلود: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("آپلود با خطا %d مواجه شد: %s", resp.StatusCode, string(body))
	}

	var uploadResp FileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("خطا در پارس کردن پاسخ آپلود: %w", err)
	}
	return &uploadResp, nil
}

// Upload آپلود فایل (مخفف)
func (c *BotClient) Upload(ctx context.Context, filePath string, fileType FileTypeEnum) (*File, error) {
	return c.UploadFileDirectly(ctx, filePath, fileType)
}

// Download دانلود فایل
func (c *BotClient) Download(ctx context.Context, fileID string, savePath string) error {
	fileInfo, err := c.GetFile(ctx, fileID)
	if err != nil {
		return err
	}

	resp, err := http.Get(fileInfo.DownloadURL)
	if err != nil {
		return fmt.Errorf("خطا در دانلود فایل: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("خطا در دریافت فایل: وضعیت %d", resp.StatusCode)
	}

	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("خطا در ایجاد فایل خروجی: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// DownloadFile دانلود فایل
func (c *BotClient) DownloadFile(ctx context.Context, fileID string, savePath string) error {
	return c.Download(ctx, fileID, savePath)
}

// DownloadProfilePicture دانلود عکس پروفایل
func (c *BotClient) DownloadProfilePicture(ctx context.Context, userID string, savePath string) error {
	userInfo, err := c.GetUserInfo(ctx, userID)
	if err != nil {
		return err
	}

	if userInfo.Avatar == nil {
		return fmt.Errorf("کاربر آواتار ندارد")
	}

	return c.Download(ctx, userInfo.Avatar.FileID, savePath)
}

// ==================== متدهای کمکی و داخلی ====================

// getBotID دریافت شناسه بات (منطق جدید و مقاوم)
func (c *BotClient) getBotID(ctx context.Context) error {
	if c.botID != "" {
		return nil
	}

	var lastErr error

	// لیست APIها برای تلاش، با API فعلی شروع می‌شود
	apisToTry := []APIType{c.apiManager.GetCurrentAPI()}

	// اگر در حالت سوییچر هستیم، API دیگر را هم به لیست اضافه می‌کنیم
	if c.connectionMode == SwitcherMode {
		if apisToTry[0] == BotAPI {
			apisToTry = append(apisToTry, MessengerAPI)
		} else {
			apisToTry = append(apisToTry, BotAPI)
		}
	}

	// تلاش برای دریافت شناسه بات از APIهای موجود
	for _, apiType := range apisToTry {
		// سوییچ به API مورد نظر
		c.apiManager.SwitchAPI(apiType, "تلاش برای دریافت شناسه بات")

		bot, err := c.GetMe(ctx)
		if err == nil {
			c.botID = bot.BotID
			return nil // موفقیت، خروج از تابع
		}

		lastErr = err
	}

	// اگر هیچ‌کدام از APIها جواب ندادند
	return fmt.Errorf("خطا در دریافت شناسه بات پس از تلاش با تمام APIهای موجود: %w", lastErr)
}

// makeRequest ارسال درخواست به API (نسخه 2.0.0 - هوشمند)
func (c *BotClient) makeRequest(ctx context.Context, method string, apiMethod string, reqBody, respBody any) error {
	if err := c.ensureSession(); err != nil {
		return err
	}

	var baseURL string
	currentAPI := c.apiManager.GetCurrentAPI()
	if currentAPI == BotAPI {
		baseURL = c.baseURL
	} else {
		baseURL = c.messengerURL
	}

	url := fmt.Sprintf("%s/%s/%s", baseURL, c.token, apiMethod)
	startTime := time.Now()

	var bodyReader io.Reader
	if reqBody != nil {
		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("خطا در ساخت درخواست JSON: %w", err)
		}

		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("خطا در ایجاد درخواست HTTP: %w", err)
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", userAgent)

	var lastErr error
	for i := 0; i < c.maxRetries; i++ {
		resp, err := c.httpClient.Do(req)
		responseTime := time.Since(startTime)

		if err != nil {
			lastErr = fmt.Errorf("خطا در ارسال درخواست: %w", err)
			c.networkStabilityManager.RecordEvent(currentAPI, responseTime, false, lastErr)

			if c.connectionMode == SwitcherMode && i == c.maxRetries-1 {
				go c.apiManager.CheckHealth(ctx, currentAPI)
			}

			if i < c.maxRetries-1 && c.networkStabilityManager.IsRetryableError(err) {
				delay := c.networkStabilityManager.CalculateBackoffDelay(i)
				time.Sleep(delay)
				continue
			}
			return lastErr
		}

		respData, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = fmt.Errorf("خطا در خواندن پاسخ: %w", err)
			c.networkStabilityManager.RecordEvent(currentAPI, responseTime, false, lastErr)

			if c.connectionMode == SwitcherMode && i == c.maxRetries-1 {
				go c.apiManager.CheckHealth(ctx, currentAPI)
			}

			if i < c.maxRetries-1 && c.isRetryableStatus(resp.StatusCode) {
				delay := c.networkStabilityManager.CalculateBackoffDelay(i)
				time.Sleep(delay)
				continue
			}
			return lastErr
		}

		if resp.StatusCode != http.StatusOK {
			var apiErr APIError
			apiErr.StatusCode = resp.StatusCode
			if json.Unmarshal(respData, &apiErr) != nil || apiErr.Message == "" {
				apiErr.Message = string(respData)
			}
			lastErr = &apiErr
			c.networkStabilityManager.RecordEvent(currentAPI, responseTime, false, lastErr)

			if resp.StatusCode == http.StatusRequestTimeout && c.ignoreTimeout {
				if i < c.maxRetries-1 {
					delay := c.networkStabilityManager.CalculateBackoffDelay(i)
					time.Sleep(delay)
					continue
				}
			}

			if c.connectionMode == SwitcherMode && (resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests) {
				go c.apiManager.CheckHealth(ctx, currentAPI)
			}

			if i < c.maxRetries-1 && c.isRetryableStatus(resp.StatusCode) {
				delay := c.networkStabilityManager.CalculateBackoffDelay(i)
				time.Sleep(delay)
				continue
			}
			return lastErr
		}

		c.networkStabilityManager.RecordEvent(currentAPI, responseTime, true, nil)

		var baseResponse struct {
			Data   json.RawMessage `json:"data"`
			Status string          `json:"status"`
		}

		if err := json.Unmarshal(respData, &baseResponse); err != nil {
			if respBody != nil {
				if err := json.Unmarshal(respData, respBody); err != nil {
					return fmt.Errorf("خطا در پارس کردن پاسخ کامل: %w", err)
				}
			}
			return nil
		}

		if baseResponse.Status != "OK" {
			return fmt.Errorf("API returned non-OK status: %s", baseResponse.Status)
		}

		if respBody != nil && baseResponse.Data != nil {
			if err := json.Unmarshal(baseResponse.Data, respBody); err != nil {
				return fmt.Errorf("خطا در پارس کردن داده‌های پاسخ: %w", err)
			}
		}

		return nil
	}

	return lastErr
}

// uploadFile آپلود فایل و دریافت fileID
func (c *BotClient) uploadFile(ctx context.Context, filePath string, fileType FileTypeEnum) (string, error) {
	uploadResp, err := c.RequestSendFile(ctx, fileType)
	if err != nil {
		return "", fmt.Errorf("خطا در درخواست آدرس آپلود: %w", err)
	}

	fileResp, err := c.UploadFile(uploadResp.UploadURL, filePath)
	if err != nil {
		return "", fmt.Errorf("خطا در آپلود فایل: %w", err)
	}

	return fileResp.FileID, nil
}

// processUpdate پردازش آپدیت (متد خصوصی)
func (c *BotClient) processUpdate(ctx context.Context, update *Update) error {
	if !c.filtersPass(ctx, update) {
		return nil
	}
	return c.dispatchUpdate(ctx, update)
}

// onStart اجرای هندلرهای شروع
func (c *BotClient) onStart(ctx context.Context) error {
	startUpdate := &Update{
		Type: StartedBot,
	}
	return c.dispatchUpdate(ctx, startUpdate)
}

// onShutdown اجرای هندلرهای خاموشی
func (c *BotClient) onShutdown(ctx context.Context) error {
	shutdownUpdate := &Update{
		Type: StoppedBot,
	}
	return c.dispatchUpdate(ctx, shutdownUpdate)
}

// close بستن کلاینت
func (c *BotClient) close() error {
	c.closeSession()
	return nil
}

// GetFile دریافت اطلاعات فایل
func (c *BotClient) GetFile(ctx context.Context, fileID string) (*GetFileResponse, error) {
	var resp GetFileResponse
	err := c.makeRequest(ctx, http.MethodPost, "getFile", &GetFileRequest{FileID: fileID}, &resp)
	return &resp, err
}

// UploadFileDirectly آپلود فایل و دریافت اطلاعات کامل
func (c *BotClient) UploadFileDirectly(ctx context.Context, filePath string, fileType FileTypeEnum) (*File, error) {
	fileID, err := c.uploadFile(ctx, filePath, fileType)
	if err != nil {
		return nil, err
	}

	fileInfo, err := c.GetFile(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت اطلاعات فایل آپلود شده: %w", err)
	}

	return &File{
		FileID:   fileID,
		FileName: fileInfo.FileName,
		Size:     fileInfo.Size,
		Type:     fileType,
	}, nil
}

// ==================== متدهای کمکی و داخلی (ادامه) ====================

// ensureSession اطمینان از وجود session
func (c *BotClient) ensureSession() error {
	c.mu.RLock()
	hasSession := c.httpClient != nil
	c.mu.RUnlock()

	if !hasSession {
		c.mu.Lock()
		c.httpClient = &http.Client{Timeout: 30 * time.Second}
		c.mu.Unlock()
	}
	return nil
}

// closeSession بستن session
func (c *BotClient) closeSession() {
	c.mu.Lock()
	if c.httpClient != nil {
		c.httpClient.CloseIdleConnections()
		c.httpClient = nil
	}
	c.mu.Unlock()
}

// sleepBackoff تأخیر نمایی برای تلاش مجدد (منسوخ شده، از NetworkStabilityManager استفاده کنید)
func (c *BotClient) sleepBackoff(retryCount int) {
	delay := time.Duration(retryCount*retryCount) * time.Second
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}
	time.Sleep(delay)
}

// applyRateLimit مدیریت نرخ ارسال درخواست
func (c *BotClient) applyRateLimit() {
	c.mu.Lock()
	lastSent := c.lastSentTime
	c.mu.Unlock()

	if time.Since(lastSent) < c.rateLimitDuration {
		time.Sleep(c.rateLimitDuration - time.Since(lastSent))
	}

	c.mu.Lock()
	c.lastSentTime = time.Now()
	c.mu.Unlock()
}

// isRetryableStatus بررسی امکان تلاش مجدد
func (c *BotClient) isRetryableStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode >= 500 ||
		statusCode == 0 // خطای شبکه
}

// dispatchUpdate ارسال آپدیت به هندلرها (اصلاح شده)
func (c *BotClient) dispatchUpdate(ctx context.Context, update *Update) error {
	var handler HandlerFunc = func(ctx context.Context, update *Update) error {
		for _, h := range c.handlers {
			if h.filter == nil || h.filter(ctx, update) {
				if err := h.handler(ctx, update); err != nil {
					// *** این خط اضافه شده است ***
					log.Printf("❌ خطا در اجرای هندلر برای آپدیت %s: %v", update.Type, err)
				}
			}
		}
		return nil
	}

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		middleware := c.middlewares[i]
		next := handler
		handler = func(ctx context.Context, update *Update) error {
			return middleware(ctx, update, next)
		}
	}

	return handler(ctx, update)
}

// parseUpdate پارس کردن آپدیت (اضافه شده)
func (c *BotClient) parseUpdate(data []byte) (*Update, error) {
	var update Update
	if err := json.Unmarshal(data, &update); err != nil {
		return nil, fmt.Errorf("خطا در پارس کردن آپدیت: %w", err)
	}
	return &update, nil
}

// extractMessageID استخراج شناسه پیام (اضافه شده)
func (c *BotClient) extractMessageID(update *Update) string {
	if update.NewMessage != nil {
		return update.NewMessage.MessageID
	}
	if update.UpdatedMessage != nil {
		return update.UpdatedMessage.MessageID
	}
	if update.RemovedMessageID != nil {
		return *update.RemovedMessageID
	}
	return ""
}

// filtersPass بررسی فیلترها (اضافه شده)
func (c *BotClient) filtersPass(ctx context.Context, update *Update) bool {
	for _, filter := range c.updateFilters {
		if !filter(ctx, update) {
			return false
		}
	}
	return true
}

// ==================== متدهای مربوط به Hot-Reload ====================

// EnableHotReload فعال‌سازی قابلیت Hot-Reload
func (c *BotClient) EnableHotReload() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hotReloadEnabled = true
}

// DisableHotReload غیرفعال‌سازی قابلیت Hot-Reload
func (c *BotClient) DisableHotReload() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hotReloadEnabled = false
}

// IsHotReloadEnabled بررسی وضعیت Hot-Reload
func (c *BotClient) IsHotReloadEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hotReloadEnabled
}

// GetReloadManager دسترسی به مدیر Hot-Reload را فراهم می‌کند
func (c *BotClient) GetReloadManager() *ReloadManager {
	return c.reloadManager
}

// ==================== متدهای مربوط به State Management ====================

// SetState ذخیره وضعیت برای کاربر
func (c *BotClient) SetState(userID, key string, value interface{}) {
	c.stateManager.SetState(userID, key, value)
}

// GetState دریافت وضعیت برای کاربر
func (c *BotClient) GetState(userID, key string) (interface{}, bool) {
	return c.stateManager.GetState(userID, key)
}

// DeleteState حذف وضعیت برای کاربر
func (c *BotClient) DeleteState(userID, key string) {
	c.stateManager.DeleteState(userID, key)
}

// DeleteUserState حذف تمام وضعیت‌های کاربر
func (c *BotClient) DeleteUserState(userID string) {
	c.stateManager.DeleteUserState(userID)
}

// ==================== متدهای مربوط به Anti-Spam ====================

// CheckAntiSpam بررسی ضد اسپم
func (c *BotClient) CheckAntiSpam(userID string) bool {
	return c.antiSpam.Check(userID)
}

// ResetAntiSpam بازنشانی ضد اسپم برای کاربر
func (c *BotClient) ResetAntiSpam(userID string) {
	c.antiSpam.Reset(userID)
}

// ==================== متدهای مربوط به Metadata ====================

// SetMetadata تنظیم متادیتا
func (c *BotClient) SetMetadata(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metadata[key] = value
}

// GetMetadata دریافت متادیتا
func (c *BotClient) GetMetadata(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.metadata[key]
	return value, exists
}

// DeleteMetadata حذف متادیتا
func (c *BotClient) DeleteMetadata(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.metadata, key)
}

// ==================== متدهای مربوط به API Manager ====================

// GetAPIManager دریافت مدیر API
func (c *BotClient) GetAPIManager() *APIManager {
	return c.apiManager
}

// GetConnectionMode دریافت حالت اتصال
func (c *BotClient) GetConnectionMode() ConnectionMode {
	return c.connectionMode
}

// SetNotificationOptions تنظیم گزینه‌های اطلاع‌رسانی
func (c *BotClient) SetNotificationOptions(opts NotificationOptions) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.notificationOpts = &opts
}

// sendAPIChangeNotification ارسال اطلاع‌رسانی تغییر API
func (c *BotClient) sendAPIChangeNotification(message string) {
	if c.notificationOpts == nil || !c.notificationOpts.Enabled {
		return
	}

	if c.notificationOpts.SendToAll {
		if c.notificationOpts.ChatID != "" {
			_, err := c.SendMessage(context.Background(), &SendMessageRequest{
				ChatID: c.notificationOpts.ChatID,
				Text:   message,
			})
			if err != nil {
			}
		}
		return
	}

	if c.notificationOpts.SendToOwner && c.notificationOpts.ChatID != "" {
		_, err := c.SendMessage(context.Background(), &SendMessageRequest{
			ChatID: c.notificationOpts.ChatID,
			Text:   message,
		})
		if err != nil {
		}
	}
}

// ==================== متدهای مربوط به ذخیره اطلاعات کاربران ====================

// loadUsersFromFile بارگذاری لیست کاربران از فایل JSON در زمان راه‌اندازی
func (c *BotClient) loadUsersFromFile() {
	c.fileMu.Lock()
	defer c.fileMu.Unlock()

	data, err := os.ReadFile(c.usersFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
		}
		// اگر فایل وجود ندارد< با یک لیست خالی شروع میشود
		return
	}

	err = json.Unmarshal(data, &c.registeredUsers)
	if err != nil {
		c.registeredUsers = make([]*User, 0) // اگر فایل خراب است، لیست را خالی کن
	}
}

// SaveUserInfoToFile ذخیره اطلاعات یک کاربر جدید در فایل JSON
func (c *BotClient) SaveUserInfoToFile(user *User) error {
	c.fileMu.Lock()
	defer c.fileMu.Unlock()

	// بررسی اینکه آیا کاربر از قبل وجود دارد یا خیر
	for _, existingUser := range c.registeredUsers {
		if existingUser.UserID == user.UserID {
			return nil // کاربر تکراری است، کاری نکن
		}
	}

	// افزودن کاربر جدید به لیست
	c.registeredUsers = append(c.registeredUsers, user)

	// تبدیل کل لیست به JSON با فرمت خوانا
	data, err := json.MarshalIndent(c.registeredUsers, "", "  ")
	if err != nil {
		return fmt.Errorf("خطا در تبدیل داده‌های کاربران به JSON: %w", err)
	}

	// نوشتن داده‌ها در فایل
	err = os.WriteFile(c.usersFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("خطا در نوشتن فایل کاربران: %w", err)
	}

	return nil
}

// ==================== متدهای مربوط به ذخیره آفست (اصلاح شده) ====================

// offsetData ساختار داده برای ذخیره آفست در فایل JSON
type offsetData struct {
	Offset string `json:"offset"`
}

// loadOffsetFromFile بارگذاری آفست از فایل
func (c *BotClient) loadOffsetFromFile() string {
	// اگر مسیر فایل هنوز تنظیم نشده (مثلاً قبل از init)، مقدار خالی برمی‌گردان
	if c.offsetFilePath == "" {
		log.Printf("اخطار: مسیر فایل آفست هنوز تنظیم نشده است.")
		return ""
	}

	data, err := os.ReadFile(c.offsetFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			// فقط خطاهای واقعی فایل سیستم را لاگ می‌کنیم
			log.Printf("خطا در خواندن فایل آفست %s: %v", c.offsetFilePath, err)
		}
		return ""
	}

	var dataStruct offsetData
	if err := json.Unmarshal(data, &dataStruct); err != nil {
		log.Printf("خطا در پارس کردن فایل آفست %s: %v. فایل بازنویسی خواهد شد.", c.offsetFilePath, err)
		// در صورت خطا، فایل را با یک آفست خالی بازنویسی می‌کنیم تا خرابی بعدی رخ ندهد
		c.saveOffsetToFile("")
		return ""
	}

	return dataStruct.Offset
}

// saveOffsetToFile ذخیره آفست در فایل
func (c *BotClient) saveOffsetToFile(offset string) error {
	if c.offsetFilePath == "" {
		// این خطا نباید هرگز رخ دهد، اما برای اطمینان چک می‌کنیم
		return fmt.Errorf("مسیر فایل آفست تنظیم نشده است")
	}

	dataStruct := offsetData{Offset: offset}
	data, err := json.MarshalIndent(dataStruct, "", "  ")
	if err != nil {
		return fmt.Errorf("خطا در تبدیل آفست به JSON: %w", err)
	}

	// استفاده از 0600 برای امنیت بیشتر فایل
	err = os.WriteFile(c.offsetFilePath, data, 0600)
	if err != nil {
		// این خطا حیاتی است و باید لاگ شود
		log.Printf("خطا در نوشتن فایل آفست %s: %v", c.offsetFilePath, err)
		return fmt.Errorf("خطا در نوشتن فایل آفست: %w", err)
	}

	return nil
}

// GetOffset دریافت آخرین آفست ذخیره شده
func (c *BotClient) GetOffset() string {
	return c.loadOffsetFromFile()
}

// SetOffset تنظیم و ذخیره آفست جدید
func (c *BotClient) SetOffset(offset string) error {
	return c.saveOffsetToFile(offset)
}

// ==================== متدهای مربوط به مارک‌داون و HTML ====================

// SendMessageWithMarkdown ارسال پیام متنی با پشتیبانی از مارک‌داون (نسخه نهایی و صحیح)
func (c *BotClient) SendMessageWithMarkdown(ctx context.Context, chatID string, markdownText string, options ...*SendMessageRequest) (string, error) {
	parser := &TextParser{}
	parsed, err := parser.CheckMarkdown(markdownText)
	if err != nil {
		return "", err
	}

	var req *SendMessageRequest
	if len(options) > 0 && options[0] != nil {
		req = options[0]
	} else {
		req = &SendMessageRequest{ChatID: chatID}
	}

	req.Text = parsed.Text

	if len(parsed.Metadata) > 0 {
		req.Metadata = map[string]interface{}{
			"meta_data_parts": parsed.Metadata,
		}
	}

	return c.SendMessage(ctx, req)
}

// SendMessageWithHTML ارسال پیام متنی با پشتیبانی از HTML
func (c *BotClient) SendMessageWithHTML(ctx context.Context, chatID string, htmlText string, options ...*SendMessageRequest) (string, error) {
	parser := &TextParser{}
	parsed, err := parser.CheckHTML(htmlText)
	if err != nil {
		return "", err
	}

	var req *SendMessageRequest
	if len(options) > 0 && options[0] != nil {
		req = options[0]
	} else {
		req = &SendMessageRequest{ChatID: chatID}
	}

	req.Text = parsed.Text

	if len(parsed.Metadata) > 0 {
		req.Metadata = map[string]interface{}{
			"meta_data_parts": parsed.Metadata,
		}
	}

	return c.SendMessage(ctx, req)
}

// EditMessageTextWithMarkdown ویرایش متن پیام با پشتیبانی از مارک‌داون
func (c *BotClient) EditMessageTextWithMarkdown(ctx context.Context, chatID string, messageID string, markdownText string, options ...*EditMessageTextRequest) error {
	parser := &TextParser{}
	parsed, err := parser.CheckMarkdown(markdownText)
	if err != nil {
		return err
	}

	var req *EditMessageTextRequest
	if len(options) > 0 && options[0] != nil {
		req = options[0]
	} else {
		req = &EditMessageTextRequest{
			ChatID:    chatID,
			MessageID: messageID,
		}
	}

	req.Text = parsed.Text
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["meta_data_parts"] = parsed.Metadata

	return c.EditMessageText(ctx, req)
}

// EditMessageTextWithHTML ویرایش متن پیام با پشتیبانی از HTML
func (c *BotClient) EditMessageTextWithHTML(ctx context.Context, chatID string, messageID string, htmlText string, options ...*EditMessageTextRequest) error {
	parser := &TextParser{}
	parsed, err := parser.CheckHTML(htmlText)
	if err != nil {
		return err
	}

	var req *EditMessageTextRequest
	if len(options) > 0 && options[0] != nil {
		req = options[0]
	} else {
		req = &EditMessageTextRequest{
			ChatID:    chatID,
			MessageID: messageID,
		}
	}

	req.Text = parsed.Text
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["meta_data_parts"] = parsed.Metadata

	return c.EditMessageText(ctx, req)
}

// SendPhotoWithMarkdownCaption ارسال عکس با کپشن مارک‌داون
func (c *BotClient) SendPhotoWithMarkdownCaption(ctx context.Context, chatID string, filePath string, markdownCaption string, options ...*SendFileRequest) (string, error) {
	parser := &TextParser{}
	parsed, err := parser.CheckMarkdown(markdownCaption)
	if err != nil {
		return "", err
	}

	fileID, err := c.uploadFile(ctx, filePath, ImageType)
	if err != nil {
		return "", err
	}

	var req *SendFileRequest
	if len(options) > 0 && options[0] != nil {
		req = options[0]
		req.FileID = fileID
	} else {
		req = &SendFileRequest{
			ChatID: chatID,
			FileID: fileID,
		}
	}

	req.Text = parsed.Text
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["meta_data_parts"] = parsed.Metadata

	return c.SendFile(ctx, req)
}

// SendPhotoWithHTMLCaption ارسال عکس با کپشن HTML
func (c *BotClient) SendPhotoWithHTMLCaption(ctx context.Context, chatID string, filePath string, htmlCaption string, options ...*SendFileRequest) (string, error) {
	parser := &TextParser{}
	parsed, err := parser.CheckHTML(htmlCaption)
	if err != nil {
		return "", err
	}

	fileID, err := c.uploadFile(ctx, filePath, ImageType)
	if err != nil {
		return "", err
	}

	var req *SendFileRequest
	if len(options) > 0 && options[0] != nil {
		req = options[0]
		req.FileID = fileID
	} else {
		req = &SendFileRequest{
			ChatID: chatID,
			FileID: fileID,
		}
	}

	req.Text = parsed.Text
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["meta_data_parts"] = parsed.Metadata

	return c.SendFile(ctx, req)
}

// SendVideoWithMarkdownCaption ارسال ویدیو با کپشن مارک‌داون
func (c *BotClient) SendVideoWithMarkdownCaption(ctx context.Context, chatID string, filePath string, markdownCaption string, options ...*SendFileRequest) (string, error) {
	parser := &TextParser{}
	parsed, err := parser.CheckMarkdown(markdownCaption)
	if err != nil {
		return "", err
	}

	fileID, err := c.uploadFile(ctx, filePath, VideoType)
	if err != nil {
		return "", err
	}

	var req *SendFileRequest
	if len(options) > 0 && options[0] != nil {
		req = options[0]
		req.FileID = fileID
	} else {
		req = &SendFileRequest{
			ChatID: chatID,
			FileID: fileID,
		}
	}

	req.Text = parsed.Text
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["meta_data_parts"] = parsed.Metadata

	return c.SendFile(ctx, req)
}

// SendVideoWithHTMLCaption ارسال ویدیو با کپشن HTML
func (c *BotClient) SendVideoWithHTMLCaption(ctx context.Context, chatID string, filePath string, htmlCaption string, options ...*SendFileRequest) (string, error) {
	parser := &TextParser{}
	parsed, err := parser.CheckHTML(htmlCaption)
	if err != nil {
		return "", err
	}

	fileID, err := c.uploadFile(ctx, filePath, VideoType)
	if err != nil {
		return "", err
	}

	var req *SendFileRequest
	if len(options) > 0 && options[0] != nil {
		req = options[0]
		req.FileID = fileID
	} else {
		req = &SendFileRequest{
			ChatID: chatID,
			FileID: fileID,
		}
	}

	req.Text = parsed.Text
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["meta_data_parts"] = parsed.Metadata

	return c.SendFile(ctx, req)
}

// SendDocumentWithMarkdownCaption ارسال سند با کپشن مارک‌داون
func (c *BotClient) SendDocumentWithMarkdownCaption(ctx context.Context, chatID string, filePath string, markdownCaption string, options ...*SendFileRequest) (string, error) {
	parser := &TextParser{}
	parsed, err := parser.CheckMarkdown(markdownCaption)
	if err != nil {
		return "", err
	}

	fileID, err := c.uploadFile(ctx, filePath, FileType)
	if err != nil {
		return "", err
	}

	var req *SendFileRequest
	if len(options) > 0 && options[0] != nil {
		req = options[0]
		req.FileID = fileID
	} else {
		req = &SendFileRequest{
			ChatID: chatID,
			FileID: fileID,
		}
	}

	req.Text = parsed.Text
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["meta_data_parts"] = parsed.Metadata

	return c.SendFile(ctx, req)
}

// SendDocumentWithHTMLCaption ارسال سند با کپشن HTML
func (c *BotClient) SendDocumentWithHTMLCaption(ctx context.Context, chatID string, filePath string, htmlCaption string, options ...*SendFileRequest) (string, error) {
	parser := &TextParser{}
	parsed, err := parser.CheckHTML(htmlCaption)
	if err != nil {
		return "", err
	}

	fileID, err := c.uploadFile(ctx, filePath, FileType)
	if err != nil {
		return "", err
	}

	var req *SendFileRequest
	if len(options) > 0 && options[0] != nil {
		req = options[0]
		req.FileID = fileID
	} else {
		req = &SendFileRequest{
			ChatID: chatID,
			FileID: fileID,
		}
	}

	req.Text = parsed.Text
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["meta_data_parts"] = parsed.Metadata

	return c.SendFile(ctx, req)
}
