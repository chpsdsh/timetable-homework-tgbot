package formatter

import (
	"fmt"
	"strings"
	"sync"
	"timetable-homework-tgbot/internal/domain"
)

type Formatter struct {
}

var instance *Formatter

var once sync.Once

func GetFormatter() *Formatter {
	once.Do(func() {
		instance = &Formatter{}
	})
	return instance
}

const telegramMaxLen = 4000

func (f *Formatter) FormatTimetable(lessons []domain.LessonBrief) string {
	if len(lessons) == 0 {
		return "расписание не найдено"
	}

	weekdayOrder := []string{
		"Понедельник",
		"Вторник",
		"Среда",
		"Четверг",
		"Пятница",
		"Суббота",
	}

	byDay := make(map[string][]domain.LessonBrief)
	for _, l := range lessons {
		byDay[l.Weekday] = append(byDay[l.Weekday], l)
	}

	var b strings.Builder

	for _, day := range weekdayOrder {
		dayLessons := byDay[day]
		if len(dayLessons) == 0 {
			continue
		}

		fmt.Fprintf(&b, "%s:\n", day)

		for _, l := range dayLessons {
			line := fmt.Sprintf("  %s %s", l.StartTime, l.Title)

			if l.LessonType != "" {
				line += fmt.Sprintf(" (%s)", l.LessonType)
			}
			if l.Room != "" {
				line += fmt.Sprintf(", ауд. %s", l.Room)
			}
			if l.Tutor != "" {
				line += fmt.Sprintf(", преп. %s", l.Tutor)
			}
			if l.Groups != "" {
				line += fmt.Sprintf(", группы: %s", l.Groups)
			}
			if l.Week != "" {
				line += fmt.Sprintf(", неделя: %s", l.Week)
			}

			b.WriteString(line)
			b.WriteByte('\n')
		}

		b.WriteByte('\n')
	}

	res := strings.TrimSpace(b.String())
	if res == "" {
		return "нет пар"
	}
	return res
}

func (f *Formatter) SplitForTelegram(text string) []string {
	if len(text) <= telegramMaxLen {
		return []string{text}
	}

	lines := strings.Split(text, "\n")

	var res []string
	var b strings.Builder

	for _, line := range lines {
		if b.Len()+len(line)+1 > telegramMaxLen {
			res = append(res, b.String())
			b.Reset()
		}

		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(line)
	}

	if b.Len() > 0 {
		res = append(res, b.String())
	}

	return res
}

func (f *Formatter) FormatHomeworks(list []domain.HWBrief) string {
	if len(list) == 0 {
		return "Домашек нет"
	}

	var b strings.Builder
	for _, h := range list {
		fmt.Fprintf(&b, "%s: %s(%s)\n", h.Subject, h.HomeworkText, h.Status)
	}

	return strings.TrimSpace(b.String())
}
