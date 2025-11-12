package handlers

import (
	"context"
	"strings"
	"timetable-homework-tgbot/internal/infrastracture/controllers"
	"timetable-homework-tgbot/internal/infrastracture/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HWHandler struct {
	hw  controllers.HomeworkController
	bot *telegram.Bot
}

func NewHWHandler(hw controllers.HomeworkController, bot *telegram.Bot) *HWHandler {
	return &HWHandler{hw: hw, bot: bot}
}

func (h *HWHandler) PinStart(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	days, err := h.hw.DaysWithLessons(ctx, userID)
	if err != nil || len(days) == 0 {
		_ = h.bot.Send(chatID, "Нет ближайших занятий или вы не присоединены к группе.", telegram.KBMember())
		return
	}
	h.bot.State.Set(chatID, telegram.StateWaitHWDay)
	_ = h.bot.Send(chatID, "К какому дню недели прикрепить ДЗ?", telegram.KBDays(days))
}

func (h *HWHandler) WaitDay(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	day := strings.TrimSpace(u.Message.Text)
	lessons, err := h.hw.LessonsByDay(ctx, userID, day)
	if err != nil || len(lessons) == 0 {
		_ = h.bot.Send(chatID, "В этот день пар нет. Выбери другой.", telegram.KBMember())
		return
	}
	h.bot.HWSessSet(chatID, telegram.HwSession{Day: day})
	curr := h.bot.State.Get(chatID)
	if strings.HasPrefix(curr, telegram.StateWaitHWEditDay) {
		h.bot.State.Set(chatID, telegram.StateWaitHWEditLesson)
	} else {
		h.bot.State.Set(chatID, telegram.StateWaitHWLesson)
	}
	_ = h.bot.Send(chatID, "Выбери пару:", telegram.KBLessons(lessons))
}

func (h *HWHandler) WaitLesson(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	id, ok := telegram.ExtractIDFromLabel(strings.TrimSpace(u.Message.Text))
	if !ok {
		_ = h.bot.Send(chatID, "Нажми кнопку ещё раз.", telegram.KBMember())
		return
	}
	s := h.bot.HWSessGet(chatID)
	s.LessonID = id
	h.bot.HWSessSet(chatID, s)
	curr := h.bot.State.Get(chatID)
	if strings.HasPrefix(curr, telegram.StateWaitHWEditLesson) {
		h.bot.State.Set(chatID, telegram.StateWaitHWTextEdit)
	} else {
		h.bot.State.Set(chatID, telegram.StateWaitHWText)
	}
	_ = h.bot.SendRemove(chatID, "Введи текст ДЗ:")
}

func (h *HWHandler) WaitText(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	text := strings.TrimSpace(u.Message.Text)
	s := h.bot.HWSessGet(chatID)
	if s.Day == "" || s.LessonID == "" {
		_ = h.bot.Send(chatID, "Сессия потерялась, начни заново.", telegram.KBMember())
		h.bot.State.Del(chatID)
		return
	}
	_ = h.hw.Pin(ctx, userID, s.Day, s.LessonID, text) // TODO: обработать err
	h.bot.State.Del(chatID)
	h.bot.HWSessDel(chatID)
	_ = h.bot.Send(chatID, "Домашнее задание сохранено ✅", telegram.KBMember())
}

func (h *HWHandler) EditStart(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	days, err := h.hw.DaysWithLessons(ctx, userID)
	if err != nil || len(days) == 0 {
		_ = h.bot.Send(chatID, "Нет доступных дней.", telegram.KBMember())
		return
	}
	h.bot.State.Set(chatID, telegram.StateWaitHWEditDay)
	_ = h.bot.Send(chatID, "Выбери день для редактирования ДЗ:", telegram.KBDays(days))
}

func (h *HWHandler) WaitTextEdit(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	newText := strings.TrimSpace(u.Message.Text)
	if h.bot.State.Get(chatID) != telegram.StateWaitHWTextEdit {
		_ = h.bot.Send(chatID, "Сессия потерялась, начни заново.", telegram.KBMember())
		h.bot.State.Del(chatID)
		return
	}
	s := h.bot.HWSessGet(chatID)
	if s.LessonID == "" {
		_ = h.bot.Send(chatID, "Не удалось определить пару.", telegram.KBMember())
		h.bot.State.Del(chatID)
		return
	}
	_ = h.hw.Update(ctx, userID, s.LessonID, newText) // TODO: обработать err
	h.bot.State.Del(chatID)
	h.bot.HWSessDel(chatID)
	_ = h.bot.Send(chatID, "Домашнее задание обновлено ✅", telegram.KBMember())
}
