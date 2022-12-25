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
	Id          int    `json:"Id,omitempty"`
	Date        string `json:"Date,omitempty"`
	BreakfastId int    `json:"BreakfastId,omitempty"`
	LunchId     int    `json:"LunchId,omitempty"`
	DinnerId    int    `json:"DinnerId,omitempty"`
}

type UpdateScheduleInput struct {
	Date        *string `json:"Date,omitempty"`
	BreakfastId *int    `json:"BreakfastId,omitempty"`
	LunchId     *int    `json:"LunchId,omitempty"`
	DinnerId    *int    `json:"DinnerId,omitempty"`
}
