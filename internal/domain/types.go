package domain

import "time"

type LessonBrief struct {
	ID    string
	Title string
}

type HWBrief struct {
	ID    string // ID записи ДЗ в БД
	Title string // как показывать пользователю
}

type Weekday = time.Weekday
