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
	AddNotification(ctx context.Context, userID int64, subject string, ts time.Time) error
	GetPendingNotifications(ctx context.Context, now time.Time) ([]domain.Notification, error)
	DeleteNotification(ctx context.Context, hwID, userID int64) error
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
	_, err := r.DB.SQL.ExecContext(ctx, q, userID, subject, ts)
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
SELECT id, user_id, subject, ts, weekday, status
FROM notifications
WHERE ts <= $1
ORDER BY ts
`

	rows, err := r.DB.SQL.QueryContext(ctx, q, now)
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
			&n.Weekday,
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

func (r *NotificationRepo) DeleteNotification(
	ctx context.Context,
	hwID, userID int64,
) error {
	const q = `
DELETE FROM notifications n
WHERE n.user_id = $2
  AND n.subject = (
      SELECT h.subject
      FROM homeworks h
      WHERE h.id_hw = $1
        AND h.id_user = $2
  );
`

	res, err := r.DB.SQL.ExecContext(ctx, q, hwID, userID)
	if err != nil {
		return fmt.Errorf("DeleteNotification exec: %w", err)
	}

	if rows, err := res.RowsAffected(); err == nil && rows == 0 {
		return fmt.Errorf("no notifications for hwID=%d userID=%d", hwID, userID)
	}

	return nil
}
