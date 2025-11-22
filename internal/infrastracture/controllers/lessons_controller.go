package controllers

import (
	"context"
	"log"
	"timetable-homework-tgbot/internal/infrastracture/formatter"
	"timetable-homework-tgbot/internal/repositories"
)

type LessonsController interface {
	GetTimetableGroup(ctx context.Context, group string) string
	GetTimetableTeacher(ctx context.Context, teacherFio string) string
	GetTimetableRoom(ctx context.Context, room string) string
}

type lessonsController struct {
	lessonsRepo repositories.LessonsRepository
}

func NewLessonController(lessonsRepo repositories.LessonsRepository) LessonsController {
	return &lessonsController{lessonsRepo: lessonsRepo}
}

func (l *lessonsController) GetTimetableGroup(ctx context.Context, group string) string {
	lessons, err := l.lessonsRepo.GetLessonsGroup(ctx, group)
	if err != nil {
		log.Println(err.Error())
		return "not valid group"
	}
	timetable := formatter.FormatTimetable(lessons)
	return timetable
}

func (l *lessonsController) GetTimetableTeacher(ctx context.Context, teacherFio string) string {
	lessons, err := l.lessonsRepo.GetLessonsTeacher(ctx, teacherFio)
	if err != nil {
		log.Println(err.Error())
		return "not valid teacher"
	}
	timetable := formatter.FormatTimetable(lessons)
	return timetable
}

func (l *lessonsController) GetTimetableRoom(ctx context.Context, room string) string {
	lessons, err := l.lessonsRepo.GetLessonsRoom(ctx, room)
	if err != nil {
		log.Println(err.Error())
		return "not valid room"
	}
	timetable := formatter.FormatTimetable(lessons)
	return timetable
}
