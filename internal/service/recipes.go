package service

import (
	"context"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/internal/repository"
	"mime/multipart"
)

type RecipesService struct {
	recipes     repository.Recipes
	spoonacular repository.SpoonacularAPI
	uploader    ImageUploader
}

func NewRecipesService(repo repository.Recipes, spoon repository.SpoonacularAPI, ingredients repository.Ingredients, uploader ImageUploader) *RecipesService {
	return &RecipesService{
		recipes:     repo,
		spoonacular: spoon,
		uploader:    uploader,
	}
}

func (r *RecipesService) CreateRecipe(userId int, recipeInput models.RecipeInput, imageFiles []*multipart.FileHeader) (int, error) {
	imageURLs := make([]string, len(imageFiles))
	errs := make(chan error, len(imageFiles))
	results := make(chan string, len(imageFiles))

	// Запускаем горутины для загрузки изображений
	for i, file := range imageFiles {
		go func(i int, file *multipart.FileHeader) {
			imageURL, err := r.uploader.UploadImage(context.TODO(), file)
			if err != nil {
				errs <- err
				return
			}
			results <- imageURL
		}(i, file)
	}

	// Собираем результаты
	for i := 0; i < len(imageFiles); i++ {
		select {
		case err := <-errs:
			return 0, err
		case result := <-results:
			imageURLs[i] = result
		}
	}

	recipeInput.ImageURLs = imageURLs
	ingredientsInfo, err := r.IngredientsInfo(recipeInput.IngredientInputs)
	if err != nil {
		return 0, err
	}

	recipe := models.Recipe{
		Id:               recipeInput.Id,
		Title:            recipeInput.Title,
		Description:      recipeInput.Description,
		IsPublic:         recipeInput.IsPublic,
		Cost:             recipeInput.Cost,
		TimeToPrepare:    recipeInput.TimeToPrepare,
		Healthy:          recipeInput.Healthy,
		ImageURLs:        recipeInput.ImageURLs,
		IngredientInputs: recipeInput.IngredientInputs,
	}
	recipe.Ingredients = ingredientsInfo

	recipeID, err := r.recipes.CreateRecipe(userId, recipe, ingredientsInfo)
	if err != nil {
		return 0, err
	}

	return recipeID, nil
}

func (r *RecipesService) GetFilteredUserRecipes(userId int, input models.RecipesFilter) ([]models.Recipe, error) {
	return r.recipes.GetFilteredUserRecipes(userId, input)
}

func (r *RecipesService) GetFilteredRecipes(input models.RecipesFilter) ([]models.Recipe, error) {
	return r.recipes.GetFilteredRecipes(input)
}
func (r *RecipesService) GetRecipeById(userId, id int) (models.RecipeOutput, error) {
	return r.recipes.GetRecipeById(userId, id)
}
func (r *RecipesService) GetAllRecipes(userId int) ([]models.RecipeOutput, error) {
	return r.recipes.GetAllRecipes(userId)
}

func (r *RecipesService) GetPublicRecipes() ([]models.Recipe, error) {
	return r.recipes.GetPublicRecipes()
}

func (r *RecipesService) UpdateRecipe(userId, id int, input models.UpdateRecipeInput, imageFiles []*multipart.FileHeader) error {
	imageURLs := make([]string, len(imageFiles))

	for i, file := range imageFiles {
		imageURL, err := r.uploader.UploadImage(context.TODO(), file)
		if err != nil {
			return err
		}
		imageURLs[i] = imageURL
	}

	input.ImageURLs = &imageURLs

	recipe := models.UpdateRecipe{
		Id:               input.Id,
		Title:            input.Title,
		Description:      input.Description,
		IsPublic:         input.IsPublic,
		Cost:             input.Cost,
		TimeToPrepare:    input.TimeToPrepare,
		Healthy:          input.Healthy,
		ImageURLs:        input.ImageURLs,
		IngredientInputs: input.IngredientInputs,
	}

	if input.IngredientInputs != nil {
		ingredientsInfo, err := r.IngredientsInfo(*input.IngredientInputs)
		if err != nil {
			return err
		}
		recipe.Ingredients = ingredientsInfo

	}

	return r.recipes.UpdateRecipe(userId, id, recipe)
}
func (r *RecipesService) DeleteRecipe(userId, id int) error {

	return r.recipes.DeleteRecipe(userId, id)
}
