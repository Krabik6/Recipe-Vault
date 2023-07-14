package service

import (
	"context"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/internal/repository"
	"mime/multipart"
)

type RecipesService struct {
	repo     repository.Recipes
	uploader ImageUploader
}

func NewRecipesService(repo repository.Recipes, uploader ImageUploader) *RecipesService {
	return &RecipesService{repo: repo, uploader: uploader}
}

func (r *RecipesService) CreateRecipe(userId int, recipe models.Recipe, imageFiles []*multipart.FileHeader) (int, error) {
	imageURLs := make([]string, len(imageFiles))

	for i, file := range imageFiles {
		imageURL, err := r.uploader.UploadImage(context.TODO(), file)
		if err != nil {
			return 0, err
		}
		imageURLs[i] = imageURL
	}

	recipe.ImageURLs = imageURLs

	recipeID, err := r.repo.CreateRecipe(userId, recipe)
	if err != nil {
		return 0, err
	}

	return recipeID, nil
}

func (r *RecipesService) GetFilteredUserRecipes(userId int, input models.RecipesFilter) ([]models.Recipe, error) {
	return r.repo.GetFilteredUserRecipes(userId, input)
}

func (r *RecipesService) GetFilteredRecipes(input models.RecipesFilter) ([]models.Recipe, error) {
	return r.repo.GetFilteredRecipes(input)
}

func (r *RecipesService) GetRecipeById(userId, id int) (models.Recipe, error) {
	return r.repo.GetRecipeById(userId, id)
}
func (r *RecipesService) GetAllRecipes(userId int) ([]models.Recipe, error) {
	return r.repo.GetAllRecipes(userId)
}

func (r *RecipesService) GetPublicRecipes() ([]models.Recipe, error) {
	return r.repo.GetPublicRecipes()
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
	return r.repo.UpdateRecipe(userId, id, input)
}
func (r *RecipesService) DeleteRecipe(userId, id int) error {

	return r.repo.DeleteRecipe(userId, id)
}
