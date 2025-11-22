package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"
	"timetable-homework-tgbot/internal/domain"
	"timetable-homework-tgbot/internal/repositories"
)

type NotificationController interface {
	SetReminder(ctx context.Context, userID int64, subject, weekday, hhmm string) error
	GetUserNotifications(ctx context.Context, userID int64) ([]domain.Notification, error)
	GetPendingNotifications(ctx context.Context) ([]domain.Notification, error)
	DeleteUserNotification(ctx context.Context, userID int64, notification string) error
	DeleteUserNotificationWithTs(ctx context.Context, userID int64, subject string, ts time.Time) error
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

func (n *notificationController) GetUserNotifications(ctx context.Context, userID int64) ([]domain.Notification, error) {
	not, err := n.notificationRepo.GetUserNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}
	return not, nil
}

func (n *notificationController) DeleteUserNotification(ctx context.Context, userID int64, notification string) error {
	subject, t, err := parseNotificationLabel(notification)
	if err != nil {
		return err
	}
	if err := n.notificationRepo.DeleteNotification(ctx, userID, subject, t); err != nil {
		return err
	}
	return nil
}

func (n *notificationController) DeleteUserNotificationWithTs(ctx context.Context, userID int64, subject string, ts time.Time) error {
	if err := n.notificationRepo.DeleteNotification(ctx, userID, subject, ts); err != nil {
		return err
	}
	return nil
}

func (n *notificationController) GetPendingNotifications(ctx context.Context) ([]domain.Notification, error) {
	not, err := n.notificationRepo.GetPendingNotifications(ctx, time.Now())
	if err != nil {
		return nil, err
	}
	return not, nil
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

func parseNotificationLabel(label string) (string, time.Time, error) {
	parts := strings.SplitN(label, " — ", 2)
	if len(parts) != 2 {
		return "", time.Time{}, fmt.Errorf("invalid label format: %q", label)
	}

	subject := strings.TrimSpace(parts[0])
	tsStr := strings.TrimSpace(parts[1])

	const layout = "02.01.2006 15:04"
	ts, err := time.ParseInLocation(layout, tsStr, time.Local)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("parse time %q: %w", tsStr, err)
	}

	return subject, ts, nil
}
