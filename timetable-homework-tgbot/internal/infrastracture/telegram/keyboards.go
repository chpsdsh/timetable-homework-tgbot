package telegram

import (
	"fmt"
	"time"

	"timetable-homework-tgbot/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BtnShowTimeTable              = "Посмотреть расписание"
	BtnPinHW                      = "Прикрепить домашнее задание"
	BtnChangeHW                   = "Редактировать домашнее задание"
	BtnConfReminder               = "Настроить напоминание о домашнем задание"
	BtnDeleteNotification         = "Удалить напоминание о домашних заданиях"
	BtnWatchHomeworks             = "Посмотреть домашние задания"
	BtnDeleteHomeworks            = "Удалить домашнее задание"
	BtnUpdateHomeworkStatus       = "Отметить домашку сделанной"
	StateWaitHomeworkUpdateChoose = "Выбрать домашку чтобы пометить сделанной"
	BtnLeave                      = "Отсоедениться от группы"
	BtnGroup                      = "Группы"
	BtnTeacher                    = "Преподаватели"
	BtnClassRoom                  = "Аудитории"
	BtnJoin                       = "Присоединиться к группе"
	BtnSkip                       = "Не присоединяться к группе"
	BtnDelete                     = "Удалить"
	BtnNotDelete                  = "Не удалять"
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
			tgbotapi.NewKeyboardButton(BtnWatchHomeworks),
			tgbotapi.NewKeyboardButton(BtnChangeHW),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnConfReminder),
			tgbotapi.NewKeyboardButton(BtnDeleteNotification),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnDeleteHomeworks),
			tgbotapi.NewKeyboardButton(BtnUpdateHomeworkStatus),
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
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(fmt.Sprintf("%s - %s", l.Title, l.LessonType))))
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

func KBNotifications(list []domain.Notification) tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, 0, len(list))

	for _, n := range list {
		label := fmt.Sprintf(
			"%s — %s",
			n.Subject,
			n.Timestamp.Format("02.01.2006 15:04"),
		)

		rows = append(rows,
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(label),
			),
		)
	}

	kb := tgbotapi.NewReplyKeyboard(rows...)
	kb.ResizeKeyboard = true
	return kb
}

func KBWeekdays(today time.Time) tgbotapi.ReplyKeyboardMarkup {
	const days = 8

	rows := make([][]tgbotapi.KeyboardButton, 0)

	for i := 0; i < days; i++ {
		d := today.AddDate(0, 0, i)
		text := d.Format("02.01.2006")

		btn := tgbotapi.NewKeyboardButton(text)

		if i%4 == 0 {
			rows = append(rows, tgbotapi.NewKeyboardButtonRow(btn))
		} else {
			rows[len(rows)-1] = append(rows[len(rows)-1], btn)
		}
	}

	kb := tgbotapi.NewReplyKeyboard(rows...)
	kb.ResizeKeyboard = true
	return kb
}

func KBConfirmDelete() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnDelete),
			tgbotapi.NewKeyboardButton(BtnNotDelete),
		),
	)
	kb.ResizeKeyboard = true
	return kb
}
