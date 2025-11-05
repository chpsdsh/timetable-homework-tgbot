package app

import (
	"context"
	"os/signal"
	"syscall"
	"timetable-homework-tgbot/internal/infrastracture/telegram"
)

func Run() error {
	bot, err := telegram.NewBot()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return bot.Run(ctx)
}
