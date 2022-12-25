package repository

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type RecipesPostgres struct {
	db *sqlx.DB
}

func NewRecipesPostgres(db *sqlx.DB) *RecipesPostgres {
	return &RecipesPostgres{db: db}
}

func (r *RecipesPostgres) CreateRecipe(userId int, recipe models.Recipe) error {
	addRecipeQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2)", recipeTable)
	log.Println(addRecipeQuery)

	_, err := r.db.Exec(addRecipeQuery, recipe.Title, recipe.Description)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

func (r *RecipesPostgres) GetRecipeById(userId, id int) (models.Recipe, error) {
	output := models.Recipe{}
	getRecipeByIdQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", recipeTable)
	err := r.db.Get(&output, getRecipeByIdQuery, id)
	if err != nil {
		return output, err
	}

	return output, err
}

func (r *RecipesPostgres) GetAllRecipes(userId int) ([]models.Recipe, error) {
	output := []models.Recipe{}
	getRecipeByIdQuery := fmt.Sprintf("SELECT * FROM %s", recipeTable)
	err := r.db.Select(&output, getRecipeByIdQuery)
	if err != nil {
		return output, err
	}

	return output, err
}

func (r *RecipesPostgres) UpdateRecipe(userId, id int, input models.UpdateRecipeInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%d", recipeTable, setQuery, id)
	args = append(args)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *RecipesPostgres) DeleteRecipe(userId, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", recipeTable)
	_, err := r.db.Exec(query, id)
	return err

}
