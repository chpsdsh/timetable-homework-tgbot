package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"timetable-homework-tgbot/internal/infrastracture/database"
)

type UsersRepository interface {
	GetGroup(ctx context.Context, userID int64) (string, error)
	SetGroup(ctx context.Context, userID int64, group string) error
	RemoveGroup(ctx context.Context, userID int64) error
}

type UserRepo struct {
	db *database.DB
}

func NewUserRepo(db *database.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetGroup(ctx context.Context, userID int64) (string, error) {
	var group sql.NullString

	err := r.db.GetSql().QueryRowContext(ctx,
		`SELECT "group" FROM users WHERE tg_id = $1`,
		userID,
	).Scan(&group)

	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	
	if err != nil {
		return "", fmt.Errorf("GetGroup query: %w", err)
	}

	if group.Valid {
		return group.String, nil
	}
	return "", nil
}

func (r *UserRepo) SetGroup(ctx context.Context, userID int64, group string) error {
	_, err := r.db.GetSql().ExecContext(ctx, `
INSERT INTO users (tg_id, "group")
VALUES ($1, $2)
ON CONFLICT (tg_id) DO UPDATE SET "group" = EXCLUDED."group"
`,
		userID, group,
	)
	if err != nil {
		return fmt.Errorf("SetGroup exec: %w", err)
	}
	return nil
}

func (r *UserRepo) RemoveGroup(ctx context.Context, userID int64) error {
	_, err := r.db.GetSql().ExecContext(ctx,
		`UPDATE users SET "group" = NULL WHERE tg_id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("RemoveGroup exec: %w", err)
	}
	return nil
}
