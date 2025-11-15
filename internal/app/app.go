package app

import (
	"context"
	"fmt"
	"os"

	"timetable-homework-tgbot/internal/infrastracture/controllers"
	"timetable-homework-tgbot/internal/infrastracture/handlers"
	"timetable-homework-tgbot/internal/infrastracture/telegram"
	"timetable-homework-tgbot/internal/repositories" // оставляем импорт для сигнатуры

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct{ bot *telegram.Bot }

// NewWithDeps — сигнатуру не меняем, но deps сейчас игнорим (всё фейково).
func NewWithDeps(
	api *tgbotapi.BotAPI,
	_ repositories.UsersRepository,
	_ repositories.LessonsRepository,
	_ repositories.HomeworkRepository,
	_ repositories.NotificationRepository,
) (*App, error) {

	// Контроллеры — ФЕЙКИ с TODO(DB)
	authCtl := controllers.NewAuthFake("Europe/Bucharest")
	hwCtl := controllers.NewHomeworkFake()
	notifCtl := controllers.NewNotificationFake(authCtl)

	// Telegram
	state := telegram.NewMemState()
	bot := telegram.NewBot(api, state)

	// Хендлеры
	cmdH := handlers.NewCommandHandler(authCtl, bot)
	ttH := handlers.NewTimetableHandler(bot)
	hwH := handlers.NewHWHandler(hwCtl, bot)
	ntH := handlers.NewNotifyHandler(hwCtl, notifCtl, bot)

	// Роутинг
	r := bot.Router()
	// команды
	r.OnCommand("start", cmdH.Start)

	// меню
	r.OnText(telegram.BtnShowTimeTable, ttH.ShowMenu)
	r.OnText(telegram.BtnJoin, cmdH.Join)
	r.OnText(telegram.BtnLeave, cmdH.Leave)
	r.OnText(telegram.BtnSkip, cmdH.Skip)
	r.OnState(telegram.StateWaitUserGroup, cmdH.WaitUserGroup)

	// расписание
	r.OnText(telegram.BtnGroup, ttH.AskGroup)
	r.OnText(telegram.BtnTeacher, ttH.AskTeacher)
	r.OnText(telegram.BtnClassRoom, ttH.AskRoom)
	r.OnState(telegram.StateWaitGroupTB, ttH.WaitGroup)
	r.OnState(telegram.StateWaitTeacherTB, ttH.WaitTeacher)
	r.OnState(telegram.StateWaitRoomTB, ttH.WaitRoom)

	// ДЗ
	r.OnText(telegram.BtnPinHW, hwH.PinStart)
	r.OnState(telegram.StateWaitHWDay, hwH.WaitDay)
	r.OnState(telegram.StateWaitHWLesson, hwH.WaitLesson)
	r.OnState(telegram.StateWaitHWText, hwH.WaitText)

	r.OnText(telegram.BtnChangeHW, hwH.EditStart)
	r.OnState(telegram.StateWaitHWEditDay, hwH.WaitDay)
	r.OnState(telegram.StateWaitHWEditLesson, hwH.WaitLesson)
	r.OnState(telegram.StateWaitHWTextEdit, hwH.WaitTextEdit)

	// Напоминания
	r.OnText(telegram.BtnConfReminder, ntH.Start)
	r.OnState(telegram.StateWaitRemindChooseHW, ntH.WaitChooseHW)
	r.OnState(telegram.StateWaitRemindChooseDay, ntH.WaitChooseDay)
	r.OnState(telegram.StateWaitRemindChooseTime, ntH.WaitChooseTime)

	// дефолт
	r.Default(func(ctx context.Context, u tgbotapi.Update) {
		_ = bot.Send(u.Message.Chat.ID, "Нажми кнопку меню или /start", telegram.KBMember())
	})

	return &App{bot: bot}, nil
}

func NewFromEnv(
	users repositories.UsersRepository,
	lessons repositories.LessonsRepository,
	hw repositories.HomeworkRepository,
	notifs repositories.NotificationRepository,
) (*App, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is empty")
	}
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return NewWithDeps(api, users, lessons, hw, notifs)
}

func (a *App) Run(ctx context.Context) error { return a.bot.Run(ctx) }
