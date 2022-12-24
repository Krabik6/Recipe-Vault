package models

type ScheduleOutput struct {
	Id                   int    `json:"Id"  `
	DateOf               string `json:"DateOf" db:"DateOf" `
	BreakfastId          int    `json:"BreakfastId" db:"BreakfastId"`
	LunchId              int    `json:"LunchId" db:"LunchId"`
	DinnerId             int    `json:"DinnerId" db:"DinnerId"`
	BreakfastTitle       string `json:"BreakfastTitle" db:"BreakfastTitle"`
	BreakfastDescription string `json:"BreakfastDescription" db:"BreakfastDescription"`
	LunchTitle           string `json:"LunchTitle" db:"LunchTitle"`
	LunchDescription     string `json:"LunchDescription" db:"LunchDescription"`
	DinnerTitle          string `json:"DinnerTitle" db:"DinnerTitle"`
	DinnerDescription    string `json:"DinnerDescription" db:"DinnerDescription"`
}
