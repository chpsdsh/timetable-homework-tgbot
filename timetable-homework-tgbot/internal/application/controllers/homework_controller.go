package controllers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"timetable-homework-tgbot/internal/infrastracture/repositories"

	"timetable-homework-tgbot/internal/domain"
)

type HomeworkController interface {
	GetRole(ctx context.Context, userID int64) (string, error)
	StudentGroup(ctx context.Context, userID int64) (string, error)
	DaysWithLessons(ctx context.Context, group string) ([]string, error)
	LessonsByDay(ctx context.Context, group, day string) ([]domain.LessonBrief, error)
	Pin(ctx context.Context, group string, day string, subject, text string) error
	Update(ctx context.Context, homeworkID int64, newText string) error
	UpdateStatus(ctx context.Context, homeworkID int64) error
	DeleteHomework(ctx context.Context, homeworkID int64) error
	ListAll(ctx context.Context) ([]domain.HWBrief, error)
	ListForLastWeek(ctx context.Context, userID int64) ([]domain.HWBrief, error)
	CheckExistence(ctx context.Context, userID int64, homeworkID int64) (bool, error)
}

type homeworkController struct {
	timeout time.Duration

	usersRepo    repositories.UsersRepository
	homeworkRepo repositories.HomeworkRepository
	lessonsRepo  repositories.LessonsRepository
}

func NewHomeworkController(usersRepo repositories.UsersRepository, homeworkRepo repositories.HomeworkRepository,
	lessonsRepo repositories.LessonsRepository) HomeworkController {
	return &homeworkController{
		timeout:      5 * time.Second,
		usersRepo:    usersRepo,
		homeworkRepo: homeworkRepo,
		lessonsRepo:  lessonsRepo,
	}
}

func (c *homeworkController) GetRole(ctx context.Context, userID int64) (string, error) {
	return c.usersRepo.GetRole(ctx, userID)
}

func (c *homeworkController) StudentGroup(ctx context.Context, userID int64) (string, error) {
	return c.usersRepo.GetGroup(ctx, userID)
}

func (c *homeworkController) DaysWithLessons(ctx context.Context, group string) ([]string, error) {
	days, err := c.lessonsRepo.GetDaysWithLessonsByGroup(ctx, group)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return days, nil
}

func (c *homeworkController) LessonsByDay(ctx context.Context, group, day string) ([]domain.LessonBrief, error) {
	lessons, err := c.lessonsRepo.LessonsByDayGroup(ctx, group, day)
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (c *homeworkController) Pin(ctx context.Context, group string, day string, subject, text string) error {
	group = strings.ReplaceAll(strings.TrimSpace(group), " ", "")
	if err := c.homeworkRepo.SaveForGroup(ctx, group, fmt.Sprintf("%s(%s)", subject, day), fmt.Sprintf("• %s", text)); err != nil {
		return err
	}
	return nil
}

func (c *homeworkController) Update(ctx context.Context, homeworkID int64, newText string) error {
	if err := c.homeworkRepo.Update(ctx, homeworkID, newText); err != nil {
		return err
	}
	return nil
}

func (c *homeworkController) UpdateStatus(ctx context.Context, homeworkID int64) error {
	if err := c.homeworkRepo.UpdateStatus(ctx, homeworkID); err != nil {
		return err
	}
	return nil
}

func (c *homeworkController) DeleteHomework(ctx context.Context, homeworkID int64) error {
	if err := c.homeworkRepo.Delete(ctx, homeworkID); err != nil {
		return err
	}
	return nil
}

func (c *homeworkController) ListForLastWeek(ctx context.Context, userID int64) ([]domain.HWBrief, error) {
	homework, err := c.homeworkRepo.ListForUser(ctx, userID)
	if err != nil {
		return []domain.HWBrief{}, err
	}
	return homework, nil
}

func (c *homeworkController) ListAll(ctx context.Context) ([]domain.HWBrief, error) {
	homework, err := c.homeworkRepo.ListAll(ctx)
	if err != nil {
		return []domain.HWBrief{}, err
	}
	return homework, nil
}

func (c *homeworkController) CheckExistence(ctx context.Context, userID int64, homeworkID int64) (bool, error) {
	return c.homeworkRepo.CheckExistence(ctx, userID, homeworkID)
}
