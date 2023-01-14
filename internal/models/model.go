package models

type Recipe struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic" db:"public"`
}

type UpdateRecipeInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsPublic    *bool   `json:"isPublic" db:"public"`
}

// date_of
// breakfast_id
// lunch_id
// dinner_id
// user_id
type Schedule struct {
	Id          int    `json:"id,omitempty"`
	Date        string `json:"date,omitempty"`
	BreakfastId int    `json:"breakfastId,omitempty"`
	LunchId     int    `json:"lunchId,omitempty"`
	DinnerId    int    `json:"dinnerId,omitempty"`
}

type UpdateScheduleInput struct {
	Date        *string `json:"date,omitempty"`
	BreakfastId *int    `json:"breakfastId,omitempty"`
	LunchId     *int    `json:"lunchId,omitempty"`
	DinnerId    *int    `json:"dinnerId,omitempty"`
}

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RecipeOutput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic" db:"public"`
}
