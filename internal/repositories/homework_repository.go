package repositories

import (
	"context"
	"fmt"
	"strconv"
	"timetable-homework-tgbot/internal/domain"
	"timetable-homework-tgbot/internal/infrastracture/database"
)

type HomeworkRepo struct {
	DB *database.DB
}

type HomeworkRepository interface {
	Save(ctx context.Context, userID int64, lessonID, text string) error
	Update(ctx context.Context, userID int64, lessonID, newText string) error
	ListForLastWeek(ctx context.Context, userID int64) ([]domain.HWBrief, error)
}

func (r *HomeworkRepo) Save(
	ctx context.Context,
	userID int64,
	lessonID, text string,
) error {
	const q = `
INSERT INTO homeworks (id_user, subject, homework_text, status)
SELECT u.id_user, $2, $3, 'new'
FROM users u
WHERE u.tg_id = $1;
`

	res, err := r.DB.SQL.ExecContext(ctx, q, userID, lessonID, text)
	if err != nil {
		return fmt.Errorf("Save homework exec: %w", err)
	}

	if rows, err := res.RowsAffected(); err == nil && rows == 0 {
		return fmt.Errorf("Save homework: user with tg_id=%d not found", userID)
	}

	return nil
}

func (r *HomeworkRepo) Update(
	ctx context.Context,
	userID int64,
	lessonID, newText string,
) error {
	const q = `
UPDATE homeworks h
SET homework_text = $3
WHERE h.subject = $2
  AND h.id_user = (
    SELECT id_user FROM users WHERE tg_id = $1
  );
`
	res, err := r.DB.SQL.ExecContext(ctx, q, userID, lessonID, newText)
	if err != nil {
		return fmt.Errorf("Update homework exec: %w", err)
	}

	if rows, err := res.RowsAffected(); err == nil && rows == 0 {
		return fmt.Errorf("Update homework: not found (user=%d, lessonID=%s)", userID, lessonID)
	}

	return nil
}

func (r *HomeworkRepo) ListForLastWeek(
	ctx context.Context,
	userID int64,
) ([]domain.HWBrief, error) {

	const q = `
SELECT h.id_hw, h.subject
FROM homeworks h
WHERE h.id_user = (
    SELECT id_user FROM users WHERE tg_id = $1
);
`

	rows, err := r.DB.SQL.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("ListForLastWeek query: %w", err)
	}
	defer rows.Close()

	var res []domain.HWBrief

	for rows.Next() {
		var id int64
		var title string

		if err := rows.Scan(&id, &title); err != nil {
			return nil, fmt.Errorf("ListForLastWeek scan: %w", err)
		}

		res = append(res, domain.HWBrief{
			ID:    strconv.FormatInt(id, 10),
			Title: title,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListForLastWeek rows: %w", err)
	}

	return res, nil
}
