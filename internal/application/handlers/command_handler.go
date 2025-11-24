package handlers

import (
	"context"
	"log"
	"strings"
	"timetable-homework-tgbot/internal/application/controllers"
	"timetable-homework-tgbot/internal/infrastracture/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	auth controllers.AuthController
	bot  *telegram.Bot
}

func NewCommandHandler(auth controllers.AuthController, bot *telegram.Bot) *CommandHandler {
	return &CommandHandler{auth: auth, bot: bot}
}

func (h *CommandHandler) Start(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "Хотите присоединиться к своей группе в НГУ?")
	msg.ReplyMarkup = telegram.KBAskJoin()
	_ = h.bot.SendWithRetry(msg)
}

func (h *CommandHandler) Join(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	joined, _ := h.auth.EnsureJoined(ctx, userID)
	if joined {
		_ = h.bot.Send(chatID, "Ты уже присоединён к группе.", telegram.KBMember())
		return
	}
	h.bot.State.Set(chatID, telegram.StateWaitUserGroup)
	_ = h.bot.SendRemove(chatID, "Введи номер своей группы (например, 23204).")
}

func (h *CommandHandler) Leave(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	h.bot.State.Del(chatID)
	if err := h.auth.LeaveGroup(ctx, userID); err != nil {
		log.Printf("leave failed: %v", err)
	}
	_ = h.bot.Send(chatID, "Вы отсоединены от группы.", telegram.KBGuest())
}

func (h *CommandHandler) WaitUserGroup(ctx context.Context, u tgbotapi.Update) {
	m := u.Message
	chatID, userID := m.Chat.ID, m.From.ID
	group := strings.TrimSpace(m.Text)

	if group == "" {
		_ = h.bot.Send(chatID, "Пусто. Введи номер своей группы (например, 23204).", telegram.KBMember())
		return
	}

	if err := h.auth.JoinGroup(ctx, userID, group); err != nil {
		log.Printf("JoinGroup failed: %v", err)
		h.bot.State.Del(chatID)
		_ = h.bot.Send(chatID, "Не удалось присоединиться. Проверь номер группы и попробуй ещё раз.", telegram.KBGuest())
		return
	}

	h.bot.State.Del(chatID)
	_ = h.bot.Send(chatID, "Группа сохранена: "+group, telegram.KBMember())
}

func (h *CommandHandler) Skip(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	h.bot.State.Del(chatID)
	_ = h.bot.Send(chatID, "Ок, продолжим без привязки.", telegram.KBGuest())
}
