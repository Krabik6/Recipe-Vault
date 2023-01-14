package models

type ScheduleOutput struct {
	BreakfastTitle       string `json:"breakfastTitle" db:"BreakfastTitle"`
	BreakfastDescription string `json:"breakfastDescription" db:"BreakfastDescription"`
	LunchTitle           string `json:"lunchTitle" db:"LunchTitle"`
	LunchDescription     string `json:"lunchDescription" db:"LunchDescription"`
	DinnerTitle          string `json:"dinnerTitle" db:"DinnerTitle"`
	DinnerDescription    string `json:"dinnerDescription" db:"DinnerDescription"`
}
