package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"timetable-homework-tgbot/internal/app"
	"timetable-homework-tgbot/internal/infrastracture/database"
	"timetable-homework-tgbot/internal/repositories"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("key.env")

	ctx := context.Background()
	db, err := database.NewDB(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	if err := db.InitSchema(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	if err := db.FillDatabase(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

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
