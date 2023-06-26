package api

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"io"
	"log"
	"net/http"
)

func GetMealPlans(client *http.Client, token string) ([]models.ScheduleOutput, error) {
	req, err := http.NewRequest("GET", "http://localhost:8000/api/schedule/all", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code %d. \nResponse: %v", resp.StatusCode, string(body))
	}

	var mealPlans []models.ScheduleOutput
	err = json.NewDecoder(resp.Body).Decode(&mealPlans)
	if err != nil {
		return nil, err
	}

	log.Println(mealPlans)

	return mealPlans, nil
}
