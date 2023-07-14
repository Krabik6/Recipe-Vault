package repository

import (
	"errors"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

	addRecipeQuery := fmt.Sprintf(`INSERT INTO %s 
		(title, description, user_id, public, "cost", "timeToPrepare", "healthy", "imageURLs")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id`,
		recipeTable)

	row := db.QueryRow(addRecipeQuery, recipe.Title, recipe.Description, userId, recipe.IsPublic, recipe.Cost, recipe.TimeToPrepare, recipe.Healthy, pq.Array(recipe.ImageURLs))

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RecipesPostgres) GetRecipeById(userId, id int) (models.Recipe, error) {
	db := r.db

	output := models.Recipe{}
	getRecipeByIdQuery := fmt.Sprintf(`
		SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy", rt."imageURLs"
		FROM %s AS rt
		WHERE rt."user_id" = $1 AND rt.id = $2`, recipeTable)

	row := db.QueryRowx(getRecipeByIdQuery, userId, id)
	if err := row.Err(); err != nil {
		return output, err
	}

	err := row.Scan(
		&output.Id,
		&output.Title,
		&output.Description,
		&output.IsPublic,
		&output.Cost,
		&output.TimeToPrepare,
		&output.Healthy,
		pq.Array(&output.ImageURLs),
	)
	if err != nil {
		return output, err
	}

	return output, nil
}
func (r *RecipesPostgres) GetFilteredRecipes(input models.RecipesFilter) ([]models.Recipe, error) {
	db := r.db

	var output []models.Recipe

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

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
		setQuery = "AND " + setQuery
	}

	query := fmt.Sprintf(`
		SELECT
			rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy", rt."imageURLs"
		FROM
			%s AS rt
		WHERE
			rt.public = true %s
	`, recipeTable, setQuery)

	rows, err := db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(
			&recipe.Id,
			&recipe.Title,
			&recipe.Description,
			&recipe.IsPublic,
			&recipe.Cost,
			&recipe.TimeToPrepare,
			&recipe.Healthy,
			pq.Array(&recipe.ImageURLs),
		)
		if err != nil {
			return nil, err
		}
		output = append(output, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
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

	query := fmt.Sprintf(`
		SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy", rt."imageURLs"
		FROM %s AS rt
		WHERE rt."public" = true AND rt."user_id" = $%d %s`, recipeTable, argId, setQuery)

	args = append(args, userId)

	rows, err := db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(
			&recipe.Id,
			&recipe.Title,
			&recipe.Description,
			&recipe.IsPublic,
			&recipe.Cost,
			&recipe.TimeToPrepare,
			&recipe.Healthy,
			pq.Array(&recipe.ImageURLs),
		)
		if err != nil {
			return nil, err
		}
		output = append(output, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func (r *RecipesPostgres) GetAllRecipes(userId int) ([]models.Recipe, error) {
	db := r.db

	var output []models.Recipe
	getAllRecipeQuery := fmt.Sprintf(`
		SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy", rt."imageURLs"
		FROM %s AS rt
		WHERE rt."user_id" = $1`, recipeTable)

	rows, err := db.Queryx(getAllRecipeQuery, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(
			&recipe.Id,
			&recipe.Title,
			&recipe.Description,
			&recipe.IsPublic,
			&recipe.Cost,
			&recipe.TimeToPrepare,
			&recipe.Healthy,
			pq.Array(&recipe.ImageURLs),
		)
		if err != nil {
			return nil, err
		}
		output = append(output, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func (r *RecipesPostgres) GetPublicRecipes() ([]models.Recipe, error) {
	db := r.db

	var output []models.Recipe
	getAllExistRecipeQuery := fmt.Sprintf(`
		SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy", rt."imageURLs"
		FROM %s AS rt
		WHERE rt.public=true`, recipeTable)

	rows, err := db.Queryx(getAllExistRecipeQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(
			&recipe.Id,
			&recipe.Title,
			&recipe.Description,
			&recipe.IsPublic,
			&recipe.Cost,
			&recipe.TimeToPrepare,
			&recipe.Healthy,
			pq.Array(&recipe.ImageURLs),
		)
		if err != nil {
			return nil, err
		}
		output = append(output, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func (r *RecipesPostgres) UpdateRecipe(userId, id int, input models.UpdateRecipeInput) error {
	db := r.db

	// Check if the recipe belongs to the user
	recipeExistsQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE id = $1 AND user_id = $2`, recipeTable)
	var count int
	err := db.Get(&count, recipeExistsQuery, id, userId)
	if err != nil {
		return err
	}
	if count == 0 {
		// Recipe does not exist or does not belong to the user
		return errors.New("recipe not found or not owned by the user")
	}

	// Rest of the code for updating the recipe
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

	if input.ImageURLs != nil {
		setValues = append(setValues, fmt.Sprintf(`"imageURLs"=$%d`, argId))
		args = append(args, pq.Array(*input.ImageURLs))
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE "id"=%d AND "user_id"=%d`, recipeTable, setQuery, id, userId)
	fmt.Println(query)
	args = append(args)

	_, err = db.Exec(query, args...)

	return err
}

func (r *RecipesPostgres) DeleteRecipe(userId, id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE "id"=$1 and "user_id"=$2 `, recipeTable)
	_, err := r.db.Exec(query, id, userId)
	return err

}
