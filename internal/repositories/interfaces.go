package repositories

import (
	"context"
	"time"
	"timetable-homework-tgbot/internal/domain"
)

type UsersRepository interface {
	GetGroup(ctx context.Context, userID int64) (string, error)
	SetGroup(ctx context.Context, userID int64, group string) error
	UnsetGroup(ctx context.Context, userID int64) error

	GetTZ(ctx context.Context, userID int64) (string, error)
	SetTZ(ctx context.Context, userID int64, tz string) error
}

type LessonsRepository interface {
	GroupExists(ctx context.Context, group string) (bool, error)
	DaysWithLessons(ctx context.Context, group string) ([]string, error)
	LessonsByDay(ctx context.Context, group, day string) ([]domain.LessonBrief, error)
}

type HomeworkRepository interface {
	Save(ctx context.Context, userID int64, group, day, lessonID, text string) error
	Update(ctx context.Context, userID int64, lessonID, newText string) error
	ListForLastWeek(ctx context.Context, userID int64) ([]domain.HWBrief, error)
}

type NotificationRepository interface {
	SaveWeekly(ctx context.Context, userID int64, homeworkID string, weekday time.Weekday, hhmm string, nextAt time.Time) error
}
