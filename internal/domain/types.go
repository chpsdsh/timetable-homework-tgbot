package domain

import "time"

type LessonBrief struct {
	Title      string
	LessonType string
	Tutor      string
	StartTime  string
	Weekday    string
	Room       string
	Groups     string
	Week       string
}

type HWBrief struct {
	Subject      string
	HomeworkText string
	Status       string
}

type Notification struct {
	UserID    int64
	Subject   string
	Timestamp time.Time
}
