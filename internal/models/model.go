package models

type Recipe struct {
	Id            int     `json:"id,omitempty" db:"id"`
	Title         string  `json:"title,omitempty" db:"title"`
	Description   string  `json:"description,omitempty" db:"description"`
	IsPublic      bool    `json:"public,omitempty" db:"public"`
	Cost          float64 `json:"cost,omitempty" db:"cost"`
	TimeToPrepare int64   `json:"timeToPrepare,omitempty" db:"timeToPrepare"`
	Healthy       int     `json:"healthy,omitempty" db:"healthy"`
}

type UpdateRecipeInput struct {
	Id            *int     `json:"id,omitempty" db:"id"`
	Title         *string  `json:"title"`
	Description   *string  `json:"description"`
	IsPublic      *bool    `json:"isPublic" db:"public"`
	Cost          *float64 `json:"cost,omitempty"`
	TimeToPrepare *int     `json:"timeToPrepare,omitempty"`
	Healthy       *int     `json:"healthy,omitempty"`
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

type Ingredient struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic" db:"public"`
}

type IngredientOutput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic" db:"public"`
}

type UpdateIngredientInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsPublic    *bool   `json:"isPublic" db:"public"`
}

//"name" varchar, --завтрак
//"at_time" timestamp not null, -- 10.04.2045, 10:00 по мск
//"user_id"

type Meal struct {
	Name    string `json:"name,omitempty"`
	AtTime  string `json:"at_time,omitempty"`
	Recipes []int  `json:"recipes"`
}

type ScheduleByDateOutput struct {
	Id            int     `json:"id,omitempty" db:"id"`
	Name          string  `json:"name,omitempty" db:"name"`
	AtTime        string  `json:"at_time,omitempty" db:"at_time"`
	Title         string  `json:"title,omitempty" db:"title"`
	Description   string  `json:"description,omitempty" db:"description"`
	Public        bool    `json:"public,omitempty" db:"public"`
	Cost          float64 `json:"cost,omitempty" db:"cost"`
	TimeToPrepare int     `json:"timeToPrepare,omitempty" db:"timeToPrepare"`
	Healthy       int     `json:"healthy,omitempty" db:"healthy"`
}

//type Filter struct {
//	Parameter   string
//	Restriction string
//	Value       int64
//}

type RecipesFilter struct {
	CostMoreThan          *float64 `json:"costMoreThan,omitempty"`
	CostLessThan          *float64 `json:"costLessThan,omitempty"`
	TimeToPrepareMoreThan *int     `json:"timeToPrepareMoreThan,omitempty"`
	TimeToPrepareLessThan *int     `json:"timeToPrepareLessThan,omitempty"`
	HealthyMoreThan       *int     `json:"healthyMoreThan,omitempty"`
	HealthyLessThan       *int     `json:"healthyLessThan,omitempty"`
}
