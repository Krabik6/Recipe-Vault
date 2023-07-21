package spoonacular

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (api *SpoonacularAPI) GetIngredientInfo(id int, options *models.IngredientInfoOptions) (*models.IngredientAPIResponse, error) {
	queryParams := url.Values{}
	queryParams.Set("amount", strconv.Itoa(options.Amount))
	queryParams.Set("unit", options.Unit)

	url := fmt.Sprintf("%s/food/ingredients/%d/information?%s", api.BaseURL, id, queryParams.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-RapidAPI-Host", "spoonacular-recipe-food-nutrition-v1.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", api.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Spoonacular API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from Spoonacular API: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var ingredientResult models.IngredientAPIResponse
	err = json.Unmarshal(body, &ingredientResult)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return &ingredientResult, nil
}
