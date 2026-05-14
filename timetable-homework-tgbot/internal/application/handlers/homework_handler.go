package handlers

import (
	"context"
	"log"
	"strconv"
	"strings"
	"timetable-homework-tgbot/internal/application/controllers"
	"timetable-homework-tgbot/internal/application/formatter"
	"timetable-homework-tgbot/internal/infrastracture/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HomeworkHandler struct {
	hw        controllers.HomeworkController
	bot       *telegram.Bot
	formatter *formatter.Formatter
	keyboard  *telegram.KeyboardController
}

func NewHomeworkHandler(hw controllers.HomeworkController, bot *telegram.Bot) *HomeworkHandler {
	return &HomeworkHandler{hw: hw, bot: bot, formatter: formatter.GetFormatter(), keyboard: telegram.GetKeyboardController()}
}

func (h *HomeworkHandler) PinStart(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	role, err := h.hw.GetRole(ctx, userID)
	if err != nil || role != "teacher" {
		_ = h.bot.Send(chatID, "Прикреплять домашние задания может только преподаватель.", h.keyboard.KBGuest())
		return
	}
	h.bot.GetState().Set(chatID, telegram.StateWaitHWGroup)
	_ = h.bot.SendRemove(chatID, "Для какой группы прикрепить ДЗ?")
}

func (h *HomeworkHandler) WaitGroup(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	group := strings.ReplaceAll(strings.TrimSpace(u.Message.Text), " ", "")
	days, err := h.hw.DaysWithLessons(ctx, group)
	if err != nil || len(days) == 0 {
		h.bot.GetState().Del(chatID)
		_ = h.bot.Send(chatID, "Не нашёл расписание для этой группы.", h.keyboard.KBTeacher())
		return
	}
	h.bot.HWSessSet(chatID, telegram.HwSession{Group: group})
	h.bot.GetState().Set(chatID, telegram.StateWaitHWDay)
	_ = h.bot.Send(chatID, "К какому дню недели прикрепить ДЗ?", h.keyboard.KBDays(days))
}

func (h *HomeworkHandler) WaitDay(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	day := strings.TrimSpace(u.Message.Text)
	s := h.bot.HWSessGet(chatID)
	lessons, err := h.hw.LessonsByDay(ctx, userID, s.Group, day)
	if err != nil || len(lessons) == 0 {
		h.bot.GetState().Del(chatID)
		_ = h.bot.Send(chatID, "В этот день нет твоих пар для выбранной группы.", h.keyboard.KBTeacher())
		return
	}
	s.Day = day
	h.bot.HWSessSet(chatID, s)

	h.bot.GetState().Set(chatID, telegram.StateWaitHWLesson)

	_ = h.bot.Send(chatID, "Выбери пару:", h.keyboard.KBLessons(lessons))
}

func (h *HomeworkHandler) WaitLesson(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	args := strings.Split(strings.TrimSpace(u.Message.Text), "-")
	lesson := strings.TrimSpace(args[0])

	s := h.bot.HWSessGet(chatID)
	s.LessonTitle = lesson
	h.bot.HWSessSet(chatID, s)
	curr := h.bot.GetState().Get(chatID)

	if strings.HasPrefix(curr, telegram.StateWaitHWEditLesson) {
		h.bot.GetState().Set(chatID, telegram.StateWaitHWTextEdit)
	} else {
		h.bot.GetState().Set(chatID, telegram.StateWaitHWText)
	}
	_ = h.bot.SendRemove(chatID, "Введи текст ДЗ:")
}

func (h *HomeworkHandler) WaitText(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	text := strings.TrimSpace(u.Message.Text)
	s := h.bot.HWSessGet(chatID)
	if s.Group == "" || s.Day == "" || s.LessonTitle == "" {
		_ = h.bot.Send(chatID, "Сессия потерялась, начни заново.", h.keyboard.KBTeacher())
		h.bot.GetState().Del(chatID)
		return
	}
	if err := h.hw.Pin(ctx, userID, s.Group, s.Day, s.LessonTitle, text); err != nil {
		_ = h.bot.Send(chatID, "Не удалось прикрепить домашнее задание", h.keyboard.KBTeacher())
		h.bot.GetState().Del(chatID)
		return
	}
	h.bot.GetState().Del(chatID)
	h.bot.HWSessDel(chatID)
	_ = h.bot.Send(chatID, "Домашнее задание сохранено ✅", h.keyboard.KBTeacher())
}

