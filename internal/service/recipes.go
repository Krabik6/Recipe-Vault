package service

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/internal/repository"
)

type RecipesService struct {
	repo repository.Recipes
}

func NewRecipesService(repo repository.Recipes) *RecipesService {
	return &RecipesService{repo: repo}
}

func (r *RecipesService) CreateRecipe(userId int, recipe models.Recipe) error {
	return r.repo.CreateRecipe(userId, recipe)
}
func (r *RecipesService) GetRecipeById(userId, id int) (models.Recipe, error) {
	return r.repo.GetRecipeById(userId, id)
}
func (r *RecipesService) GetAllRecipes(userId int) ([]models.Recipe, error) {
	return r.repo.GetAllRecipes(userId)
}
func (r *RecipesService) UpdateRecipe(userId, id int, input models.UpdateRecipeInput) error {
	return r.repo.UpdateRecipe(userId, id, input)
}
func (r *RecipesService) DeleteRecipe(userId, id int) error {
	return r.repo.DeleteRecipe(userId, id)
}
