package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StateWaitUserGroup        = "wait_user_group"
	StateWaitGroupTB          = "wait_tb_group"
	StateWaitTeacherTB        = "wait_tb_teacher"
	StateWaitRoomTB           = "wait_tb_room"
	StateWaitHWDay            = "wait_hw_day"
	StateWaitHWLesson         = "wait_hw_lesson"
	StateWaitHWText           = "wait_hw_text"
	StateWaitConfirmDelete    = "wait_confirm_delete"
	StateWaitHWEditLesson     = "wait_edit_hw_lesson"
	StateWaitHWTextEdit       = "wait_edit_hw_text"
	StateWaitHWTable          = "wait_hw_table"
	StateWaitHWTableToDelete  = "wait_hw_table_to_delete"
	StateWaitRemindChooseHW   = "wait_remind_hw"
	StateWaitRemindChooseDay  = "wait_remind_day"
	StateWaitRemindChooseTime = "wait_remind_time"
	StateWaitRemindChoose     = "wait_remind_choose"
)

var ErrTooManyRequests = fmt.Errorf("too many requests")

type HwSession struct {
	Day         string
	LessonTitle string
}

type RemindSession struct {
	SubjectWithTask string
	Date            string
	TimeHHMM        string
}

type Bot struct {
	api    *tgbotapi.BotAPI
	state  StateStore
	router *Router

	hwMu sync.Mutex
	hw   map[int64]HwSession

	remMu sync.Mutex
	rem   map[int64]RemindSession
}

func NewBot(api *tgbotapi.BotAPI, state StateStore) *Bot {
	return &Bot{
		api: api, state: state, router: NewRouter(),
		hw: map[int64]HwSession{}, rem: map[int64]RemindSession{},
	}
}

func (b *Bot) Router() *Router { return b.router }

func (b *Bot) GetState() StateStore { return b.state }

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

func (b *Bot) handleMessage(parent context.Context, upd tgbotapi.Update) {
	m := upd.Message
	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()

	chatID := m.Chat.ID
	text := strings.TrimSpace(m.Text)

	if st := b.state.Get(chatID); st != "" && !m.IsCommand() {
		if h, ok := matchState(b.router, st); ok {
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

	log.Printf("default: state=%q text=%q", b.state.Get(chatID), text)
	b.router.def(ctx, upd)
}

func (b *Bot) Send(chatID int64, text string, kb tgbotapi.ReplyKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = kb
	_, err := b.api.Send(msg)
	return err
}
func (b *Bot) SendRemove(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) HWSessSet(chatID int64, s HwSession) { b.hwMu.Lock(); b.hw[chatID] = s; b.hwMu.Unlock() }
func (b *Bot) HWSessGet(chatID int64) HwSession {
	b.hwMu.Lock()
	defer b.hwMu.Unlock()
	return b.hw[chatID]
}
func (b *Bot) HWSessDel(chatID int64) { b.hwMu.Lock(); delete(b.hw, chatID); b.hwMu.Unlock() }

func (b *Bot) RemSessSet(chatID int64, s RemindSession) {
	b.remMu.Lock()
	b.rem[chatID] = s
	b.remMu.Unlock()
}
func (b *Bot) RemSessGet(chatID int64) (RemindSession, bool) {
	b.remMu.Lock()
	defer b.remMu.Unlock()
	s, ok := b.rem[chatID]
	return s, ok
}
func (b *Bot) RemSessDel(chatID int64) { b.remMu.Lock(); delete(b.rem, chatID); b.remMu.Unlock() }

func (b *Bot) SendWithRetry(msg tgbotapi.Chattable) error {
	const maxAttempts = 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		_, err := b.api.Send(msg)
		if err == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return ErrTooManyRequests
}
