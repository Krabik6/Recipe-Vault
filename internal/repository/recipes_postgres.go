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

func (r *RecipesPostgres) CreateRecipe(userId int, recipe models.Recipe) (int, error) {
	db := r.db

	addRecipeQuery := fmt.Sprintf(`INSERT INTO %s 
		(title, description, user_id, public, "cost", "timeToPrepare", "healthy")
		values ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`,
		recipeTable)

	row := db.QueryRow(addRecipeQuery, recipe.Title, recipe.Description, userId, recipe.IsPublic, recipe.Cost, recipe.TimeToPrepare, recipe.Healthy)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RecipesPostgres) GetRecipeById(userId, id int) (models.Recipe, error) {
	db := r.db

	output := models.Recipe{}
	getRecipeByIdQuery := fmt.Sprintf(`SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy" FROM %s as rt WHERE rt."user_id" = $1 and rt.id=$2`, recipeTable)
	err := db.Get(&output, getRecipeByIdQuery, userId, id)
	if err != nil {
		return output, err
	}
	//todo
	return output, err
}

func (r *RecipesPostgres) GetFilteredRecipes(input models.RecipesFilter) ([]models.Recipe, error) {
	db := r.db
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	var output []models.Recipe

	if input.CostMoreThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"cost" > $%d`, argId))
		args = append(args, *input.CostMoreThan)
		argId++
	}

	if input.CostLessThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"cost" < $%d`, argId))
		args = append(args, *input.CostLessThan)
		argId++
	}

	if input.TimeToPrepareMoreThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"timeToPrepare" > $%d`, argId))
		args = append(args, *input.TimeToPrepareMoreThan)
		argId++
	}

	if input.TimeToPrepareLessThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"timeToPrepare" < $%d`, argId))
		args = append(args, *input.TimeToPrepareLessThan)
		argId++
	}

	if input.HealthyMoreThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"healthy" > $%d`, argId))
		args = append(args, *input.HealthyMoreThan)
		argId++
	}

	if input.HealthyLessThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"healthy" < $%d`, argId))
		args = append(args, *input.HealthyLessThan)
		argId++
	}

	setQuery := strings.Join(setValues, " and ")
	if len(setQuery) > 0 {
		setQuery = "and " + setQuery
	}
	log.Println(setQuery, 100)
	log.Println(args...)

	query := fmt.Sprintf(`SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy" FROM  %s as rt  WHERE rt.public=true %s`, recipeTable, setQuery)
	args = append(args)

	err := db.Select(&output, query, args...)
	if err != nil {
		return nil, err
	}

	return output, err
}

// func GetFilteredUserRecipes
func (r *RecipesPostgres) GetFilteredUserRecipes(userId int, input models.RecipesFilter) ([]models.Recipe, error) {
	db := r.db
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	var output []models.Recipe

	if input.CostMoreThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"cost" > $%d`, argId))
		args = append(args, *input.CostMoreThan)
		argId++
	}

	if input.CostLessThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"cost" < $%d`, argId))
		args = append(args, *input.CostLessThan)
		argId++
	}

	if input.TimeToPrepareMoreThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"timeToPrepare" > $%d`, argId))
		args = append(args, *input.TimeToPrepareMoreThan)
		argId++
	}

	if input.TimeToPrepareLessThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"timeToPrepare" < $%d`, argId))
		args = append(args, *input.TimeToPrepareLessThan)
		argId++
	}

	if input.HealthyMoreThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"healthy" > $%d`, argId))
		args = append(args, *input.HealthyMoreThan)
		argId++
	}

	if input.HealthyLessThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"healthy" < $%d`, argId))
		args = append(args, *input.HealthyLessThan)
		argId++
	}

	setQuery := strings.Join(setValues, " and ")
	if len(setQuery) > 0 {
		setQuery = "and " + setQuery
	}
	log.Println(setQuery)
	log.Println(args...)

	query := fmt.Sprintf(`SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy" FROM  %s as rt  WHERE rt.public=true %s and rt."user_id"=%d`, recipeTable, setQuery, userId)
	args = append(args)

	err := db.Select(&output, query, args...)
	if err != nil {
		return nil, err
	}

	return output, err

}

func (r *RecipesPostgres) GetAllRecipes(userId int) ([]models.Recipe, error) {
	db := r.db

	var output []models.Recipe
	getAllRecipeQuery := fmt.Sprintf(`SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy" FROM  %s as rt  WHERE rt."user_id" = $1`, recipeTable)
	err := db.Select(&output, getAllRecipeQuery, userId)
	if err != nil {
		return output, err
	}

	return output, err
}

func (r *RecipesPostgres) GetPublicRecipes() ([]models.Recipe, error) {
	db := r.db

	var output []models.Recipe
	getAllExistRecipeQuery := fmt.Sprintf(`SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy" FROM %s as rt WHERE rt.public=true`, recipeTable)
	err := db.Select(&output, getAllExistRecipeQuery)
	if err != nil {
		return output, err
	}

	return output, err
}

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

	if input.IsPublic != nil {
		setValues = append(setValues, fmt.Sprintf("public=$%d", argId))
		args = append(args, *input.IsPublic)
		argId++
	}

	if input.Cost != nil {
		setValues = append(setValues, fmt.Sprintf(`"cost"=$%d`, argId))
		args = append(args, *input.Cost)
		argId++
	}

	if input.TimeToPrepare != nil {
		setValues = append(setValues, fmt.Sprintf(`"timeToPrepare"=$%d`, argId))
		args = append(args, *input.TimeToPrepare)
		argId++
	}

	if input.Healthy != nil {
		setValues = append(setValues, fmt.Sprintf(`"healthy"=$%d`, argId))
		args = append(args, *input.Healthy)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE "id"=%d AND "user_id"=%d`, recipeTable, setQuery, id, userId)
	fmt.Println(query)
	args = append(args)

	_, err := db.Exec(query, args...)

	return err
}

func (r *RecipesPostgres) DeleteRecipe(userId, id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE "id"=$1 and "user_id"=$2 `, recipeTable)
	_, err := r.db.Exec(query, id, userId)
	return err

}
