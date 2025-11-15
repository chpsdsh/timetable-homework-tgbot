package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"timetable-homework-tgbot/internal/domain"
	"timetable-homework-tgbot/internal/infrastracture/database"
)

type LessonsRepository interface {
	GetLessonsGroup(ctx context.Context, group string) ([]string, error)
	GetLessonsTeacher(ctx context.Context, teacherFio string) ([]string, error)
	GetLessonsRoom(ctx context.Context, roomName string) ([]string, error)
	LessonsByDayGroup(ctx context.Context, group, day string) ([]domain.LessonBrief, error)
	LessonsByDayTeacher(ctx context.Context, teacherFio, day string) ([]domain.LessonBrief, error)
	LessonsByDayRoom(ctx context.Context, roomName, day string) ([]domain.LessonBrief, error)
}

type LessonsRepo struct {
	DB *database.DB
}

func (r *LessonsRepo) GetLessonsGroup(ctx context.Context, group string) ([]string, error) {
	const q = `
SELECT  subject
FROM group_schedule
WHERE group_name = $1
ORDER BY subject;
`
	rows, err := r.DB.SQL.QueryContext(ctx, q, group)
	if err != nil {
		return nil, fmt.Errorf("GetLessonsGroup query: %w", err)
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var subj string
		if err := rows.Scan(&subj); err != nil {
			return nil, fmt.Errorf("GetLessonsGroup scan: %w", err)
		}
		res = append(res, subj)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetLessonsGroup rows: %w", err)
	}
	return res, nil
}

func (r *LessonsRepo) GetLessonsTeacher(ctx context.Context, teacherFio string) ([]string, error) {
	const q = `
SELECT  subject
FROM teacher_schedule
WHERE teacher_fio = $1
ORDER BY subject;
`
	rows, err := r.DB.SQL.QueryContext(ctx, q, teacherFio)
	if err != nil {
		return nil, fmt.Errorf("GetLessonsTeacher query: %w", err)
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var subj string
		if err := rows.Scan(&subj); err != nil {
			return nil, fmt.Errorf("GetLessonsTeacher scan: %w", err)
		}
		res = append(res, subj)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetLessonsTeacher rows: %w", err)
	}
	return res, nil
}

func (r *LessonsRepo) GetLessonsRoom(ctx context.Context, roomName string) ([]string, error) {
	const q = `
SELECT subject
FROM room_schedule
WHERE room_name = $1
ORDER BY subject;
`
	rows, err := r.DB.SQL.QueryContext(ctx, q, roomName)
	if err != nil {
		return nil, fmt.Errorf("GetLessonsRoom query: %w", err)
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var subj string
		if err := rows.Scan(&subj); err != nil {
			return nil, fmt.Errorf("GetLessonsRoom scan: %w", err)
		}
		res = append(res, subj)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetLessonsRoom rows: %w", err)
	}
	return res, nil
}

func (r *LessonsRepo) LessonsByDayGroup(
	ctx context.Context,
	group, day string,
) ([]domain.LessonBrief, error) {
	const q = `
SELECT subject, lesson_type, tutor, start_time, weekday, room, week
FROM group_schedule
WHERE group_name = $1
  AND weekday   = $2
ORDER BY start_time;
`

	rows, err := r.DB.SQL.QueryContext(ctx, q, group, day)
	if err != nil {
		return nil, fmt.Errorf("LessonsByDayGroup query: %w", err)
	}
	defer rows.Close()

	var res []domain.LessonBrief

	for rows.Next() {
		var (
			subject    string
			lessonType sql.NullString
			tutor      sql.NullString
			startTime  time.Time
			weekday    string
			room       sql.NullString
			week       sql.NullString
		)

		if err := rows.Scan(
			&subject,
			&lessonType,
			&tutor,
			&startTime,
			&weekday,
			&room,
			&week,
		); err != nil {
			return nil, fmt.Errorf("LessonsByDayGroup scan: %w", err)
		}

		res = append(res, domain.LessonBrief{
			Title:      subject,
			LessonType: nullToString(lessonType),
			Tutor:      nullToString(tutor),
			StartTime:  startTime.Format("15:04"),
			Weekday:    weekday,
			Room:       nullToString(room),
			Groups:     []string{group}, // для группы можно подставить саму группу
			Week:       nullToString(week),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("LessonsByDayGroup rows: %w", err)
	}

	return res, nil
}

func (r *LessonsRepo) LessonsByDayTeacher(
	ctx context.Context,
	teacherFio, day string,
) ([]domain.LessonBrief, error) {
	const q = `
SELECT subject, lesson_type, "groups", start_time, weekday, room, week
FROM teacher_schedule
WHERE teacher_fio = $1
  AND weekday     = $2
ORDER BY start_time;
`

	rows, err := r.DB.SQL.QueryContext(ctx, q, teacherFio, day)
	if err != nil {
		return nil, fmt.Errorf("LessonsByDayTeacher query: %w", err)
	}
	defer rows.Close()

	var res []domain.LessonBrief

	for rows.Next() {
		var (
			subject    string
			lessonType sql.NullString
			groups     []string
			startTime  time.Time
			weekday    string
			room       sql.NullString
			week       sql.NullString
		)

		if err := rows.Scan(
			&subject,
			&lessonType,
			&groups,
			&startTime,
			&weekday,
			&room,
			&week,
		); err != nil {
			return nil, fmt.Errorf("LessonsByDayTeacher scan: %w", err)
		}

		res = append(res, domain.LessonBrief{
			Title:      subject,
			LessonType: nullToString(lessonType),
			Tutor:      teacherFio,
			StartTime:  startTime.Format("15:04"),
			Weekday:    weekday,
			Room:       nullToString(room),
			Groups:     groups,
			Week:       nullToString(week),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("LessonsByDayTeacher rows: %w", err)
	}

	return res, nil
}

func (r *LessonsRepo) LessonsByDayRoom(
	ctx context.Context,
	roomName, day string,
) ([]domain.LessonBrief, error) {
	const q = `
SELECT subject, lesson_type, tutor, start_time, weekday, "groups", week
FROM room_schedule
WHERE room_name = $1
  AND weekday   = $2
ORDER BY start_time;
`

	rows, err := r.DB.SQL.QueryContext(ctx, q, roomName, day)
	if err != nil {
		return nil, fmt.Errorf("LessonsByDayRoom query: %w", err)
	}
	defer rows.Close()

	var res []domain.LessonBrief

	for rows.Next() {
		var (
			subject    string
			lessonType sql.NullString
			tutor      sql.NullString
			startTime  time.Time
			weekday    string
			groups     []string
			week       sql.NullString
		)

		if err := rows.Scan(
			&subject,
			&lessonType,
			&tutor,
			&startTime,
			&weekday,
			&groups,
			&week,
		); err != nil {
			return nil, fmt.Errorf("LessonsByDayRoom scan: %w", err)
		}

		res = append(res, domain.LessonBrief{
			Title:      subject,
			LessonType: nullToString(lessonType),
			Tutor:      nullToString(tutor),
			StartTime:  startTime.Format("15:04"),
			Weekday:    weekday,
			Room:       roomName,
			Groups:     groups,
			Week:       nullToString(week),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("LessonsByDayRoom rows: %w", err)
	}

	return res, nil
}

func nullToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
