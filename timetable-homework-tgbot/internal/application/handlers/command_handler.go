package handlers

import (
	"context"
	"log"
	"strings"
	"timetable-homework-tgbot/internal/application/controllers"
	"timetable-homework-tgbot/internal/infrastracture/repositories"
	"timetable-homework-tgbot/internal/infrastracture/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	auth     controllers.AuthController
	bot      *telegram.Bot
	keyboard *telegram.KeyboardController
}

func NewCommandHandler(auth controllers.AuthController, bot *telegram.Bot) *CommandHandler {
	return &CommandHandler{auth: auth, bot: bot, keyboard: telegram.GetKeyboardController()}
}

func (h *CommandHandler) Start(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "Выбери роль для входа в бот.")
	msg.ReplyMarkup = h.keyboard.KBAskJoin()
	_ = h.bot.SendWithRetry(msg)
}

func (h *CommandHandler) Join(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	registered, _ := h.auth.EnsureRegistered(ctx, userID)
	if registered {
		role, _ := h.auth.GetRole(ctx, userID)
		_ = h.bot.Send(chatID, "Ты уже вошёл в бот.", keyboardByRole(h.keyboard, role))
		return
	}
	h.bot.AuthSessSet(chatID, telegram.AuthSession{Role: repositories.RoleStudent})
	h.bot.GetState().Set(chatID, telegram.StateWaitLogin)
	_ = h.bot.SendRemove(chatID, "Введи логин студента.")
}

func (h *CommandHandler) JoinTeacher(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	registered, _ := h.auth.EnsureRegistered(ctx, userID)
	if registered {
		role, _ := h.auth.GetRole(ctx, userID)
		_ = h.bot.Send(chatID, "Ты уже вошёл в бот.", keyboardByRole(h.keyboard, role))
		return
	}
	h.bot.AuthSessSet(chatID, telegram.AuthSession{Role: repositories.RoleTeacher})
	h.bot.GetState().Set(chatID, telegram.StateWaitLogin)
	_ = h.bot.SendRemove(chatID, "Введи логин преподавателя.")
}

func (h *CommandHandler) Leave(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	h.bot.GetState().Del(chatID)
	h.bot.AuthSessDel(chatID)
	if err := h.auth.LeaveGroup(ctx, userID); err != nil {
		log.Printf("leave failed: %v", err)
	}
	_ = h.bot.Send(chatID, "Вы вышли из аккаунта.", h.keyboard.KBGuest())
}

func (h *CommandHandler) WaitUserGroup(ctx context.Context, u tgbotapi.Update) {
	m := u.Message
	chatID, userID := m.Chat.ID, m.From.ID
	group := strings.TrimSpace(m.Text)

	if group == "" {
		_ = h.bot.Send(chatID, "Пусто. Введи номер своей группы (например, 23204).", h.keyboard.KBMember())
		return
	}

	if err := h.auth.RegisterStudent(ctx, userID, group); err != nil {
		log.Printf("RegisterStudent failed: %v", err)
		h.bot.GetState().Del(chatID)
		_ = h.bot.Send(chatID, "Не удалось войти. Проверь номер группы и попробуй ещё раз.", h.keyboard.KBGuest())
		return
	}

	h.bot.GetState().Del(chatID)
	_ = h.bot.Send(chatID, "Группа сохранена: "+group, h.keyboard.KBStudent())
}

func (h *CommandHandler) WaitLogin(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	login := strings.TrimSpace(u.Message.Text)
	if login == "" {
		_ = h.bot.SendRemove(chatID, "Логин пустой. Введи логин.")
		return
	}

	s, ok := h.bot.AuthSessGet(chatID)
	if !ok || s.Role == "" {
		h.bot.GetState().Del(chatID)
		_ = h.bot.Send(chatID, "Сессия входа потерялась. Начни заново через /start.", h.keyboard.KBGuest())
		return
	}
	s.Login = login
	h.bot.AuthSessSet(chatID, s)
	h.bot.GetState().Set(chatID, telegram.StateWaitPassword)
	_ = h.bot.SendRemove(chatID, "Введи пароль.")
}

func (h *CommandHandler) WaitPassword(ctx context.Context, u tgbotapi.Update) {
	chatID, userID := u.Message.Chat.ID, u.Message.From.ID
	password := strings.TrimSpace(u.Message.Text)
	s, ok := h.bot.AuthSessGet(chatID)
	if !ok || s.Role == "" || s.Login == "" {
		h.bot.GetState().Del(chatID)
		h.bot.AuthSessDel(chatID)
		_ = h.bot.Send(chatID, "Сессия входа потерялась. Начни заново через /start.", h.keyboard.KBGuest())
		return
	}

	if err := h.auth.Login(ctx, userID, s.Role, s.Login, password); err != nil {
		log.Printf("login failed: role=%s login=%s err=%v", s.Role, s.Login, err)
		h.bot.GetState().Del(chatID)
		h.bot.AuthSessDel(chatID)
		_ = h.bot.Send(chatID, "Неверный логин или пароль.", h.keyboard.KBGuest())
		return
	}

	h.bot.GetState().Del(chatID)
	h.bot.AuthSessDel(chatID)
	if s.Role == repositories.RoleTeacher {
		_ = h.bot.Send(chatID, "Готово, ты вошёл как преподаватель.", h.keyboard.KBTeacher())
		return
	}
	_ = h.bot.Send(chatID, "Готово, ты вошёл как студент.", h.keyboard.KBStudent())
}

func (h *CommandHandler) Skip(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	h.bot.GetState().Del(chatID)
	h.bot.AuthSessDel(chatID)
	_ = h.bot.Send(chatID, "Ок, продолжим без привязки.", h.keyboard.KBGuest())
}

func keyboardByRole(k *telegram.KeyboardController, role string) tgbotapi.ReplyKeyboardMarkup {
	switch role {
	case "teacher":
		return k.KBTeacher()
	case "student":
		return k.KBStudent()
	default:
		return k.KBGuest()
	}
}
