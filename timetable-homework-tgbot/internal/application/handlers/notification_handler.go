package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"timetable-homework-tgbot/internal/application/controllers"
	"timetable-homework-tgbot/internal/infrastracture/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NotifyHandler struct {
	hw       controllers.HomeworkController
	ctl      controllers.NotificationController
	bot      *telegram.Bot
	keyboard *telegram.KeyboardController
}

func NewNotifyHandler(hw controllers.HomeworkController, ctl controllers.NotificationController, bot *telegram.Bot) *NotifyHandler {
	return &NotifyHandler{hw: hw, ctl: ctl, bot: bot, keyboard: telegram.GetKeyboardController()}
}

func (h *NotifyHandler) Start(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	list, err := h.hw.ListForLastWeek(ctx, userID)
	if err != nil || len(list) == 0 {
		_ = h.bot.Send(chatID, "За последнюю неделю ДЗ не найдено.", h.keyboard.KBMember())
		return
	}
	h.bot.GetState().Set(chatID, telegram.StateWaitRemindChooseHW)
	_ = h.bot.Send(chatID, "Выбери ДЗ, для которого поставить напоминание:", h.keyboard.KBHomeworks(list))
}

func (h *NotifyHandler) WaitChooseHW(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	userID := u.Message.From.ID
	homework := strings.TrimSpace(u.Message.Text)
	arg := strings.Split(homework, ":")

	exist, err := h.hw.CheckExistence(ctx, userID, strings.TrimSpace(arg[0]))
	if err != nil {
		log.Println(err)
		h.bot.HWSessDel(chatID)
		h.bot.GetState().Del(chatID)
		return
	}
	if !exist {
		log.Println(homework)
		h.bot.GetState().Del(chatID)
		h.bot.HWSessDel(chatID)
		_ = h.bot.Send(chatID, "Некорректное домашнее задание", h.keyboard.KBMember())
		return
	}

	h.bot.RemSessSet(chatID, telegram.RemindSession{SubjectWithTask: homework})
	h.bot.GetState().Set(chatID, telegram.StateWaitRemindChooseDay)
	_ = h.bot.Send(chatID, "В какой день напоминать?", h.keyboard.KBWeekdays(time.Now()))
}

func (h *NotifyHandler) WaitChooseDay(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	date := strings.TrimSpace(u.Message.Text)
	if !isDate(date) {
		h.bot.GetState().Del(chatID)
		_ = h.bot.SendRemove(chatID, "Неверный формат даты")
		return
	}

	s, _ := h.bot.RemSessGet(chatID)
	s.Date = date
	h.bot.RemSessSet(chatID, s)
	h.bot.GetState().Set(chatID, telegram.StateWaitRemindChooseTime)
	_ = h.bot.SendRemove(chatID, "Во сколько напоминать(HH:MM)?")
}

func (h *NotifyHandler) WaitChooseTime(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	tStr := strings.TrimSpace(u.Message.Text)
	if !isHHMM(tStr) {
		_ = h.bot.SendRemove(chatID, "Формат времени HH:MM.")
		return
	}

	s, ok := h.bot.RemSessGet(chatID)
	if !ok || s.SubjectWithTask == "" {
		_ = h.bot.Send(chatID, "Сессия потерялась. Начни заново.", h.keyboard.KBMember())
		h.bot.GetState().Del(chatID)
		return
	}
	s.TimeHHMM = tStr

	log.Println("subject:", s.SubjectWithTask)
	if err := h.ctl.SetReminder(ctx, userID, s.SubjectWithTask, s.Date, s.TimeHHMM); err != nil {
		log.Println(err.Error())
		_ = h.bot.Send(chatID, "Не удалось создать напоминание.", h.keyboard.KBMember())
		h.bot.RemSessDel(chatID)
		h.bot.GetState().Del(chatID)
		return
	}

	h.bot.RemSessDel(chatID)
	h.bot.GetState().Del(chatID)
	_ = h.bot.Send(chatID, fmt.Sprintf("Напоминание поставлено: %s в %s ✅", s.Date, s.TimeHHMM), h.keyboard.KBMember())
}

func (h *NotifyHandler) StartDeleteNotification(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID

	list, err := h.ctl.GetUserNotifications(ctx, userID)

	if err != nil || len(list) == 0 {
		h.bot.GetState().Del(chatID)
		_ = h.bot.Send(chatID, "Напоминаний не найдено.", h.keyboard.KBMember())
		return
	}
	h.bot.GetState().Set(chatID, telegram.StateWaitRemindChoose)
	_ = h.bot.Send(chatID, "Выбери напоминание, для удаления:", h.keyboard.KBNotifications(list))

}

func (h *NotifyHandler) WaitDeleteNotification(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	not := strings.TrimSpace(u.Message.Text)
	log.Println(not)
	if err := h.ctl.DeleteUserNotification(ctx, userID, not); err != nil {
		_ = h.bot.Send(chatID, "Не удалось удалить напоминание.", h.keyboard.KBMember())
		h.bot.GetState().Del(chatID)
		return
	}
	h.bot.GetState().Del(chatID)
	_ = h.bot.Send(chatID, fmt.Sprintf("Напоминание удалено: %s ✅", not), h.keyboard.KBMember())
}

func (h *NotifyHandler) StartNotificationWorker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		log.Println("ticker Started")
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h.checkPendingNotifications(ctx)
			}
		}
	}()
}

func (h *NotifyHandler) checkPendingNotifications(ctx context.Context) {
	pending, err := h.ctl.GetPendingNotifications(ctx)
	if err != nil {
		log.Println("GetPendingNotifications:", err)
		return
	}

	for _, n := range pending {

		text := fmt.Sprintf(
			"Напоминание по предмету %s на %s",
			n.Subject,
			n.Timestamp.Format("02.01.2006 15:04"),
		)
		log.Println(text)
		if err := h.bot.Send(n.UserID, text, h.keyboard.KBMember()); err != nil {
			log.Println("send notif:", err)
			continue
		}

		if err := h.ctl.DeleteUserNotificationWithTs(ctx, n.UserID, n.Subject, n.Timestamp); err != nil {
			log.Println("delete notif:", err)
		}
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

func isDate(dateStr string) bool {
	const layout = "02.01.2006"

	_, err := time.Parse(layout, dateStr)
	if err != nil {
		return false
	}

	return true
}
