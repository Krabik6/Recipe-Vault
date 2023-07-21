package models

type RecipeInput struct {
	Id               int               `json:"id,omitempty" db:"id"`
	Title            string            `json:"title,omitempty" db:"title"`
	Description      string            `json:"description,omitempty" db:"description"`
	IsPublic         bool              `json:"public,omitempty" db:"public"`
	Cost             float64           `json:"cost,omitempty" db:"cost"`
	TimeToPrepare    int64             `json:"timeToPrepare,omitempty" db:"timeToPrepare"`
	Healthy          int64             `json:"healthy,omitempty" db:"healthy"`
	ImageURLs        []string          `json:"imageURLs,omitempty" db:"imageURLs"`
	IngredientInputs []IngredientInput `json:"ingredients"`
}

type Recipe struct {
	Id               int               `json:"id,omitempty" db:"id"`
	Title            string            `json:"title,omitempty" db:"title"`
	Description      string            `json:"description,omitempty" db:"description"`
	IsPublic         bool              `json:"public,omitempty" db:"public"`
	Cost             float64           `json:"cost,omitempty" db:"cost"`
	TimeToPrepare    int64             `json:"timeToPrepare,omitempty" db:"timeToPrepare"`
	Healthy          int64             `json:"healthy,omitempty" db:"healthy"`
	ImageURLs        []string          `json:"imageURLs,omitempty" db:"imageURLs"`
	Ingredients      []Ingredient      `json:"ingredients"`
	IngredientInputs []IngredientInput `json:"ingredient_inputs,omitempty" db:"ingredient_inputs"`
}

type RecipeOutput struct {
	Id                int                `json:"id,omitempty" db:"id"`
	Title             string             `json:"title,omitempty" db:"title"`
	Description       string             `json:"description,omitempty" db:"description"`
	IsPublic          bool               `json:"public,omitempty" db:"public"`
	Cost              float64            `json:"cost,omitempty" db:"cost"`
	TimeToPrepare     int64              `json:"timeToPrepare,omitempty" db:"timeToPrepare"`
	Healthy           int64              `json:"healthy,omitempty" db:"healthy"`
	ImageURLs         []string           `json:"imageURLs,omitempty" db:"imageURLs"`
	IngredientOutputs []IngredientOutput `json:"ingredients"`
}

type UpdateRecipeInput struct {
	Id               *int               `json:"id,omitempty" db:"id"`
	Title            *string            `json:"title,omitempty" db:"title"`
	Description      *string            `json:"description,omitempty" db:"description"`
	IsPublic         *bool              `json:"public,omitempty" db:"public"`
	Cost             *float64           `json:"cost,omitempty" db:"cost"`
	TimeToPrepare    *int               `json:"timeToPrepare,omitempty" db:"timeToPrepare"`
	Healthy          *int               `json:"healthy,omitempty" db:"healthy"`
	ImageURLs        *[]string          `json:"imageURLs,omitempty" db:"imageURLs"`
	IngredientInputs *[]IngredientInput `json:"ingredients"`
}

type UpdateRecipe struct {
	Id               *int               `json:"id,omitempty" db:"id"`
	Title            *string            `json:"title,omitempty" db:"title"`
	Description      *string            `json:"description,omitempty" db:"description"`
	IsPublic         *bool              `json:"public,omitempty" db:"public"`
	Cost             *float64           `json:"cost,omitempty" db:"cost"`
	TimeToPrepare    *int               `json:"timeToPrepare,omitempty" db:"timeToPrepare"`
	Healthy          *int               `json:"healthy,omitempty" db:"healthy"`
	ImageURLs        *[]string          `json:"imageURLs,omitempty" db:"imageURLs"`
	IngredientInputs *[]IngredientInput `json:"ingredientInputs"`
	Ingredients      []Ingredient       `json:"ingredients"`
}
