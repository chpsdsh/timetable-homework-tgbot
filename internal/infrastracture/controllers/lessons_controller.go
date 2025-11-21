package controllers

import (
	"context"
	"timetable-homework-tgbot/internal/repositories"
)

type LessonsController interface {
	GetTimetableGroup(ctx context.Context, group string) string
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
		return "not valid group"
	}

}
