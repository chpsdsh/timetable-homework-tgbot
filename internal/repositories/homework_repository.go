package repositories

import (
	"context"
	"fmt"
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
INSERT INTO homeworks (id_user, subject, homework_text)
VALUES ($1, $2, $3);
`
	_, err := r.DB.SQL.ExecContext(ctx, q, userID, lessonID, text)
	if err != nil {
		return fmt.Errorf("Save homework exec: %w", err)
	}

	return nil
}

func (r *HomeworkRepo) Update(
	ctx context.Context,
	userID int64, // tg_id
	lessonID, newText string,
) error {
	const q = `
UPDATE homeworks
SET homework_text = $3
WHERE subject = $2
  AND id_user = $1;
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
SELECT h.subject, h.homework_text
FROM homeworks h
WHERE h.id_user = $1;
`

	rows, err := r.DB.SQL.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("ListForLastWeek query: %w", err)
	}
	defer rows.Close()

	var res []domain.HWBrief

	for rows.Next() {
		var subj string
		var text string

		if err := rows.Scan(&subj, &text); err != nil {
			return nil, fmt.Errorf("ListForLastWeek scan: %w", err)
		}

		res = append(res, domain.HWBrief{
			Subject:      subj,
			HomeworkText: text,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListForLastWeek rows: %w", err)
	}

	return res, nil
}
