package spoonacular

import "github.com/Krabik6/meal-schedule/internal/models"

type SpoonacularAPI struct {
	BaseURL string
	APIKey  string
	Options *models.IngredientSearchOptions
}

func NewSpoonacularAPI(
	baseURL string,
	APIKey string,
	Options *models.IngredientSearchOptions,
) *SpoonacularAPI {
	return &SpoonacularAPI{
		BaseURL: baseURL,
		APIKey:  APIKey,
		Options: Options,
	}
}
