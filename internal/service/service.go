package service

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Recipes interface {
	CreateRecipe(userId int, recipe models.Recipe) (int, error)
	GetRecipeById(userId, id int) (models.RecipeOutput, error)
	GetAllRecipes(userId int) ([]models.Recipe, error)
	UpdateRecipe(userId, id int, input models.UpdateRecipeInput) error
	DeleteRecipe(userId, id int) error
	GetPublicRecipes() ([]models.Recipe, error)
}

type Schedule interface {
	FillSchedule(userId int, schedule models.Schedule) (int, error)
	GetAllSchedule(userId int) ([]models.ScheduleOutput, error)
	GetScheduleByDate(userId int, date string) (models.ScheduleOutput, error)
	UpdateSchedule(userId int, date string, input models.UpdateScheduleInput) error
	DeleteSchedule(userId int, date string) error
}

type Service struct {
	Authorization
	Recipes
	Schedule
}

//repos *Repository.Repository

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Recipes:       NewRecipesService(repos.Recipes),
		Schedule:      NewScheduleService(repos.Schedule),
	}
}
