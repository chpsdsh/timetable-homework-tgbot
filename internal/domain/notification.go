package domain

import "time"

type Notification struct {
	ID         int64
	HomeworkID int64
	UserID     int64
	Timestamp  time.Time
	Weekday    string
	Status     string
}
