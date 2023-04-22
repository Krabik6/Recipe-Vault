package models

type ScheduleOutput struct {
	Id       int    `json:"id,omitempty" db:"id"`
	Name     string `json:"name,omitempty" db:"name"`
	AtTime   string `json:"at_time,omitempty" db:"at_time"`
	UserId   int    `json:"user_Id,omitempty" db:"user_id"`
	RecipeID int    `json:"recipeID,omitempty" db:"recipeId"`
	MealId   int    `json:"mealId,omitempty" db:"mealId"`
}
