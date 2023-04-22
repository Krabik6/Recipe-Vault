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
	GetRecipeById(userId, id int) (models.Recipe, error)
	GetAllRecipes(userId int) ([]models.Recipe, error)
	UpdateRecipe(userId, id int, input models.UpdateRecipeInput) error
	DeleteRecipe(userId, id int) error
	GetPublicRecipes() ([]models.Recipe, error)
	GetFilteredRecipes(input models.RecipesFilter) ([]models.Recipe, error)
	GetFilteredUserRecipes(userId int, input models.RecipesFilter) ([]models.Recipe, error)
}

type Schedule interface {
	FillSchedule(userId int, meal models.Meal) (int, error)
	GetAllSchedule(userId int) ([]models.ScheduleOutput, error)
	GetScheduleByDate(userId int, date string) ([]models.ScheduleOutput, error)
	UpdateSchedule(userId int, date string, input models.UpdateScheduleInput) error
	DeleteSchedule(userId int, date string) error
	CreateMeal(userId int, meal models.Meal) (int, error)
}

type Ingredients interface {
	CreateIngredient(userId int, ingredient models.Ingredient) (int, error)
	GetIngredientById(userId, id int) (models.IngredientOutput, error)
	GetAllIngredients(userId int) ([]models.Ingredient, error)
	GetPublicIngredients() ([]models.Ingredient, error)
	UpdateIngredient(userId, id int, input models.UpdateIngredientInput) error
	DeleteIngredient(userId, id int) error
}

type Service struct {
	Authorization
	Recipes
	Schedule
	Ingredients
}

//repos *Repository.Repository

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Recipes:       NewRecipesService(repos.Recipes),
		Schedule:      NewScheduleService(repos.Schedule),
		Ingredients:   NewIngredientsService(repos.Ingredient),
	}
}
