package models

type Recipe struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRecipeInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
