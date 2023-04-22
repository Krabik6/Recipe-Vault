package service

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/internal/repository"
)

type ScheduleService struct {
	repo repository.Schedule
}

func NewScheduleService(repo repository.Schedule) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) FillSchedule(userId int, meal models.Meal) (int, error) {
	return s.repo.FillSchedule(userId, meal)
}

func (s *ScheduleService) GetAllSchedule(userId int) ([]models.ScheduleOutput, error) {
	return s.repo.GetAllSchedule(userId)
}
func (s *ScheduleService) GetScheduleByDate(userId int, date string) ([]models.ScheduleOutput, error) {
	return s.repo.GetScheduleByDate(userId, date)

}
func (s *ScheduleService) UpdateSchedule(userId int, date string, input models.UpdateScheduleInput) error {
	return s.repo.UpdateSchedule(userId, date, input)

}
func (s *ScheduleService) DeleteSchedule(userId int, date string) error {
	return s.repo.DeleteSchedule(userId, date)
}

func (s *ScheduleService) CreateMeal(userId int, meal models.Meal) (int, error) {
	return s.repo.CreateMeal(userId, meal)
}
