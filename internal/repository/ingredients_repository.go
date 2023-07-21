package repository

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func (r *RecipesPostgres) AddOrUpdateIngredientTx(tx *sqlx.Tx, ingredient models.Ingredient) (int, error) {
	addIngredientQuery := fmt.Sprintf(`INSERT INTO %s ("name", "price", "unit", "possible_units", "protein", "fat", "carbs", "aisle", "image", "categoryPath", "consistency", "external_id", "unitLong", "unitShort") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) 
		RETURNING id`,
		ingredientTable)

	row := tx.QueryRow(addIngredientQuery, ingredient.Name, ingredient.Price, ingredient.Unit, pq.Array(ingredient.PossibleUnits), ingredient.Protein, ingredient.Fat, ingredient.Carbs, ingredient.Aisle, ingredient.Image, pq.Array(ingredient.CategoryPath), ingredient.Consistency, ingredient.ExternalID, ingredient.UnitLong, ingredient.UnitShort)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// todo user's amount and unit
func (r *RecipesPostgres) AddRecipeIngredientTx(tx *sqlx.Tx, recipeId int, ingredient models.Ingredient, amount float64) error {
	addRecipeIngredientQuery := fmt.Sprintf(`
		INSERT INTO %s (recipe_id, ingredient_id, amount, unit, price)
		VALUES ($1, $2, $3, $4, $5)`,
		recipeIngredientsTable)

	_, err := tx.Exec(addRecipeIngredientQuery, recipeId, ingredient.ID, amount, ingredient.Unit, ingredient.Price)
	if err != nil {
		return err
	}

	return nil
}

func (r *RecipesPostgres) RemoveAllRecipeIngredientsTx(tx *sqlx.Tx, recipeId int) error {
	removeRecipeIngredientsQuery := fmt.Sprintf(`
		DELETE FROM %s WHERE recipe_id = $1`,
		recipeIngredientsTable)

	_, err := tx.Exec(removeRecipeIngredientsQuery, recipeId)
	if err != nil {
		return err
	}

	return nil
}
