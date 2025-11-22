package controllers

import (
	"context"
	"time"
	"timetable-homework-tgbot/internal/repositories"
)

type NotificationController interface {
	SetReminder(ctx context.Context, userID int64, subject, weekday, hhmm string) error
}

type notificationController struct {
	notificationRepo repositories.NotificationRepository
}

func NewNotificationController(notificationRepo repositories.NotificationRepository) NotificationController {
	return &notificationController{notificationRepo: notificationRepo}
}

func (n *notificationController) SetReminder(ctx context.Context, userID int64, subject, date, hhmm string) error {
	t, err := parseNotificationTime(date, hhmm)
	if err != nil {
		return err
	}
	if err := n.notificationRepo.AddNotification(ctx, userID, subject, t); err != nil {
		return err
	}
	return nil
}

func parseNotificationTime(dateStr, timeStr string) (time.Time, error) {
	const layout = "02.01.2006 15:04"

	combined := dateStr + " " + timeStr

	t, err := time.ParseInLocation(layout, combined, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