func (h *HomeworkHandler) EditStart(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	hw, err := h.hw.ListForLastWeek(ctx, userID)
	if err != nil || len(hw) == 0 {
		h.bot.GetState().Del(chatID)
		_ = h.bot.Send(chatID, "Нет домашек для редактирования.", h.keyboard.KBMember())
		return
	}
	h.bot.GetState().Set(chatID, telegram.StateWaitHWTable)
	_ = h.bot.Send(chatID, "Выбери день для редактирования ДЗ:", h.keyboard.KBHomeworks(hw))
}

func (h *HomeworkHandler) WaitTextEdit(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	newText := strings.TrimSpace(u.Message.Text)
	if h.bot.GetState().Get(chatID) != telegram.StateWaitHWTextEdit {
		_ = h.bot.Send(chatID, "Сессия потерялась, начни заново.", h.keyboard.KBMember())
		h.bot.GetState().Del(chatID)
		return
	}
	s := h.bot.HWSessGet(chatID)
	if s.LessonTitle == "" {
		_ = h.bot.Send(chatID, "Не удалось определить пару.", h.keyboard.KBMember())
		h.bot.GetState().Del(chatID)
		return
	}

	homeworkID, err := strconv.ParseInt(s.LessonTitle, 10, 64)
	if err != nil {
		_ = h.bot.Send(chatID, "Не удалось определить домашнее задание", h.keyboard.KBMember())
		h.bot.GetState().Del(chatID)
		return
	}

	if err := h.hw.Update(ctx, homeworkID, newText); err != nil {
		_ = h.bot.Send(chatID, "Не удалось обновить домашнее задание", h.keyboard.KBMember())
		h.bot.GetState().Del(chatID)
		return
	}

	h.bot.GetState().Del(chatID)
	h.bot.HWSessDel(chatID)
	_ = h.bot.Send(chatID, "Домашнее задание обновлено ✅", h.keyboard.KBMember())
}

func (h *HomeworkHandler) ListHomeworks(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	role, err := h.hw.GetRole(ctx, userID)
	if err != nil {
		log.Println(err)
		h.bot.GetState().Del(chatID)
		_ = h.bot.Send(chatID, "Не удалось определить роль.", h.keyboard.KBGuest())
		return
	}

	if u.Message.Text == telegram.BtnDeleteHomeworks {
		if role != "teacher" {
			h.bot.GetState().Del(chatID)
			_ = h.bot.Send(chatID, "Удалять домашние задания может только преподаватель.", h.keyboard.KBStudent())
			return
		}
		hw, err := h.hw.ListForTeacher(ctx, userID)
		if err != nil || len(hw) == 0 {
			log.Println(err)
			h.bot.GetState().Del(chatID)
			_ = h.bot.Send(chatID, "Нет домашек для удаления.\n ", h.keyboard.KBTeacher())
			return
		}
		h.bot.GetState().Set(chatID, telegram.StateWaitHWTableToDelete)
		_ = h.bot.Send(chatID, "Выбери домашку для удаления:\n ", h.keyboard.KBHomeworks(hw))
		return
	}

	hw, err := h.hw.ListForLastWeek(ctx, userID)
	if err != nil || len(hw) == 0 {
		h.bot.GetState().Del(chatID)
		_ = h.bot.Send(chatID, "Нет домашних заданий", h.keyboard.KBMember())
		return
	}
	switch u.Message.Text {
	case telegram.BtnWatchHomeworks:
		formHw := h.formatter.FormatHomeworks(hw)
		h.bot.GetState().Del(chatID)
		_ = h.bot.Send(chatID, "Список домашек:\n "+formHw, h.keyboard.KBMember())
	case telegram.BtnUpdateHomeworkStatus:
		h.bot.GetState().Set(chatID, telegram.StateWaitHomeworkUpdateChoose)
		_ = h.bot.Send(chatID, "Выбери домашку для обновления статуса:\n ", h.keyboard.KBHomeworks(hw))
	}
}

