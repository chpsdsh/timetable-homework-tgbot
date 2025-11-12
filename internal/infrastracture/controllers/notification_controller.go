package controllers

import (
	"context"
	"sync"
	"time"
)

type NotificationController interface {
	SetWeeklyReminder(ctx context.Context, userID int64, homeworkID string, weekday time.Weekday, hhmm string) error
}

// ФАЛЬСИФИЦИРОВАННО: in-memory.
type reminder struct {
	UserID    int64
	HWID      string
	Weekday   time.Weekday
	HHMM      string
	NextAtUTC time.Time
}

type notify struct {
	mu   sync.Mutex
	data []reminder

	auth AuthController // нужен, чтобы взять TZ пользователя (как в реальной версии)
}

func NewNotificationFake(auth AuthController) NotificationController {
	return &notify{auth: auth}
}

func (n *notify) SetWeeklyReminder(ctx context.Context, userID int64, homeworkID string, weekday time.Weekday, hhmm string) error {
	// TODO(DB): сохранить напоминание (weekly) в БД с расчётом next_at с учётом TZ пользователя
	tz, err := n.auth.UserTZ(ctx, userID)
	if err != nil {
		return err
	}
	loc, _ := time.LoadLocation(tz)
	now := time.Now().In(loc)
	next := computeNextWeekly(weekday, hhmm, now).UTC()

	n.mu.Lock()
	n.data = append(n.data, reminder{
		UserID: userID, HWID: homeworkID, Weekday: weekday, HHMM: hhmm, NextAtUTC: next,
	})
	n.mu.Unlock()
	return nil
}

// утилита (как у тебя)
func computeNextWeekly(wd time.Weekday, hhmm string, from time.Time) time.Time {
	hh := int((hhmm[0]-'0')*10 + (hhmm[1] - '0'))
	mm := int((hhmm[3]-'0')*10 + (hhmm[4] - '0'))
	delta := (int(wd) - int(from.Weekday()) + 7) % 7
	cand := time.Date(from.Year(), from.Month(), from.Day(), hh, mm, 0, 0, from.Location()).AddDate(0, 0, delta)
	if !cand.After(from) {
		cand = cand.AddDate(0, 0, 7)
	}
	return cand
}
