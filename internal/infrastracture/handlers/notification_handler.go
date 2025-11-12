package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"timetable-homework-tgbot/internal/domain"
	"timetable-homework-tgbot/internal/infrastracture/controllers"
	"timetable-homework-tgbot/internal/infrastracture/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NotifyHandler struct {
	hw  controllers.HomeworkController
	ctl controllers.NotificationController
	bot *telegram.Bot
}

func NewNotifyHandler(hw controllers.HomeworkController, ctl controllers.NotificationController, bot *telegram.Bot) *NotifyHandler {
	return &NotifyHandler{hw: hw, ctl: ctl, bot: bot}
}

func (h *NotifyHandler) Start(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	list, err := h.hw.ListForLastWeek(ctx, userID)
	if err != nil || len(list) == 0 {
		_ = h.bot.Send(chatID, "За последнюю неделю ДЗ не найдено.", telegram.KBMember())
		return
	}
	h.bot.State.Set(chatID, telegram.StateWaitRemindChooseHW)
	_ = h.bot.Send(chatID, "Выбери ДЗ, для которого поставить напоминание:", telegram.KBHomeworks(toDomain(list)))
}

func (h *NotifyHandler) WaitChooseHW(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	id, ok := telegram.ExtractIDFromLabel(strings.TrimSpace(u.Message.Text))
	if !ok {
		_ = h.bot.Send(chatID, "Не понял, какую запись ДЗ выбрали.", telegram.KBMember())
		return
	}
	h.bot.RemSessSet(chatID, telegram.RemindSession{HomeworkID: id})
	h.bot.State.Set(chatID, telegram.StateWaitRemindChooseDay)
	_ = h.bot.Send(chatID, "В какой день недели напоминать?", telegram.KBWeekdays())
}

func (h *NotifyHandler) WaitChooseDay(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	wd, ok := parseWeekday(strings.TrimSpace(u.Message.Text))
	if !ok {
		_ = h.bot.Send(chatID, "Выбери день из клавиатуры (Пн..Вс).", telegram.KBWeekdays())
		return
	}
	s, _ := h.bot.RemSessGet(chatID)
	s.Weekday = wd
	h.bot.RemSessSet(chatID, s)
	h.bot.State.Set(chatID, telegram.StateWaitRemindChooseTime)
	_ = h.bot.Send(chatID, "Во сколько напоминать?", telegram.KBTimeSlots())
}

func (h *NotifyHandler) WaitChooseTime(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	tStr := strings.TrimSpace(u.Message.Text)
	if !isHHMM(tStr) {
		_ = h.bot.Send(chatID, "Формат времени HH:MM.", telegram.KBTimeSlots())
		return
	}

	s, ok := h.bot.RemSessGet(chatID)
	if !ok || s.HomeworkID == "" {
		_ = h.bot.Send(chatID, "Сессия потерялась. Начни заново.", telegram.KBMember())
		h.bot.State.Del(chatID)
		return
	}
	s.TimeHHMM = tStr

	if err := h.ctl.SetWeeklyReminder(ctx, userID, s.HomeworkID, s.Weekday, s.TimeHHMM); err != nil {
		_ = h.bot.Send(chatID, "Не удалось создать напоминание.", telegram.KBMember())
		h.bot.RemSessDel(chatID)
		h.bot.State.Del(chatID)
		return
	}

	h.bot.RemSessDel(chatID)
	h.bot.State.Del(chatID)
	_ = h.bot.Send(chatID, fmt.Sprintf("Напоминание поставлено: %s в %s ✅", ruWeekdayShort(s.Weekday), s.TimeHHMM), telegram.KBMember())
}

func toDomain(in []domain.HWBrief) []domain.HWBrief { return in } // заглушка, если тип совпадает

func parseWeekday(s string) (time.Weekday, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "пн":
		return time.Monday, true
	case "вт":
		return time.Tuesday, true
	case "ср":
		return time.Wednesday, true
	case "чт":
		return time.Thursday, true
	case "пт":
		return time.Friday, true
	case "сб":
		return time.Saturday, true
	case "вс":
		return time.Sunday, true
	default:
		return time.Sunday, false
	}
}

func isHHMM(s string) bool {
	if len(s) != 5 || s[2] != ':' {
		return false
	}
	hh, err1 := strconv.Atoi(s[:2])
	mm, err2 := strconv.Atoi(s[3:])
	return err1 == nil && err2 == nil && hh >= 0 && hh < 24 && mm >= 0 && mm < 60
}

func ruWeekdayShort(wd time.Weekday) string {
	switch wd {
	case time.Monday:
		return "Пн"
	case time.Tuesday:
		return "Вт"
	case time.Wednesday:
		return "Ср"
	case time.Thursday:
		return "Чт"
	case time.Friday:
		return "Пт"
	case time.Saturday:
		return "Сб"
	default:
		return "Вс"
	}
}
