package repositories

import (
	"context"
	"fmt"
	"time"
	"timetable-homework-tgbot/internal/domain"
	"timetable-homework-tgbot/internal/infrastracture/database"
)

type NotificationRepo struct {
	DB *database.DB
}

type NotificationRepository interface {
	AddNotification(ctx context.Context, hwID, userID int64, ts time.Time, weekday string) error
	GetNotifications(ctx context.Context) ([]domain.Notification, error)
}

func (r *NotificationRepo) AddNotification(
	ctx context.Context,
	hwID, userID int64,
	ts time.Time,
	weekday string,
) error {
	const q = `
INSERT INTO notifications (id_hw, user_id, ts, weekday, status)
VALUES ($1, $2, $3, $4, 'pending')
`
	_, err := r.DB.SQL.ExecContext(ctx, q, hwID, userID, ts, weekday)
	if err != nil {
		return fmt.Errorf("AddNotification exec: %w", err)
	}
	return nil
}

func (r *NotificationRepo) GetNotifications(
	ctx context.Context,
) ([]domain.Notification, error) {
	const q = `
SELECT id, id_hw, user_id, ts, weekday, status
FROM notifications
WHERE status = 'pending'
  AND ts <= now()
ORDER BY ts
`
	rows, err := r.DB.SQL.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("GetNotifications query: %w", err)
	}
	defer rows.Close()

	var res []domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(
			&n.ID,
			&n.HomeworkID,
			&n.UserID,
			&n.Timestamp,
			&n.Weekday,
			&n.Status,
		); err != nil {
			return nil, fmt.Errorf("GetNotifications scan: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetNotifications rows: %w", err)
	}

	return res, nil
}
