package controllers

import (
	"context"
	"log"
	"timetable-homework-tgbot/internal/application/formatter"
	"timetable-homework-tgbot/internal/infrastracture/repositories"
)

type LessonsController interface {
	EnsureJoined(ctx context.Context, userID int64) (bool, error)
	GetTimetableGroup(ctx context.Context, group string) string
	GetTimetableTeacher(ctx context.Context, teacherFio string) string
	GetTimetableRoom(ctx context.Context, room string) string
}

type lessonsController struct {
	auth        AuthController
	lessonsRepo repositories.LessonsRepository
	formatter   formatter.Formatter
}

func NewLessonController(lessonsRepo repositories.LessonsRepository, userRepo repositories.UsersRepository, controller AuthController) LessonsController {
	return &lessonsController{lessonsRepo: lessonsRepo, auth: controller, formatter: *formatter.NewFormatter()}
}

func (l *lessonsController) EnsureJoined(ctx context.Context, userID int64) (bool, error) {
	b, err := l.auth.EnsureJoined(ctx, userID)
	if err != nil {
		return false, err
	}
	return b, nil
}

func (l *lessonsController) GetTimetableGroup(ctx context.Context, group string) string {
	lessons, err := l.lessonsRepo.GetLessonsGroup(ctx, group)
	if err != nil {
		log.Println(err.Error())
		return "not valid group"
	}
	timetable := l.formatter.FormatTimetable(lessons)
	return timetable
}

func (l *lessonsController) GetTimetableTeacher(ctx context.Context, teacherFio string) string {
	lessons, err := l.lessonsRepo.GetLessonsTeacher(ctx, teacherFio)
	if err != nil {
		log.Println(err.Error())
		return "not valid teacher"
	}
	timetable := l.formatter.FormatTimetable(lessons)
	return timetable
}

func (l *lessonsController) GetTimetableRoom(ctx context.Context, room string) string {
	lessons, err := l.lessonsRepo.GetLessonsRoom(ctx, room)
	if err != nil {
		log.Println(err.Error())
		return "not valid room"
	}
	timetable := l.formatter.FormatTimetable(lessons)
	return timetable
}
