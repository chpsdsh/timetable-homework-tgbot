package controllers

import (
	"context"
	"errors"
	"time"
	"timetable-homework-tgbot/internal/infrastracture/repositories"
)

type AuthController interface {
	GetRole(ctx context.Context, userID int64) (string, error)
	EnsureRegistered(ctx context.Context, userID int64) (bool, error)
	Login(ctx context.Context, userID int64, role, login, password string) error
	LeaveGroup(ctx context.Context, userID int64) error
}

type authenticationController struct {
	timeout  time.Duration
	userRepo repositories.UsersRepository
}

func NewAuthController(userRepo repositories.UsersRepository) AuthController {
	return &authenticationController{
		timeout:  5 * time.Second,
		userRepo: userRepo,
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

func (a *authenticationController) Login(ctx context.Context, userID int64, role, login, password string) error {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	switch role {
	case repositories.RoleStudent:
		return a.userRepo.LoginStudent(ctx, userID, login, password)
	case repositories.RoleTeacher:
		return a.userRepo.LoginTeacher(ctx, userID, login, password)
	default:
		return errors.New("unknown role")
	}
}

func (a *authenticationController) LeaveGroup(ctx context.Context, userID int64) error {
	err := a.userRepo.Leave(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
