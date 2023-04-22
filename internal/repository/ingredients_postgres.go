package repository

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strings"
)

type IngredientsPostgres struct {
	db *sqlx.DB
}

func NewIngredientsPostgres(db *sqlx.DB) *IngredientsPostgres {
	return &IngredientsPostgres{db: db}
}

func (r *IngredientsPostgres) CreateIngredient(userId int, ingredient models.Ingredient) (int, error) {
	db := r.db

	addIngredientQuery := fmt.Sprintf(`INSERT INTO %s (title, description, user_id, public) values ($1, $2, $3, $4) RETURNING id`, ingredientsTable)

	row := db.QueryRow(addIngredientQuery, ingredient.Title, ingredient.Description, userId, ingredient.IsPublic)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *IngredientsPostgres) GetIngredientById(userId, id int) (models.IngredientOutput, error) {
	db := r.db

	output := models.IngredientOutput{}
	getIngredientByIdQuery := fmt.Sprintf(`SELECT rt.title, rt.description, rt.public FROM %s as rt WHERE rt."user_id" = $1 and rt.id=$2`, ingredientsTable)
	err := db.Get(&output, getIngredientByIdQuery, userId, id)
	if err != nil {
		return output, err
	}
	//todo
	return output, err
}

func (r *IngredientsPostgres) GetAllIngredients(userId int) ([]models.Ingredient, error) {
	db := r.db

	var output []models.Ingredient
	getAllRecipeQuery := fmt.Sprintf(`SELECT rt.id, rt.title, rt.description, rt.public FROM  %s as rt  WHERE rt."user_id" = $1`, ingredientsTable)
	err := db.Select(&output, getAllRecipeQuery, userId)
	if err != nil {
		return output, err
	}

	return output, err
}

func (r *IngredientsPostgres) GetPublicIngredients() ([]models.Ingredient, error) {
	db := r.db

	var output []models.Ingredient
	getAllExistIngredientQuery := fmt.Sprintf(`SELECT rt.id, rt.title, rt.description, rt.public FROM %s as rt WHERE rt.public=true`, ingredientsTable)
	err := db.Select(&output, getAllExistIngredientQuery)
	if err != nil {
		return output, err
	}

	return output, err
}

func (r *IngredientsPostgres) UpdateIngredient(userId, id int, input models.UpdateIngredientInput) error {
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

	if input.IsPublic != nil {
		setValues = append(setValues, fmt.Sprintf("public=$%d", argId))
		args = append(args, *input.IsPublic)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE "id"=%d AND "user_id"=%d`, ingredientsTable, setQuery, id, userId)
	args = append(args)

	_, err := db.Exec(query, args...)

	return err
}

func (r *IngredientsPostgres) DeleteIngredient(userId, id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE "id"=$1 and "user_id"=$2 `, ingredientsTable)
	_, err := r.db.Exec(query, id, userId)
	return err
}
