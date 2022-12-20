package models

type Recipe struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateRecipeInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
