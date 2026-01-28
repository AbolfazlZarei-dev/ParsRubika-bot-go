package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v2"
)

func main() {
	token := "YOUR_BOT_TOKEN_HERE"

	// ایجاد پوشه آفست‌ها
	os.MkdirAll("offsets", 0755)

	log.Println("Starting Client...")
	client := ParsRubika.NewClient(token)

	// Middleware
	client.Middleware(func(ctx context.Context, update *ParsRubika.Update, next ParsRubika.HandlerFunc) error {
		if update.NewMessage != nil {
			log.Printf("📩 Message: %s", update.NewMessage.Text)
		}
		return next(ctx, update)
	})

	// ===============================
	// ارسال عکس (Photo)
	// ===============================
	client.OnCommand("photo", func(ctx context.Context, update *ParsRubika.Update) error {
		// فایل عکس را مسیر دهید. می‌توانید مسیر کامل (C:/path/to/image.jpg) هم بدهید.
		// در اینجا فرض شده که فایل کنار فایل exe است.
		filePath := "image.jpg"

		_, err := client.SendPhoto(ctx, update.ChatID, filePath, "This is a photo caption! 📷")
		if err != nil {
			return err
		}
		log.Printf("Photo sent: %s", filePath)
		return nil
	})

	// ===============================
	// ارسال ویدیو (Video)
	// ===============================
	client.OnCommand("video", func(ctx context.Context, update *ParsRubika.Update) error {
		filePath := "video.mp4"

		_, err := client.SendVideo(ctx, update.ChatID, filePath, "This is a video! 🎥")
		if err != nil {
			return err
		}
		log.Printf("Video sent: %s", filePath)
		return nil
	})

	// ===============================
	// ارسال سند/فایل (Document)
	// ===============================
	client.OnCommand("doc", func(ctx context.Context, update *ParsRubika.Update) error {
		filePath := "document.pdf"

		_, err := client.SendDocument(ctx, update.ChatID, filePath, "This is a document file! 📄")
		if err != nil {
			return err
		}
		log.Printf("Document sent: %s", filePath)
		return nil
	})

	// ===============================
	// ارسال موزیک (Audio/Music)
	// ===============================
	client.OnCommand("audio", func(ctx context.Context, update *ParsRubika.Update) error {
		filePath := "audio.mp3"

		_, err := client.SendMusic(ctx, update.ChatID, filePath, "This is a music file! 🎵")
		if err != nil {
			return err
		}
		log.Printf("Music sent: %s", filePath)
		return nil
	})

	// ===============================
	// ارسال صدای ضبط شده (Voice)
	// ===============================
	client.OnCommand("voice", func(ctx context.Context, update *ParsRubika.Update) error {
		filePath := "voice.ogg" // یا ogg

		_, err := client.SendVoice(ctx, update.ChatID, filePath, "This is a voice note! 🎙️")
		if err != nil {
			return err
		}
		log.Printf("Voice sent: %s", filePath)
		return nil
	})

	// ===============================
	// ارسال استیکر (Sticker)
	// ===============================
	client.OnCommand("sticker", func(ctx context.Context, update *ParsRubika.Update) error {
		filePath := "sticker.webp" // فرمت وب پی معمولا برای استیکر است

		_, err := client.SendSticker(ctx, update.ChatID, filePath)
		if err != nil {
			return err
		}
		log.Printf("Sticker sent: %s", filePath)
		return nil
	})

	// ===============================
	// راهنما (Help)
	// ===============================
	client.OnCommand("help", func(ctx context.Context, update *ParsRubika.Update) error {
		text := "Send Media Commands:\n\n" +
			"/photo - Send a photo (need image.jpg)\n" +
			"/video - Send a video (need video.mp4)\n" +
			"/doc - Send a document (need document.pdf)\n" +
			"/audio - Send music (need audio.mp3)\n" +
			"/voice - Send voice note (need voice.ogg)\n" +
			"/sticker - Send sticker (need sticker.webp)"

		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   text,
		})
		return err
	})

	ctx := context.Background()
	log.Println("Media Bot is running...")
	log.Println("Make sure you have test files (image.jpg, video.mp4, etc.) in the folder!")

	if err := client.Run(ctx, ParsRubika.PollingOptions{
		Limit:        100,
		PollInterval: 2 * time.Second,
	}); err != nil {
		log.Fatalf("❌ Error: %v", err)
	}
}
