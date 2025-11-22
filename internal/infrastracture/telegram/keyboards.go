package telegram

import (
	"fmt"

	"timetable-homework-tgbot/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BtnShowTimeTable = "Посмотреть расписание"
	BtnPinHW         = "Прикрепить домашнее задание"
	BtnChangeHW      = "Редактировать домашнее задание"
	BtnConfReminder  = "Настроить напоминание о домашнем задание"
	BtnLeave         = "Отсоедениться от группы"
	BtnGroup         = "Группы"
	BtnTeacher       = "Преподаватели"
	BtnClassRoom     = "Аудитории"
	BtnJoin          = "Присоединиться к группе"
	BtnSkip          = "Не присоединяться к группе"
)

func KBAskJoin() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnJoin),
			tgbotapi.NewKeyboardButton(BtnSkip),
		),
	)
	kb.ResizeKeyboard = true
	return kb
}

func KBMember() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnShowTimeTable),
			tgbotapi.NewKeyboardButton(BtnPinHW),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnChangeHW),
			tgbotapi.NewKeyboardButton(BtnConfReminder),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnLeave),
		),
	)
	kb.ResizeKeyboard = true
	return kb
}

func KBGuest() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnShowTimeTable),
			tgbotapi.NewKeyboardButton(BtnJoin),
		),
	)
	kb.ResizeKeyboard = true
	return kb
}

func KBChooseTimetable() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnGroup),
			tgbotapi.NewKeyboardButton(BtnTeacher),
			tgbotapi.NewKeyboardButton(BtnClassRoom),
		),
	)
	kb.ResizeKeyboard = true
	return kb
}

func KBDays(days []string) tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, 0, len(days))
	for _, d := range days {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(d)))
	}
	kb := tgbotapi.NewReplyKeyboard(rows...)
	kb.ResizeKeyboard = true
	return kb
}

func KBLessons(list []domain.LessonBrief) tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, 0, len(list))
	for _, l := range list {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(l.Title)))
	}
	kb := tgbotapi.NewReplyKeyboard(rows...)
	kb.ResizeKeyboard = true
	return kb
}

func KBHomeworks(list []domain.HWBrief) tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, 0, len(list))
	for _, h := range list {
		label := fmt.Sprintf("%s : %s", h.Subject, h.HomeworkText)
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(label)))
	}
	kb := tgbotapi.NewReplyKeyboard(rows...)
	kb.ResizeKeyboard = true
	return kb
}

func KBWeekdays() tgbotapi.ReplyKeyboardMarkup {
	d := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(d[0]), tgbotapi.NewKeyboardButton(d[1]), tgbotapi.NewKeyboardButton(d[2]), tgbotapi.NewKeyboardButton(d[3])),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(d[4]), tgbotapi.NewKeyboardButton(d[5]), tgbotapi.NewKeyboardButton(d[6])),
	)
	kb.ResizeKeyboard = true
	return kb
}

func KBTimeSlots() tgbotapi.ReplyKeyboardMarkup {
	s := []string{"08:00", "10:00", "12:00", "14:00", "16:00", "18:00", "20:00", "22:00"}
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(s[0]), tgbotapi.NewKeyboardButton(s[1]), tgbotapi.NewKeyboardButton(s[2]), tgbotapi.NewKeyboardButton(s[3])),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(s[4]), tgbotapi.NewKeyboardButton(s[5]), tgbotapi.NewKeyboardButton(s[6]), tgbotapi.NewKeyboardButton(s[7])),
	)
	kb.ResizeKeyboard = true
	return kb
}
