package models

type RecipeInput struct {
	Title         string  `form:"title" binding:"required"`
	Description   string  `form:"description" binding:"required"`
	IsPublic      bool    `form:"public"`
	Cost          float64 `form:"cost"`
	TimeToPrepare int64   `form:"timeToPrepare"`
	Healthy       int     `form:"healthy"`
}
