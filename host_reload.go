package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

import (
	"context"
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

// HostReloadWatcher نظارت بر تغییرات فایل‌ها برای Host-Reload واقعی
type HostReloadWatcher struct {
	client   *BotClient
	watcher  *fsnotify.Watcher
	isActive bool
	mu       sync.Mutex
}

// NewHostReloadWatcher ایجاد یک نمونه جدید از HostReloadWatcher
func NewHostReloadWatcher(client *BotClient) *HostReloadWatcher {
	return &HostReloadWatcher{
		client: client,
	}
}

// StartWatching شروع نظارت بر تغییرات فایل‌ها
func (hrw *HostReloadWatcher) StartWatching(ctx context.Context) error {
	hrw.mu.Lock()
	defer hrw.mu.Unlock()

	if hrw.isActive {
		log.Println("Host-Reload در حال حاضر فعال است.")
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("خطا در ایجاد watcher: %w", err)
	}

	hrw.watcher = watcher
	hrw.isActive = true

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
		return fmt.Errorf("خطا در راه‌اندازی نظارت بر فایل‌ها: %w", err)
	}

	log.Println("Host-Reload watcher فعال شد.")

	go hrw.watchLoop(ctx)

	return nil
}

// StopWatching توقف نظارت بر تغییرات فایل‌ها
func (hrw *HostReloadWatcher) StopWatching() error {
	hrw.mu.Lock()
	defer hrw.mu.Unlock()

	if !hrw.isActive || hrw.watcher == nil {
		return nil
	}

	hrw.isActive = false
	err := hrw.watcher.Close()
	hrw.watcher = nil

	log.Println("Host-Reload watcher غیرفعال شد.")
	return err
}

// watchLoop حلقه نظارت بر تغییرات
func (hrw *HostReloadWatcher) watchLoop(ctx context.Context) {
	debounceTimer := time.NewTimer(0)
	<-debounceTimer.C

	var changedFiles []string

	for {
		select {
		case <-ctx.Done():
			hrw.StopWatching()
			return

		case event, ok := <-hrw.watcher.Events:
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
					hrw.handleFileChanges(changedFiles)
					changedFiles = nil
				}
			})

		case err, ok := <-hrw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("خطا در Host-Reload watcher: %v", err)
		}
	}
}

// handleFileChanges پردازش تغییرات فایل‌ها
func (hrw *HostReloadWatcher) handleFileChanges(files []string) {
	log.Printf("Host-Reload: تغییرات در فایل‌های زیر شناسایی شد: %v", files)

	if err := hrw.rebuildAndRestart(); err != nil {
		log.Printf("Host-Reload: خطا در کامپایل و ری‌استارت پروژه: %v", err)
		return
	}

	log.Println("Host-Reload: پروژه با موفقیت ری‌استارت شد.")
}

// rebuildAndRestart کامپایل و ری‌استارت کامل برنامه (منطق جدید برای ویندوز)
func (hrw *HostReloadWatcher) rebuildAndRestart() error {
	// 1. بررسی وجود فایل go.mod
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("فایل go.mod یافت نشد، این یک پروژه Go نیست")
	}

	// 2. دریافت مسیر فایل اجرایی فعلی
	currentExecutablePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("خطا در دریافت مسیر فایل اجرایی: %w", err)
	}

	// 3. تعیین نام فایل‌های باینری قدیمی و جدید
	newExecutablePath := currentExecutablePath + "_new"
	oldExecutablePath := currentExecutablePath + "_old"

	// 4. کامپایل پروژه به فایل جدید
	log.Printf("Host-Reload: در حال کامپایل پروژه به %s...", newExecutablePath)
	cmd := exec.Command("go", "build", "-o", newExecutablePath, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		// پاک کردن فایل جدید در صورت شکست در کامپایل
		os.Remove(newExecutablePath)
		return fmt.Errorf("خطا در کامپایل: %w", err)
	}
	log.Println("Host-Reload: پروژه با موفقیت کامپایل شد.")

	// 5. جایگزینی اتمی فایل قدیمی با فایل جدید (برای ویندوز)
	// ابتدا فایل فعلی را به نام قدیمی منتقل می‌کنیم
	log.Printf("Host-Reload: در حال انتقال فایل قدیمی به %s...", oldExecutablePath)
	if err := os.Rename(currentExecutablePath, oldExecutablePath); err != nil {
		os.Remove(newExecutablePath)
		return fmt.Errorf("خطا در انتقال فایل قدیمی: %w", err)
	}

	// سپس فایل جدید را به نام اصلی منتقل می‌کنیم
	log.Printf("Host-Reload: در حال جایگزینی با فایل جدید...")
	if err := os.Rename(newExecutablePath, currentExecutablePath); err != nil {
		// تلاش برای بازگرداندن فایل قدیمی در صورت شکست
		log.Printf("Host-Reload: شکست در جایگزینی. تلاش برای بازگرداندن فایل قدیمی...")
		os.Rename(oldExecutablePath, currentExecutablePath)
		os.Remove(newExecutablePath)
		return fmt.Errorf("خطا در جایگزینی فایل جدید: %w", err)
	}
	log.Println("Host-Reload: فایل باینری با موفقیت جایگزین شد.")

	// 6. راه‌اندازی مجدد برنامه و خروج از پروسه فعلی
	log.Println("Host-Reload: در حال راه‌اندازی مجدد برنامه...")

	args := os.Args
	newProcess := exec.Command(currentExecutablePath, args[1:]...)
	newProcess.Stdout = os.Stdout
	newProcess.Stderr = os.Stderr
	newProcess.Stdin = os.Stdin

	err = newProcess.Start()
	if err != nil {
		log.Printf("Host-Reload: خطا در راه‌اندازی مجدد برنامه: %v", err)
		return fmt.Errorf("خطا در راه‌اندازی مجدد برنامه: %w", err)
	}

	log.Println("Host-Reload: برنامه جدید با موفقیت آغاز شد. برنامه فعلی در حال خروج است.")

	// 7. خروج از پروسه فعلی تا برنامه جدید جایگزین آن شود
	os.Exit(0)

	return nil // این کد هرگز اجرا نخواهد شد
}
