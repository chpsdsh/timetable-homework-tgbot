package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"timetable-homework-tgbot/internal/app"
	"timetable-homework-tgbot/internal/repositories"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("key.env")

	// Инициализируй реальные репозитории тут (db := sql.Open(...)):
	var usersRepo repositories.UsersRepository        // = repositories.NewUsersRepo(db)
	var lessonsRepo repositories.LessonsRepository    // = repositories.NewLessonsRepo(db/парсер)
	var hwRepo repositories.HomeworkRepository        // = repositories.NewHWRepo(db)
	var notifRepo repositories.NotificationRepository // = repositories.NewNotifRepo(db)

	// Собираем приложение из ENV (создаст tgbotapi.BotAPI сам)
	a, err := app.NewFromEnv(usersRepo, lessonsRepo, hwRepo, notifRepo)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
