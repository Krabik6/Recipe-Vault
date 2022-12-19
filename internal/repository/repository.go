package repository

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
)

type Recipes interface {
	CreateRecipe(userId int, recipe models.Recipe) error
	GetRecipeById(userId, id int) (models.Recipe, error)
	GetAllRecipes(userId int) ([]models.Recipe, error)
	UpdateRecipe(userId, id int, input models.UpdateRecipeInput) error
	DeleteRecipe(userId, id int) error
}

type Repository struct {
	Recipes
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Recipes: NewRecipesPostgres(db),
	}
}
