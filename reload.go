package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ReloadManager مدیریت قابلیت بارگذاری مجدد (Hot-Reload) به صورت واقعی
type ReloadManager struct {
	client      *BotClient
	isWatching  bool
	watcher     *fsnotify.Watcher
	mu          sync.Mutex
	reloadFuncs []func() // لیستی از توابعی که پس از هر بارگذاری مجدد اجرا می‌شوند
	lastReload  time.Time
}

// NewReloadManager ایجاد یک نمونه جدید از ReloadManager
func NewReloadManager(client *BotClient) *ReloadManager {
	return &ReloadManager{
		client:      client,
		reloadFuncs: make([]func(), 0),
	}
}

// StartWatching شروع نظارت بر تغییرات فایل‌ها
func (rm *ReloadManager) StartWatching() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.isWatching {
		log.Println("Hot-Reload در حال حاضر در حال اجراست.")
		return
	}

	if !rm.client.IsHotReloadEnabled() {
		log.Println("Hot-Reload در کلاینت غیرفعال است. برای فعال‌سازی از گزینه WithHotReload استفاده کنید.")
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("خطا در ایجاد watcher: %v", err)
		return
	}

	rm.watcher = watcher
	rm.isWatching = true
	log.Println("Hot-Reload watcher شروع به کار کرد.")

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			if strings.Contains(path, "vendor/") || strings.Contains(path, ".git/") {
				return nil
			}

			err = watcher.Add(path)
			if err != nil {
				log.Printf("خطا در افزودن فایل به watcher: %s, خطا: %v", path, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("خطا در راه‌اندازی نظارت بر فایل‌ها: %v", err)
		rm.isWatching = false
		watcher.Close()
		return
	}

	go rm.watchLoop()
}

// StopWatching توقف نظارت بر تغییرات فایل‌ها
func (rm *ReloadManager) StopWatching() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if !rm.isWatching || rm.watcher == nil {
		return
	}

	rm.isWatching = false
	err := rm.watcher.Close()
	rm.watcher = nil

	log.Println("Hot-Reload watcher متوقف شد.")
	if err != nil {
		log.Printf("خطا در بستن watcher: %v", err)
	}
}

// watchLoop حلقه نظارت بر تغییرات
func (rm *ReloadManager) watchLoop() {
	debounceTimer := time.NewTimer(0)
	<-debounceTimer.C

	var changedFiles []string

	for {
		select {
		case event, ok := <-rm.watcher.Events:
			if !ok {
				return
			}

			if !strings.HasSuffix(event.Name, ".go") {
				continue
			}

			if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove {
				continue
			}

			changedFiles = append(changedFiles, event.Name)

			debounceTimer.Stop()
			debounceTimer = time.AfterFunc(2*time.Second, func() {
				if len(changedFiles) > 0 {
					rm.handleFileChanges(changedFiles)
					changedFiles = nil
				}
			})

		case err, ok := <-rm.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("خطا در Hot-Reload watcher: %v", err)
		}
	}
}

// handleFileChanges پردازش تغییرات فایل‌ها
func (rm *ReloadManager) handleFileChanges(files []string) {
	log.Printf("تغییرات در فایل‌های زیر شناسایی شد: %v", files)

	if time.Since(rm.lastReload) < 5*time.Second {
		log.Println("بارگذاری مجدد بیش از حد، نادیده گرفته شد.")
		return
	}

	rm.lastReload = time.Now()

	if err := rm.rebuildProject(); err != nil {
		log.Printf("خطا در کامپایل مجدد پروژه: %v", err)
		return
	}

	log.Println("پروژه با موفقیت کامپایل شد.")

	for _, fn := range rm.reloadFuncs {
		if fn != nil {
			fn()
		}
	}

	log.Println("بارگذاری مجدد (Hot-Reload) با موفقیت انجام شد.")
}

// rebuildProject کامپایل مجدد پروژه
func (rm *ReloadManager) rebuildProject() error {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("فایل go.mod یافت نشد")
	}

	cmd := exec.Command("go", "build", "-o", "tmp_binary")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("خطا در کامپایل: %w", err)
	}

	defer os.Remove("tmp_binary")

	return nil
}

// TriggerReload اجرای فرآیند بارگذاری مجدد
func (rm *ReloadManager) TriggerReload() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	log.Println("بارگذاری مجدد (Hot-Reload) فعال شد...")

	for _, fn := range rm.reloadFuncs {
		if fn != nil {
			fn()
		}
	}

	log.Println("بارگذاری مجدد با موفقیت انجام شد.")
}

// OnReload ثبت یک تابع برای اجرا پس از هر بارگذاری مجدد
func (rm *ReloadManager) OnReload(fn func()) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.reloadFuncs = append(rm.reloadFuncs, fn)
}
