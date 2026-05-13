package controllers

import (
	"context"
	"errors"
	"strings"
	"time"
	"timetable-homework-tgbot/internal/infrastracture/repositories"
)

var (
	ErrGroupEmpty    = errors.New("group is empty")
	ErrGroupNotExist = errors.New("group does not exist")
)

type AuthController interface {
	GetRole(ctx context.Context, userID int64) (string, error)
	EnsureRegistered(ctx context.Context, userID int64) (bool, error)
	RegisterStudent(ctx context.Context, userID int64, group string) error
	RegisterTeacher(ctx context.Context, userID int64) error
	LeaveGroup(ctx context.Context, userID int64) error
}

type authenticationController struct {
	timeout    time.Duration
	userRepo   repositories.UsersRepository
	lessonRepo repositories.LessonsRepository
}

func NewAuthController(userRepo repositories.UsersRepository, lessonRepo repositories.LessonsRepository) AuthController {
	return &authenticationController{
		timeout:    5 * time.Second,
		userRepo:   userRepo,
		lessonRepo: lessonRepo,
	}
}

func (a *authenticationController) GetRole(ctx context.Context, userID int64) (string, error) {
	return a.userRepo.GetRole(ctx, userID)
}

func (a *authenticationController) EnsureRegistered(ctx context.Context, userID int64) (bool, error) {
	role, err := a.userRepo.GetRole(ctx, userID)
	if err != nil {
		return false, err
	}
	return role != repositories.RoleGuest, nil
}

func (a *authenticationController) RegisterStudent(ctx context.Context, userID int64, group string) error {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	group = strings.ReplaceAll(strings.TrimSpace(group), " ", "")
	if group == "" {
		return ErrGroupEmpty
	}

	lessons, err := a.lessonRepo.GetLessonsGroup(ctx, group)
	if err != nil || len(lessons) == 0 {
		return ErrGroupNotExist
	}

	if err := a.userRepo.RegisterStudent(ctx, userID, group); err != nil {
		return err
	}

	return nil
}

func (a *authenticationController) RegisterTeacher(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	return a.userRepo.RegisterTeacher(ctx, userID)
}

func (a *authenticationController) LeaveGroup(ctx context.Context, userID int64) error {
	err := a.userRepo.Leave(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
