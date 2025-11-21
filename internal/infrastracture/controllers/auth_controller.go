package controllers

import (
	"context"
	"errors"
	"strings"
	"time"
	"timetable-homework-tgbot/internal/repositories"
)

type AuthController interface {
	EnsureJoined(ctx context.Context, userID int64) (bool, error)
	JoinGroup(ctx context.Context, userID int64, group string) error
	LeaveGroup(ctx context.Context, userID int64) error
}

var (
	ErrGroupEmpty    = errors.New("group is empty")
	ErrGroupNotExist = errors.New("group does not exist")
)

// ФАЛЬСИФИЦИРОВАННО: in-memory состояние вместо БД.
type auth struct {
	timeout    time.Duration
	userRepo   repositories.UsersRepository
	lessonRepo repositories.LessonsRepository
}

func NewAuthFake(defaultTZ string, userRepo repositories.UsersRepository, lessonRepo repositories.LessonsRepository) AuthController {
	if defaultTZ == "" {
		defaultTZ = "Europe/Bucharest"
	}
	return &auth{
		timeout:    5 * time.Second,
		userRepo:   userRepo,
		lessonRepo: lessonRepo,
	}
}

func (a *auth) EnsureJoined(ctx context.Context, userID int64) (bool, error) {
	group, err := a.userRepo.GetGroup(ctx, userID)
	if err != nil {
		return false, err
	}
	return group != "", nil
}

func (a *auth) JoinGroup(ctx context.Context, userID int64, group string) error {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	group = strings.ReplaceAll(strings.TrimSpace(group), " ", "")
	if group == "" {
		return ErrGroupEmpty
	}

	//проверка существования группы
	if _, err := a.lessonRepo.GetLessonsGroup(ctx, group); err != nil {
		return ErrGroupNotExist
	}

	//добавление группы для юзера
	if err := a.userRepo.SetGroup(ctx, userID, group); err != nil {
		return err
	}

	return nil
}

func (a *auth) LeaveGroup(ctx context.Context, userID int64) error {
	err := a.userRepo.RemoveGroup(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
