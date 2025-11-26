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

type TimetableHandler struct {
	bot               *telegram.Bot
	lessonsController controllers.LessonsController
}

func NewTimetableHandler(bot *telegram.Bot, lessonCtl controllers.LessonsController) *TimetableHandler {
	return &TimetableHandler{bot: bot, lessonsController: lessonCtl}
}

func (h *TimetableHandler) ShowMenu(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	_ = h.bot.Send(chatID, "Какое расписание ты хочешь увидеть?", telegram.KBChooseTimetable())
}

func (h *TimetableHandler) AskGroup(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	h.bot.GetState().Set(chatID, telegram.StateWaitGroupTB)
	_ = h.bot.SendRemove(chatID, "Введи номер группы")
}

func (h *TimetableHandler) WaitGroup(ctx context.Context, u tgbotapi.Update) {
	m := u.Message
	chatID := m.Chat.ID
	userId := u.Message.From.ID
	group := strings.TrimSpace(m.Text)

	timetable := h.lessonsController.GetTimetableGroup(ctx, group)

	h.bot.GetState().Del(chatID)

	joined, err := h.lessonsController.EnsureJoined(ctx, userId)
	if err != nil {
		log.Println(err)
		h.bot.GetState().Del(chatID)
		return
	}
	var keyboardFunc func() tgbotapi.ReplyKeyboardMarkup
	if joined {
		keyboardFunc = telegram.KBMember
	} else {
		keyboardFunc = telegram.KBGuest
	}

	parts := formatter.SplitForTelegram(timetable)

	for _, part := range parts {
		if err := h.bot.Send(chatID, part, keyboardFunc()); err != nil {
			log.Println("send timetable part:", err)
			break
		}
	}
}

func (h *TimetableHandler) AskTeacher(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	h.bot.GetState().Set(chatID, telegram.StateWaitTeacherTB)
	_ = h.bot.SendRemove(chatID, "Введи ФИО преподавателя")
}

func (h *TimetableHandler) WaitTeacher(ctx context.Context, u tgbotapi.Update) {
	m := u.Message
	chatID := m.Chat.ID
	userId := u.Message.From.ID
	teacher := strings.TrimSpace(m.Text)

	timetable := h.lessonsController.GetTimetableTeacher(ctx, teacher)
	h.bot.GetState().Del(chatID)

	joined, err := h.lessonsController.EnsureJoined(ctx, userId)
	if err != nil {
		log.Println(err)
		h.bot.GetState().Del(chatID)
		return
	}
	var keyboardFunc func() tgbotapi.ReplyKeyboardMarkup
	if joined {
		keyboardFunc = telegram.KBMember
	} else {
		keyboardFunc = telegram.KBGuest
	}

	parts := formatter.SplitForTelegram(timetable)

	for _, part := range parts {
		if err := h.bot.Send(chatID, part, keyboardFunc()); err != nil {
			log.Println("send timetable part:", err)
			break
		}
	}
}

func (h *TimetableHandler) AskRoom(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	h.bot.GetState().Set(chatID, telegram.StateWaitRoomTB)
	_ = h.bot.SendRemove(chatID, "Введи номер аудитории")
}

func (h *TimetableHandler) WaitRoom(ctx context.Context, u tgbotapi.Update) {
	m := u.Message
	chatID := m.Chat.ID
	userId := u.Message.From.ID
	room := strings.TrimSpace(m.Text)

	timetable := h.lessonsController.GetTimetableRoom(ctx, room)
	h.bot.GetState().Del(chatID)
	joined, err := h.lessonsController.EnsureJoined(ctx, userId)
	if err != nil {
		log.Println(err)
		h.bot.GetState().Del(chatID)
		return
	}
	var keyboardFunc func() tgbotapi.ReplyKeyboardMarkup
	if joined {
		keyboardFunc = telegram.KBMember
	} else {
		keyboardFunc = telegram.KBGuest
	}
	parts := formatter.SplitForTelegram(timetable)
	for _, part := range parts {
		if err := h.bot.Send(chatID, part, keyboardFunc()); err != nil {
			log.Println("send timetable part:", err)
			break
		}
	}
}
