package api

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	"io"
	"net/http"
	"strings"
)

// CreateRecipe function takes the bot, update, client, and recipe as arguments, and returns an error.
// It is used to create a recipe using the API.
// CreateRecipe function takes the client, recipe, and token as arguments and returns an error.
// It is used to create a recipe using the API.
func CreateRecipe(client *http.Client, recipe model.CreateRecipeInput, token string) (int64, error) {
	requestBody, err := json.Marshal(recipe)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/api/recipes/", strings.NewReader(string(requestBody)))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("status code %d. \nResponse: %v", resp.StatusCode, string(body))
	}

	// Get the ID of the created recipe
	var createRecipeResponse model.CreateRecipeResponse
	err = json.NewDecoder(resp.Body).Decode(&createRecipeResponse)
	if err != nil {
		return 0, err
	}

	return createRecipeResponse.Id, nil
}
