package controllers

import (
	"context"
	"fmt"
	"log"
	"time"
	"timetable-homework-tgbot/internal/repositories"

	"timetable-homework-tgbot/internal/domain"
)

type HomeworkController interface {
	DaysWithLessons(ctx context.Context, userID int64) ([]string, error)
	LessonsByDay(ctx context.Context, userID int64, day string) ([]domain.LessonBrief, error)
	Pin(ctx context.Context, userID int64, day string, subject, text string) error
	Update(ctx context.Context, userID int64, subject, newText string) error
	UpdateStatus(ctx context.Context, userID int64, subject string) error
	DeleteHomework(ctx context.Context, userID int64, subject string) error
	ListForLastWeek(ctx context.Context, userID int64) ([]domain.HWBrief, error)
	CheckExistence(ctx context.Context, userID int64, subject string) (bool, error)
}

type hw struct {
	timeout time.Duration

	usersRepo    repositories.UsersRepository
	homeworkRepo repositories.HomeworkRepository
	lessonsRepo  repositories.LessonsRepository
}

func NewHomeworkController(usersRepo repositories.UsersRepository, homeworkRepo repositories.HomeworkRepository,
	lessonsRepo repositories.LessonsRepository) HomeworkController {
	return &hw{
		timeout:      5 * time.Second,
		usersRepo:    usersRepo,
		homeworkRepo: homeworkRepo,
		lessonsRepo:  lessonsRepo,
	}
}

func (c *hw) DaysWithLessons(ctx context.Context, userID int64) ([]string, error) {
	group, err := c.usersRepo.GetGroup(ctx, userID)
	if err != nil {
		return nil, err
	}

	days, err := c.lessonsRepo.GetDaysWithLessonsByGroup(ctx, group)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return days, nil
}

func (c *hw) LessonsByDay(ctx context.Context, userID int64, day string) ([]domain.LessonBrief, error) {
	group, err := c.usersRepo.GetGroup(ctx, userID)
	if err != nil {
		return nil, err
	}

	lessons, err := c.lessonsRepo.LessonsByDayGroup(ctx, group, day)
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (c *hw) Pin(ctx context.Context, userID int64, day string, subject, text string) error {
	if err := c.homeworkRepo.Save(ctx, userID, fmt.Sprintf("%s(%s)", subject, day), fmt.Sprintf("• %s", text)); err != nil {
		return err
	}
	return nil
}

func (c *hw) Update(ctx context.Context, userID int64, subject, newText string) error {
	if err := c.homeworkRepo.Update(ctx, userID, subject, newText); err != nil {
		return err
	}
	return nil
}

func (c *hw) UpdateStatus(ctx context.Context, userID int64, subject string) error {
	if err := c.homeworkRepo.UpdateStatus(ctx, userID, subject); err != nil {
		return err
	}
	return nil
}

func (c *hw) DeleteHomework(ctx context.Context, userID int64, subject string) error {
	if err := c.homeworkRepo.Delete(ctx, userID, subject); err != nil {
		return err
	}
	return nil
}

func (c *hw) ListForLastWeek(ctx context.Context, userID int64) ([]domain.HWBrief, error) {
	homework, err := c.homeworkRepo.ListForLastWeek(ctx, userID)
	if err != nil {
		return []domain.HWBrief{}, err
	}
	return homework, nil
}

func (c *hw) CheckExistence(ctx context.Context, userID int64, subject string) (bool, error) {
	return c.homeworkRepo.CheckExistence(ctx, userID, subject)
}
