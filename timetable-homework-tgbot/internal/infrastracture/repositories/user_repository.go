package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"timetable-homework-tgbot/internal/infrastracture/database"
)

const (
	RoleGuest   = ""
	RoleStudent = "student"
	RoleTeacher = "teacher"
)

type UsersRepository interface {
	GetRole(ctx context.Context, userID int64) (string, error)
	GetGroup(ctx context.Context, userID int64) (string, error)
	RegisterStudent(ctx context.Context, userID int64, group string) error
	RegisterTeacher(ctx context.Context, userID int64) error
	Leave(ctx context.Context, userID int64) error
}

type UserRepo struct {
	db *database.DB
}

func NewUserRepo(db *database.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetRole(ctx context.Context, userID int64) (string, error) {
	var exists int
	err := r.db.GetSql().QueryRowContext(ctx, `SELECT 1 FROM students WHERE user_id = $1`, userID).Scan(&exists)
	if err == nil {
		return RoleStudent, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("GetRole student query: %w", err)
	}

	err = r.db.GetSql().QueryRowContext(ctx, `SELECT 1 FROM teachers WHERE user_id = $1`, userID).Scan(&exists)
	if err == nil {
		return RoleTeacher, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("GetRole teacher query: %w", err)
	}

	return RoleGuest, nil
}

func (r *UserRepo) GetGroup(ctx context.Context, userID int64) (string, error) {
	var group string

	err := r.db.GetSql().QueryRowContext(ctx, `
SELECT g.group_name
FROM students s
JOIN groups g ON g.id = s.group_id
WHERE s.user_id = $1
`, userID).Scan(&group)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("GetGroup query: %w", err)
	}

	return group, nil
}

func (r *UserRepo) RegisterStudent(ctx context.Context, userID int64, group string) error {
	group = strings.TrimSpace(group)
	tx, err := r.db.GetSql().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("RegisterStudent begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `INSERT INTO users (tg_id) VALUES ($1) ON CONFLICT DO NOTHING`, userID); err != nil {
		return fmt.Errorf("RegisterStudent user: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM teachers WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("RegisterStudent delete teacher: %w", err)
	}

	var groupID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM groups WHERE group_name = $1`, group).Scan(&groupID); err != nil {
		return fmt.Errorf("RegisterStudent group: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO students (user_id, group_id)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE SET group_id = EXCLUDED.group_id
`, userID, groupID); err != nil {
		return fmt.Errorf("RegisterStudent student: %w", err)
	}

	return tx.Commit()
}

func (r *UserRepo) RegisterTeacher(ctx context.Context, userID int64) error {
	tx, err := r.db.GetSql().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("RegisterTeacher begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `INSERT INTO users (tg_id) VALUES ($1) ON CONFLICT DO NOTHING`, userID); err != nil {
		return fmt.Errorf("RegisterTeacher user: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM students WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("RegisterTeacher delete student: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO teachers (user_id)
VALUES ($1)
ON CONFLICT (user_id) DO NOTHING
`, userID); err != nil {
		return fmt.Errorf("RegisterTeacher teacher: %w", err)
	}

	return tx.Commit()
}

func (r *UserRepo) Leave(ctx context.Context, userID int64) error {
	tx, err := r.db.GetSql().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Leave begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM students WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("Leave delete student: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM teachers WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("Leave delete teacher: %w", err)
	}

	return tx.Commit()
}
