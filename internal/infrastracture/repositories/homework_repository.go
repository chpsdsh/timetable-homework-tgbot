package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"timetable-homework-tgbot/internal/domain"
	"timetable-homework-tgbot/internal/infrastracture/database"
)

type HomeworkRepository interface {
	Save(ctx context.Context, userID int64, subject, text string) error
	Update(ctx context.Context, userID int64, subject, newText string) error
	UpdateStatus(ctx context.Context, userID int64, subject string) error
	Delete(ctx context.Context, userID int64, subject string) error
	ListForLastWeek(ctx context.Context, userID int64) ([]domain.HWBrief, error)
	CheckExistence(ctx context.Context, userID int64, subject string) (bool, error)
}

type HomeworkRepo struct {
	db *database.DB
}

func NewHomeworkRepo(db *database.DB) *HomeworkRepo {
	return &HomeworkRepo{db: db}
}

func (r *HomeworkRepo) Save(
	ctx context.Context,
	userID int64,
	subject, text string,
) error {
	const q = `
INSERT INTO homeworks (id_user, subject, homework_text,status)
VALUES ($1, $2, $3, $4);
`
	_, err := r.db.SQL.ExecContext(ctx, q, userID, subject, text, "new")
	if err != nil {
		return fmt.Errorf("Save homework exec: %w", err)
	}

	return nil
}

func (r *HomeworkRepo) Update(
	ctx context.Context,
	userID int64,
	subject, newText string,
) error {
	const q = `
UPDATE homeworks
SET homework_text = $3
WHERE subject = $2
  AND id_user = $1;
`
	res, err := r.db.SQL.ExecContext(ctx, q, userID, subject, newText)
	if err != nil {
		return fmt.Errorf("Update homework exec: %w", err)
	}

	if rows, err := res.RowsAffected(); err == nil && rows == 0 {
		return fmt.Errorf("Update homework: not found (user=%d, subject=%s)", userID, subject)
	}

	return nil
}

func (r *HomeworkRepo) ListForLastWeek(
	ctx context.Context,
	userID int64,
) ([]domain.HWBrief, error) {

	const q = `
SELECT h.subject, h.homework_text, h.status
FROM homeworks h
WHERE h.id_user = $1;
`

	rows, err := r.db.SQL.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("ListForLastWeek query: %w", err)
	}
	defer rows.Close()

	var res []domain.HWBrief

	for rows.Next() {
		var subj string
		var text string
		var status string

		if err := rows.Scan(&subj, &text, &status); err != nil {
			return nil, fmt.Errorf("ListForLastWeek scan: %w", err)
		}

		res = append(res, domain.HWBrief{
			Subject:      subj,
			HomeworkText: text,
			Status:       status,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListForLastWeek rows: %w", err)
	}

	return res, nil
}

func (r *HomeworkRepo) UpdateStatus(ctx context.Context, userID int64, subject string) error {
	const q = `
UPDATE homeworks
SET status = $3
WHERE subject = $2
  AND id_user = $1;
`
	res, err := r.db.SQL.ExecContext(ctx, q, userID, subject, "done")
	if err != nil {
		return fmt.Errorf("Update homework exec: %w", err)
	}

	if rows, err := res.RowsAffected(); err == nil && rows == 0 {
		return fmt.Errorf("Update homework: not found (user=%d, subject=%s)", userID, subject)
	}
	return nil
}

func (r *HomeworkRepo) Delete(
	ctx context.Context,
	userID int64,
	subject string,
) error {
	const q = `
DELETE FROM homeworks
WHERE id_user = $1
  AND subject = $2;
`
	_, err := r.db.SQL.ExecContext(ctx, q, userID, subject)
	if err != nil {
		return fmt.Errorf("homeworks.Delete exec: %w", err)
	}

	return nil
}

func (r *HomeworkRepo) CheckExistence(
	ctx context.Context,
	userID int64,
	subject string,
) (bool, error) {
	const q = `
SELECT 1
FROM homeworks
WHERE id_user = $1
  AND subject = $2
LIMIT 1;
`

	var dummy int
	err := r.db.SQL.QueryRowContext(ctx, q, userID, subject).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("homeworks.CheckExistence query: %w", err)
	}

	return true, nil
}
