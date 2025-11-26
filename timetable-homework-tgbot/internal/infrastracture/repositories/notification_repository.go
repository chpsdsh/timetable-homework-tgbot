package repositories

import (
	"context"
	"fmt"
	"time"
	"timetable-homework-tgbot/internal/domain"
	"timetable-homework-tgbot/internal/infrastracture/database"
)

type NotificationRepository interface {
	AddNotification(ctx context.Context, userID int64, subject string, ts time.Time) error
	GetPendingNotifications(ctx context.Context, now time.Time) ([]domain.Notification, error)
	DeleteNotification(ctx context.Context, userID int64, subject string, ts time.Time) error
	GetUserNotifications(ctx context.Context, userID int64) ([]domain.Notification, error)
}

type NotificationRepo struct {
	db *database.DB
}

func NewNotificationRepo(db *database.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) AddNotification(
	ctx context.Context,
	userID int64,
	subject string,
	ts time.Time,
) error {
	const q = `
INSERT INTO notifications (user_id, subject, ts)
VALUES ($1, $2, $3)
`
	_, err := r.db.GetSql().ExecContext(ctx, q, userID, subject, ts)
	if err != nil {
		return fmt.Errorf("AddNotification exec: %w", err)
	}
	return nil
}

func (r *NotificationRepo) GetPendingNotifications(
	ctx context.Context,
	now time.Time,
) ([]domain.Notification, error) {

	const q = `
SELECT  user_id, subject, ts
FROM notifications
WHERE ts <= $1
ORDER BY ts
`

	rows, err := r.db.GetSql().QueryContext(ctx, q, now)
	if err != nil {
		return nil, fmt.Errorf("GetPendingNotifications query: %w", err)
	}
	defer rows.Close()

	var res []domain.Notification

	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(
			&n.UserID,
			&n.Subject,
			&n.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("GetPendingNotifications scan: %w", err)
		}
		res = append(res, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetPendingNotifications rows: %w", err)
	}

	return res, nil
}

func (r *NotificationRepo) GetUserNotifications(
	ctx context.Context,
	userID int64,
) ([]domain.Notification, error) {
	const q = `
SELECT user_id, subject, ts
FROM notifications
WHERE user_id = $1
ORDER BY ts;
`

	rows, err := r.db.GetSql().QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserNotifications query: %w", err)
	}
	defer rows.Close()

	var res []domain.Notification

	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(
			&n.UserID,
			&n.Subject,
			&n.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("GetUserNotifications scan: %w", err)
		}
		res = append(res, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUserNotifications rows: %w", err)
	}

	return res, nil
}

func (r *NotificationRepo) DeleteNotification(
	ctx context.Context,
	userID int64,
	subject string,
	ts time.Time,
) error {
	const q = `
DELETE FROM notifications
WHERE user_id = $1
  AND subject = $2
  AND ts = $3;
`
	_, err := r.db.GetSql().ExecContext(ctx, q, userID, subject, ts)
	if err != nil {
		return fmt.Errorf("DeleteNotification exec: %w", err)
	}

	return nil
}
