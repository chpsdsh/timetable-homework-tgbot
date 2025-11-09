package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func kbAskJoin() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(btnJoin),
			tgbotapi.NewKeyboardButton(btnSkip),
		),
	)
	kb.ResizeKeyboard = true

	return kb
}

func kbMember() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Посмотреть расписание"),
			tgbotapi.NewKeyboardButton("Прикрепить домашнее задание"),
			tgbotapi.NewKeyboardButton("Редактировать домашнее задание"),
			tgbotapi.NewKeyboardButton("Настроить напоминание о домашнем задание"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отсоедениться от группы"),
		),
	)
	kb.ResizeKeyboard = true
	return kb
}

func kbGuest() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Посмотреть расписание"),
			tgbotapi.NewKeyboardButton("Присоедениться к группе"),
		),
	)
	kb.ResizeKeyboard = true
	return kb
}
