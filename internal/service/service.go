package service

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/internal/repository"
)

type Authorization interface {
}

type Recipes interface {
	CreateRecipe(userId int, recipe models.Recipe) error
	GetRecipeById(userId, id int) (models.Recipe, error)
	GetAllRecipes(userId int) ([]models.Recipe, error)
	UpdateRecipe(userId, id int, input models.UpdateRecipeInput) error
	DeleteRecipe(userId, id int) error
}

type Schedule interface {
	FillSchedule(userId int, schedule models.Schedule) (int, error)
	GetAllSchedule(userId int) ([]models.ScheduleOutput, error)
	GetScheduleByDate(userId int, date string) (models.ScheduleOutput, error)
	UpdateSchedule(userId int, date string, input models.UpdateScheduleInput) error
	DeleteSchedule(userId int, date string) error
}

type Service struct {
	//Authorization
	Recipes
	Schedule
}

//repos *Repository.Repository

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Recipes:  NewRecipesService(repos.Recipes),
		Schedule: NewScheduleService(repos.Schedule),
	}
}
