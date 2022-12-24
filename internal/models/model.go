package models

type Recipe struct {
	Id          int    `json:"Id"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type UpdateRecipeInput struct {
	Title       *string `json:"Title"`
	Description *string `json:"Description"`
}

type Schedule struct {
	Id        int    `json:"Id"`
	Date      string `json:"date"`
	Breakfast int    `json:"breakfast"`
	Lunch     int    `json:"lunch"`
	Dinner    int    `json:"dinner"`
}

type UpdateScheduleInput struct {
	Date      *string `json:"date"`
	Breakfast *int    `json:"breakfast"`
	Lunch     *int    `json:"lunch"`
	Dinner    *int    `json:"dinner"`
}
