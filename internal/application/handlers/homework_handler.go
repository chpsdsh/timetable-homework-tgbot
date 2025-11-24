package handlers

import (
	"context"
	"log"
	"strings"
	"timetable-homework-tgbot/internal/application/controllers"
	"timetable-homework-tgbot/internal/application/formatter"
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
		h.bot.State.Del(chatID)
		_ = h.bot.Send(chatID, "В этот день пар нет. Выбери другой.", telegram.KBMember())
		return
	}
	h.bot.HWSessSet(chatID, telegram.HwSession{Day: day})

	h.bot.State.Set(chatID, telegram.StateWaitHWLesson)

	_ = h.bot.Send(chatID, "Выбери пару:", telegram.KBLessons(lessons))
}

func (h *HWHandler) WaitLesson(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	args := strings.Split(strings.TrimSpace(u.Message.Text), "-")
	lesson := strings.TrimSpace(args[0])

	s := h.bot.HWSessGet(chatID)
	s.LessonTitle = lesson
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
	if s.Day == "" || s.LessonTitle == "" {
		_ = h.bot.Send(chatID, "Сессия потерялась, начни заново.", telegram.KBMember())
		h.bot.State.Del(chatID)
		return
	}
	if err := h.hw.Pin(ctx, userID, s.Day, s.LessonTitle, text); err != nil {
		_ = h.bot.Send(chatID, "Не удалось прикрепить домашнее задание", telegram.KBMember())
		h.bot.State.Del(chatID)
		return
	}
	h.bot.State.Del(chatID)
	h.bot.HWSessDel(chatID)
	_ = h.bot.Send(chatID, "Домашнее задание сохранено ✅", telegram.KBMember())
}

func (h *HWHandler) EditStart(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	hw, err := h.hw.ListForLastWeek(ctx, userID)
	if err != nil || len(hw) == 0 {
		h.bot.State.Del(chatID)
		_ = h.bot.Send(chatID, "Нет домашек для редактирования.", telegram.KBMember())
		return
	}
	h.bot.State.Set(chatID, telegram.StateWaitHWTable)
	_ = h.bot.Send(chatID, "Выбери день для редактирования ДЗ:", telegram.KBHomeworks(hw))
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
	if s.LessonTitle == "" {
		_ = h.bot.Send(chatID, "Не удалось определить пару.", telegram.KBMember())
		h.bot.State.Del(chatID)
		return
	}

	if err := h.hw.Update(ctx, userID, s.LessonTitle, newText); err != nil {
		_ = h.bot.Send(chatID, "Не удалось обновить домашнее задание", telegram.KBMember())
		h.bot.State.Del(chatID)
		return
	}

	h.bot.State.Del(chatID)
	h.bot.HWSessDel(chatID)
	_ = h.bot.Send(chatID, "Домашнее задание обновлено ✅", telegram.KBMember())
}

func (h *HWHandler) ListHomeworks(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	hw, err := h.hw.ListForLastWeek(ctx, userID)
	if err != nil || len(hw) == 0 {
		h.bot.State.Del(chatID)
		_ = h.bot.Send(chatID, "Нет домашних заданий", telegram.KBMember())
		return
	}
	switch u.Message.Text {
	case telegram.BtnWatchHomeworks:
		formHw := formatter.FormatHomeworks(hw)
		h.bot.State.Del(chatID)
		_ = h.bot.Send(chatID, "Список домашек:\n "+formHw, telegram.KBMember())
	case telegram.BtnDeleteHomeworks:
		list, err := h.hw.ListForLastWeek(ctx, userID)
		if err != nil || len(list) == 0 {
			log.Println(err)
			h.bot.State.Del(chatID)
			_ = h.bot.Send(chatID, "Нет домашек для удаления.\n ", telegram.KBMember())
			return
		}

		h.bot.State.Set(chatID, telegram.StateWaitHWTableToDelete)
		_ = h.bot.Send(chatID, "Выбери домашку для удаления:\n ", telegram.KBHomeworks(hw))
	case telegram.BtnUpdateHomeworkStatus:
		h.bot.State.Set(chatID, telegram.StateWaitHomeworkUpdateChoose)
		_ = h.bot.Send(chatID, "Выбери домашку для обновления статуса:\n ", telegram.KBHomeworks(hw))
	}
}

func (h *HWHandler) WaitHomeWorkTable(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	userID := u.Message.From.ID
	s := h.bot.HWSessGet(chatID)
	arr := strings.Split(strings.TrimSpace(u.Message.Text), ":")
	s.LessonTitle = strings.TrimSpace(arr[0])
	log.Println(s.LessonTitle)
	exist, err := h.hw.CheckExistence(ctx, userID, s.LessonTitle)
	if err != nil {
		log.Println(err)
		h.bot.HWSessDel(chatID)
		h.bot.State.Del(chatID)
		return
	}
	if !exist {
		h.bot.State.Del(chatID)
		h.bot.HWSessDel(chatID)
		_ = h.bot.Send(chatID, "Некорректное домашнее задание", telegram.KBMember())
		return
	}

	h.bot.HWSessSet(chatID, s)
	log.Println(h.bot.State.Get(chatID))
	switch h.bot.State.Get(chatID) {
	case telegram.StateWaitHWTable:
		h.bot.State.Set(chatID, telegram.StateWaitHWTextEdit)
		_ = h.bot.SendRemove(chatID, "Введи новый текст ДЗ:")
	case telegram.StateWaitHWTableToDelete:
		h.bot.State.Set(chatID, telegram.StateWaitConfirmDelete)
		_ = h.bot.Send(chatID, "Подтвердите удаление:", telegram.KBConfirmDelete())
	case telegram.StateWaitHomeworkUpdateChoose:
		if err := h.hw.UpdateStatus(ctx, userID, s.LessonTitle); err != nil {
			log.Println(err)
			h.bot.HWSessDel(chatID)
			h.bot.State.Del(chatID)
			return
		}
		h.bot.HWSessDel(chatID)
		h.bot.State.Del(chatID)
		_ = h.bot.Send(chatID, "Статус домашнего задания обновлен", telegram.KBMember())

	}
}

func (h *HWHandler) WaitConfirmDelete(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	userID := u.Message.From.ID
	s := h.bot.HWSessGet(chatID)
	if u.Message.Text == telegram.BtnDelete {
		if err := h.hw.DeleteHomework(ctx, userID, s.LessonTitle); err != nil {
			log.Println(err)
			h.bot.State.Del(chatID)
			return
		}
		h.bot.State.Del(chatID)
		h.bot.HWSessDel(chatID)
		_ = h.bot.Send(chatID, "Домашнее задание удалено ✅", telegram.KBMember())
	} else if u.Message.Text == telegram.BtnNotDelete {
		h.bot.State.Del(chatID)
		h.bot.HWSessDel(chatID)
		_ = h.bot.Send(chatID, "Домашнее задание оставлено ✅", telegram.KBMember())
	}

}
