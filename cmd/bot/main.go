package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"timetable-homework-tgbot/internal/application"
	"timetable-homework-tgbot/internal/infrastracture/database"
	"timetable-homework-tgbot/internal/infrastracture/repositories"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("key.env")

	ctx := context.Background()
	db, err := database.NewDB(ctx)
	defer func() { _ = db.Close() }()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	if err := db.InitSchema(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	if err := db.FillDatabase(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	var usersRepo repositories.UsersRepository = repositories.NewUserRepo(db)
	var lessonsRepo repositories.LessonsRepository = repositories.NewLessonsRepo(db)
	var hwRepo repositories.HomeworkRepository = repositories.NewHomeworkRepo(db)
	var notifRepo repositories.NotificationRepository = repositories.NewNotificationRepo(db)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := application.NewFromEnv(usersRepo, lessonsRepo, hwRepo, notifRepo, ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
