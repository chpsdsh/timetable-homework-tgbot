package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) handleStart(chatID int64) error {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Присоединиться к группе"),
			tgbotapi.NewKeyboardButton("Не присоединяться к группе"),
		),
	)

	kb.ResizeKeyboard = true
	kb.OneTimeKeyboard = false

	msg := tgbotapi.NewMessage(chatID, "Привет. Ты зашел в телеграм бота с расписанием НГУ. Хочешь присоединиться к своей группе в НГУ?")
	msg.ReplyMarkup = kb
	return b.sendWithRetry(msg)
}
