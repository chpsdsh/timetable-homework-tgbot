package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"timetable-homework-tgbot/internal/domain/lesson"
	"timetable-homework-tgbot/internal/domain/urlselector"

	httpparser "timetable-homework-tgbot/internal/infrastracture/parser"
)

const workers = 40

func (d *DB) FillDatabase() error {
	ctx := context.Background()

	log.Println("Filling groups ...")
	if err := d.fillGroupsParallel(ctx); err != nil {
		return err
	}

	log.Println("Filling teachers ...")
	if err := d.fillTeachersParallel(ctx); err != nil {
		return err
	}

	log.Println("Filling rooms ...")
	if err := d.fillRoomsParallel(ctx); err != nil {
		return err
	}

	log.Println("database fill done")
	return nil
}

func (d *DB) fillGroupsParallel(ctx context.Context) error {
	faculties := httpparser.ParseFaculties()
	var allGroups []urlselector.Group
	for _, faculty := range faculties {
		groups := httpparser.ParseGroups(faculty.FullUrl)
		allGroups = append(allGroups, groups...)
	}
	jobs := make(chan urlselector.Group, workers)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for g := range jobs {
				lessons := httpparser.ParseLessonsStudent(g.GroupUrl)
				if err := d.fillGroupSchedule(ctx, g.GroupName, lessons); err != nil {
					select {
					case errCh <- fmt.Errorf("fill group schedule %s: %w", g.GroupName, err):
					default:
					}
					return
				}
			}
		}()
	}

	for _, g := range allGroups {
		jobs <- g
	}
	close(jobs)

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (d *DB) fillTeachersParallel(ctx context.Context) error {
	teachers := httpparser.ParseTeachers()

	jobs := make(chan urlselector.Teacher, workers)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for t := range jobs {
				lessonsTeacher := httpparser.ParseLessonsTeacher(t.FullURL)
				if err := d.fillTeacherSchedule(ctx, t.FullName, lessonsTeacher); err != nil {
					select {
					case errCh <- fmt.Errorf("fill teacher schedule %s: %w", t.FullName, err):
					default:
					}
					return
				}
			}
		}()
	}

	for _, t := range teachers {
		jobs <- t
	}
	close(jobs)
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (d *DB) fillRoomsParallel(ctx context.Context) error {
	rooms := httpparser.ParseRooms()

	jobs := make(chan urlselector.Room, workers)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for r := range jobs {
				lessonsRoom := httpparser.ParseLessonsRoom(r.FllURL)
				if err := d.fillRoomSchedule(ctx, r.Room, lessonsRoom); err != nil {
					select {
					case errCh <- fmt.Errorf("fill room schedule %s: %w", r.Room, err):
					default:
					}
					return
				}
			}
		}()
	}

	for _, r := range rooms {
		jobs <- r
	}
	close(jobs)
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (d *DB) fillGroupSchedule(ctx context.Context, groupName string, lessons []lesson.LessonStudent) error {
	stmt, err := d.sql.PrepareContext(ctx, `
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
	stmt, err := d.sql.PrepareContext(ctx, `
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
	stmt, err := d.sql.PrepareContext(ctx, `
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
