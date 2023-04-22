package model

type SignUpCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name" binding:"required"`
}

type SignInCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ScheduleResponse struct {
	Id int64 `json:"id"`
}

type CreateRecipeResponse struct {
	Id int64 `json:"id"`
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

type UpdateRecipeResponse struct {
	Response string `json:"response"`
}

type RecipesFilter struct {
	CostMoreThan          *float64 `json:"costMoreThan,omitempty"`
	CostLessThan          *float64 `json:"costLessThan,omitempty"`
	TimeToPrepareMoreThan *int     `json:"timeToPrepareMoreThan,omitempty"`
	TimeToPrepareLessThan *int     `json:"timeToPrepareLessThan,omitempty"`
	HealthyMoreThan       *int     `json:"healthyMoreThan,omitempty"`
	HealthyLessThan       *int     `json:"healthyLessThan,omitempty"`
}
