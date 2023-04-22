package utils

import (
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