func (h *HomeworkHandler) WaitHomeWorkTable(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	userID := u.Message.From.ID
	s := h.bot.HWSessGet(chatID)
	homeworkID, err := parseHomeworkID(u.Message.Text)
	if err != nil {
		h.bot.GetState().Del(chatID)
		h.bot.HWSessDel(chatID)
		_ = h.bot.Send(chatID, "Некорректное домашнее задание", h.keyboard.KBMember())
		return
	}
	s.LessonTitle = strconv.FormatInt(homeworkID, 10)
	h.bot.HWSessSet(chatID, s)
	switch h.bot.GetState().Get(chatID) {
	case telegram.StateWaitHWTable:
		exist, err := h.hw.CheckExistence(ctx, userID, homeworkID)
		if err != nil {
			log.Println(err)
			h.bot.HWSessDel(chatID)
			h.bot.GetState().Del(chatID)
			return
		}
		if !exist {
			h.bot.GetState().Del(chatID)
			h.bot.HWSessDel(chatID)
			_ = h.bot.Send(chatID, "Некорректное домашнее задание", h.keyboard.KBMember())
			return
		}
		h.bot.GetState().Set(chatID, telegram.StateWaitHWTextEdit)
		_ = h.bot.SendRemove(chatID, "Введи новый текст ДЗ:")
	case telegram.StateWaitHWTableToDelete:
		role, err := h.hw.GetRole(ctx, userID)
		if err != nil || role != "teacher" {
			h.bot.GetState().Del(chatID)
			h.bot.HWSessDel(chatID)
			_ = h.bot.Send(chatID, "Удалять домашние задания может только преподаватель.", h.keyboard.KBStudent())
			return
		}
		h.bot.GetState().Set(chatID, telegram.StateWaitConfirmDelete)
		_ = h.bot.Send(chatID, "Подтвердите удаление:", h.keyboard.KBConfirmDelete())
	case telegram.StateWaitHomeworkUpdateChoose:
		exist, err := h.hw.CheckExistence(ctx, userID, homeworkID)
		if err != nil {
			log.Println(err)
			h.bot.HWSessDel(chatID)
			h.bot.GetState().Del(chatID)
			return
		}
		if !exist {
			h.bot.GetState().Del(chatID)
			h.bot.HWSessDel(chatID)
			_ = h.bot.Send(chatID, "Некорректное домашнее задание", h.keyboard.KBMember())
			return
		}
		if err := h.hw.UpdateStatus(ctx, homeworkID); err != nil {
			log.Println(err)
			h.bot.HWSessDel(chatID)
			h.bot.GetState().Del(chatID)
			return
		}
		h.bot.HWSessDel(chatID)
		h.bot.GetState().Del(chatID)
		_ = h.bot.Send(chatID, "Статус домашнего задания обновлен", h.keyboard.KBMember())

	}
}

func (h *HomeworkHandler) WaitConfirmDelete(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	s := h.bot.HWSessGet(chatID)
	if u.Message.Text == telegram.BtnDelete {
		role, err := h.hw.GetRole(ctx, userID)
		if err != nil || role != "teacher" {
			h.bot.GetState().Del(chatID)
			h.bot.HWSessDel(chatID)
			_ = h.bot.Send(chatID, "Удалять домашние задания может только преподаватель.", h.keyboard.KBStudent())
			return
		}
		homeworkID, err := strconv.ParseInt(s.LessonTitle, 10, 64)
		if err != nil {
			log.Println(err)
			h.bot.GetState().Del(chatID)
			return
		}
		if err := h.hw.DeleteHomework(ctx, userID, homeworkID); err != nil {
			log.Println(err)
			h.bot.GetState().Del(chatID)
			return
		}
		h.bot.GetState().Del(chatID)
		h.bot.HWSessDel(chatID)
		_ = h.bot.Send(chatID, "Домашнее задание удалено ✅", h.keyboard.KBTeacher())
	} else if u.Message.Text == telegram.BtnNotDelete {
		h.bot.GetState().Del(chatID)
		h.bot.HWSessDel(chatID)
		_ = h.bot.Send(chatID, "Домашнее задание оставлено ✅", h.keyboard.KBTeacher())
	}

}

func parseHomeworkID(label string) (int64, error) {
	label = strings.TrimSpace(label)
	label = strings.TrimPrefix(label, "#")
	idPart := strings.TrimSpace(strings.SplitN(label, ":", 2)[0])
	return strconv.ParseInt(idPart, 10, 64)
}
