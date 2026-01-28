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
	os.MkdirAll("offsets", 0755)

	client := ParsRubika.NewClient(token)

	// Handle all text messages (Echo Bot)
	client.OnMessageUpdates(func(ctx context.Context, update *ParsRubika.Update) error {
		if update.NewMessage == nil {
			return nil
		}

		// Ignore commands (they start with /)
		if len(update.NewMessage.Text) > 0 && update.NewMessage.Text[0] == '/' {
			return nil
		}

		// Echo the text back
		_, err := client.SendMessage(ctx, &ParsRubika.SendMessageRequest{
			ChatID: update.ChatID,
			Text:   "You said: " + update.NewMessage.Text,
		})
		return err
	})

	ctx := context.Background()
	log.Println("Echo Bot running...")

	client.Run(ctx, ParsRubika.PollingOptions{
		Limit:        100,
		PollInterval: 2 * time.Second,
	})
}
