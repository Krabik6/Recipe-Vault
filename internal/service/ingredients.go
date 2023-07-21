package service

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
)

// import (
//
//	"github.com/Krabik6/meal-schedule/internal/models"
//
// )
func (r *RecipesService) IngredientsInfo(ingredients []models.IngredientInput) ([]models.Ingredient, error) {
	processedIngredients := make([]models.Ingredient, len(ingredients))
	errs := make(chan error, len(ingredients))
	results := make(chan models.Ingredient, len(ingredients))

	// Запускаем горутины для получения информации об ингредиентах
	for i, ingredient := range ingredients {
		go func(i int, ingredient models.IngredientInput) {
			// Search for the ingredient
			searchResult, err := r.spoonacular.SearchIngredient(ingredient.Name)
			if err != nil {
				errs <- fmt.Errorf("failed to search ingredient: %w", err)
				return
			}

			// Check if the desired ingredient is found
			if len(searchResult.Results) == 0 {
				errs <- fmt.Errorf("ingredient not found: %s", ingredient.Name)
				return
			}

			// Get the first search result (assuming it's the desired ingredient)
			ingredientID := searchResult.Results[0].ID

			ingredientInfoOptions := &models.IngredientInfoOptions{
				Amount: 100,
				Unit:   "grams",
			}

			// Get the ingredient info
			info, err := r.spoonacular.GetIngredientInfo(ingredientID, ingredientInfoOptions)
			if err != nil {
				errs <- fmt.Errorf("failed to get ingredient info: %w", err)
				return
			}

			// Convert IngredientAPIResponse to Ingredient
			processedIngredient := models.Ingredient{
				ID:            info.ID,
				Name:          info.Name,
				Price:         info.EstimatedCost.Value,
				Unit:          info.Unit,
				UnitShort:     info.UnitShort,
				UnitLong:      info.UnitLong,
				PossibleUnits: info.PossibleUnits,
				Protein:       info.Nutrition.CaloricBreakdown.PercentProtein,
				Fat:           info.Nutrition.CaloricBreakdown.PercentFat,
				Carbs:         info.Nutrition.CaloricBreakdown.PercentCarbs,
				Aisle:         info.Aisle,
				Image:         info.Image,
				CategoryPath:  info.CategoryPath,
				Consistency:   info.Consistency,
				ExternalID:    info.ID,
				Amount:        info.Amount,
			}
			results <- processedIngredient
		}(i, ingredient)
	}

	// Собираем результаты
	for i := 0; i < len(ingredients); i++ {
		select {
		case err := <-errs:
			return nil, err
		case result := <-results:
			processedIngredients[i] = result
		}
	}

	return processedIngredients, nil
}

//
//func (r *RecipesService) AddOrUpdateIngredient(ingredient models.Ingredient) (int, error) {
//	// Check if the ingredient already exists in the database
//	existingIngredient, err := r.ingredients.GetIngredientByName(ingredient.Name)
//	if err != nil {
//		// If the ingredient does not exist, add it to the database
//		ingredientID, err := r.ingredients.AddIngredient(ingredient)
//		if err != nil {
//			return 0, err
//		}
//		return ingredientID, nil
//	} else {
//		// If the ingredient exists, use its ID
//		return existingIngredient.ID, nil
//	}
//}
//
//func (r *RecipesService) ProcessIngredients(ingredientsInfo []models.Ingredient) ([]models.Ingredient, error) {
//	processedIngredients := make([]models.Ingredient, len(ingredientsInfo))
//
//	for i, ingredient := range ingredientsInfo {
//		ingredientID, err := r.AddOrUpdateIngredient(ingredient)
//		if err != nil {
//			return nil, err
//		}
//		ingredient.ID = ingredientID
//		processedIngredients[i] = ingredient
//	}
//
//	return processedIngredients, nil
//}
