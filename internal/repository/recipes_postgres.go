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
	addRecipeQuery := fmt.Sprintf("INSERT INTO recipes (name, description) values ($1, $2)")
	_, err := r.db.Exec(addRecipeQuery, recipe.Name, recipe.Description)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (r *RecipesPostgres) GetRecipeById(userId, id int) (models.Recipe, error) {
	output := models.Recipe{}
	getRecipeByIdQuery := fmt.Sprintf("SELECT * FROM recipes WHERE id=$1")
	err := r.db.Get(&output, getRecipeByIdQuery, id)
	if err != nil {
		return output, err
	}
	fmt.Println(output)

	return output, err
}

func (r *RecipesPostgres) GetAllRecipes(userId int) ([]models.Recipe, error) {
	output := []models.Recipe{}
	getRecipeByIdQuery := fmt.Sprintf("SELECT * FROM recipes")
	err := r.db.Select(&output, getRecipeByIdQuery)
	if err != nil {
		return output, err
	}
	fmt.Println(output)

	return output, err
}
func (r *RecipesPostgres) UpdateRecipe(userId, id int, input models.UpdateRecipeInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE recipes SET %s WHERE id=%d", setQuery, id)
	args = append(args)

	_, err := r.db.Exec(query, args...)
	return err
}
func (r *RecipesPostgres) DeleteRecipe(userId, id int) error {
	query := fmt.Sprintf("DELETE FROM recipes WHERE id=$1")
	_, err := r.db.Exec(query, id)
	return err

}
