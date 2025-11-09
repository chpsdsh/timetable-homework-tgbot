package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) hStart(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	kb := kbAskJoin()
	msg := tgbotapi.NewMessage(chatID, "Хотите присоединиться к своей группе в НГУ?")
	msg.ReplyMarkup = kb
	_ = b.sendWithRetry(msg)
}

func (b *Bot) hJoin(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	// TODO(DB): проверить не привязан ли уже пользователь к группе, если привязан,
	// то показываем msg.ReplyMarkup = kbMember() //клавиатуру авторизованного пользователя
	b.state.Set(chatID, stateWait)
	msg := tgbotapi.NewMessage(chatID, "Введи номер своей группы (например, 23204).")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_ = b.sendWithRetry(msg)
}

func (b *Bot) hSkip(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "Ок, продолжим без привязки.")
	msg.ReplyMarkup = kbGuest() // клавиатура пользователя
	_ = b.sendWithRetry(msg)
}

func (b *Bot) hWaitGroup(ctx context.Context, u tgbotapi.Update) {
	m := u.Message
	chatID, userID := m.Chat.ID, m.From.ID
	fmt.Println(userID)
	group := strings.TrimSpace(m.Text)

	//TODO(DB) : валидация существования группы через бд
	//if ok := bd.validateGroup(group); ok != nil {
	//	msg := tgbotapi.NewMessage(chatID, "Такой группы в НГУ нет, попробуй снова.")
	//  b.sendWithRetry(msg)
	//  continue
	//}

	//// TODO(DB): сохранить привязку: bd.SetGroup(ctx, userID, group)

	log.Printf("UserID : %d in group: %s", userID, group)
	b.state.Del(chatID)
	msg := tgbotapi.NewMessage(chatID, "Группа сохранена: "+group)
	msg.ReplyMarkup = kbMember() //клавиатура авторизованного пользователя
	_ = b.sendWithRetry(msg)
}

func (b *Bot) hDefault(ctx context.Context, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	_ = b.sendWithRetry(tgbotapi.NewMessage(chatID, "Нажми кнопку меню или /start"))
}
