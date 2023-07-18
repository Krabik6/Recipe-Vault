package service

import (
	"context"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/internal/repository"
	"mime/multipart"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Recipes interface {
	CreateRecipe(userId int, recipe models.Recipe, imageFiles []*multipart.FileHeader) (int, error)
	GetRecipeById(userId, id int) (models.Recipe, error)
	GetAllRecipes(userId int) ([]models.Recipe, error)
	UpdateRecipe(userId, id int, input models.UpdateRecipeInput, imageFiles []*multipart.FileHeader) error
	DeleteRecipe(userId, id int) error
	GetPublicRecipes() ([]models.Recipe, error)
	GetFilteredRecipes(input models.RecipesFilter) ([]models.Recipe, error)
	GetFilteredUserRecipes(userId int, input models.RecipesFilter) ([]models.Recipe, error)
}

type Schedule interface {
	GetAllSchedule(userId int) ([]models.ScheduleOutput, error)
	GetScheduleByPeriod(userId int, date string, dayPeriod int) ([]models.ScheduleOutput, error)
	UpdateSchedule(userId int, date string, input models.UpdateScheduleInput) error
	DeleteSchedule(userId int, date string) error
	CreateMeal(userId int, meal models.Meal) (int, error)
}

type ImageUploader interface {
	UploadImage(ctx context.Context, imageFile *multipart.FileHeader) (string, error)
}

type Service struct {
	Authorization
	Recipes
	Schedule
}

//repos *Repository.Repository

func NewService(repos *repository.Repository, uploader ImageUploader) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Recipes:       NewRecipesService(repos.Recipes, uploader),
		Schedule:      NewScheduleService(repos.Schedule),
	}
}
