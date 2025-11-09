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
	api    *tgbotapi.BotAPI
	state  StateStore
	router *Router
}

func NewBot(state StateStore) (*Bot, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	b := &Bot{api: api, state: state, router: NewRouter()}
	b.registerRoutes()

	return b, nil
}

func (b *Bot) registerRoutes() {
	b.router.OnCommand("start", b.hStart)
	b.router.OnText(btnJoin, b.hJoin)
	b.router.OnText(btnSkip, b.hSkip)
	b.router.OnState(stateWait, b.hWaitGroup)
	b.router.Default(b.hDefault)
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
	btnJoin   = "Присоединиться к группе"
	btnSkip   = "Не присоединяться к группе"
	stateWait = "wait_group"
)

func (b *Bot) handleMessage(parent context.Context, upd tgbotapi.Update) {
	m := upd.Message
	if m == nil {
		return
	}

	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	chatID := m.Chat.ID
	text := strings.TrimSpace(m.Text)

	if st := b.state.Get(chatID); st != "" && !m.IsCommand() {
		if h, ok := b.router.state[st]; ok {
			h(ctx, upd)
			return
		}
	}

	if m.IsCommand() {
		if h, ok := b.router.cmd[m.Command()]; ok {
			h(ctx, upd)
			return
		}
	}

	if h, ok := b.router.text[text]; ok {
		h(ctx, upd)
		return
	}

	b.router.def(ctx, upd)
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
