package repository

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
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
	GetAllSchedule(userId int) ([]models.ScheduleOutput, error)
	GetScheduleByPeriod(userId int, date string, dayPeriod int) ([]models.ScheduleOutput, error)
	UpdateSchedule(userId int, date string, input models.UpdateScheduleInput) error
	DeleteSchedule(userId int, date string) error
	CreateMeal(userId int, meal models.Meal) (int, error)
}

type Repository struct {
	Authorization
	Recipes
	Schedule
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Recipes:       NewRecipesPostgres(db),
		Schedule:      NewSchedulePostgres(db),
	}
}
