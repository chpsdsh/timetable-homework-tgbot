package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

func New(ctx context.Context) (*DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &DB{SQL: db}, nil
}

func (d *DB) Close() error {
	if d == nil || d.SQL == nil {
		return nil
	}
	return d.SQL.Close()
}
func (d *DB) InitSchema(ctx context.Context) error {
	const schema = `
-- Пользователи
CREATE TABLE IF NOT EXISTS users (
  id_user     BIGSERIAL PRIMARY KEY,
  tg_id       BIGINT UNIQUE NOT NULL,
  "group"     TEXT
);

CREATE INDEX IF NOT EXISTS idx_users_group ON users("group");

-- Домашки
CREATE TABLE IF NOT EXISTS homeworks (
  id_hw         BIGSERIAL PRIMARY KEY,
  id_user       BIGINT NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  subject       TEXT NOT NULL,
  homework_text TEXT NOT NULL,
  status        TEXT NOT NULL DEFAULT 'new'
);

CREATE INDEX IF NOT EXISTS idx_hw_user ON homeworks(id_user);
CREATE INDEX IF NOT EXISTS idx_hw_status ON homeworks(status);

-- Уведомления по домашкам
CREATE TABLE IF NOT EXISTS notifications (
  id            BIGSERIAL PRIMARY KEY,
  id_hw         BIGINT NOT NULL REFERENCES homeworks(id_hw) ON DELETE CASCADE,
  user_id       BIGINT NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  ts            TIMESTAMPTZ NOT NULL,
  status        TEXT NOT NULL DEFAULT 'pending'
);

CREATE INDEX IF NOT EXISTS idx_notif_user_ts ON notifications(user_id, ts);
CREATE INDEX IF NOT EXISTS idx_notif_status ON notifications(status);

-- Расписание групп
CREATE TABLE IF NOT EXISTS group_schedule (
  group_name  TEXT NOT NULL,
  subject     TEXT NOT NULL,
  lesson_type TEXT,
  tutor       TEXT,
  start_time  TIME NOT NULL,
  weekday     TEXT NOT NULL CHECK (weekday IN
               ('Понедельник','Вторник','Среда','Четверг','Пятница','Суббота')),
  room        TEXT,
  week        TEXT
);

CREATE INDEX IF NOT EXISTS idx_group_sched_lookup
  ON group_schedule(group_name, weekday, start_time);

-- Расписание преподавателей
CREATE TABLE IF NOT EXISTS teacher_schedule (
  teacher_fio TEXT NOT NULL,
  subject     TEXT NOT NULL,
  lesson_type TEXT,
  "groups"    TEXT[] NOT NULL,
  start_time  TIME NOT NULL,
  weekday     TEXT NOT NULL CHECK (weekday IN
               ('Понедельник','Вторник','Среда','Четверг','Пятница','Суббота')),
  room        TEXT,
  week        TEXT
);

CREATE INDEX IF NOT EXISTS idx_teacher_sched_lookup
  ON teacher_schedule(teacher_fio, weekday, start_time);

-- Расписание аудиторий
CREATE TABLE IF NOT EXISTS room_schedule (
  room_name   TEXT NOT NULL,
  subject     TEXT NOT NULL,
  lesson_type TEXT,
  tutor       TEXT,
  start_time  TIME NOT NULL,
  weekday     TEXT NOT NULL CHECK (weekday IN
               ('Понедельник','Вторник','Среда','Четверг','Пятница','Суббота')),
  "groups"    TEXT[],
  week        TEXT
);

CREATE INDEX IF NOT EXISTS idx_room_sched_lookup
  ON room_schedule(room_name, weekday, start_time);
`

	_, err := d.SQL.ExecContext(ctx, schema)
	return err
}
