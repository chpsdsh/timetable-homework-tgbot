package database

import (
	"context"
	"fmt"
	"timetable-homework-tgbot/internal/domain/lesson"
)

func (d *DB) fillDatabase() {

}

func (d *DB) fillGroupSchedule(ctx context.Context, groupName string, lessons []lesson.LessonStudent) error {
	stmt, err := d.SQL.PrepareContext(ctx, `
INSERT INTO group_schedule
(group_name, subject, lesson_type, tutor, start_time, weekday, room, week)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, it := range lessons {
		if _, err := stmt.ExecContext(ctx,
			groupName, it.Subject, it.LessonType, it.Tutor, it.StartTime, it.Weekday, it.Room, it.Week,
		); err != nil {
			return fmt.Errorf("insert group_schedule: %w", err)
		}
	}
	return nil
}

func (d *DB) fillTeacherSchedule(ctx context.Context, teacherFIO string, lessons []lesson.LessonTeacher) error {
	stmt, err := d.SQL.PrepareContext(ctx, `
INSERT INTO teacher_schedule
(teacher_fio, subject, lesson_type, "groups", start_time, weekday, room, week)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, it := range lessons {
		if _, err := stmt.ExecContext(ctx,
			teacherFIO, it.Subject, it.LessonType, it.Groups, it.StartTime, it.Weekday, it.Room, it.Week,
		); err != nil {
			return fmt.Errorf("insert teacher_schedule: %w", err)
		}
	}
	return nil
}

func (d *DB) fillRoomSchedule(ctx context.Context, roomName string, lessons []lesson.LessonRoom) error {
	stmt, err := d.SQL.PrepareContext(ctx, `
INSERT INTO room_schedule
(room_name, subject, lesson_type, tutor, start_time, weekday, "groups", week)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, it := range lessons {
		if _, err := stmt.ExecContext(ctx,
			roomName, it.Subject, it.LessonType, it.Tutor, it.StartTime, it.Weekday, it.Groups, it.Week,
		); err != nil {
			return fmt.Errorf("insert room_schedule: %w", err)
		}
	}
	return nil
}
