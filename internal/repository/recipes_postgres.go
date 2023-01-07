package repository

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strings"
)

type RecipesPostgres struct {
	db *sqlx.DB
}

func NewRecipesPostgres(db *sqlx.DB) *RecipesPostgres {
	return &RecipesPostgres{db: db}
}

func (r *RecipesPostgres) CreateRecipe(userId int, recipe models.Recipe) (int, error) {
	db := r.db

	addRecipeQuery := fmt.Sprintf(`INSERT INTO %s (title, description, "userId") values ($1, $2, $3) RETURNING id`, recipeTable)

	row := db.QueryRow(addRecipeQuery, recipe.Title, recipe.Description, userId)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RecipesPostgres) GetRecipeById(userId, id int) (models.Recipe, error) {
	db := r.db

	output := models.Recipe{}
	getRecipeByIdQuery := fmt.Sprintf(`SELECT rt.id, rt.title, rt.description FROM %s as rt WHERE rt."userId" = $1 and rt.id=$2`, recipeTable)
	err := db.Get(&output, getRecipeByIdQuery, userId, id)
	if err != nil {
		return output, err
	}
	//todo
	return output, err
}

func (r *RecipesPostgres) GetAllRecipes(userId int) ([]models.Recipe, error) {
	db := r.db

	output := []models.Recipe{}
	getAllRecipeQuery := fmt.Sprintf(`SELECT rt.id, rt.title, rt.description FROM as %s rt  WHERE rt."userId" = $1`, recipeTable)
	err := db.Select(&output, getAllRecipeQuery, userId)
	if err != nil {
		return output, err
	}

	return output, err
}

//func (r *RecipesPostgres) GetAllExistRecipes() ([]models.Recipe, error) {
//	db := r.db
//
//	output := []models.Recipe{}
//	getAllExistRecipeQuery := fmt.Sprintf(`SELECT rt.id, rt.title, rt.description FROM %s`, recipeTable)
//	err := db.Select(&output, getAllExistRecipeQuery)
//	if err != nil {
//		return output, err
//	}
//
//	return output, err
//} todo interfaces etc

func (r *RecipesPostgres) UpdateRecipe(userId, id int, input models.UpdateRecipeInput) error {
	db := r.db

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

	query := fmt.Sprintf(`UPDATE %s rt SET %s FROM %s ut WHERE rt.id = ut.id AND rt.id=%d AND ut."userId"=%s`, recipeTable, setQuery, userTable, id, userId)
	args = append(args)

	_, err := db.Exec(query, args...)

	return err
}

func (r *RecipesPostgres) DeleteRecipe(userId, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", recipeTable)
	_, err := r.db.Exec(query, id)
	return err

}
