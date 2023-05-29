package utils

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"runtime"
	"strconv"
)

func StringArrayToIntArray(strArr []string) ([]int, error) {
	intArr := make([]int, len(strArr))

	for i, str := range strArr {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		intArr[i] = num
	}

	return intArr, nil
}

func ArrayToString(meal []models.Meal) string {
	var str string
	for _, v := range meal {
		str += v.Name + " " + v.AtTime + " " + ArrayIntToString(v.Recipes) + "; \n"
	}
	return str
}

/*
type ScheduleMealsOutput struct {
	Id       int    `json:"id,omitempty" db:"id"`
	Name     string `json:"name,omitempty" db:"name"`
	AtTime   string `json:"at_time,omitempty" db:"at_time"`
	UserId   int    `json:"user_Id,omitempty" db:"user_id"`
	RecipeID int    `json:"recipeID,omitempty" db:"recipeId"`
	MealId   int    `json:"mealId,omitempty" db:"mealId"`
}
*/

func ScheduleMealsOutputToString(meal []models.ScheduleByDateOutput) string {
	var str string
	for _, v := range meal {
		//cost to string
		stringCost := strconv.FormatFloat(v.Cost, 'f', 2, 64)
		stringTimeToPrepare := strconv.Itoa(v.TimeToPrepare)
		stringHealthy := strconv.Itoa(v.Healthy)
		str += fmt.Sprintf("Time: %s \nName: %s \nTitle: %s \nDescription: %s \nCost: %s \nCooking time: %s \nHealthy: %s \n\n", v.AtTime, v.Name, v.Title, v.Description, stringCost, stringTimeToPrepare, stringHealthy)
	}
	return str
}
func ArrayIntToString(arr []int) string {
	var str string
	for _, v := range arr {
		str += strconv.Itoa(v) + ","
	}
	return str
}

func GetErrorLocation() (string, int) {
	_, file, line, _ := runtime.Caller(2)
	return file, line
}
