package controllers

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

type AuthController interface {
	EnsureJoined(ctx context.Context, userID int64) (bool, error)
	JoinGroup(ctx context.Context, userID int64, group string) error
	LeaveGroup(ctx context.Context, userID int64) error
	UserTZ(ctx context.Context, userID int64) (string, error)
}

var (
	ErrGroupEmpty      = errors.New("group is empty")
	ErrGroupNotExist   = errors.New("group does not exist")
	ErrAlreadyJoined   = errors.New("user already joined")
	ErrInvalidTimezone = errors.New("invalid timezone")
)

// ФАЛЬСИФИЦИРОВАННО: in-memory состояние вместо БД.
type auth struct {
	defaultTZ string
	timeout   time.Duration

	mu     sync.RWMutex
	groups map[int64]string // userID -> group
	tz     map[int64]string // userID -> tz
}

func NewAuthFake(defaultTZ string) AuthController {
	if defaultTZ == "" {
		defaultTZ = "Europe/Bucharest"
	}
	return &auth{
		defaultTZ: defaultTZ,
		timeout:   5 * time.Second,
		groups:    map[int64]string{},
		tz:        map[int64]string{},
	}
}

func (a *auth) EnsureJoined(ctx context.Context, userID int64) (bool, error) {
	// TODO(DB): SELECT group_code FROM users WHERE user_id=$1
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.groups[userID] != "", nil
}

func (a *auth) JoinGroup(ctx context.Context, userID int64, group string) error {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	group = strings.ReplaceAll(strings.TrimSpace(group), " ", "")
	if group == "" {
		return ErrGroupEmpty
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// уже присоединён?
	if cur := a.groups[userID]; cur != "" {
		if strings.EqualFold(cur, group) {
			return nil // идемпотентно
		}
		return ErrAlreadyJoined
	}

	// TODO(DB): проверить существование группы (SELECT 1 FROM groups WHERE code=$1)
	// Временно: считаем, что "00000" — не существует.
	if group == "00000" {
		return ErrGroupNotExist
	}

	// TODO(DB): UPSERT users(user_id, group_code)
	a.groups[userID] = group

	// TODO(DB): if tz is NULL => set default
	if a.tz[userID] == "" {
		a.tz[userID] = a.defaultTZ
	}
	return nil
}

func (a *auth) LeaveGroup(ctx context.Context, userID int64) error {
	// TODO(DB): UPDATE users SET group_code=NULL WHERE user_id=$1
	a.mu.Lock()
	delete(a.groups, userID)
	a.mu.Unlock()
	return nil
}

func (a *auth) UserTZ(ctx context.Context, userID int64) (string, error) {
	// TODO(DB): SELECT tz FROM users WHERE user_id=$1
	a.mu.RLock()
	tz := a.tz[userID]
	a.mu.RUnlock()
	if tz == "" {
		tz = a.defaultTZ
	}
	// базовая проверка валидности зоны
	if _, err := time.LoadLocation(tz); err != nil {
		return "", ErrInvalidTimezone
	}
	return tz, nil
}
