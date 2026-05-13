package controllers

import (
	"context"
	"log"
	"timetable-homework-tgbot/internal/application/formatter"
	"timetable-homework-tgbot/internal/infrastracture/repositories"
)

type LessonsController interface {
	GetRole(ctx context.Context, userID int64) (string, error)
	GetTimetableGroup(ctx context.Context, group string) string
	GetTimetableTeacher(ctx context.Context, teacherFio string) string
	GetTimetableRoom(ctx context.Context, room string) string
}

type lessonsController struct {
	auth        AuthController
	lessonsRepo repositories.LessonsRepository
	formatter   *formatter.Formatter
}

func NewLessonController(lessonsRepo repositories.LessonsRepository, userRepo repositories.UsersRepository, controller AuthController) LessonsController {
	return &lessonsController{lessonsRepo: lessonsRepo, auth: controller, formatter: formatter.GetFormatter()}
}

func (l *lessonsController) GetRole(ctx context.Context, userID int64) (string, error) {
	return l.auth.GetRole(ctx, userID)
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
