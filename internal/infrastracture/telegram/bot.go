package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot() (*Bot, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{api: api}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return nil
		case upd := <-updates:
			if upd.Message == nil {
				continue
			}
			go b.handleMessage(ctx, upd)
		}
	}
}

const (
	btnJoin = "Присоединиться к группе"
	btnSkip = "Не присоединяться к группе"
)

func (b *Bot) handleMessage(ctx context.Context, upd tgbotapi.Update) {
	m := upd.Message
	chatID := m.Chat.ID
	if m.IsCommand() {
		switch m.Command() {
		case "start":
			if err := b.handleStart(chatID); err != nil {
				log.Printf("handlestart failed chat=%d: %v", chatID, err)
			}
			return
		}
	}

	switch strings.TrimSpace(m.Text) {
	case btnJoin:
		msg := tgbotapi.NewMessage(chatID, "Введи номер своей группы (например, 23204).")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		_, _ = b.api.Send(msg)
	}

}

var ErrTooManyRequests = fmt.Errorf("too many requests")

func (b *Bot) sendWithRetry(msg tgbotapi.Chattable) error {
	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		_, err := b.api.Send(msg)
		if err == nil {
			return nil
		}

		// Логируем
		log.Printf("send failed (attempt %d/%d): %v", attempt, maxAttempts, err)
		time.Sleep(200 * time.Millisecond)
	}

	return ErrTooManyRequests
}
