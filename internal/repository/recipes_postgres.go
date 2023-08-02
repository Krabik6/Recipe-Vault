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
func (r *RecipesPostgres) CreateRecipeTx(tx *sqlx.Tx, userId int, recipe models.Recipe) (int, error) {
	addRecipeQuery := fmt.Sprintf(`INSERT INTO %s 
		("title", "description", "user_id", "public", "cost", "timeToPrepare", "healthy", "imageURLs")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id`,
		recipeTable)

	row := tx.QueryRow(addRecipeQuery, recipe.Title, recipe.Description, userId, recipe.IsPublic, recipe.Cost, recipe.TimeToPrepare, recipe.Healthy, pq.Array(recipe.ImageURLs))
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RecipesPostgres) CreateRecipe(userId int, recipe models.Recipe, ingredients []models.Ingredient) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}

	recipeID, err := r.CreateRecipeTx(tx, userId, recipe)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, ingredient := range ingredients {
		ingredientID, err := r.AddOrUpdateIngredientTx(tx, ingredient)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		ingredient.ID = ingredientID

		err = r.AddRecipeIngredientTx(tx, recipeID, ingredient, ingredient.Amount)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return recipeID, nil
}

func (r *RecipesPostgres) GetRecipeById(userId, id int) (models.RecipeOutput, error) {
	db := r.db

	output := models.RecipeOutput{}
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

	// Получение ингредиентов
	ingredients, err := r.GetIngredientsByRecipeId(id)
	if err != nil {
		return output, err
	}

	output.IngredientOutputs = ingredients

	return output, nil
}
func (r *RecipesPostgres) GetIngredientsByRecipeId(recipeId int) ([]models.IngredientOutput, error) {
	db := r.db

	var ingredients []models.IngredientOutput
	getIngredientsQuery := fmt.Sprintf(`
		SELECT i."id", i."name", i."price", i."unit", i."unitShort", i."unitLong", i."protein", i."fat", i."carbs", i."aisle", i."image", i."consistency", i."external_id", ri."amount", i."possible_units", i."categoryPath"
		FROM %s AS i
		INNER JOIN %s AS ri ON i."id" = ri."ingredient_id"
		WHERE ri."recipe_id" = $1`, ingredientsTable, recipeIngredientsTable)

	rows, err := db.Queryx(getIngredientsQuery, recipeId)
	if err != nil {
		return ingredients, err
	}
	defer rows.Close()

	for rows.Next() {
		var ingredient models.IngredientOutput
		err = rows.Scan(
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.Price,
			&ingredient.Unit,
			&ingredient.UnitShort,
			&ingredient.UnitLong,
			&ingredient.Protein,
			&ingredient.Fat,
			&ingredient.Carbs,
			&ingredient.Aisle,
			&ingredient.Image,
			&ingredient.Consistency,
			&ingredient.ExternalID,
			&ingredient.Amount,
			pq.Array(&ingredient.PossibleUnits),
			pq.Array(&ingredient.CategoryPath),
		)
		if err != nil {
			return ingredients, err
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
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

func (r *RecipesPostgres) GetAllRecipes(userId int) ([]models.RecipeOutput, error) {
	db := r.db

	var output []models.RecipeOutput
	getAllRecipeQuery := fmt.Sprintf(`
		SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy", rt."imageURLs", 
		i.*, ri."amount"
		FROM %s AS rt
		INNER JOIN %s AS ri ON rt."id" = ri."recipe_id"
		INNER JOIN %s AS i ON ri."ingredient_id" = i."id"
		WHERE rt."user_id" = $1`, recipeTable, recipeIngredientsTable, ingredientsTable)

	rows, err := db.Queryx(getAllRecipeQuery, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipeMap := make(map[int]*models.RecipeOutput)

	for rows.Next() {
		var recipe models.RecipeOutput
		var ingredient models.IngredientOutput
		err := rows.Scan(
			&recipe.Id,
			&recipe.Title,
			&recipe.Description,
			&recipe.IsPublic,
			&recipe.Cost,
			&recipe.TimeToPrepare,
			&recipe.Healthy,
			pq.Array(&recipe.ImageURLs),
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.Price,
			&ingredient.Unit,
			pq.Array(&ingredient.PossibleUnits),
			&ingredient.UnitShort,
			&ingredient.UnitLong,
			&ingredient.Protein,
			&ingredient.Fat,
			&ingredient.Carbs,
			&ingredient.Aisle,
			&ingredient.Image,
			pq.Array(&ingredient.CategoryPath),
			&ingredient.Consistency,
			&ingredient.ExternalID,
			&ingredient.Amount,
		)
		if err != nil {
			return nil, err
		}

		if existingRecipe, exists := recipeMap[recipe.Id]; exists {
			existingRecipe.IngredientOutputs = append(existingRecipe.IngredientOutputs, ingredient)
		} else {
			recipe.IngredientOutputs = append(recipe.IngredientOutputs, ingredient)
			recipeMap[recipe.Id] = &recipe
		}
	}

	for _, recipe := range recipeMap {
		output = append(output, *recipe)
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

func (r *RecipesPostgres) UpdateRecipe(userId, recipeId int, input models.UpdateRecipe) error {
	db := r.db

	// Start a new transaction
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Check if the recipe belongs to the user
	recipeExistsQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE id = $1 AND user_id = $2`, recipeTable)
	var count int
	err = tx.Get(&count, recipeExistsQuery, recipeId, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if count == 0 {
		// Recipe does not exist or does not belong to the user
		tx.Rollback()
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

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE "id"=%d AND "user_id"=%d`, recipeTable, setQuery, recipeId, userId)
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Remove all existing ingredient relationships
	err = r.RemoveAllRecipeIngredientsTx(tx, recipeId)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, ingredient := range input.Ingredients {
		ingredientID, err := r.AddOrUpdateIngredientTx(tx, ingredient)
		if err != nil {
			tx.Rollback()
			return err
		}
		ingredient.ID = ingredientID

		err = r.AddRecipeIngredientTx(tx, recipeId, ingredient, ingredient.Amount)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *RecipesPostgres) DeleteRecipe(userId, id int) error {
	// Start a new transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// Delete all ingredient relationships for the recipe
	err = r.RemoveAllRecipeIngredientsTx(tx, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the recipe
	query := fmt.Sprintf(`DELETE FROM %s WHERE "id"=$1 and "user_id"=$2 `, recipeTable)
	_, err = tx.Exec(query, id, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
