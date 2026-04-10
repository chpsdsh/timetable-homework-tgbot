package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"timetable-homework-tgbot/internal/infrastracture/config"
	"timetable-homework-tgbot/internal/infrastracture/database"

	"github.com/joho/godotenv"
)

const envFilename = ".env"

func main() {
	if err := godotenv.Load(envFilename); err != nil {
		slog.Error("error loading file", slog.String("file", envFilename), slog.String("err", err.Error()))
		os.Exit(1)
	}
	dbConf, err := config.ParsePostgresConfig()
	if err != nil {
		slog.Error("error parsing config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	db, err := database.NewDB(dbConf)
	if err != nil {
		slog.Error("error connecting to database", slog.String("err", err.Error()))
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()
	db.Close()
}
