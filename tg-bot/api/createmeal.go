package api

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"io"
	"net/http"
	"strings"
)

// CreateMealPlan function takes the client, meal plan, and token as arguments and returns an error.
// It is used to create a meal plan for a specific date using the API.
func (a *Api) CreateMealPlan(mealPlan models.Meal, token string) error {
	requestBody, err := json.Marshal(mealPlan)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/api/schedule/meal", strings.NewReader(string(requestBody)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("status code %d. \nResponse: %v", resp.StatusCode, string(body))
	}

	return nil
}
