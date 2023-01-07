package models

type ScheduleOutput struct {
	BreakfastTitle       string `json:"BreakfastTitle" db:"BreakfastTitle"`
	BreakfastDescription string `json:"BreakfastDescription" db:"BreakfastDescription"`
	LunchTitle           string `json:"LunchTitle" db:"LunchTitle"`
	LunchDescription     string `json:"LunchDescription" db:"LunchDescription"`
	DinnerTitle          string `json:"DinnerTitle" db:"DinnerTitle"`
	DinnerDescription    string `json:"DinnerDescription" db:"DinnerDescription"`
}
