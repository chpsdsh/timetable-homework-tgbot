package formatter

import (
	"fmt"
	"sort"
	"strings"
	"timetable-homework-tgbot/internal/domain/lesson"
)

func FormatGroupTimetable(lessons []lesson.LessonStudent) string {
	if len(lessons) == 0 {
		return "нет пар"
	}

	weekdayOrder := []string{
		"Понедельник",
		"Вторник",
		"Среда",
		"Четверг",
		"Пятница",
		"Суббота",
	}

	byDay := make(map[string][]lesson.LessonStudent)
	for _, l := range lessons {
		byDay[l.Weekday] = append(byDay[l.Weekday], l)
	}

	var b strings.Builder

	for _, day := range weekdayOrder {
		dayLessons := byDay[day]
		if len(dayLessons) == 0 {
			continue
		}

		sort.Slice(dayLessons, func(i, j int) bool {
			return dayLessons[i].StartTime < dayLessons[j].StartTime
		})

		fmt.Fprintf(&b, "%s:\n", day)

		for _, l := range dayLessons {
			line := fmt.Sprintf("  %s %s", l.StartTime, l.Subject)

			if l.LessonType != "" {
				line += fmt.Sprintf(" (%s)", l.LessonType)
			}
			if l.Room != "" {
				line += fmt.Sprintf(", ауд. %s", l.Room)
			}
			if l.Tutor != "" {
				line += fmt.Sprintf(", преп. %s", l.Tutor)
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

func FormatTeacherTimetable(lessons []lesson.LessonTeacher) string {
	if len(lessons) == 0 {
		return "нет пар"
	}

	weekdayOrder := []string{
		"Понедельник",
		"Вторник",
		"Среда",
		"Четверг",
		"Пятница",
		"Суббота",
	}

	byDay := make(map[string][]lesson.LessonTeacher)
	for _, l := range lessons {
		byDay[l.Weekday] = append(byDay[l.Weekday], l)
	}

	var b strings.Builder

	for _, day := range weekdayOrder {
		dayLessons := byDay[day]
		if len(dayLessons) == 0 {
			continue
		}

		sort.Slice(dayLessons, func(i, j int) bool {
			return dayLessons[i].StartTime < dayLessons[j].StartTime
		})

		fmt.Fprintf(&b, "%s:\n", day)

		for _, l := range dayLessons {
			line := fmt.Sprintf("  %s %s", l.StartTime, l.Subject)

			if l.LessonType != "" {
				line += fmt.Sprintf(" (%s)", l.LessonType)
			}

			if len(l.Groups) > 0 {
				line += fmt.Sprintf(", группы: %s", strings.Join(l.Groups, ", "))
			}

			if l.Room != "" {
				line += fmt.Sprintf(", ауд. %s", l.Room)
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

func FormatRoomTimetable(lessons []lesson.LessonRoom) string {
	if len(lessons) == 0 {
		return "нет пар"
	}

	weekdayOrder := []string{
		"Понедельник",
		"Вторник",
		"Среда",
		"Четверг",
		"Пятница",
		"Суббота",
	}

	byDay := make(map[string][]lesson.LessonRoom)
	for _, l := range lessons {
		byDay[l.Weekday] = append(byDay[l.Weekday], l)
	}

	var b strings.Builder

	for _, day := range weekdayOrder {
		dayLessons := byDay[day]
		if len(dayLessons) == 0 {
			continue
		}

		// сортировка по времени внутри дня
		sort.Slice(dayLessons, func(i, j int) bool {
			return dayLessons[i].StartTime < dayLessons[j].StartTime
		})

		fmt.Fprintf(&b, "%s:\n", day)

		for _, l := range dayLessons {
			line := fmt.Sprintf("  %s %s", l.StartTime, l.Subject)

			if l.LessonType != "" {
				line += fmt.Sprintf(" (%s)", l.LessonType)
			}

			if l.Tutor != "" {
				line += fmt.Sprintf(", преп. %s", l.Tutor)
			}

			if len(l.Groups) > 0 {
				line += fmt.Sprintf(", группы: %s", strings.Join(l.Groups, ", "))
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
