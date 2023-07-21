package repository

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type SpoonacularAPI interface {
	SearchIngredient(Query string) (models.IngredientSearchResponse, error)
	GetIngredientInfo(id int, options *models.IngredientInfoOptions) (*models.IngredientAPIResponse, error)
	ExtractIngredient(text string) (*models.ExtractedIngredient, error)
	ConvertAmounts(ingredientName string, sourceAmount float64, sourceUnit string, targetUnit string) (*models.ConversionResult, error)
}

type Recipes interface {
	CreateRecipe(userId int, recipe models.Recipe, ingredients []models.Ingredient) (int, error)
	GetRecipeById(userId, id int) (models.RecipeOutput, error)
	GetAllRecipes(userId int) ([]models.RecipeOutput, error)
	UpdateRecipe(userId, id int, input models.UpdateRecipe) error
	DeleteRecipe(userId, id int) error
	GetPublicRecipes() ([]models.Recipe, error)
	GetFilteredRecipes(input models.RecipesFilter) ([]models.Recipe, error)
	GetFilteredUserRecipes(userId int, input models.RecipesFilter) ([]models.Recipe, error)
}

type Ingredients interface {
	GetIngredientByName(name string) (models.Ingredient, error)
	AddIngredient(ingredient models.Ingredient) (int, error)
	AddRecipeIngredient(recipeId int, ingredient models.Ingredient, amount float64) error
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
	SpoonacularAPI
	Ingredients
}

func NewRepository(db *sqlx.DB, SpoonacularAPI SpoonacularAPI) *Repository {
	return &Repository{
		Authorization:  NewAuthPostgres(db),
		Recipes:        NewRecipesPostgres(db),
		Schedule:       NewSchedulePostgres(db),
		SpoonacularAPI: SpoonacularAPI,
	}
}
