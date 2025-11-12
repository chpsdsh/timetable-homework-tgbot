package controllers

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"timetable-homework-tgbot/internal/domain"
)

type HomeworkController interface {
	DaysWithLessons(ctx context.Context, userID int64) ([]string, error)
	LessonsByDay(ctx context.Context, userID int64, day string) ([]domain.LessonBrief, error)

	Pin(ctx context.Context, userID int64, day string, lessonID, text string) error
	Update(ctx context.Context, userID int64, lessonID, newText string) error

	ListForLastWeek(ctx context.Context, userID int64) ([]domain.HWBrief, error)
}

// ФАЛЬСИФИЦИРОВАННО: in-memory.
type hw struct {
	timeout time.Duration

	mu     sync.RWMutex
	hwText map[int64]map[string]string // userID -> lessonID -> text
	hwLog  map[int64][]domain.HWBrief  // userID -> последние записи (для меню напоминаний)
}

func NewHomeworkFake() HomeworkController {
	return &hw{
		timeout: 5 * time.Second,
		hwText:  map[int64]map[string]string{},
		hwLog:   map[int64][]domain.HWBrief{},
	}
}

func (c *hw) DaysWithLessons(ctx context.Context, userID int64) ([]string, error) {
	// TODO(DB): достать из LessonsRepository дни, где у группы пользователя есть пары
	return []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница"}, nil
}

func (c *hw) LessonsByDay(ctx context.Context, userID int64, day string) ([]domain.LessonBrief, error) {
	// TODO(DB): SELECT пары для группы пользователя в указанный день
	day = strings.TrimSpace(strings.ToLower(day))
	switch day {
	case "понедельник":
		return []domain.LessonBrief{
			{ID: "123", Title: "09:00–10:35 • Математика"},
			{ID: "124", Title: "10:45–12:20 • Физика"},
		}, nil
	default:
		return []domain.LessonBrief{
			{ID: "125", Title: "09:00–10:35 • Информатика"},
		}, nil
	}
}

func (c *hw) Pin(ctx context.Context, userID int64, day string, lessonID, text string) error {
	// TODO(DB): INSERT INTO homework(user_id, day, lesson_id, text, ...)
	c.mu.Lock()
	if c.hwText[userID] == nil {
		c.hwText[userID] = map[string]string{}
	}
	c.hwText[userID][lessonID] = text
	c.hwLog[userID] = append(c.hwLog[userID], domain.HWBrief{
		ID:    "h" + lessonID,
		Title: fmt.Sprintf("%s • %s", day, text),
	})
	c.mu.Unlock()
	return nil
}

func (c *hw) Update(ctx context.Context, userID int64, lessonID, newText string) error {
	// TODO(DB): UPDATE homework SET text=$1 WHERE user_id=$2 AND lesson_id=$3
	c.mu.Lock()
	if c.hwText[userID] == nil {
		c.hwText[userID] = map[string]string{}
	}
	c.hwText[userID][lessonID] = newText
	c.mu.Unlock()
	return nil
}

func (c *hw) ListForLastWeek(ctx context.Context, userID int64) ([]domain.HWBrief, error) {
	// TODO(DB): SELECT ... WHERE created_at > now()-7d ORDER BY created_at
	c.mu.RLock()
	defer c.mu.RUnlock()
	logs := c.hwLog[userID]
	n := len(logs)
	if n > 5 {
		return append([]domain.HWBrief(nil), logs[n-5:]...), nil
	}
	return append([]domain.HWBrief(nil), logs...), nil
}
