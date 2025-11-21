package domain

import "time"

type LessonBrief struct {
	Title      string
	LessonType string
	Tutor      string
	StartTime  string
	Weekday    string
	Room       string
	Groups     []string
	Week       string
}

type HWBrief struct {
	ID    string // ID записи ДЗ в БД
	Title string // Как показывать пользователю
}

type Weekday = time.Weekday
