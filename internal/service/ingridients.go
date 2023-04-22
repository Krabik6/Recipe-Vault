package service

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/internal/repository"
)

type IngredientsService struct {
	repo repository.Ingredient
}

func NewIngredientsService(repo repository.Ingredient) *IngredientsService {
	return &IngredientsService{repo: repo}
}

func (r *IngredientsService) CreateIngredient(userId int, ingredient models.Ingredient) (int, error) {
	return r.repo.CreateIngredient(userId, ingredient)
}
func (r *IngredientsService) GetIngredientById(userId, id int) (models.IngredientOutput, error) {
	return r.repo.GetIngredientById(userId, id)
}
func (r *IngredientsService) GetAllIngredients(userId int) ([]models.Ingredient, error) {
	return r.repo.GetAllIngredients(userId)
}

func (r *IngredientsService) GetPublicIngredients() ([]models.Ingredient, error) {
	return r.repo.GetPublicIngredients()
}

func (r *IngredientsService) UpdateIngredient(userId, id int, input models.UpdateIngredientInput) error {
	return r.repo.UpdateIngredient(userId, id, input)
}
func (r *IngredientsService) DeleteIngredient(userId, id int) error {
	return r.repo.DeleteIngredient(userId, id)
}
